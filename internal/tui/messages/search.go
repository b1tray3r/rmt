package messages

import "github.com/b1tray3r/rmt/internal/domain"

type SearchSubmittedMsg struct {
	Query string
}

type SearchCompletedMsg struct {
	Query   string
	Results []domain.Issue
	Error   error
}
