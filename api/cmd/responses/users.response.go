package response

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UsersResponse []struct {
	ID          bson.ObjectID   `json:"id"`
	FullName    string          `json:"full_name"`
	Email       string          `json:"email"`
	Bio         string          `json:"bio"`
	ProfilePic  string          `json:"profile_pic"`
	NativeLng   string          `json:"native_lng"`
	LearningLng string          `json:"learning_lng"`
	Location    string          `json:"location"`
	IsOnboarded bool            `json:"is_onboarded"`
	FriendIDs   []bson.ObjectID `json:"friend_ids"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
