package validator

import z "github.com/Oudwins/zog"

var getAllFromFriendRequestSchema = z.Struct(z.Schema{
	"Page":         z.Int().Required().GTE(1).LTE(100),
	"PageSize":     z.Int().Required().GTE(1).LTE(1000),
	"SearchSender": z.String().Trim(),
	"Status":       z.String().OneOf([]string{"All", "Pending", "Accepted"}),
})
