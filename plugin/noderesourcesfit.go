package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-score/proto"
)

const (
	ErrReasonResourcesFit = "NodeResourcesFit: node(s) didn't fit the resources"
)

type NodeResourcesFit struct{}

func (n *NodeResourcesFit) Score(args *proto.Args) proto.Status {
	var status proto.Status

	// TODO
	status.Error = ErrReasonResourcesFit

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
		"NodeResourcesFit": &proto.ScorePlugin{Impl: &NodeResourcesFit{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
