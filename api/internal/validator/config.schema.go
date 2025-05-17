package validator

import (
	z "github.com/Oudwins/zog"
)

var configSchema = z.Struct(z.Schema{
	"port": z.Int().Required().LT(65535, z.Message("Port must be at most 65535")),
	"env":  z.String().Required().OneOf([]string{"development", "staging", "production"}),
	"log": z.Struct(z.Schema{
		"level": z.String().Required().OneOf([]string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}),
	}),
	"cors": z.Struct(z.Schema{
		"origins": z.Slice(z.String().URL(z.Message("Each origins item must be valid url"))).Optional(),
	}),
})
