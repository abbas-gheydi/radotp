package confs

import (
	"log"

	"github.com/Abbas-gheydi/radotp/pkgs/rad"

	"github.com/spf13/viper"
)

type Configurations struct {
	Web      webconf
	Database databaseconf
	Ldap     LdapProvider
	Radius   rad.RadConfs
	Metrics  metrics
}

func (c *Configurations) Load() {

	viper.SetConfigType("toml")

	viper.SetConfigName("radiusd.conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/radotp/")
	viper.AddConfigPath("/etc/motp/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// panic(fmt.Errorf("fatal error config file: %w ", err))
		// Detailed error message
		log.Fatalf("fatal error config file: %s ", err.Error())
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Println(err)
	}
}

type webconf struct {
	ListenHTTP                string
	ListenHTTPS               string
	RedirectToHTTPS           bool
	RedirectToHTTPSPortNumber string
	Isuuer                    string
	Apikey                    string
	EnableRestApi             bool
}
type metrics struct {
	EnablePrometheusExporter bool
	Listen                   string
	PromethuesAddress        string
}
type databaseconf struct {
	Server,
	Port,
	Username,
	Password,
	Dbname,
	Sslmode,
	Timezone string
	MaxOpenConns,
	MaxIdleConns,
	ConnMaxLifetimeInMiuntes int
}

type LdapProvider struct {
	FortiGroups                  []string
	LdapGroupsFilter             string
	LdapServers                  []string
	Basedn                       string
	Port                         int
	Security                     int
	ForceSearchForSamAccountName bool
	PreWin2kLogonNameDomain      string
}
