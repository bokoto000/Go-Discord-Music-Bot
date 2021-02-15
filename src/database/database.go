package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

var err error

/*
func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		"golang", "golang", "discordbot")
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()
	query, err := ioutil.ReadFile("configDatabase.sql")
	queryString := string(query)
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(queryString); err != nil {
		panic(err)
	}

	fmt.Println("# Inserting values")

	var lastInsertId int
	err = db.QueryRow("INSERT INTO keywords(guild_id,keyword_key,keyword_value) VALUES($1,$2,$3) returning id;", "0", "test", "test value").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)

	fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM keywords")
	checkErr(err)

	for rows.Next() {
		var id int
		var guild_id string
		var keyword_key string
		var keyword_value string
		err = rows.Scan(&id, &guild_id, &keyword_key, &keyword_value)
		checkErr(err)
		fmt.Println("id | guild_id | keyword_key | keyword_value ")
		fmt.Printf("%3v | %8v | %8v | %8v\n", id, guild_id, keyword_key, keyword_value)
	}

}*/

type Db struct {
	db *sql.DB
}

func SetDb() *Db {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		"golang", "golang", "discordbot")
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	query, err := ioutil.ReadFile("configDatabase.sql")
	queryString := string(query)
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(queryString); err != nil {
		panic(err)
	}
	database := new(Db)
	database.db = db
	return database
}

func (db *Db) Close() {
	db.db.Close()
}

func (database *Db) GetKeywordValue(guild_id string, key string) (bool, string) {
	db := database.db
	rows, err := db.Query("SELECT * FROM keywords WHERE keyword")
	if checkErrBool(err) {
		fmt.Println("Error working with database")
		return false, "Error running query for get keyword value"
	}

	for rows.Next() {
		var id int
		var guild_id string
		var keyword_key string
		var keyword_value string
		err = rows.Scan(&id, &guild_id, &keyword_key, &keyword_value)
		checkErr(err)
		fmt.Println("id | guild_id | keyword_key | keyword_value ")
		fmt.Printf("%3v | %8v | %8v | %8v\n", id, guild_id, keyword_key, keyword_value)
		return true, keyword_value
	}
	return false, "Result not found"
}

func (database *Db) Insert(guild_id string, key string, value string) bool {
	db := database.db
	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(guild_id,keyword_key,keyword_value) VALUES($1,$2,$3) returning id;", guild_id, key, value).Scan(&lastInsertId)
	if checkErrBool(err) {
		fmt.Println("Error working with database")
		return false
	}
	fmt.Println("last inserted id =", lastInsertId)
	return true
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func checkErrBool(err error) bool {
	if err != nil {
		return true
	}
	return false
}
