package structure

// 작성 후기 시발 개 뭐라는겨
type SchoolCode struct {
	SchoolInfo []struct {
		Head   []int `json:"list_total_count"`
		Result struct {
			Code    string `json:"CODE"`
			Message string `json:"MESSAGE"`
		}
	} `json:"schoolInfo"`

	Row []struct {
		// 교육청 코드
		ATPT_OFCDC_SC_CODE string
		// 교육청 이름
		ATPT_OFCDC_SC_NM string
		// 학교 코드
		SD_SCHUL_CODE string
		// 학교 이름
		SCHUL_NM string
		// 학교 영어 이름
		ENG_SCHUL_NM string
		// 초, 중, 고 ex)고등학교
		SCHUL_KND_SC_NM string
		// 시 도
		LCTN_SC_NM string

		JU_ORG_NM string
		// 사립 공립
		FOND_SC_NM string
		// 우편번호
		ORG_RDNZC string
		// 주소
		ORG_RDNMA string

		ORG_RDNDA string
		// 학교 전번
		ORG_TELNO string
		// 홈페이지
		HMPG_ADRES string
		// 학교 특성 ex) 남여공학
		COEDU_SC_NM string
		// 팩스 번호
		ORG_FAXNO string
		// ex) 특성화고
		HS_SC_NM string

		INDST_SPECL_CCCCL_EXST_YN string

		// 계열 ex) 전문계
		HS_GNRL_BUSNS_SC_NM string

		SPCLY_PURPS_HS_ORD_NM string
		ENE_BFE_SEHF_SC_NM    string
		// 학교 시간 ex) 주간
		DGHT_SC_NM string
		// 개교일
		FOND_YMD string
		// 개교일 2
		FOAS_MEMRD string
		// 마지막 수정일
		LOAD_DTM string
	} `json:"row"`
}
