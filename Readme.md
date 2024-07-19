# Temporama

Temporama is a distributed cache system using in-memory storage and the Redis Serialization Protocol (RESP). Since Temporama uses RESP as its communication protocol, you can use any Redis client driver with any programming language to communicate with Temporama.

Disclaimer: Temporama is still a work in progress. While it is functional, it currently supports only a limited set of commands and features. Expect frequent updates and potential changes as development continues.

## Current Supported Command
While Temporama is still a work in progress, it currently supports only a limited set of commands:
- GET
- SET
- DEL
- HMGET
- HMSET
- HGET
- HSET
- HGETALL
- HELLO (for handshake)
- PING

Similar to Redis, Temporama also supports command pipelining, allowing multiple commands to be sent in a single request by the client.

## Start Temporama
To start Temporama, you need to build it first using `make`:
```
make build
```

Then, start Temporama:
```
./bin/temporama
```

By default, it uses port 6379. Use the command below to override the port number:
```
PORT=8029 ./bin/temporama
```

## Connect using `redis-cli`

Use the redis-cli command to connect to Temporama. If Temporama is running on the default Redis port (6379) on localhost, you can connect with:
```
redis-cli -h localhost -p 6379
```

If Temporama is running on a different host or port, replace localhost and 6379 with the appropriate values:
```
redis-cli -h <temporama-host> -p <temporama-port>
```

Once connected, you can start executing Redis commands to interact with Temporama. For example:
```
SET mykey "Hello, Temporama"
GET mykey
```

## Connect using `redigo`
Use the following example code to connect to Temporama and execute some Redis commands.

```go
package main

import (
    "fmt"
    "github.com/gomodule/redigo/redis"
)

func main() {
    // Connect to Temporama
    conn, err := redis.Dial("tcp", "localhost:6379")
    if err != nil {
        fmt.Println("Error connecting to Temporama:", err)
        return
    }
    defer conn.Close()

    // Set a key
    _, err = conn.Do("SET", "mykey", "Hello, Temporama")
    if err != nil {
        fmt.Println("Error setting key:", err)
        return
    }

    // Get the key
    value, err := redis.String(conn.Do("GET", "mykey"))
    if err != nil {
        fmt.Println("Error getting key:", err)
        return
    }

    fmt.Println("mykey:", value)
}
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)