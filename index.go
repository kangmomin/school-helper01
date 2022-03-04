package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"schoolHelper/router"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Tocken = os.Getenv("TOCKEN")

func main() {
	sc := make(chan os.Signal, 1)

	client, err := discordgo.New("BOT" + Tocken)
	defer client.Close()

	if err != nil {
		log.Fatalln(err)
		fmt.Println(err)
		return
	}

	//route
	client.AddHandler(router.Route)

	client.Identify.Intents = discordgo.IntentsGuildMessages

	err = client.Open()

	if err != nil {
		log.Println(err)
		return
	}

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
