package config

type FunAppConfig struct {
	Server    Server
	RateLimit RateLimit
	Db        Db
	Log       Log
	Http      HttpClientConfig
	Tracing   Tracing
}
