package factory

import (
	"github.com/kozmod/progen/internal/entity"
	"slices"
	"strings"
)

type (
	actionFilter interface {
		MatchString(s string) bool
	}
)

type ActionFilter struct {
	skipFilter     actionFilter
	selectedGroups map[string]struct{}
	groupsByAction map[string]map[string]struct{}
	manualActions  map[string]struct{}

	logger entity.Logger
}

func NewActionFilter(
	skipActions []string,
	selectedGroups []string,
	groupsByAction map[string]map[string]struct{},
	manualActionsSet map[string]struct{},
	logger entity.Logger,

) *ActionFilter {
	selectedGroupsSet := entity.SliceSet(selectedGroups)
	skipActions = slices.Compact(skipActions)

	switch {
	case len(selectedGroups) > 0:
		logger.Infof("groups will be execute: [%s]", strings.Join(selectedGroups, entity.LogSliceSep))
	case len(manualActionsSet) > 0:
		manualActions := make([]string, 0, len(manualActionsSet))
		for action := range manualActionsSet {
			manualActions = append(manualActions, action)
		}
		logger.Infof("manual actions will be skipped: [%s]", strings.Join(manualActions, entity.LogSliceSep))
	}

	return &ActionFilter{
		skipFilter:     entity.NewRegexpChain(skipActions...),
		selectedGroups: selectedGroupsSet,
		groupsByAction: groupsByAction,
		manualActions:  manualActionsSet,
		logger:         logger,
	}
}

func (f *ActionFilter) MatchString(action string) bool {
	if f.skipFilter.MatchString(action) {
		f.logger.Infof("action will be skipped: [%s]", action)
		return false
	}

	switch {
	case len(f.selectedGroups) > 0:
		if groups, ok := f.groupsByAction[action]; ok {
			for group := range groups {
				if _, ok = f.selectedGroups[group]; ok {
					return true
				}
			}
		}
		return false
	default:
		_, ok := f.manualActions[action]
		return !ok
	}
}
