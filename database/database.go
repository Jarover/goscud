//
// Package database - инициализация баз данных
//
package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *sql.DB   //база данных FB
var dbp *gorm.DB //база данных PostgreSQL

// InitDB Инициализация базы данных FB
func InitDB(dbHost, dbName, username, password string) {

	DbaUrl := username + ":" + password + "@" + dbHost + "/" + dbName

	conn, err := sql.Open("firebirdsql", DbaUrl)

	if err != nil {
		log.Fatal(err)
	}

	db = conn

	//db.Debug().AutoMigrate(&Account{}, &Contact{}) //Миграция базы данных
}

// InitDBP Инициализация базы данных Postgre
func InitDBP(dbHost, dbName, username, password string) {

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Создать строку подключения
	//fmt.Println(dbUri)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	}

	dbp = conn

	//db.Debug().AutoMigrate(&Account{}, &Contact{}) //Миграция базы данных
}

// GetDB возвращает дескриптор объекта DB FB
func GetDB() *sql.DB {
	return db
}

// GetDBP возвращает дескриптор объекта DB PostgreSQL
func GetDBP() *gorm.DB {
	return dbp
}
