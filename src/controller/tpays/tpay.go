package tpays

import (
	"biz-web/query/mngOrder"
	paymentsql "biz-web/query/payment"
	"biz-web/src/controller"
	apiPush "biz-web/src/controller/api/push"
	"biz-web/src/controller/cls"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	//apiPush "biz-web/src/controller/api/push"
)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

type Tpay struct {
	MerchantKey string
	EncKey      string
	EdiDate     string
	Mid         string
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

type TpayGenBillingData struct {
	Result_cd  string `json:"result_cd"`  //
	Result_msg string `json:"result_msg"` //
	Card_token string `json:"card_token"` //
	Card_num   string `json:"card_num"`   //
	Card_code  string `json:"card_code"`  //
	Card_name  string `json:"card_name"`  //

	Tid      string `json:"tid"`      //
	Paid_amt string `json:"paid_amt"` //
	Moid     string `json:"moid"`     //

	App_co      string `json:"app_co"`      //
	App_co_name string `json:"app_co_name"` //
	App_no      string `json:"app_no"`      //

	CancelDate string `json:"CancelDate"` //
	CancelTime string `json:"CancelTime"` //

}

type MocaOrderInfo struct {
	OrderNo    string `json:"orderNo"`    //
	CpNo       string `json:"cpNo"`       //
	OrdNo      string `json:"ordNo"`      //
	ItemPrice  string `json:"itemPrice"`  //
	ItemName   string `json:"itemName"`   //
	ExpireDate string `json:"expireDate"` //
	ItemImg    string `json:"itemImg"`    //
	ItemDesc   string `json:"itemDesc"`   //
}

type MocaOrderResult struct {
	Result_cd  string        `json:"resultCode"` //
	Result_msg string        `json:"resultMsg"`  //
	ResultData MocaOrderInfo `json:"resultData"` //
}

// 결제 준비
func TpayReady(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	mid, _ := cls.GetTokenValue("TPAY.TPAY_MID", fname)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)
	payActionUrl, _ := cls.GetTokenValue("TPAY.PAY_ACTION_URL", fname)
	payActionWEBUrl, _ := cls.GetTokenValue("TPAY.PAY_ACTION_WEB_URL", fname)
	payLocalUrl, _ := cls.GetTokenValue("TPAY.PAY_LOCAL_URL", fname)
	returnUrl, _ := cls.GetTokenValue("TPAY.RETURN_URL", fname)
	cancelUrl, _ := cls.GetTokenValue("TPAY.CANCEL_URL", fname)

	paymentTy := params["paymentTy"]

	if params["searchTy"] == "1" {
		paymentTy = "0"
		params["paymentTy"] = "0"
	} else {
		paymentTy = "3"
		params["paymentTy"] = "3"
	}

	selectedDate := ""
	appPrefix := "mngAdmin://"
	selectedDate = params["selectedDate"]
	userAgent := c.Request().Header.Get("User-Agent")

	//if paymentTy == "2" || paymentTy == "3" {
	//	selectedDate = selectedDate + "235959"
	//}
	connType := "1"
	ediDate, vbankExpDate := getEdiDate()
	moid := getMoid(params["mallUserId"])
	pgCd := "01"
	userTy := "0"

	var tpay Tpay
	tpay.MerchantKey = merchantKey
	tpay.EdiDate = ediDate

	key := fmt.Sprintf("%s%s", ediDate, merchantKey)
	tpay.EncKey = getMD5HashHandler(key)

	input := fmt.Sprintf("%s%s%s", params["amt"], mid, moid)
	encryptData := base64.StdEncoding.EncodeToString(AesHandlerEncrypt(input, tpay.EncKey, merchantKey))

	userInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	goodsName := ""
	if params["searchTy"] == "1" {
		goodsName = userInfo[0]["REST_NM"] + " 충전"
	} else {
		goodsName = userInfo[0]["REST_NM"] + " 정산"
	}

	params["buyerEmail"] = userInfo[0]["EMAIL"]
	params["buyerTel"] = userInfo[0]["HP_NO"]
	params["userNm"] = userInfo[0]["USER_NM"]
	params["goodsName"] = goodsName

	mallReserved := moid + "," + params["amt"] + "," + params["restId"] + "," + params["grpId"] + "," + params["searchTy"] + "," + params["addAmt"] + "," + params["payChannel"] + "," + pgCd + "," + userTy + "," + selectedDate + "," + params["paymentTy"] + "," + params["osTy"]

	cip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)

	params["selectedDate"] = selectedDate
	params["encryptData"] = encryptData
	params["ediDate"] = ediDate
	params["vbankExpDate"] = vbankExpDate
	params["connType"] = connType
	params["mid"] = mid
	params["moid"] = moid
	params["merchantKey"] = merchantKey
	params["payActionUrl"] = payActionUrl
	params["payLocalUrl"] = payLocalUrl
	params["returnUrl"] = returnUrl
	params["cancelUrl"] = cancelUrl
	params["mallReserved"] = mallReserved
	params["paymentTy"] = paymentTy
	params["selectedDate"] = selectedDate
	params["appPrefix"] = appPrefix
	params["userAgent"] = userAgent
	params["clientIp"] = cip
	params["mallIp"] = ""

	page := "tpays/tpayReady.htm"

	if params["osTy"] == "web" {
		params["payActionWEBUrl"] = payActionWEBUrl
		page = "tpays/tpayReadyWeb.htm"
	}

	m := make(map[string]interface{})
	m["tpay"] = params

	return c.Render(http.StatusOK, page, m)
}

func TpayResult(c echo.Context) error {

	dprintf(4, c, "call TpayResult\n")

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)

	mallReserved := params["mallReserved"]
	ediDate := params["ediDate"]

	decAmt := AesHandlerDecrypt(params["amt"], ediDate, merchantKey)
	decMoid := AesHandlerDecrypt(params["moid"], ediDate, merchantKey)

	arrMallReserved := strings.Split(mallReserved, ",")

	//println(mallReserved)

	sndMoId := arrMallReserved[0]
	sndAmt := arrMallReserved[1]
	sndRestId := arrMallReserved[2]
	sndGrpId := arrMallReserved[3]
	sndSearchTy := arrMallReserved[4]
	sndAddamt := arrMallReserved[5]
	sndPayChannel := arrMallReserved[6]
	sndPgCd := arrMallReserved[7]
	sndUserTy := arrMallReserved[8]
	sndSelectedDate := arrMallReserved[9]
	sndPaymentTy := arrMallReserved[10]
	osTy := arrMallReserved[11]

	pushAcDate := ""

	m := make(map[string]interface{})
	m["osTy"] = osTy
	m["chargeAmt"] = sndAmt
	//return c.Render(http.StatusOK, "tpays/payResult.htm", m)

	if decMoid != sndMoId {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (moid 다름) :", decMoid, sndMoId)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
	}

	if decAmt != sndAmt {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (금액 다름) :", decAmt, sndAmt)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
	}

	if sndSearchTy == "2" {
		//sndPaymentTy="3"
	}

	paymethod := params["payMethod"]
	resultCd := params["resultCd"]
	resultMsg := params["resultMsg"]
	params["restId"] = sndRestId
	params["grpId"] = sndGrpId
	params["decMoid"] = decMoid

	//중복 체크
	dupChk, err := cls.GetSelectDataRequire(paymentsql.SelectTPayDupCheck, params, c)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "결제 결과 반영 오류"
		return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
	}

	restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestType, params, c)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "결제 결과 반영 오류"
		return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
	}

	restType := restTypeInfo[0]["REST_TYPE"]
	chargeAmt := restTypeInfo[0]["CHARGE_AMT"]

	dupCnt, _ := strconv.Atoi(dupChk[0]["DUP_CNT"])
	if dupCnt > 0 {
		m["resultCode"] = "99"
		return c.Render(http.StatusOK, "tpays/resultDup.htm", m)
	}

	if (paymethod == "CARD" && resultCd == "3001") || (paymethod == "BANK" && resultCd == "4000") {

		pgParam := make(map[string]string)

		if paymethod == "CARD" {
			pgParam["payInfo"] = "1"
		} else {
			pgParam["payInfo"] = "0"
		}

		tx, err := cls.DBc.Begin()
		if err != nil {
			dprintf(4, c, " 결제 결과 반영 TpayResult 오류  \n")
			return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
		}

		txErr := err

		defer func() {
			if txErr != nil {
				// transaction rollback
				dprintf(4, c, "do rollback - 결제 결과 반영 TpayResult)  \n")

				payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
				tx.Rollback()
			}
		}()

		pgParam["moid"] = sndMoId
		pgParam["paymethod"] = paymethod
		pgParam["transtype"] = "0"
		pgParam["goodsname"] = params["goodsName"]
		pgParam["amt"] = sndAmt
		pgParam["addAmt"] = sndAddamt
		pgParam["userip"] = "0"
		pgParam["tid"] = params["tid"]
		pgParam["state"] = "0000"
		pgParam["statecd"] = params["stateCd"]
		pgParam["cardno"] = params["cardNo"]
		pgParam["authcode"] = params["authCode"]
		pgParam["authdate"] = params["authDate"]
		pgParam["cardquota"] = params["cardQuota"]
		pgParam["fncd"] = params["fnCd"]
		pgParam["fnname"] = params["fnName"]
		pgParam["resultcd"] = params["resultCd"]
		pgParam["resultmsg"] = params["resultMsg"]
		pgParam["pgCd"] = sndPgCd

		pgParam["histId"] = sndMoId
		pgParam["restId"] = sndRestId
		pgParam["grpId"] = sndGrpId
		pgParam["userId"] = params["mallUserId"]
		pgParam["creditAmt"] = sndAmt
		pgParam["userTy"] = sndUserTy
		pgParam["searchTy"] = sndSearchTy
		pgParam["paymentTy"] = sndPaymentTy
		pgParam["payChannel"] = sndPayChannel
		pgParam["selectedDate"] = sndSelectedDate

		// 결제정보 DB 저장1 : DAR_PAYMENT_REPORT
		insertPaymentReportQuery, err := cls.GetQueryJson(paymentsql.InsertTpayPayment, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentReportQuery parameter fail"))
		}
		_, err = tx.Exec(insertPaymentReportQuery)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentReportQuery, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
		}

		insertPaymentHistory, err := cls.GetQueryJson(paymentsql.InsertTpayPaymentHistory, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentHistory parameter fail"))
		}
		_, err = tx.Exec(insertPaymentHistory)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentHistory, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
		}

		if restType == "G" {

			isndAmt, _ := strconv.Atoi(sndAmt)
			ichargeAmt, _ := strconv.Atoi(chargeAmt)
			pgParam["chargeAmt"] = strconv.Itoa(ichargeAmt + isndAmt)

			UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombineChargeAmt, pgParam)
			if err != nil {
				txErr = err
				return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
			}
			_, err = tx.Exec(UpdateCombineChargeAmtQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}

		}

		if sndSearchTy == "1" {
			//선불처리

			agrmInfo, err := cls.GetSelectDataRequire(paymentsql.SelectAgrmInfo, pgParam, c)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}

			prepaidAmt := agrmInfo[0]["PREPAID_AMT"]
			prepaidPoint := agrmInfo[0]["PREPAID_POINT"]
			pgParam["agrmId"] = agrmInfo[0]["AGRM_ID"]

			intAmt, _ := strconv.Atoi(sndAmt)
			intAddAmt, _ := strconv.Atoi(sndAddamt)
			intPrepaidAmt, _ := strconv.Atoi(prepaidAmt)
			pgParam["prepaidAmt"] = strconv.Itoa(intPrepaidAmt + intAmt + intAddAmt)

			intPrepaidPoint, _ := strconv.Atoi(prepaidPoint)
			pgParam["prepaidPoint"] = strconv.Itoa(intPrepaidPoint + intAddAmt)

			updateAgrmPrepaidAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateAgrmPrepaidAmt, pgParam)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}
			_, err = tx.Exec(updateAgrmPrepaidAmtQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", updateAgrmPrepaidAmtQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}

			pgParam["jobTy"] = "0"
			pgParam["prepaidAmt"] = strconv.Itoa(intAmt + intAddAmt)
			InsertPrepaidQuery, err := cls.GetQueryJson(paymentsql.InsertPrepaid, pgParam)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}
			_, err = tx.Exec(InsertPrepaidQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPrepaidQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}

		} else if sndSearchTy == "2" {
			//후불 처리
			unpaidInfo, err := cls.GetSelectDataRequire(paymentsql.SelectUnpaidPaymentInfo, pgParam, c)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}

			pushAcDate = unpaidInfo[0]["AC_DATE"]

			if sndAmt != unpaidInfo[0]["CREDIT_AMT"] {

				payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
				tx.Rollback()
				dprintf(1, c, "정산금액과 결제 금액이 다릅니다.\n", err)
				m["resultCode"] = "98"
				m["resultMsg"] = "정산금액과 결제 금액이 다릅니다"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)

			}

			UnpaidPaymentCreditNowQuery, err := cls.GetQueryJson(paymentsql.UnpaidPaymentCreditNow, pgParam)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UnpaidPaymentCreditNowQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}
			_, err = tx.Exec(UnpaidPaymentCreditNowQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UnpaidPaymentCreditNowQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
			}

		} else {
			payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
			tx.Rollback()

			dprintf(1, c, "결제 방식 오류로 인하여 취소하였습니다.\n", err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 방식 오류"
			return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
		}

		err = tx.Commit()
		if err != nil {
			txErr = err
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
		}

	} else {
		m["resultCode"] = "99"
		m["resultMsg"] = resultMsg
		return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
	}

	sendResultConfirm(params["tid"], "000")

	params["sndGrpId"] = sndGrpId
	payInfo, err := cls.GetSelectDataRequire(paymentsql.SelectPayInfo, params, c)
	if err != nil {
		m["resultCode"] = "98"
		return c.Render(http.StatusOK, "tpays/resultFail.htm", m)
	}

	m["goodsName"] = payInfo[0]["GOODSNAME"]
	m["userNm"] = payInfo[0]["USER_NM"]

	// 푸쉬 전송 시작

	if restType == "N" {
		pushMsg := ""
		if sndSearchTy == "1" {
			pushMsg = payInfo[0]["GRP_NM"] + "의 " + payInfo[0]["USER_NM"] + "님이 " + sndAmt + "원을 충전(결제)하였습니다."
			apiPush.SendPush_Msg_V1("충전", pushMsg, "M", "1", sndRestId, "", "bookpay")

			lprintf(4, "충전 push 메세지  :", pushMsg)
		} else if sndSearchTy == "2" {
			pushMsg = payInfo[0]["GRP_NM"] + "의 " + payInfo[0]["USER_NM"] + "님이 " + pushAcDate + "까지 사용한 금액을 정산(결제)하였습니다."
			apiPush.SendPush_Msg_V1("정산", pushMsg, "M", "1", sndRestId, "", "bookpay")
			lprintf(4, "정산 push 메세지  :", pushMsg)
		}

	}

	// 푸쉬 전송 완료

	page := "tpays/payResult.htm"

	if osTy == "web" {
		page = "tpays/payResultWeb.htm"
		m["grpId"] = sndGrpId
		m["restId"] = sndRestId
		m["userId"] = params["mallUserId"]
		m["searchTy"] = sndSearchTy
	}

	return c.Render(http.StatusOK, page, m)
}

func sendResultConfirm(tid string, result string) {

	payload := strings.NewReader("tid=" + tid + "&result=" + result)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", "https://webtx.tpay.co.kr/resultConfirm", payload)
	if err != nil {
		lprintf(1, "[ERROR] Tpay resultConfirm Send1 : %s\n", err)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		lprintf(1, "[ERROR] Tpay resultConfirm Send2 : %s\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] Tpay resultConfirm Send3 : %s\n", err)
	}
	lprintf(4, "[INFO] Tpay resultConfirm resp(%s)\n", string(body))
	println(string(body))

}

func sendCancelResultConfirm(apikey string, mid string, tid string, result string) {

	payload := strings.NewReader("api_key=" + apikey + "&mid=" + mid + "&tid=" + tid + "&result_code=" + result + "&result_cd=" + result + "&result=" + result)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", "https://webtx.tpay.co.kr/api/v1/result_confirm", payload)
	if err != nil {
		lprintf(1, "[ERROR] Tpay Cancel resultConfirm Send1 : %s\n", err)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		lprintf(1, "[ERROR] Tpay Cancel resultConfirm Send2 : %s\n", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] Tpay Cancel resultConfirm Send3 : %s\n", err)
	}
	lprintf(4, "[INFO] Tpay Cancel resultConfirm resp(%s)\n", string(body))
	println(string(body))

}

func TpayCancel(c echo.Context) error {

	payCancel("", "", "", "", "")

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	return c.JSON(http.StatusOK, m)

}

func billingPayCancel(moid string, cancel_amt string, cancel_msg string, tid string) map[string]string {

	cancel_pw := controller.TPAY_SIMPLE_PAY_CANCEL_PWD
	mid := controller.TPAY_MID_SIMPLE_PAY
	merchantKey := controller.TPAY_SIMPLE_PAY_MERCHANT_KEY
	cancelUrl := controller.TPAY_SIMPLE_PAY_CANCEL_URL

	partial_cancel := "0"

	uValue := url.Values{
		"api_key":        {merchantKey},
		"mid":            {mid},
		"moid":           {moid},
		"cancel_pw":      {cancel_pw},
		"cancel_amt":     {cancel_amt},
		"cancel_msg":     {cancel_msg},
		"partial_cancel": {partial_cancel},
		"tid":            {tid},
	}

	m := make(map[string]string)
	resp, err := http.PostForm(cancelUrl, uValue)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링결제 취소 오류"
		return m
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링결제 취소 오류"
		return m
	}

	println(string(respBody))

	var result TpayRecvData
	err = json.Unmarshal(respBody, &result)
	if err != nil {

		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링결제 취소 오류"
		return m
	}

	if result.Result_cd != "2001" {
		m["resultCode"] = "99"
		m["resultMsg"] = result.Result_msg
	} else {
		m["resultCode"] = "00"
		m["resultMsg"] = result.Result_msg
		m["tid"] = result.Tid
		m["CancelDate"] = result.CancelDate
		m["CancelTime"] = result.CancelTime

		sendCancelResultConfirm(merchantKey, mid, result.Tid, "000")
	}
	return m

}

func payCancel(moid string, cancelAmt string, cancelMsg string, partialCancelCode string, tid string) map[string]string {

	lprintf(4, "call 결제 취소요청  \n")
	lprintf(4, "moid:  \n", moid)
	lprintf(4, "tid:  \n", tid)

	TpayMap := make(map[string]string)

	fname := cls.Cls_conf(os.Args)
	mid, _ := cls.GetTokenValue("TPAY.TPAY_MID", fname)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)
	cancelPw, _ := cls.GetTokenValue("TPAY.TPAY_CANCEL_PW", fname)
	restfulCancelUrl, _ := cls.GetTokenValue("TPAY.RESTFUL_CANCEL_URL", fname)

	uValue := url.Values{
		"api_key":        {merchantKey},
		"mid":            {mid},
		"moid":           {moid},
		"cancel_pw":      {cancelPw},
		"cancel_amt":     {cancelAmt},
		"cancel_msg":     {cancelMsg},
		"partial_cancel": {partialCancelCode},
		"tid":            {tid},
	}

	resp, err := http.PostForm(restfulCancelUrl, uValue)
	if err != nil {
		lprintf(4, "결제 취소 요청 오류 :  ", err)
		panic(err)
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		TpayMap["result_cd"] = "9999"
		TpayMap["result_msg"] = "전송 요청 오류"
		lprintf(4, "결제 취소 요청 오류 : ", err)
		return TpayMap
	}

	var result TpayRecvData
	err = json.Unmarshal(respBody, &result)
	if err != nil {

		TpayMap["result_cd"] = "9999"
		TpayMap["result_msg"] = "전송 요청 오류"
		return TpayMap
	}

	lprintf(4, "결제 취소 결과 : ", result.Result_msg)

	TpayMap["result_cd"] = result.Result_cd
	TpayMap["result_msg"] = result.Result_msg
	TpayMap["PayMethod"] = result.PayMethod
	TpayMap["CancelDate"] = result.CancelDate
	TpayMap["CancelTime"] = result.CancelTime
	TpayMap["tid"] = result.Tid
	TpayMap["moid"] = result.Moid

	sendCancelResultConfirm(merchantKey, mid, result.Tid, "000")

	return TpayMap

}

func getEdiDate() (string, string) {
	now := time.Now()

	bankDate := now.AddDate(0, 0, 1).Format("20060102")

	return now.Format("20060102150405"), bankDate
}

func getMoid(userId string) string {

	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000

	return fmt.Sprintf("%d%s", millis, userId)
}

func getMD5HashHandler(key string) string {

	hash := md5.Sum([]byte(key))
	hash2 := md5.Sum([]byte(key))
	var hexString string

	//fmt.Println(text, "  hash len : ", len(hash))

	for i := 0; i < len(hash); i++ {
		var tmp []byte

		//fmt.Println("hash i : ", hash[i], " 0xFF : ", 0xFF&hash[i])

		tmp = append(tmp, uint8(hash[i]&0xFF))
		hexEncode := hex.EncodeToString(tmp)
		if len(hexEncode) == 1 {
			hexEncode = "0" + hexEncode
		}

		hexString += hexEncode
	}

	//fmt.Println("hex    ", hex.EncodeToString(hash2[:]))
	//fmt.Println("encKey  ", hexString, " len ", len(hexString))

	if hex.EncodeToString(hash2[:]) == hexString {
		//fmt.Println("hex true")
	} else {
		//fmt.Println("hex false")
	}

	return hexString
}

// tpay decrypt
func AesHandlerDecrypt(input, ediDate, merchantKey string) string {

	key := fmt.Sprintf("%s%s", ediDate, merchantKey)
	encKey := getMD5HashHandler(key)
	decode, _ := base64.StdEncoding.DecodeString(input)

	result := AesDecrypt(byteTostring(decode), encKey, merchantKey)

	return string(result)
}

func byteTostring(buf []byte) string {

	var tmp string

	for i := 0; i < len(buf); i++ {
		if buf[i]&0xff < 0x10 {
			tmp += "0"
		}

		tmp += strconv.FormatInt(int64(buf[i]&0xff), 16)
	}

	return tmp
}

func AesDecrypt(input, key, merchantKey string) []byte {

	ips := []byte(merchantKey[:16])
	keyBytes := hexToByteArray(key)
	text := hexToByteArray(input)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		fmt.Println("key error1", err)
		return nil
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, ips[:blockSize])

	crypted := make([]byte, len(text))
	blockMode.CryptBlocks(crypted, text)

	return PKCS5UnPadding(crypted)
}

func AesHandlerEncrypt(input, key, merchantKey string) []byte {

	inbytes := []byte(merchantKey[:16])

	return AesEncrypt([]byte(input), hexToByteArray(key), inbytes)
}

func hexToByteArray(hex string) []byte {

	bb := make([]byte, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		pi, err := strconv.ParseInt(hex[i:i+2], 16, 16)
		if err != nil {
			fmt.Println(err.Error())
		}
		bb[int(math.Floor(float64(i/2)))] = uint8(pi)
	}

	return bb
}

func AesEncrypt(origData, key, IV []byte) []byte {

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
		return nil
	}

	//fmt.Println("blockSize : ", block.BlockSize())

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, IV[:blockSize])

	origData = PKCS5Padding(origData, blockSize)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return crypted
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func TpayCancelPage(c echo.Context) error {

	//params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.Render(http.StatusOK, "tpays/payCancel.htm", m)
}

func MngPayCharge(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["grpId"] = params["grpId"]
	m["restId"] = params["restId"]
	m["searchTy"] = params["searchTy"]

	return c.Render(http.StatusOK, "mngOrder/payCharge.htm", m)
}

// 웹결제 준비 데이터
func TpayReadyWebData(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	mid, _ := cls.GetTokenValue("TPAY.TPAY_MID", fname)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)
	payLocalUrl, _ := cls.GetTokenValue("TPAY.PAY_LOCAL_URL", fname)
	returnUrl, _ := cls.GetTokenValue("TPAY.RETURN_URL", fname)
	cancelUrl, _ := cls.GetTokenValue("TPAY.CANCEL_URL", fname)
	domain, _ := cls.GetTokenValue("TPAY.DOMAIN", fname)

	paymentTy := params["paymentTy"]
	selectedDate := ""
	appPrefix := "mngAdmin://"
	selectedDate = params["selectedDate"]
	userAgent := c.Request().Header.Get("User-Agent")

	connType := "1"
	ediDate, vbankExpDate := getEdiDate()
	moid := getMoid(params["userId"])
	//pgCd := "01"
	//userTy := "0"

	var tpay Tpay
	tpay.MerchantKey = merchantKey
	tpay.EdiDate = ediDate

	key := fmt.Sprintf("%s%s", ediDate, merchantKey)
	tpay.EncKey = getMD5HashHandler(key)

	input := fmt.Sprintf("%s%s%s", params["amt"], mid, moid)
	encryptData := base64.StdEncoding.EncodeToString(AesHandlerEncrypt(input, tpay.EncKey, merchantKey))

	params["mallUserId"] = params["userId"]
	userInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	goodsName := ""

	if params["searchTy"] == "1" {
		paymentTy = "0"
		params["paymentTy"] = "0"
	} else {
		paymentTy = "3"
		params["paymentTy"] = "3"
	}

	if params["searchTy"] == "1" {
		goodsName = userInfo[0]["REST_NM"] + " 충전"
	} else {
		goodsName = userInfo[0]["REST_NM"] + " 정산"
	}

	params["buyerEmail"] = userInfo[0]["EMAIL"]
	params["buyerTel"] = userInfo[0]["HP_NO"]
	params["userNm"] = userInfo[0]["USER_NM"]
	params["goodsName"] = goodsName

	//mallReserved := moid + "," + params["amt"] + "," + params["restId"] + "," + params["grpId"] + "," + params["searchTy"] + "," + params["addAmt"] + "," + params["payChannel"] + "," + pgCd + "," + userTy + "," + selectedDate + "," + params["paymentTy"] + "," + params["osTy"]

	cip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)

	params["encryptData"] = encryptData
	params["ediDate"] = ediDate
	params["vbankExpDate"] = vbankExpDate
	params["connType"] = connType
	params["mid"] = mid
	params["moid"] = moid
	params["merchantKey"] = merchantKey
	params["payLocalUrl"] = payLocalUrl
	params["returnUrl"] = returnUrl
	params["cancelUrl"] = cancelUrl
	//params["mallReserved"] = mallReserved
	params["paymentTy"] = paymentTy
	params["selectedDate"] = selectedDate
	params["appPrefix"] = appPrefix
	params["userAgent"] = userAgent
	params["clientIp"] = cip
	params["mallIp"] = ""
	params["domain"] = domain

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["tpay"] = params

	return c.JSON(http.StatusOK, m)
}

func TpayUnpaid(c echo.Context) error {

	dprintf(4, c, "call TpayUnpaidResult\n")

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)

	mallReserved := params["mallReserved"]
	ediDate := params["ediDate"]

	decAmt := AesHandlerDecrypt(params["amt"], ediDate, merchantKey)
	decMoid := AesHandlerDecrypt(params["moid"], ediDate, merchantKey)

	arrMallReserved := strings.Split(mallReserved, ",")

	//println(mallReserved)

	sndMoId := arrMallReserved[0]
	sndAmt := arrMallReserved[1]
	sndRestId := arrMallReserved[2]
	sndGrpId := arrMallReserved[3]
	sndSearchTy := arrMallReserved[4]
	sndAddamt := arrMallReserved[5]
	sndPayChannel := arrMallReserved[6]
	sndPgCd := arrMallReserved[7]
	sndUserTy := arrMallReserved[8]
	sndSelectedDate := arrMallReserved[9]
	sndPaymentTy := arrMallReserved[10]
	osTy := arrMallReserved[11]

	pushAcDate := ""

	m := make(map[string]interface{})
	m["osTy"] = osTy
	m["chargeAmt"] = sndAmt
	//return c.Render(http.StatusOK, "tpays/payResult.htm", m)

	if decMoid != sndMoId {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (moid 다름) :", decMoid, sndMoId)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	if decAmt != sndAmt {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (금액 다름) :", decAmt, sndAmt)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	if sndSearchTy == "2" {
		//sndPaymentTy="3"
	}

	paymethod := params["payMethod"]
	resultCd := params["resultCd"]
	resultMsg := params["resultMsg"]
	params["restId"] = sndRestId
	params["grpId"] = sndGrpId
	params["decMoid"] = decMoid

	//중복 체크
	dupChk, err := cls.GetSelectDataRequire(paymentsql.SelectTPayDupCheck, params, c)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "결제 결과 반영 오류"
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	dupCnt, _ := strconv.Atoi(dupChk[0]["DUP_CNT"])
	if dupCnt > 0 {
		m["resultCode"] = "99"
		return c.Render(http.StatusOK, "tpays/resultDupWeb.htm", m)
	}

	if (paymethod == "CARD" && resultCd == "3001") || (paymethod == "BANK" && resultCd == "4000") {

		pgParam := make(map[string]string)

		if paymethod == "CARD" {
			pgParam["payInfo"] = "1"
		} else {
			pgParam["payInfo"] = "0"
		}

		tx, err := cls.DBc.Begin()
		if err != nil {
			dprintf(4, c, " 결제 결과 반영 TpayResult 오류  \n")
			//return "5100", errors.New("begin error")

		}
		txErr := err
		defer func() {
			if txErr != nil {
				// transaction rollback
				dprintf(4, c, "do rollback - 결제 결과 반영 실패 TpayUnpaidResult)  \n")
				dprintf(4, c, "do rollback - 결제 결과 반영 실패 메세지 TpayUnpaidResult)  \n", txErr.Error())
				payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
				tx.Rollback()
			}
		}()

		pgParam["moid"] = sndMoId
		pgParam["paymethod"] = paymethod
		pgParam["transtype"] = "0"
		pgParam["goodsname"] = params["goodsName"]
		pgParam["amt"] = sndAmt
		pgParam["addAmt"] = sndAddamt
		pgParam["userip"] = "0"
		pgParam["tid"] = params["tid"]
		pgParam["state"] = "0000"
		pgParam["statecd"] = params["stateCd"]
		pgParam["cardno"] = params["cardNo"]
		pgParam["authcode"] = params["authCode"]
		pgParam["authdate"] = params["authDate"]
		pgParam["cardquota"] = params["cardQuota"]
		pgParam["fncd"] = params["fnCd"]
		pgParam["fnname"] = params["fnName"]
		pgParam["resultcd"] = params["resultCd"]
		pgParam["resultmsg"] = params["resultMsg"]
		pgParam["pgCd"] = sndPgCd

		pgParam["histId"] = sndMoId
		//pgParam["restId"] = sndRestId
		pgParam["grpId"] = sndGrpId
		pgParam["userId"] = params["mallUserId"]
		pgParam["creditAmt"] = sndAmt
		pgParam["userTy"] = sndUserTy
		pgParam["searchTy"] = sndSearchTy
		pgParam["paymentTy"] = sndPaymentTy
		pgParam["payChannel"] = sndPayChannel
		pgParam["selectedDate"] = sndSelectedDate

		// 결제정보 DB 저장1 : DAR_PAYMENT_REPORT
		insertPaymentReportQuery, err := cls.GetQueryJson(paymentsql.InsertTpayPayment, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentReportQuery parameter fail"))
		}
		_, err = tx.Exec(insertPaymentReportQuery)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentReportQuery, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류(report)"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
		}

		chkAmt := 0

		if sndSearchTy == "2" {

			beneficonUpdate := false
			beneficonCargeAmt := 0

			wincubeUpdate := false
			wincubeCargeAmt := 0

			strArray := strings.Split(sndRestId, "@")

			for _, str := range strArray {

				pgParam["restId"] = str

				//	println(str)
				//후불 처리
				unpaidInfo, err := cls.GetSelectDataRequire(paymentsql.SelectUnpaidPaymentInfo, pgParam, c)
				if err != nil {
					txErr = err
					m["resultCode"] = "98"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}
				pgParam["creditAmt"] = unpaidInfo[0]["CREDIT_AMT"]
				restType := unpaidInfo[0]["REST_TYPE"]
				creditAmt, _ := strconv.Atoi(unpaidInfo[0]["CREDIT_AMT"])

				//기프티콘 충전금액 데이터 모으기
				if restType == "G" {
					company := unpaidInfo[0]["CEO_NM"]
					if company == "Wincube" {
						wincubeUpdate = true
						wincubeCargeAmt = wincubeCargeAmt + creditAmt
					} else if company == "BENEPICON" {
						beneficonUpdate = true
						beneficonCargeAmt = beneficonCargeAmt + creditAmt
					} else {
						tx.Rollback()
						dprintf(1, c, "[Error] 잘못된 기프티콘 회사 데이터  \n", err)
						//payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
						m["resultCode"] = "98"
						m["resultMsg"] = "결제 결과 반영 오류(gift)"
						return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
					}
				}

				chkAmt = chkAmt + creditAmt
				pgParam["histId"] = sndMoId + "_" + str

				//pgParam["histId"] = pgParam["histId"]+"_"+str

				insertPaymentHistory, err := cls.GetQueryJson(paymentsql.InsertTpayPaymentHistory, pgParam)
				if err != nil {
					txErr = err
					return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentHistory parameter fail"))
				}
				_, err = tx.Exec(insertPaymentHistory)
				if err != nil {
					txErr = err
					dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentHistory, err)
					//payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
					m["resultCode"] = "98"
					m["resultMsg"] = "결제 결과 반영 오류(hist)"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}

				pushAcDate = unpaidInfo[0]["AC_DATE"]

				UnpaidPaymentCreditNowQuery, err := cls.GetQueryJson(paymentsql.UnpaidPaymentCreditNow, pgParam)
				if err != nil {
					txErr = err
					m["resultCode"] = "98"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}
				_, err = tx.Exec(UnpaidPaymentCreditNowQuery)
				if err != nil {
					txErr = err
					dprintf(1, c, "Query(%s) -> error (%s) \n", UnpaidPaymentCreditNowQuery, err)
					//payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
					m["resultCode"] = "98"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}
			}

			if wincubeUpdate == true {

				wincubeParam := make(map[string]string)
				wincubeParam["restId"] = "S0000000649"

				restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestCombine, wincubeParam, c)
				if err != nil {
					txErr = err
					m["resultCode"] = "99"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}
				chargeAmt := restTypeInfo[0]["CHARGE_AMT"]
				ichargeAmt, _ := strconv.Atoi(chargeAmt)
				wincubeParam["chargeAmt"] = strconv.Itoa(ichargeAmt + wincubeCargeAmt)

				UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombinePaidChargeAmt, wincubeParam)
				if err != nil {
					txErr = err
					return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
				}
				_, err = tx.Exec(UpdateCombineChargeAmtQuery)
				if err != nil {
					txErr = err
					dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
					m["resultCode"] = "98"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}

			}

			if beneficonUpdate == true {

				beneficonParam := make(map[string]string)
				beneficonParam["restId"] = "S0000000615"

				restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestCombine, beneficonParam, c)
				if err != nil {
					txErr = err
					m["resultCode"] = "99"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}
				chargeAmt := restTypeInfo[0]["CHARGE_AMT"]
				ichargeAmt, _ := strconv.Atoi(chargeAmt)
				beneficonParam["chargeAmt"] = strconv.Itoa(ichargeAmt + beneficonCargeAmt)

				UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombinePaidChargeAmt, beneficonParam)
				if err != nil {
					txErr = err
					return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
				}
				_, err = tx.Exec(UpdateCombineChargeAmtQuery)
				if err != nil {
					txErr = err
					dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
					m["resultCode"] = "98"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}

			}

		} else {
			payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
			tx.Rollback()

			dprintf(1, c, "결제 방식 오류로 인하여 취소하였습니다.\n", err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 방식 오류"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
		}

		if sndAmt != strconv.Itoa(chkAmt) {
			payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소(금액상이)", "0", params["tid"])
			tx.Rollback()

			dprintf(1, c, "정산금액과 결제 금액이 다릅니다.\n", err)
			m["resultCode"] = "98"
			m["resultMsg"] = "정산금액과 결제 금액이 다릅니다"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)

		}

		err = tx.Commit()
		if err != nil {
			txErr = err
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
		}

	} else {

		m["resultCode"] = "99"
		m["resultMsg"] = resultMsg
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	sendResultConfirm(params["tid"], "000")

	params["sndGrpId"] = sndGrpId
	payInfo, err := cls.GetSelectDataRequire(paymentsql.SelectPayInfo, params, c)
	if err != nil {
		m["resultCode"] = "98"
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	m["goodsName"] = payInfo[0]["GOODSNAME"]
	m["userNm"] = payInfo[0]["USER_NM"]

	// 푸쉬 전송 시작

	if len(os.Getenv("SERVER_TYPE")) > 0 {

		pushMsg := ""

		strArray2 := strings.Split(sndRestId, "@")

		for _, str2 := range strArray2 {

			pushMsg = payInfo[0]["GRP_NM"] + "의 " + payInfo[0]["USER_NM"] + "님이 " + pushAcDate + "까지 사용한 금액을 정산(결제)하였습니다."
			apiPush.SendPush_Msg_V1("정산", pushMsg, "M", "1", str2, "", "bookpay")
			lprintf(4, "정산 push 메세지  :", pushMsg)

		}
	}

	// 푸쉬 전송 완료

	page := "tpays/payResultWebUnpaid.htm"
	m["grpId"] = sndGrpId
	m["restId"] = sndRestId
	m["userId"] = params["mallUserId"]
	m["searchTy"] = sndSearchTy

	return c.Render(http.StatusOK, page, m)
}

func TpayUnpaidNew(c echo.Context) error {

	dprintf(4, c, "call TpayUnpaidResult\n")

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)

	mallReserved := params["mallReserved"]
	ediDate := params["ediDate"]

	decAmt := AesHandlerDecrypt(params["amt"], ediDate, merchantKey)
	decMoid := AesHandlerDecrypt(params["moid"], ediDate, merchantKey)

	arrMallReserved := strings.Split(mallReserved, ",")

	//println(decAmt)

	sndMoId := arrMallReserved[0]
	sndAmt := arrMallReserved[1]
	sndRestId := arrMallReserved[2]
	sndGrpId := arrMallReserved[3]
	sndSearchTy := arrMallReserved[4]
	sndAddamt := arrMallReserved[5]
	sndPayChannel := arrMallReserved[6]
	sndPgCd := arrMallReserved[7]
	sndUserTy := arrMallReserved[8]
	sndSelectedDate := arrMallReserved[9]
	sndPaymentTy := arrMallReserved[10]
	osTy := arrMallReserved[11]

	pushAcDate := ""

	m := make(map[string]interface{})
	m["osTy"] = osTy
	m["chargeAmt"] = sndAmt
	//return c.Render(http.StatusOK, "tpays/payResult.htm", m)

	if decMoid != sndMoId {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (moid 다름) :", decMoid, sndMoId)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	if decAmt != sndAmt {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (금액 다름) :", decAmt, sndAmt)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제 금액이 잘못되었습니다."
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	if sndSearchTy == "2" {
		//sndPaymentTy="3"
	}

	paymethod := params["payMethod"]
	resultCd := params["resultCd"]
	resultMsg := params["resultMsg"]
	params["restId"] = sndRestId
	params["grpId"] = sndGrpId
	params["decMoid"] = decMoid

	//중복 체크
	dupChk, err := cls.GetSelectDataRequire(paymentsql.SelectTPayDupCheck, params, c)
	if err != nil {
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
		m["resultCode"] = "99"
		m["resultMsg"] = "정산금액과 결제 금액이 다릅니다."
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	dupCnt, _ := strconv.Atoi(dupChk[0]["DUP_CNT"])
	if dupCnt > 0 {
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
		m["resultCode"] = "99"
		m["resultMsg"] = "중복 결제 입니다."
		return c.Render(http.StatusOK, "tpays/resultDupWeb.htm", m)
	}

	if (paymethod == "CARD" && resultCd == "3001") || (paymethod == "BANK" && resultCd == "4000") {

		//결제 금액 체크
		PaidOrderAmtCheck, err := cls.GetSelectDataRequire(paymentsql.SelectPaidOrderAmtCheck, params, c)
		if err != nil {
			payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
			m["resultCode"] = "99"
			m["resultMsg"] = "결제 금액 체크 오류"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
		}
		chkTotalAmt, _ := PaidOrderAmtCheck[0]["TOTAL_AMT"]
		if chkTotalAmt != sndAmt {
			lprintf(4, "결제 취소 요청 (금액 다름) :", chkTotalAmt, sndAmt)
			payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

			m["resultCode"] = "99"
			m["resultMsg"] = "결제 금액 오류"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
		}

		pgParam := make(map[string]string)

		if paymethod == "CARD" {
			pgParam["payInfo"] = "1"
		} else {
			pgParam["payInfo"] = "0"
		}

		tx, err := cls.DBc.Begin()
		if err != nil {
			dprintf(4, c, " 일괄 결제 결과 반영 TpayUnpaidResult 시작 오류 \n")
			//return "5100", errors.New("begin error")

		}

		pgParam["moid"] = sndMoId
		pgParam["paymethod"] = paymethod
		pgParam["transtype"] = "0"
		pgParam["goodsname"] = params["goodsName"]
		pgParam["amt"] = sndAmt
		pgParam["addAmt"] = sndAddamt
		pgParam["userip"] = "0"
		pgParam["tid"] = params["tid"]
		pgParam["state"] = "0000"
		pgParam["statecd"] = params["stateCd"]
		pgParam["cardno"] = params["cardNo"]
		pgParam["authcode"] = params["authCode"]
		pgParam["authdate"] = params["authDate"]
		pgParam["cardquota"] = params["cardQuota"]
		pgParam["fncd"] = params["fnCd"]
		pgParam["fnname"] = params["fnName"]
		pgParam["resultcd"] = params["resultCd"]
		pgParam["resultmsg"] = params["resultMsg"]
		pgParam["pgCd"] = sndPgCd

		pgParam["histId"] = sndMoId
		//pgParam["restId"] = sndRestId
		pgParam["grpId"] = sndGrpId
		pgParam["userId"] = params["mallUserId"]
		pgParam["creditAmt"] = sndAmt
		pgParam["userTy"] = sndUserTy
		pgParam["searchTy"] = sndSearchTy
		pgParam["paymentTy"] = sndPaymentTy
		pgParam["payChannel"] = sndPayChannel
		pgParam["selectedDate"] = sndSelectedDate

		beneficonUpdate := false
		beneficonCargeAmt := 0

		wincubeUpdate := false
		wincubeCargeAmt := 0

		strArray := strings.Split(sndRestId, "@")

		for _, str := range strArray {

			giftParam := make(map[string]string)

			giftParam["restId"] = str
			giftParam["grpId"] = sndGrpId
			giftParam["selectedDate"] = sndSelectedDate

			//	println(str)
			//후불 처리
			unpaidInfo, err := cls.GetSelectDataRequire(paymentsql.SelectUnpaidPaymentInfo, giftParam, c)
			if err != nil {
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
			}
			restType := unpaidInfo[0]["REST_TYPE"]
			creditAmt, _ := strconv.Atoi(unpaidInfo[0]["CREDIT_AMT"])

			//기프티콘 충전금액 데이터 모으기
			if restType == "G" {
				company := unpaidInfo[0]["CEO_NM"]
				if company == "Wincube" {
					wincubeUpdate = true
					wincubeCargeAmt = wincubeCargeAmt + creditAmt
				} else if company == "BENEPICON" {
					beneficonUpdate = true
					beneficonCargeAmt = beneficonCargeAmt + creditAmt
				} else {
					dprintf(1, c, "[Error] 잘못된 기프티콘 회사 데이터  \n", err)
					payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
					m["resultCode"] = "98"
					m["resultMsg"] = "결제 결과 반영 오류(gift)"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}
			}
			pushAcDate = unpaidInfo[0]["AC_DATE"]

		}

		txErr := err
		defer func() {
			if txErr != nil {
				// transaction rollback
				dprintf(4, c, "do rollback - 결제 결과 반영 실패 메세지 TpayUnpaidResult)  \n", txErr.Error())
				payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
				tx.Rollback()
			}
		}()

		// 결제정보 DB 저장1 : DAR_PAYMENT_REPORT
		insertPaymentReportQuery, err := cls.GetQueryJson(paymentsql.InsertTpayPayment, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentReportQuery parameter fail"))
		}
		_, err = tx.Exec(insertPaymentReportQuery)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentReportQuery, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류(report)"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
		}

		if sndSearchTy == "2" {

			//결제 내역 히스토리 테이블 INSERT
			insertUnpaidPaymentHistory, err := cls.GetQueryJson(paymentsql.InsertTpayUnpaidPaymentHistory, pgParam)
			if err != nil {
				txErr = err
				return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentHistory parameter fail"))
			}
			_, err = tx.Exec(insertUnpaidPaymentHistory)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", insertUnpaidPaymentHistory, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류(hist)"
				return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
			}

			//ORDER 테이블 업데이트
			UnpaidPaymentOkQuery, err := cls.GetQueryJson(paymentsql.UnpaidPaymentOk, pgParam)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
			}
			_, err = tx.Exec(UnpaidPaymentOkQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UnpaidPaymentOkQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
			}

			if wincubeUpdate == true {

				wincubeParam := make(map[string]string)
				wincubeParam["restId"] = "S0000000649"

				restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestCombine, wincubeParam, c)
				if err != nil {
					txErr = err
					m["resultCode"] = "99"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}
				chargeAmt := restTypeInfo[0]["CHARGE_AMT"]
				ichargeAmt, _ := strconv.Atoi(chargeAmt)
				wincubeParam["chargeAmt"] = strconv.Itoa(ichargeAmt + wincubeCargeAmt)

				UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombinePaidChargeAmt, wincubeParam)
				if err != nil {
					txErr = err
					return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
				}
				_, err = tx.Exec(UpdateCombineChargeAmtQuery)
				if err != nil {
					txErr = err
					dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
					m["resultCode"] = "98"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}

			}

			if beneficonUpdate == true {

				beneficonParam := make(map[string]string)
				beneficonParam["restId"] = "S0000000615"

				restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestCombine, beneficonParam, c)
				if err != nil {
					txErr = err
					m["resultCode"] = "99"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}
				chargeAmt := restTypeInfo[0]["CHARGE_AMT"]
				ichargeAmt, _ := strconv.Atoi(chargeAmt)
				beneficonParam["chargeAmt"] = strconv.Itoa(ichargeAmt + beneficonCargeAmt)

				UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombinePaidChargeAmt, beneficonParam)
				if err != nil {
					txErr = err
					return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
				}
				_, err = tx.Exec(UpdateCombineChargeAmtQuery)
				if err != nil {
					txErr = err
					dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
					m["resultCode"] = "98"
					m["resultMsg"] = "결제 결과 반영 오류"
					return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
				}

			}

		} else {
			payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
			tx.Rollback()

			dprintf(1, c, "결제 방식 오류로 인하여 취소하였습니다.\n", err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 방식 오류"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
		}

		err = tx.Commit()
		if err != nil {
			txErr = err
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
		}

	} else {

		m["resultCode"] = "99"
		m["resultMsg"] = resultMsg
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	sendResultConfirm(params["tid"], "000")

	params["sndGrpId"] = sndGrpId
	payInfo, err := cls.GetSelectDataRequire(paymentsql.SelectPayInfo, params, c)
	if err != nil {
		m["resultCode"] = "98"
		return c.Render(http.StatusOK, "tpays/resultFailWeb.htm", m)
	}

	m["goodsName"] = payInfo[0]["GOODSNAME"]
	m["userNm"] = payInfo[0]["USER_NM"]

	// 푸쉬 전송 시작

	if len(os.Getenv("SERVER_TYPE")) > 0 {
		pushMsg := ""

		strArray2 := strings.Split(sndRestId, "@")

		for _, str2 := range strArray2 {

			pushMsg = payInfo[0]["GRP_NM"] + "의 " + payInfo[0]["USER_NM"] + "님이 " + pushAcDate + "까지 사용한 금액을 정산(결제)하였습니다."
			apiPush.SendPush_Msg_V1("정산", pushMsg, "M", "1", str2, "", "bookpay")
			lprintf(4, "정산 push 메세지  :", pushMsg)

		}
	}

	// 푸쉬 전송 완료

	page := "tpays/payResultWebUnpaid.htm"
	m["grpId"] = sndGrpId
	m["restId"] = sndRestId
	m["userId"] = params["mallUserId"]
	m["searchTy"] = sndSearchTy

	return c.Render(http.StatusOK, page, m)
}

func TpayUnpaidReady(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	mid, _ := cls.GetTokenValue("TPAY.TPAY_MID", fname)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)
	payLocalUrl, _ := cls.GetTokenValue("TPAY.PAY_LOCAL_URL", fname)
	returnUrl, _ := cls.GetTokenValue("TPAY.UNPAID_RETURN_URL", fname)
	cancelUrl, _ := cls.GetTokenValue("TPAY.CANCEL_URL", fname)
	domain, _ := cls.GetTokenValue("TPAY.DOMAIN", fname)

	paymentTy := params["paymentTy"]
	selectedDate := ""
	appPrefix := "mngAdmin://"
	selectedDate = params["selectedDate"]
	userAgent := c.Request().Header.Get("User-Agent")

	connType := "1"
	ediDate, vbankExpDate := getEdiDate()
	moid := getMoid(params["userId"])
	//pgCd := "01"
	//userTy := "0"

	var tpay Tpay
	tpay.MerchantKey = merchantKey
	tpay.EdiDate = ediDate

	key := fmt.Sprintf("%s%s", ediDate, merchantKey)
	tpay.EncKey = getMD5HashHandler(key)

	input := fmt.Sprintf("%s%s%s", params["amt"], mid, moid)
	encryptData := base64.StdEncoding.EncodeToString(AesHandlerEncrypt(input, tpay.EncKey, merchantKey))

	params["amt"] = "0"

	params["mallUserId"] = params["userId"]
	userInfo, err := cls.GetSelectData(paymentsql.SelectRestUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	goodsName := "일괄 정산"

	params["buyerEmail"] = userInfo[0]["EMAIL"]
	params["buyerTel"] = userInfo[0]["HP_NO"]
	params["userNm"] = userInfo[0]["USER_NM"]
	params["goodsName"] = goodsName
	params["moid"] = moid

	//mallReserved := moid + "," + params["amt"] + "," + params["restId"] + "," + params["grpId"] + "," + params["searchTy"] + "," + params["addAmt"] + "," + params["payChannel"] + "," + pgCd + "," + userTy + "," + selectedDate + "," + params["paymentTy"] + "," + params["osTy"]

	//println(moid)

	//후불 처리 키값 업데이트
	unpaidReadyData, err := cls.GetQueryJson(mngOrder.UpdateUnpaidReadyData, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(unpaidReadyData)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	unpaidList, err := cls.GetSelectTypeRequire(mngOrder.SelectTpayUnpaidList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if unpaidList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "정산 내역이 없습니다."))
	}

	cip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)

	params["encryptData"] = encryptData
	params["ediDate"] = ediDate
	params["vbankExpDate"] = vbankExpDate
	params["connType"] = connType
	params["mid"] = mid
	params["moid"] = moid
	params["merchantKey"] = merchantKey
	params["payLocalUrl"] = payLocalUrl
	params["returnUrl"] = returnUrl
	params["cancelUrl"] = cancelUrl
	//params["mallReserved"] = mallReserved
	params["paymentTy"] = paymentTy
	params["selectedDate"] = selectedDate
	params["appPrefix"] = appPrefix
	params["userAgent"] = userAgent
	params["clientIp"] = cip
	params["mallIp"] = ""
	params["domain"] = domain

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["tpay"] = params
	m["unpaidList"] = unpaidList

	return c.JSON(http.StatusOK, m)
}

// 결제 준비
func TpayReadyBilling(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	fname := cls.Cls_conf(os.Args)
	mid, _ := cls.GetTokenValue("TPAY.TPAY_BILLING_MID", fname)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_BILLING_MERCHANT_KEY", fname)

	params["buyerEmail"] = mid
	params["buyerTel"] = merchantKey

	m := make(map[string]interface{})
	m["tpay"] = params

	return c.Render(http.StatusOK, "tpays/billing/billingReadyWeb.htm", m)
}

// 결제 준비
func TpayGenBillikey(c echo.Context) error {

	dprintf(4, c, "call TpayGenBillikey\n")

	params := cls.GetParamJsonMap(c)
	mid := "darayo003m"                                                                                       //commons.TPAY_MID
	merchantKey := "KVMb4iMUI20FIQ+wVEbaHUljTY6sl62FISECLw1Hd56sl1JB/N7zPBMEUK+v2iNKsnOUQzz+wX7VsWC50RsVcA==" //commons.TPAY_MERCHANT_KEY

	card_num := params["card_num"]
	buyer_auth_num := params["buyer_auth_num"]
	card_exp := params["card_exp"]
	card_pwd := params["card_pwd"]
	//card_cvc:=params["card_cvc"]
	//card_code:=params["card_code"]

	uValue := url.Values{
		"api_key": {merchantKey},
		"mid":     {mid},
		//"card_code": {card_code},
		"card_num":       {card_num},
		"buyer_auth_num": {buyer_auth_num},
		"card_exp":       {card_exp},
		//	"card_cvc": {card_cvc},
		"card_pwd": {card_pwd},
	}

	resp, err := http.PostForm("https://webtx.tpay.co.kr/api/v1/gen_billkey", uValue)
	if err != nil {
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 발급 요청 오류"
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "전송 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	println(string(respBody))

	var result TpayGenBillingData
	err = json.Unmarshal(respBody, &result)
	if err != nil {

		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "전송 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	m := make(map[string]interface{})
	if result.Result_cd != "0000" {
		m["resultCode"] = "99"
		m["resultMsg"] = result.Result_msg
	} else {
		m["resultCode"] = "00"
		m["resultMsg"] = result.Result_msg
		m["Card_token"] = result.Card_token
		m["Card_name"] = result.Card_name
		m["Card_num"] = result.Card_num
		m["Card_exp"] = card_exp

	}

	return c.JSON(http.StatusOK, m)
}

// 빌링키 삭제
func TpayDelbillkey(c echo.Context) error {

	dprintf(4, c, "call TpayDelbillkey\n")

	params := cls.GetParamJsonMap(c)
	mid := "darayo003m"                                                                                       //commons.TPAY_MID
	merchantKey := "KVMb4iMUI20FIQ+wVEbaHUljTY6sl62FISECLw1Hd56sl1JB/N7zPBMEUK+v2iNKsnOUQzz+wX7VsWC50RsVcA==" //commons.TPAY_MERCHANT_KEY

	card_token := params["card_token"]

	uValue := url.Values{
		"api_key":    {merchantKey},
		"mid":        {mid},
		"card_token": {card_token},
	}

	resp, err := http.PostForm("https://webtx.tpay.co.kr/api/v1/del_billkey", uValue)
	if err != nil {
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 삭제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 삭제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	println(string(respBody))

	var result TpayGenBillingData
	err = json.Unmarshal(respBody, &result)
	if err != nil {

		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 삭제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	m := make(map[string]interface{})
	if result.Result_cd != "0000" {
		m["resultCode"] = "99"
		m["resultMsg"] = result.Result_msg
	} else {
		m["resultCode"] = "00"
		m["resultMsg"] = result.Result_msg
	}

	return c.JSON(http.StatusOK, m)
}

// 빌링키 결제
func TpayBillingPay(c echo.Context) error {

	dprintf(4, c, "call TpayBillingPay\n")

	params := cls.GetParamJsonMap(c)
	mid := "darayo003m"                                                                                       //commons.TPAY_MID
	merchantKey := "KVMb4iMUI20FIQ+wVEbaHUljTY6sl62FISECLw1Hd56sl1JB/N7zPBMEUK+v2iNKsnOUQzz+wX7VsWC50RsVcA==" //commons.TPAY_MERCHANT_KEY

	card_token := params["card_token"]
	goods_nm := params["goods_nm"]
	amt := params["amt"]
	moid := params["moid"]
	mall_user_id := params["mall_user_id"]
	buyer_name := params["buyer_name"]
	buyer_tel := params["buyer_tel"]
	buyer_email := params["buyer_email"]
	batch_div := "1"

	uValue := url.Values{
		"api_key":      {merchantKey},
		"mid":          {mid},
		"card_token":   {card_token},
		"goods_nm":     {goods_nm},
		"amt":          {amt},
		"moid":         {moid},
		"mall_user_id": {mall_user_id},
		"buyer_name":   {buyer_name},
		"buyer_tel":    {buyer_tel},
		"buyer_email":  {buyer_email},
		"batch_div":    {batch_div},
	}

	resp, err := http.PostForm("https://webtx.tpay.co.kr/api/v1/payments_token", uValue)
	if err != nil {
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	println(string(respBody))

	var result TpayGenBillingData
	err = json.Unmarshal(respBody, &result)
	if err != nil {

		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	m := make(map[string]interface{})
	if result.Result_cd != "0000" {
		m["resultCode"] = "99"
		m["resultMsg"] = result.Result_msg
	} else {
		m["resultCode"] = "00"
		m["resultMsg"] = result.Result_msg
		m["tid"] = result.Tid
		m["moid"] = result.Moid
		m["amt"] = result.Paid_amt

		sendResultConfirm(result.Tid, "3001")
	}

	return c.JSON(http.StatusOK, m)
}

// 빌링키 결제 취소
func TpayBillingPayCancel(c echo.Context) error {

	dprintf(4, c, "call TpayBillingPayCancel\n")

	params := cls.GetParamJsonMap(c)
	cancel_pw := controller.TPAY_SIMPLE_PAY_CANCEL_PWD
	mid := controller.TPAY_MID_SIMPLE_PAY
	merchantKey := controller.TPAY_SIMPLE_PAY_MERCHANT_KEY
	cancelUrl := controller.TPAY_SIMPLE_PAY_CANCEL_URL

	cancel_amt := params["cancel_amt"]
	moid := params["moid"]
	cancel_msg := params["cancel_msg"]
	partial_cancel := "0"
	tid := params["tid"]

	uValue := url.Values{
		"api_key":        {merchantKey},
		"mid":            {mid},
		"moid":           {moid},
		"cancel_pw":      {cancel_pw},
		"cancel_amt":     {cancel_amt},
		"cancel_msg":     {cancel_msg},
		"partial_cancel": {partial_cancel},
		"tid":            {tid},
	}

	resp, err := http.PostForm(cancelUrl, uValue)
	if err != nil {
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	println(string(respBody))

	var result TpayRecvData
	err = json.Unmarshal(respBody, &result)
	if err != nil {

		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "Tpay 빌링키 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	m := make(map[string]interface{})
	if result.Result_cd != "2001" {
		m["resultCode"] = "99"
		m["resultMsg"] = result.Result_msg
	} else {
		m["resultCode"] = "00"
		m["resultMsg"] = result.Result_msg
		m["tid"] = result.Tid
		m["CancelDate"] = result.CancelDate
		m["CancelTime"] = result.CancelTime

		sendCancelResultConfirm(merchantKey, mid, result.Tid, "000")
	}

	return c.JSON(http.StatusOK, m)
}

//빌링키 결제
func TpayBillingPayment(c echo.Context) error {

	dprintf(4, c, "call TpayBillingPayment\n")

	params := cls.GetParamJsonMap(c)
	mid := controller.TPAY_MID_SIMPLE_PAY
	merchantKey := controller.TPAY_SIMPLE_PAY_MERCHANT_KEY
	payUrl := controller.TPAY_SIMPLE_PAY_PAYMENT_URL

	restId := params["restId"]
	grpId := params["grpId"]
	goods_nm := params["goodsName"]
	amt := params["amt"]
	addAmt := params["addAmt"]
	sndPgCd := "01"

	paymentTy := params["paymentTy"]
	selectedDate := params["selectedDate"]
	payChannel := params["payChannel"]
	searchTy := params["searchTy"]
	userTy := params["userTy"]

	moid := getMoid(params["mallUserId"])
	mall_user_id := params["mallUserId"]

	batch_div := "1"

	params["userId"] = params["mallUserId"]

	cardInfo, err := cls.GetSelectDataRequire(paymentsql.SelectTpayBillingCardInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if cardInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "카드 정보가 없습니다."))
	}

	BILLING_PWD := cardInfo[0]["BILLING_PWD"]
	card_token := cardInfo[0]["CARD_TOKEN"]
	buyer_name := cardInfo[0]["EMAIL"]
	buyer_tel := cardInfo[0]["HP_NO"]
	buyer_email := cardInfo[0]["USER_NM"]

	if BILLING_PWD == "NONE" || BILLING_PWD == "" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "간편결제 비밀번호를 설정후 사용해주세요."))
	}

	if BILLING_PWD != params["billPwd"] {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "비밀번호를 확인해주세요."))
	}

	uValue := url.Values{
		"api_key":      {merchantKey},
		"mid":          {mid},
		"card_token":   {card_token},
		"goods_nm":     {goods_nm},
		"amt":          {amt},
		"moid":         {moid},
		"mall_user_id": {mall_user_id},
		"buyer_name":   {buyer_name},
		"buyer_tel":    {buyer_tel},
		"buyer_email":  {buyer_email},
		"batch_div":    {batch_div},
	}

	m := make(map[string]interface{})

	resp, err := http.PostForm(payUrl, uValue)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "간편 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "간편 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	println(string(respBody))

	var result TpayGenBillingData
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "간편 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	if result.Result_cd != "0000" {
		m["resultCode"] = "99"
		m["resultMsg"] = result.Result_msg
		return c.JSON(http.StatusOK, m)
	} else {

		m["resultCode"] = "00"
		m["resultMsg"] = result.Result_msg
		//m["tid"] = result.Tid
		//m["moid"] = result.Moid
		//m["amt"] = result.Paid_amt

		pushAcDate := ""
		paymethod := "BILLING"
		params["restId"] = restId
		params["grpId"] = grpId
		params["decMoid"] = moid
		params["tid"] = result.Tid

		//중복 체크
		dupChk, err := cls.GetSelectDataRequire(paymentsql.SelectTPayDupCheck, params, c)
		if err != nil {
			m["resultCode"] = "99"
			m["resultMsg"] = "간편 결제 결과 반영 오류"
			billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
			return c.JSON(http.StatusOK, m)
		}

		dupCnt, _ := strconv.Atoi(dupChk[0]["DUP_CNT"])
		if dupCnt > 0 {
			m["resultCode"] = "99"
			m["resultMsg"] = "중복결제-간편"
			return c.JSON(http.StatusOK, m)
		}

		restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestType, params, c)
		if err != nil {
			m["resultCode"] = "99"
			m["resultMsg"] = "간편 결제 결과 반영 오류"
			billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
			return c.JSON(http.StatusOK, m)
		}

		restType := restTypeInfo[0]["REST_TYPE"]
		chargeAmt := restTypeInfo[0]["CHARGE_AMT"]

		pgParam := make(map[string]string)
		pgParam["payInfo"] = "5"

		tx, err := cls.DBc.Begin()
		if err != nil {
			m["resultCode"] = "99"
			m["resultMsg"] = "간편 결제 결과 반영 오류"
			billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
			return c.JSON(http.StatusOK, m)
		}

		txErr := err

		defer func() {
			if txErr != nil {
				// transaction rollback
				dprintf(4, c, "do rollback -간편 결제 결과 반영 TpayResult)  \n")
				billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
				tx.Rollback()
			}
		}()

		pgParam["moid"] = moid
		pgParam["paymethod"] = paymethod
		pgParam["transtype"] = "0"
		pgParam["goodsname"] = params["goodsName"]
		pgParam["amt"] = amt
		pgParam["addAmt"] = addAmt
		pgParam["userip"] = "0"
		pgParam["tid"] = params["tid"]
		pgParam["state"] = "0000"
		pgParam["statecd"] = "0"
		pgParam["cardno"] = result.Card_num
		pgParam["authcode"] = result.App_no

		now := time.Now()
		authDate := now.Format("060102150405")
		pgParam["authdate"] = authDate

		pgParam["cardquota"] = params["cardQuota"]
		pgParam["fncd"] = result.App_co
		pgParam["fnname"] = result.App_co_name
		pgParam["resultcd"] = result.Result_cd
		pgParam["resultmsg"] = result.Result_msg
		pgParam["pgCd"] = sndPgCd

		pgParam["histId"] = moid
		pgParam["restId"] = restId
		pgParam["grpId"] = grpId
		pgParam["userId"] = params["mallUserId"]
		pgParam["creditAmt"] = amt
		pgParam["userTy"] = userTy
		pgParam["searchTy"] = searchTy
		pgParam["paymentTy"] = paymentTy
		pgParam["payChannel"] = payChannel
		pgParam["selectedDate"] = selectedDate

		// 결제정보 DB 저장1 : DAR_PAYMENT_REPORT
		insertPaymentReportQuery, err := cls.GetQueryJson(paymentsql.InsertTpayPayment, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentReportQuery parameter fail"))
		}
		_, err = tx.Exec(insertPaymentReportQuery)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentReportQuery, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "간편결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		insertPaymentHistory, err := cls.GetQueryJson(paymentsql.InsertTpayPaymentHistory, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentHistory parameter fail"))
		}
		_, err = tx.Exec(insertPaymentHistory)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentHistory, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "간편결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		if restType == "G" {

			isndAmt, _ := strconv.Atoi(amt)
			ichargeAmt, _ := strconv.Atoi(chargeAmt)
			pgParam["chargeAmt"] = strconv.Itoa(ichargeAmt + isndAmt)

			UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombineChargeAmt, pgParam)
			if err != nil {
				txErr = err
				return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
			}
			_, err = tx.Exec(UpdateCombineChargeAmtQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

		}

		if searchTy == "1" {
			//선불처리

			agrmInfo, err := cls.GetSelectDataRequire(paymentsql.SelectAgrmInfo, pgParam, c)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

			prepaidAmt := agrmInfo[0]["PREPAID_AMT"]
			prepaidPoint := agrmInfo[0]["PREPAID_POINT"]
			pgParam["agrmId"] = agrmInfo[0]["AGRM_ID"]

			intAmt, _ := strconv.Atoi(amt)
			intAddAmt, _ := strconv.Atoi(addAmt)
			intPrepaidAmt, _ := strconv.Atoi(prepaidAmt)
			pgParam["prepaidAmt"] = strconv.Itoa(intPrepaidAmt + intAmt + intAddAmt)

			intPrepaidPoint, _ := strconv.Atoi(prepaidPoint)
			pgParam["prepaidPoint"] = strconv.Itoa(intPrepaidPoint + intAddAmt)

			updateAgrmPrepaidAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateAgrmPrepaidAmt, pgParam)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				return c.JSON(http.StatusOK, m)
			}
			_, err = tx.Exec(updateAgrmPrepaidAmtQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", updateAgrmPrepaidAmtQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

			pgParam["jobTy"] = "0"
			pgParam["prepaidAmt"] = strconv.Itoa(intAmt + intAddAmt)
			InsertPrepaidQuery, err := cls.GetQueryJson(paymentsql.InsertPrepaid, pgParam)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}
			_, err = tx.Exec(InsertPrepaidQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPrepaidQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

		} else if searchTy == "2" {
			//후불 처리
			unpaidInfo, err := cls.GetSelectDataRequire(paymentsql.SelectUnpaidPaymentInfo, pgParam, c)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

			pushAcDate = unpaidInfo[0]["AC_DATE"]

			if amt != unpaidInfo[0]["CREDIT_AMT"] {

				billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
				tx.Rollback()
				dprintf(1, c, "정산금액과 결제 금액이 다릅니다.\n", err)
				m["resultCode"] = "98"
				m["resultMsg"] = "정산금액과 결제 금액이 다릅니다"
				return c.JSON(http.StatusOK, m)

			}

			UnpaidPaymentCreditNowQuery, err := cls.GetQueryJson(paymentsql.UnpaidPaymentCreditNow, pgParam)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UnpaidPaymentCreditNowQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}
			_, err = tx.Exec(UnpaidPaymentCreditNowQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UnpaidPaymentCreditNowQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

		} else {
			billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
			tx.Rollback()

			dprintf(1, c, "결제 방식 오류로 인하여 취소하였습니다.\n", err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 방식 오류"
			return c.JSON(http.StatusOK, m)
		}

		err = tx.Commit()
		if err != nil {
			txErr = err
			m["resultCode"] = "98"
			m["resultMsg"] = "간편결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		sendResultConfirm(result.Tid, "000")
		params["sndGrpId"] = grpId
		payInfo, err := cls.GetSelectDataRequire(paymentsql.SelectPayInfo, params, c)
		if err != nil {
			m["resultCode"] = "98"
			return c.JSON(http.StatusOK, m)
		}

		m["goodsName"] = payInfo[0]["GOODSNAME"]
		m["userNm"] = payInfo[0]["USER_NM"]
		m["amt"] = amt

		// 푸쉬 전송 시작

		if restType == "N" {
			pushMsg := ""
			if searchTy == "1" {
				pushMsg = payInfo[0]["GRP_NM"] + "의 " + payInfo[0]["USER_NM"] + "님이 " + amt + "원을 충전(결제)하였습니다."
				apiPush.SendPush_Msg_V1("충전", pushMsg, "M", "1", restId, "", "bookpay")

				lprintf(4, "충전 push 메세지  :", pushMsg)
			} else if searchTy == "2" {
				pushMsg = payInfo[0]["GRP_NM"] + "의 " + payInfo[0]["USER_NM"] + "님이 " + pushAcDate + "까지 사용한 금액을 정산(결제)하였습니다."
				apiPush.SendPush_Msg_V1("정산", pushMsg, "M", "1", restId, "", "bookpay")
				lprintf(4, "정산 push 메세지  :", pushMsg)
			}

		}

	}

	// 푸쉬 전송 완료

	return c.JSON(http.StatusOK, m)
}

func TpayReadyNew(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	mid, _ := cls.GetTokenValue("TPAY.TPAY_MID", fname)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)
	payActionUrl, _ := cls.GetTokenValue("TPAY.PAY_ACTION_URL", fname)
	payActionWEBUrl, _ := cls.GetTokenValue("TPAY.PAY_ACTION_WEB_URL", fname)
	payLocalUrl, _ := cls.GetTokenValue("TPAY.PAY_LOCAL_URL", fname)
	returnUrl, _ := cls.GetTokenValue("TPAY.RETURN_URL_NEW", fname)
	cancelUrl, _ := cls.GetTokenValue("TPAY.CANCEL_URL", fname)

	paymentTy := params["paymentTy"]

	if params["searchTy"] == "1" {
		paymentTy = "0"
		params["paymentTy"] = "0"
	} else {
		paymentTy = "3"
		params["paymentTy"] = "3"
	}

	selectedDate := ""
	appPrefix := "mngAdmin://"
	selectedDate = params["selectedDate"]
	userAgent := c.Request().Header.Get("User-Agent")

	//if paymentTy == "2" || paymentTy == "3" {
	//	selectedDate = selectedDate + "235959"
	//}
	connType := "1"
	ediDate, vbankExpDate := getEdiDate()
	moid := getMoid(params["mallUserId"])
	pgCd := "01"
	userTy := "0"

	var tpay Tpay
	tpay.MerchantKey = merchantKey
	tpay.EdiDate = ediDate

	key := fmt.Sprintf("%s%s", ediDate, merchantKey)
	tpay.EncKey = getMD5HashHandler(key)

	input := fmt.Sprintf("%s%s%s", params["amt"], mid, moid)
	encryptData := base64.StdEncoding.EncodeToString(AesHandlerEncrypt(input, tpay.EncKey, merchantKey))

	userInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	goodsName := ""
	if params["searchTy"] == "1" {
		goodsName = userInfo[0]["REST_NM"] + " 충전"
	} else {
		goodsName = userInfo[0]["REST_NM"] + " 정산"
	}

	params["buyerEmail"] = userInfo[0]["EMAIL"]
	params["buyerTel"] = userInfo[0]["HP_NO"]
	params["userNm"] = userInfo[0]["USER_NM"]
	params["goodsName"] = goodsName

	mallReserved := moid + "," + params["amt"] + "," + params["restId"] + "," + params["grpId"] + "," + params["searchTy"] + "," + params["addAmt"] + "," + params["payChannel"] + "," + pgCd + "," + userTy + "," + selectedDate + "," + params["paymentTy"] + "," + params["osTy"]

	cip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)

	params["selectedDate"] = selectedDate
	params["encryptData"] = encryptData
	params["ediDate"] = ediDate
	params["vbankExpDate"] = vbankExpDate
	params["connType"] = connType
	params["mid"] = mid
	params["moid"] = moid
	params["merchantKey"] = merchantKey
	params["payActionUrl"] = payActionUrl
	params["payLocalUrl"] = payLocalUrl
	params["returnUrl"] = returnUrl
	params["cancelUrl"] = cancelUrl
	params["mallReserved"] = mallReserved
	params["paymentTy"] = paymentTy
	params["selectedDate"] = selectedDate
	params["appPrefix"] = appPrefix
	params["userAgent"] = userAgent
	params["clientIp"] = cip
	params["mallIp"] = ""

	page := "tpays/tpayReadyNew.htm"

	if params["osTy"] == "web" {
		params["payActionWEBUrl"] = payActionWEBUrl
		page = "tpays/tpayReadyWeb.htm"
	}

	m := make(map[string]interface{})
	m["tpay"] = params

	return c.Render(http.StatusOK, page, m)
}

func TpayResultNew(c echo.Context) error {

	dprintf(4, c, "call TpayResultNew\n")

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)

	mallReserved := params["mallReserved"]
	ediDate := params["ediDate"]

	decAmt := AesHandlerDecrypt(params["amt"], ediDate, merchantKey)
	decMoid := AesHandlerDecrypt(params["moid"], ediDate, merchantKey)

	arrMallReserved := strings.Split(mallReserved, ",")

	println(mallReserved)

	sndMoId := arrMallReserved[0]
	sndAmt := arrMallReserved[1]
	sndRestId := arrMallReserved[2]
	sndGrpId := arrMallReserved[3]
	sndSearchTy := arrMallReserved[4]
	sndAddamt := arrMallReserved[5]
	sndPayChannel := arrMallReserved[6]
	sndPgCd := arrMallReserved[7]
	sndUserTy := arrMallReserved[8]
	sndSelectedDate := arrMallReserved[9]
	sndPaymentTy := arrMallReserved[10]
	osTy := arrMallReserved[11]

	pushAcDate := ""

	m := make(map[string]interface{})
	m["osTy"] = osTy
	m["chargeAmt"] = sndAmt
	//return c.Render(http.StatusOK, "tpays/payResult.htm", m)

	if decMoid != sndMoId {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (moid 다름) :", decMoid, sndMoId)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.JSON(http.StatusOK, m)
	}

	if decAmt != sndAmt {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (금액 다름) :", decAmt, sndAmt)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.JSON(http.StatusOK, m)
	}

	if sndSearchTy == "2" {
		//sndPaymentTy="3"
	}

	paymethod := params["payMethod"]
	resultCd := params["resultCd"]
	resultMsg := params["resultMsg"]
	params["restId"] = sndRestId
	params["grpId"] = sndGrpId
	params["decMoid"] = decMoid

	//중복 체크
	dupChk, err := cls.GetSelectDataRequire(paymentsql.SelectTPayDupCheck, params, c)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "결제 결과 반영 오류"
		return c.JSON(http.StatusOK, m)
	}

	restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestType, params, c)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "결제 결과 반영 오류"
		return c.JSON(http.StatusOK, m)
	}

	restType := restTypeInfo[0]["REST_TYPE"]
	chargeAmt := restTypeInfo[0]["CHARGE_AMT"]

	dupCnt, _ := strconv.Atoi(dupChk[0]["DUP_CNT"])
	if dupCnt > 0 {
		m["resultCode"] = "99"
		m["resultMsg"] = "중복 결제 입니다"
		return c.JSON(http.StatusOK, m)
	}

	if (paymethod == "CARD" && resultCd == "3001") || (paymethod == "BANK" && resultCd == "4000") {

		pgParam := make(map[string]string)

		if paymethod == "CARD" {
			pgParam["payInfo"] = "1"
		} else {
			pgParam["payInfo"] = "0"
		}

		tx, err := cls.DBc.Begin()
		if err != nil {
			dprintf(4, c, " 결제 결과 반영 TpayResult 오류  \n")
			return c.JSON(http.StatusOK, m)
		}

		txErr := err

		defer func() {
			if txErr != nil {
				// transaction rollback
				dprintf(4, c, "do rollback - 결제 결과 반영 TpayResult)  \n")

				payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
				tx.Rollback()
			}
		}()

		pgParam["moid"] = sndMoId
		pgParam["paymethod"] = paymethod
		pgParam["transtype"] = "0"
		pgParam["goodsname"] = params["goodsName"]
		pgParam["amt"] = sndAmt
		pgParam["addAmt"] = sndAddamt
		pgParam["userip"] = "0"
		pgParam["tid"] = params["tid"]
		pgParam["state"] = "0000"
		pgParam["statecd"] = params["stateCd"]
		pgParam["cardno"] = params["cardNo"]
		pgParam["authcode"] = params["authCode"]
		pgParam["authdate"] = params["authDate"]
		pgParam["cardquota"] = params["cardQuota"]
		pgParam["fncd"] = params["fnCd"]
		pgParam["fnname"] = params["fnName"]
		pgParam["resultcd"] = params["resultCd"]
		pgParam["resultmsg"] = params["resultMsg"]
		pgParam["pgCd"] = sndPgCd

		pgParam["histId"] = sndMoId
		pgParam["restId"] = sndRestId
		pgParam["grpId"] = sndGrpId
		pgParam["userId"] = params["mallUserId"]
		pgParam["creditAmt"] = sndAmt
		pgParam["userTy"] = sndUserTy
		pgParam["searchTy"] = sndSearchTy
		pgParam["paymentTy"] = sndPaymentTy
		pgParam["payChannel"] = sndPayChannel
		pgParam["selectedDate"] = sndSelectedDate

		// 결제정보 DB 저장1 : DAR_PAYMENT_REPORT
		insertPaymentReportQuery, err := cls.GetQueryJson(paymentsql.InsertTpayPayment, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentReportQuery parameter fail"))
		}
		_, err = tx.Exec(insertPaymentReportQuery)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentReportQuery, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		insertPaymentHistory, err := cls.GetQueryJson(paymentsql.InsertTpayPaymentHistory, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentHistory parameter fail"))
		}
		_, err = tx.Exec(insertPaymentHistory)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentHistory, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		if restType == "G" {

			isndAmt, _ := strconv.Atoi(sndAmt)
			ichargeAmt, _ := strconv.Atoi(chargeAmt)
			pgParam["chargeAmt"] = strconv.Itoa(ichargeAmt + isndAmt)

			UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombineChargeAmt, pgParam)
			if err != nil {
				txErr = err
				return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
			}
			_, err = tx.Exec(UpdateCombineChargeAmtQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

		}

		if sndSearchTy == "1" {
			//선불처리

			agrmInfo, err := cls.GetSelectDataRequire(paymentsql.SelectAgrmInfo, pgParam, c)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

			prepaidAmt := agrmInfo[0]["PREPAID_AMT"]
			prepaidPoint := agrmInfo[0]["PREPAID_POINT"]
			pgParam["agrmId"] = agrmInfo[0]["AGRM_ID"]

			intAmt, _ := strconv.Atoi(sndAmt)
			intAddAmt, _ := strconv.Atoi(sndAddamt)
			intPrepaidAmt, _ := strconv.Atoi(prepaidAmt)
			pgParam["prepaidAmt"] = strconv.Itoa(intPrepaidAmt + intAmt + intAddAmt)

			intPrepaidPoint, _ := strconv.Atoi(prepaidPoint)
			pgParam["prepaidPoint"] = strconv.Itoa(intPrepaidPoint + intAddAmt)

			updateAgrmPrepaidAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateAgrmPrepaidAmt, pgParam)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				return c.JSON(http.StatusOK, m)
			}
			_, err = tx.Exec(updateAgrmPrepaidAmtQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", updateAgrmPrepaidAmtQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

			pgParam["jobTy"] = "0"
			pgParam["prepaidAmt"] = strconv.Itoa(intAmt + intAddAmt)
			InsertPrepaidQuery, err := cls.GetQueryJson(paymentsql.InsertPrepaid, pgParam)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}
			_, err = tx.Exec(InsertPrepaidQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPrepaidQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

		} else if sndSearchTy == "2" {
			//후불 처리
			unpaidInfo, err := cls.GetSelectDataRequire(paymentsql.SelectUnpaidPaymentInfo, pgParam, c)
			if err != nil {
				txErr = err
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

			pushAcDate = unpaidInfo[0]["AC_DATE"]

			if sndAmt != unpaidInfo[0]["CREDIT_AMT"] {

				payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
				tx.Rollback()
				dprintf(1, c, "정산금액과 결제 금액이 다릅니다.\n", err)
				m["resultCode"] = "98"
				m["resultMsg"] = "정산금액과 결제 금액이 다릅니다"
				return c.JSON(http.StatusOK, m)

			}

			UnpaidPaymentCreditNowQuery, err := cls.GetQueryJson(paymentsql.UnpaidPaymentCreditNow, pgParam)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UnpaidPaymentCreditNowQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}
			_, err = tx.Exec(UnpaidPaymentCreditNowQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UnpaidPaymentCreditNowQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

		} else {
			payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
			tx.Rollback()

			dprintf(1, c, "결제 방식 오류로 인하여 취소하였습니다.\n", err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 방식 오류"
			return c.JSON(http.StatusOK, m)
		}

		err = tx.Commit()
		if err != nil {
			txErr = err
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

	} else {
		m["resultCode"] = "99"
		m["resultMsg"] = resultMsg
		return c.JSON(http.StatusOK, m)
	}

	sendResultConfirm(params["tid"], "000")

	params["sndGrpId"] = sndGrpId
	payInfo, err := cls.GetSelectDataRequire(paymentsql.SelectPayInfo, params, c)
	if err != nil {
		m["resultCode"] = "98"
		return c.JSON(http.StatusOK, m)
	}

	m["goodsName"] = payInfo[0]["GOODSNAME"]
	m["userNm"] = payInfo[0]["USER_NM"]
	m["amt"] = sndAmt

	// 푸쉬 전송 시작

	if restType == "N" {
		pushMsg := ""
		if sndSearchTy == "1" {
			pushMsg = payInfo[0]["GRP_NM"] + "의 " + payInfo[0]["USER_NM"] + "님이 " + sndAmt + "원을 충전(결제)하였습니다."
			apiPush.SendPush_Msg_V1("충전", pushMsg, "M", "1", sndRestId, "", "bookpay")

			lprintf(4, "충전 push 메세지  :", pushMsg)
		} else if sndSearchTy == "2" {
			pushMsg = payInfo[0]["GRP_NM"] + "의 " + payInfo[0]["USER_NM"] + "님이 " + pushAcDate + "까지 사용한 금액을 정산(결제)하였습니다."
			apiPush.SendPush_Msg_V1("정산", pushMsg, "M", "1", sndRestId, "", "bookpay")
			lprintf(4, "정산 push 메세지  :", pushMsg)
		}

	}

	// 푸쉬 전송 완료

	return c.JSON(http.StatusOK, m)
}

func SimplePayNorder(c echo.Context) error {

	dprintf(4, c, "call SimplePayNorder\n")

	params := cls.GetParamJsonMap(c)
	mid := controller.TPAY_MID_SIMPLE_PAY
	merchantKey := controller.TPAY_SIMPLE_PAY_MERCHANT_KEY
	payUrl := controller.TPAY_SIMPLE_PAY_PAYMENT_URL

	restId := params["restId"]
	grpId := params["grpId"]

	amt := params["chargeAmt"]
	params["amt"] = amt
	params["addAmt"] = "0"
	sndPgCd := "01"
	paymentTy := "5" // 즉시결제
	payChannel := "02"
	searchTy := "1"
	userTy := "0"

	moid := getMoid(params["userId"])
	mall_user_id := params["userId"]

	batch_div := "1"

	params["mallUserId"] = params["userId"]

	cardInfo, err := cls.GetSelectDataRequire(paymentsql.SelectTpayBillingCardInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if cardInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "카드 정보가 없습니다."))
	}

	BILLING_PWD := cardInfo[0]["BILLING_PWD"]
	card_token := cardInfo[0]["CARD_TOKEN"]
	buyer_name := cardInfo[0]["EMAIL"]
	buyer_tel := cardInfo[0]["HP_NO"]
	buyer_email := cardInfo[0]["USER_NM"]
	goods_nm := cardInfo[0]["REST_NM"] + " 상품 구입"
	params["goodsName"] = goods_nm

	if BILLING_PWD == "NONE" || BILLING_PWD == "" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "간편결제 비밀번호를 설정후 사용해주세요."))
	}

	if BILLING_PWD != params["billPwd"] {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "비밀번호를 확인해주세요."))
	}

	uValue := url.Values{
		"api_key":      {merchantKey},
		"mid":          {mid},
		"card_token":   {card_token},
		"goods_nm":     {goods_nm},
		"amt":          {amt},
		"moid":         {moid},
		"mall_user_id": {mall_user_id},
		"buyer_name":   {buyer_name},
		"buyer_tel":    {buyer_tel},
		"buyer_email":  {buyer_email},
		"batch_div":    {batch_div},
	}

	m := make(map[string]interface{})

	resp, err := http.PostForm(payUrl, uValue)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "간편 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "간편 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	println(string(respBody))

	var result TpayGenBillingData
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "간편 결제 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	if result.Result_cd != "0000" {
		m["resultCode"] = "99"
		m["resultMsg"] = result.Result_msg
		return c.JSON(http.StatusOK, m)
	} else {
		paymethod := "BILLING"
		params["restId"] = restId
		params["grpId"] = grpId
		params["decMoid"] = moid
		params["tid"] = result.Tid

		//중복 체크
		dupChk, err := cls.GetSelectDataRequire(paymentsql.SelectTPayDupCheck, params, c)
		if err != nil {
			m["resultCode"] = "99"
			m["resultMsg"] = "간편 결제 결과 반영 오류"
			billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
			return c.JSON(http.StatusOK, m)
		}

		dupCnt, _ := strconv.Atoi(dupChk[0]["DUP_CNT"])
		if dupCnt > 0 {
			m["resultCode"] = "99"
			m["resultMsg"] = "중복결제-간편"
			return c.JSON(http.StatusOK, m)
		}

		restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestType, params, c)
		if err != nil {
			m["resultCode"] = "99"
			m["resultMsg"] = "간편 결제 결과 반영 오류"
			billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
			return c.JSON(http.StatusOK, m)
		}

		restType := restTypeInfo[0]["REST_TYPE"]
		chargeAmt := restTypeInfo[0]["CHARGE_AMT"]

		pgParam := make(map[string]string)
		pgParam["payInfo"] = "5"

		tx, err := cls.DBc.Begin()
		if err != nil {
			m["resultCode"] = "99"
			m["resultMsg"] = "간편 결제 결과 반영 오류"
			billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
			return c.JSON(http.StatusOK, m)
		}

		txErr := err

		defer func() {
			if txErr != nil {
				// transaction rollback
				dprintf(4, c, "do rollback -간편 결제 결과 반영 TpayResult)  \n")
				billingPayCancel(moid, amt, "간편 결제반영 실패로 인한 취소", params["tid"])
				tx.Rollback()
			}
		}()

		pgParam["moid"] = moid
		pgParam["paymethod"] = paymethod
		pgParam["transtype"] = "0"
		pgParam["goodsname"] = params["goodsName"]
		pgParam["amt"] = amt
		pgParam["addAmt"] = "0"
		pgParam["userip"] = "0"
		pgParam["tid"] = params["tid"]
		pgParam["state"] = "0000"
		pgParam["statecd"] = "0"
		pgParam["cardno"] = result.Card_num
		pgParam["authcode"] = result.App_no

		now := time.Now()
		authDate := now.Format("060102150405")
		pgParam["authdate"] = authDate

		pgParam["cardquota"] = params["cardQuota"]
		pgParam["fncd"] = result.App_co
		pgParam["fnname"] = result.App_co_name
		pgParam["resultcd"] = result.Result_cd
		pgParam["resultmsg"] = result.Result_msg
		pgParam["pgCd"] = sndPgCd

		pgParam["histId"] = moid
		pgParam["restId"] = restId
		pgParam["grpId"] = grpId
		pgParam["userId"] = params["mallUserId"]
		pgParam["creditAmt"] = amt
		pgParam["userTy"] = userTy
		pgParam["searchTy"] = searchTy
		pgParam["paymentTy"] = paymentTy
		pgParam["payChannel"] = payChannel
		pgParam["selectedDate"] = ""

		if params["unlink"] == "Y" {
			pgParam["grpId"] = "00000000000"
		}

		// 결제정보 DB 저장1 : DAR_PAYMENT_REPORT
		insertPaymentReportQuery, err := cls.GetQueryJson(paymentsql.InsertTpayPayment, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentReportQuery parameter fail"))
		}
		_, err = tx.Exec(insertPaymentReportQuery)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentReportQuery, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "간편결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		insertPaymentHistory, err := cls.GetQueryJson(paymentsql.InsertTpayPaymentHistory, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentHistory parameter fail"))
		}
		_, err = tx.Exec(insertPaymentHistory)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentHistory, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "간편결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		if restType == "G" {

			isndAmt, _ := strconv.Atoi(amt)
			ichargeAmt, _ := strconv.Atoi(chargeAmt)
			pgParam["chargeAmt"] = strconv.Itoa(ichargeAmt + isndAmt)

			UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombineChargeAmt, pgParam)
			if err != nil {
				txErr = err
				return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
			}
			_, err = tx.Exec(UpdateCombineChargeAmtQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "간편결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}
		}

		err = tx.Commit()
		if err != nil {
			txErr = err
			m["resultCode"] = "98"
			m["resultMsg"] = "간편결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		m := make(map[string]interface{})

		mocalUrl := controller.CONFIG_MOCA_URL

		P_URL := mocalUrl + "/api/moca/v2/order/" + restId + "/instantPayOrder"

		payload := url.Values{
			"payTy":       {params["payTy"]},
			"userId":      {params["userId"]},
			"grpId":       {params["grpId"]},
			"orderTy":     {params["orderTy"]},
			"orderAmt":    {params["orderAmt"]},
			"qrOrderTy":   {params["qrOrderTy"]},
			"itemNo":      {params["itemNo"]},
			"itemPrice":   {params["itemPrice"]},
			"itemCount":   {params["itemCount"]},
			"chargeAmt":   {params["chargeAmt"]},
			"remainAmt":   {params["remainAmt"]},
			"instantMoid": {moid},
			"unlink":      {params["unlink"]},
		}

		req, err := http.PostForm(P_URL, payload)
		if err != nil {
			billingPayCancel(moid, amt, "즉시 결제 주문 요청 실패로 인한 취소", params["tid"])
			m["resultCode"] = "99"
			m["resultMsg"] = "주문 요청 오류"
			return c.JSON(http.StatusOK, m)
		}
		defer req.Body.Close()
		// Response 체크.
		respOrder, err := ioutil.ReadAll(req.Body)
		if err != nil {
			billingPayCancel(moid, amt, "즉시 결제 주문 요청 실패로 인한 취소", params["tid"])
			m["resultCode"] = "99"
			m["resultMsg"] = "주문 요청 오류"
			return c.JSON(http.StatusOK, m)
		}

		println("주문결과", string(respOrder))

		var orderResult MocaOrderResult
		err = json.Unmarshal(respOrder, &orderResult)
		if err != nil {
			billingPayCancel(moid, amt, "즉시 결제 주문 요청 실패로 인한 취소", params["tid"])
			m := make(map[string]interface{})
			m["resultCode"] = "99"
			m["resultMsg"] = "주문 요청 오류"
			return c.JSON(http.StatusOK, m)
		}

		if orderResult.Result_cd != "00" {
			billingPayCancel(moid, amt, "즉시 결제 주문 요청 실패로 인한 취소", params["tid"])
			m := make(map[string]interface{})
			m["resultCode"] = "99"
			m["resultMsg"] = orderResult.Result_msg
			return c.JSON(http.StatusOK, m)
		}

		sendResultConfirm(result.Tid, "000")

		//	billingPayCancel(moid, amt, "테스트 취소", params["tid"])

		m["resultCode"] = "00"
		m["resultMsg"] = "정상"
		m["resultData"] = orderResult.ResultData

		return c.JSON(http.StatusOK, m)
	}

}

func PayNorderReady(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	mid, _ := cls.GetTokenValue("TPAY.TPAY_MID", fname)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)
	payActionUrl, _ := cls.GetTokenValue("TPAY.PAY_ACTION_URL", fname)
	payActionWEBUrl, _ := cls.GetTokenValue("TPAY.PAY_ACTION_WEB_URL", fname)
	payLocalUrl, _ := cls.GetTokenValue("TPAY.PAY_LOCAL_URL", fname)
	cancelUrl, _ := cls.GetTokenValue("TPAY.CANCEL_URL", fname)

	returnUrl, _ := cls.GetTokenValue("TPAY.RETURN_URL_PAYNORDER", fname)

	paymentTy := "5"

	params["paymentTy"] = paymentTy
	params["searchTy"] = params["payTy"]

	params["mallUserId"] = params["userId"]
	params["addAmt"] = "0"
	params["payChannel"] = "02"

	params["amt"] = params["chargeAmt"]

	selectedDate := ""
	appPrefix := "mngAdmin://"
	selectedDate = ""
	userAgent := c.Request().Header.Get("User-Agent")

	//if paymentTy == "2" || paymentTy == "3" {
	//	selectedDate = selectedDate + "235959"
	//}
	connType := "1"
	ediDate, vbankExpDate := getEdiDate()
	moid := getMoid(params["mallUserId"])

	var tpay Tpay
	tpay.MerchantKey = merchantKey
	tpay.EdiDate = ediDate

	key := fmt.Sprintf("%s%s", ediDate, merchantKey)
	tpay.EncKey = getMD5HashHandler(key)

	input := fmt.Sprintf("%s%s%s", params["amt"], mid, moid)
	encryptData := base64.StdEncoding.EncodeToString(AesHandlerEncrypt(input, tpay.EncKey, merchantKey))

	userInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	goodsName := userInfo[0]["REST_NM"] + " 상품 구입"

	params["buyerEmail"] = userInfo[0]["EMAIL"]
	params["buyerTel"] = userInfo[0]["HP_NO"]
	params["userNm"] = userInfo[0]["USER_NM"]
	params["goodsName"] = goodsName
	mallReserved := moid +
		"," + params["amt"] +
		"," + params["restId"] +
		"," + params["grpId"] +
		"," + params["searchTy"] +
		"," + params["osTy"] +
		"," + params["chargeAmt"] +
		"," + params["remainAmt"] +
		"," + params["payTy"] +
		"," + params["userId"] +
		"," + params["orderTy"] +
		"," + params["orderAmt"] +
		"," + params["qrOrderTy"] +
		"," + params["itemNo"] +
		"," + params["itemPrice"] +
		"," + params["itemCount"]

	cip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)

	params["selectedDate"] = selectedDate
	params["encryptData"] = encryptData
	params["ediDate"] = ediDate
	params["vbankExpDate"] = vbankExpDate
	params["connType"] = connType
	params["mid"] = mid
	params["moid"] = moid
	params["merchantKey"] = merchantKey
	params["payActionUrl"] = payActionUrl
	params["payLocalUrl"] = payLocalUrl
	params["returnUrl"] = returnUrl
	params["cancelUrl"] = cancelUrl
	params["mallReserved"] = mallReserved
	params["paymentTy"] = paymentTy
	params["selectedDate"] = selectedDate
	params["appPrefix"] = appPrefix
	params["userAgent"] = userAgent
	params["clientIp"] = cip
	params["mallIp"] = ""

	page := "tpays/tpayReadyNew.htm"

	if params["osTy"] == "web" {
		params["payActionWEBUrl"] = payActionWEBUrl
		page = "tpays/tpayReadyWeb.htm"
	}

	m := make(map[string]interface{})
	m["tpay"] = params

	return c.Render(http.StatusOK, page, m)
}

func PayNorder(c echo.Context) error {

	dprintf(4, c, "call PayNorder\n")

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	merchantKey, _ := cls.GetTokenValue("TPAY.TPAY_MERCHANT_KEY", fname)

	mallReserved := params["mallReserved"]
	ediDate := params["ediDate"]

	decAmt := AesHandlerDecrypt(params["amt"], ediDate, merchantKey)
	decMoid := AesHandlerDecrypt(params["moid"], ediDate, merchantKey)

	arrMallReserved := strings.Split(mallReserved, ",")

	//println(mallReserved)

	sndMoId := arrMallReserved[0]
	sndAmt := arrMallReserved[1]
	sndRestId := arrMallReserved[2]
	sndGrpId := arrMallReserved[3]
	sndSearchTy := arrMallReserved[4]
	osTy := arrMallReserved[5]
	sndChargeAmt := arrMallReserved[6]
	sndRemainAmt := arrMallReserved[7]
	sndPayTy := arrMallReserved[8]
	sndUserId := arrMallReserved[9]
	sndOrderTy := arrMallReserved[10]
	sndOrderAmt := arrMallReserved[11]
	sndQrOrderTy := arrMallReserved[12]
	sndItemNo := arrMallReserved[13]
	sndItemPrice := arrMallReserved[14]
	sndItemCount := arrMallReserved[15]

	sndPayChannel := "02"
	sndAddamt := "0"
	sndPaymentTy := "5"
	sndPgCd := "01"
	sndUserTy := "0"
	sndSelectedDate := "0"

	m := make(map[string]interface{})
	m["osTy"] = osTy
	m["chargeAmt"] = sndAmt
	//return c.Render(http.StatusOK, "tpays/payResult.htm", m)

	if decMoid != sndMoId {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (moid 다름) :", decMoid, sndMoId)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.JSON(http.StatusOK, m)
	}

	if decAmt != sndAmt {

		// 취소 호출
		lprintf(4, "결제 취소 요청 (금액 다름) :", decAmt, sndAmt)
		payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])

		m["resultCode"] = "99"
		m["resultMsg"] = "결제반영 실패"
		return c.JSON(http.StatusOK, m)
	}

	if sndSearchTy == "2" {
		//sndPaymentTy="3"
	}

	paymethod := params["payMethod"]
	resultCd := params["resultCd"]
	resultMsg := params["resultMsg"]
	params["restId"] = sndRestId
	params["grpId"] = sndGrpId
	params["decMoid"] = decMoid

	//중복 체크
	dupChk, err := cls.GetSelectDataRequire(paymentsql.SelectTPayDupCheck, params, c)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "결제 결과 반영 오류"
		return c.JSON(http.StatusOK, m)
	}

	restTypeInfo, err := cls.GetSelectDataRequire(paymentsql.SelectRestType, params, c)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = "결제 결과 반영 오류"
		return c.JSON(http.StatusOK, m)
	}

	restType := restTypeInfo[0]["REST_TYPE"]
	chargeAmt := restTypeInfo[0]["CHARGE_AMT"]

	dupCnt, _ := strconv.Atoi(dupChk[0]["DUP_CNT"])
	if dupCnt > 0 {
		m["resultCode"] = "99"
		m["resultMsg"] = "중복 결제 입니다"
		return c.JSON(http.StatusOK, m)
	}

	if (paymethod == "CARD" && resultCd == "3001") || (paymethod == "BANK" && resultCd == "4000") {

		pgParam := make(map[string]string)

		if paymethod == "CARD" {
			pgParam["payInfo"] = "1"
		} else {
			pgParam["payInfo"] = "0"
		}

		tx, err := cls.DBc.Begin()
		if err != nil {
			dprintf(4, c, " 결제 결과 반영 TpayResult 오류  \n")
			return c.JSON(http.StatusOK, m)
		}

		txErr := err

		defer func() {
			if txErr != nil {
				// transaction rollback
				dprintf(4, c, "do rollback - 결제 결과 반영 TpayResult)  \n")

				payCancel(sndMoId, sndAmt, "결제반영 실패로 인한 취소", "0", params["tid"])
				tx.Rollback()
			}
		}()

		pgParam["moid"] = sndMoId
		pgParam["paymethod"] = paymethod
		pgParam["transtype"] = "0"
		pgParam["goodsname"] = params["goodsName"]
		pgParam["amt"] = sndAmt
		pgParam["addAmt"] = sndAddamt
		pgParam["userip"] = "0"
		pgParam["tid"] = params["tid"]
		pgParam["state"] = "0000"
		pgParam["statecd"] = params["stateCd"]
		pgParam["cardno"] = params["cardNo"]
		pgParam["authcode"] = params["authCode"]
		pgParam["authdate"] = params["authDate"]
		pgParam["cardquota"] = params["cardQuota"]
		pgParam["fncd"] = params["fnCd"]
		pgParam["fnname"] = params["fnName"]
		pgParam["resultcd"] = params["resultCd"]
		pgParam["resultmsg"] = params["resultMsg"]
		pgParam["pgCd"] = sndPgCd

		pgParam["histId"] = sndMoId
		pgParam["restId"] = sndRestId
		pgParam["grpId"] = sndGrpId
		pgParam["userId"] = params["mallUserId"]
		pgParam["creditAmt"] = sndAmt
		pgParam["userTy"] = sndUserTy
		pgParam["searchTy"] = sndSearchTy
		pgParam["paymentTy"] = sndPaymentTy
		pgParam["payChannel"] = sndPayChannel
		pgParam["selectedDate"] = sndSelectedDate

		// 결제정보 DB 저장1 : DAR_PAYMENT_REPORT
		insertPaymentReportQuery, err := cls.GetQueryJson(paymentsql.InsertTpayPayment, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentReportQuery parameter fail"))
		}
		_, err = tx.Exec(insertPaymentReportQuery)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentReportQuery, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		insertPaymentHistory, err := cls.GetQueryJson(paymentsql.InsertTpayPaymentHistory, pgParam)
		if err != nil {
			txErr = err
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "insertPaymentHistory parameter fail"))
		}
		_, err = tx.Exec(insertPaymentHistory)
		if err != nil {
			txErr = err
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertPaymentHistory, err)
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

		if restType == "G" {

			isndAmt, _ := strconv.Atoi(sndAmt)
			ichargeAmt, _ := strconv.Atoi(chargeAmt)
			pgParam["chargeAmt"] = strconv.Itoa(ichargeAmt + isndAmt)

			UpdateCombineChargeAmtQuery, err := cls.GetQueryJson(paymentsql.UpdateCombineChargeAmt, pgParam)
			if err != nil {
				txErr = err
				return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateCombineChargeAmtQuery parameter fail"))
			}
			_, err = tx.Exec(UpdateCombineChargeAmtQuery)
			if err != nil {
				txErr = err
				dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateCombineChargeAmtQuery, err)
				m["resultCode"] = "98"
				m["resultMsg"] = "결제 결과 반영 오류"
				return c.JSON(http.StatusOK, m)
			}

		}

		err = tx.Commit()
		if err != nil {
			txErr = err
			m["resultCode"] = "98"
			m["resultMsg"] = "결제 결과 반영 오류"
			return c.JSON(http.StatusOK, m)
		}

	} else {
		m["resultCode"] = "99"
		m["resultMsg"] = resultMsg
		return c.JSON(http.StatusOK, m)
	}

	restId := params["restId"]
	mocalUrl := controller.CONFIG_MOCA_URL

	P_URL := mocalUrl + "/api/moca/v2/order/" + restId + "/instantPayOrder"

	payload := url.Values{
		"payTy":       {sndPayTy},
		"userId":      {sndUserId},
		"grpId":       {sndGrpId},
		"orderTy":     {sndOrderTy},
		"orderAmt":    {sndOrderAmt},
		"qrOrderTy":   {sndQrOrderTy},
		"itemNo":      {sndItemNo},
		"itemPrice":   {sndItemPrice},
		"itemCount":   {sndItemCount},
		"chargeAmt":   {sndChargeAmt},
		"remainAmt":   {sndRemainAmt},
		"instantMoid": {sndMoId},
		"unlink":      {params["unlink"]},
	}

	req, err := http.PostForm(P_URL, payload)
	if err != nil {
		payCancel(sndMoId, sndAmt, "즉시 결제 주문 요청 실패로 인한 취소", "0", params["tid"])
		m["resultCode"] = "99"
		m["resultMsg"] = "주문 요청 오류"
		return c.JSON(http.StatusOK, m)
	}
	defer req.Body.Close()
	// Response 체크.
	respOrder, err := ioutil.ReadAll(req.Body)
	if err != nil {
		payCancel(sndMoId, sndAmt, "즉시 결제 주문 요청 실패로 인한 취소", "0", params["tid"])
		m["resultCode"] = "99"
		m["resultMsg"] = "주문 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	println("주문결과", string(respOrder))

	var orderResult MocaOrderResult
	err = json.Unmarshal(respOrder, &orderResult)
	if err != nil {
		payCancel(sndMoId, sndAmt, "즉시 결제 주문 요청 실패로 인한 취소", "0", params["tid"])
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "주문 요청 오류"
		return c.JSON(http.StatusOK, m)
	}

	if orderResult.Result_cd != "00" {
		payCancel(sndMoId, sndAmt, "즉시 결제 주문 요청 실패로 인한 취소", "0", params["tid"])
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = orderResult.Result_msg
		return c.JSON(http.StatusOK, m)
	}

	sendResultConfirm(params["tid"], "000")

	m["resultCode"] = "00"
	m["resultMsg"] = "정상"
	m["resultData"] = orderResult.ResultData

	return c.JSON(http.StatusOK, m)
}
