package confs

import (
	"fmt"

	"github.com/Abbas-gheydi/radotp/pkgs/monitoring"
	"github.com/Abbas-gheydi/radotp/pkgs/rad"
	"github.com/Abbas-gheydi/radotp/pkgs/storage"
	"github.com/Abbas-gheydi/radotp/pkgs/web"
	ldapAuth "github.com/korylprince/go-ad-auth/v3"
)

var Cfg Configurations

func LoadConfigs() {
	//load configs from file
	Cfg.Load()

	//radius configs
	rad.RadiusConfigs = Cfg.Radius

	//web configs
	web.ListenAddr = Cfg.Web.Listen
	web.QrIssuer = Cfg.Web.Isuuer
	web.ApiKey = Cfg.Web.Apikey
	web.PromethuesServerAddress = Cfg.Metrics.PromethuesAddress
	web.EnableRestApi = Cfg.Web.EnableRestApi

	//ldap configs
	//rad.Auth_Provider = cfg.Ldap
	rad.Auth_Provider.Groups = Cfg.Ldap.Groups
	rad.Auth_Provider.LdapConfig = &ldapAuth.Config{}
	rad.Auth_Provider.LdapConfig.BaseDN = Cfg.Ldap.Basedn
	rad.Auth_Provider.LdapConfig.Port = Cfg.Ldap.Port
	rad.Auth_Provider.LdapConfig.Security = ldapAuth.SecurityType(Cfg.Ldap.Security)
	rad.Auth_Provider.LdapConfig.Server = Cfg.Ldap.LdapServers[0]
	rad.Auth_Provider.LdapServers = Cfg.Ldap.LdapServers

	//database configs
	storage.Dsn = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", Cfg.Database.Server, Cfg.Database.Username, Cfg.Database.Password, Cfg.Database.Dbname, Cfg.Database.Port, Cfg.Database.Sslmode, Cfg.Database.Timezone)
	//monitoring
	monitoring.Listen = Cfg.Metrics.Listen
	rad.EnableMetrics = Cfg.Metrics.EnablePrometheusExporter

}
