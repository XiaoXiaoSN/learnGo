package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	if err := rwFile(); err != nil {
		fmt.Println(err)
	}

	if err := writeByStruct(); err != nil {
		fmt.Println(err)
	}
}

func rwFile() error {
	defer fmt.Println("== done ==")

	// file in ./config.yml
	viper.AddConfigPath("./configs")  // optionally look for config in the working directory
	viper.SetConfigName("config.yml") // name of config file (without extension)
	// viper.SetConfigType("yaml")   // or viper.SetConfigType("YAML")

	viper.SetDefault("default", true)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Config file was found but another error was produced
	fmt.Println(">>", viper.Get("env"))

	err := viper.WriteConfigAs("./.config.yml")
	if err != nil {
		return err
	}

	return nil
}

type conf struct {
	Env          string `yaml:"env" json:"env"`
	Service      []svc  `yaml:"service" json:"service"`
	DefaultValue string `yaml:"defaultValue"`
}

type svc struct {
	Host   string `yaml:"host" json:"host"`
	Enable bool   `yaml:"enable" json:"enable"`
}

func writeByStruct() error {
	cfg := conf{
		Env: "code",
		Service: []svc{
			svc{Host: "h1"},
			svc{Enable: false},
		},
	}

	jsonBytes, _ := json.Marshal(&cfg)

	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.ReadConfig(bytes.NewBuffer(jsonBytes))

	// try to read
	fmt.Println(viper.Get("env"))

	err := viper.WriteConfigAs("./.config.yml")
	if err != nil {
		return err
	}

	cfg2 := conf{}
	err = viper.Unmarshal(&cfg2)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", cfg2)

	return nil
}
