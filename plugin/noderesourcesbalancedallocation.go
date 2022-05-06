package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/pipego/plugin-score/proto"
	"math"
)

var (
	resourceToWeightMapAllocation = map[string]int64{
		proto.ResourceCPU:     proto.DefaultCPUWeight,
		proto.ResourceMemory:  proto.DefaultMemoryWeight,
		proto.ResourceStorage: proto.DefaultStorageWeight,
	}
)

type NodeResourcesBalancedAllocation struct{}
type resourceToValueMapAllocation map[string]int64

func (n *NodeResourcesBalancedAllocation) Score(args *proto.Args) proto.Result {
	requested := make(resourceToValueMapAllocation)
	allocatable := make(resourceToValueMapAllocation)

	for resource := range resourceToWeightMapAllocation {
		alloc, req := n.calculateResourceAllocatableRequest(&args.Node, &args.Task, resource)
		if alloc != 0 {
			allocatable[resource], requested[resource] = alloc, req
		}
	}

	return proto.Result{
		Score: n.balancedResourceScorer(requested, allocatable),
	}
}

func (n *NodeResourcesBalancedAllocation) calculateResourceAllocatableRequest(node *proto.Node, task *proto.Task, resource string) (int64, int64) {
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

func (n *NodeResourcesBalancedAllocation) calculateTaskResourceRequest(task *proto.Task, resource string) int64 {
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

func (n *NodeResourcesBalancedAllocation) balancedResourceScorer(requested, allocable resourceToValueMapAllocation) int64 {
	var resourceToFractions []float64
	var totalFraction float64

	for name, value := range requested {
		fraction := float64(value) / float64(allocable[name])
		if fraction > 1 {
			fraction = 1
		}
		totalFraction += fraction
		resourceToFractions = append(resourceToFractions, fraction)
	}

	std := 0.0

	// For most cases, resources are limited to cpu and memory, the std could be simplified to std := (fraction1-fraction2)/2
	// len(fractions) > 2: calculate std based on the well-known formula - root square of Σ((fraction(i)-mean)^2)/len(fractions)
	// Otherwise, set the std to zero is enough.
	if len(resourceToFractions) == 2 {
		std = math.Abs((resourceToFractions[0] - resourceToFractions[1]) / 2)
	} else if len(resourceToFractions) > 2 {
		mean := totalFraction / float64(len(resourceToFractions))
		var sum float64
		for _, fraction := range resourceToFractions {
			sum = sum + (fraction-mean)*(fraction-mean)
		}
		std = math.Sqrt(sum / float64(len(resourceToFractions)))
	}

	// STD (standard deviation) is always a positive value. 1-deviation lets the score to be higher for node which has least deviation and
	// multiplying it with `MaxNodeScore` provides the scaling factor needed.
	return int64((1 - std) * float64(proto.MaxNodeScore))
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
