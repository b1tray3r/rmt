package domain

import (
	"strings"

	"github.com/b1tray3r/rmt/internal/redmine"
	"github.com/b1tray3r/rmt/internal/redmine/models"
	"github.com/charmbracelet/bubbles/list"
)

type Project struct {
	id   int
	name string
}

func (p *Project) ID() int {
	return p.id
}

func (p *Project) Name() string {
	return p.name
}

type Issue struct {
	id          int
	link        string
	author      string
	title       string
	project     *Project
	description string
}

// Ensure Issue implements the list.Item interface so it can be used in the list view
var _ list.Item = (*Issue)(nil)

func NewIssue(id int, link, author, title, description string, project *Project) *Issue {
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
	if len(runes) > 55 {
		return string(runes[:55])
	}
	return i.title
}

func (i *Issue) FullTitle() string {
	return i.title
}

// Description returns the issue description as a single line, replacing "\n*" with ", " and all newlines with spaces, and truncates to 75 runes if necessary.
func (i *Issue) Description() string {
	processed := strings.ReplaceAll(i.description, "\n*", ", ")
	processed = strings.ReplaceAll(processed, "\n", " ")
	processed = strings.Join(strings.Fields(processed), " ")
	runes := []rune(processed)
	if len(runes) > 55 {
		return string(runes[:55])
	}
	if len(runes) == 0 {
		return "-- no description --"
	}
	return processed
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

func (i *Issue) Project() *Project {
	return i.project
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
	for _, issue := range issues.Results {
		title := s.cleanTitle(issue.Title)

		ni := NewIssue(
			issue.ID,
			issue.URL,
			"none",
			title,
			issue.Description,
			nil,
		)
		result = append(result, ni)
	}

	return result, nil
}
