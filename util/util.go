package util

import (
	"schoolHelper/structure"
)

// 유저의 인덱스 번호를 찾아줌
func FindUser(userInfoes []structure.User, id string) (idx int, isFind bool) {
	idx = 0
	for idx < len(userInfoes) {
		if userInfoes[idx].UserId == id {
			return idx, true
		}
		idx++
	}
	return idx, false
}
