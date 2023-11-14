package config

type FunAppConfig struct {
	Server    Server
	RateLimit RateLimit
	Db        Db
	Http      HttpClientConfig
	Tracing   Tracing
}
