package response

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FriendRequestWithSenderResponse struct {
	ID          bson.ObjectID `json:"id"`
	SenderID    bson.ObjectID `json:"sender_id"`
	RecipientID bson.ObjectID `json:"recipient_id"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Sender      struct {
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
	} `json:"sender"`
}
