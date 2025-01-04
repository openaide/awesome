package main

import (
	"context"
	"testing"
)

func TestBuildRunGptr(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"Test", "Renewable energy sources and their potential"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if err := BuildGPTRImage(ctx); err != nil {
				t.Errorf("BuildGPTRImage() error = %v", err)
				return
			}
			if err := RunGPTRContainer(ctx, tt.query, "outputs/gptr"); err != nil {
				t.Errorf("RunGPTRContainer() error = %v", err)
				return
			}
		})
	}
}
