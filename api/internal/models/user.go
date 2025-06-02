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
	ID          bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	FullName    string          `bson:"full_name" json:"full_name"`
	Email       string          `bson:"email" json:"email"`
	Password    password        `bson:"inline" json:"-"`
	Bio         string          `bson:"bio" json:"bio"`
	ProfilePic  string          `bson:"profile_pic" json:"profile_pic"`
	NativeLng   string          `bson:"native_lng" json:"native_lng"`
	LearningLng string          `bson:"learning_lng" json:"learning_lng"`
	Location    string          `bson:"location" json:"location"`
	IsOnboarded bool            `bson:"is_onboarded" json:"is_onboarded"`
	FriendIDs   []bson.ObjectID `bson:"friend_ids" json:"friend_ids"`
	CreatedAt   time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at" json:"updated_at"`
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
	user.FriendIDs = []bson.ObjectID{}

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

func (m *UserModel) Recommended(param RecommendedUserParam) ([]*UserWithFriendRequest, Metadata, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{"$and", bson.A{
			bson.D{{"_id", bson.D{{"$ne", param.CurrentUser.ID}}}},
			bson.D{{"_id", bson.D{{"$nin", param.CurrentUser.FriendIDs}}}},
			bson.D{{"is_onboarded", true}},
		}},
	}}}

	lookupStageSentFriendRequest := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "friend_request"},
		{Key: "let", Value: bson.D{
			{Key: "recipientId", Value: "$_id"},
		}},
		{Key: "pipeline", Value: bson.A{
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "$expr", Value: bson.D{
					{Key: "$and", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$sender_id", param.CurrentUser.ID}}},
						bson.D{{Key: "$eq", Value: bson.A{"$recipient_id", "$$recipientId"}}},
					}},
				}},
			}}},
		}},
		{Key: "as", Value: "sent_friend_request"},
	}}}

	lookupStageFromFriendRequest := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "friend_request"},
		{Key: "let", Value: bson.D{
			{Key: "recipientId", Value: "$_id"},
		}},
		{Key: "pipeline", Value: bson.A{
			bson.D{{Key: "$match", Value: bson.D{
				{Key: "$expr", Value: bson.D{
					{Key: "$and", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$sender_id", "$$recipientId"}}},
						bson.D{{Key: "$eq", Value: bson.A{"$recipient_id", param.CurrentUser.ID}}},
					}},
				}},
			}}},
		}},
		{Key: "as", Value: "from_friend_request"},
	}}}

	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "has_friend_request", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$sent_friend_request"}}, 0}}},
				bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$from_friend_request"}}, 0}}},
			}},
		}},
	}}}

	sortStage := bson.D{{Key: "$sort", Value: bson.D{
		{Key: "has_friend_request", Value: 1}, // sort if !has_friend_request appear first
	}}}

	skipStage := bson.D{{Key: "$skip", Value: (param.Page - 1) * param.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: param.PageSize}}

	// Result and count in a single aggregation
	resultsPipeline := mongo.Pipeline{
		matchStage,
		lookupStageSentFriendRequest,
		lookupStageFromFriendRequest,
		addFieldsStage,
		sortStage,
		skipStage,
		limitStage,
	}

	countPipeline := mongo.Pipeline{matchStage, bson.D{{Key: "$count", Value: "total"}}}

	facetStage := bson.D{{Key: "$facet", Value: bson.M{
		"data":  resultsPipeline,
		"count": countPipeline,
	}}}

	pipeline := mongo.Pipeline{facetStage}

	// Execute
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer cursor.Close(ctx)

	var rawResult []struct {
		Data  []*UserWithFriendRequest `bson:"data"`
		Count []struct {
			Total int64 `bson:"total"`
		} `bson:"count"`
	}

	if err := cursor.All(ctx, &rawResult); err != nil {
		return nil, Metadata{}, err
	}

	if len(rawResult) == 0 {
		return nil, Metadata{}, nil
	}
	result := rawResult[0]

	var totalCount int64
	if len(result.Count) > 0 {
		totalCount = result.Count[0].Total
	}

	metadata := calculateMetadata(totalCount, param.Page, param.PageSize)

	return result.Data, metadata, nil
}

type MyFriendsParam struct {
	CurrentUser *User
	Search      string
	Page        int64
	PageSize    int64
}

func (m *UserModel) MyFriends(param MyFriendsParam) ([]*User, Metadata, error) {
	// Step 1: Match hanya user yang merupakan teman dan sudah onboarded
	matchStage := bson.D{{Key: "$match", Value: bson.M{
		"_id":          bson.M{"$in": param.CurrentUser.FriendIDs},
		"is_onboarded": true,
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
