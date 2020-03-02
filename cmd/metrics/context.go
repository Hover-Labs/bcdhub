package main

import (
	"fmt"
	"time"

	"github.com/aopoltorzhicky/bcdhub/internal/elastic"
	"github.com/aopoltorzhicky/bcdhub/internal/mq"
	"github.com/aopoltorzhicky/bcdhub/internal/noderpc"
)

// Context -
type Context struct {
	ES  *elastic.Elastic
	RPC map[string]*noderpc.NodeRPC
	MQ  *mq.MQ
}

func newContext(cfg config) (*Context, error) {
	es, err := elastic.New([]string{cfg.Search.URI})
	if err != nil {
		return nil, err
	}
	RPCs := createRPCs(cfg)
	messageQueue, err := mq.New(cfg.Mq.URI, cfg.Mq.Queues)
	if err != nil {
		return nil, err
	}
	return &Context{
		ES:  es,
		RPC: RPCs,
		MQ:  messageQueue,
	}, nil
}

// Close -
func (ctx *Context) Close() {
	ctx.MQ.Close()
}

func createRPCs(cfg config) map[string]*noderpc.NodeRPC {
	rpc := make(map[string]*noderpc.NodeRPC)
	for i := range cfg.NodeRPC {
		nodeCfg := cfg.NodeRPC[i]
		rpc[nodeCfg.Network] = noderpc.NewNodeRPC(nodeCfg.Host)
		rpc[nodeCfg.Network].SetTimeout(time.Second * 30)
	}
	return rpc
}

// GetRPC -
func (ctx *Context) GetRPC(network string) (*noderpc.NodeRPC, error) {
	if rpc, ok := ctx.RPC[network]; ok {
		return rpc, nil
	}
	return nil, fmt.Errorf("Unknown rpc network %s", network)
}