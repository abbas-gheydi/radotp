package auth

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// SecurityType specifies the type of security to use when connecting to an Active Directory Server.
type SecurityType int

// Security will default to SecurityNone if not given.
const (
	SecurityNone SecurityType = iota
	SecurityTLS
	SecurityStartTLS
	SecurityInsecureTLS
	SecurityInsecureStartTLS
)

// Config contains settings for connecting to an Active Directory server.
type Config struct {
	Server                       string
	Port                         int
	BaseDN                       string
	Security                     SecurityType
	RootCAs                      *x509.CertPool
	ForceSearchForSamAccountName bool
	PreWin2kLogonNameDomain      string
}

// Domain returns the domain derived from BaseDN or an error if misconfigured.
func (c *Config) Domain() (string, error) {
	domain := ""
	for _, v := range strings.Split(strings.ToLower(c.BaseDN), ",") {
		if trimmed := strings.TrimSpace(v); strings.HasPrefix(trimmed, "dc=") {
			domain = domain + "." + trimmed[3:]
		}
	}
	if len(domain) <= 1 {
		return "", errors.New("Configuration error: invalid BaseDN")
	}
	return domain[1:], nil
}

// UPN returns the userPrincipalName for the given username or an error if misconfigured.
func (c *Config) UPN(username string) (string, error) {
	if _, err := mail.ParseAddress(username); err == nil {
		return username, nil
	}

	domain, err := c.Domain()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s@%s", username, domain), nil
}

func (c *Config) SamAccountName(username string) (fullUserName, user string, err error) {
	// Split the username into user and domain parts
	tmpList := strings.SplitN(username, "@", 2)
	if len(tmpList) != 2 {
		return "", "", errors.New("invalid username format")
	}

	// Extract user and domain
	user = tmpList[0]
	domainParts := strings.Split(tmpList[1], ".")
	if len(domainParts) < 2 {
		return "", "", errors.New("invalid domain format")
	}
	domain := domainParts[0]
	if c.PreWin2kLogonNameDomain != "" {
		domain = c.PreWin2kLogonNameDomain
	}

	// Construct the full user name using domain and username
	fullUserName = domain + `\` + user
	return fullUserName, user, nil
}

func (c *Config) ExtractUserName(username string) (fullUserName, user string, err error) {
	// Extract full user name using User Principal Name (UPN)
	fullUserName, err = c.UPN(username)
	if err != nil {
		return "", "", err
	}

	user = fullUserName

	// If forced to search for SamAccountName, extract it
	if c.ForceSearchForSamAccountName {
		var samFullUserName, samUser string
		samFullUserName, samUser, err = c.SamAccountName(user)
		if err != nil {
			return "", "", err
		}

		// Update values only if SamAccountName was successfully extracted
		if samUser != "" && samFullUserName != "" {
			fullUserName = samFullUserName
			user = samUser
		}
	}

	if user == "" || fullUserName == "" {
		err = errors.New("ldap error parsing username")
	}

	return fullUserName, user, nil
}
