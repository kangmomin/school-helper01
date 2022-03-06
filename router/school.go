package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	menu := strings.ReplaceAll(data.MealServiceDietInfo[0].Row[0].DdishNm, ".", "")
	menu = strings.ReplaceAll(menu, "/^[1-9][0-9]*$/", "")
	menu = strings.ReplaceAll(menu, " ", "")
	menu = strings.ReplaceAll(menu, ":", "")
	menuList := strings.Split(menu, "<br/>")

	// make embed message
	embed := embed.NewEmbed()
	embed.SetTitle(date.Format("2006-01-02") + "일")
	embed.SetDescription(data.MealServiceDietInfo[0].Row[0].SchulNm)
	for idx := range menuList {
		if len(menuList) < idx*3+2 {
			break
		}

		embed.AddField(menuList[idx*3], menuList[idx*3+1]+"\n"+menuList[idx*3+2])
	}

	embed.InlineAllFields().
		Truncate().
		SetColor(0x7FD5E9)

	res.ChannelMessageSendEmbed(req.ChannelID, embed.MessageEmbed)
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

	if len(body.SchoolInfo) > 2 || len(body.SchoolInfo[1].Row) > 1 {
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
