package payment

import (
	querysql "biz-web/query"
	commons "biz-web/query/commons"
	paymentsql "biz-web/query/payment"
	"biz-web/src/controller"
	Daily "biz-web/src/controller"
	"biz-web/src/controller/cls"
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf



func GetCombineStoreList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	resultCnt, err := cls.GetSelectData(paymentsql.SelectCombineStoreListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	//select
	resultList, err := cls.GetSelectTypeRequire(paymentsql.SelectCombineStoreList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultCnt"] = resultCnt[0]["TOTAL_COUNT"]
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}



func GetCombineSubStoreList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	//select
	resultList, err := cls.GetSelectTypeRequire(paymentsql.SelectCombineSubStoreList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}


func GetCombineStoreData(c echo.Context) error {



	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset


	restData, err := cls.GetSelectType(paymentsql.SelectCombineStoreData, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	paidList, err := cls.GetSelectTypeRequire(paymentsql.SelectCombinePaidList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	resultCnt, err := cls.GetSelectType(paymentsql.SelectCombineOrderListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}


	orderList, err := cls.GetSelectType(paymentsql.SelectCombineOrderList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	resultData := make(map[string]interface{})
	resultData["restData"] = restData[0]
	resultData["resultCnt"] = resultCnt[0]
	resultData["orderList"] = orderList
	resultData["paidList"] = paidList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData


	return c.JSON(http.StatusOK, m)

}


func GetCombineStoreDataWincube(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset




	resultCnt, err := cls.GetSelectType(paymentsql.SelectCombineOrderWincubeListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}


	orderList, err := cls.GetSelectType(paymentsql.SelectCombineOrderWincubeList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	resultData := make(map[string]interface{})
	resultData["resultCnt"] = resultCnt[0]
	resultData["orderList"] = orderList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData


	return c.JSON(http.StatusOK, m)

}


func GetPaidMngList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	resultCnt, err := cls.GetSelectData(paymentsql.SelectPaidMngListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	//select
	resultList, err := cls.GetSelectTypeRequire(paymentsql.SelectPaidMngList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultCnt"] = resultCnt[0]["TOTAL_COUNT"]
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}

func GetPaidList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	resultData, err := cls.GetSelectData(paymentsql.SelectPaidListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	//select
	resultList, err := cls.GetSelectData(paymentsql.SelectPaidList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = resultData
		m["resultList"] = []string{}

		return c.JSON(http.StatusOK, m)
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}

func GetPaidExcel(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	//select
	resultList, err := cls.GetSelectType(paymentsql.SelectPaidExcelList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}

func GetPaidIngList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	//println(tpayPwEncode("dalja1010&"))

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	//select
	selectQurey := ``

	switch params["select"] {
	case "SETTLMNT_DT":
		selectQurey = `AND A.SETTLMNT_DT >= '#{startDate}'
		AND A.SETTLMNT_DT <= '#{endDate}'
`
		break
	case "PAYMENT_DT":
		selectQurey = `
			AND A.PAYMENT_DT >= '#{startDate}' 
			AND A.PAYMENT_DT <= '#{endDate}'
`
		break
	case "SEND_DATE":
		selectQurey = ` 
			AND A.SEND_DATE >= '#{startDate}' 
			AND A.SEND_DATE <= '#{endDate}'
`
		break
	}
	selectQurey += `AND B.REST_NM LIKE '%#{searchText}%'
					AND A.RESULT_PAY_YN='N'
					ORDER BY A.SEND_DATE DESC,A.PAYMENT_DT ASC
`
	resultCnt, err := cls.GetSelectData(paymentsql.SelectPaidIngListCnt+selectQurey, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	resultList, err := cls.GetSelectData(paymentsql.SelectPaidIngList+selectQurey+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {
		m := make(map[string]interface{})
		m["resultCode"] = "01"
		m["resultMsg"] = "응답 성공"
		return c.JSON(http.StatusOK, m)
	}

	pgInfo, err := cls.SelectData(paymentsql.SelectPgInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	reqInfo := make(map[string]string)
	reqInfo["mid"] = pgInfo[0]["PG_MID"]
	reqInfo["uid"] = pgInfo[0]["PG_UID"]
	reqInfo["pw"] = pgInfo[0]["PG_PSWD"]
	reqInfo["fr_dt"] = params["startDate"]
	reqInfo["to_dt"] = params["endDate"]

	code, data := sendPayInfo(reqInfo)

	if code == "00" {
		for _, x := range data.List {
			//	println(x.Seq)
			for i := range resultList {
				r_seq, _ := strconv.Atoi(resultList[i]["RESULT_SEQ"])

				if x.Seq == r_seq {
					//println(r_seq)
					resultList[i]["status_nm"] = x.Status_nm
					resultList[i]["err_msg"] = x.Err_msg
				}
			}
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultCnt"] = resultCnt[0]["TOTAL_COUNT"]
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}

func GetPaidInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	//select
	resultData, err := cls.GetSelectType(paymentsql.SelectPaidInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData[0]

	return c.JSON(http.StatusOK, m)
}

func GetPaidOkList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	resultCnt, err := cls.GetSelectData(paymentsql.SelectPaidOkListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	//select
	resultList, err := cls.GetSelectData(paymentsql.SelectPaidOkList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	//	if resultList == nil {
	//	return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	//	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultCnt"] = resultCnt[0]["TOTAL_COUNT"]
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}

func GetPaidOkExcel(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	//select
	resultList, err := cls.GetSelectType(paymentsql.SelectPaidOkExcelList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	//	if resultList == nil {
	//	return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	//	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}

func SetPaidReq(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	pgInfo, err := cls.SelectData(paymentsql.SelectPgInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	sendList, err := cls.GetSelectData(paymentsql.SelectRestPaymentSendList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if sendList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	for i := range sendList {

		reqInfo := make(map[string]string)
		reqInfo["mid"] = pgInfo[0]["PG_MID"]
		reqInfo["uid"] = pgInfo[0]["PG_UID"]
		reqInfo["pw"] = pgInfo[0]["PG_PSWD"]



		reqInfo["sub_id"] = sendList[i]["REST_ID"]
		reqInfo["settlmnt_dt"] = sendList[i]["SETTLMNT_DT"]
		reqInfo["amt"] = sendList[i]["REST_PAYMENT_AMT"]
		reqInfo["moid"] = sendList[i]["REST_PAYMENT_ID"]

		resultMap := sendPayReq(reqInfo)
		resultMap["restPaymentId"] = sendList[i]["REST_PAYMENT_ID"]

		rQuery := ""
		if resultMap["result_cd"] == "0000" {
			rQuery = paymentsql.UpdateDarRestPaymentTpayRegResult
		} else {
			rQuery = paymentsql.UpdateDarRestPaymentTpayRegResultFail
		}

		tx, err := cls.DBc.Begin()
		if err != nil {
			//return "5100", errors.New("begin error")
		}

		// 오류 처리
		defer func() {
			if err != nil {
				// transaction rollback
				dprintf(4, c, "do rollback -지급 요청 결과 반영 SetPaidReq)  \n")
				tx.Rollback()
			}
		}()

		paidResultQuery, err := cls.GetQueryJson(rQuery, resultMap)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "paidResultQuery parameter fail"))
		}

		println(rQuery)
		_, err = tx.Exec(paidResultQuery)
		if err != nil {
			dprintf(1, c, "Query(%s) -> error (%s) \n", paidResultQuery, err)
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		err = tx.Commit()
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	//sendPayInfo(reqInfo)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)
}

func ManualTpay(c echo.Context) error {

	Daily.TpayMakePayStep1()
	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)
}

func SetMakeFee(c echo.Context) error {

	lprintf(4, "[Start] SetMakeFee \n")

	params := make(map[string]string)

	makeList, err := cls.SelectData(querysql.SelectFeeMakeList, params)
	if err != nil {
		lprintf(1, "[ERROR] SetMakeFee SelectFeeMakeList error(%s) \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "[ERROR] do rollback -Tpay 수수료 생성 error(%s) \n", err.Error())
			tx.Rollback()
		}
	}()

	for i := range makeList {

		mparam := make(map[string]string)
		mparam["moid"] = makeList[i]["HIST_ID"]
		mparam["histId"] = makeList[i]["MOID"]
		mparam["payMethod"] = makeList[i]["PAYMETHOD"]

		amt := makeList[i]["AMT"]
		mparam["amt"] = amt
		mparam["restId"] = makeList[i]["REST_ID"]
		paymethod := makeList[i]["PAYMENT_DT"]
		mparam["paymentDt"] = paymethod

		//가맹점 수수료
		restFeesInfo, err := cls.SelectData(paymentsql.SelectRestFeesInfo, mparam)
		if err != nil {
			lprintf(1, "[ERROR] SetMakeFee SelectRestFeesInfo error(%s) \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		//PG 수수료
		pgFeesInfo, err := cls.SelectData(paymentsql.SelectPgFeesInfo, mparam)
		if err != nil {
			lprintf(1, "[ERROR] SetMakeFee SelectPgFeesInfo error(%s) \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		//가맹점 수수료
		newRestFeesInfo := Daily.NewRestFees(restFeesInfo[0]["REST_FEES"], amt)

		//PG사 수수료
		newPgFeesInfo := Daily.NewPgFees(paymethod, pgFeesInfo[0]["FEE"], amt, pgFeesInfo[0]["VAT_YN"])

		// 회사 수수료
		newFitFeesInfo := Daily.NewFitFees(newRestFeesInfo["RestFeesAmt"], newPgFeesInfo["PgFeesAmt"])

		// 지급 금액
		restPayAmt := Daily.RealAmt(amt, newRestFeesInfo["RestFeesAmt"])

		mparam["restPayAmt"] = restPayAmt

		mparam["totSuplyAmt"] = newRestFeesInfo["RestSupplyAmt"]
		mparam["totVat"] = newRestFeesInfo["RestVatAmt"]
		mparam["totFee"] = newRestFeesInfo["RestFeesAmt"]

		mparam["pgSuplyAmt"] = newPgFeesInfo["PgSupplyAmt"]
		mparam["pgVat"] = newPgFeesInfo["PgVatAmt"]
		mparam["pgFee"] = newPgFeesInfo["PgFeesAmt"]

		mparam["fitSuplyAmt"] = newFitFeesInfo["FitSupplyAmt"]
		mparam["fitVat"] = newFitFeesInfo["FitVatAmt"]
		mparam["fitFee"] = newFitFeesInfo["FitFeesAmt"]

		UpdatePaymentFeesQuery, err := cls.GetQueryJson(querysql.UpdatePaymentFees, mparam)
		if err != nil {
			lprintf(1, "[ERROR] SetMakeFee UpdatePaymentFeesQuery error(%s) \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = tx.Exec(UpdatePaymentFeesQuery)
		if err != nil {
			lprintf(1, "[ERROR] SetMakeFee UpdatePaymentFeesQuery error(%s) \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		UpdatePaymentHistFeesQuery, err := cls.GetQueryJson(querysql.UpdatePaymentHistFees, mparam)
		if err != nil {
			lprintf(1, "[ERROR] SetMakeFee UpdatePaymentHistFeesQuery error(%s) \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = tx.Exec(UpdatePaymentHistFeesQuery)
		if err != nil {
			lprintf(1, "[ERROR] SetMakeFee UpdatePaymentHistFeesQuery error(%s) \n", err.Error())
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		lprintf(1, "[ERROR] SetMakeFee UpdateStoreInfoQuery error(%s) \n", err.Error())
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)
}

func SetMakePaid(c echo.Context) error {

	Daily.TpayMakePayStep2()

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)
}

func SetPaidInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -지급 상세 수정  SetPaidInfo)  \n")
			tx.Rollback()
		}
	}()

	ResultPayYnQuery, err := cls.GetQueryJson(paymentsql.UpdateDarRestPayment, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateDarRestPayment parameter fail"))
	}

	_, err = tx.Exec(ResultPayYnQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", ResultPayYnQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//sendPayInfo(reqInfo)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)
}

func SetSettlmntDtChange(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -지급요청일 일괄 변경 SetSettlmntDtChange)  \n")
			tx.Rollback()
		}
	}()

	ResultPayYnQuery, err := cls.GetQueryJson(paymentsql.UpdateDarRestPaymentSettlmentDt, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "ResultPayYnQuery parameter fail"))
	}

	_, err = tx.Exec(ResultPayYnQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", ResultPayYnQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//sendPayInfo(reqInfo)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)
}


func GetAccount(c echo.Context) error {

	params := cls.GetParamJsonMap(c)


	pgInfo, err := cls.SelectData(paymentsql.SelectPgInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	mid := pgInfo[0]["PG_MID"]
	uid := pgInfo[0]["PG_UID"]
	pw := pgInfo[0]["PG_PSWD"]
	fr_dt := params["startDate"]
	to_dt := params["endDate"]
	row_cnt := "10000"
	page := "1"



	TPAY_ACCOUNT_URL := "https://mms.tpay.co.kr/api/v1/om/account"

	payload := strings.NewReader("mid=" + mid +
		"&uid=" + uid +
		"&pw=" + tpayPwEncode(pw) +
		"&fr_dt=" + fr_dt +
		"&to_dt=" + to_dt +
		"&row_cnt=" + row_cnt +
		"&page=" + page)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", TPAY_ACCOUNT_URL, payload)
	if err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	//println(string(body))
	var response TpayAccountData
	err = json.Unmarshal(body, &response)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	resultMap := make(map[string]interface{})

	resultMap["result_cd"] = response.Result_cd
	resultMap["result_msg"] = response.Result_msg
	resultMap["result_count"] = response.Tot_cnt
	resultMap["list"] = response.List



	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)

}

func SetPaidOk(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	resultList, err := cls.GetSelectData(paymentsql.SelectPaidIngList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {
		m := make(map[string]interface{})
		m["resultCode"] = "01"
		m["resultMsg"] = "응답 성공"
		return c.JSON(http.StatusOK, m)
	}

	pgInfo, err := cls.SelectData(paymentsql.SelectPgInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	reqInfo := make(map[string]string)
	reqInfo["mid"] = pgInfo[0]["PG_MID"]
	reqInfo["uid"] = pgInfo[0]["PG_UID"]
	reqInfo["pw"] = pgInfo[0]["PG_PSWD"]
	reqInfo["fr_dt"] = params["startDate"]
	reqInfo["to_dt"] = params["endDate"]

	code, data := sendPayInfo(reqInfo)

	if code == "00" {
		for _, x := range data.List {

			for i := range resultList {
				r_seq, _ := strconv.Atoi(resultList[i]["RESULT_SEQ"])

				if x.Seq == r_seq && x.Status_nm == "성공" {

					//println(resultList[i]["restPaymentId"])

					payResult := make(map[string]string)
					payResult["restPaymentId"] = resultList[i]["restPaymentId"]

					tx, err := cls.DBc.Begin()
					if err != nil {
						//return "5100", errors.New("begin error")
					}

					// 오류 처리
					defer func() {
						if err != nil {
							// transaction rollback
							dprintf(4, c, "do rollback -지급 완료 결과 반영 SetPaidOk)  \n")
							tx.Rollback()
						}
					}()

					ResultPayYnQuery, err := cls.GetQueryJson(paymentsql.UpdateDarRestPaymentTpayResultPay, payResult)
					if err != nil {
						return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "ResultPayYnQuery parameter fail"))
					}
					println(ResultPayYnQuery)
					_, err = tx.Exec(ResultPayYnQuery)
					if err != nil {
						dprintf(1, c, "Query(%s) -> error (%s) \n", ResultPayYnQuery, err)
						return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
					}
					err = tx.Commit()
					if err != nil {
						return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
					}

				}
			}
		}
	}

	//sendPayInfo(reqInfo)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)
}

type PaidResponseData struct {
	Result_cd   string `json:"result_cd"`   //
	Result_msg  string `json:"result_msg"`  //
	Seq         int    `json:"seq"`         //
	Settlmnt_dt string `json:"settlmnt_dt"` //
	Sub_id      string `json:"sub_id"`      //
	Moid        string `json:"moid "`       //

}


type TpayAccountData struct {
	Result_cd   string `json:"result_cd"`   //
	Result_msg  string `json:"result_msg"`  //
	Tot_cnt     int    `json:"tot_cnt"`         //
	Page     	int    `json:"page"`         //
	Row_cnt     int    `json:"row_cnt"`         //
	List	 	[]TpayAccountDataList `json:"list"` //


}

type TpayAccountDataList struct {
	Id   			string `json:"id"`   //
	Tr_dt   		string `json:"tr_dt"`   //
	Tr_cl_nm   		string `json:"tr_cl_nm"`   //
	In_amt   		int `json:"in_amt"`   //
	Out_amt   		int `json:"out_amt"`   //
	Remain_amt   	int `json:"remain_amt"`   //
}


type TpaySubmallRegRecvData struct {
	Result_cd  string `json:"result_cd"`  //
	Result_msg string `json:"result_msg"` //
	mid        string `json:"mid"`        //
	sub_id     string `json:"sub_id"`     //
}

type TpayRecvData struct {
	Result_cd  string `json:"result_cd"`  //
	Result_msg string `json:"result_msg"` //
	PayMethod  string `json:"PayMethod"`  //
	CancelDate string `json:"CancelDate"` //
	CancelTime string `json:"CancelTime"` //
	Tid        string `json:"tid"`        //
	Moid       string `json:"moid"`       //
}

type TpayPaidInfoData struct {
	Result_cd  string         `json:"result_cd"`  //
	Result_msg string         `json:"result_msg"` //
	Page       int            `json:"page"`       //
	Row_cnt    int            `json:"row_cnt"`    //
	Tot_cnt    int            `json:"tot_cnt"`    //
	Req_cnt    int            `json:"req_cnt"`    //
	Succ_cnt   int            `json:"succ_cnt"`   //
	Fail_cnt   int            `json:"fail_cnt"`   //
	List       []PaidInfoList `json:"list"`       //
}

type PaidInfoList struct {
	Seq          int    `json:"seq"`          //
	Id           string `json:"id"`           //
	Settlmnt_dt  string `json:"settlmnt_dt"`  //
	Sub_co_no    string `json:"sub_co_no"`    //
	Sub_nm       string `json:"sub_nm"`       //
	Sub_id       string `json:"sub_id"`       //
	Settlmnt_amt int    `json:"settlmnt_amt"` //
	Accnt_nm     string `json:"accnt_nm"`     //
	Bank_nm      string `json:"bank_nm"`      //
	Accnt_no     string `json:"accnt_no"`     //
	Status_nm    string `json:"status_nm"`    //
	Err_msg      string `json:"err_msg"`      //
}

func sendPayReq(reqParam map[string]string) map[string]string {

	mid := reqParam["mid"]
	uid := reqParam["uid"]
	pw := reqParam["pw"]
	sub_id := reqParam["sub_id"]
	settlmnt_dt := reqParam["settlmnt_dt"]
	amt := reqParam["amt"]
	moid := reqParam["moid"]

	resultMap := make(map[string]string)

	TPAY_SETTLMNT_URL := "https://mms.tpay.co.kr/api/v1/om/settlmnt/register"

	payload := strings.NewReader("mid=" + mid +
		"&uid=" + uid +
		"&pw=" + tpayPwEncode(pw) +
		"&sub_id=" + sub_id +
		"&settlmnt_dt=" + settlmnt_dt +
		"&amt=" + amt +
		"&moid=" + moid)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", TPAY_SETTLMNT_URL, payload)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] error sendPayReq: %s\n", err)
		resultMap["result_cd"] = "9999"
		resultMap["result_msg"] = "전송 요청 오류"
		return resultMap
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		resultMap["result_cd"] = "9999"
		resultMap["result_msg"] = "전송 요청 오류"
		return resultMap
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		resultMap["result_cd"] = "9999"
		resultMap["result_msg"] = "전송 요청 오류"
		return resultMap
	}

	println(string(body))

	var response PaidResponseData
	err = json.Unmarshal(body, &response)
	if err != nil {
		lprintf(1, "[ERROR] error sendPayReq: %s\n", err.Error())
		resultMap["result_cd"] = "9999"
		resultMap["result_msg"] = "전송 요청 오류"
		return resultMap
	}

	resultMap["result_cd"] = response.Result_cd
	resultMap["result_msg"] = response.Result_msg
	resultMap["seq"] = strconv.Itoa(response.Seq)
	resultMap["sub_id"] = response.Sub_id
	resultMap["settlmnt_dt"] = response.Settlmnt_dt
	resultMap["moid"] = response.Moid

	return resultMap

}

func sendPayInfo(reqParam map[string]string) (string, TpayPaidInfoData) {

	//TpayMap := make(map[string]string)

	mid := reqParam["mid"]
	uid := reqParam["uid"]
	pw := reqParam["pw"]

	TPAY_SETTLMNT_URL := "https://mms.tpay.co.kr/api/v1/om/settlmnt"

	row_cnt := "10000"
	page := "1"
	fr_dt := reqParam["fr_dt"]
	to_dt := reqParam["to_dt"]

	payload := strings.NewReader("mid=" + mid +
		"&uid=" + uid +
		"&pw=" + tpayPwEncode(pw) +
		"&row_cnt=" + row_cnt +
		"&page=" + page +
		"&fr_dt=" + fr_dt +
		"&to_dt=" + to_dt)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", TPAY_SETTLMNT_URL, payload)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] error sms send : %s\n", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] error sms send : %s\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] error sms send : %s\n", err)
	}

	//println(string(body))

	var response TpayPaidInfoData
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "99", response
	}

	return "00", response

}

func tpayPwEncode(str string) string {

	hasher := md5.New()
	hasher.Write([]byte(str))
	md5 := string(base64.StdEncoding.EncodeToString(hasher.Sum(nil)))
	h := sha512.New()
	h.Write([]byte(md5))
	sha512 := string(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	return url.QueryEscape(sha512)

}

func GetPaidIngExcel(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	//select
	resultList, err := cls.GetSelectType(paymentsql.SelectPaidIngExcelList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	//	if resultList == nil {
	//	return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	//	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)
}



func SetStoreReg(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	pgInfo, err := cls.SelectData(paymentsql.SelectPgInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	mid := pgInfo[0]["PG_MID"]
	uid := pgInfo[0]["PG_UID"]
	pw := pgInfo[0]["PG_PSWD"]

	storePayInfo, err := cls.SelectData(paymentsql.SelectStorePayInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if storePayInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "가맹점 정보가 부족합니다."))
	}


	if storePayInfo[0]["ACCOUNT_CERT_YN"]=="N"{
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "계좌정보가 없습니다."))
	}

	sub_id := params["restId"]
	co_cl:="1"
	sub_co_no:=storePayInfo[0]["BUSID"]
	sub_co_nm:=storePayInfo[0]["REST_NM"]
	accnt_no:=storePayInfo[0]["ACCOUNT_NO"]
	bank_cd:=storePayInfo[0]["BANK_CD"]
	accnt_nm:=storePayInfo[0]["ACCOUNT_NM"]

	resultMap := make(map[string]string)

	TPAY_SETTLMNT_URL := "https://mms.tpay.co.kr/api/v1/om/submall/register"

	payload := strings.NewReader("mid=" + mid +
		"&uid=" + uid +
		"&pw=" + tpayPwEncode(pw) +
		"&sub_id=" + sub_id +
		"&co_cl=" + co_cl +
		"&sub_co_no=" + sub_co_no +
		"&sub_co_nm=" + sub_co_nm +
		"&accnt_no=" + accnt_no +
		"&bank_cd=" + bank_cd +
		"&accnt_nm=" + accnt_nm)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", TPAY_SETTLMNT_URL, payload)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] error sendSubmallReq: %s\n", err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//println(string(body))

	var response TpaySubmallRegRecvData
	err = json.Unmarshal(body, &response)
	if err != nil {
		lprintf(1, "[ERROR] error sendSubmallReq: %s\n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultMap["resultCd"] = response.Result_cd
	resultMap["result_msg"] = response.Result_msg
	resultMap["restId"] = response.sub_id
	resultMap["mid"] = response.mid

	lprintf(4, "[INFO]  sendSubmallReq Result_cd: %s\n", response.Result_cd)
	lprintf(4, "[INFO]  sendSubmallReq result_msg: %s\n", response.Result_msg)

	if response.Result_cd=="0000"{

		tx, err := cls.DBc.Begin()
		if err != nil {
			//return "5100", errors.New("begin error")
		}

		// 오류 처리
		defer func() {
			if err != nil {
				// transaction rollback
				dprintf(4, c, "do rollback -tpay 서브몰 등록 SetStoreReg)  \n")
				tx.Rollback()
			}
		}()

		submallRegQuery, err := cls.GetQueryJson(paymentsql.UpdateTpayStoreRegResult, resultMap)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "submallRegQuery parameter fail"))
		}

		_, err = tx.Exec(submallRegQuery)
		if err != nil {
			dprintf(1, c, "Query(%s) -> error (%s) \n", submallRegQuery, err)
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		err = tx.Commit()
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

	}


	//sendPayInfo(reqInfo)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)
}

func SetCombinePaidMake(c echo.Context) error {


	lprintf(4, "[Start] SetCombinePaidMake \n")

	params := cls.GetParamJsonMap(c)
	params["dayBefore"] = "1"

	paymentDt, err := cls.GetSelectData(querysql.SelectDayBefore, params, c )
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	params["paymentDt"] = paymentDt[0]["beforeDate"]


	settlmtDt, err := cls.GetSelectData(querysql.SelectBizDay, params, c )
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if settlmtDt == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	params["settlmtDt"] = settlmtDt[0]["totalDate"]




	//  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	txErr := err
	defer func() {
		if txErr != nil {
			// transaction rollback
			lprintf(4, "[ERROR] do rollback -SetCombinePaidMake error(%s) \n", err.Error())
			tx.Rollback()
		}
	}()

		mparam := make(map[string]string)

		paymethod := "CARD"
		amt := params["amt"]
		chargeAmt := params["chargeAmt"]
		mparam["amt"]=amt
		mparam["restId"]=params["restId"]
		mparam["restNm"]=params["restNm"]
		mparam["payMethod"]=paymethod
		mparam["settlmtDt"]=params["settlmtDt"]




		//가맹점 수수료
		restFeesInfo, err := cls.GetSelectData(paymentsql.SelectRestFeesInfo, mparam, c)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		//PG 수수료
		pgFeesInfo, err := cls.GetSelectData(paymentsql.SelectPgFeesInfo, mparam, c)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}


		//가맹점 수수료
		newRestFeesInfo := Daily.NewRestFees(restFeesInfo[0]["REST_FEES"],amt)
		//PG사 수수료
		newPgFeesInfo := Daily.NewPgFees(paymethod,pgFeesInfo[0]["FEE"],amt,pgFeesInfo[0]["VAT_YN"])
		// 회사 수수료
		newFitFeesInfo := Daily.NewFitFees(newRestFeesInfo["RestFeesAmt"],newPgFeesInfo["PgFeesAmt"])
		// 지급 금액
		restPayAmt := Daily.RealAmt(amt,newRestFeesInfo["RestFeesAmt"])

		mparam["restPayAmt"]=restPayAmt
		mparam["totFee"]=newRestFeesInfo["RestFeesAmt"]
		mparam["pgFee"]=newPgFeesInfo["PgFeesAmt"]
		mparam["fitFee"]=newFitFeesInfo["FitFeesAmt"]




		InserCombineDarRestPaymentQuery, err := cls.GetQueryJson(querysql.InserCombineDarRestPayment, mparam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = tx.Exec(InserCombineDarRestPaymentQuery)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}



		intAmt, _ := strconv.Atoi(amt)
		intchargeAmt, _ := strconv.Atoi(chargeAmt)
		mparam["chargeAmt"] = strconv.Itoa(intchargeAmt - intAmt)

		UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombinePaidChargeAmt, mparam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
		}
		_, err = tx.Exec(UpdateCombineChargeAmtQuery)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}



	// transaction commit
	err = tx.Commit()
	if err != nil {
		lprintf(1, "[ERROR] SetCombinePaidMake error(%s) \n", err.Error())
	}


	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)

}




func GetPaymentList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	resultCnt, err := cls.GetSelectData(paymentsql.SelectPaymentListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultList, err := cls.GetSelectType(paymentsql.SelectPaymentList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = resultCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]

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



func SetGiftconUpdate(c echo.Context) error {


	lprintf(4, "[Start] SetGiftconUpdate \n")

	params := cls.GetParamJsonMap(c)
	store := params["store"]
	restId:= params["restId"]

	P_URL := ""
	resultMap := make(map[string]string)

	fname := cls.Cls_conf(os.Args)
	mocalUrl, _ := cls.GetTokenValue("MOCA_URL", fname)

	if store=="b"{
		P_URL=mocalUrl+"/api/moca/v2/benepicons/check?restId=" + restId
	}else{
		P_URL=mocalUrl+"/api/moca/v2/wincubes/status?restId=" + restId
	}

	payload := strings.NewReader("")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", P_URL, payload)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] error call SetGifticonCheck: %s\n", err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	res, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	println(string(body))




	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultMap

	return c.JSON(http.StatusOK, m)

}


