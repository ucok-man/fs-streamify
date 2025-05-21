package response

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FriendRequestResponse struct {
	ID          bson.ObjectID `json:"id"`
	SenderID    bson.ObjectID `json:"sender_id"`
	RecipientID bson.ObjectID `json:"recipient_id"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
