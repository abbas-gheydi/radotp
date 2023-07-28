package authentiate

type Auth_Provider interface {
	IsUserAuthenticated(username string, password string, checkForvendorFortinetGroup bool) (authstat bool, groups []string)
}
