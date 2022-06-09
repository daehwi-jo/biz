package controller

import (
	"biz-web/src/controller/cls"
)

//결과값 관련 구조체

var lprintf func(int, string, ...interface{}) = cls.Lprintf

func SetErrResult(code, msg string) ResponseHeader {
	var resCode ResponseHeader
	resCode.ResultCode = code
	resCode.ResultMsg = msg
	cls.Fprint(code, msg)
	return resCode
}

func SetResult(code string, data []byte) Response {
	var resCode Response
	resCode.ResultCode = code
	resCode.ResultMsg = "성공"
	resCode.ResultData = data
	return resCode
}

/*
func getPage(rowNum, sPage, lPage, cPage int) Page {
	var pageInfo Page
	pageInfo.TotalBlock = 1
	pageInfo.StartPage = sPage
	pageInfo.LastPage = lPage
	pageInfo.CurrentBlock = 1
	pageInfo.TotalPage = 1
	pageInfo.NextPage = 0
	pageInfo.PrevPage = 0
	pageInfo.TotalCount = rowNum
	pageInfo.CurrentPage = cPage

	var pageInfoList PageList
	pageInfoList.PageNoText = cPage
	pageInfoList.PageNo = cPage
	pageInfoList.ClassName = "on"

	pageInfo.PagingList = append(pageInfo.PagingList, pageInfoList)

	return pageInfo
}

*/

// commons
type ResponseHeader struct {
	ResultCode string `json:"resultCode"` // result code
	ResultMsg  string `json:"resultMsg"`  // result msg
}

// commons
type Response struct {
	ResultCode string      `json:"resultCode"` // result
	ResultMsg  string      `json:"resultMsg"`  // result code
	ResultData interface{} `json:"resultData"` // result data
}

type PageList struct {
	PageNoText int    `json:"pageNoText"`
	PageNo     int    `json:"pageNo"`
	ClassName  string `json:"className"`
}

type UseHistory struct {
	UserId    string `json:"UserId"`    // mngUser id
	UserNm    string `json:"UserNm"`    // mngUser name
	RestNm    string `json:"RestNm"`    // 가맹점 이름
	ItemNm    string `json:"ItemNm"`    // 구매 아이템 이름
	OrderDate string `json:"OrderDate"` // 주문날자시간
	UseAmt    string `json:"UseAmt"`    // 주문 금액
}




var TPAY_MID string
var TPAY_MERCHANT_KEY string
var TPAY_PAY_ACTION_URL string
var TPAY_PAY_ACTION_WEB_URL string
var TPAY_PAY_LOCAL_URL string
var TPAY_RETURN_URL string
var TPAY_CANCEL_URL string

var TPAY_MID_SIMPLE_PAY string
var TPAY_SIMPLE_PAY_MERCHANT_KEY string
var TPAY_SIMPLE_PAY_GEN_URL string
var TPAY_SIMPLE_PAY_DEL_URL string
var TPAY_SIMPLE_PAY_PAYMENT_URL string
var TPAY_SIMPLE_PAY_CANCEL_URL string
var TPAY_SIMPLE_PAY_CANCEL_PWD string
var CONFIG_MOCA_URL string



// Tpay설정
func Tpay_conf(fname string) int {
	lprintf(4, "[INFO] Tpay_conf start (%s)\n", fname)

	TPAY_MID, _ = cls.GetTokenValue("TPAY.TPAY_MID", fname)
	TPAY_MERCHANT_KEY, _ = cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)
	TPAY_PAY_ACTION_URL, _ = cls.GetTokenValue("TPAY.PAY_ACTION_URL", fname)
	TPAY_PAY_ACTION_WEB_URL, _ = cls.GetTokenValue("TPAY.PAY_ACTION_WEB_URL", fname)
	TPAY_PAY_LOCAL_URL, _ = cls.GetTokenValue("TPAY.PAY_LOCAL_URL", fname)
	TPAY_RETURN_URL, _ = cls.GetTokenValue("TPAY.RETURN_URL", fname)
	TPAY_CANCEL_URL, _ = cls.GetTokenValue("TPAY.CANCEL_URL", fname)


	TPAY_MID_SIMPLE_PAY, _ = cls.GetTokenValue("TPAY.MID_SIMPLE_PAY", fname)
	TPAY_SIMPLE_PAY_MERCHANT_KEY, _ = cls.GetTokenValue("TPAY.SIMPLE_PAY_MERCHANT_KEY", fname)
	TPAY_SIMPLE_PAY_GEN_URL, _ = cls.GetTokenValue("TPAY.SIMPLE_PAY_GEN_URL", fname)
	TPAY_SIMPLE_PAY_DEL_URL, _ = cls.GetTokenValue("TPAY.SIMPLE_PAY_DEL_URL", fname)
	TPAY_SIMPLE_PAY_PAYMENT_URL, _ = cls.GetTokenValue("TPAY.SIMPLE_PAY_PAYMENT_URL", fname)
	TPAY_SIMPLE_PAY_CANCEL_URL, _ = cls.GetTokenValue("TPAY.SIMPLE_PAY_CANCEL_URL", fname)
	TPAY_SIMPLE_PAY_CANCEL_PWD, _ = cls.GetTokenValue("TPAY.SIMPLE_PAY_CANCEL_PWD", fname)
	CONFIG_MOCA_URL, _ = cls.GetTokenValue("MOCA_URL", fname)

	return 0
}

