package rediskeys

import "fmt"

const (
	invalidTokenFormat = "invalid_token:%s"
)

func InvalidTokenKey(token string) string {
	return fmt.Sprintf(invalidTokenFormat, token)
}
