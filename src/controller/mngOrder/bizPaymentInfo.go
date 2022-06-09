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

func GetPaymentHistory(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	paymentListCnt, err := cls.GetSelectDataRequire(mngOrder.SelectPaymentListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	paymentList, err := cls.GetSelectTypeRequire(mngOrder.SelectPaymentList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData := make(map[string]interface{}) //두개 이상의 결과값을 리턴 할 경우
	resultData["totalCount"] = paymentListCnt[0]["TOTAL_COUNT"]
	resultData["totalAmt"] = paymentListCnt[0]["ALL_AMT"]
	resultData["totalPage"] = utils.GetTotalPage(paymentListCnt[0]["TOTAL_COUNT"])

	if paymentList == nil {
		resultData["paymentList"] = []string{}
	} else {
		resultData["paymentList"] = paymentList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}
