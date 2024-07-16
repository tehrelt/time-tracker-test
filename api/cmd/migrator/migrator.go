package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var schemasDir string

func init() {
	flag.StringVar(&schemasDir, "dir", "migrations", "path to schemas dir")
}

func main() {
	flag.Parse()

	fmt.Printf("schemas dir: %s\n", schemasDir)

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	cs := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
	fmt.Printf("cs: %s\n", cs)

	m, err := migrate.New("file://"+schemasDir, cs)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		panic(err)
	}

	fmt.Println("migrations started")
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			fmt.Printf("err: %s\n", err)
			panic(err)
		}
		fmt.Println("no change")
	}
	fmt.Println("migrations done")

}