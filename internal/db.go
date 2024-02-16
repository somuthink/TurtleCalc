package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}

var glob_db = CreateAndOpen("tutrle_calc")

func get_servers() int {
	servers, _ := strconv.Atoi(os.Getenv("MACHINES"))
	return servers
}

func dsn(dbName string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?multiStatements=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		dbName,
	)
}

func CreateAndOpen(name string) *DB {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	db, err = sql.Open("mysql", dsn(name))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(fmt.Sprintf("USE %s", name))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS `%s`.`problems` (`id` TINYINT NOT NULL AUTO_INCREMENT , `text` TINYTEXT NOT NULL , `interm_val` INT NOT NULL DEFAULT '0' , `answer` INT NOT NULL DEFAULT '0' , `operations_left` TINYINT NOT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB;",
			name,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS `%s`.`operations` (`id` TINYINT NOT NULL AUTO_INCREMENT , `operation` CHAR(1) NOT NULL , `exec_time` INT NOT NULL DEFAULT '200' , PRIMARY KEY (`id`)) ENGINE = InnoDB;",
			name,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	var cnt int
	db.QueryRow("SELECT COUNT(*) FROM `operations`").Scan(&cnt)

	if cnt == 0 {
		db.Exec("INSERT IGNORE INTO `operations` (operation) VALUES ('+'),('-'),('*'),('/') ")
	}

	_, err = db.Exec(
		fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS `%s`.`servers` (`id` TINYINT NOT NULL AUTO_INCREMENT , `operation` TINYTEXT NULL DEFAULT NULL , PRIMARY KEY (`id`)) ENGINE = InnoDB;",
			name,
		),
	)

	db.QueryRow("SELECT COUNT(*) FROM `servers`").Scan(&cnt)

	if cnt == 0 {

		stmt, err := db.Prepare("INSERT INTO `servers` (operation) VALUES ('nothing')")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		for i := 0; i < get_servers(); i++ {
			_, err := stmt.Exec()
			if err != nil {
				log.Fatal(err)
			}
		}

	}

	return &DB{db: db}
}
