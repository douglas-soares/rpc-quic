package rpc

// rpcData contains information about the data trade between client and server
type rpcData struct {
	Function string
	Args     []interface{}
	Err      error
}
