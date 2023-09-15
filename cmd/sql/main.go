package main

import (
	"fmt"
	"os"

	"github.com/dmcclung/pixelparade/db"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	db, err := db.DefaultPostgresConfig.Open()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		panic(err)
	}

	sqlFiles := []string{"models/sql/user.sql", "models/sql/session.sql"}

	for _, sqlFile := range sqlFiles {
		content, err := os.ReadFile(sqlFile)
		if err != nil {
			panic(err)
		}

		sql := string(content)
		_, err = db.Exec(sql)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Done!")
}
