package config

type FunAppConfig struct {
	Server Server
	Db     Db
	Http   HttpClientConfig
}
