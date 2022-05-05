package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-score/proto"
)

type NodeResourcesBalancedAllocation struct{}

func (n *NodeResourcesBalancedAllocation) Score(args *proto.Args) proto.Result {
	var result proto.Result

	// TODO

	return result
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
