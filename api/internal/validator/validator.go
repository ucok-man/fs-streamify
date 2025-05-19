package validator

import (
	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/internals"
)

type schema struct {
	Config    *z.StructSchema
	SignupDTO *z.StructSchema
	SigninDTO *z.StructSchema
}

func Schema() schema {
	return schema{
		Config:    configSchema,
		SignupDTO: signupDTOSchema,
		SigninDTO: signinDTOSchema,
	}
}

func Sanitize(m map[string][]*internals.ZogIssue) map[string][]string {
	return z.Issues.SanitizeMap(m)
}
