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
	// Define the aggregation pipeline
	lookupSenderStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "sender_id",
		"foreignField": "_id",
		"as":           "sender",
	}}}

	var searchSenderStage bson.D
	if param.SearchSender != "" {
		searchSenderStage = bson.D{{Key: "$search", Value: bson.D{
			{Key: "index", Value: "user_full_name_index"},
			{Key: "$text", Value: bson.M{
				"query": param.SearchSender,
				"path":  "sender.full_name",
			}},
		}}}
	}

	var matchStage bson.D
	if param.Status == "All" {
		matchStage = bson.D{{Key: "$match", Value: bson.M{
			"recipient_id": param.CurrentUserId,
		}}}
	} else {
		matchStage = bson.D{{Key: "$match", Value: bson.M{
			"recipient_id": param.CurrentUserId,
			"status":       FriendRequestStatus(param.Status),
		}}}
	}

	skipStage := bson.D{{Key: "$skip", Value: (param.Page - 1) * param.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: param.PageSize}}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := m.coll.Aggregate(ctx, mongo.Pipeline{lookupSenderStage, searchSenderStage, matchStage, skipStage, limitStage})
	if err != nil {
		return nil, err
	}
	var results []*FriendRequestWithSender
	if err = cursor.All(context.Background(), &results); err != nil {
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

func (m *FriendRequestModel) GetAllSendFriendRequest(param GetAllSendFriendRequestParam) ([]*FriendRequestWithSender, error) {
	// Define the aggregation pipeline
	lookupRecipientStage := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "users",
		"localField":   "recipient_id",
		"foreignField": "_id",
		"as":           "recipient",
	}}}

	var searchRecipient bson.D
	if param.SearchRecipient != "" {
		searchRecipient = bson.D{{Key: "$search", Value: bson.D{
			{Key: "index", Value: "user_full_name_index"},
			{Key: "$text", Value: bson.M{
				"query": param.SearchRecipient,
				"path":  "recipient.full_name",
			}},
		}}}
	}

	var matchStage bson.D
	if param.Status == "All" {
		matchStage = bson.D{{Key: "$match", Value: bson.M{
			"recipient_id": param.CurrentUserId,
		}}}
	} else {
		matchStage = bson.D{{Key: "$match", Value: bson.M{
			"recipient_id": param.CurrentUserId,
			"status":       FriendRequestStatus(param.Status),
		}}}
	}

	skipStage := bson.D{{Key: "$skip", Value: (param.Page - 1) * param.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: param.PageSize}}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := m.coll.Aggregate(ctx, mongo.Pipeline{lookupRecipientStage, searchRecipient, matchStage, skipStage, limitStage})
	if err != nil {
		return nil, err
	}
	var results []*FriendRequestWithSender
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
