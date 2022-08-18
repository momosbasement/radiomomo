package gui

import (
	model "github.com/momosbasement/radiomomo/pkg/models"
)

type Page struct {
	Title   string
	Content string
	Tracks  []model.Track
}
