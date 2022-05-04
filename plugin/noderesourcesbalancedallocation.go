package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-score/proto"
)

const (
	ErrReasonResourcesBalancedAllocation = "NodeResourcesBalancedAllocation: node(s) didn't match the resources balanced allocation"
)

type NodeResourcesBalancedAllocation struct{}

func (n *NodeResourcesBalancedAllocation) Score(args *proto.Args) proto.Status {
	var status proto.Status

	// TODO
	status.Error = ErrReasonResourcesBalancedAllocation

	return status
}

// nolint:typecheck
func main() {
	config := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-score",
		MagicCookieValue: "plugin-score",
	}

	pluginMap := map[string]plugin.Plugin{
		"NodeResourcesBalancedAllocation": &proto.ScorePlugin{Impl: &NodeResourcesBalancedAllocation{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
