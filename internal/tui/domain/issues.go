package domain

import (
	"fmt"
	"strconv"
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

// IssueBaseURLProvider defines an interface for retrieving the base URL for issues.
type IssueBaseURLProvider interface {
	GetBaseURL() string
}

// IssueGetter defines an interface for retrieving a single issue by its ID.
type IssueGetter interface {
	GetIssue(id int) (*Issue, error)
}

// IssueSearcher defines an interface for searching issues by a query string.
type IssueSearcher interface {
	Search(query string) ([]*Issue, error)
}

type TimeEntryCreator interface {
	CreateTimeEntry(params models.CreateTimeEntryParams) (*models.TimeEntry, error)
}

type ProjectActivityGetter interface {
	GetProjectActivities(projectID int) (map[int]string, error)
}

// IssueRepository composes the one-purpose interfaces for issue operations.
type IssueRepository interface {
	IssueBaseURLProvider
	IssueGetter
	IssueSearcher
	TimeEntryCreator
	ProjectActivityGetter
}

type RedmineIssueRepository struct {
	client redmine.RedmineAPI
}

func NewRedmineIssueRepository(client redmine.RedmineAPI) *RedmineIssueRepository {
	return &RedmineIssueRepository{
		client: client,
	}
}

func (s *RedmineIssueRepository) GetBaseURL() string {
	return s.client.GetBaseURL()
}

func (s *RedmineIssueRepository) cleanTitle(title string) string {
	parts := strings.Split(title, ": ")
	if len(parts) > 1 {
		return strings.Join(parts[1:], ": ")
	}

	return title
}

func (s *RedmineIssueRepository) GetProjectActivities(projectID int) (map[int]string, error) {
	project, err := s.client.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	result := make(map[int]string)
	for _, activity := range project.TimeEntryActivities {
		result[activity.ID] = activity.Name
	}

	return result, nil
}

func (s *RedmineIssueRepository) CreateTimeEntry(params models.CreateTimeEntryParams) (*models.TimeEntry, error) {
	return s.client.CreateTimeEntry(params)
}

func (s *RedmineIssueRepository) GetIssue(id int) (*Issue, error) {
	issue, err := s.client.GetIssue(id)
	if err != nil {
		return nil, err
	}

	return NewIssue(
		issue.ID,
		fmt.Sprintf("%s/issues/%d", s.GetBaseURL(), issue.ID),
		issue.Author.Name,
		s.cleanTitle(issue.Subject),
		issue.Description,
		&Project{
			id:   issue.Project.ID,
			name: issue.Project.Name,
		},
	), nil
}

func (s *RedmineIssueRepository) Search(query string) ([]*Issue, error) {
	if strings.HasPrefix(query, "#") {
		idStr := strings.TrimPrefix(query, "#")
		if id, err := strconv.Atoi(idStr); err == nil {
			issue, err := s.client.GetIssue(id)
			if err != nil {
				return nil, err
			}

			ni := NewIssue(
				issue.ID,
				fmt.Sprintf("%s/issues/%d", s.GetBaseURL(), issue.ID),
				issue.Author.Name,
				s.cleanTitle(issue.Subject),
				issue.Description,
				&Project{
					id:   issue.Project.ID,
					name: issue.Project.Name,
				},
			)
			return []*Issue{ni}, nil
		}
	}

	// Regular search for non-ID queries
	issues, err := s.client.Search(models.SearchParams{
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	var result []*Issue
	for _, issue := range issues.Results {
		i, err := s.client.GetIssue(issue.ID)
		if err != nil {
			return nil, err
		}

		ni := NewIssue(
			issue.ID,
			issue.URL,
			i.Author.Name,
			s.cleanTitle(i.Subject),
			i.Description,
			&Project{
				id:   i.Project.ID,
				name: i.Project.Name,
			},
		)
		result = append(result, ni)
	}

	return result, nil
}
