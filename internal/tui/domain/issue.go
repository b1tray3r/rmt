package domain

import "github.com/charmbracelet/bubbles/list"

type Issue struct {
	id          int
	link        string
	author      string
	title       string
	description string
}

var _ list.Item = (*Issue)(nil)

func NewIssue(id int, link, author, title, description string) *Issue {
	return &Issue{
		id:          id,
		link:        link,
		author:      author,
		title:       title,
		description: description,
	}
}

func (i *Issue) FilterValue() string {
	return i.link + " " + i.author + " " + i.title + " " + i.description
}

func (i *Issue) Title() string {
	return i.title
}

func (i *Issue) Description() string {
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
