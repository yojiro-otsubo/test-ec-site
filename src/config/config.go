package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	Port       int
	DBdriver   string
	DBhost     string
	DBname     string
	DBuser     string
	DBpassword string
	StripeKey  string
	PK         string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		Port:       cfg.Section("web").Key("port").MustInt(),
		DBdriver:   cfg.Section("db").Key("driver").String(),
		DBhost:     cfg.Section("db").Key("db_host").String(),
		DBname:     cfg.Section("db").Key("name").String(),
		DBuser:     cfg.Section("db").Key("user").String(),
		DBpassword: cfg.Section("db").Key("password").String(),
		StripeKey:  cfg.Section("stripe").Key("stripe_key").String(),
		PK:         cfg.Section("stripe").Key("publish_key").String(),
	}
}
