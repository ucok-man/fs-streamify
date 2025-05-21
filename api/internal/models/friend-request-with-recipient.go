package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FriendRequestWithRecipient struct {
	ID          bson.ObjectID       `bson:"_id,omitempty"`
	SenderID    bson.ObjectID       `bson:"sender_id"`
	RecipientID bson.ObjectID       `bson:"recipient_id"`
	Status      FriendRequestStatus `bson:"status"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
	Recipient   struct {
		ID          bson.ObjectID   `bson:"_id,omitempty"`
		FullName    string          `bson:"full_name"`
		Email       string          `bson:"email"`
		Bio         string          `bson:"bio"`
		ProfilePic  string          `bson:"profile_pic"`
		NativeLng   string          `bson:"native_lng"`
		LearningLng string          `bson:"learning_lng"`
		Location    string          `bson:"location"`
		IsOnboarded bool            `bson:"is_onboarded"` // TODO:default false
		FriendIDs   []bson.ObjectID `bson:"friend_ids"`
		CreatedAt   time.Time       `bson:"created_at"`
		UpdatedAt   time.Time       `bson:"updated_at"`
	} `bson:"recipient"`
}
