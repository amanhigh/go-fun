package tools

import "fmt"

func RunQuery(host string, database string, query string) string {
	dbCmd := fmt.Sprintf(`ssh %v "echo 'use %v; %v' | sudo mysql -uroot"`, host, database, query)
	return RunCommandPrintError(dbCmd)
}
