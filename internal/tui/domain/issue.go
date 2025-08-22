package domain

import (
	"fmt"
	"strings"

	"github.com/b1tray3r/rmt/internal/redmine"
	"github.com/b1tray3r/rmt/internal/redmine/models"
	"github.com/charmbracelet/bubbles/list"
)

type Issue struct {
	id          int
	link        string
	author      string
	title       string
	project     string
	description string
}

// Ensure Issue implements the list.Item interface so it can be used in the list view
var _ list.Item = (*Issue)(nil)

func NewIssue(id int, link, author, title, project, description string) *Issue {
	return &Issue{
		id:          id,
		link:        link,
		author:      author,
		title:       title,
		project:     project,
		description: description,
	}
}

func (i *Issue) FilterValue() string {
	return i.link + " " + i.author + " " + i.title
}

func (i *Issue) Title() string {
	runes := []rune(i.title)
	if len(runes) > 75 {
		return string(runes[:75])
	}
	return i.title
}

func (i *Issue) FullTitle() string {
	return i.title
}

func (i *Issue) Description() string {
	runes := []rune(i.description)
	for idx, r := range runes {
		if r == '\n' {
			runes[idx] = ' '
		}
	}
	processed := strings.Join(strings.Fields(string(runes)), " ")
	runes = []rune(processed)
	if len(runes) > 75 {

		return string(runes[:75])
	}
	return i.description
}

func (i *Issue) FullDescription() string {
	return i.description
}

func (i *Issue) Author() string {
	return i.author
}

func (i *Issue) Link() string {
	return i.link
}

func (i *Issue) ID() int {
	return i.id
}

type IssueService struct {
	client redmine.RedmineAPI
}

func NewIssueService(client redmine.RedmineAPI) *IssueService {
	return &IssueService{
		client: client,
	}
}

func (s *IssueService) cleanTitle(title string) string {
	parts := strings.Split(title, ": ")
	if len(parts) > 1 {
		return strings.Join(parts[1:], ": ")
	}

	return title
}

func (s *IssueService) Search(query string) ([]*Issue, error) {
	issues, err := s.client.Search(models.SearchParams{
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	var result []*Issue
	for i, issue := range issues.Results {
		title := s.cleanTitle(issue.Title)

		ni := NewIssue(
			issue.ID,
			issue.URL,
			"none",
			fmt.Sprintf("%d: %s", i, title),
			issue.Project,
			issue.Description,
		)
		result = append(result, ni)
	}

	return result, nil
}
