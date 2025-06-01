package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserWithFriendRequest struct {
	ID                 bson.ObjectID    `bson:"_id,omitempty" json:"id"`
	FullName           string           `bson:"full_name" json:"full_name"`
	Email              string           `bson:"email" json:"email"`
	Password           password         `bson:"inline" json:"-"`
	Bio                string           `bson:"bio" json:"bio"`
	ProfilePic         string           `bson:"profile_pic" json:"profile_pic"`
	NativeLng          string           `bson:"native_lng" json:"native_lng"`
	LearningLng        string           `bson:"learning_lng" json:"learning_lng"`
	Location           string           `bson:"location" json:"location"`
	IsOnboarded        bool             `bson:"is_onboarded" json:"is_onboarded"`
	FriendIDs          []bson.ObjectID  `bson:"friend_ids" json:"friend_ids"`
	CreatedAt          time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt          time.Time        `bson:"updated_at" json:"updated_at"`
	SendtFriendRequest []*FriendRequest `bson:"sent_friend_request" json:"sent_friend_request"`
	FromFriendRequest  []*FriendRequest `bson:"from_friend_request" json:"from_friend_request"`
}
