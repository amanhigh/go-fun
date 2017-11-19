package models

import "time"

type AuthNConfig struct {
	ClientId       string
	Secret         string
	TokenUrl       string
	RequestTimeout time.Duration
}
