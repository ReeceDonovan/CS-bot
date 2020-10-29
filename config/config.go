package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() error {
	// Viper
	initDefaults()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	printAll()
	return nil
}

func printAll() {
	fmt.Println("Startup variables:")
	for k, v := range viper.AllSettings() {
		fmt.Println(k + ":")
		for sk, sv := range v.(map[string]interface{}) {
			if strval, ok := sv.(string); ok {
				if len(strval) > 5 {
					fmt.Printf("%s: %s...\n", sk, strval[:5])
				} else {
					fmt.Printf("%s: %s\n", sk, strval)
				}
			} else {
				fmt.Printf("%s: %v\n", sk, sv)
			}
		}
	}
}
