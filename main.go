package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/raspiantoro/temporama/command"
	_ "github.com/raspiantoro/temporama/memstore"
	"github.com/raspiantoro/temporama/resp"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "6379"
	}

	server := resp.NewServer("0.0.0.0", port)
	server.Handler(command.Registers())

	err := server.ServeAndListen()
	if err != nil {
		log.Fatalln("failed to listen to network address: ", err)
		return
	}

	log.Println("bye")
}
