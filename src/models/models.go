package models

// Request contains information about the client request
type Request struct {
	Function string
	Args     []interface{}
}

// Response contains information about the remote server request response
type Response struct {
	Result []interface{}
	Err    error
}
