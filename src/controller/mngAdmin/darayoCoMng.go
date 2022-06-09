package mngAdmin

import (
	"biz-web/query/commons"
	darayoquery "biz-web/query/mngAdmin"
	bizMngQuery "biz-web/query/mngUser"
	"biz-web/src/controller"
	Daily "biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetBizCompanyList(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	whereQuery := ``

	if params["useYn"] != "" {
		whereQuery += ` AND a.USE_YN = '#{useYn}'`
	}

	if params["searchKey"] == "companyNm" {
		whereQuery += ` AND a.COMPANY_NM LIKE '%#{searchKeyword}%'`
	}

	orderQuery := `order by a.COMPANY_ID desc`

	companyCnt, err := cls.GetSelectDataRequire(darayoquery.SelectCompanyCnt+whereQuery+orderQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	companyList, err := cls.GetSelectTypeRequire(darayoquery.SelectCompanyList+whereQuery+orderQuery+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = companyCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(companyCnt[0]["TOTAL_COUNT"])

	if companyList == nil {
		result["companyList"] = []string{}
	} else {
		result["companyList"] = companyList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetBIZCompanyCsInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	companyInfo, err := cls.GetSelectType(darayoquery.SelectCompanyInfo, params, c) //장부 정보
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "SelectCompanyInfo"+err.Error()))
	}

	managerInfo, err := cls.GetSelectTypeRequire(bizMngQuery.SelectManagerInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "SelectManagerInfo"+err.Error()))
	}

	m := make(map[string]interface{})
	result := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	if companyInfo != nil {
		result["companyInfo"] = companyInfo[0]
	} else {
		result["companyInfo"] = nil
	}

	if managerInfo != nil {
		result["managerInfo"] = managerInfo[0]
	} else {
		result["managerInfo"] = nil
	}

	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetCompanyIdList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	companyList, err := cls.GetSelectData(darayoquery.SelectCompanyId, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData := make(map[string]interface{})
	resultData["companyList"] = companyList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

func GetCompanyManager(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	resultData := make(map[string]interface{})
	list, err := cls.GetSelectDataRequire(darayoquery.SelectCompanyManagers, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if list == nil {
		resultData["managerList"] = make(map[string]interface{})
	} else {
		resultData["managerList"] = list
	}

	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

func GetCompanyUsers(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	m := make(map[string]interface{})
	m["resultCode"] = "00"

	m["resultMsg"] = "응답 성공"

	whereQuery := ``
	if params["searchKey"] == "userNm" {
		whereQuery = `and B.USER_NM LIKE '%#{searchKeyword}%'
`
	} else {
		whereQuery = `and B.HP_NO LIKE '%#{searchKeyword}%'
`
	}

	listCnt, err := cls.GetSelectDataRequire(darayoquery.SelectCompanyUsersCnt+whereQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	list, err := cls.GetSelectDataRequire(darayoquery.SelectCompanyUsers+whereQuery+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData := make(map[string]interface{})
	if list == nil {
		resultData["CUserList"] = []string{}
	} else {
		resultData["CUserList"] = list
	}

	resultData["totalCnt"] = listCnt[0]["total"]
	resultData["pageSize"] = params["pageSize"]
	resultData["pageNo"] = params["pageNo"]
	resultData["totalPage"] = utils.GetTotalPage(listCnt[0]["total"])
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

func AddCompanyManager(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	Query, err := cls.GetQueryJson(darayoquery.InsertCompanyManagers, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}
	if params["authorCd"] == "CM" {
		Query2, err := cls.SetUpdateParam(darayoquery.UpdateCompanyManagerAuthorAllBM, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		} else {
			_, err = cls.QueryDB(Query2) // 쿼리 실행
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			}
		}
	}

	return c.JSON(http.StatusOK, m)
}

func ModCompanyManager(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	Query, err := cls.SetUpdateParam(darayoquery.UpdateCompanyManagerAuthorCM, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	Query2, err := cls.SetUpdateParam(darayoquery.UpdateCompanyManagerAuthorAllBM, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query2) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	return c.JSON(http.StatusOK, m)
}

// companyId = 기업 아이디
// userId = 기념일 충전 기업 직원 아이디
func CompayPriceDay(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	Daily.CompayPriceDay(true, params["companyId"], params["userId"])

	return c.JSON(http.StatusOK, nil)

}

func TransBizUse(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})

	//0. selectData 비즈전환을 위한 데이터
	selectData, err := cls.GetSelectData(darayoquery.SelectDataForTransBiz, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "회사정보 오류 :"+err.Error()))
	}
	if selectData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "회사데이터 없음 :"+err.Error()))
	}

	sysFlag := true
	checkSys, err := cls.GetSelectData(darayoquery.SelectCheckLoginIdInSysInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "SelectCheckLoginIdInSysInfo :"+err.Error()))
	}
	if checkSys != nil {
		sysFlag = false
	}

	params["grpId"] = selectData[0]["grpId"]
	params["loginId"] = selectData[0]["loginId"]
	params["loginPw"] = selectData[0]["loginPw"]
	params["userNm"] = selectData[0]["userNm"]
	params["codeId"] = selectData[0]["codeId"]

	selectGrpUser, err := cls.GetSelectData(darayoquery.SelectGrpUsers, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "회사정보 오류 :"+err.Error()))
	}

	// TRNAN 시작
	tx, err := cls.DBc.Begin() //디비 시작
	if err != nil {
		//return "5100", errors.New("begin error")
	}
	// 트랜젝션 롤백
	defer func() {
		if err != nil {
			// transaction rollback
			tx.Rollback()
		}
	}()

	//1. 기업 부서 데이터 없으면 데이터 추가
	params["codeTy"] = "0"
	params["codeInfo"] = "본사"
	companyDept, err := cls.GetSelectData(darayoquery.SelectCompanyDept, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if companyDept == nil {
		InsertCompanyDept, err := cls.SetUpdateParam(darayoquery.InsertCompanyDept, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyDept "+err.Error()))
		} else {
			data, err := tx.Exec(InsertCompanyDept)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyDept "+err.Error()))
			}
			codeId, formatErr := data.LastInsertId()
			if formatErr != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyDept codeId formatErr "+err.Error()))
			}
			params["codeId"] = strconv.FormatInt(codeId, 10)
		}
	}

	//2. SysUserInfo 데이터 추가
	params["authorCd"] = "CM"
	params["modifyBy"] = "<null>"
	//sysUserPw 설정
	if selectData[0]["appleKey"] != "" && selectData[0]["naverKey"] != "" && selectData[0]["kakaoKey"] != "" {
		params["loginPw"] = params["loginId"]
	}

	if sysFlag {
		InsertSysUserInfo, err := cls.SetUpdateParam(darayoquery.InsertSysUserInfo, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertSysUserInfo "+err.Error()))
		}
		_, err = tx.Exec(InsertSysUserInfo)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertSysUserInfo "+err.Error()))
		}
	}

	//3. InsertCompanyBook 데이터 추가 -> 모카에서 데이터 추가중
	//InsertCompanyBook, err := cls.SetUpdateParam(darayoquery.InsertCompanyBook, params)
	//if err != nil {
	//	return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyBook "+err.Error()))
	//}
	//_, err = tx.Exec(InsertCompanyBook)
	//if err != nil {
	//	return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyBook "+err.Error()))
	//}

	//3. InsertCompanyUser 데이터 추가
	userData := map[string]string{}
	for _, user := range selectGrpUser {
		userData["companyId"] = params["companyId"]
		userData["grpId"] = params["grpId"]
		userData["userId"] = user["userId"]
		userData["userNm"] = user["userNm"]
		userData["hpNo"] = user["hpNo"]
		userData["useYn"] = user["useYn"]
		userData["codeId"] = params["codeId"]

		InsertCompanyUser, err := cls.SetUpdateParam(darayoquery.InsertCompanyUser, userData)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyUser "+err.Error()))
		}
		_, err = tx.Exec(InsertCompanyUser)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyUser "+err.Error()))
		}

		//4. InsertCompanyUser 매니저 추가
		if user["grpAuth"] == "0" {
			manager := map[string]string{}
			manager["companyId"] = params["companyId"]
			manager["codeInfo"] = params["codeInfo"]
			manager["authorCd"] = params["authorCd"]
			manager["userId"] = user["userId"]
			manager["email"] = user["email"]
			manager["hpNo"] = user["hpNo"]

			InsertCompanyManager, err := cls.SetUpdateParam(darayoquery.InsertCompanyManager, manager)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyUser "+err.Error()))
			}
			_, err = tx.Exec(InsertCompanyManager)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyManager "+err.Error()))
			}
		}
	}

	params["useYn"] = "Y"
	UpdateCompanyUseYn, err := cls.SetUpdateParam(darayoquery.UpdateCompanyUseYn, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "UpdateCompanyUseYn "+err.Error()))
	}
	_, err = tx.Exec(UpdateCompanyUseYn)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "UpdateCompanyUseYn "+err.Error()))
	}

	//트랜잭션 커밋
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m["resultCode"] = "00"
	if !sysFlag {
		m["resultMsg"] = "장부장의 login 아이디로 sysUser 존재함"
	} else {
		m["resultMsg"] = "응답 성공"
	}

	return c.JSON(http.StatusOK, m)
}

func AddCompany(c echo.Context) error { //장부생성 등록 안만들었음....
	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})

	//0. SelectDataForAddCompany ,  필요한 데이터 추출
	selectData, err := cls.GetSelectData(darayoquery.SelectDataForAddCompany, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "SelectDataForAddCompany "+err.Error()))
	}
	if selectData == nil { //등록되지 않은 사용자
		m["resultCode"] = "99"
		m["resultMsg"] = "담당자 정보 오류 (입력된 정보로 가입된 사용자가 없습니다.)"
		return c.JSON(http.StatusOK, m)
	}

	sysFlag := true
	checkSys, err := cls.GetSelectData(darayoquery.SelectCheckLoginIdInSysInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "SelectCheckLoginIdInSysInfo :"+err.Error()))
	}
	if checkSys != nil {
		sysFlag = false
	}

	params["grpId"] = selectData[0]["grpId"]
	params["hpNo"] = selectData[0]["hpNo"]
	params["userId"] = selectData[0]["userId"]
	params["userNm"] = selectData[0]["userNm"]
	params["email"] = selectData[0]["email"]
	params["codeId"] = selectData[0]["codeId"]
	params["companyId"] = selectData[0]["companyId"]
	params["loginId"] = selectData[0]["loginId"]
	if params["kakaoKey"] == "" && params["naverKey"] == "" && params["appleKey"] == "" {
		params["loginPw"] = selectData[0]["loginId"]
	} else {
		params["loginPw"] = selectData[0]["loginPw"]
	}

	tx, err := cls.DBc.Begin() //디비 시작
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	defer func() {
		if err != nil {
			// transaction rollback
			tx.Rollback()
		}
	}()

	//회사
	//1. b_company 데이터 삽입
	InsertCompany, err := cls.SetUpdateParam(darayoquery.InsertCompany, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompany "+err.Error()))
	} else {
		data, err := tx.Exec(InsertCompany)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompany "+err.Error()))
		}
		companyId, formatErr := data.LastInsertId()
		if formatErr != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompany companyId formatErr "+err.Error()))
		}
		params["companyId"] = strconv.FormatInt(companyId, 10)
	}

	//2. b_company_code 부서 추가
	params["codeTy"] = "0"
	params["codeInfo"] = "본사"
	InsertCompanyDept, err := cls.GetQueryJson(darayoquery.InsertCompanyDept, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyDept "+err.Error()))
	} else {
		data, err := tx.Exec(InsertCompanyDept)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyDept "+err.Error()))
		}
		codeId, formatErr := data.LastInsertId()
		if formatErr != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyDept codeId formatErr "+err.Error()))
		}
		params["codeId"] = strconv.FormatInt(codeId, 10)
	}

	//3. b_company_manager 회사 장부 매니저 추가
	params["authorCd"] = "CM"
	InsertCompanyManager, err := cls.SetUpdateParam(darayoquery.InsertCompanyManager, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyManager "+err.Error()))
	}
	_, err = tx.Exec(InsertCompanyManager)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyManager "+err.Error()))
	}

	//장부
	//4. priv_grp_info 장부 생성
	params["grpNm"] = params["companyNm"] + "_장부"
	params["addr"] = ""
	params["addr2"] = ""
	params["authStat"] = "1"
	params["limitYn"] = "N"
	params["limitAmt"] = "0"
	params["limitDayAmt"] = "0"
	params["detailViewYn"] = "Y"
	params["grpTypeCd"] = "1"
	params["supportAmt"] = "0"
	params["grpPayTy"] = "1"

	InsertGrp, err := cls.SetUpdateParam(darayoquery.InsertGrp, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertGrp "+err.Error()))
	}
	_, err = tx.Exec(InsertGrp) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertGrp "+err.Error()))
	}

	//장부유저 등록
	params["joinTy"] = "0"
	params["grpAuth"] = "0"
	params["authStat"] = "1"

	InsertGrpUser, err := cls.SetUpdateParam(darayoquery.InsertGrpUser, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertGrpUser "+err.Error()))
	}
	_, err = tx.Exec(InsertGrpUser) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertGrpUser "+err.Error()))
	}

	//회사 + 장부 연결
	//회사 장부 등록
	InsertCompanyBook, err := cls.SetUpdateParam(darayoquery.InsertCompanyBook, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyBook "+err.Error()))
	}
	_, err = tx.Exec(InsertCompanyBook) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyBook "+err.Error()))
	}

	//관리자 등록(회사 유저에 새로등록하는 장부아이디로 신규 유저 등록)
	params["useYn"] = "Y"
	InsertCompanyUser, err := cls.SetUpdateParam(darayoquery.InsertCompanyUser, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyUser "+err.Error()))
	}
	_, err = tx.Exec(InsertCompanyUser) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertCompanyUser "+err.Error()))
	}

	//biz 유저 추가
	//4. sys_user_info 시스템 유저 추가
	if sysFlag {
		InsertSysUserInfo, err := cls.SetUpdateParam(darayoquery.InsertSysUserInfo, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertSysUserInfo "+err.Error()))
		}
		_, err = tx.Exec(InsertSysUserInfo)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "InsertSysUserInfo "+err.Error()))
		}
	}

	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m["resultCode"] = "00"
	if !sysFlag {
		m["resultMsg"] = "장부장의 login 아이디로 sysUser 존재함"
	} else {
		m["resultMsg"] = "응답 성공"
	}

	return c.JSON(http.StatusOK, m)
}
