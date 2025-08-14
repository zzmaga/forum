package view

import "forum/architecture/models"

type Page struct {
	User  *models.User
	Users []*models.User
	// Post *models.Post
	// Categories []*models.Category

	// Comments []models.Comment
	Error   error // Notification Error
	Warn    error // Notfication Warning
	Info    error // Notification Info
	Success error // Notification Success
}
