package validator

import z "github.com/Oudwins/zog"

var myFriendsSchema = z.Struct(z.Schema{
	"Page":     z.Int().Required().GTE(1).LTE(100),
	"PageSize": z.Int().Required().GTE(1).LTE(1000),
	"Search":   z.String().Trim(),
})
