package admin

// AdminUserName is the reserved username for the platform administrator.
const AdminUserName = "admin"

// IsAdmin reports whether username belongs to the platform administrator.
func IsAdmin(username string) bool {
	return username == AdminUserName
}
