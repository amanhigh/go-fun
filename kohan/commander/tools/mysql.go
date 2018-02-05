package tools

import "fmt"

func RunRootQuery(host string, database string, query string) string {
	dbCmd := fmt.Sprintf(`ssh %v "%v | sudo mysql -uroot"`, host, getDatabaseQuery(database, query))
	return RunCommandPrintError(dbCmd)
}

func RunQuery(host string, user string, pass string, database string, query string) string {
	cmd := fmt.Sprintf("%v | mysql -u %v -p%v -h %v", getDatabaseQuery(database, query), user, pass, host)
	return RunCommandPrintError(cmd)
}

func getDatabaseQuery(database string, query string) string {
	return fmt.Sprintf("echo 'use %v; %v'", database, query)
}
