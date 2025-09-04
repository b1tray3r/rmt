package domain

import "github.com/charmbracelet/bubbles/list"

type Favorite struct {
	id     int
	name   string
	config string
}

// FilterValue implements list.Item.
func (f *Favorite) FilterValue() string {
	return f.name
}

var _ list.Item = (*Favorite)(nil)

func NewFavorite(id int, name, config string) *Favorite {
	return &Favorite{
		id:     id,
		name:   name,
		config: config,
	}
}

// ID returns the favorite's ID
func (f *Favorite) ID() int {
	return f.id
}

// Name returns the favorite's name
func (f *Favorite) Name() string {
	return f.name
}

// Config returns the favorite's config
func (f *Favorite) Config() string {
	return f.config
}
