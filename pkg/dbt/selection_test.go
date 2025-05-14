package dbt

import (
	"reflect"
	"testing"
)

func TestParseSelectionSyntax(t *testing.T) {
	tests := []struct {
		name     string
		selectors []string
		want     []SelectionCriteria
		wantErr  bool
	}{
		{
			name:     "simple model name",
			selectors: []string{"my_model"},
			want: []SelectionCriteria{
				{
					Type:      SelectionName,
					Value:     "my_model",
					Modifier:  ModifierNone,
					Depth:     0,
					Exclude:   false,
					Direction: "",
				},
			},
			wantErr: false,
		},
		{
			name:     "tag selection",
			selectors: []string{"tag:daily"},
			want: []SelectionCriteria{
				{
					Type:      SelectionTag,
					Value:     "daily",
					Modifier:  ModifierNone,
					Depth:     0,
					Exclude:   false,
					Direction: "",
				},
			},
			wantErr: false,
		},
		{
			name:     "downstream modifier",
			selectors: []string{"my_model+downstream"},
			want: []SelectionCriteria{
				{
					Type:      SelectionName,
					Value:     "my_model",
					Modifier:  ModifierDownstream,
					Depth:     0,
					Exclude:   false,
					Direction: "downstream",
				},
			},
			wantErr: false,
		},
		{
			name:     "upstream modifier",
			selectors: []string{"my_model+upstream"},
			want: []SelectionCriteria{
				{
					Type:      SelectionName,
					Value:     "my_model",
					Modifier:  ModifierUpstream,
					Depth:     0,
					Exclude:   false,
					Direction: "upstream",
				},
			},
			wantErr: false,
		},
		{
			name:     "depth modifier",
			selectors: []string{"my_model+3"},
			want: []SelectionCriteria{
				{
					Type:      SelectionName,
					Value:     "my_model",
					Modifier:  ModifierPlusN,
					Depth:     3,
					Exclude:   false,
					Direction: "",
				},
			},
			wantErr: false,
		},
		{
			name:     "exclusion",
			selectors: []string{"!my_model"},
			want: []SelectionCriteria{
				{
					Type:      SelectionName,
					Value:     "my_model",
					Modifier:  ModifierNone,
					Depth:     0,
					Exclude:   true,
					Direction: "",
				},
			},
			wantErr: false,
		},
		{
			name:     "multiple selectors",
			selectors: []string{"tag:daily", "my_model+downstream", "!exclude_model"},
			want: []SelectionCriteria{
				{
					Type:      SelectionTag,
					Value:     "daily",
					Modifier:  ModifierNone,
					Depth:     0,
					Exclude:   false,
					Direction: "",
				},
				{
					Type:      SelectionName,
					Value:     "my_model",
					Modifier:  ModifierDownstream,
					Depth:     0,
					Exclude:   false,
					Direction: "downstream",
				},
				{
					Type:      SelectionName,
					Value:     "exclude_model",
					Modifier:  ModifierNone,
					Depth:     0,
					Exclude:   true,
					Direction: "",
				},
			},
			wantErr: false,
		},
		{
			name:     "unknown selection type",
			selectors: []string{"unknown:value"},
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSelectionSyntax(tt.selectors)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSelectionSyntax() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSelectionSyntax() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplySelection(t *testing.T) {
	// Create a test manifest
	manifest := &DBTManifest{
		Nodes: map[string]*DBTModel{
			"model.test.model1": {
				Name:         "model1",
				Schema:       "test",
				ResourceType: "model",
				Tags:         []string{"daily", "important"},
			},
			"model.test.model2": {
				Name:         "model2",
				Schema:       "test",
				ResourceType: "model",
				Tags:         []string{"weekly"},
			},
			"model.test.model3": {
				Name:         "model3",
				Schema:       "test",
				ResourceType: "model",
				Tags:         []string{"daily"},
			},
			"snapshot.test.snapshot1": {
				Name:         "snapshot1",
				Schema:       "test",
				ResourceType: "snapshot",
			},
		},
	}

	tests := []struct {
		name     string
		criteria []SelectionCriteria
		want     int // Number of models expected
		wantErr  bool
	}{
		{
			name:     "empty criteria (select all models)",
			criteria: []SelectionCriteria{},
			want:     3, // All models, no snapshots
			wantErr:  false,
		},
		{
			name: "select by name",
			criteria: []SelectionCriteria{
				{
					Type:  SelectionName,
					Value: "model1",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "select by tag",
			criteria: []SelectionCriteria{
				{
					Type:  SelectionTag,
					Value: "daily",
				},
			},
			want:    2, // model1 and model3
			wantErr: false,
		},
		{
			name: "select by name with exclusion",
			criteria: []SelectionCriteria{
				{
					Type:  SelectionName,
					Value: "model1",
				},
				{
					Type:    SelectionName,
					Value:   "model1",
					Exclude: true,
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "multiple criteria",
			criteria: []SelectionCriteria{
				{
					Type:  SelectionTag,
					Value: "daily",
				},
				{
					Type:    SelectionTag,
					Value:   "important",
					Exclude: true,
				},
			},
			want:    1, // model3 (daily but not important)
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ApplySelection(manifest, tt.criteria)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplySelection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ApplySelection() returned %d models, want %d", len(got), tt.want)
			}
		})
	}
}
