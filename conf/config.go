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

type FastCGIConf struct {
	FastCGI        string
	FastCGIProto   string
	FastCGIAddress string
	FastCGIRoot    string
}

var (
	defaultPort = "80"
	port        = flag.String("port", "", "Port to run the server.")

	root    = flag.String("root", "", "specify the root.")
	backend = flag.String("backend", "", "default backend server URL.")

	config map[string]interface{}
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
	config = make(map[string]interface{})

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		fmt.Println("read config from embed data")
		viper.SetConfigType("toml")
		viper.ReadConfig(bytes.NewReader(embedConfig))
	} else {
		fmt.Println("viper get config file:", viper.ConfigFileUsed())
	}
	config = viper.AllSettings()
	b, _ := json.MarshalIndent(config, "", "  ")
	fmt.Print(string(b))
}

func GetRoot() string {
	if *root != "" {
		return *root
	}

	if rootHere := os.Getenv("ROOT"); rootHere != "" {
		return rootHere
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Getwd err : %v", err)
		os.Exit(1)
	}
	return pwd
}

func GetPort() string {
	if *port != "" {
		return *port
	}

	if ePort := os.Getenv("HERE_PORT"); ePort != "" {
		return ePort
	}

	return defaultPort
}

func GetBackend() string {
	if *backend != "" {
		return *backend
	}

	if eBackend := os.Getenv("HERE_BACKEND"); eBackend != "" {
		return eBackend
	}

	return ""
}

func GetConfig() map[string]interface{} {
	return config
}

//func GetRewrite() map[string]string {
//	return rewrite
//}
