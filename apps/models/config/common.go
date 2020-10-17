package config

import (
	"errors"
	"fmt"
	"time"
)

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (self *Server) GetUrl(uri string) string {
	return fmt.Sprintf("http://%v:%v%v", self.Host, self.Port, uri)
}

type Vault struct {
	Server `yaml:",inline"`
	Token  string `yaml:"token"`
}

type Db struct {
	Env         string `yaml:"env"`
	Url         string `yaml:"url"`
	AutoMigrate bool   `yaml:"auto_migrate"`
}

type HttpClientConfig struct {
	/* Timeouts */
	DialTimeout           time.Duration `yaml:"dial_timeout"`
	RequestTimeout        time.Duration `yaml:"request_timeout"`
	IdleConnectionTimeout time.Duration `yaml:"idle_connection_timeout"`

	/* Flags */
	KeepAlive   bool `yaml:"keep_alive"`
	Compression bool `yaml:"compression"`

	IdleConnectionsPerHost int `yaml:"idle_connections_per_host"`
}

type ZoneMap map[string]Server

func (self ZoneMap) GetUrl(zone, uri string) (url string, err error) {
	if server, ok := self[zone]; ok {
		return server.GetUrl(uri), nil
	} else {
		return "", errors.New("INVALID_ZONE: " + zone)
	}
}
