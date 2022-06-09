package mngOrder

import (
	"biz-web/query/commons"
	"biz-web/query/mngOrder"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetCalculateBookList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	pageSize,_ := strconv.Atoi(params["pageSize"])
	pageNo,_ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo-1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	resultCnt, err := cls.GetSelectDataRequire(mngOrder.SelectGrpCalculateCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultList, err := cls.GetSelectTypeRequire(mngOrder.SelectGrpCalculateList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = resultCnt[0]["TOTAL_COUNT"]


	if resultList == nil {
		result["resultList"] = []string{}
	} else {
		result["resultList"] = resultList
	}



	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}
