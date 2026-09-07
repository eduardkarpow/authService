package main

import (
	"authService/internal/config"
	"authService/internal/database"
	"flag"
	"log"
)

func main() {
	conf := config.Load()
	migrateFlag := flag.String("migrate", "", "Run database migrations: up/down")
	flag.Parse()

	db, err := database.Connect(conf.DBConnStr)
	if err != nil {
		log.Fatal(err)
		return
	}
	if *migrateFlag != "" {
		switch *migrateFlag {
		case "up":
			if err := db.RunMigrationsUp("auth"); err != nil {
				log.Fatal(err)
			}
			log.Println("Database migrations up")
		case "down":
			if err := db.RunMigrationsDown("auth"); err != nil {
				log.Fatal(err)
			}
			log.Println("Database migrations down")
		default:
			log.Fatal("Unsupported migrate flag")
		}
		return
	}
}
