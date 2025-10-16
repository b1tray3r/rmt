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
	SearchWithFilter(query string) ([]*Issue, error) // New method for Redmine issue queries
}

type TimeEntryCreator interface {
	CreateTimeEntry(params models.CreateTimeEntryParams) (*models.TimeEntry, error)
}

type ProjectActivityGetter interface {
	GetProjectActivities(projectID int, activityPatterns []string) (map[int]string, error)
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

func (s *RedmineIssueRepository) GetProjectActivities(projectID int, activityPatterns []string) (map[int]string, error) {
	project, err := s.client.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	result := make(map[int]string)
	for _, activity := range project.TimeEntryActivities {
		if s.matchesActivityPatterns(activity.Name, activityPatterns) {
			result[activity.ID] = activity.Name
		}
	}

	return result, nil
}

// matchesActivityPatterns checks if an activity name matches any of the given patterns.
// matchesActivityPatterns returns true if no patterns are provided (matches all).
func (s *RedmineIssueRepository) matchesActivityPatterns(activityName string, activityPatterns []string) bool {
	if len(activityPatterns) == 0 {
		return true
	}

	for _, pattern := range activityPatterns {
		if strings.HasPrefix(activityName, pattern) {
			return true
		}
	}

	return false
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

// SearchWithFilter searches issues using Redmine issue query format (actual Redmine format)
func (s *RedmineIssueRepository) SearchWithFilter(query string) ([]*Issue, error) {
	// Use the raw query string directly with Redmine's issues API
	results, err := s.client.SearchIssuesRaw(query)
	if err != nil {
		return nil, err
	}

	var issueList []*Issue
	for _, issue := range results.Issues {
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
		issueList = append(issueList, ni)
	}

	return issueList, nil
} // parseRedmineQuery parses a Redmine query string into an IssueFilter
func (s *RedmineIssueRepository) parseRedmineQuery(query string) models.IssueFilter {
	filter := models.IssueFilter{
		Limit:        100,
		CustomFields: make(map[int]string),
	}

	// Split query by & to get individual parameters
	parts := strings.Split(query, "&")
	for _, part := range parts {
		if strings.Contains(part, "=") {
			keyValue := strings.SplitN(part, "=", 2)
			key := keyValue[0]
			value := keyValue[1]

			// Handle different parameter types
			switch {
			case strings.HasPrefix(key, "cf_"):
				// Custom field parameter: cf_10=*this week*
				if fieldIDStr := strings.TrimPrefix(key, "cf_"); fieldIDStr != "" {
					if fieldID, err := strconv.Atoi(fieldIDStr); err == nil {
						filter.CustomFields[fieldID] = value
					}
				}
			case key == "assigned_to_id":
				filter.AssignedTo = value
			case key == "status_id":
				if statusID, err := strconv.Atoi(value); err == nil {
					filter.StatusID = []int{statusID}
				}
			case key == "tracker_id":
				if trackerID, err := strconv.Atoi(value); err == nil {
					filter.TrackerID = []int{trackerID}
				}
			case key == "project_id":
				if projectID, err := strconv.Atoi(value); err == nil {
					filter.ProjectID = []int{projectID}
				}
			case key == "subject":
				filter.Subject = value
			case key == "limit":
				if limit, err := strconv.Atoi(value); err == nil {
					filter.Limit = limit
				}
			case key == "offset":
				if offset, err := strconv.Atoi(value); err == nil {
					filter.Offset = offset
				}
			}
		}
	}

	return filter
}
