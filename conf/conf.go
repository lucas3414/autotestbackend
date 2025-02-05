package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("setting")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./conf/")

	err := viper.ReadInConfig()

	if err != nil {
		panic("读取配置文件失败:" + err.Error())
	}

	fmt.Println(viper.GetString("server.port"))

}
