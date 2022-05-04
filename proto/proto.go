package proto

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type Score interface {
	Score(*Args) Status
}

type Args struct {
	Node Node
	Task Task
}

type Node struct {
	Name          string
	Unschedulable bool
}

type Task struct {
	NodeName               string
	ToleratesUnschedulable bool
}

type Status struct {
	Error string
}

type ScoreRPC struct {
	client *rpc.Client
}

func (n *ScoreRPC) Score(args *Args) Status {
	var resp Status
	if err := n.client.Call("Plugin.Score", args, &resp); err != nil {
		panic(err)
	}
	return resp
}

type ScoreRPCServer struct {
	Impl Score
}

func (n *ScoreRPCServer) Score(args *Args, resp *Status) error {
	*resp = n.Impl.Score(args)
	return nil
}

type ScorePlugin struct {
	Impl Score
}

func (n *ScorePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ScoreRPCServer{Impl: n.Impl}, nil
}

func (ScorePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ScoreRPC{client: c}, nil
}
