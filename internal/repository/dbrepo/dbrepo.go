package dbrepo

import (
	"context"
	"database/sql"
	"html"
	"time"

	"github.com/elkcityhazard/am-form/internal/config"
	"github.com/elkcityhazard/am-form/internal/models"
)

// databaseRepo holds the database connection and receiever functions for mocking later

var databaseRepo *DBRepo

// DBRepo is a type that accesses the app config so we can use the database
type DBRepo struct {
	App *config.AppConfig
}

// NewDBRepo creates a new repository
func NewDBRepo(a *config.AppConfig) *DBRepo {
	repo := &DBRepo{
		App: a,
	}
	databaseRepo = repo
	return repo
}

// GetDBRepo returns the database repo
func GetDBRepo() *DBRepo {
	return databaseRepo
}

// InsertMessage looks up the message email in the database, creates a user if it doesn't exist, and inserts the message
func (d *DBRepo) InsertMessage(msg *models.Message) error {

	ctx, cancel := context.WithTimeout(d.App.Context, 10*time.Second)

	defer cancel()

	d.App.WG.Add(1)
	errorChan := make(chan error, 1)
	go func() {

		defer d.App.WG.Done()
		defer close(errorChan)

		tx, err := d.App.DB.BeginTx(ctx, nil)

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		err = tx.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", msg.Email).Scan(&msg.UserID)

		if err != nil {

			if err == sql.ErrNoRows {
				newUser := &models.User{
					Email:     msg.Email,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Version:   1,
				}

				userID, err := d.InsertUser(newUser)

				if err != nil {
					errorChan <- err
					tx.Rollback()
					return
				}

				msg.UserID = userID

			}
		}

		result, err := tx.ExecContext(ctx, "INSERT INTO messages (user_id, email, message) VALUES (?, ?, ?)", msg.UserID, msg.Email, html.EscapeString(msg.Message))

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		messageID, err := result.LastInsertId()

		if err != nil {
			errorChan <- err
			tx.Rollback()
			return
		}

		msg.ID = messageID
		tx.Commit()
		errorChan <- nil

	}()

	if err := <-errorChan; err != nil {
		return err
	} else {
		return nil
	}
}

/**********/
/* USERS */
/**********/

func (d *DBRepo) InsertUser(user *models.User) (int64, error) {

	ctx, cancel := context.WithTimeout(d.App.Context, 10*time.Second)

	defer cancel()

	d.App.WG.Add(1)

	idChan, errorChan := make(chan int64, 1), make(chan error, 1)

	go func() {

		defer d.App.WG.Done()
		defer close(idChan)
		defer close(errorChan)
		d.App.Mutex.Lock()

		tx, err := d.App.DB.BeginTx(ctx, nil)

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return

		}

		result, err := tx.ExecContext(ctx, "INSERT INTO users (email, created_at, updated_at, version) VALUES (?, ?, ?, ?)", user.Email, user.CreatedAt, user.UpdatedAt, user.Version)

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		id, err := result.LastInsertId()

		if err != nil {
			tx.Rollback()
			errorChan <- err
			return
		}

		tx.Commit()
		idChan <- id
		errorChan <- nil

		defer d.App.Mutex.Unlock()

	}()

	if err := <-errorChan; err != nil {
		return 0, err
	} else {
		return <-idChan, nil
	}

}
