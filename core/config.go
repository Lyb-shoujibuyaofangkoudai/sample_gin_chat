package core

import (
	"fmt"
	"github.com/spf13/viper"
)

const ConfigFilePath = "config/settings.yaml"

func InitConfig() {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config") // 在config目录下查找配置文件
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print("配置文件读取失败：", err)
	}
	//dataBytes, err := os.ReadFile(ConfigFilePath)
	//if err != nil {
	//	fmt.Println("读取文件失败：", err)
	//	return
	//}
	////fmt.Println("yaml 文件的内容: \n", string(dataBytes))
	//c := config.Config{}
	//err = yaml.Unmarshal(dataBytes, &c)
	//if err != nil {
	//	fmt.Println("解析 yaml 文件失败：", err)
	//	return
	//}
	//global.Config = &c // 存入全局变量
}
