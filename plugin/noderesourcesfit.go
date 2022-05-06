package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-score/proto"
)

var (
	resourceToWeightMap = map[string]int64{
		proto.ResourceCPU:     proto.DefaultCPUWeight,
		proto.ResourceMemory:  proto.DefaultMemoryWeight,
		proto.ResourceStorage: proto.DefaultStorageWeight,
	}
)

type NodeResourcesFit struct{}
type resourceToValueMap map[string]int64

func (n *NodeResourcesFit) Score(args *proto.Args) proto.Result {
	requested := make(resourceToValueMap)
	allocatable := make(resourceToValueMap)

	for resource := range resourceToWeightMap {
		alloc, req := n.calculateResourceAllocatableRequest(&args.Node, &args.Task, resource)
		if alloc != 0 {
			allocatable[resource], requested[resource] = alloc, req
		}
	}

	return proto.Result{
		Score: n.leastResourceScorer(requested, allocatable),
	}
}

func (n *NodeResourcesFit) calculateResourceAllocatableRequest(node *proto.Node, task *proto.Task, resource string) (int64, int64) {
	taskRequest := n.calculateTaskResourceRequest(task, resource)

	switch resource {
	case proto.ResourceCPU:
		return node.AllocatableResource.MilliCPU, node.RequestedResource.MilliCPU + taskRequest
	case proto.ResourceMemory:
		return node.AllocatableResource.Memory, node.RequestedResource.Memory + taskRequest
	case proto.ResourceStorage:
		return node.AllocatableResource.Storage, node.RequestedResource.Storage + taskRequest
	default:
		// BYPASS
	}

	return 0, 0
}

func (n *NodeResourcesFit) calculateTaskResourceRequest(task *proto.Task, resource string) int64 {
	switch resource {
	case proto.ResourceCPU:
		if task.RequestedResource.MilliCPU == 0 {
			return proto.DefaultMilliCPURequest
		}
		return task.RequestedResource.MilliCPU
	case proto.ResourceMemory:
		if task.RequestedResource.Memory == 0 {
			return proto.DefaultMemoryRequest
		}
		return task.RequestedResource.Memory
	case proto.ResourceStorage:
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
func (n *NodeResourcesFit) leastResourceScorer(requested, allocable resourceToValueMap) int64 {
	var nodeScore, weightSum int64

	for resource := range requested {
		weight := resourceToWeightMap[resource]
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

	return ((capacity - requested) * proto.MaxNodeScore) / capacity
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
