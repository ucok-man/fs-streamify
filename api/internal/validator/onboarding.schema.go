package validator

import z "github.com/Oudwins/zog"

var onboardingDTOSchema = z.Struct(z.Schema{
	"Fullname":    z.String().Trim().Required().Min(3).Max(255),
	"Bio":         z.String().Trim().Required().Min(10).Max(255),
	"NativeLng":   z.String().Trim().Required(),
	"LearningLng": z.String().Trim().Required(),
	"Location":    z.String().Trim().Required(),
})
