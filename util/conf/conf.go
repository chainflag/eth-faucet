package conf

import (
	"github.com/spf13/viper"
)

var c *conf

func init() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	c = &conf{viper: v}
}

type conf struct {
	viper *viper.Viper
}

func GetInt(key string) int       { return c.viper.GetInt(key) }
func GetString(key string) string { return c.viper.GetString(key) }
