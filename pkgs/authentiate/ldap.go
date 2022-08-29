package authentiate

import (
	"log"
	"strings"
	"sync"
	"time"

	ldapAuth "github.com/korylprince/go-ad-auth/v3"
)

var (
	ldapMutex                         sync.RWMutex
	lastTimeLdapProviderChanged       time.Time
	switchToAnotherLdapServerInterval = time.Second * 30
)

type LdapProvider struct {
	LdapConfig  *ldapAuth.Config
	Groups      []string
	LdapServers []string
}

func (l LdapProvider) changeLdapSegver() {
	if len(l.LdapServers) > 1 && time.Now().After(lastTimeLdapProviderChanged.Add(switchToAnotherLdapServerInterval)) {
		ldapMutex.Lock()

		for _, srv := range l.LdapServers {
			if srv == l.LdapConfig.Server {
				continue
			}

			l.LdapConfig.Server = srv
			log.Println("change ldap server to ", srv)
			lastTimeLdapProviderChanged = time.Now()

			break

		}
		ldapMutex.Unlock()

	}

}

func (l LdapProvider) IsUserAuthenticated(username string, password string) (authStat bool, groups []string) {
	ldapMutex.RLock()
	//log.Println("ldap server address", l.LdapConfig.Server)
	authStat, _, groups, err := ldapAuth.AuthenticateExtended(l.LdapConfig, username, password, []string{"cn"}, l.Groups)
	//log.Printf("status %v entry %v groups %v", authStat, entry, groups)
	defer ldapMutex.RLocker().Unlock()

	if err != nil {

		log.Println(err)
		//if group name in settings is not true
		if strings.Contains(err.Error(), "Search error") {

			l.Groups = nil
			authStat, _, groups, _ = ldapAuth.AuthenticateExtended(l.LdapConfig, username, password, []string{"cn"}, l.Groups)
		}
		//try another ldap server
		if strings.Contains(err.Error(), "Connection error") {
			go l.changeLdapSegver()
		}

		return

	}

	return

}
