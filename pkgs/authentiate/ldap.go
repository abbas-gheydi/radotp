package authentiate

import (
	"log"
	"strings"

	ldapAuth "github.com/korylprince/go-ad-auth/v3"
)

type LdapProvider struct {
	LdapConfig  *ldapAuth.Config
	Groups      []string
	LdapServers []string
}

func (l LdapProvider) changeLapSegver() {
	if len(l.LdapServers) > 1 {

		for _, srv := range l.LdapServers {
			if srv == l.LdapConfig.Server {
				continue
			}
			l.LdapConfig.Server = srv
			log.Println("change ldap server to ", srv)
			break

		}

	}

}

func (l LdapProvider) IsUserAuthenticated(username string, password string) (authStat bool, groups []string) {
	//log.Println("ldap server address", l.LdapConfig.Server)
	authStat, _, groups, err := ldapAuth.AuthenticateExtended(l.LdapConfig, username, password, []string{"cn"}, l.Groups)
	//log.Printf("status %v entry %v groups %v", authStat, entry, groups)

	if err != nil {

		log.Println(err)
		//if group name in settings is not true
		if strings.Contains(err.Error(), "Search error") {

			l.Groups = nil
			authStat, _, groups, _ = ldapAuth.AuthenticateExtended(l.LdapConfig, username, password, []string{"cn"}, l.Groups)
		}
		//try another ldap server
		if strings.Contains(err.Error(), "Connection error") {
			l.changeLapSegver()
		}

		return

	}

	return

}
