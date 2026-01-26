package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/stretchr/testify/require"
)

func TestNewSpinner(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "create spinner with message",
			message: "Loading...",
		},
		{
			name:    "create spinner with empty message",
			message: "",
		},
		{
			name:    "create spinner with long message",
			message: "Downloading packages from remote repository...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// #given: 一个消息字符串
			// #when: 调用 NewSpinner 创建 spinner
			s := NewSpinner(tt.message)

			// #then: 返回非空的 spinner.Model
			require.NotNil(t, s)
			require.IsType(t, &spinner.Model{}, s)
		})
	}
}

func TestNewProgress(t *testing.T) {
	tests := []struct {
		name  string
		total int
	}{
		{
			name:  "create progress with positive total",
			total: 100,
		},
		{
			name:  "create progress with small total",
			total: 10,
		},
		{
			name:  "create progress with large total",
			total: 1000,
		},
		{
			name:  "create progress with zero total",
			total: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// #given: 一个总数值
			// #when: 调用 NewProgress 创建 progress bar
			p := NewProgress(tt.total)

			// #then: 返回非空的 progress.Model
			require.NotNil(t, p)
			require.IsType(t, &progress.Model{}, p)
		})
	}
}
