package DB

import (
	"database/sql"
	"log"
)

func createUrlsTable(db *sql.DB) {
	st, err0 := db.Prepare("CREATE TABLE IF NOT EXISTS urls(shorten TEXT, url TEXT, expire DATE, cntr INT(11) DEFAULT \"0\")")
	if err0 != nil {
		panic(err0.Error())
	}
	_, err1 := st.Exec()
	if err1 != nil {
		panic(err1.Error())
	}

}

func ConnectDB(username, password, address, dbName string) *sql.DB {
	// Setup db connection
	db, err := sql.Open("mysql",
		username+":"+password+"@tcp("+address+")/"+dbName)
	if err != nil {
		log.Fatal(err)
	}
	createUrlsTable(db)
	return db
}

func Add(shorten string, url string, date string, db *sql.DB) {
	add, err := db.Prepare("INSERT INTO urls(shorten, url, expire) VALUES(\"" + shorten + "\",\"" + url + "\",'" + date + "')")
	if err != nil {
		log.Println(err.Error())
	}
	defer add.Close()
	if err != nil {
		log.Fatal(err)
	}
	add.Exec()
}

func Exist(shorten string, db *sql.DB) bool {
	var id int
	err := db.QueryRow("SELECT shorten FROM urls WHERE shorten = '" + shorten + "'").Scan(&id)
	if !(err != nil && err == sql.ErrNoRows) {
		return true
	} else {
		return false
	}
}

func Select(selected string, id string, db *sql.DB) string {
	var value string
	er := db.QueryRow("SELECT " + selected + " FROM urls WHERE shorten = '" + id + "'").Scan(&value)
	if er != nil && er != sql.ErrNoRows {
		log.Fatal(er)
	}
	return value
}

func Used(shorten string, db *sql.DB) {
	var value int
	er := db.QueryRow("SELECT cntr FROM urls WHERE shorten = '" + shorten + "'").Scan(&value)
	if er != nil && er != sql.ErrNoRows {
		log.Fatal(er)
	}
	value++
	insForm, err := db.Prepare("UPDATE urls SET cntr=? WHERE shorten=?")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(value, shorten)
}
