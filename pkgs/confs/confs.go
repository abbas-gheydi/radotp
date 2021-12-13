package confs

import (
	"fmt"
	"log"

	"github.com/Abbas-gheydi/radotp/pkgs/authentiate"
	"github.com/Abbas-gheydi/radotp/pkgs/rad"

	"github.com/spf13/viper"
)

type Configurations struct {
	Web      webconf
	Database databaseconf
	Ldap     authentiate.LdapProvider
	Radius   rad.RadConfs
	Metrics  metrics
}

func (c *Configurations) Load() {

	viper.SetConfigType("toml")

	viper.SetConfigName("radotp.conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/radotp/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w ", err))
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Println(err)
	}
}

type webconf struct {
	Listen string
	Isuuer string
}
type metrics struct {
	EnablePrometheusExporter bool
	Listen                   string
}
type databaseconf struct {
	Server,
	Port,
	Username,
	Password,
	Dbname,
	Connection_max,
	Sslmode,
	Timezone string
}
