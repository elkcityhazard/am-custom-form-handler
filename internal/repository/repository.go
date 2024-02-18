package repository

import "github.com/elkcityhazard/am-form/internal/models"

type AppRepository interface {
	InsertMessage(msg *models.Message) error
	// GetMessages(offset, limit int) ([]*models.Message, error)
	// GetMessage(id int64) (*models.Message, error)
	// UpdateMessage(msg *models.Message) error
	// DeleteMessage(id int64) error

	InsertUser(user *models.User) (id int64, err error)
	// GetUsers(offset, limit int) ([]*models.User, error)
	// GetUser(id int64) (*models.User, error)
	// UpdateUser(user *models.User) error
	// DeleteUser(id int64) error
}
