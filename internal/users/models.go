package users

import "time"

type User struct {
	ID         string    `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	ExternalID int       `db:"external_id" json:"external_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
