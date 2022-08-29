package main

import (
	"fmt"

	"github.com/Abbas-gheydi/radotp/pkgs/confs"
	"github.com/Abbas-gheydi/radotp/pkgs/monitoring"
	"github.com/Abbas-gheydi/radotp/pkgs/rad"
	"github.com/Abbas-gheydi/radotp/pkgs/storage"
	"github.com/Abbas-gheydi/radotp/pkgs/web"
	ldapAuth "github.com/korylprince/go-ad-auth/v3"
)

var cfg confs.Configurations

func art() {

	fmt.Print(`	 
                                                                  
    _/_/_/      _/_/    _/_/_/      _/_/    _/_/_/_/_/  _/_/_/    
   _/    _/  _/    _/  _/    _/  _/    _/      _/      _/    _/   
  _/_/_/    _/_/_/_/  _/    _/  _/    _/      _/      _/_/_/      
 _/    _/  _/    _/  _/    _/  _/    _/      _/      _/           
_/    _/  _/    _/  _/_/_/      _/_/        _/      _/            
                                                                  
                                                                  
                                                              
`)
}
func loadConfigs() {
	//load configs from file
	cfg.Load()

	//radius configs
	rad.RadiusConfigs = cfg.Radius

	//web configs
	web.ListenAddr = cfg.Web.Listen
	web.QrIssuer = cfg.Web.Isuuer
	web.ApiKey = cfg.Web.Apikey
	web.PromethuesServerAddress = cfg.Metrics.PromethuesAddress
	web.EnableRestApi = cfg.Web.EnableRestApi

	//ldap configs
	//rad.Auth_Provider = cfg.Ldap
	rad.Auth_Provider.Groups = cfg.Ldap.Groups
	rad.Auth_Provider.LdapConfig = &ldapAuth.Config{}
	rad.Auth_Provider.LdapConfig.BaseDN = cfg.Ldap.Basedn
	rad.Auth_Provider.LdapConfig.Port = cfg.Ldap.Port
	rad.Auth_Provider.LdapConfig.Security = ldapAuth.SecurityType(cfg.Ldap.Security)
	rad.Auth_Provider.LdapConfig.Server = cfg.Ldap.LdapServer[0]
	rad.Auth_Provider.LdapServers = cfg.Ldap.LdapServer

	//database configs
	storage.Dsn = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", cfg.Database.Server, cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Port, cfg.Database.Sslmode, cfg.Database.Timezone)
	//monitoring
	monitoring.Listen = cfg.Metrics.Listen
	rad.EnableMetrics = cfg.Metrics.EnablePrometheusExporter
}

func main() {
	loadConfigs()

	art()

	//database configuraion
	storage.Initialize()

	if cfg.Metrics.EnablePrometheusExporter {
		go monitoring.Start()
	}

	go web.StartRouter()

	rad.StartRadius()

}
