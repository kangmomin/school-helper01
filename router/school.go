package router

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"schoolHelper/structure"

	"github.com/bwmarrin/discordgo"
)

func SchoolRouter(res *discordgo.Session, req *discordgo.MessageCreate, cmd []string) {
	if cmd[1] == "등록" {
		addSchool(res, req, cmd)
	}
	if cmd[1] == "급식" {
		callCafeteria(res, req)
	}
}

func callCafeteria(res *discordgo.Session, req *discordgo.MessageCreate) {

}

func addSchool(res *discordgo.Session, req *discordgo.MessageCreate, cmd []string) {
	// 학교 등록 학교명
	schoolName := cmd[2]

	resp, err := http.Get("https://open.neis.go.kr/hub/schoolInfo?Type=json&pIndex=1&pSize=100&SCHUL_NM=" + schoolName)

	if err != nil {
		log.Println(err)
		res.ChannelMessageSend(req.ChannelID, "학교를 찾을 수 없습니다.")
		return
	}

	var body structure.SchoolCode
	err = json.NewDecoder(resp.Body).Decode(&body)

	if err != nil {
		log.Println(err)
		res.ChannelMessageSend(req.ChannelID, "api의 응답을 읽지 못했습니다.")
		return
	}

	resp.Body.Close()
	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		res.ChannelMessageSend(req.ChannelID, "error during marshalling")
		return
	}
	os.WriteFile("./data/user.json", jsonData, 0644)
}
