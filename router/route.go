package router

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Route(res *discordgo.Session, req *discordgo.MessageCreate) {
	if req.Author.ID == res.State.User.ID {
		return
	}

	// 만약 !로 시작하지 않는다면(cmd가 아니면) return
	if cmd := strings.Split(req.Content, ""); cmd[0] != "!" {
		return
	}

	// 명령어를 구별해내기 위한 문자열 자르기
	msg := strings.Split(req.Content, "")

	if len(msg) < 2 {
		return
	}

	// 문자열의 시작이 해당 봇의 명령어가 아니라면
	if msg[0] != "!" {
		return
	}
	// 한글자씩 나뉜 문자열을 첫 자를 제외하고 다시 ""로 붙이고 붙인 값에서 스플릿;;
	cmd := strings.Split(strings.Join(msg[1:], ""), " ")

	if cmd[0] == "학교" {
		SchoolRouter(res, req, cmd)
	}
}
