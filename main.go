package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-score/common"
	"github.com/pipego/scheduler/plugin"
)

type config struct {
	args *plugin.Args
	name string
	path string
}

var (
	configs = []config{
		// Plugin: NodeResourcesBalancedAllocation
		{
			args: &plugin.Args{
				Node: plugin.Node{
					AllocatableResource: plugin.Resource{
						MilliCPU: 2000,
						Memory:   3000,
					},
					RequestedResource: plugin.Resource{
						MilliCPU: 256,
						Memory:   512,
					},
				},
				Task: plugin.Task{
					RequestedResource: plugin.Resource{
						MilliCPU: 1024,
						Memory:   2048,
					},
				},
			},
			name: "NodeResourcesBalancedAllocation",
			path: "./plugin/score-noderesourcesbalancedallocation",
		},
		{
			args: &plugin.Args{
				Node: plugin.Node{
					AllocatableResource: plugin.Resource{
						MilliCPU: 1024,
						Memory:   2048,
						Storage:  4096,
					},
					RequestedResource: plugin.Resource{
						MilliCPU: 512,
						Memory:   1024,
						Storage:  2048,
					},
				},
				Task: plugin.Task{
					RequestedResource: plugin.Resource{
						MilliCPU: 256,
						Memory:   512,
						Storage:  1024,
					},
				},
			},
			name: "NodeResourcesBalancedAllocation",
			path: "./plugin/score-noderesourcesbalancedallocation",
		},
		// Plugin: NodeResourcesFit
		{
			args: &plugin.Args{
				Node: plugin.Node{
					AllocatableResource: plugin.Resource{
						MilliCPU: 1024,
						Memory:   2048,
						Storage:  4096,
					},
					RequestedResource: plugin.Resource{
						MilliCPU: 512,
						Memory:   1024,
						Storage:  2048,
					},
				},
				Task: plugin.Task{
					RequestedResource: plugin.Resource{
						MilliCPU: 1024,
						Memory:   2048,
						Storage:  4096,
					},
				},
			},
			name: "NodeResourcesFit",
			path: "./plugin/score-noderesourcesfit",
		},
		{
			args: &plugin.Args{
				Node: plugin.Node{
					AllocatableResource: plugin.Resource{
						MilliCPU: 1024,
						Memory:   2048,
						Storage:  4096,
					},
					RequestedResource: plugin.Resource{
						MilliCPU: 512,
						Memory:   1024,
						Storage:  2048,
					},
				},
				Task: plugin.Task{
					RequestedResource: plugin.Resource{
						MilliCPU: 256,
						Memory:   512,
						Storage:  1024,
					},
				},
			},
			name: "NodeResourcesFit",
			path: "./plugin/score-noderesourcesfit",
		},
	}
)

func main() {
	for _, item := range configs {
		result := helper(item.path, item.name, item.args)
		fmt.Printf("%s: %d\n", item.name, result.Score)
	}
}

func helper(path, name string, args *plugin.Args) plugin.ScoreResult {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-score",
		MagicCookieValue: "plugin-score",
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin-score",
		Output: os.Stderr,
		Level:  hclog.Error,
	})

	plugins := map[string]gop.Plugin{
		name: &common.ScorePlugin{},
	}

	client := gop.NewClient(&gop.ClientConfig{
		Cmd:             exec.Command(path),
		HandshakeConfig: config,
		Logger:          logger,
		Plugins:         plugins,
	})
	defer client.Kill()

	rpcClient, _ := client.Client()
	raw, _ := rpcClient.Dispense(name)
	n := raw.(plugin.ScorePlugin)
	status := n.Run(args)

	return status
}
