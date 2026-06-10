package regex

import "regexp"

var (
	// Email 匹配常见邮箱格式，如 user@example.com
	Email = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// Phone 匹配国际手机号（E.164），可选 + 前缀，总长 7–15 位
	Phone = regexp.MustCompile(`^\+?[1-9]\d{6,14}$`)
)
