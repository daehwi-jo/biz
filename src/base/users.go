package base

import (
	usersql "biz-web/query"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"bytes"
	"encoding/xml"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

/*
작성일 : 21.06.28
작성자 : 김형곤
*/

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf //??
var lprintf func(int, string, ...interface{}) = cls.Lprintf               //??

type ResultMsgMap struct {
	XMLName   xml.Name `xml:"map"`
	Id        string   `xml:"id,attr"`
	DetailMsg string   `xml:"detailMsg"`
	Msg       string   `xml:"msg"`
	Code      string   `xml:"code"`
	Result    string   `xml:"result"`
}

type Map struct {
	XMLName             xml.Name     `xml:"map"`
	ResultMsgMap        ResultMsgMap `xml:"map"`
	TrtEndCd            string       `xml:"trtEndCd"`
	SmpcBmanEnglTrtCntn string       `xml:"smpcBmanEnglTrtCntn"`
	NrgtTxprYn          string       `xml:"nrgtTxprYn"`
	SmpcBmanTrtCntn     string       `xml:"smpcBmanTrtCntn"`
	TrtCntn             string       `xml:"trtCntn"`
}

// 사업자 번호 조회
func BizNumCheck(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	bizNum := params["bizNum"]

	url := "https://teht.hometax.go.kr/wqAction.do?actionId=ATTABZAA001R08&screenId=UTEABAAA13&popupYn=false&realScreenId="
	xmlData := "<map id='ATTABZAA001R08'>" +
		"<pubcUserNo/>" +
		"<mobYn>N</mobYn>" +
		"<inqrTrgtClCd>1</inqrTrgtClCd>" +
		"<txprDscmNo>" + bizNum + "</txprDscmNo>" +
		"<dongCode>05</dongCode" +
		"><psbSearch>Y</psbSearch>" +
		"<map id='userReqInfoVO'/>" +
		"</map>"
	buf := bytes.NewBufferString(xmlData)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		//fmt.Println(err)
	}
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var rMap Map
	xml.Unmarshal(body, &rMap)

	checkResult := rMap.SmpcBmanTrtCntn

	if strings.Replace(checkResult, " ", "", -1) != "등록되어있는사업자등록번호입니다." {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", checkResult))
	}

	idChk, err := cls.GetSelectData(usersql.SelectLoginIdDupCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	idCnt, _ := strconv.Atoi(idChk[0]["loginIdCnt"])
	if idCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("01", "이미 가입된 사업자 번호입니다."))
	}

	bizNumCnt, err := cls.GetSelectData(usersql.SelectBizNumDupCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	bizCnt, _ := strconv.Atoi(bizNumCnt[0]["bizCnt"])

	println(bizCnt)
	if bizCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("01", "이미 가입된 사업자 번호입니다."))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

// 회원 가입 1단계
func SetUserJoin(c echo.Context) error {

	dprintf(4, c, "call setUserJoin\n")

	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectData(usersql.SelectCreatUserSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "유저아이디 생성 실패(2)"))
	}
	userId := resultData[0]["newUserId"]
	params["userId"] = userId
	// 유저 가입  TRNAN 시작
	tx, err := cls.DBc.Begin() //디비 시작
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -유저 가입(setUserJoin)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	socialType := params["socialType"]
	socialToken := params["socialToken"]
	loginPw := params["loginPw"]
	loginId := params["loginId"]

	if socialType == "kakao" {
		params["kakaoKey"] = socialToken
		params["kakaoPw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	} else if socialType == "naver" {
		params["naverKey"] = socialToken
		params["naverPw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	} else if socialType == "apple" {
		params["appleKey"] = socialToken
		params["applePw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	}

	termsOfBenefit := params["termsOfBenefit"]
	pushYn := "N"
	if termsOfBenefit == "Y" {
		pushYn = "Y"
	}

	// 유저 생성
	params["userTy"] = "1"
	params["atLoginYn"] = "Y"
	params["pushYn"] = pushYn
	UserCreateQuery, err := cls.SetUpdateParam(usersql.InserCreateUser, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UserCreateQuery parameter fail"))
	}

	_, err = tx.Exec(UserCreateQuery)
	dprintf(4, c, "call set Query (%s)\n", UserCreateQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UserCreateQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	TermsCreateQuery, err := cls.SetUpdateParam(usersql.InsertTermsUser, params) //쿼리문 생성
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "TermsCreateQuery parameter fail"))
	}

	_, err = tx.Exec(TermsCreateQuery)                       //쿼리 실행
	dprintf(4, c, "call set Query (%s)\n", TermsCreateQuery) //쿼리 올려줬다는 로그
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", TermsCreateQuery, err) // 에러 발생시
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	//c = cls.SetLoginJWT(c, userId)

	userData := make(map[string]interface{})
	userData["userId"] = userId
	userData["bizNum"] = loginId

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}
