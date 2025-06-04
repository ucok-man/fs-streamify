package validator

import z "github.com/Oudwins/zog"

var signupDTOSchema = z.Struct(z.Schema{
	"Fullname": z.String().Required().Min(3).Max(255),
	"Email":    z.String().Email(),
	"Password": z.String().Min(8).Max(32).ContainsUpper().ContainsDigit().ContainsSpecial(),
})
