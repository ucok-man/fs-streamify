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
	/* ------------------------- unique email ------------------------- */
	uniqeEmailIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	name, err := coll.Indexes().CreateOne(context.TODO(), uniqeEmailIdx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error creating unique email index")
	}
	logger.Info().Str("index_name", name).Msg("Success creating index")

	/* ------------------ text search index fullname ------------------ */
	name, err = coll.SearchIndexes().CreateOne(context.Background(), mongo.SearchIndexModel{
		Options: options.SearchIndexes().SetName("user_full_name_index"),
		Definition: bson.D{{Key: "mappings", Value: bson.D{
			{Key: "dynamic", Value: true},
		}}},
	})
	if err != nil {
		panic(err)
	}
	logger.Info().Str("index_name", name).Msg("Success creating index")

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.coll.FindOne(ctx, filter).Decode(&user)
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

func (m *UserModel) GetById(id bson.ObjectID) (*User, error) {
	filter := bson.D{{
		Key:   "_id",
		Value: id,
	}}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.coll.FindOne(ctx, filter).Decode(&user)
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
	current := time.Now()
	user.CreatedAt = current
	user.UpdatedAt = current

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
	current := time.Now()
	user.UpdatedAt = current

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

type RecommendedUserParam struct {
	CurrentUser *User
	Page        int64
	PageSize    int64
}

func (m *UserModel) Recommended(param RecommendedUserParam) ([]*User, Metadata, error) {
	// Validasi nilai Page dan PageSize
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.PageSize <= 0 {
		param.PageSize = 10
	}

	filter := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{"_id", bson.D{{"$ne", param.CurrentUser.ID}}}},
		bson.D{{"_id", bson.D{{"$nin", param.CurrentUser.FriendIDs}}}},
		bson.D{{"isOnboaded", true}},
	}}}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSkip((param.Page - 1) * param.PageSize)
	findOptions.SetLimit(param.PageSize)

	cursor, err := m.coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer cursor.Close(ctx)

	totalCount, err := m.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalCount, param.Page, param.PageSize)

	var users []*User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, Metadata{}, err
	}

	return users, metadata, nil
}

type MyFriendsParam struct {
	CurrentUser *User
	Search      string
	Page        int64
	PageSize    int64
}

func (m *UserModel) MyFriends(param MyFriendsParam) ([]*User, Metadata, error) {
	// pipeline := mongo.Pipeline{}

	// Step 1: Match hanya user yang merupakan teman dan sudah onboarded
	matchStage := bson.D{{Key: "$match", Value: bson.M{
		"_id":        bson.M{"$in": param.CurrentUser.FriendIDs},
		"isOnboaded": true,
	}}}

	// Step 2: Optional search by full_name using Atlas Search
	var searchStage bson.D
	if param.Search != "" {
		searchStage = bson.D{{Key: "$search", Value: bson.M{
			"index": "user_full_name_index",
			"text": bson.M{
				"query": param.Search,
				"path":  "full_name",
			},
		}}}
	}

	// Step 3: Pagination
	skipStage := bson.D{{Key: "$skip", Value: (param.Page - 1) * param.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: param.PageSize}}

	// $facet untuk split antara hasil dan count
	resultsPipeline := mongo.Pipeline{}
	if searchStage != nil {
		resultsPipeline = append(resultsPipeline, searchStage)
	}
	resultsPipeline = append(resultsPipeline, skipStage, limitStage)

	countPipeline := mongo.Pipeline{}
	if searchStage != nil {
		countPipeline = append(countPipeline, searchStage)
	}
	countPipeline = append(countPipeline, bson.D{{Key: "$count", Value: "total"}})

	facetStage := bson.D{{Key: "$facet", Value: bson.M{
		"data":  resultsPipeline,
		"count": countPipeline,
	}}}

	pipeline := mongo.Pipeline{matchStage, facetStage}

	// Execute
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer cursor.Close(ctx)

	var rawResult []struct {
		Data  []*User `bson:"data"`
		Count []struct {
			Total int64 `bson:"total"`
		} `bson:"count"`
	}

	if err := cursor.All(ctx, &rawResult); err != nil {
		return nil, Metadata{}, err
	}

	if len(rawResult) <= 0 {
		return nil, Metadata{}, err
	}
	result := rawResult[0]

	var totalCount int64
	if len(result.Count) > 0 {
		totalCount = result.Count[0].Total
	}

	metadata := calculateMetadata(totalCount, param.Page, param.PageSize)

	return result.Data, metadata, nil
}

func (m *UserModel) AddFriends(id bson.ObjectID, friendId bson.ObjectID) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$push", bson.D{{"friend_ids", friendId}}}}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.coll.UpdateOne(ctx, filter, update, options.UpdateOne())
	if err != nil {
		return err
	}
	return nil
}
