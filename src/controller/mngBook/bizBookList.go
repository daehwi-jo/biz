package mngBook

import (
	"biz-web/query/commons"
	bizBookMngQuery "biz-web/query/mngBook"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetBookList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	bookListCnt, err := cls.GetSelectDataRequire(bizBookMngQuery.SelectBookGrpListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	bookList, err := cls.GetSelectTypeRequire(bizBookMngQuery.SelectBookGrpList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = bookListCnt[0]["TOTAL_COUNT"]
	result["totalPage"] = utils.GetTotalPage(bookListCnt[0]["TOTAL_COUNT"])

	if bookList == nil { //nil 이나옴 이유는 모름
		result["bookList"] = []string{}
	} else {
		result["bookList"] = bookList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetBookInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)
	grpBookData, err := cls.GetSelectTypeRequire(bizBookMngQuery.SelectBookInfo, params, c)

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if grpBookData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	result := make(map[string]interface{})
	result["grpBookData"] = grpBookData[0]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ModifyBookInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	query, err := cls.SetUpdateParam(bizBookMngQuery.UpdateGrp, params) //쿼리문 작성

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetGrpUserList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	bookList, err := cls.GetSelectTypeRequire(bizBookMngQuery.SelectGrpUserList, params, c)

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	if bookList == nil {
		result["bookList"] = []string{}
	} else {
		result["bookList"] = bookList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func UpdateGrpUser(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	deptCode, err := cls.GetSelectDataRequire(bizBookMngQuery.SelectCompanyDeptCode, params, c) //유저아이디 찾기
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	new := make(map[string]string) //새로운 관리자
	new["grpAuth"] = "0"
	new["userId"] = params["newUserId"]
	new["uId"] = params["newUserId"]
	new["grpId"] = params["grpId"]
	new["companyId"] = params["companyId"]
	new["deptId"] = deptCode[0]["codeId"]
	new["authStat"] = "1"
	new["joinTy"] = "0"

	old := make(map[string]string) //기존 관리자
	old["grpAuth"] = "1"
	old["userId"] = params["userId"]
	old["grpId"] = params["grpId"]

	Query, Query2 := "", ""

	//해당장부에 새로 변경할 유저가 있는지 확인
	grpUserState, err := cls.GetSelectDataRequire(bizBookMngQuery.SelectGrpUserInfo, new, c) //유저 데이터 가져옴
	if grpUserState[0]["state"] == "0" {
		//장부에 유저 추가
		insertGrpUser, err := cls.SetUpdateParam(bizBookMngQuery.InsertGrpUser, new)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		_, err = cls.QueryDB(insertGrpUser) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	companyUserId, err := cls.GetSelectDataRequire(bizBookMngQuery.SearchCompanyUser, new, c) //유저아이디 찾기
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if len(companyUserId) == 0 {
		//회사유저에 인설트
		userData, err := cls.GetSelectDataRequire(bizBookMngQuery.SelectUserInfo, new, c) //유저 데이터 가져옴

		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		new["userNm"] = userData[0]["userNm"]
		new["hpNo"] = userData[0]["hpNo"]

		Query, err = cls.SetUpdateParam(bizBookMngQuery.InsertCompanyUser, new)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		Query2, err = cls.SetUpdateParam(bizBookMngQuery.InsertConnectBook, new) //장부연결
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

	} else {
		//회사유저에 업데이트
		Query, err = cls.SetUpdateParam(bizBookMngQuery.UpdateCompanyUser, new)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		Query2, err = cls.SetUpdateParam(bizBookMngQuery.UpdateConnectBook, new) //장부연결
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	_, err = cls.QueryDB(Query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_, err = cls.QueryDB(Query2) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	oldQuery, err := cls.SetUpdateParam(bizBookMngQuery.UpdateGrpUser, old) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(oldQuery) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	newQuery, err := cls.SetUpdateParam(bizBookMngQuery.UpdateGrpUser, new) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(newQuery) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func DeleteGrp(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	query, err := cls.SetUpdateParam(bizBookMngQuery.DeleteGrp, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_, err = cls.QueryDB(query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

//장부 추기
func CreateCoBookGrp(c echo.Context) error { //이거문제만 해결하면 끝남
	params := cls.GetParamJsonMap(c)

	//새로운 번호 추가
	groupId, err := cls.QueryMapColumn(bizBookMngQuery.SelectCreateGrpSeq, c)
	if groupId == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}
	params["grpId"] = groupId[0]["GRP_ID"]

	//회사정보
	companyInfo, err := cls.GetSelectDataRequire(bizBookMngQuery.SelectCompanyInfo, params, c)
	if companyInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	params["addr"] = companyInfo[0]["ADDR"] + ""
	params["addr2"] = companyInfo[0]["ADDR2"] + ""
	params["busid"] = companyInfo[0]["BUSID"] + ""
	params["ceoNm"] = companyInfo[0]["CEO_NM"] + ""
	params["email"] = companyInfo[0]["EMAIL"] + ""
	params["lat"] = companyInfo[0]["LAT"] + ""
	params["lng"] = companyInfo[0]["LNG"] + ""
	params["tel"] = companyInfo[0]["TEL"] + ""

	//없는값 강제로 삽입
	params["authStat"] = "1"
	params["detailViewYn"] = "Y"
	params["paymentId"] = ""

	//장부생성 쿼리
	Query, err := cls.SetUpdateParam(bizBookMngQuery.InsertCreateGroup, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		//여기서 디비로 삽입 해야함
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	//장부 기본관리자 정보
	userInfo, err := cls.GetSelectDataRequire(bizBookMngQuery.GetDefaultGrpBookMng, params, c)
	params["userId"] = userInfo[0]["USER_ID"]
	params["userNm"] = userInfo[0]["USER_NM"]
	params["dept"] = userInfo[0]["DEPT"]
	params["hpNo"] = userInfo[0]["TEL"]
	params["joinTy"] = "0"
	params["grpAuth"] = "0"
	params["deptId"] = `NULL` //이것도 null로 처리

	//장부유저 등록
	Query2, err := cls.SetUpdateParam(bizBookMngQuery.InsertCreateGroupUser, params) //쿼리문 작성
	fmt.Println(Query2)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		//여기서 디비로 삽입 해야함
		_, err = cls.QueryDB(Query2) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	//회사 장부 등록
	Query3, err := cls.SetUpdateParam(bizBookMngQuery.InsertCompanyBook, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		//여기서 디비로 삽입 해야함
		_, err = cls.QueryDB(Query3) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	//관리자 등록(회사 유저에 새로등록하는 장부아이디로 신규 유저 등록)
	Query4, err := cls.SetUpdateParam(bizBookMngQuery.InsertCompanyUser, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		//여기서 디비로 삽입 해야함
		_, err = cls.QueryDB(Query4) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetCompanyUserList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	list, err := cls.GetSelectDataRequire(bizBookMngQuery.SelectCompanyUser, params, c)

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	if list == nil {
		m["resultData"] = []string{}
	} else {
		m["resultData"] = list
	}

	return c.JSON(http.StatusOK, m)
}

func AddBookManager(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	//쿼리문 만들고
	Query, err := cls.SetUpdateParam(bizBookMngQuery.InsertGrpManager, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		fmt.Println(Query)
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetCompanyBookList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	addQurey := ``

	if params["userId"] != "" {
		addQurey = `AND c.user_id = '#{userId}'`
	}

	result := make(map[string]interface{})

	//쿼리문 만들고
	listCnt, err := cls.GetSelectData(bizBookMngQuery.SelectNotBizCompanyBookCnt+addQurey, params, c) //쿼리문 작성

	list, err := cls.GetSelectData(bizBookMngQuery.SelectNotBizCompanyBook+addQurey+commons.PagingQuery, params, c) //쿼리문 작성

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if list == nil {
		result["grpList"] = []string{}
	} else {
		result["grpList"] = list
	}

	result["totalCnt"] = listCnt[0]["total"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]

	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ModifyBookCompanyId(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	dept, err := cls.GetSelectData(bizBookMngQuery.SelectDeptCodeFirst, params, c)

	grpUserList, err := cls.GetSelectData(bizBookMngQuery.SelectCompanyGrpUser, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	for i := range grpUserList {
		d := make(map[string]string)
		d["companyId"] = params["companyId"]
		d["grpId"] = params["grpId"]
		d["hpNo"] = grpUserList[i]["hpNo"]
		d["useYn"] = grpUserList[i]["useYn"]
		d["userId"] = grpUserList[i]["userId"]
		d["userNm"] = grpUserList[i]["userNm"]
		d["dept"] = dept[0]["deptCode"]

		query, err := cls.SetUpdateParam(bizBookMngQuery.InsertCompanyBookUser, d) //쿼리문 작성
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		} else {
			_, err = cls.QueryDB(query) // 쿼리 실행
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			}
		}
	}

	query2, err := cls.SetUpdateParam(bizBookMngQuery.UpdateBookCompanyId, params) //쿼리문 작성

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(query2) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}
