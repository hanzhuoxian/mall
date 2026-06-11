package admin

const AdminUserName = "admin"

func IsAdmin(username string) bool {
	return username == AdminUserName
}
