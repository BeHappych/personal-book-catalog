package main

import (
	"strings"

	"github.com/spf13/viper"
)

func main() {
	initConfigs()

	print(viper.GetString("db.address"))
}

func initConfigs() {
	println("init configs")

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}

		//добавляем функционал чтения конфигурации из переменных окружения
		//будет использоваться если удалить/переименовать config.yaml
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
	}
}
