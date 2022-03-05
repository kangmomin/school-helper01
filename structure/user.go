package structure

type User struct {
	UserId     string `json:"userId"`
	SchoolCode string `json:"schoolCode"`
	// 교육청 코드
	AtptOfcdcScCode string `json:"ATPT_OFCDC_SC_CODE"`
}
