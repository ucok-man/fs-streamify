package seed

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/0x6flab/namegenerator"
	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/spf13/cobra"
	"github.com/ucok-man/streamify-api/internal/config"
	"github.com/ucok-man/streamify-api/internal/logger"
	"github.com/ucok-man/streamify-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/exp/rand"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

var SeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run database seed migration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		logger, err := logger.New(cfg.Log.Level, cfg.Env)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed initialize logger")
		}

		conn, err := cfg.OpenDB()
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed initialize db connection")
		}
		defer conn.Disconnect(context.Background())

		db := conn.Database(cfg.DB.DatabaseName)
		ng := namegenerator.NewGenerator()
		userColl := db.Collection("users")
		friendRequestColl := db.Collection("friend_request")

		/* ---------------------------------------------------------------- */
		/*                             SEED USER                            */
		/* ---------------------------------------------------------------- */

		userIds := []bson.ObjectID{}
		logger.Info().Msg("Begin seeding users...")
		mc := &models.User{}

		streamclient, err := stream.NewClient(cfg.GetStreamIO.ApiKey, cfg.GetStreamIO.ApiSecret)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed initialize stream chat client")
		}

		for i := 0; i < 100; i++ {
			rndname := ng.Generate()
			name := strings.Join(strings.Split(rndname, "-"), " ")
			user := &models.User{
				FullName:    name,
				Email:       fmt.Sprintf("%s@dummy.com", strings.ToLower(rndname)),
				Bio:         fmt.Sprintf("Hello, I'am %v", name),
				ProfilePic:  getRandomPicturePlaceholder(),
				NativeLng:   getRandomLng(),
				LearningLng: getRandomLng(),
				Location:    "Some City, Country",
				IsOnboarded: true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				FriendIDs:   []bson.ObjectID{},
			}

			result, err := userColl.InsertOne(context.Background(), user)
			if err != nil {
				logger.Fatal().Err(err).Msg("Error inserting user")
			}

			err = user.Password.Set("@Password123")
			if err != nil {
				logger.Fatal().Err(err).Msg("Error hashing password")
			}
			_, err = userColl.UpdateByID(context.Background(), result.InsertedID, bson.D{{
				Key: "$set",
				Value: bson.D{{
					Key:   "password",
					Value: user.Password.Hash,
				}},
			}})
			if err != nil {
				logger.Fatal().Err(err).Msg("Error updating password")
			}

			userID, ok := result.InsertedID.(bson.ObjectID)
			if !ok {
				logger.Fatal().Err(fmt.Errorf("inserted ID is not ObjectID")).Msg("Error insert user result")
			}

			// Create user in getstream.io
			_, err = streamclient.UpsertUser(context.Background(), &stream.User{
				ID:    userID.Hex(),
				Name:  user.FullName,
				Image: user.ProfilePic,
			})
			if err != nil {
				logger.Fatal().Err(err).Msg("Error creating stream user")
			}

			if i == 0 {
				mc.ID = userID
				mc.Password = user.Password
				mc.Email = user.Email
			}
			userIds = append(userIds, userID)
		}
		logger.Info().Msgf("Success seeding %v users", len(userIds))
		logger.Info().Str("password", *mc.Password.Plaintext).Str("email", mc.Email).Msgf("MC Id: %v", mc.ID)

		/* ---------------------------------------------------------------- */
		/*                        SEED FRIEND REQUEST                       */
		/* ---------------------------------------------------------------- */

		logger.Info().Msg("Begin seeding friend request...")
		for i := 1; i <= 61; i++ {
			fr := &models.FriendRequest{}
			if i <= 21 {
				fr.SenderID = mc.ID
				fr.RecipientID = userIds[i]
				fr.Status = models.FriendRequestStatusPending
			} else if i <= 41 {
				fr.SenderID = mc.ID
				fr.RecipientID = userIds[i]
				fr.Status = models.FriendRequestStatusAccepted

				_, err = userColl.UpdateByID(context.Background(), mc.ID, bson.M{
					"$push": bson.M{
						"friend_ids": userIds[i],
					},
				})
				if err != nil {
					logger.Fatal().Err(err).Msg("Error updating friend ids on MC")
				}

				_, err = userColl.UpdateByID(context.Background(), userIds[i], bson.M{
					"$push": bson.M{
						"friend_ids": mc.ID,
					},
				})
				if err != nil {
					logger.Fatal().Err(err).Msg("Error updating friend ids on recipient")
				}
			} else {
				fr.SenderID = userIds[i]
				fr.RecipientID = mc.ID
				fr.Status = models.FriendRequestStatusPending
			}
			fr.CreatedAt = time.Now()
			fr.UpdatedAt = time.Now()

			_, err := friendRequestColl.InsertOne(context.Background(), fr)
			if err != nil {
				logger.Fatal().Err(err).Msg("Error inserting friend request")
			}
		}
		logger.Info().Msg("Success seeding friend request")
	},
}

func getRandomPicturePlaceholder() string {
	idx := rand.Intn(100) + 1
	url := fmt.Sprintf("https://avatar.iran.liara.run/public/%v.png", idx)
	return url
}

func getRandomLng() string {
	idx := rand.Intn(len(LANGUAGES))
	return LANGUAGES[idx]
}

var LANGUAGES = []string{
	"English",
	"Spanish",
	"French",
	"German",
	"Mandarin",
	"Japanese",
	"Korean",
	"Hindi",
	"Russian",
	"Portuguese",
	"Arabic",
	"Italian",
	"Turkish",
	"Dutch",
}
