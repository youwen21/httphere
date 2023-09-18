package conf

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var (
	host = flag.String("host", "", "host the server listen.")
	port = flag.String("port", "", "port the server listen.")

	root = flag.String("root", "", "specify the root.")

	Here HereConf
)

func init() {
	initConf()
}

func initConf() {
	viper.SetConfigName("httphere") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		fmt.Println("ReadInConfig err: ", err)
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

func GetHost() string {
	if *host != "" {
		return *host
	}

	return Here.Base.ListenHost
}

func GetPort() string {
	if *port != "" {
		return *port
	}

	return Here.Base.ListenPort
}

func GetRoot() string {
	if *root != "" {
		return *root
	}

	return Here.Base.StaticRoot
}
