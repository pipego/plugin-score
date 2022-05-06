package proto

import (
	"math"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

const (
	// ResourceCPU CPU, in cores. (500m = .5 cores)
	ResourceCPU = "cpu"
	// ResourceMemory Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory = "memory"
	// ResourceStorage Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage = "storage"
)

const (
	// DefaultMilliCPURequest defines default milli cpu request number.
	DefaultMilliCPURequest int64 = 100 // 0.1 core
	// DefaultMemoryRequest defines default memory request size.
	DefaultMemoryRequest int64 = 200 * 1024 * 1024 // 200 MB
)

// Resources to consider when scoring.
// The default resource set includes "cpu" and "memory" with an equal weight.
const (
	// DefaultCPUWeight defines default cpu weight (allowed weights go from 1 to 100)
	DefaultCPUWeight int64 = 1
	// DefaultMemoryWeight defines default memory weight (allowed weights go from 1 to 100)
	DefaultMemoryWeight int64 = 1
	// DefaultStorageWeight defines default storage weight (allowed weights go from 1 to 100)
	DefaultStorageWeight int64 = 1
)

const (
	// MaxNodeScore is the maximum score a Score plugin is expected to return.
	MaxNodeScore int64 = 100

	// MinNodeScore is the minimum score a Score plugin is expected to return.
	MinNodeScore int64 = 0

	// MaxTotalScore is the maximum total score.
	MaxTotalScore int64 = math.MaxInt64
)

type Args struct {
	Node Node
	Task Task
}

type Node struct {
	AllocatableResource Resource `json:"allocatableResource"`
	Label               Label    `json:"label"`
	Name                string   `json:"name"`
	RequestedResource   Resource `json:"requestedResource"`
	Unschedulable       bool     `json:"unschedulable"`
}

type Task struct {
	NodeName               string   `json:"nodeName"`
	NodeSelector           Selector `json:"nodeSelector"`
	RequestedResource      Resource `json:"requestedResource"`
	ToleratesUnschedulable bool     `json:"toleratesUnschedulable"`
}

type Label map[string]string

type Resource struct {
	MilliCPU int64 `json:"milliCPU"`
	Memory   int64 `json:"memory"`
	Storage  int64 `json:"storage"`
}

type Selector map[string][]string

type Score interface {
	Score(*Args) Result
}

type Result struct {
	Score int64
}

type ScoreRPC struct {
	client *rpc.Client
}

func (n *ScoreRPC) Score(args *Args) Result {
	var resp Result
	if err := n.client.Call("Plugin.Score", args, &resp); err != nil {
		panic(err)
	}
	return resp
}

type ScoreRPCServer struct {
	Impl Score
}

func (n *ScoreRPCServer) Score(args *Args, resp *Result) error {
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