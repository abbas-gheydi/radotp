package authentiate

import (
	"log"
	"strings"
	"sync"
	"time"

	ldapAuth "github.com/abbas-gheydi/go-ad-auth/v3"
)

var (
	ldapMutex                         sync.RWMutex
	lastTimeLdapProviderChanged       time.Time
	switchToAnotherLdapServerInterval = time.Second * 30
)

type LdapProvider struct {
	LdapConfig       *ldapAuth.Config
	FortiGroups      []string
	LdapServers      []string
	LdapGroupsFilter []string
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

func (l LdapProvider) isUserAuthorized(groups []string) bool {
	if len(l.LdapGroupsFilter) == 0 {
		return true
	}
	for _, g := range groups {
		if g == l.LdapGroupsFilter[0] {
			return true
		}
	}
	return false
}

func (l LdapProvider) IsUserAuthenticated(username string, password string, checkForVendorFortinetGroup bool) (isAuthenticated bool, vendorFortinetGroupName []string) {
	winNTSplitChar := "\\"
	if strings.Contains(username, winNTSplitChar) && strings.Split(username, winNTSplitChar)[1] != "" {
		username = strings.Split(username, winNTSplitChar)[1]
	}

	verifyPasswordAndRetrieveGroupsFromLdap := func(groups []string) (isAuthenticated bool, joinedGroupsName []string, err error) {
		ldapMutex.RLock()
		defer ldapMutex.RLocker().Unlock()
		isAuthenticated, _, joinedGroupsName, err = ldapAuth.AuthenticateExtended(l.LdapConfig, username, password, []string{"cn"}, groups)
		return
	}

	isAuthenticated, joinedGroupsName, err := verifyPasswordAndRetrieveGroupsFromLdap(l.LdapGroupsFilter)

	if isAuthenticated {
		isAuthenticated = l.isUserAuthorized(joinedGroupsName)
	}

	if checkForVendorFortinetGroup {
		if isAuthenticated {
			_, vendorFortinetGroupName, err = verifyPasswordAndRetrieveGroupsFromLdap(l.FortiGroups)
		}
	}

	if err != nil {

		log.Println(err)
		//if group name in settings is not true
		if strings.Contains(err.Error(), "Search error") {
			log.Println(l.FortiGroups, "Group name invalid. Check settings.")
		}
		//try another ldap server
		if strings.Contains(err.Error(), "Connection error") {
			go l.changeLdapSegver()
		}
		return
	}
	return

}
