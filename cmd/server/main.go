package main

import (
	"authService/internal/config"
	"authService/internal/database"
	"authService/internal/service"
	"authService/internal/transport"
	"flag"
	"log"
	"net"
	"os"

	authv1 "github.com/eduardkarpow/auth-proto/gen/go/auth/v1"
	"google.golang.org/grpc"
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
	userService := service.NewUserService(db.Pool, conf)
	grpcHandler := transport.NewUserGrpcHandler(userService)
	stop := make(chan os.Signal, 1)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServer(grpcServer, grpcHandler)
	go func() {
		log.Println("Listening on :8080")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server gracefully stopped")
}
