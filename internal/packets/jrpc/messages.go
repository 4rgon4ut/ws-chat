package jrpc

import "encoding/json"

// Request encloses method name and all params
type Request struct {
	Format string          `json:"jsonrpc"`
	ID     uint64          `json:"id"`
	Method string          `json:"method"`           // method (function) name
	Params json.RawMessage `json:"params,omitempty"` // function arguments
}
