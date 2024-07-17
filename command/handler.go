package command

import (
	"fmt"
	"os"

	"github.com/raspiantoro/temporama/info"
	"github.com/raspiantoro/temporama/memstore"
	"github.com/raspiantoro/temporama/resp"
)

func Registers() resp.CommandHandler {
	mux := resp.NewCommandMux()

	mux.HandleFunc("ping", Ping)
	mux.HandleFunc("hello", Hello)
	mux.HandleFunc("get", Get)
	mux.HandleFunc("set", Set)
	mux.HandleFunc("del", Delete)
	mux.HandleFunc("hmget", HmGet)
	mux.HandleFunc("hmset", HmSet)

	return mux
}

func Ping(cmd resp.Command) resp.ValueNode {
	if cmd.Key() != "" || len(cmd.Args()) > 0 {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue("ERROR: wrong arguments number for 'ping' command"),
		)
	}

	return resp.NewValueNode(
		resp.ValueNodeTypeSimpleString,
		resp.WithValue("PONG"),
	)
}

func Hello(cmd resp.Command) resp.ValueNode {
	if len(cmd.Args()) > 0 {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue("ERROR: wrong arguments number for 'hello' command"),
		)
	}

	key := cmd.Key()

	if key != "" && key != "2" && key != "3" {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue("ERROR: unsupport protocol version"),
		)
	}

	var response resp.ValueNode
	var proto string

	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "standalone"
	}

	role := os.Getenv("ROLE")
	if role == "" {
		role = "master"
	}

	if key == "3" {
		response = resp.NewValueNode(resp.ValueNodeTypeMaps)
		proto = "3"
	} else {
		response = resp.NewValueNode(resp.ValueNodeTypeArray)
		proto = "2"
	}

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("server"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(info.Server),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("version"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(info.Version),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("proto"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(proto),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("mode"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(mode),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("role"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(role),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("git commit"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(info.GitCommit),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("build date"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(info.BuildDate),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("go version"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(info.GoVersion),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue("os arch"),
	))

	response.Append(resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(info.OsArch),
	))

	return response
}

func Get(cmd resp.Command) resp.ValueNode {
	if cmd.Key() == "" || len(cmd.Args()) > 0 {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue("ERROR: wrong arguments number for 'get' command"),
		)
	}

	val, err := memstore.Get(memstore.ValueTypeString, cmd.Key())
	if err == memstore.ErrNilEntries {
		return resp.NewValueNode(
			resp.ValueNodeTypeBulkString,
			resp.WithValue("-1"),
		)
	}
	if err != nil {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue(fmt.Sprintf("ERROR: %s", err.Error())),
		)
	}

	return resp.NewValueNode(
		resp.ValueNodeTypeBulkString,
		resp.WithValue(val.(string)),
	)
}

func Set(cmd resp.Command) resp.ValueNode {
	if cmd.Key() == "" || len(cmd.Args()) <= 0 {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue("ERROR: wrong arguments number for 'set' command"),
		)
	}

	err := memstore.Set(memstore.ValueTypeString, cmd.Key(), cmd.Args()...)
	if err != nil {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue(fmt.Sprintf("ERROR: %s", err.Error())),
		)
	}

	return resp.NewValueNode(
		resp.ValueNodeTypeSimpleString,
		resp.WithValue("OK"),
	)
}

func Delete(cmd resp.Command) resp.ValueNode {
	response := memstore.Delete(cmd.Key())
	return resp.NewValueNode(
		resp.ValueNodeTypeIntegers,
		resp.WithValue(response),
	)
}

func HmGet(cmd resp.Command) resp.ValueNode {
	if cmd.Key() == "" || len(cmd.Args()) <= 0 {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue("ERROR: wrong arguments number for 'hmget' command"),
		)
	}

	val, err := memstore.Get(memstore.ValueTypeMap, cmd.Key(), cmd.Args()...)
	if err == memstore.ErrNilEntries {
		return resp.NewValueNode(
			resp.ValueNodeTypeBulkString,
			resp.WithValue("-1"),
		)
	}
	if err != nil {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue(fmt.Sprintf("ERROR: %s", err.Error())),
		)
	}

	vals := val.([]string)

	response := resp.NewValueNode(resp.ValueNodeTypeArray)

	for _, v := range vals {
		child := resp.NewValueNode(
			resp.ValueNodeTypeBulkString,
			resp.WithValue(v),
		)

		response.Append(child)
	}

	return response
}

func HmSet(cmd resp.Command) resp.ValueNode {
	if cmd.Key() == "" || len(cmd.Args()) <= 0 {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue("ERROR: wrong arguments number for 'hmset' command"),
		)
	}

	err := memstore.Set(memstore.ValueTypeMap, cmd.Key(), cmd.Args()...)
	if err != nil {
		return resp.NewValueNode(
			resp.ValueNodeTypeSimpleError,
			resp.WithValue(fmt.Sprintf("ERROR: %s", err.Error())),
		)
	}

	return resp.NewValueNode(
		resp.ValueNodeTypeSimpleString,
		resp.WithValue("OK"),
	)
}
