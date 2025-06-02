package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FriendRequestWithRecipient struct {
	ID          bson.ObjectID       `bson:"_id,omitempty" json:"id"`
	SenderID    bson.ObjectID       `bson:"sender_id" json:"sender_id"`
	RecipientID bson.ObjectID       `bson:"recipient_id" json:"recipient_id"`
	Status      FriendRequestStatus `bson:"status" json:"status"`
	CreatedAt   time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updated_at"`
	Recipient   struct {
		ID          bson.ObjectID   `bson:"_id,omitempty" json:"id"`
		FullName    string          `bson:"full_name" json:"full_name"`
		Email       string          `bson:"email" json:"email"`
		Bio         string          `bson:"bio" json:"bio"`
		ProfilePic  string          `bson:"profile_pic" json:"profile_pic"`
		NativeLng   string          `bson:"native_lng" json:"native_lng"`
		LearningLng string          `bson:"learning_lng" json:"learning_lng"`
		Location    string          `bson:"location" json:"location"`
		IsOnboarded bool            `bson:"is_onboarded" json:"is_onboarded"`
		FriendIDs   []bson.ObjectID `bson:"friend_ids" json:"friend_ids"`
		CreatedAt   time.Time       `bson:"created_at" json:"created_at"`
		UpdatedAt   time.Time       `bson:"updated_at" json:"updated_at"`
	} `bson:"recipient" json:"recipient"`
}
