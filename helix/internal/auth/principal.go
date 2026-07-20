package auth

type Principal struct {
	UserID string
	Email  string
	Roles  []string
}
