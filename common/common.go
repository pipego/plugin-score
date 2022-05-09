package common

import (
	"net/rpc"

	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/plugin"
)

type Score interface {
	Score(*plugin.Args) Result
}

type Result struct {
	Score int64
}

type ScoreRPC struct {
	client *rpc.Client
}

func (n *ScoreRPC) Score(args *plugin.Args) Result {
	var resp Result
	if err := n.client.Call("Plugin.Score", args, &resp); err != nil {
		panic(err)
	}
	return resp
}

type ScoreRPCServer struct {
	Impl Score
}

func (n *ScoreRPCServer) Score(args *plugin.Args, resp *Result) error {
	*resp = n.Impl.Score(args)
	return nil
}

type ScorePlugin struct {
	Impl Score
}

func (n *ScorePlugin) Server(*gop.MuxBroker) (interface{}, error) {
	return &ScoreRPCServer{Impl: n.Impl}, nil
}

func (ScorePlugin) Client(b *gop.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ScoreRPC{client: c}, nil
}
