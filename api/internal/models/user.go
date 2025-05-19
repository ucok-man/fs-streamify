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
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m *UserModel) GetById(id string) (*User, error) {
	idbson, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{
		Key:   "_id",
		Value: idbson,
	}}

	var user User
	err = m.coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
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

func (m *UserModel) Update(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	current := time.Now()
	user.UpdatedAt = current

	update := bson.D{
		{Key: "$set", Value: bson.D{{"full_name", user.FullName}}},
		{Key: "$set", Value: bson.D{{"bio", user.Bio}}},
		{Key: "$set", Value: bson.D{{"profile_pic", user.ProfilePic}}},
		{Key: "$set", Value: bson.D{{"native_lng", user.NativeLng}}},
		{Key: "$set", Value: bson.D{{"learning_lng", user.LearningLng}}},
		{Key: "$set", Value: bson.D{{"location", user.Location}}},
		{Key: "$set", Value: bson.D{{"is_onboarded", user.IsOnboarded}}},
		{Key: "$set", Value: bson.D{{"updated_at", user.UpdatedAt}}},
		{Key: "$set", Value: bson.D{{"friend_ids", user.FriendIDs}}},
	}

	_, err := m.coll.UpdateByID(ctx, user.ID, update)
	if err != nil {
		return nil, err
	}

	return user, nil
}
