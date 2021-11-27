package migrations

import (
	"database/sql"
	"fmt"
	model2 "github.com/amanhigh/go-fun/components/learn/frameworks/orm/model"
)

// Up is executed when this migration is applied
func Up_20170503165552(txn *sql.Tx) {
	fmt.Printf("%+v\n", txn)
	txn.Exec("ALTER TABLE aman.verticals ADD my_column VARCHAR(255);")
	migrateData()
}
func migrateData() {
	db := orm.DB
	defer db.Close()
	verticals := new([]model2.Vertical)
	db.Find(verticals)
	for _, vertical := range *verticals {
		vertical.MyColumn = "My" + vertical.Name
		db.Save(vertical)
	}
}

// Down is executed when this migration is rolled back
func Down_20170503165552(txn *sql.Tx) {
	txn.Exec("ALTER TABLE aman.verticals DROP COLUMN my_column;")

}
