package dbt

import (
	"fmt"
	"strings"
)

// SelectionType represents the type of selection
type SelectionType int

const (
	// SelectionName selects models by name
	SelectionName SelectionType = iota
	// SelectionTag selects models by tag
	SelectionTag
	// SelectionPath selects models by path
	SelectionPath
	// SelectionPackage selects models by package
	SelectionPackage
	// SelectionConfig selects models by config
	SelectionConfig
	// SelectionExposure selects models by exposure
	SelectionExposure
	// SelectionSource selects models by source
	SelectionSource
	// SelectionState selects models by state
	SelectionState
	// SelectionVersion selects models by version
	SelectionVersion
)

// SelectionModifier represents a modifier for selection
type SelectionModifier int

const (
	// ModifierNone indicates no modifier
	ModifierNone SelectionModifier = iota
	// ModifierPlus indicates the + modifier (include selected models)
	ModifierPlus
	// ModifierPlusN indicates the +N modifier (include N generations of children)
	ModifierPlusN
	// ModifierAt indicates the @ modifier (include selected models)
	ModifierAt
	// ModifierUpstream indicates the + modifier with upstream direction
	ModifierUpstream
	// ModifierDownstream indicates the + modifier with downstream direction
	ModifierDownstream
	// ModifierChildren indicates the children modifier
	ModifierChildren
	// ModifierParents indicates the parents modifier
	ModifierParents
)

// SelectionCriteria represents a selection criteria
type SelectionCriteria struct {
	Type      SelectionType
	Value     string
	Modifier  SelectionModifier
	Depth     int  // For ModifierPlusN
	Exclude   bool // For exclusion with ! prefix
	Direction string
}

// ParseSelectionSyntax parses dbt selection syntax
func ParseSelectionSyntax(selectors []string) ([]SelectionCriteria, error) {
	criteria := make([]SelectionCriteria, 0, len(selectors))

	for _, selector := range selectors {
		// Check for exclusion
		exclude := false
		if strings.HasPrefix(selector, "!") {
			exclude = true
			selector = selector[1:]
		}

		// Check for modifiers
		var modifier SelectionModifier
		var direction string
		var depth int
		var value string

		// Check for +N modifier
		if strings.Contains(selector, "+") {
			parts := strings.Split(selector, "+")
			value = parts[0]

			// Check for upstream/downstream
			if len(parts) > 1 {
				if parts[1] == "upstream" {
					modifier = ModifierUpstream
					direction = "upstream"
				} else if parts[1] == "downstream" {
					modifier = ModifierDownstream
					direction = "downstream"
				} else if parts[1] == "children" {
					modifier = ModifierChildren
					direction = "children"
				} else if parts[1] == "parents" {
					modifier = ModifierParents
					direction = "parents"
				} else {
					// Try to parse as number for +N
					if _, err := fmt.Sscanf(parts[1], "%d", &depth); err == nil {
						modifier = ModifierPlusN
					} else {
						modifier = ModifierPlus
					}
				}
			} else {
				modifier = ModifierPlus
			}
		} else if strings.Contains(selector, "@") {
			parts := strings.Split(selector, "@")
			value = parts[0]
			modifier = ModifierAt
		} else {
			value = selector
			modifier = ModifierNone
		}

		// Determine selection type
		var selType SelectionType
		if strings.Contains(value, ":") {
			parts := strings.Split(value, ":")
			typeStr := parts[0]
			value = parts[1]

			switch typeStr {
			case "tag":
				selType = SelectionTag
			case "path":
				selType = SelectionPath
			case "package":
				selType = SelectionPackage
			case "config":
				selType = SelectionConfig
			case "exposure":
				selType = SelectionExposure
			case "source":
				selType = SelectionSource
			case "state":
				selType = SelectionState
			case "version":
				selType = SelectionVersion
			default:
				return nil, fmt.Errorf("unknown selection type: %s", typeStr)
			}
		} else {
			selType = SelectionName
		}

		criteria = append(criteria, SelectionCriteria{
			Type:      selType,
			Value:     value,
			Modifier:  modifier,
			Depth:     depth,
			Exclude:   exclude,
			Direction: direction,
		})
	}

	return criteria, nil
}

// ApplySelection applies selection criteria to models
func ApplySelection(manifest *DBTManifest, criteria []SelectionCriteria) ([]*DBTModel, error) {
	if len(criteria) == 0 {
		// If no criteria, select all models
		models := make([]*DBTModel, 0, len(manifest.Nodes))
		for _, node := range manifest.Nodes {
			if node.ResourceType == "model" {
				models = append(models, node)
			}
		}
		return models, nil
	}

	// Apply each selection criteria
	selectedModels := make(map[string]*DBTModel)
	excludedModels := make(map[string]*DBTModel)

	for _, criterion := range criteria {
		// Get initial selection
		initialSelection, err := getInitialSelection(manifest, criterion)
		if err != nil {
			return nil, err
		}

		// Apply modifiers
		selection, err := applyModifiers(manifest, initialSelection, criterion)
		if err != nil {
			return nil, err
		}

		// Add to selected or excluded models
		if criterion.Exclude {
			for _, model := range selection {
				excludedModels[model.Name] = model
			}
		} else {
			for _, model := range selection {
				selectedModels[model.Name] = model
			}
		}
	}

	// Remove excluded models from selected models
	for name := range excludedModels {
		delete(selectedModels, name)
	}

	// Convert map to slice
	result := make([]*DBTModel, 0, len(selectedModels))
	for _, model := range selectedModels {
		result = append(result, model)
	}

	return result, nil
}

// getInitialSelection gets the initial selection based on criteria
func getInitialSelection(manifest *DBTManifest, criterion SelectionCriteria) ([]*DBTModel, error) {
	selection := make([]*DBTModel, 0)

	for _, node := range manifest.Nodes {
		if node.ResourceType != "model" {
			continue
		}

		switch criterion.Type {
		case SelectionName:
			if node.Name == criterion.Value {
				selection = append(selection, node)
			}
		case SelectionTag:
			for _, tag := range node.Tags {
				if tag == criterion.Value {
					selection = append(selection, node)
					break
				}
			}
		case SelectionPath:
			// TODO: Implement path selection
			// This would check if the model's path matches the criterion value
		case SelectionPackage:
			// TODO: Implement package selection
			// This would check if the model's package matches the criterion value
		case SelectionConfig:
			// TODO: Implement config selection
			// This would check if the model has a specific config value
		case SelectionExposure:
			// TODO: Implement exposure selection
			// This would check if the model is exposed in a specific way
		case SelectionSource:
			// TODO: Implement source selection
			// This would check if the model uses a specific source
		case SelectionState:
			// TODO: Implement state selection
			// This would check the model's state
		case SelectionVersion:
			// TODO: Implement version selection
			// This would check the model's version
		}
	}

	return selection, nil
}

// applyModifiers applies modifiers to the initial selection
func applyModifiers(manifest *DBTManifest, initialSelection []*DBTModel, criterion SelectionCriteria) ([]*DBTModel, error) {
	switch criterion.Modifier {
	case ModifierNone:
		return initialSelection, nil
	case ModifierPlus:
		// Include the model and all its children
		return initialSelection, nil // TODO: Implement
	case ModifierPlusN:
		// Include the model and N generations of children
		return initialSelection, nil // TODO: Implement
	case ModifierAt:
		// Include the model at a specific version
		return initialSelection, nil // TODO: Implement
	case ModifierUpstream:
		// Include the model and all its upstream dependencies
		return getUpstreamModels(manifest, initialSelection), nil
	case ModifierDownstream:
		// Include the model and all its downstream dependencies
		return getDownstreamModels(manifest, initialSelection), nil
	case ModifierChildren:
		// Include the model's immediate children
		return getChildModels(manifest, initialSelection), nil
	case ModifierParents:
		// Include the model's immediate parents
		return getParentModels(manifest, initialSelection), nil
	default:
		return initialSelection, nil
	}
}

// getUpstreamModels gets all upstream models for the given models
func getUpstreamModels(manifest *DBTManifest, models []*DBTModel) []*DBTModel {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Build a dependency graph from the manifest
	// 2. Traverse the graph to find all upstream dependencies

	// For now, just return the input models
	return models
}

// getDownstreamModels gets all downstream models for the given models
func getDownstreamModels(manifest *DBTManifest, models []*DBTModel) []*DBTModel {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Build a dependency graph from the manifest
	// 2. Traverse the graph to find all downstream dependencies

	// For now, just return the input models
	return models
}

// getChildModels gets immediate child models for the given models
func getChildModels(manifest *DBTManifest, models []*DBTModel) []*DBTModel {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Build a dependency graph from the manifest
	// 2. Find all immediate children

	// For now, just return the input models
	return models
}

// getParentModels gets immediate parent models for the given models
func getParentModels(manifest *DBTManifest, models []*DBTModel) []*DBTModel {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Build a dependency graph from the manifest
	// 2. Find all immediate parents

	// For now, just return the input models
	return models
}
