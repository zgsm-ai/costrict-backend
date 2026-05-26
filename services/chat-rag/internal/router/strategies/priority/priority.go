package priority

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/zgsm-ai/chat-rag/internal/bootstrap"
	"github.com/zgsm-ai/chat-rag/internal/config"
	"github.com/zgsm-ai/chat-rag/internal/logger"
	"github.com/zgsm-ai/chat-rag/internal/model"
	"github.com/zgsm-ai/chat-rag/internal/types"
	"go.uber.org/zap"
)

// Strategy implements the priority-based round-robin routing strategy
type Strategy struct {
	cfg            config.PriorityConfig
	priorityGroups map[int]*PriorityGroup
	mu             sync.RWMutex
}

// New creates a new priority strategy instance
func New(cfg config.PriorityConfig) (*Strategy, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid priority config: %w", err)
	}

	s := &Strategy{
		cfg:            cfg,
		priorityGroups: make(map[int]*PriorityGroup),
	}

	// Initialize priority groups from config
	if err := s.initializePriorityGroups(); err != nil {
		return nil, fmt.Errorf("failed to initialize priority groups: %w", err)
	}

	return s, nil
}

// Name returns the strategy name
func (s *Strategy) Name() string {
	return "priority"
}

// Run implements the Strategy interface
func (s *Strategy) Run(
	ctx context.Context,
	svcCtx *bootstrap.ServiceContext,
	headers *http.Header,
	req *types.ChatCompletionRequest,
) (string, string, []string, error) {
	if req == nil || len(req.Messages) == 0 {
		return "", "", nil, nil
	}

	if !strings.EqualFold(req.Model, "auto") {
		return "", "", nil, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	identity, _ := model.GetIdentityFromContext(ctx)
	selection := s.visibleSelectionFor(identity, time.Now())
	if len(selection.allVisible) == 0 {
		logger.WarnC(ctx, "priority router: no visible candidates found")
		return "", "", nil, errors.New("no visible candidates configured for current user")
	}
	if selection.selectedGroup == nil {
		logger.ErrorC(ctx, "priority router: visible priority group not found")
		return "", "", nil, errors.New("visible priority group not found")
	}

	selectedModel := selection.selectedGroup.selectVisibleModelByRoundRobin(selection.selectedVisible)

	logger.InfoC(ctx, "priority router: model selected",
		zap.String("selectedModel", selectedModel),
		zap.Int("priority", selection.selectedPriority),
		zap.Int("visibleGroupSize", len(selection.selectedVisible)),
	)

	orderedCandidates := s.buildOrderedCandidates(selectedModel, selection.allVisible)
	orderedCandidates = s.appendVisibleFallback(ctx, orderedCandidates, selection.allVisible)

	return selectedModel, "", orderedCandidates, nil
}

// initializePriorityGroups initializes priority groups from configuration
func (s *Strategy) initializePriorityGroups() error {
	for _, candidate := range s.cfg.Candidates {
		if !candidate.Enabled {
			continue
		}

		// Get or create priority group
		group, exists := s.priorityGroups[candidate.Priority]
		if !exists {
			group = newPriorityGroup(candidate.Priority)
			s.priorityGroups[candidate.Priority] = group
		}

		// Add model to group
		model := &ModelCandidate{
			modelName:   candidate.ModelName,
			priority:    candidate.Priority,
			weight:      candidate.Weight,
			enabled:     candidate.Enabled,
			minVipLevel: candidate.MinVipLevel,
		}
		group.addModel(model)
	}

	if len(s.priorityGroups) == 0 {
		return errors.New("no enabled candidates configured")
	}

	return nil
}

// buildOrderedCandidates builds an ordered list of candidates for degradation
// The selectedModel is always first, followed by other visible models sorted by priority (ascending)
// and weight (descending) within the same priority
func (s *Strategy) buildOrderedCandidates(selectedModel string, visibleCandidates []*ModelCandidate) []string {
	result := []string{selectedModel}

	ordered := append([]*ModelCandidate(nil), visibleCandidates...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].priority == ordered[j].priority {
			return ordered[i].weight > ordered[j].weight
		}
		return ordered[i].priority < ordered[j].priority
	})

	for _, candidate := range ordered {
		if candidate.modelName != selectedModel {
			result = append(result, candidate.modelName)
		}
	}

	return result
}

func containsModel(models []string, modelName string) bool {
	for _, m := range models {
		if m == modelName {
			return true
		}
	}
	return false
}

func findCandidateByName(candidates []*ModelCandidate, modelName string) *ModelCandidate {
	for _, c := range candidates {
		if c.modelName == modelName {
			return c
		}
	}
	return nil
}

func (s *Strategy) appendVisibleFallback(
	ctx context.Context,
	ordered []string,
	visibleCandidates []*ModelCandidate,
) []string {
	if s.cfg.FallbackModelName == "" {
		return ordered
	}
	if containsModel(ordered, s.cfg.FallbackModelName) {
		return ordered
	}
	fallback := findCandidateByName(visibleCandidates, s.cfg.FallbackModelName)
	if fallback == nil {
		if !s.isCandidateDeclared(s.cfg.FallbackModelName) {
			logger.WarnC(ctx, "priority router: fallback model is not declared in candidates",
				zap.String("fallbackModelName", s.cfg.FallbackModelName))
		} else {
			logger.WarnC(ctx, "priority router: fallback model is not visible to current user",
				zap.String("fallbackModelName", s.cfg.FallbackModelName))
		}
		return ordered
	}
	return append(ordered, fallback.modelName)
}

func (s *Strategy) isCandidateDeclared(modelName string) bool {
	for _, c := range s.cfg.Candidates {
		if c.ModelName == modelName {
			return true
		}
	}
	return false
}

func isCandidateVisible(candidate *ModelCandidate, identity *model.Identity, now time.Time) bool {
	if candidate == nil {
		return false
	}
	required := candidate.minVipLevel
	if required == 0 {
		return true
	}
	if identity == nil || identity.UserInfo == nil {
		return false
	}
	user := identity.UserInfo
	hasVipLevel := user.Vip >= required
	isVipUnexpired := user.VipExpire == nil || now.Before(*user.VipExpire)
	return hasVipLevel && isVipUnexpired
}

func filterVisibleModels(models []*ModelCandidate, identity *model.Identity, now time.Time) []*ModelCandidate {
	visible := make([]*ModelCandidate, 0, len(models))
	for _, candidate := range models {
		if isCandidateVisible(candidate, identity, now) {
			visible = append(visible, candidate)
		}
	}
	return visible
}

type visibleSelection struct {
	selectedPriority int
	selectedGroup    *PriorityGroup
	selectedVisible  []*ModelCandidate
	allVisible       []*ModelCandidate
}

func (s *Strategy) visibleSelectionFor(identity *model.Identity, now time.Time) visibleSelection {
	priorities := make([]int, 0, len(s.priorityGroups))
	for priority := range s.priorityGroups {
		priorities = append(priorities, priority)
	}
	sort.Ints(priorities)

	selection := visibleSelection{
		selectedPriority: -1,
		allVisible:       make([]*ModelCandidate, 0),
	}

	for _, priority := range priorities {
		group := s.priorityGroups[priority]
		if group == nil {
			continue
		}
		visible := filterVisibleModels(group.getModels(), identity, now)
		if len(visible) == 0 {
			continue
		}
		if selection.selectedGroup == nil {
			selection.selectedPriority = priority
			selection.selectedGroup = group
			selection.selectedVisible = visible
		}
		selection.allVisible = append(selection.allVisible, visible...)
	}

	return selection
}

// validateConfig validates the priority configuration
func validateConfig(cfg config.PriorityConfig) error {
	if len(cfg.Candidates) == 0 {
		return errors.New("no candidates configured")
	}

	hasEnabled := false
	for _, candidate := range cfg.Candidates {
		if candidate.ModelName == "" {
			return errors.New("candidate modelName is empty")
		}
		if candidate.Priority < 0 || candidate.Priority > 999 {
			return fmt.Errorf("priority must be between 0 and 999 for model %s", candidate.ModelName)
		}
		if candidate.Weight < 1 || candidate.Weight > 100 {
			return fmt.Errorf("weight must be between 1 and 100 for model %s", candidate.ModelName)
		}
		if candidate.MinVipLevel < 0 {
			return fmt.Errorf("minVipLevel must be >= 0 for model %s", candidate.ModelName)
		}
		if candidate.Enabled {
			hasEnabled = true
		}
	}

	if !hasEnabled {
		return errors.New("at least one candidate must be enabled")
	}

	return nil
}
