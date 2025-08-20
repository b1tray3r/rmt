package domain

import "github.com/charmbracelet/bubbles/list"

type Issue struct {
	ID      int
	Link    string
	Author  string
	Subject string
	Content string
}

var _ list.Item = (*Issue)(nil)

func (i *Issue) FilterValue() string {
	return i.Link + " " + i.Author + " " + i.Subject + " " + i.Content
}

func (i *Issue) Title() string {
	return i.Subject
}

func (i *Issue) Description() string {
	return i.Content
}
