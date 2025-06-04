package validator

import z "github.com/Oudwins/zog"

var signinDTOSchema = z.Struct(z.Schema{
	"Email":    z.String().Email(),
	"Password": z.String().Min(8).Max(32).ContainsUpper().ContainsDigit().ContainsSpecial(),
})
