package bmreqs

import (
	"context"
)

// Requirements types
const (
	// ObjectSet is the set-alike requirements
	ObjectSet = uint8(0) + iota
	ObjectMax
)

// Requirements operations
const (
	OpAdd = uint8(0) + iota
	OpGet
	OpDump
)

type ReqRequest struct {
	Node  string
	Name  string
	T     uint8
	Op    uint8
	Value string
}

type ReqResponse struct {
	Value string
	Error error
}

// ReqRoot is the entry point data structure for all the bmreqs operations
type ReqRoot struct {
	ask    chan ReqRequest
	answer chan ReqResponse
	ctx    context.Context
	cancel context.CancelFunc
	obj    *bmReqObj
}

// NewReqRoot creates a new *ReqGroup object
func NewReqRoot() *ReqRoot {
	rg := new(ReqRoot)
	rg.ask = make(chan ReqRequest)
	rg.answer = make(chan ReqResponse)
	obj := new(bmReqObj)
	obj.init()
	rg.obj = obj
	ctx, cancel := context.WithCancel(context.Background())
	rg.ctx = ctx
	rg.cancel = cancel
	go rg.run()
	return rg
}

func (rg *ReqRoot) Requirement(req ReqRequest) ReqResponse {
	rg.ask <- req
	return <-rg.answer

}

func (rg *ReqRoot) Close() {
	rg.cancel()
}
