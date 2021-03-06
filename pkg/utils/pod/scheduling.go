/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pod

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func FailedToSchedule(pod *v1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodScheduled && condition.Reason == v1.PodReasonUnschedulable {
			return true
		}
	}
	return false
}

// IsSchedulable returns true if the pod can schedule to the node
func IsSchedulable(pod *v1.PodSpec, node *v1.Node) bool {
	// Tolerate Taints
	if !ToleratesAllTaints(pod, node.Spec.Taints) {
		return false
	}
	// Match Node Selector labels
	if !labels.SelectorFromSet(pod.NodeSelector).Matches(labels.Set(node.Labels)) {
		return false
	}
	// TODO, support node affinity
	return true
}

// ToleratesAllTaints returns true if the pod tolerates all taints
func ToleratesAllTaints(pod *v1.PodSpec, taints []v1.Taint) bool {
	for _, taint := range taints {
		if !ToleratesTaint(pod, taint) {
			return false
		}
	}
	return true
}

// ToleratesTaint returns true if the pod tolerates the taint
// https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/#concepts
func ToleratesTaint(pod *v1.PodSpec, taint v1.Taint) bool {
	// Soft constraints are consider to be always tolerated.
	if taint.Effect == v1.TaintEffectPreferNoSchedule {
		return true
	}
	for _, toleration := range pod.Tolerations {
		if toleration.Key == taint.Key {
			if toleration.Operator == v1.TolerationOpExists {
				return true
			}
			if toleration.Operator == v1.TolerationOpEqual && toleration.Value == taint.Value {
				return true
			}
		}
	}
	return false
}
