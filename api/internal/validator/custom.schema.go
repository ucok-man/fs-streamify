package validator

import (
	"time"

	z "github.com/Oudwins/zog"
)

func Uint64() *z.NumberSchema[uint64] {
	num := &z.NumberSchema[uint64]{}
	return num
}

func Duration() *z.Custom[time.Duration] {
	return z.CustomFunc(func(ptr *time.Duration, ctx z.Ctx) bool {
		_, err := time.ParseDuration((*ptr).String())
		return err == nil
	}, z.Message("Invalid time duration format"))
}
