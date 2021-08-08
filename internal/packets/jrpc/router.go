package jrpc

import (
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// AnonWrapper ...
type AnonWrapper func(id uint64, params json.RawMessage)

// Router ...
type Router struct {
	funcs map[string]AnonWrapper

	once sync.Once
	mux  sync.Mutex
}

// NewRouter returns jrpc router instance
func NewRouter() *Router {
	var o sync.Once
	return &Router{
		funcs: make(map[string]AnonWrapper),
		once:  o,
		mux:   sync.Mutex{},
	}
}

// Add method handler. Handler will be called on matching method (Request.Method)
func (r *Router) Add(method string, fn AnonWrapper) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.once.Do(func() {
		r.funcs = map[string]AnonWrapper{}
	})
	r.funcs[method] = fn
	log.Info("added handler for %s", method)
}

// Process ...
func (r *Router) Process(req *Request) {
	method := r.funcs[req.Method]
	method(req.ID, req.Params)
}

// CheckMethodExist ...
func (r *Router) CheckMethodExist(req *Request) error {
	if _, ok := r.funcs[req.Method]; ok != true {
		return fmt.Errorf("no [%s] function registered in JRPC router", req.Method)
	}
	return nil
}
