package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"schoolHelper/logger"
	"schoolHelper/structure"

	"github.com/bwmarrin/discordgo"
)

func GetUserData(userInfoes *[]structure.User) error {
	data, err := os.Open("./data/user.json")
	if err != nil {
		return err
	}
	byteValue, _ := ioutil.ReadAll(data)
	json.Unmarshal(byteValue, &userInfoes)
	return nil
}

func ErrProcesser(message string, err error, res *discordgo.Session, req *discordgo.MessageCreate) {
	log := logger.Logger

	log.Println(err)
	res.ChannelMessageSend(req.ChannelID, message)
}
