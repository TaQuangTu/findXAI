package lockdb

import "errors"

var (
	OurLockErrorInvalidRedisDb = errors.New("invalid redis database")
)
