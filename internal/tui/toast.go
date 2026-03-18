package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ToastType represents the type of notification
type ToastType string

const (
	ToastSuccess ToastType = "success"
	ToastError   ToastType = "error"
	ToastWarning ToastType = "warning"
	ToastInfo    ToastType = "info"
)

// Toast represents a notification message
type Toast struct {
	Message   string
	Type      ToastType
	CreatedAt time.Time
	Duration  time.Duration
}

// ToastManager handles multiple toast notifications
type ToastManager struct {
	toasts   []Toast
	maxCount int
}

var defaultToastManager = &ToastManager{maxCount: 5}

func (tm *ToastManager) Add(message string, toastType ToastType, duration time.Duration) {
	if duration == 0 {
		duration = 3 * time.Second
	}

	toast := Toast{
		Message:   message,
		Type:      toastType,
		CreatedAt: time.Now(),
		Duration:  duration,
	}

	tm.toasts = append([]Toast{toast}, tm.toasts...)

	if len(tm.toasts) > tm.maxCount {
		tm.toasts = tm.toasts[:tm.maxCount]
	}
}

func (tm *ToastManager) RemoveExpired() {
	var active []Toast
	for _, t := range tm.toasts {
		if time.Since(t.CreatedAt) < t.Duration {
			active = append(active, t)
		}
	}
	tm.toasts = active
}

func (tm *ToastManager) IsActive() bool {
	tm.RemoveExpired()
	return len(tm.toasts) > 0
}

func (tm *ToastManager) Render(width int) string {
	tm.RemoveExpired()
	if len(tm.toasts) == 0 {
		return ""
	}

	var lines []string
	for _, toast := range tm.toasts {
		elapsed := time.Since(toast.CreatedAt)
		remaining := toast.Duration - elapsed
		if remaining <= 0 {
			continue
		}

		age := int(remaining.Seconds())
		progress := strings.Repeat("█", age) + strings.Repeat("░", int(toast.Duration.Seconds())-age)
		if len(progress) > 10 {
			progress = progress[:10]
		}

		style := tm.getStyle(toast.Type)
		icon := tm.getIcon(toast.Type)

		maxWidth := width - 20
		message := toast.Message
		if len(message) > maxWidth {
			message = message[:maxWidth-3] + "..."
		}

		line := fmt.Sprintf("%s %s %s", icon, message, style.Render("["+progress+"]"))
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return ""
	}

	result := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E293B")).
		Padding(1, 2).
		Width(width).
		Render(strings.Join(lines, "\n"))

	return result
}

func (tm *ToastManager) getStyle(toastType ToastType) lipgloss.Style {
	switch toastType {
	case ToastSuccess:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981"))
	case ToastError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	case ToastWarning:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6"))
	}
}

func (tm *ToastManager) getIcon(toastType ToastType) string {
	switch toastType {
	case ToastSuccess:
		return "✓"
	case ToastError:
		return "✗"
	case ToastWarning:
		return "⚠"
	default:
		return "ℹ"
	}
}

// Convenience functions
func ShowToast(message string, toastType ToastType) {
	defaultToastManager.Add(message, toastType, 0)
}

func ShowSuccess(message string) {
	ShowToast(message, ToastSuccess)
}

func ShowError(message string) {
	ShowToast(message, ToastError)
}

func ShowWarning(message string) {
	ShowToast(message, ToastWarning)
}

func ShowInfo(message string) {
	ShowToast(message, ToastInfo)
}
