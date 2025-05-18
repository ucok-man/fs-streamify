package models

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("error duplicate email")
)

type Models struct {
	Logger *zerolog.Logger
	User   *UserModel
}

func NewModels(db *mongo.Database, logger *zerolog.Logger) Models {
	return Models{
		User: NewUserModel(
			db.Collection("users"),
			logger.With().Str("context", "user_model_service").Logger(),
		),
	}
}
