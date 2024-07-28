package memstore

import (
	"fmt"
	"hash/crc32"
	"sync"
)

const (
	MaxShardBlock uint32 = 20
)

const (
	// IEEE is by far and away the most common CRC-32 polynomial.
	// Used by ethernet (IEEE 802.3), v.42, fddi, gzip, zip, png, ...
	IEEE = 0xedb88320

	// Castagnoli's polynomial, used in iSCSI.
	// Has better error detection characteristics than IEEE.
	// https://dx.doi.org/10.1109/26.231911
	Castagnoli = 0x82f63b78

	// Koopman's polynomial.
	// Also has better error detection characteristics than IEEE.
	// https://dx.doi.org/10.1109/DSN.2002.1028931
	Koopman = 0xeb31d82e
)

var (
	shard  *Shard
	doOnce sync.Once
)

func init() {
	doOnce.Do(func() {
		shard = new(Shard)
		shard.init()
	})
}

func Get(valueType ValueType, key string, args ...string) (any, error) {
	return shard.Get(valueType, key, args...)
}

func Set(valueType ValueType, key string, args ...string) (int, error) {
	return shard.Set(valueType, key, args...)
}

func Delete(key string) string {
	shard.Delete(key)
	return "1"
}

var (
	NumNodes uint32 = 10
)

type Range struct {
	start uint32
	end   uint32
}

func (r *Range) InRange(num uint32) bool {
	return num >= r.start && num <= r.end
}

func (r *Range) Length() uint32 {
	return (r.end + 1) - r.start
}

type Shard struct {
	doOnce sync.Once
	nodes  map[uint32]ShardNode
}

// useful zero value,
// see: https://dave.cheney.net/2013/01/19/what-is-the-zero-value-and-why-is-it-useful
func (s *Shard) init() {
	s.doOnce.Do(func() {

		remain := MaxShardBlock % NumNodes
		fairNumBlocks := (MaxShardBlock - remain) / NumNodes
		rRemain := remain % 2
		hRemain := (remain - rRemain) / 2

		s.nodes = make(map[uint32]ShardNode, NumNodes)

		startRange := uint32(0)
		for i := uint32(0); i < NumNodes; i++ {
			var extra, rests uint32

			if hRemain > 0 && rests == 0 {
				extra = 1
				hRemain--
			} else if i >= NumNodes-1 {
				rests = ((remain - rRemain) / 2) + rRemain
				extra = rests
			}

			end := s.assignNode(i, startRange, fairNumBlocks, extra)
			startRange = end
			startRange++
		}
	})
}

// for debug only
func (s *Shard) printInfo() {
	for i := 0; i < len(s.nodes); i++ {
		n := s.nodes[uint32(i)].(localShard)
		fmt.Printf("node: %d, start: %d, end: %d\n", i, n.blockRange.start, n.blockRange.end)
		for x := 0; x < len(n.blocks); x++ {
			fmt.Printf("\t\tblock: %d, storage: %+v\n", x, n.blocks[uint32(x)].storage)
		}
	}

	fmt.Println("------------------------------------------------------")
}

func (s *Shard) getNode(blockNum uint32) ShardNode {
	for _, node := range s.nodes {
		if node.InRange(blockNum) {
			return node
		}
	}

	return nil
}

func (s *Shard) assignNode(numNode, startRange, fairNumBlocks, extra uint32) uint32 {
	var end uint32

	if extra > 0 {
		end = extra
	}

	end += startRange + fairNumBlocks - 1

	if startRange+extra >= MaxShardBlock {
		end = startRange + extra
	}

	s.nodes[numNode] = newLocalShard(Range{start: startRange, end: end})

	return end
}

func (s *Shard) Get(valueType ValueType, key string, args ...string) (any, error) {
	crcTable := crc32.MakeTable(IEEE)
	hashKey := crc32.Checksum([]byte(key), crcTable)
	blockNum := hashKey % MaxShardBlock

	node := s.getNode(blockNum)
	if node == nil {
		return nil, ErrNilEntries
	}

	return node.Get(valueType, blockNum, hashKey, key, args...)
}

func (s *Shard) Set(valueType ValueType, key string, args ...string) (int, error) {
	crcTable := crc32.MakeTable(IEEE)
	hashKey := crc32.Checksum([]byte(key), crcTable)
	blockNum := hashKey % MaxShardBlock

	node := s.getNode(blockNum)
	if node == nil {
		return 0, ErrNilEntries
	}

	n, e := node.Set(valueType, blockNum, hashKey, key, args...)
	// s.printInfo()
	return n, e
}

func (s *Shard) Delete(key string) {
	crcTable := crc32.MakeTable(IEEE)
	hashKey := crc32.Checksum([]byte(key), crcTable)
	blockNum := hashKey % MaxShardBlock

	node := s.getNode(blockNum)
	if node == nil {
		return
	}

	node.Delete(blockNum, hashKey, key)
	// s.printInfo()
}

type ShardNode interface {
	InRange(blockNum uint32) bool
	Get(valueType ValueType, blockNum, hashKey uint32, key string, args ...string) (any, error)
	Set(valueType ValueType, blockNum, hashKey uint32, key string, args ...string) (int, error)
	Delete(blockNum, hashKey uint32, key string)
}

type localShard struct {
	blockRange Range
	blocks     map[uint32]shardBlock
}

func newLocalShard(blockRange Range) localShard {
	blocks := make(map[uint32]shardBlock, blockRange.Length())

	for i := uint32(0); i < blockRange.Length(); i++ {
		blocks[i] = shardBlock{
			storage: &Storage{
				entries: make(map[uint32]EntryNode),
			},
		}
	}

	return localShard{
		blockRange: blockRange,
		blocks:     blocks,
	}
}

func (l localShard) InRange(blockNum uint32) bool {
	return l.blockRange.InRange(blockNum)
}

func (l localShard) Get(valueType ValueType, blockNum, hashKey uint32, key string, args ...string) (any, error) {
	return l.blocks[blockNum-l.blockRange.start].Get(valueType, hashKey, key, args...)
}

func (l localShard) Set(valueType ValueType, blockNum, hashKey uint32, key string, args ...string) (int, error) {
	n, e := l.blocks[blockNum-l.blockRange.start].Set(valueType, hashKey, key, args...)
	return n, e
}

func (l localShard) Delete(blockNum, hashKey uint32, key string) {
	l.blocks[blockNum-l.blockRange.start].Delete(hashKey, key)
}

type RemoteConfig struct{}

type remoteShard struct {
	blockRange Range
	config     RemoteConfig
}

func newRemoteShard(blockRange Range, config RemoteConfig) remoteShard {
	return remoteShard{
		blockRange: blockRange,
		config:     config,
	}
}

func (r *remoteShard) InRange(blockNum uint32) bool {
	return r.blockRange.InRange(blockNum)
}

func (r *remoteShard) Get(valueType ValueType, hashKey uint32, key string, args ...string) (any, error) {
	return nil, nil
}

type shardBlock struct {
	storage *Storage
}

func (sb shardBlock) Get(valueType ValueType, hashKey uint32, key string, args ...string) (any, error) {
	return sb.storage.Get(valueType, hashKey, key, args...)
}

func (sb shardBlock) Set(valueType ValueType, hashKey uint32, key string, args ...string) (int, error) {
	return sb.storage.Set(valueType, hashKey, key, args...)
}

func (sb shardBlock) Delete(hashKey uint32, key string) {
	sb.storage.Delete(hashKey, key)
}
