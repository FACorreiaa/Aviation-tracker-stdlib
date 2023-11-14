package structs

type SessionManagerKey struct{}

type UserSession struct {
	Id       int
	Username string
	Email    string
}

type User struct {
	Id       int
	Username string
	Email    string
	Password string
}
