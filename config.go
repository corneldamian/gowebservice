package gowebservice

import (
	"os"
	"path/filepath"

	"gopkg.in/gcfg.v1"
)

type WebServer struct {
	Address           string
	Log_Dir           string
	Enable_Access_Log bool
	Access_Log        string
	System_Log        string
	Log_Level         string
}

var cfg WebServerConfiger

func initConfig(ocfg WebServerConfiger) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg = ocfg

	return gcfg.ReadFileInto(cfg, filepath.Join(dir, "config.ini"))
}

func GetConfig() WebServerConfiger {
	return cfg
}

type WebServerConfiger interface {
	WebServerConfig() *WebServer
}
