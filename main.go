package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	gop "github.com/hashicorp/go-plugin"
	"github.com/pkg/errors"

	"github.com/pipego/scheduler/common"
	"github.com/pipego/scheduler/plugin"
)

type config struct {
	args *common.Args
	name string
	path string
}

var (
	configs = []config{
		// Plugin: NodeResourcesBalancedAllocation
		{
			args: &common.Args{
				Node: common.Node{
					AllocatableResource: common.Resource{
						MilliCPU: 2000,
						Memory:   3000,
					},
					RequestedResource: common.Resource{
						MilliCPU: 256,
						Memory:   512,
					},
				},
				Task: common.Task{
					RequestedResource: common.Resource{
						MilliCPU: 1024,
						Memory:   2048,
					},
				},
			},
			name: "NodeResourcesBalancedAllocation",
			path: "./plugin/score-noderesourcesbalancedallocation",
		},
		{
			args: &common.Args{
				Node: common.Node{
					AllocatableResource: common.Resource{
						MilliCPU: 1024,
						Memory:   2048,
						Storage:  4096,
					},
					RequestedResource: common.Resource{
						MilliCPU: 512,
						Memory:   1024,
						Storage:  2048,
					},
				},
				Task: common.Task{
					RequestedResource: common.Resource{
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
			args: &common.Args{
				Node: common.Node{
					AllocatableResource: common.Resource{
						MilliCPU: 1024,
						Memory:   2048,
						Storage:  4096,
					},
					RequestedResource: common.Resource{
						MilliCPU: 512,
						Memory:   1024,
						Storage:  2048,
					},
				},
				Task: common.Task{
					RequestedResource: common.Resource{
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
			args: &common.Args{
				Node: common.Node{
					AllocatableResource: common.Resource{
						MilliCPU: 1024,
						Memory:   2048,
						Storage:  4096,
					},
					RequestedResource: common.Resource{
						MilliCPU: 512,
						Memory:   1024,
						Storage:  2048,
					},
				},
				Task: common.Task{
					RequestedResource: common.Resource{
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

	handshake = gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin",
		MagicCookieValue: "plugin",
	}

	logger = hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stderr,
		Level:  hclog.Error,
	})
)

func main() {
	for _, item := range configs {
		p, _ := filepath.Abs(item.path)
		if result, err := helper(p, item.name, item.args); err == nil {
			fmt.Printf("%s: %d\n", item.name, result.Score)
		} else {
			fmt.Printf(err.Error())
		}
	}
}

func helper(path, name string, args *common.Args) (plugin.ScoreResult, error) {
	plugins := map[string]gop.Plugin{
		name: &plugin.Score{},
	}

	client := gop.NewClient(&gop.ClientConfig{
		Cmd:             exec.Command(path),
		HandshakeConfig: handshake,
		Logger:          logger,
		Plugins:         plugins,
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		return plugin.ScoreResult{}, errors.Wrap(err, "failed to init client")
	}

	raw, err := rpcClient.Dispense(name)
	if err != nil {
		return plugin.ScoreResult{}, errors.Wrap(err, "failed to dispense instance")
	}

	n := raw.(plugin.ScoreImpl)
	status := n.Run(args)

	return status, nil
}
