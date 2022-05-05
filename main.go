package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-score/proto"
)

type config struct {
	args *proto.Args
	name string
	path string
}

var (
	configs = []config{
		// Plugin: NodeResourcesBalancedAllocation
		{
			args: &proto.Args{},
			name: "NodeResourcesBalancedAllocation",
			path: "./plugin/score-noderesourcesbalancedallocation",
		},
		// Plugin: NodeResourcesFit
		{
			args: &proto.Args{},
			name: "NodeResourcesFit",
			path: "./plugin/score-noderesourcesfit",
		},
	}
)

func main() {
	for _, item := range configs {
		result := helper(item.path, item.name, item.args)
		fmt.Println(result.Score)
	}
}

func helper(path, name string, args *proto.Args) proto.Result {
	config := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-score",
		MagicCookieValue: "plugin-score",
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin-score",
		Output: os.Stderr,
		Level:  hclog.Error,
	})

	plugins := map[string]plugin.Plugin{
		name: &proto.ScorePlugin{},
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		Cmd:             exec.Command(path),
		HandshakeConfig: config,
		Logger:          logger,
		Plugins:         plugins,
	})
	defer client.Kill()

	rpcClient, _ := client.Client()
	raw, _ := rpcClient.Dispense(name)
	n := raw.(proto.Score)
	status := n.Score(args)

	return status
}
