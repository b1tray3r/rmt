package messages

import (
	"github.com/b1tray3r/rmt/internal/tui/domain"
)

// IssueSelectedMsg is sent when an issue is selected from the list.
type IssueSelectedMsg struct {
	Issue *domain.Issue
}

// ReturnToIssueMsg indicates the user wants to return to issue view
// Parent applications should handle this message to navigate back to issue view
type ReturnToIssueMsg struct{}
