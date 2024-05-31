package confs

import (
	"fmt"

	"github.com/Abbas-gheydi/radotp/pkgs/monitoring"
	"github.com/Abbas-gheydi/radotp/pkgs/rad"
	"github.com/Abbas-gheydi/radotp/pkgs/storage"
	"github.com/Abbas-gheydi/radotp/pkgs/web"
	ldapAuth "github.com/abbas-gheydi/go-ad-auth/v3"
)

var Cfg Configurations

func LoadConfigs() {
	//load configs from file
	Cfg.Load()

	//radius configs
	rad.RadiusConfigs = Cfg.Radius

	//web configs
	web.HTTPListenAddr = Cfg.Web.ListenHTTP
	web.HTTPSListenAddr = Cfg.Web.ListenHTTPS
	web.RedirectToHTTPS = Cfg.Web.RedirectToHTTPS
	web.RedirectToHTTPSPortNumber = Cfg.Web.RedirectToHTTPSPortNumber
	web.QrIssuer = Cfg.Web.Isuuer
	web.ApiKey = Cfg.Web.Apikey
	web.PromethuesServerAddress = Cfg.Metrics.PromethuesAddress
	web.EnableRestApi = Cfg.Web.EnableRestApi

	//ldap configs
	//rad.Auth_Provider = cfg.Ldap
	rad.Auth_Provider.FortiGroups = Cfg.Ldap.FortiGroups
	if Cfg.Ldap.LdapGroupsFilter != "" {
		rad.Auth_Provider.LdapGroupsFilter = []string{Cfg.Ldap.LdapGroupsFilter}
	}
	rad.Auth_Provider.LdapConfig = &ldapAuth.Config{}
	rad.Auth_Provider.LdapConfig.BaseDN = Cfg.Ldap.Basedn
	rad.Auth_Provider.LdapConfig.Port = Cfg.Ldap.Port
	rad.Auth_Provider.LdapConfig.Security = ldapAuth.SecurityType(Cfg.Ldap.Security)
	rad.Auth_Provider.LdapConfig.Server = Cfg.Ldap.LdapServers[0]
	rad.Auth_Provider.LdapServers = Cfg.Ldap.LdapServers
	rad.Auth_Provider.LdapConfig.ForceSearchForSamAccountName = Cfg.Ldap.ForceSearchForSamAccountName

	//database configs
	storage.Dsn = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", Cfg.Database.Server, Cfg.Database.Username, Cfg.Database.Password, Cfg.Database.Dbname, Cfg.Database.Port, Cfg.Database.Sslmode, Cfg.Database.Timezone)
	storage.MaxOpenConns = Cfg.Database.MaxOpenConns
	storage.MaxIdleConns = Cfg.Database.MaxIdleConns
	storage.ConnMaxLifetimeInMiuntes = Cfg.Database.ConnMaxLifetimeInMiuntes
	//monitoring
	monitoring.Listen = Cfg.Metrics.Listen
	rad.EnableMetrics = Cfg.Metrics.EnablePrometheusExporter

}
