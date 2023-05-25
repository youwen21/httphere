package conf

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

var (
	port = flag.String("port", "", "Port to run the server.")

	root    = flag.String("root", "", "specify the root.")
	backend = flag.String("backend", "", "default backend server URL.")

	Here HereConf
)

func init() {
	initConf()
	//initRewrite()
}

//func initRewrite() {
//	rewriteRules := os.Getenv("REWRITE")
//	if rewriteRules == "" {
//		return
//	}
//	rewrite = make(map[string]string)
//
//	rules := strings.Split(rewriteRules, ",")
//	for _, v := range rules {
//		if !strings.Contains(v, "=") {
//			fmt.Println("rewrite rule invalid: " + v)
//			continue
//		}
//
//		rule := strings.Split(v, "=")
//		rewrite[rule[0]] = rule[1]
//	}
//}

func initConf() {
	viper.SetConfigName("httphere") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println("read config from embed data")
		//viper.SetConfigType("yaml")
		viper.ReadConfig(bytes.NewReader(embedConfig))
	} else {
		fmt.Println("viper get config file:", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&Here)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	settings := viper.AllSettings()
	b, _ := json.MarshalIndent(settings, "", "  ")
	fmt.Println(string(b))
}

func GetRoot() string {
	if *root != "" {
		return *root
	}

	return Here.Base.StaticRoot
}

func GetPort() string {
	if *port != "" {
		return *port
	}

	return Here.Base.ListenPort
}

func GetBackend() string {
	if *backend != "" {
		return *backend
	}

	return ""
}
