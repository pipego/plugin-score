package main

import (
	"math"

	gop "github.com/hashicorp/go-plugin"

	"github.com/pipego/plugin-score/common"
	"github.com/pipego/scheduler/plugin"
)

var (
	resourceToWeightMapAllocation = map[string]int64{
		plugin.ResourceCPU:     plugin.DefaultCPUWeight,
		plugin.ResourceMemory:  plugin.DefaultMemoryWeight,
		plugin.ResourceStorage: plugin.DefaultStorageWeight,
	}
)

type NodeResourcesBalancedAllocation struct{}
type resourceToValueMapAllocation map[string]int64

func (n *NodeResourcesBalancedAllocation) Score(args *plugin.Args) common.Result {
	requested := make(resourceToValueMapAllocation)
	allocatable := make(resourceToValueMapAllocation)

	for resource := range resourceToWeightMapAllocation {
		alloc, req := n.calculateResourceAllocatableRequest(&args.Node, &args.Task, resource)
		if alloc != 0 {
			allocatable[resource], requested[resource] = alloc, req
		}
	}

	return common.Result{
		Score: n.balancedResourceScorer(requested, allocatable),
	}
}

func (n *NodeResourcesBalancedAllocation) calculateResourceAllocatableRequest(node *plugin.Node, task *plugin.Task, resource string) (int64, int64) {
	taskRequest := n.calculateTaskResourceRequest(task, resource)

	switch resource {
	case plugin.ResourceCPU:
		return node.AllocatableResource.MilliCPU, node.RequestedResource.MilliCPU + taskRequest
	case plugin.ResourceMemory:
		return node.AllocatableResource.Memory, node.RequestedResource.Memory + taskRequest
	case plugin.ResourceStorage:
		return node.AllocatableResource.Storage, node.RequestedResource.Storage + taskRequest
	default:
		// BYPASS
	}

	return 0, 0
}

func (n *NodeResourcesBalancedAllocation) calculateTaskResourceRequest(task *plugin.Task, resource string) int64 {
	switch resource {
	case plugin.ResourceCPU:
		if task.RequestedResource.MilliCPU == 0 {
			return plugin.DefaultMilliCPURequest
		}
		return task.RequestedResource.MilliCPU
	case plugin.ResourceMemory:
		if task.RequestedResource.Memory == 0 {
			return plugin.DefaultMemoryRequest
		}
		return task.RequestedResource.Memory
	case plugin.ResourceStorage:
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
	return int64((1 - std) * float64(plugin.MaxNodeScore))
}

// nolint:typecheck
func main() {
	config := gop.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "plugin-score",
		MagicCookieValue: "plugin-score",
	}

	pluginMap := map[string]gop.Plugin{
		"NodeResourcesBalancedAllocation": &common.ScorePlugin{Impl: &NodeResourcesBalancedAllocation{}},
	}

	gop.Serve(&gop.ServeConfig{
		HandshakeConfig: config,
		Plugins:         pluginMap,
	})
}
