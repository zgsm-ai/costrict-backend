package priority

import (
	"fmt"
	"strings"
	"sync"
)

type ModelCandidate struct {
	modelName   string
	priority    int
	weight      int
	enabled     bool
	minVipLevel int
}

type PriorityGroup struct {
	priority int
	models   []*ModelCandidate
	states   map[string][]int
	mu       sync.Mutex
}

func newPriorityGroup(priority int) *PriorityGroup {
	return &PriorityGroup{
		priority: priority,
		models:   make([]*ModelCandidate, 0),
		states:   make(map[string][]int),
	}
}

func (pg *PriorityGroup) addModel(model *ModelCandidate) {
	pg.models = append(pg.models, model)
}

func buildVisibleSetKey(visible []*ModelCandidate) string {
	var sb strings.Builder
	for _, model := range visible {
		fmt.Fprintf(&sb, "%d:%s", len(model.modelName), model.modelName)
	}
	return sb.String()
}

func (pg *PriorityGroup) selectVisibleModelByRoundRobin(visible []*ModelCandidate) string {
	if len(visible) == 0 {
		return ""
	}
	if len(visible) == 1 {
		return visible[0].modelName
	}

	key := buildVisibleSetKey(visible)

	pg.mu.Lock()
	defer pg.mu.Unlock()

	weights := pg.states[key]
	if len(weights) != len(visible) {
		weights = make([]int, len(visible))
	}

	totalWeight := 0
	for _, model := range visible {
		totalWeight += model.weight
	}

	maxWeight := -1
	selectedIdx := 0
	for i, model := range visible {
		weights[i] += model.weight
		if weights[i] > maxWeight {
			maxWeight = weights[i]
			selectedIdx = i
		}
	}

	weights[selectedIdx] -= totalWeight
	pg.states[key] = weights

	return visible[selectedIdx].modelName
}

func (pg *PriorityGroup) getModels() []*ModelCandidate {
	return pg.models
}
