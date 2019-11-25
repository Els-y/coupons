package redis

func GenAuthorizationKey(token string) string {
	return "token:" + token
}

func GenUserKey(username string) string {
	return "user:" + username
}

func GenCouponKey(username, couponName string) string {
	return "coupon:" + username + ":" + couponName
}
