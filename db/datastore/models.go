// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"time"
)

type Account struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at"`
}