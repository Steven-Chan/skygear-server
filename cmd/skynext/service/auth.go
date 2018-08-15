package service

type UserAuthenticator interface {
	Authenticate(authData interface{}, password string) bool
}

type MemoryUserAuthenticator struct {
	Store map[string]string
}

func (a MemoryUserAuthenticator) Authenticate(authData interface{}, password string) bool {
	username, ok := authData.(string)
	if !ok {
		return false
	}

	return a.Store[username] == password
}
