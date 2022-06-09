package mngOrder

import (
	"biz-web/query/commons"
	"biz-web/query/mngOrder"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetBookList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	resultData := make(map[string]interface{})

	if params["companyId"] == "" {
		resultData["bookList"] = []string{}
	} else {
		bookList, err := cls.GetSelectData(mngOrder.SelectGrpCode, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		if bookList != nil {
			resultData["bookList"] = bookList
		} else {
			resultData["bookList"] = []string{}
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

func GetOrderHistory(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	if params["companyId"] == "" {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "params error"))
	}

	orderListCnt, err := cls.GetSelectData(mngOrder.SelectOrderListCount+mngOrder.SelectOrderListCountWhere(params), params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	orderList, err := cls.GetSelectType(mngOrder.SelectOrderList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData := make(map[string]interface{})
	resultData["totalCount"] = orderListCnt[0]["TOTAL_COUNT"]
	resultData["totalAmt"] = orderListCnt[0]["ALL_AMT"]
	resultData["totalPage"] = utils.GetTotalPage(orderListCnt[0]["TOTAL_COUNT"])

	if orderList == nil {
		resultData["orderList"] = []string{}
	} else {
		resultData["orderList"] = orderList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

func GetOrderHistoryExcel(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	listData, err := cls.GetSelectType(mngOrder.SelectOrderExcelList, params, c)

	resultData := make(map[string]interface{})
	if listData == nil {
		resultData["listData"] = make(map[string]interface{})
	} else {
		resultData["listData"] = listData
	}

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}
