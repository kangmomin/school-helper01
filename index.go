package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"schoolHelper/router"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
		fmt.Println(err)
		return
	}

	Tocken := os.Getenv("TOCKEN")
	sc := make(chan os.Signal, 1)

	client, err := discordgo.New("Bot " + Tocken)

	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//route
	client.AddHandler(router.Route)

	client.Identify.Intents = discordgo.IntentsGuildMessages

	err = client.Open()

	if err != nil {
		log.Println(err)
		return
	}
	setActivity(client)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func setActivity(s *discordgo.Session) {
	err := s.UpdateListeningStatus("!설명")
	if err != nil {
		panic(err)
	}
}
