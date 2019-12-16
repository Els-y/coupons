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

func GenCouponLeftKey(salerName, couponName string) string {
	return "coupon:left:" + salerName + ":" + couponName
}

func GenCouponOwnersKey(couponName string) string {
	return "coupon:owners:" + couponName
}
