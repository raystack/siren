package notification

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/odpf/siren/core/alert"
	"github.com/stretchr/testify/require"
)

func Test_removeDuplicateStringValues(t *testing.T) {
	tests := []struct {
		name     string
		strSlice []string
		want     []string
	}{
		{
			name:     "should remove duplicated string value in slice",
			strSlice: []string{"a", "b", "c", "a"},
			want:     []string{"a", "b", "c"},
		},
		{
			name:     "should return as-is if no duplicated string",
			strSlice: []string{"a", "b", "c"},
			want:     []string{"a", "b", "c"},
		},
		{
			name:     "should return empty slice if empty",
			strSlice: []string{},
			want:     []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicateStringValues(tt.strSlice)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("got diff = %v", diff)
			}
		})
	}
}

func Test_groupByLabels(t *testing.T) {
	hashKey1, err := hashstructure.Hash(map[string]string{
		"key1": "val1",
		"key2": "val2",
	}, hashstructure.FormatV2, nil)
	require.NoError(t, err)

	hashKey2, err := hashstructure.Hash(map[string]string{
		"key2": "val2",
		"key3": "val3",
	}, hashstructure.FormatV2, nil)
	require.NoError(t, err)

	hashKey3, err := hashstructure.Hash(map[string]string{
		"key1": "val1",
		"key3": "val3",
	}, hashstructure.FormatV2, nil)
	require.NoError(t, err)

	tests := []struct {
		name    string
		alerts  []alert.Alert
		want    map[uint64][]alert.Alert
		wantErr bool
	}{
		{
			name: "shoudl group alerts if labels are same",
			alerts: []alert.Alert{
				{
					ID: 12,
					Labels: map[string]string{
						"key1": "val1",
						"key2": "val2",
					},
				},
				{
					ID: 34,
					Labels: map[string]string{
						"key1": "val1",
						"key2": "val2",
					},
				},
				{
					ID: 56,
					Labels: map[string]string{
						"key2": "val2",
						"key3": "val3",
					},
				},
				{
					ID: 78,
					Labels: map[string]string{
						"key3": "val3",
						"key2": "val2",
					},
				},
				{
					ID: 910,
					Labels: map[string]string{
						"key1": "val1",
						"key3": "val3",
					},
				},
			},
			want: map[uint64][]alert.Alert{
				hashKey1: {
					{
						ID: 12,
						Labels: map[string]string{
							"key1": "val1",
							"key2": "val2",
						},
					},
					{
						ID: 34,
						Labels: map[string]string{
							"key1": "val1",
							"key2": "val2",
						},
					},
				},
				hashKey2: {
					{
						ID: 56,
						Labels: map[string]string{
							"key2": "val2",
							"key3": "val3",
						},
					},
					{
						ID: 78,
						Labels: map[string]string{
							"key3": "val3",
							"key2": "val2",
						},
					},
				},
				hashKey3: {
					{
						ID: 910,
						Labels: map[string]string{
							"key1": "val1",
							"key3": "val3",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := groupByLabels(tt.alerts)
			if (err != nil) != tt.wantErr {
				t.Errorf("groupByLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("got diff = %v", diff)
			}
		})
	}
}
