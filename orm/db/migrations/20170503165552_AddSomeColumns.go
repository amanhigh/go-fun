package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20170503165552(txn *sql.Tx) {
	txn.Exec("ALTER TABLE aman.verticals ADD my_column VARCHAR(255);")
}

// Down is executed when this migration is rolled back
func Down_20170503165552(txn *sql.Tx) {
	txn.Exec("ALTER TABLE aman.verticals DROP COLUMN my_column;")

}
