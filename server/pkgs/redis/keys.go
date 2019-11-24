package redis

func GenAuthorizationKey(token string) string {
	return "token:" + token
}
