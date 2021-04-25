# RPC-QUIC

RPC-QUIC is a middleware based on Remote Procedure Calls (RPC) with the QUIC protocol and it is implemented using the Go programming language.
This project uses the [quic-go](https://github.com/lucas-clemente/quic-go) implementation.

# Installation

RPC-QUIC requires a Go version with modules support. So make sure to initialize a Go module:

```ssh
go mod init github.com/my/repo
```

And then add RPC-QUIC into your project:

```ssh
go get github.com/douglas-soares/rpc-quic
```

# Usage

The RPC-QUIC contains functions for clients and server.
For more information about available functions check out the [RPC-QUIC interfaces](src/rpc.go).

It also contains a implementation of a Naming Service in order to offer location transparency.
For more information about available functions check out the [Naming Service implementations](src/naming_service/naming.go).

## Server

```Go
tlsConfig: = &tls.Config{
	...
	NextProtos:   []string{"server"},
}

// Create a server instance
server := rpc.NewServer()

// Register the server remote function (fibonacci)
server.Register("fibonacci", fibonacci)

// Start to listen the port
server.ListenAndServe("localhost:8080", tlsConfig, &quic.Config{})
```

For more details, take a look at [this server example](example/server/server.go).

## Client

```Go
tlsConfig := &tls.Config{
	...
	NextProtos:         []string{"server"},
}

// Create a client with the server address
client := rpc.NewClient("localhost:8080", tlsConf, &quic.Config{})

// call fibonacci remote function
var resp int
err := client.Call("fibonacci", 1, &resp)

fmt.Println("Fibonacci result:", resp)
```

In order to enable the QUIC 0-RTT feature, you must add following configuration in the `tlsConf`:

```Go
tlsConfig := &tls.Config{
	...
	ClientSessionCache: tls.NewLRUClientSessionCache(100), // enables 0-RTT
}
```

For more details, take a look at [this client example](example/client/client.go).

## Naming Service

This service offers two functions:
- Bind: registers the server name with its location
- Lookup: returns the server location for the given server name

### Server

```Go
// Create the naming server
namingServer := naming.NewNamingService("localhost:4040")
// Start to listen
err := namingServer.ListenAndServe(generateTLSConfig())
```

For more details, take a look at [this naming server example](example/naming/naming.go).

### Client

```Go
// Create the naming client
namingServer := naming.NewNamingService("localhost:4040")

// Start naming client
namingServer.StartClient(tlsConfN)

// Register the server name and its location
err := namingServer.Bind("Servidor", serverAddr)

// Look up for the desired server location
serverAddr, err := n.LookUp("Servidor")
```

For more details, take a look at [this client example](example/client/client.go) and [this server example](example/server/server.go).
