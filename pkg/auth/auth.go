package auth

import "golang.org/x/crypto/bcrypt"

// GeneratePassword 加密密码（注册时使用）
// 参数：明文密码
// 返回：加密后的哈希字符串、错误
func GeneratePassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes), err
}

// ComparePassword 对比密码（登录时使用）
// 参数：db中存储的哈希密码、用户输入的明文密码
// 返回： error
func ComparePassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
