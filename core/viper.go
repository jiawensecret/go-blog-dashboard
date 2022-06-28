package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"platform/global"
)

func Viper() *viper.Viper {
	config := "config.yaml"
	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件出错: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置发生变化:", e.Name)
		if err := v.Unmarshal(&global.Config); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.Config); err != nil {
		fmt.Println(err)
	}

	return v
}
