package common

import (
	"net/rpc"

	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/plugin"
)

type ScoreRPC struct {
	client *rpc.Client
}

func (n *ScoreRPC) Run(args *plugin.Args) plugin.ScoreResult {
	var resp plugin.ScoreResult
	if err := n.client.Call("Plugin.Run", args, &resp); err != nil {
		panic(err)
	}
	return resp
}

type ScoreRPCServer struct {
	Impl plugin.ScorePlugin
}

func (n *ScoreRPCServer) Run(args *plugin.Args, resp *plugin.ScoreResult) error {
	*resp = n.Impl.Run(args)
	return nil
}

type ScorePlugin struct {
	Impl plugin.ScorePlugin
}

func (n *ScorePlugin) Server(*gop.MuxBroker) (interface{}, error) {
	return &ScoreRPCServer{Impl: n.Impl}, nil
}

func (ScorePlugin) Client(b *gop.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ScoreRPC{client: c}, nil
}
