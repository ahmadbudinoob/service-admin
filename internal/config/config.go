package config

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/spf13/viper"
)

var OracleInfo struct {
	Host        string
	Port        int
	ServiceName string
	User        string
	Pass        string
}

var MainApp struct {
	Host string
	Port string
}

var (
	CryptoKey []byte
)

func Init() {
	var err error

	envConfig := viper.New()
	envConfig.SetEnvPrefix("idx_dataref_app")
	envConfig.AutomaticEnv()

	envConfig.SetConfigName("admin")
	envConfig.SetConfigType("yaml")
	envConfig.AddConfigPath("$HOME/config")
	envConfig.AddConfigPath(".")

	if err = envConfig.ReadInConfig(); err != nil {
		panic(err)
	}

	hasher := md5.New()
	hasher.Write([]byte(envConfig.GetString("crypto_key")))
	CryptoKey = []byte(hex.EncodeToString(hasher.Sum(nil)))

	MainApp.Host = envConfig.GetString("listen_on_host")
	MainApp.Port = envConfig.GetString("listen_on_port")

	OracleInfo.Host = envConfig.GetString("oracle_host")
	OracleInfo.Port = envConfig.GetInt("oracle_port")
	OracleInfo.User = envConfig.GetString("oracle_user")
	OracleInfo.Pass = envConfig.GetString("oracle_pass")
	OracleInfo.ServiceName = envConfig.GetString("oracle_service_name")

}
