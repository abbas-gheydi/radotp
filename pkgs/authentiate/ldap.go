package authentiate

import (
	"log"
	"strings"

	ldapAuth "github.com/korylprince/go-ad-auth/v3"
)

type LdapProvider struct {
	LdapConfig ldapAuth.Config
	Groups     []string
}

func (l LdapProvider) IsUserAuthenticated(username string, password string) (authStat bool, groups []string) {

	authStat, _, groups, err := ldapAuth.AuthenticateExtended(&l.LdapConfig, username, password, []string{"cn"}, l.Groups)
	//log.Printf("status %v entry %v groups %v", authStat, entry, groups)

	if err != nil {

		log.Println(err)
		//if group name in settings is not true
		if strings.Contains(err.Error(), "Search error") {

			l.Groups = nil
			authStat, _, groups, _ = ldapAuth.AuthenticateExtended(&l.LdapConfig, username, password, []string{"cn"}, l.Groups)
		}

		return

	}

	return

}
