package auth

type User struct {
	Username string
	Password string
	Role     string
}

var users = []User{
	{Username: "user1", Password: "passwd", Role: "user"},
	{Username: "user2", Password: "passwd", Role: "user"},
	{Username: "admin", Password: "passwd", Role: "admin"},
}
