package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"schoolHelper/structure"
	"schoolHelper/util"
	"strings"
	"time"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

func SchoolRouter(res *discordgo.Session, req *discordgo.MessageCreate, cmd []string) {
	if cmd[1] == "등록" {
		addSchool(res, req, cmd)
	}
	if cmd[1] == "급식" {
		callLunch(res, req)
	}
	if cmd[1] == "취소" {
		deleteSchool(res, req, cmd)
	}
}

func callLunch(res *discordgo.Session, req *discordgo.MessageCreate) {
	var userInfoes []structure.User
	util.GetUserData(&userInfoes)

	// 유저(커맨드 작성자)의 아이디 위치
	idx, isFind := util.FindUser(userInfoes, req.Author.ID)
	if !isFind {
		res.ChannelMessageSend(req.ChannelID, "아직 학교가 등록되지 않았습니다. !학교 등록")
		return
	}

	date := time.Now()
	if date.Hour() > 12 {
		date = date.AddDate(0, 0, 1)
	}
	reqPath := "https://open.neis.go.kr/hub/mealServiceDietInfo?Type=json&MLSV_YMD=" + date.Format("20060102") + "&ATPT_OFCDC_SC_CODE=" + url.QueryEscape(userInfoes[idx].AtptOfcdcScCode) + "&SD_SCHUL_CODE=" + url.QueryEscape(userInfoes[idx].SchoolCode) // url파싱이 자동으로 안되서 미리 함
	resp, err := http.Get(reqPath)
	if err != nil {
		util.ErrProcesser("api가 응답하지 않습니다.", err, res, req)
		return
	}

	var data structure.Lunch
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.ErrProcesser("api가 응답하지 않습니다.", err, res, req)
		return
	}

	json.Unmarshal(byteValue, &data)
	if len(data.MealServiceDietInfo) < 1 {
		res.ChannelMessageSend(req.ChannelID, "오늘은 급식이 없습니다.")
		return
	}

	// menu string data parse
	menu := strings.ReplaceAll(data.MealServiceDietInfo[1].Row[0].DdishNm, ".", "")
	reg := regexp.MustCompile("[0-9]+")
	menu = reg.ReplaceAllString(menu, "")
	menu = strings.ReplaceAll(menu, " ", "")
	menu = strings.ReplaceAll(menu, ":", "")
	menuList := strings.Split(menu, "<br/>")

	// make embed message
	embed := embed.NewEmbed()
	embed.SetTitle(date.Format("2006-01-02") + "일")
	embed.SetDescription(data.MealServiceDietInfo[1].Row[0].SchulNm)

	// menu list string version
	var menuListUp string
	for idx := 0; idx < len(menuList); idx++ {
		menuListUp += "\n" + menuList[idx]
	}
	embed = embed.AddField("중식", menuListUp)

	embed.InlineAllFields().
		Truncate().
		SetColor(0x7FD5E9)

	res.ChannelMessageSendEmbed(req.ChannelID, embed.MessageEmbed)
}

func addSchool(res *discordgo.Session, req *discordgo.MessageCreate, cmd []string) {
	if len(cmd) < 3 {
		return
	}
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

	if len(body.SchoolInfo) > 2 || len(body.SchoolInfo[1].Row) > 1 {
		res.ChannelMessageSend(req.ChannelID, "특정 학교를 찾을 수 없습니다.\n더 자세한 이름을 적어주십시오.\n띄어쓰기는 빼주십시오.")
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

	_, isFind := util.FindUser(userInfoes, req.Author.ID)
	if isFind {
		res.ChannelMessageSend(req.ChannelID, "이미 등록된 학교가 있습니다.")
		return
	}

	// 유저 정보 업데이트
	addSchoolInfo.SchoolCode = body.SchoolInfo[1].Row[0].SdSchulCode
	addSchoolInfo.AtptOfcdcScCode = body.SchoolInfo[1].Row[0].AtptOfcdcScCode
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

func deleteSchool(res *discordgo.Session, req *discordgo.MessageCreate, cmd []string) {
	if len(cmd) < 2 {
		return
	}

	file, err := os.Open("./data/user.json")
	if err != nil {
		util.ErrProcesser("유저 정보를 가저오지 못하였습니다.", err, res, req)
	}

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		util.ErrProcesser("유저 정보를 가저오지 못하였습니다.", err, res, req)
	}

	var users []structure.User

	json.Unmarshal(byteValue, &users)
	for i := 0; i < len(users); i++ {
		if users[i].UserId == req.Author.ID {
			users = append(users[:i], users[i+1:]...)
			data, _ := json.Marshal(users)
			ioutil.WriteFile("./data/user.json", data, 0644)
			res.ChannelMessageSend(req.ChannelID, "등록 되어있던 학교를 취소하였습니다.")
			break
		}
	}

	res.ChannelMessageSend(req.ChannelID, "아직 학교가 등록되지 않았습니다. !학교 등록")
}
