package mngAdmin

import (
	"biz-web/query/commons"
	"biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetBizBookList (Java = BizUserMngDesc)
func GetBizBookList(c echo.Context) error {

	cls.Lprintf(4, "[INFO] GetBizBookList call \n")

	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	addQuery := ""
	if params["searchKey"] == "userNm" {
		addQuery = "AND C.USER_NM LIKE '%#{searchKeyword}%'"
	} else {
		addQuery = "AND A.GRP_NM LIKE '%#{searchKeyword}%'"
	}

	bookListCnt, err := cls.GetSelectDataRequire(mngAdmin.SelectBookListCnt+addQuery, params, c)
	if err != nil {
		cls.Lprintf(1, "[ERROR] query  err(%s) \n", err.Error())
		cls.Lprintf(1, "[ERROR] query  err(%d) \n", 3)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	bookList, err := cls.GetSelectTypeRequire(mngAdmin.SelectBookList+addQuery+commons.RegDateOrderBy+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = bookListCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(bookListCnt[0]["TOTAL_COUNT"])

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

func GetBookStoreList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectTypeRequire(mngAdmin.SelectAddStoreList, params, c)

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

// GetBizBookCsInfo
// notice : (Java = BizUserMngDesc)
func GetBizBookCsInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	bookInfo, err := cls.GetSelectData(mngAdmin.SelectBookInfo, params, c) //장부 정보
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if bookInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	linkStore, err := cls.GetSelectTypeRequire(mngAdmin.SelectBookLinkStore, params, c) //장부에 연결된 가맹점
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	bookInfoData := bookInfo[0]
	query := ""

	if bookInfoData["grpTypeCd"] == "1" { //회사장부
		if bookInfoData["supportYn"] == "Y" { //지원한도 있을때
			query = mngAdmin.SelectCompanyGrpUserList
		} else {
			query = mngAdmin.SelectPrivateGrpUserList
		}
	} else { //개인장부
		query = mngAdmin.SelectPrivateGrpUserList
	}

	bookUser, err := cls.GetSelectTypeRequire(query, params, c) //장부를 사용하는 사용자
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["userInfo"] = bookInfoData
	result["linkStore"] = linkStore
	result["bookUser"] = bookUser

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetBookData(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	bookData, err := cls.GetSelectType(mngAdmin.SelectBookData, params, c) //장부 정보
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if bookData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	result := make(map[string]interface{})
	result["bookData"] = bookData

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ModBookData(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	modBookQuery, err := cls.SetUpdateParam(mngAdmin.UpdateBookData, params) //장부 정보
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(modBookQuery) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}
