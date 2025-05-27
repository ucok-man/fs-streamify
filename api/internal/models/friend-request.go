package models

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FriendRequestStatus = string

const (
	FriendRequestStatusPending  FriendRequestStatus = "Pending"
	FriendRequestStatusAccepted FriendRequestStatus = "Accepted"
)

type FriendRequest struct {
	ID          bson.ObjectID       `bson:"_id,omitempty"`
	SenderID    bson.ObjectID       `bson:"sender_id"`
	RecipientID bson.ObjectID       `bson:"recipient_id"`
	Status      FriendRequestStatus `bson:"status"`
	CreatedAt   time.Time           `bson:"created_at"`
	UpdatedAt   time.Time           `bson:"updated_at"`
}

type FriendRequestModel struct {
	logger zerolog.Logger
	coll   *mongo.Collection
}

func NewFriendRequestModel(coll *mongo.Collection, logger zerolog.Logger) *FriendRequestModel {
	return &FriendRequestModel{
		coll:   coll,
		logger: logger,
	}
}

func (m *FriendRequestModel) GetById(id bson.ObjectID) (*FriendRequest, error) {
	filter := bson.D{{
		Key:   "_id",
		Value: id,
	}}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var friendRequest FriendRequest
	err := m.coll.FindOne(ctx, filter).Decode(&friendRequest)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &friendRequest, nil
}

func (m *FriendRequestModel) CheckExisting(senderId, receipentId bson.ObjectID) (bool, error) {
	filter := bson.D{{
		"$or", bson.A{
			bson.D{
				{"sender_id", senderId},
				{"recipient_id", receipentId},
			},
			bson.D{
				{"sender_id", receipentId},
				{"recipient_id", senderId},
			},
		},
	}}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var friendRequest *FriendRequest
	if err := m.coll.FindOne(ctx, filter).Decode(friendRequest); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return false, nil
		default:
			return false, err
		}
	}

	return friendRequest == nil, nil
}

func (m *FriendRequestModel) Create(friendRequest *FriendRequest) (*FriendRequest, error) {
	friendRequest.CreatedAt = time.Now()
	friendRequest.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.coll.InsertOne(ctx, friendRequest)
	if err != nil {
		return nil, err
	}

	idrecord, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, errors.New("ID is not ObjectID, you should let mongo manage the ID")
	}

	friendRequest.ID = idrecord
	return friendRequest, nil
}

func (m *FriendRequestModel) Update(friendRequest *FriendRequest) (*FriendRequest, error) {
	current := time.Now()
	friendRequest.UpdatedAt = current

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: friendRequest.Status},
		}},
	}

	_, err := m.coll.UpdateByID(ctx, friendRequest.ID, update)
	if err != nil {
		return nil, err
	}

	return friendRequest, nil
}

type GetAllFromFriendRequestParam struct {
	CurrentUserId bson.ObjectID
	Status        string
	Page          int64
	PageSize      int64
	SearchSender  string
}

func (m *FriendRequestModel) GetAllFromFriendRequest(param GetAllFromFriendRequestParam) ([]*FriendRequestWithSender, error) {
	pipeline := mongo.Pipeline{}

	// Step 1: Match by recipient_id and status
	matchStage := bson.D{
		{Key: "$match", Value: bson.M{
			"recipient_id": param.CurrentUserId,
		}},
	}
	if param.Status != "All" {
		matchStage[0].Value.(bson.M)["status"] = FriendRequestStatus(param.Status)
	}
	pipeline = append(pipeline, matchStage)

	// Step 2: Lookup sender user data
	lookupStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "sender_id",
		"foreignField": "_id",
		"as":           "sender",
	}}}
	pipeline = append(pipeline, lookupStage)

	// Step 3: Unwind sender array
	unwindStage := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$sender",
		"preserveNullAndEmptyArrays": false,
	}}}
	pipeline = append(pipeline, unwindStage)

	// Step 4: Optional search by sender.full_name using Atlas Search
	if param.SearchSender != "" {
		searchStage := bson.D{{Key: "$search", Value: bson.M{
			"index": "user_full_name_index", // dynamic index
			"text": bson.M{
				"query": param.SearchSender,
				"path":  "sender.full_name",
			},
		}}}
		pipeline = append(pipeline, searchStage)
	}

	// Step 5–6: Pagination
	skipStage := bson.D{{Key: "$skip", Value: (param.Page - 1) * param.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: param.PageSize}}

	pipeline = append(pipeline, skipStage, limitStage)

	// Execute pipeline
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := m.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*FriendRequestWithSender
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

type GetAllSendFriendRequestParam struct {
	CurrentUserId   bson.ObjectID
	Status          string
	Page            int64
	PageSize        int64
	SearchRecipient string
}

func (m *FriendRequestModel) GetAllSendFriendRequest(param GetAllSendFriendRequestParam) ([]*FriendRequestWithRecipient, error) {
	pipeline := mongo.Pipeline{}

	// Step 1: Match by sender_id and status
	matchStage := bson.D{
		{Key: "$match", Value: bson.M{
			"sender_id": param.CurrentUserId,
		}},
	}
	if param.Status != "All" {
		matchStage[0].Value.(bson.M)["status"] = FriendRequestStatus(param.Status)
	}
	pipeline = append(pipeline, matchStage)

	// Step 2: Lookup recipient user data
	lookupStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "recipient_id",
		"foreignField": "_id",
		"as":           "recipient",
	}}}
	pipeline = append(pipeline, lookupStage)

	// Step 3: Unwind recipient array
	unwindStage := bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$recipient",
		"preserveNullAndEmptyArrays": false,
	}}}
	pipeline = append(pipeline, unwindStage)

	// Step 4: Optional search by recipient.full_name using Atlas Search
	if param.SearchRecipient != "" {
		searchStage := bson.D{{Key: "$search", Value: bson.M{
			"index": "user_full_name_index",
			"text": bson.M{
				"query": param.SearchRecipient,
				"path":  "recipient.full_name",
			},
		}}}
		pipeline = append(pipeline, searchStage)
	}

	// Step 5–6: Pagination
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.PageSize <= 0 {
		param.PageSize = 10
	}
	skipStage := bson.D{{Key: "$skip", Value: (param.Page - 1) * param.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: param.PageSize}}

	pipeline = append(pipeline, skipStage, limitStage)

	// Execute pipeline
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := m.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*FriendRequestWithRecipient
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
