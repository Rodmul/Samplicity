package model

type Sample struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Author   string `db:"author"`
	AuthorID int    `db:"author_id"`
	Path     string `db:"path"`
	Type     string `db:"type"`
}
