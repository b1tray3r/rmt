package views

import (
	"fmt"
	"strings"

	"github.com/b1tray3r/rmt/internal/tui/themes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HoursSelector handles hours selection with predefined values
type HoursSelector struct {
	options       []float64
	selectedIndex int
	focused       bool
}

// NewHoursSelector creates a new hours selector
func NewHoursSelector() *HoursSelector {
	options := []float64{
		0.25, 0.50, 0.75, 1.00, 1.25, 1.50, 1.75, 2.00,
		2.25, 2.50, 2.75, 3.00, 3.25, 3.50, 3.75, 4.00,
		4.25, 4.50, 4.75, 5.00, 5.25, 5.50, 5.75, 6.00,
		6.25, 6.50, 6.75, 7.00, 7.25, 7.50, 7.75, 8.00,
	}

	return &HoursSelector{
		options:       options,
		selectedIndex: 3, // Default to 1 hour
		focused:       false,
	}
}

// Update handles input for the hours selector
func (hs *HoursSelector) Update(msg tea.KeyMsg) {
	if !hs.focused {
		return
	}

	switch msg.String() {
	// Numeric input for selecting hours to enable quick selection with fewer key presses
	case "1", "2", "3", "4", "5", "6", "7", "8":
		for i, option := range hs.options {
			if fmt.Sprintf("%d", int(option)) == msg.String() {
				hs.selectedIndex = i
				break
			}
		}
	case "left":
		if hs.selectedIndex > 0 {
			hs.selectedIndex--
		}
	case "right":
		if hs.selectedIndex < len(hs.options)-1 {
			hs.selectedIndex++
		}
	case "home":
		hs.selectedIndex = 0
	case "end":
		hs.selectedIndex = len(hs.options) - 1
	}
}

// Focus enables input handling
func (hs *HoursSelector) Focus() {
	hs.focused = true
}

// Blur disables input handling
func (hs *HoursSelector) Blur() {
	hs.focused = false
}

// SelectedHours returns the currently selected hours
func (hs *HoursSelector) SelectedHours() float64 {
	return hs.options[hs.selectedIndex]
}

// Render returns the horizontal hours selector view
func (hs *HoursSelector) Render() string {
	var items []string

	start := hs.selectedIndex - 3
	end := hs.selectedIndex + 4

	if start < 0 {
		start = 0
		end = 7
	}
	if end > len(hs.options) {
		end = len(hs.options)
		start = end - 7
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end && i < len(hs.options); i++ {
		hourStr := fmt.Sprintf("%.2f", hs.options[i])

		var style lipgloss.Style
		if i == hs.selectedIndex {
			if hs.focused {
				style = focusedStyle.Padding(0, 1)
			} else {
				style = fieldValueStyle.Bold(true).Foreground(themes.TokyoNight.Info).Padding(0, 1)
			}
		} else {
			style = fieldValueStyle.Padding(0, 1)
		}

		items = append(items, style.Render(hourStr))
	}

	leftArrow := ""
	rightArrow := ""
	if start > 0 {
		leftArrow = fieldValueStyle.Foreground(themes.TokyoNight.Info).Render("< ")
	} else {
		leftArrow = "  "
	}
	if end < len(hs.options) {
		rightArrow = fieldValueStyle.Foreground(themes.TokyoNight.Info).Render(" >")
	} else {
		rightArrow = "  "
	}

	content := leftArrow + strings.Join(items, " ") + rightArrow

	helpText := ""
	if hs.focused {
		helpText = "\n" + helpStyle.Render("Left/Right: select hours | Home/End: first/last")
	}

	return content + helpText
}
