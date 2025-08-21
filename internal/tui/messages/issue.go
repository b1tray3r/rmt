package messages

import (
	"github.com/b1tray3r/rmt/internal/tui/domain"
)

// IssueSelectedMsg is sent when an issue is selected from the list.
type IssueSelectedMsg struct {
	Issue *domain.Issue
}
