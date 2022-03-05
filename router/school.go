package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"schoolHelper/structure"
	"schoolHelper/util"

	"github.com/bwmarrin/discordgo"
)

func SchoolRouter(res *discordgo.Session, req *discordgo.MessageCreate, cmd []string) {
	if cmd[1] == "등록" {
		addSchool(res, req, cmd)
	}
	if cmd[1] == "급식" {
		callLunch(res, req)
	}
}

func callLunch(res *discordgo.Session, req *discordgo.MessageCreate) {

}

func addSchool(res *discordgo.Session, req *discordgo.MessageCreate, cmd []string) {
	// 학교 등록 학교명
	schoolName := cmd[2]
	reqPath := "https://open.neis.go.kr/hub/schoolInfo?Type=json&pSize=3&SCHUL_NM=" + url.QueryEscape(schoolName) // url파싱이 자동으로 안되서 미리 함
	resp, err := http.Get(reqPath)

	if err != nil {
		res.ChannelMessageSend(req.ChannelID, "학교를 찾을 수 없습니다.")
		return
	}

	// json parse
	var body structure.SchoolCode
	data, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &body)
	if err != nil || len(body.SchoolInfo) < 1 {
		util.ErrProcesser("api의 응답을 읽지 못했습니다.", err, res, req)
		return
	}

	if len(body.SchoolInfo) > 2 || len(body.SchoolInfo[0].Row) > 1 {
		res.ChannelMessageSend(req.ChannelID, "특정 학교를 찾을 수 없습니다.\n더 자세한 이름을 적어주세요.")
		return
	}
	defer resp.Body.Close()

	// data to information
	var userInfoes []structure.User
	var addSchoolInfo structure.User

	err = util.GetUserData(&userInfoes)
	if err != nil {
		util.ErrProcesser("유저의 정보를 불러오지 못했습니다.", err, res, req)
		return
	}

	// 유저 정보 업데이트
	addSchoolInfo.SchoolCode = body.SchoolInfo[1].Row[0].SdSchulCode
	addSchoolInfo.UserId = req.Author.ID

	userInfoes = append(userInfoes, addSchoolInfo)

	jsonData, err := json.Marshal(userInfoes)
	if err != nil {
		util.ErrProcesser("유저 정보를 인코딩하지 못했습니다.", err, res, req)
		return
	}

	res.ChannelMessageSend(req.ChannelID, "등록을 완료했습니다. "+body.SchoolInfo[1].Row[0].AtptOfcdcScNm+" "+body.SchoolInfo[1].Row[0].SchulNm)
	os.WriteFile("./data/user.json", jsonData, 0644)
}
