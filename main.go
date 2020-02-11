package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/tv2169145/golang-grpc/api"
)

func main() {
	db, err := xorm.NewEngine("mysql", "root:12345678@tcp(localhost:3306)/grpc")
	if err != nil {
		panic(err)
	}
	_, err = db.DBMetas()
	if err != nil {
		panic(err)
	}
	api.Run(8080, db)
}
