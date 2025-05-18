package models

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
	ID          bson.ObjectID   `bson:"_id,omitempty"`
	FullName    string          `bson:"full_name"`
	Email       string          `bson:"email"`
	Password    password        `bson:"inline"`
	Bio         string          `bson:"bio"`
	ProfilePic  string          `bson:"profile_pic"`
	NativeLng   string          `bson:"native_lng"`
	LearningLng string          `bson:"learning_lng"`
	Location    string          `bson:"location"`
	IsOnboarded bool            `bson:"is_onboarded"` // TODO:default false
	FriendIDs   []bson.ObjectID `bson:"friend_ids"`
	CreatedAt   time.Time       `bson:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at"`
}

type UserModel struct {
	logger zerolog.Logger
	coll   *mongo.Collection
}

func NewUserModel(coll *mongo.Collection, logger zerolog.Logger) *UserModel {
	uniqeEmailIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	idxname, err := coll.Indexes().CreateOne(context.TODO(), uniqeEmailIdx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error creating unique email index")
	}
	logger.Info().Str("index_name", idxname).Msg("Success creating index")

	return &UserModel{
		coll:   coll,
		logger: logger,
	}
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	filter := bson.D{{
		Key:   "email",
		Value: email,
	}}

	var user User
	err := m.coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Insert(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	current := time.Now()
	user.CreatedAt = current
	user.UpdatedAt = current

	result, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		switch {
		case mongo.IsDuplicateKeyError(err):
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	userID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, errors.New("ID is not ObjectID, you should let mongo manage the ID")
	}

	user.ID = userID
	return user, nil
}
