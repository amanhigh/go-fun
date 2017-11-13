package tools

import "fmt"

func RunQuery(host string, database string, cmd string) string {
	dbCmd := fmt.Sprintf(`ssh %v "echo 'use %v; %v' | sudo mysql -uroot"`, host, database, cmd)
	return RunCommandPrintError(dbCmd)
}
