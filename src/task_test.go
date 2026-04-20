package main

import (
	"reflect"
	"testing"
)

func TestGetLinks(t *testing.T) {
	tests := []struct {
		description string
		want        []string
	}{
		{
			description: "No links here",
			want:        nil,
		},
		{
			description: "One link: https://google.com",
			want:        []string{"https://google.com"},
		},
		{
			description: "Two links: http://example.com and https://test.org/path",
			want:        []string{"http://example.com", "https://test.org/path"},
		},
	}

	for _, tt := range tests {
		task := Task{description: tt.description}
		got := task.GetLinks()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetLinks() for %q = %v, want %v", tt.description, got, tt.want)
		}
	}
}
