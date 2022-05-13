package main

import (
	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/scheduler/common"
	"github.com/pipego/scheduler/plugin"
)

var (
	resourceToWeightMapFit = map[string]int64{
		common.ResourceCPU:     common.DefaultCPUWeight,
		common.ResourceMemory:  common.DefaultMemoryWeight,
		common.ResourceStorage: common.DefaultStorageWeight,
	}
)

type NodeResourcesFit struct{}
type resourceToValueMapFit map[string]int64

func (n *NodeResourcesFit) Run(args *common.Args) plugin.ScoreResult {
	requested := make(resourceToValueMapFit)
	allocatable := make(resourceToValueMapFit)

	for resource := range resourceToWeightMapFit {
		alloc, req := n.calculateResourceAllocatableRequest(&args.Node, &args.Task, resource)
		if alloc != 0 {
			allocatable[resource], requested[resource] = alloc, req
		}
	}

	return plugin.ScoreResult{
		Score: n.leastResourceScorer(requested, allocatable),
	}
}

func (n *NodeResourcesFit) calculateResourceAllocatableRequest(node *common.Node, task *common.Task, resource string) (int64, int64) {
	taskRequest := n.calculateTaskResourceRequest(task, resource)

	switch resource {
	case common.ResourceCPU:
		return node.AllocatableResource.MilliCPU, node.RequestedResource.MilliCPU + taskRequest
	case common.ResourceMemory:
		return node.AllocatableResource.Memory, node.RequestedResource.Memory + taskRequest
	case common.ResourceStorage:
		return node.AllocatableResource.Storage, node.RequestedResource.Storage + taskRequest
	default:
		// BYPASS
	}

	return 0, 0
}

func (n *NodeResourcesFit) calculateTaskResourceRequest(task *common.Task, resource string) int64 {
	switch resource {
	case common.ResourceCPU:
		if task.RequestedResource.MilliCPU == 0 {
			return common.DefaultMilliCPURequest
		}
		return task.RequestedResource.MilliCPU
	case common.ResourceMemory:
		if task.RequestedResource.Memory == 0 {
			return common.DefaultMemoryRequest
		}
		return task.RequestedResource.Memory
	case common.ResourceStorage:
		return task.RequestedResource.Storage
	default:
		// BYPASS
	}

	return 0
}

// leastResourceScorer favors nodes with fewer requested resources.
// It calculates the percentage of memory, CPU and other resources requested by tasks scheduled on the node, and
// prioritizes based on the minimum of the average of the fraction of requested to capacity.
//
// Details:
// (cpu((capacity-requested)*MaxNodeScore*cpuWeight/capacity) + memory((capacity-requested)*MaxNodeScore*memoryWeight/capacity) + ...)/weightSum
func (n *NodeResourcesFit) leastResourceScorer(requested, allocable resourceToValueMapFit) int64 {
	var nodeScore, weightSum int64

	for resource := range requested {
		weight := resourceToWeightMapFit[resource]
		resourceScore := n.leastRequestedScore(requested[resource], allocable[resource])
		nodeScore += resourceScore * weight
		weightSum += weight
	}

	if weightSum == 0 {
		return 0
	}

	return nodeScore / weightSum
}

func (n *NodeResourcesFit) leastRequestedScore(requested, capacity int64) int64 {
	if capacity == 0 {
		return 0
	}

	if requested > capacity {
		return 0
	}

	return ((capacity - requested) * common.MaxNodeScore) / capacity
}

// nolint:typecheck
func main() {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin",
		MagicCookieValue: "plugin",
	}

	pluginMap := map[string]gop.Plugin{
		"NodeResourcesFit": &plugin.Score{Impl: &NodeResourcesFit{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
