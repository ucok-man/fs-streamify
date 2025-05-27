package validator

import (
	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/internals"
)

type schema struct {
	Config                  *z.StructSchema
	SignupDTO               *z.StructSchema
	SigninDTO               *z.StructSchema
	OnboardingDTO           *z.StructSchema
	RecommendedUser         *z.StructSchema
	MyFriendsSchema         *z.StructSchema
	GetAllFromFriendRequest *z.StructSchema
	GetAllSendFriendRequest *z.StructSchema
}

func Schema() schema {
	return schema{
		Config:                  configSchema,
		SignupDTO:               signupDTOSchema,
		SigninDTO:               signinDTOSchema,
		OnboardingDTO:           onboardingDTOSchema,
		RecommendedUser:         recommendedUserSchema,
		MyFriendsSchema:         myFriendsSchema,
		GetAllFromFriendRequest: getAllFromFriendRequestSchema,
		GetAllSendFriendRequest: getAllSendFriendRequestSchema,
	}
}

func Sanitize(m map[string][]*internals.ZogIssue) map[string][]string {
	return z.Issues.SanitizeMap(m)
}
