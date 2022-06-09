package mngBook

import (
	"biz-web/query/commons"
	bizBookMngQuery "biz-web/query/mngBook"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/struct"
	"biz-web/src/controller/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
)

func GetBookUserList(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	userListCnt, err := cls.GetSelectDataRequire(bizBookMngQuery.SelectGrpUserMngCnt, params, c) //쿼리문 바꿔야힘
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	userList, err := cls.GetSelectTypeRequire(bizBookMngQuery.SelectGrpUserMng+bizBookMngQuery.OrderBySort(params)+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = userListCnt[0]["TOTAL_COUNT"]
	result["totalPage"] = utils.GetTotalPage(userListCnt[0]["TOTAL_COUNT"])

	if userList == nil {
		result["userList"] = []string{}
	} else {
		result["userList"] = userList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetBookUserInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	userList, err := cls.GetSelectTypeRequire(bizBookMngQuery.SelectGrpCompanyUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if userList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	result := make(map[string]interface{})
	result["userList"] = userList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetGrpBookInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	dataList, err := cls.GetSelectTypeRequire(bizBookMngQuery.SelectGrpBookInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if dataList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	result := make(map[string]interface{})
	result["infoData"] = dataList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func UpdateBookUserDel(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	//사용자 권한 해제 가입되어있을 경우

	if params["rType"] == "U" {
		first, err := cls.SetUpdateParam(bizBookMngQuery.UpdateGroupUserInfo, params) //쿼리문 작성
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		_, err = cls.QueryDB(first)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	second, err := cls.SetUpdateParam(bizBookMngQuery.UpdateCompanyUserDel, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_, err = cls.QueryDB(second) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func UpdateBookUserDisconnect(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	first, err := cls.SetUpdateParam(bizBookMngQuery.UpdateGroupUserInfo, params) //쿼리문 작성

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_, err = cls.QueryDB(first) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

// 2020.03.31 수정
func ModifyBookUserInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})

	//유저 아이디 찾아서 장부에 등록
	if params["userId"] == "" {
		grpUserInfo, err := cls.GetSelectData(bizBookMngQuery.SelectCheckGrpUser, params, c) //쿼리문 작성
		if err != nil {
			cls.Lprintf(1, "[Error] Query SelectCheckGrpUser : %s \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		if len(grpUserInfo) != 0 { //장부 유저 수정시 회원 조회후 회사 회원과 맵핑하고 장부유저 추가 하는 로직
			params["userId"] = grpUserInfo[0]["userId"]
			if grpUserInfo[0]["grpId"] == "" {
				params["joinTy"] = "1"
				params["grpAuth"] = "1"
				params["authStat"] = "0"
				query, err := cls.SetUpdateParam(bizBookMngQuery.InsertCreateGroupUser, params) //쿼리문 작성
				if err != nil {
					cls.Lprintf(1, "[Error] Mapping InsertCreateGroupUser : %s \n", err.Error())
					return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
				}
				_, err = cls.QueryDB(query) // 쿼리 실행
				if err != nil {
					cls.Lprintf(1, "[Error] Query InsertCreateGroupUser : %s \n", err.Error())
					return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
				}
			} else {
				params["userId"] = ""
			}
		}
	}

	query, err := cls.SetUpdateParam(bizBookMngQuery.UpdateCompanyUserInfo, params) //쿼리문 작성
	if err != nil {
		cls.Lprintf(1, "[Error] Mapping UpdateCompanyUserInfo : %s \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(query) // 쿼리 실행
	if err != nil {
		cls.Lprintf(1, "[Error] Query UpdateCompanyUserInfo : %s \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if params["userId"] != "" && params["grpId"] != "" {
		query, err = cls.SetUpdateParam(bizBookMngQuery.UpdateCompanyUserSupportAmt, params) //쿼리문 작성
		if err != nil {
			cls.Lprintf(1, "[Error] Mapping UpdateCompanyUserSupportAmt : %s \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		_, err = cls.QueryDB(query) // 쿼리 실행
		if err != nil {
			cls.Lprintf(1, "[Error] Query UpdateCompanyUserSupportAmt : %s \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
	} else if params["userId"] == "" {
		m["resultCode"] = "98"
		m["resultMsg"] = "이미 등록된 회원 이거나 가입을 하지 않은 회원 입니다."
	} else {
		m["resultCode"] = "98"
		m["resultMsg"] = "사용자 정보가 수정되었습니다.\n (장부 미가입자는 잔여지원금이 갱신 되지 않습니다.)"
	}

	return c.JSON(http.StatusOK, m)
}

// 2020.03.31 수정
func InsertBookUserConnect(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})

	params["joinTy"] = "1"
	params["grpAuth"] = "1"
	params["authStat"] = "1"

	userData, err := cls.GetSelectData(bizBookMngQuery.SelectGrpCompanyUserInfo, params, c)
	if err != nil {
		cls.Lprintf(1, "[Error] Query SelectGrpCompanyUserInfo \n")
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	params["userId"] = userData[0]["USER_ID"]

	query := ""
	if params["userId"] == "" && params["rType"] == "I" { //userId 없고 파라미터로 들어온 타입이 I인 경우 새로 장부연결 Row 생성
		query, err = cls.SetUpdateParam(bizBookMngQuery.InsertConnectBook, params) //쿼리문 작성
	} else { //Update
		if params["userId"] == "" {
			cls.Lprintln(1, "[Error] userId is Null Error")
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		query, err = cls.SetUpdateParam(bizBookMngQuery.UpdateConnectBook, params) //쿼리문 작성
	}

	if err != nil {
		cls.Lprintln(1, "[Error] SetUpdateParam Mapping Error : %s", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_, err = cls.QueryDB(query) // 쿼리 실행
	if err != nil {
		cls.Lprintln(1, "[Error] QueryDB Error : %s", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 지원금 초기금액 세팅
	query, err = cls.SetUpdateParam(bizBookMngQuery.UpdateSupportBalance, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func ParsingExcel(c echo.Context) error {

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	fData, _, err := c.Request().FormFile("upload_file")
	if err != nil {
		cls.Lprintf(1, "[ERROR] %s", err.Error())
		return c.JSON(http.StatusOK, m)
	}

	// xlsx.OpenBinary(bodyBytes)
	file, err := xlsx.OpenReaderAt(fData, c.Request().ContentLength)
	if err != nil {
		cls.Lprintf(1, "[ERROR] %s", err.Error())
		return c.JSON(http.StatusOK, m)
	}

	rData := make(map[string]string)
	var rCnt int

	for _, sheet := range file.Sheets {
		for i, row := range sheet.Rows {
			//fmt.Println(len(row.Cells))
			//fmt.Println(row.Cells[0].Value, " ", row.Cells[1].Value)
			var data string
			for j, cell := range row.Cells {
				if j == 0 {
					data = cell.String()
				} else {
					data = fmt.Sprintf("%s,%s", data, cell.String())
				}
			}

			rData[fmt.Sprintf("%d", i)] = data
			rCnt++
		}
	}

	m["resultData"] = rData
	m["resultCnt"] = rCnt

	return c.JSON(http.StatusOK, m)
}

func InsertBookUser(c echo.Context) error {
	var data _struct.BookUser

	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)

	err := json.Unmarshal(bodyBytes, &data)

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	params := cls.GetParamJsonMap(c) //json을 map으로 컨버팅

	Query := ""

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	for i := 0; i < len(data.UserList); i++ {
		userData := make(map[string]string)
		userData["companyId"] = params["companyId"]
		userData["grpId"] = params["grpId"]
		userData["userNm"] = data.UserList[i].UserNm
		userData["hpNo"] = data.UserList[i].HpNo
		userData["deptId"] = data.UserList[i].DeptId
		userData["joinTy"] = "1"
		userData["grpAuth"] = "1"
		userData["authStat"] = "0"

		userId, err := cls.GetSelectDataRequire(bizBookMngQuery.SearchUserId, userData, c) //유저아이디 찾기
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		if len(userId) == 0 {
			userData["userId"] = "<nil>"
		} else {
			userData["userId"] = userId[0]["userId"] //조회해서 찾아넣음
		}

		companyUserId, _ := cls.GetSelectDataRequire(bizBookMngQuery.SearchCompanyUser, userData, c) //회사장부에 등록되었는지 확인
		if len(companyUserId) == 0 {                                                                 //없을경우 //회사유저 추가
			Query, err = cls.SetUpdateParam(bizBookMngQuery.InsertCompanyUser, userData)
		} else { //있을경우 //회사유저 업데이트
			Query, err = cls.SetUpdateParam(bizBookMngQuery.UpdateCompanyUser, userData)
		}
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		//장부유저 추가
		if userData["userId"] != "<nil>" {
			Query2, err := cls.SetUpdateParam(bizBookMngQuery.InsertCreateGroupUser, userData)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			} else {
				_, err = cls.QueryDB(Query2) // 쿼리 실행
				if err != nil {
					return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
				}
			}
		}
	}

	return c.JSON(http.StatusOK, m)
}

func CreateInviteURL(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	appKey := "AIzaSyB4jaVLPdwWU28SzyKIU8DngpPF5zbFLLg"
	var jsonData map[string]string

	resp, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"dynamicLinkInfo":
										{"domainUriPrefix":"https://darayo.page.link",
										"link":"https://darayo.com/grpAdd?bookId=` + params["grpId"] +
			`","androidInfo":{"androidPackageName":"com.fit.darayo"},
										"iosInfo":{"iosBundleId":"kr.co.darayo",
													"iosAppStoreId":"1158745361"}
										},
						"suffix":{"option":"SHORT"}
						}`).
		Post("https://firebasedynamiclinks.googleapis.com/v1/shortLinks?key=" + appKey) //파이어 베이스로 api 통신

	json.Unmarshal([]byte(resp.String()), &jsonData)
	//fmt.Println(resp)
	//fmt.Println(jsonData["shortLink"])
	//fmt.Println(reflect.TypeOf(jsonData["shortLink"]))

	params["inviteLink"] = jsonData["shortLink"]

	//파이어 베이스로 링크 만든거 삽입
	query, err := cls.SetUpdateParam(bizBookMngQuery.UpdateGrpBookInvite, params) //쿼리문 작성
	_, err = cls.QueryDB(query)

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["inviteLink"] = params["inviteLink"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetDeptCode(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	deptCode, err := cls.GetSelectDataRequire(bizBookMngQuery.SelectDeptCode, params, c)

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["deptCode"] = deptCode

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}
