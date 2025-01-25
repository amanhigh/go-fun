package models

type contextKey string

const (
	XRequestID contextKey = "X-Request-ID"
	ContextTx  contextKey = "contextTx"
)

// Ginkgo
const GINKGO_INEGRATION = "it"
const GINKGO_SLOW = "slow"   // Slower Tests
const GINKGO_SETUP = "setup" // Requires external Setup
const VAULT_ROOT_TOKEN = "root-token"

// DB
const MYSQL = "mysql"
const POSTGRES = "postgres"
const SQLITE = "sqlite"
