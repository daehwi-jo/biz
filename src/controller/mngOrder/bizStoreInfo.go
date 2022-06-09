package mngOrder

import (
	"biz-web/query/commons"
	bizOrderMngQuery "biz-web/query/mngOrder"
	paymentSql "biz-web/query/payment"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetStoreInfoList(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	storeDataCnt, err := cls.GetSelectData(bizOrderMngQuery.SelectStoreMngCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	storeData, err := cls.GetSelectType(bizOrderMngQuery.SelectStoreMng+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultCnt"] = storeDataCnt[0]["TOTAL_COUNT"]
	m["resultOrderCnt"] = storeDataCnt[0]["TOTAL_ORDER_COUNT"]
	m["resultList"] = storeData

	return c.JSON(http.StatusOK, m)
}

func GetStoreChargeInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	linkChk, err := cls.GetSelectDataRequire(paymentSql.SelectAgrmInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if linkChk == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "연결 정보가 없습니다."))
	}
	reqStat := linkChk[0]["REQ_STAT"]
	if reqStat != "1" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "연결 정보를 확인해주세요."))
	}
	prepaidAmt, _ := strconv.Atoi(linkChk[0]["PREPAID_AMT"])

	paymentUseYn := linkChk[0]["PAYMENT_USE_YN"]

	storeInfo, err := cls.GetSelectData(bizOrderMngQuery.SelectStoreInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	store := make(map[string]interface{})
	store["storeInfo"] = storeInfo[0]
	store["prepaidAmt"] = prepaidAmt

	if paymentUseYn == "N" {
		store["storeChargeInfo"] = []string{}
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = store

		return c.JSON(http.StatusOK, m)
	}

	amtList, err := cls.GetSelectType(bizOrderMngQuery.SelectChargeAmtList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if amtList == nil {
		store["storeChargeInfo"] = []string{}
	}
	store["storeChargeInfo"] = amtList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = store

	return c.JSON(http.StatusOK, m)
}

func GetStoreUnpaidInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	storeInfo, err := cls.GetSelectData(bizOrderMngQuery.SelectStoreInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData, err := cls.GetSelectDataRequire(bizOrderMngQuery.SelectUnpaidListCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	totalCount, _ := strconv.Atoi(resultData[0]["orderCnt"])
	totalAmt, _ := strconv.Atoi(resultData[0]["TOTAL_AMT"])

	unPaidList, err := cls.GetSelectType(bizOrderMngQuery.SelectUnpaidList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	store := make(map[string]interface{})
	store["totalCount"] = totalCount
	store["totalAmt"] = totalAmt
	store["unPaidList"] = unPaidList
	store["storeInfo"] = storeInfo[0]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = store

	return c.JSON(http.StatusOK, m)
}
