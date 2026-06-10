package auth

import (
	"testing"
)

// TestGeneratePassword 测试密码加密功能
func TestGeneratePassword(t *testing.T) {
	// 测试用例1：正常密码
	password1 := "123456"
	hash1, err := GeneratePassword(password1)
	if err != nil {
		t.Fatalf("加密正常密码失败: %v", err)
	}
	if len(hash1) == 0 {
		t.Error("加密结果为空")
	}

	// 测试用例2：空密码（边界测试）
	password2 := ""
	hash2, err := GeneratePassword(password2)
	if err != nil {
		t.Fatalf("加密空密码失败: %v", err)
	}
	if len(hash2) == 0 {
		t.Error("空密码加密结果为空")
	}

	// 重点：相同密码两次加密结果必须不同（bcrypt 自动加盐）
	hash1Again, _ := GeneratePassword(password1)
	if hash1 == hash1Again {
		t.Error("相同密码两次加密结果相同，不符合 bcrypt 特性")
	}
}

// TestComparePassword 测试密码对比功能
func TestComparePassword(t *testing.T) {
	// 1. 先生成哈希
	rawPwd := "my-password-123"
	hashPwd, err := GeneratePassword(rawPwd)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	// 测试用例1：正确密码 → 必须返回 true
	err = ComparePassword(hashPwd, rawPwd)
	if err != nil {
		t.Error("正确密码对比失败，应该返回 true")
	}

	// 测试用例2：错误密码 → 必须返回 false
	err = ComparePassword(hashPwd, "wrong-password")
	if err != nil {
		t.Error("错误密码对比成功，应该返回 false")
	}

	// 测试用例3：空哈希 + 任意密码 → false
	err = ComparePassword("", rawPwd)
	if err != nil {
		t.Error("空哈希应该对比失败")
	}

	// 测试用例4：正确空密码对比
	emptyPwd := ""
	emptyHash, _ := GeneratePassword(emptyPwd)
	err = ComparePassword(emptyHash, emptyPwd)
	if err != nil {
		t.Error("空密码正确对比应该返回 true")
	}
}
