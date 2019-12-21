package redis

func GenAuthorizationKey(token string) string {
	return "token:" + token
}

func GenUserKey(username string) string {
	return "user:" + username
}

func GenCouponKey(username, couponName string) string {
	return "user:" + username + ":coupon:" + couponName
}

func GenCouponLeftKey(username, couponName string) string {
	return "user:" + username + ":coupon:" + couponName + ":left"
}

func GenCouponOwnersKey(username, couponName string) string {
	return "user:" + username + ":coupon:" + couponName + ":owners"
}
