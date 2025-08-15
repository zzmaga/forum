package view

import "forum/architecture/models"

type Page struct {
	User       *models.User
	Users      []*models.User
	Post       *models.Post
	Posts      []*models.Post
	Categories []*models.Category

	// Comments           []models.Comment
	Error   error // Error - Notification Error
	Warn    error // Warn - Notification Warning
	Info    error // Info - Notification Info
	Success error // Success - Notification Success
}

type View struct {
	templatesDir string
}
