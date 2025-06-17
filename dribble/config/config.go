package config

import (
	"github.com/ctrl-alt-boop/gooldb/pkg/logging"
	"github.com/spf13/viper"
)

var logger = logging.NewLogger("config.log")

var Cfg *Config
var viperInstance *viper.Viper

// var defaults map[string]any = map[string]any{
// 	"showDetails": false,
// }

// var Defined map[string]any = map[string]any{
// 	"showDetails": true,
// }

// var mergedConfigs map[string]any = mergeConfigs()

// func mergeConfigs() map[string]any {
// 	merged := make(map[string]any)
// 	maps.Copy(merged, defaults)
// 	maps.Copy(merged, Defined)
// 	return merged
// }

func GetConfigMap() map[string]any {
	return viperInstance.AllSettings()
}

func Get(key string) any {
	return viperInstance.AllSettings()[key]
}

func Set(key string, value any) {
	viperInstance.Set(key, value)
}

func Save() error {
	return viperInstance.WriteConfig()
}

func LoadConfig() {
	var err error
	Cfg, viperInstance, err = loadConfig()
	if err != nil {
		panic(err)
	}
}
