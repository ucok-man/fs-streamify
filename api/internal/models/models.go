package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	// ErrRecordNotFound is returned when a movie record doesn't exist in database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
}

func NewModels(db *mongo.Database) Models {
	return Models{}
}
