package authentiate

type Auth_Provider interface {
	IsUserAuthenticated(username string, password string) (authstat bool, groups []string)
}
