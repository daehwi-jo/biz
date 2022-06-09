package controller

import (
	paymentsql "biz-web/query/payment"
	"biz-web/src/controller/cls"
	"bytes"
	"encoding/json"
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"github.com/golang-module/carbon"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	querysql "biz-web/query"
	tasksql "biz-web/query"
)

type REPORT uint

const (
	YReport REPORT = iota // 일일 리포트
	WReport               // 주간 리포트
	MReport               // 월간 리포트
)

const (
	//템플릿 코드					//템플릿 명

	CASH_016 = "cash_016" // 어제분석_주말실패
	CASH_015 = "cash_015" // 어제분석_2차실패
	CASH_014 = "cash_014" // 어제분석_1차실패
	CASH_013 = "cash_013" // 어제분석_성공
	CASH_019 = "cash_019" // 어제분석_성공4
	CASH_005 = "cash_005" // 여신인증
)

type mocaResp struct {
	ResultCode string `json:"resultCode"`
	ResultData struct {
		AvailableBalance string `json:"Available_balance"`
		ResultCd         string `json:"RESULT_CD"`
	} `json:"resultData"`
	ResultMsg string `json:"resultMsg"`
}

var HttpType, Host, Port string
var KakaoFlag bool
var RedisAddr []string

func DailyConfig(fname string) {
	if len(os.Getenv("SERVER_TYPE")) > 0 {
		HttpType = "HTTPS"
		Host = "cashapi.darayo.com"
		Port = "7788"
	}

	v, r := cls.GetTokenValue("KAKAOFLAG", fname) // value & return
	if r != cls.CONF_ERR {
		if v == "y" || v == "Y" {
			KakaoFlag = true
		}

		if v == "yes" {
			KakaoFlag = true
			HttpType = "HTTPS"
			Host = "cashapi.darayo.com"
			Port = "7788"
		}
	}

	v, r = cls.GetTokenValue("REDIS_INFO", fname)
	if r != cls.CONF_ERR {
		rCnt, err := strconv.Atoi(v)
		if err == nil && rCnt > 0 {
			for i := 0; i < rCnt; i++ {
				v, r = cls.GetTokenValue(fmt.Sprintf("REDIS_INF0%d", i), fname)
				if r != cls.CONF_ERR {
					rConfig := strings.Split(v, ",")
					if len(rConfig) == 2 {
						RedisAddr = append(RedisAddr, fmt.Sprintf("%s:%s", rConfig[1], rConfig[0]))
					}
				}
			}
		}
	}
}

// 지원금 초기화
// 매월 1일 02시 02 분 실행
func MonthResetSupportBalance() {

	if len(os.Getenv("SERVER_TYPE")) > 0 {
		lprintf(4, "[Start] MonthResetSupportBalance \n")
	} else {
		lprintf(4, "[stop] MonthResetSupportBalance \n")
		return
	}

	seq := DailyTaksLogInsert("MonthResetSupportBalance", "N", "1M", "I", 0)

	params := make(map[string]string)

	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.SetUpdateParam(querysql.UpdateResetSupportBalance, params)
	if err != nil {
		lprintf(1, "[ERROR] MonthResetSupportBalance query error(%s) \n", err.Error())
		return
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		lprintf(1, "[ERROR] MonthResetSupportBalance execute error(%s) \n", err.Error())
		return
	}

	lprintf(4, "[Finish] MonthResetSupportBalance \n")
	DailyTaksLogInsert("MonthResetSupportBalance", "Y", "1M", "U", seq)

}

func addComma(v string) string {

	if len(v) == 0 || v == "0" {
		return "0"
	}

	m := strings.TrimSpace(v)
	m1 := strings.ReplaceAll(m, ",", "")
	m2 := strings.ReplaceAll(m1, ".", "")

	amt, err := strconv.ParseInt(m2, 10, 64)
	if err != nil {
		cls.Lprintf(1, "[ERROR] parse int 64 err(%s)\n", err.Error())
		return "0"
	}

	return fmt.Sprintf("%s", humanize.Comma(amt))
}

// 매일 오전 9:30분 스케줄 작업 결과 여부 전송
func SendDailyJobResult() {

	if len(os.Getenv("SERVER_TYPE")) == 0 {
		return
	}

	params := make(map[string]string)
	var cnt int

	// b_dailytask 결과
	dResult, err := cls.SelectData(querysql.SelectDailyResult, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectDailyResult error(%s) \n", err.Error())
		return
	}

	for _, v := range dResult {
		if v["task_succ_yn"] == "Y" {
			cnt++
		}
	}

	title := fmt.Sprintf("일일 작업 결과")
	roomNumber := fmt.Sprintf("655403")
	msg := fmt.Sprintf("[BIZ_WEB]일일 작업 성공 여부 (%d/%d)", cnt, len(dResult))

	sendChannel(title, msg, roomNumber)

	resp, err := cls.HttpsGet("mocaapi.darayo.com", "7777", "api/moca/v2/wincubes/checkBalance")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var m mocaResp

	err = json.Unmarshal(data, &m)
	if err != nil {
		return
	}

	cls.Lprintf(4, "[INFO] wincube(%s)\n", m.ResultData.AvailableBalance)

	// 윈큐브 잔여금액 결과
	title = fmt.Sprintf("달아요 잔여 금액")
	msg = fmt.Sprintf("윈큐브 : %s원", addComma(m.ResultData.AvailableBalance))
	roomNumber = fmt.Sprintf("655403")

	sendChannel(title, msg, roomNumber)
}

func sendChannel(title, msg, roomNumber string) {
	body := `{ "conversation_id": ${ROOM}, "text": "캐시컴바인 알림",	"blocks": [ { "type": "header",	"text": "${TITLE}", "style": "blue" },  { "type": "text", "text": "${MSG}", "markdown": true } ] }`
	body = strings.Replace(body, "${TITLE}", title, -1)
	body = strings.Replace(body, "${MSG}", msg, -1)
	body = strings.Replace(body, "${ROOM}", roomNumber, -1)

	urlStr := "https://api.kakaowork.com/v1/messages.send?Content-Type=application/json"
	lprintf(4, "[INFO][go] url str(%s) \n", urlStr)
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(body)))
	if err != nil {
		lprintf(1, "[ERROR] http NewRequest (%s) \n", err.Error())
		return
	}

	req.Header.Set("Authorization", "Bearer 177f6c7f.dfa16ed40fd1493782f308ac9d15ce25")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		lprintf(1, "[ERROR] do error: http (%s) \n", err)
		return
	}
	defer resp.Body.Close()

	return
}

// 지정된 일자까지 미사용된 선물 취소
// 매일 01시 50 분 실행
func UnusedGiftCancel() {

	if len(os.Getenv("SERVER_TYPE")) > 0 {
		lprintf(4, "[Start] UnusedGiftCancel \n")
	} else {
		lprintf(4, "[stop] UnusedGiftCancel \n")
		return
	}

	seq := DailyTaksLogInsert("UnusedGiftCancel", "N", "1D", "I", 0)

	params := make(map[string]string)
	params["cancelDur"] = "7"

	giftList, err := cls.SelectData(querysql.SelectUnusedGiftList, params)
	if err != nil {
		lprintf(1, "[ERROR] UnusedGiftCancel SelectUnusedGiftList error(%s) \n", err.Error())
		return
	}

	if giftList == nil {
		lprintf(4, "[Finish] UnusedGiftCancel No Data \n")
		DailyTaksLogInsert("UnusedGiftCancel", "Y", "1D", "U", seq)
		return
	}

	for i := range giftList {

		params["grpId"] = giftList[i]["SND_GRP_ID"]
		params["restId"] = giftList[i]["REST_ID"]
		params["userId"] = giftList[i]["SND_USER_ID"]
		params["giftAmt"] = giftList[i]["GIFT_AMT"]
		params["giftMoid"] = giftList[i]["MOID"]

		//println(giftList[i]["MOID"])

		GiftCancel(params)
	}

	lprintf(4, "[Finish] UnusedGiftCancel \n")

	DailyTaksLogInsert("UnusedGiftCancel", "Y", "1D", "U", seq)

}

// 매일 지원금 백업
func BackupSupportBalance() {

	if len(os.Getenv("SERVER_TYPE")) > 0 {
		lprintf(4, "[Start] BackupSupportBalance \n")
	} else {
		lprintf(4, "[stop] BackupSupportBalance \n")
		return
	}

	seq := DailyTaksLogInsert("BackupSupportBalance", "N", "1D", "I", 0)

	params := make(map[string]string)
	params["cancelDur"] = "7"

	insertDailyTaskQuery, err := cls.GetQueryJson(querysql.InsertBackupSupportBalance, params)
	if err != nil {
		lprintf(1, "[ERROR] dailyTaksLogInsert InsertBackupSupportBalance error(%s) \n", err.Error())
		return
	}
	// 쿼리 실행
	_, err = cls.QueryDB(insertDailyTaskQuery)
	if err != nil {
		lprintf(1, "[ERROR] dailyTaksLogInsert InsertBackupSupportBalance error(%s) \n", err.Error())
		return
	}

	lprintf(4, "[Finish] BackupSupportBalance \n")

	DailyTaksLogInsert("BackupSupportBalance", "Y", "1D", "U", seq)

}

// 선물 취소
func GiftCancel(params map[string]string) {

	giftInfo, err := cls.SelectData(querysql.SelectGiftInFo, params)
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel SelectGiftInFo error(%s) \n", err.Error())
		return
	}
	if giftInfo == nil {
		lprintf(4, "[ERROR] GiftCancel 선물정보 없음 error(%s) \n", err.Error())
		return
	}

	giftStsCd := giftInfo[0]["GIFT_STS_CD"]
	orderNo := giftInfo[0]["ORDER_NO"]
	rcvStsCd := giftInfo[0]["RCV_STS_CD"]

	if giftStsCd == "2" {
		lprintf(4, "[ERROR] GiftCancel 이미 취소된 선물 \n")
		return

	} else if giftStsCd == "3" {
		lprintf(4, "[ERROR] GiftCancel 회수된 선물 \n")
		return

	} else if rcvStsCd == "1" {
		lprintf(4, "[ERROR] GiftCancel 사용된 선물 \n")
		return

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
			lprintf(4, "[ERROR] do rollback -선물 취소 error(%s) \n", err.Error())
			tx.Rollback()
		}
	}()

	// 선물 정보 변경
	params["giftStsCd"] = "2"
	UpdateGiftCancelQuery, err := cls.GetQueryJson(querysql.UpdateGiftCancel, params)
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel UpdateGiftCancelQuery error(%s) \n", err.Error())
		return
	}
	// 쿼리 실행

	_, err = tx.Exec(UpdateGiftCancelQuery)
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel UpdateGiftCancelQuery error(%s) \n", err.Error())
		return
	}

	agrmInfo, err := cls.SelectData(querysql.SelectAgrmInFo, params)
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel SelectGiftInFo error(%s) \n", err.Error())
		return
	}
	if agrmInfo == nil {
		lprintf(4, "[ERROR] GiftCancel 협약정보 없음 error(%s) \n", err.Error())
		return
	}

	// 선불 잔액차감 계산

	prepaidAmt, _ := strconv.Atoi(agrmInfo[0]["PREPAID_AMT"])
	agrmId := agrmInfo[0]["AGRM_ID"]
	giftAmt, _ := strconv.Atoi(params["giftAmt"])

	params["prepaidAmt"] = strconv.Itoa(prepaidAmt + giftAmt)
	params["agrmId"] = agrmId
	UpdateAgrmQuery, err := cls.GetQueryJson(querysql.UpdateAgrm, params)
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel UpdateAgrmQuery error(%s) \n", err.Error())
		return
	}
	// 쿼리 실행
	_, err = tx.Exec(UpdateAgrmQuery)
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel UpdateAgrmQuery error(%s) \n", err.Error())
		return
	}

	/// 선물 주문 취소

	params["orderNo"] = orderNo
	params["orderStat"] = "21"
	UpdateOrderCanCelQuery, err := cls.GetQueryJson(querysql.UpdateOrderCanCel, params)
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel UpdateAgrmQuery error(%s) \n", err.Error())
		return
	}
	// 쿼리 실행
	_, err = tx.Exec(UpdateOrderCanCelQuery)
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel UpdateAgrmQuery error(%s) \n", err.Error())
		return
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		lprintf(1, "[ERROR] GiftCancel  error(%s) \n", err.Error())
	}

}

// Tpay 지급요청 - 대상 선별
// 매일 01시 01 분 실행
func TpayMakePayStep1() {

	if len(os.Getenv("SERVER_TYPE")) > 0 {
		lprintf(4, "[Start] TpayMakePayStep1 \n")
	} else {
		lprintf(4, "[stop] TpayMakePayStep1 \n")
		return
	}

	seq := DailyTaksLogInsert("TpayMakePayStep1", "N", "1D", "I", 0)

	params := make(map[string]string)
	params["dayBefore"] = "1"

	//  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "[ERROR] do rollback -Tpay 지급요청 - 대상 선별 error(%s) \n", err.Error())
			tx.Rollback()
		}
	}()

	// 파라메터 맵으로 쿼리 변환
	UpdateMakePaymentIdQuery, err := cls.GetQueryJson(querysql.UpdateMakePaymentId, params)
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay UpdateMakePaymentIdQuery error(%s) \n", err.Error())
		return
	}
	// 쿼리 실행
	_, err = tx.Exec(UpdateMakePaymentIdQuery)
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay UpdateMakePaymentIdQuery error(%s) \n", err.Error())
		return
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay UpdateStoreInfoQuery error(%s) \n", err.Error())
	}

	lprintf(4, "[Finish] TpayMakePay \n")
	DailyTaksLogInsert("TpayMakePayStep1", "Y", "1D", "U", seq)

	TpayMakePayStep2()
}

// Tpay 지급요청 - 수수료 처리 및 대상 확정
// 매일 01시 01 분 실행
func TpayMakePayStep2() {

	lprintf(4, "[Start] TpayMakePayStep2 \n")
	seq := DailyTaksLogInsert("TpayMakePayStep2", "N", "1D", "I", 0)

	params := make(map[string]string)
	params["dayBefore"] = "1"

	paymentDt, err := cls.SelectData(querysql.SelectDayBefore, params)
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay SelectDayBefore error(%s) \n", err.Error())
		return
	}
	params["paymentDt"] = paymentDt[0]["beforeDate"]

	settlmtDt, err := cls.SelectData(querysql.SelectBizDay, params)
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay SelectBizDay error(%s) \n", err.Error())
		return
	}
	if settlmtDt == nil {
		lprintf(1, "[ERROR] TpayMakePay SelectBizDay error(%s) \n", err.Error())
		return
	}
	params["settlmtDt"] = settlmtDt[0]["totalDate"]

	//  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "[ERROR] do rollback -Tpay 지급요청 error(%s) \n", err.Error())
			tx.Rollback()
		}
	}()

	makeList, err := cls.SelectData(querysql.SelectMakeList, params)
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay SelectMakeList error(%s) \n", err.Error())
		return
	}

	for i := range makeList {

		mparam := make(map[string]string)
		mparam["moid"] = makeList[i]["MOID"]
		mparam["histId"] = makeList[i]["HIST_ID"]

		paymethod := makeList[i]["PAYMETHOD"]
		mparam["payMethod"] = paymethod

		amt := makeList[i]["AMT"]
		mparam["amt"] = amt
		mparam["restId"] = makeList[i]["REST_ID"]
		mparam["paymentDt"] = makeList[i]["PAYMENT_DT"]

		//가맹점 수수료
		restFeesInfo, err := cls.SelectData(paymentsql.SelectRestFeesInfo, mparam)
		if err != nil {
			lprintf(1, "[ERROR] TpayMakePay SelectRestFeesInfo error(%s) \n", err.Error())
			return
		}

		//PG 수수료
		pgFeesInfo, err := cls.SelectData(paymentsql.SelectPgFeesInfo, mparam)
		if err != nil {
			lprintf(1, "[ERROR] TpayMakePay SelectPgFeesInfo error(%s) \n", err.Error())
			return
		}

		//가맹점 수수료
		newRestFeesInfo := NewRestFees(restFeesInfo[0]["REST_FEES"], amt)

		//PG사 수수료
		newPgFeesInfo := NewPgFees(paymethod, pgFeesInfo[0]["FEE"], amt, pgFeesInfo[0]["VAT_YN"])

		// 회사 수수료
		newFitFeesInfo := NewFitFees(newRestFeesInfo["RestFeesAmt"], newPgFeesInfo["PgFeesAmt"])

		// 지급 금액
		restPayAmt := RealAmt(amt, newRestFeesInfo["RestFeesAmt"])

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
			lprintf(1, "[ERROR] TpayMakePay UpdatePaymentFeesQuery error(%s) \n", err.Error())
			return
		}
		// 쿼리 실행
		_, err = tx.Exec(UpdatePaymentFeesQuery)
		if err != nil {
			lprintf(1, "[ERROR] TpayMakePay UpdatePaymentFeesQuery error(%s) \n", err.Error())
			return
		}

		UpdatePaymentHistFeesQuery, err := cls.GetQueryJson(querysql.UpdatePaymentHistFees, mparam)
		if err != nil {
			lprintf(1, "[ERROR] TpayMakePay UpdatePaymentHistFeesQuery error(%s) \n", err.Error())
			return
		}
		// 쿼리 실행
		_, err = tx.Exec(UpdatePaymentHistFeesQuery)
		if err != nil {
			lprintf(1, "[ERROR] TpayMakePay UpdatePaymentHistFeesQuery error(%s) \n", err.Error())
			return
		}
	}

	UpdateMakePaymentQuery, err := cls.GetQueryJson(querysql.UpdateMakePayment, params)
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay UpdateMakePaymentQuery error(%s) \n", err.Error())
		return
	}
	// 쿼리 실행
	_, err = tx.Exec(UpdateMakePaymentQuery)
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay UpdateMakePaymentQuery error(%s) \n", err.Error())
		return
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		lprintf(1, "[ERROR] TpayMakePay UpdateStoreInfoQuery error(%s) \n", err.Error())
	}

	lprintf(4, "[Finish] TpayMakePayStep2 \n")
	DailyTaksLogInsert("TpayMakePayStep2", "Y", "1D", "U", seq)

}

// 기업 기념일 체크
func CompayPriceDay(oFlag bool, cId, uId string) {

	params := make(map[string]string)
	var sQuery string

	if oFlag {

		cls.Lprintf(4, "[INFO] cId(%s), uId(%s)", cId, uId)

		params["cId"] = cId
		sQuery = tasksql.SelectCompanyDayDetailInfo
	} else {
		sQuery = tasksql.SelectCompanyDayInfo
	}

	// 기념일 기업
	comps, err := cls.SelectData(sQuery, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCompanyDayInfo err(%s) \n", err.Error())
		return
	}

	for _, comp := range comps {

		params["cId"] = comp["cId"]
		params["eCode"] = comp["eCode"]
		params["price"] = comp["price"]

		// EVENT_CODE	0(직원생일)
		if comp["eCode"] == "0" {

			// d-day : 0(d-day), 1(d-1), 7(d-7)
			sDate, err := strconv.Atoi(comp["sDate"])
			if err != nil {
				cls.Lprintf(1, "[ERROR] err(%s)", err.Error())
				continue
			}

			// mmDD 형태
			shipDate := time.Now().AddDate(0, 0, sDate).Format("0102")
			_, cM, cD := carbon.Now().AddDays(sDate).Lunar().Date()
			shipDateLunar := fmt.Sprintf("%02d%02d", cM, cD)

			params["shipDate"] = shipDate
			params["shipDateLunar"] = shipDateLunar

			if oFlag {
				params["uId"] = uId
				sQuery = tasksql.SelectCompanyBirthdayUserDetailInfo
			} else {
				sQuery = tasksql.SelectCompanyBirthdayUserInfo
			}

			// 기념일 생일 직원
			compUsers, err := cls.SelectData(sQuery, params)
			if err != nil {
				lprintf(1, "[ERROR] SelectCompanyBirthdayUserInfo err(%s) \n", err.Error())
				continue
			}

			for _, compUser := range compUsers {

				uId := compUser["uId"]

				if len(uId) == 0 {
					continue
				}

				params["uId"] = uId
				params["bId"] = compUser["bId"]
				params["uName"] = compUser["uName"]
				params["hp"] = compUser["hp"]
				params["thisYear"] = time.Now().Format("2006")

				// 기존 기념일 지원금 히스토리 검색
				compHis, err := cls.SelectData(tasksql.SelectCompanyDayHis, params)
				if err != nil {
					lprintf(1, "[ERROR] SelectCompanyDayHis err(%s) \n", err.Error())
					continue
				}

				if len(compHis) >= 1 {
					continue
				}

				tx, _ := cls.DBc.Begin()

				params["cPrice"] = comp["price"]

				// 잔여 지원금 업데이트
				uQuery, err := cls.GetQueryJson(tasksql.UpdateCompanyUserDetail, params)
				if err != nil {
					lprintf(1, "[ERROR] UpdateCompanyUserDetail error(%s) \n", err.Error())
					tx.Rollback()
					continue
				}
				_, err = tx.Exec(uQuery)
				if err != nil {
					lprintf(1, "[ERROR] UpdateCompanyUserDetail exec error(%s) \n", err.Error())
					tx.Rollback()
					continue
				}

				params["priceDate"] = time.Now().Format("20060102")
				params["sDate"] = comp["sDate"]

				// 잔여 지원금 히스토리 Insert
				iQuery, err := cls.GetQueryJson(tasksql.InsertCompanyDayHis, params)
				if err != nil {
					lprintf(1, "[ERROR] InsertCompanyDayHis error(%s) \n", err.Error())
					tx.Rollback()
					continue
				}
				_, err = tx.Exec(iQuery)
				if err != nil {
					lprintf(1, "[ERROR] InsertCompanyDayHis exec error(%s) \n", err.Error())
					tx.Rollback()
					continue
				}

				err = tx.Commit()
				if err != nil {
					tx.Rollback()
					continue
				}

				cls.Lprintf(4, "[INFO] update balance -- companyId(%s), eventCode(%s), userId(%s)", params["cId"], params["eCode"], params["uId"])

			}

		}

	}

}

func WincubeItemUpdate() {

	if len(os.Getenv("SERVER_TYPE")) > 0 {
		lprintf(4, "[Start] WincubeItemUpdate \n")
	} else {
		lprintf(4, "[stop] WincubeItemUpdate \n")
		return
	}

	lprintf(4, "[Start] WincubeItemUpdate \n")

	seq := DailyTaksLogInsert("WincubeItemUpdate", "N", "1D", "I", 0)

	P_URL := ""
	fname := cls.Cls_conf(os.Args)
	mocalUrl, _ := cls.GetTokenValue("MOCA_URL", fname)

	P_URL = mocalUrl + "/api/moca/v2/wincubes/itemList"

	payload := strings.NewReader("")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", P_URL, payload)
	if err != nil {
		//fmt.Println(err)
		lprintf(1, "[ERROR] error call WincubeItemUpdate: %s\n", err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	println(string(body))

	lprintf(4, "[Finish] WincubeItemUpdate \n")
	DailyTaksLogInsert("WincubeItemUpdate", "Y", "1D", "U", seq)

}

func DailyTaksLogInsert(taskName, taskSuccYn, cycle, execute string, seq int) int {

	params := make(map[string]string)
	params["taskName"] = taskName
	params["taskSuccYn"] = taskSuccYn
	params["cycle"] = cycle
	params["execute"] = execute

	rId := 0
	if execute == "I" {
		InsertDailyTaskQuery, err := cls.GetQueryJson(querysql.InsertDailyTask, params)
		if err != nil {
			lprintf(1, "[ERROR] dailyTaksLogInsert InsertDailyTaskQuery error(%s) \n", err.Error())
			return 0
		}
		// 쿼리 실행
		_, seqId, err := cls.ExecDB(InsertDailyTaskQuery)
		if err != nil {
			lprintf(1, "[ERROR] dailyTaksLogInsert InsertDailyTaskQuery error(%s) \n", err.Error())
			return 0
		}
		rId = int(seqId)

	} else if execute == "U" {
		params["seq"] = strconv.Itoa(seq)
		UpdateDailyTaskQuery, err := cls.GetQueryJson(querysql.UpdateDailyTask, params)
		if err != nil {
			lprintf(1, "[ERROR] dailyTaksLogInsert UpdateDailyTaskQuery error(%s) \n", err.Error())
			return 0
		}
		// 쿼리 실행
		_, err = cls.QueryDB(UpdateDailyTaskQuery)
		if err != nil {
			lprintf(1, "[ERROR] dailyTaksLogInsert UpdateDailyTaskQuery error(%s) \n", err.Error())
			return 0
		}

	}

	return rId

}

func NewRestFees(restFees string, chargeAmt string) map[string]string {

	amt, _ := strconv.Atoi(chargeAmt)

	doubleChargeAmt := float64(amt)
	floatRestFees, _ := strconv.ParseFloat(restFees, 8)
	totalFeesAmt := math.Round(doubleChargeAmt * floatRestFees)
	supplyAmt := math.Round(totalFeesAmt / 1.1)
	vatAmt := math.Round(totalFeesAmt / 1.1 * 0.1)

	newRestFeesMap := make(map[string]string)
	newRestFeesMap["RestFeesAmt"] = fmt.Sprintf("%v", totalFeesAmt)
	newRestFeesMap["RestSupplyAmt"] = fmt.Sprintf("%v", supplyAmt)
	newRestFeesMap["RestVatAmt"] = fmt.Sprintf("%v", vatAmt)

	return newRestFeesMap

}

func NewPgFees(paymethod string, pgFees string, chargeAmt string, vatYn string) map[string]string {

	amt, _ := strconv.Atoi(chargeAmt)

	doubleChargeAmt := float64(amt)

	PgFeesAmt := 0.00
	PgSupplyAmt := 0.00
	PgVatAmt := 0.00

	if vatYn == "Y" {
		floatPgFees, _ := strconv.ParseFloat(pgFees, 8)
		PgSupplyAmt = math.Round((doubleChargeAmt * floatPgFees) / 1.1)
		PgVatAmt = math.Round(PgSupplyAmt * 0.1)

		if paymethod == "BANK" && amt <= 11600 {
			PgFeesAmt = 200
		} else {
			floatPgFees, _ := strconv.ParseFloat(pgFees, 8)
			PgFeesAmt = math.Round(doubleChargeAmt * floatPgFees)
		}
	} else {
		floatPgFees, _ := strconv.ParseFloat(pgFees, 8)
		PgSupplyAmt = math.Round(doubleChargeAmt * floatPgFees)
		PgVatAmt = math.Round(PgSupplyAmt * 0.1)
		if paymethod == "BANK" && amt <= 11600 {
			PgFeesAmt = 200
		} else {
			PgFeesAmt = math.Round(PgSupplyAmt * 1.1)
		}
	}

	pgFeesMap := make(map[string]string)
	pgFeesMap["PgSupplyAmt"] = fmt.Sprintf("%v", PgSupplyAmt)
	pgFeesMap["PgFeesAmt"] = fmt.Sprintf("%v", PgFeesAmt)
	pgFeesMap["PgVatAmt"] = fmt.Sprintf("%v", PgVatAmt)

	return pgFeesMap

}

func NewFitFees(restFeesAmt string, pgFeesAmt string) map[string]string {

	restAmt, _ := strconv.ParseFloat(restFeesAmt, 8)
	pgAmt, _ := strconv.ParseFloat(pgFeesAmt, 8)

	FitFeesAmt := math.Round(restAmt - pgAmt)
	FitSupplyAmt := math.Round(FitFeesAmt / 1.1)
	FitVatAmt := math.Round(FitFeesAmt / 1.1 * 0.1)

	newFitFeesMap := make(map[string]string)
	newFitFeesMap["FitFeesAmt"] = fmt.Sprintf("%v", FitFeesAmt)
	newFitFeesMap["FitSupplyAmt"] = fmt.Sprintf("%v", FitSupplyAmt)
	newFitFeesMap["FitVatAmt"] = fmt.Sprintf("%v", FitVatAmt)

	return newFitFeesMap

}

func RealAmt(amt string, restFeesAmt string) string {

	floatAmt, _ := strconv.ParseFloat(amt, 8)
	floatRestAmt, _ := strconv.ParseFloat(restFeesAmt, 8)

	real_Amt := floatAmt - floatRestAmt

	return fmt.Sprintf("%v", real_Amt)
}

func YesterdayReport1() {

	if !KakaoFlag {
		return
	}

	params := make(map[string]string)
	t := time.Now()
	params["bsDt"] = t.AddDate(0, 0, -1).Format("20060102")
	params["bsDt2"] = t.AddDate(0, 0, -8).Format("20060102")

	// 11시 실패
	failComps, err := cls.SelectData(tasksql.SelectFailDailyAlim, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectSuccessAlim err(%s) \n", err.Error())
		return
	}

	var failUrl string

	for _, comp := range failComps {
		params["userId"] = comp["userId"]
		params["restId"] = comp["restId"]

		// AlimTalkYN
		if !AlimTalkYN(params, WReport) {
			continue
		}

		failUrl = fmt.Sprintf("api/etc/kakao/yesterdayReport?bizNum=%s&code=%s", comp["bizNum"], CASH_014)

		// 11시 실패
		resp, err := cls.HttpRequest(HttpType, "GET", Host, Port, failUrl, true)
		if err != nil {
			lprintf(1, "[ERROR] WeekReport err(%s)\n", err.Error())
			continue
		}
		resp.Body.Close()
	}

	// 11시 성공
	successComps, err := cls.SelectData(tasksql.SelectSuccessDailyAlim, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectSuccessAlim err(%s) \n", err.Error())
		return
	}

	var sucessUrl string
	params["today"] = t.Format("20060102")
	params["code"] = "cash_014"

	var haliFailComps []map[string]string

	for _, comp := range successComps {
		params["userId"] = comp["userId"]
		params["restId"] = comp["restId"]
		params["bizNum"] = comp["bizNum"]

		// 11시에 실패 리포트 보냈는지
		// 성공을 했어도 11시에 성공 리포트 안보내고 2시에 성공 리포트 전송
		sendComp, err := cls.SelectData(tasksql.SelectSuccessReportSendCheck, params)
		if err != nil {
			lprintf(1, "[ERROR] SelectSuccessReportSendCheck err(%s) \n", err.Error())
			continue
		}

		if len(sendComp) > 0 && sendComp[0]["result"] == "success" {
			continue
		}

		// 일요일인데 지난주 수집결과는 0이 아닌데 어제 수집결과 0이면 fail
		if t.Weekday() == time.Sunday {
			sync, err := cls.SelectData(tasksql.SelectSyncAmt, params)
			if err != nil {
				lprintf(1, "[ERROR] SelectSyncAmt err(%s) \n", err.Error())
				continue
			}

			if len(sync) == 2 && sync[0]["amt"] != "0" && sync[1]["amt"] == "0" {
				haliFailComps = append(haliFailComps, comp)
				continue
			}
		}

		// AlimTalkYN
		if !AlimTalkYN(params, WReport) {
			continue
		}

		sucessUrl = fmt.Sprintf("api/etc/kakao/yesterdayReport?bizNum=%s&code=%s", comp["bizNum"], CASH_013)

		// 11시 성공
		resp, err := cls.HttpRequest(HttpType, "GET", Host, Port, sucessUrl, true)
		if err != nil {
			lprintf(1, "[ERROR] WeekReport err(%s)\n", err.Error())
			continue
		}
		resp.Body.Close()
	}

	// 일요일인데 수집결과 0이면 fail
	for _, comp := range haliFailComps {
		params["userId"] = comp["userId"]
		params["restId"] = comp["restId"]

		// AlimTalkYN
		if !AlimTalkYN(params, WReport) {
			continue
		}

		failUrl = fmt.Sprintf("api/etc/kakao/yesterdayReport?bizNum=%s&code=%s", comp["bizNum"], CASH_016)

		// 11시 실패
		resp, err := cls.HttpRequest(HttpType, "GET", Host, Port, failUrl, true)
		if err != nil {
			lprintf(1, "[ERROR] WeekReport err(%s)\n", err.Error())
			continue
		}
		resp.Body.Close()
	}
}

func YesterdayReport2() {

	if !KakaoFlag {
		return
	}

	params := make(map[string]string)
	params["bsDt"] = time.Now().AddDate(0, 0, -1).Format("20060102")
	params["yesterDt"] = time.Now().AddDate(0, 0, -2).Format("20060102")

	// 14시에 실패 했는데, 실패 사유가 로그인 에러고
	// 전 전날에 성공 했으면 -> 여신인증 알림
	failCardComps, err := cls.SelectData(tasksql.SelectFailDailyAlimCard, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectSuccessAlim err(%s) \n", err.Error())
		return
	}

	var failUrl string

	for _, comp := range failCardComps {
		params["userId"] = comp["userId"]
		params["restId"] = comp["restId"]

		ySucessCardComp, err := cls.SelectData(tasksql.SelectSucessDailyAlimCard, params)
		if err != nil {
			lprintf(1, "[ERROR] SelectSucessDailyAlimCard err(%s) \n", err.Error())
			continue
		}

		if len(ySucessCardComp) == 0 {
			continue
		}

		// AlimTalkYN
		if !AlimTalkYN(params, WReport) {
			continue
		}

		failUrl = fmt.Sprintf("api/etc/kakao/yesterdayReport?bizNum=%s&code=%s", comp["bizNum"], CASH_005)

		// 14시 실패
		resp, err := cls.HttpRequest(HttpType, "GET", Host, Port, failUrl, true)
		if err != nil {
			lprintf(1, "[ERROR] WeekReport err(%s)\n", err.Error())
			continue
		}
		resp.Body.Close()
	}

	// 14시 실패
	failComps, err := cls.SelectData(tasksql.SelectFailDailyAlim, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectSuccessAlim err(%s) \n", err.Error())
		return
	}

	for _, comp := range failComps {
		params["userId"] = comp["userId"]
		params["restId"] = comp["restId"]

		// AlimTalkYN
		if !AlimTalkYN(params, WReport) {
			continue
		}

		failUrl = fmt.Sprintf("api/etc/kakao/yesterdayReport?bizNum=%s&code=%s", comp["bizNum"], CASH_015)

		// 14시 실패
		resp, err := cls.HttpRequest(HttpType, "GET", Host, Port, failUrl, true)
		if err != nil {
			lprintf(1, "[ERROR] WeekReport err(%s)\n", err.Error())
			continue
		}
		resp.Body.Close()
	}

	// 14시 성공
	successComps, err := cls.SelectData(tasksql.SelectSuccessDailyAlim, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectSuccessAlim err(%s) \n", err.Error())
		return
	}

	var sucessUrl string
	params["today"] = time.Now().Format("20060102")
	params["code"] = CASH_019

	for _, comp := range successComps {
		params["userId"] = comp["userId"]
		params["restId"] = comp["restId"]

		// 이미 성공 리포트 보냈는지
		sendComp, err := cls.SelectData(tasksql.SelectSuccessReportSendCheck, params)
		if err != nil {
			lprintf(1, "[ERROR] SelectSuccessReportSendCheck err(%s) \n", err.Error())
			continue
		}

		if len(sendComp) > 0 && sendComp[0]["result"] == "success" {
			continue
		}

		// AlimTalkYN
		if !AlimTalkYN(params, WReport) {
			continue
		}

		sucessUrl = fmt.Sprintf("api/etc/kakao/yesterdayReport?bizNum=%s&code=%s", comp["bizNum"], CASH_013)

		// 14시 성공
		resp, err := cls.HttpRequest(HttpType, "GET", Host, Port, sucessUrl, true)
		if err != nil {
			lprintf(1, "[ERROR] WeekReport err(%s)\n", err.Error())
			continue
		}
		resp.Body.Close()
	}
}

func WeekReport() {

	if !KakaoFlag {
		return
	}

	params := make(map[string]string)

	dt := time.Now().AddDate(0, 0, -7)
	startWeek := cls.GetFirstOfWeek(dt)
	endWeek := cls.GetEndOfWeek(dt)
	wStartDt := startWeek.Format("20060102")
	wEndDt := endWeek.Format("20060102")

	params["startDt"] = wStartDt
	params["endDt"] = wEndDt

	comps, err := cls.SelectData(tasksql.SelectSuccessWeekAlim, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectSuccessAlim err(%s) \n", err.Error())
		return
	}

	var weekUrl string

	for _, comp := range comps {
		params["userId"] = comp["userId"]
		params["restId"] = comp["restId"]

		// AlimTalkYN
		if !AlimTalkYN(params, WReport) {
			continue
		}

		weekUrl = fmt.Sprintf("api/etc/kakao/lastWeekReport?bizNum=%s", comp["bizNum"])

		// 주간 발송
		resp, err := cls.HttpRequest(HttpType, "GET", Host, Port, weekUrl, true)
		if err != nil {
			lprintf(1, "[ERROR] WeekReport err(%s)\n", err.Error())
			continue
		}
		resp.Body.Close()
	}
}

func MonthReport() {

	if !KakaoFlag {
		return
	}

	params := make(map[string]string)
	params["bsDt"] = time.Now().AddDate(0, -1, 0).Format("200601")

	comps, err := cls.SelectData(tasksql.SelectSuccessMonthAlim, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectSuccessAlim err(%s) \n", err.Error())
		return
	}

	var monthUrl string

	for _, comp := range comps {
		params["userId"] = comp["userId"]
		params["restId"] = comp["restId"]

		// AlimTalkYN
		if !AlimTalkYN(params, WReport) {
			continue
		}

		monthUrl = fmt.Sprintf("api/etc/kakao/lastMonthReport?bizNum=%s", comp["bizNum"])

		// 월간 발송
		resp, err := cls.HttpRequest(HttpType, "GET", Host, Port, monthUrl, true)
		if err != nil {
			lprintf(1, "[ERROR] WeekReport err(%s)\n", err.Error())
			continue
		}
		resp.Body.Close()
	}
}

// userId, restId
func AlimTalkYN(params map[string]string, report REPORT) bool {

	// 사용자 USE YN
	user, err := cls.SelectData(tasksql.SelectUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s) \n", err.Error())
		return false
	}

	if len(user) == 0 {
		lprintf(4, "[INFO] param userId(%s) user len 0\n", params["userId"])
		return false
	}

	if user[0]["use_yn"] != "Y" {
		lprintf(4, "[INFO] userId(%s) not use \n", params["userId"])
		return false
	}

	// billing 기간 검색
	billing, err := cls.SelectData(tasksql.SelectBillingInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s) \n", err.Error())
		return false
	}

	if len(billing) == 0 {
		lprintf(4, "[INFO] param userId(%s) billing len 0\n", params["userId"])
		return false
	}

	if !(billing[0]["startYN"] == "Y" && billing[0]["endYN"] == "Y") {
		lprintf(4, "[INFO] param userId(%s) billing finish\n", params["userId"])
		return false
	}

	// kakao 알림 여부
	kakao, err := cls.SelectData(tasksql.SelectKakaoAlimYN, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s) \n", err.Error())
		return false
	}

	if len(kakao) == 0 {
		lprintf(4, "[INFO] param restId(%s) kakao len 0\n", params["restId"])
		return false
	}

	switch report {
	case YReport:
		if kakao[0]["kakao_daily"] != "Y" {
			lprintf(4, "[INFO] param restId(%s) kakao daily not Y\n", params["restId"])
			return false
		}
	case MReport:
		if kakao[0]["kakao_month"] != "Y" {
			lprintf(4, "[INFO] param restId(%s) kakao month not Y\n", params["restId"])
			return false
		}
	case WReport:
		if kakao[0]["kakao_week"] != "Y" {
			lprintf(4, "[INFO] param restId(%s) kakao week not Y\n", params["restId"])
			return false
		}
	}

	return true
}

func redisConnect(addr string) redis.Conn {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		lprintf(1, "[ERROR] redis con err(%s) \n", err.Error())
		return nil
	}

	pong, err := redis.String(c.Do("PING"))
	if err != nil {
		lprintf(1, "[ERROR] redis ping pong err(%s) \n", err.Error())
		c.Close()
		return nil
	}

	lprintf(4, "[INFO] redis con ping pong(%s)\n", pong)
	return c
}

func RedisInsert() {
	for _, addr := range RedisAddr {
		c := redisConnect(addr)
		if c == nil {
			return
		}
		defer c.Close()

		comps, err := cls.SelectData(tasksql.SelectRedisComp, nil)
		if err != nil {
			lprintf(1, "[ERROR] SelectRedisComp err (%s)\n", err.Error())
			return
		}

		for _, comp := range comps {
			bizNum := comp["biz_num"]
			setRedisCalendar(c, bizNum)
		}
	}
}

func setRedisCalendar(c redis.Conn, bizNum string) {
	payCalender := GetPayCalendar(bizNum)
	if payCalender != nil {
		payKey := fmt.Sprintf("%spayCalender", bizNum)
		lprintf(4, "[INFO] payKey(%s), payCalender(%v)\n", payKey, payCalender)

		pData, _ := json.Marshal(payCalender)

		reply, err := c.Do("SET", payKey, pData)
		if err != nil {
			lprintf(1, "[ERROR] redis con do err(%s) \n", err.Error())
		} else {
			lprintf(4, "[INFO] redis do reply(%v) \n", reply)
		}
	}

	aprvCalender := getAprvCalendar(bizNum)
	if aprvCalender != nil {
		aprvKey := fmt.Sprintf("%saprvCalender", bizNum)
		lprintf(4, "[INFO] aprvKey(%s), aprvCalender(%v)\n", aprvKey, aprvCalender)

		aData, _ := json.Marshal(aprvCalender)

		reply, err := c.Do("SET", aprvKey, aData)
		if err != nil {
			lprintf(1, "[ERROR] redis con do err(%s) \n", err.Error())
		} else {
			lprintf(4, "[INFO] redis do reply(%v) \n", reply)
		}
	}
}

func GetPayCalendar(bizNum string) map[string]interface{} {

	params := make(map[string]string)
	params["bizNum"] = bizNum

	t := time.Now()
	params["startDt"] = fmt.Sprintf("%s01", t.AddDate(0, -8, 0).Format("200601"))
	params["endDt"] = t.AddDate(0, 0, -1).Format("20060102")

	paySumList, err := cls.SelectData(tasksql.SelectPayCalendarSumList, params)
	if err != nil {
		return nil
	}
	delete(params, "startDt")
	delete(params, "endDt")

	data := make(map[string]interface{})

	var summary []map[string]interface{}
	for _, sumData := range paySumList {
		tmp := make(map[string]interface{})
		expt, _ := strconv.Atoi(sumData["outpExptAmt"])
		realIn, _ := strconv.Atoi(sumData["realInAmt"])
		diff, _ := strconv.Atoi(sumData["diffAmt"])

		tmp["trMonth"] = sumData["trMonth"][4:]
		tmp["outpExptAmt"] = expt
		tmp["realInAmt"] = realIn
		tmp["diffAmt"] = diff
		tmp["diffColor"] = sumData["diffColor"]

		summary = append(summary, tmp)
	}
	data["summary"] = summary

	var monthList []map[string]interface{}
	for idx, sumData := range paySumList {
		if idx == len(paySumList)-2 {
			break
		}

		// 날짜변경을 위해 Time 값으로 변경
		timeTrDt, err := time.Parse("20060102", fmt.Sprintf("%s01", sumData["trMonth"]))
		timeFirst, timeLast := cls.GetFirstAndLastOfMonth(timeTrDt)
		firstDay := cls.GetFirstOfWeek(timeFirst).Format("20060102")
		lastDay := cls.GetEndOfWeek(timeLast).Format("20060102")

		params["startDt"] = firstDay
		params["endDt"] = lastDay
		payList, err := cls.SelectData(tasksql.SelectPayCalendarList, params)
		if err != nil {
			return nil
		}
		delete(params, "startDt")
		delete(params, "endDt")

		var dayList []map[string]interface{}
		for _, payData := range payList {
			tmp := make(map[string]interface{})
			row, _ := strconv.Atoi(payData["rNum"])
			expt, _ := strconv.Atoi(payData["outpExptAmt"])
			realIn, _ := strconv.Atoi(payData["realInAmt"])
			diff, _ := strconv.Atoi(payData["diffAmt"])

			tmp["rNum"] = row
			tmp["trDt"] = payData["trDt"]
			tmp["outpExptAmt"] = expt
			tmp["realInAmt"] = realIn
			tmp["diffAmt"] = diff
			tmp["diffColor"] = payData["diffColor"]
			tmp["dayColor"] = payData["dayColor"]

			dayList = append(dayList, tmp)
		}
		monthData := make(map[string]interface{})
		monthData["trMonth"] = sumData["trMonth"]
		monthData["dayList"] = dayList
		monthList = append(monthList, monthData)
	}
	data["monthList"] = monthList

	return data
}

func getAprvCalendar(bizNum string) map[string]interface{} {

	params := make(map[string]string)
	params["bizNum"] = bizNum

	t := time.Now()
	params["startDt"] = fmt.Sprintf("%s01", t.AddDate(0, -8, 0).Format("200601"))
	params["endDt"] = t.AddDate(0, 0, -1).Format("20060102")

	aprvSumList, err := cls.SelectData(tasksql.SelectAprvCalendarSumList, params)
	if err != nil {
		return nil
	}
	delete(params, "startDt")
	delete(params, "endDt")

	data := make(map[string]interface{})

	var summary []map[string]interface{}
	for _, sumData := range aprvSumList {
		tmp := make(map[string]interface{})
		aprv, _ := strconv.Atoi(sumData["aprvAmt"])
		cash, _ := strconv.Atoi(sumData["cashAmt"])
		pca, _ := strconv.Atoi(sumData["pcaAmt"])
		tot, _ := strconv.Atoi(sumData["totAmt"])

		tmp["trMonth"] = sumData["trMonth"][4:]
		tmp["aprvAmt"] = aprv
		tmp["cashAmt"] = cash
		tmp["pcaAmt"] = pca
		tmp["totAmt"] = tot

		summary = append(summary, tmp)
	}
	data["summary"] = summary

	var monthList []map[string]interface{}
	for idx, sumData := range aprvSumList {
		if idx == len(aprvSumList)-2 {
			break
		}

		// 날짜변경을 위해 Time 값으로 변경
		timeTrDt, err := time.Parse("20060102", fmt.Sprintf("%s01", sumData["trMonth"]))
		timeFirst, timeLast := cls.GetFirstAndLastOfMonth(timeTrDt)
		firstDay := cls.GetFirstOfWeek(timeFirst).Format("20060102")
		lastDay := cls.GetEndOfWeek(timeLast).Format("20060102")

		params["startDt"] = firstDay
		params["endDt"] = lastDay
		aprvList, err := cls.SelectData(tasksql.SelectAprvCalendarList, params)
		if err != nil {
			return nil
		}
		delete(params, "startDt")
		delete(params, "endDt")

		var dayList []map[string]interface{}
		for _, aprvData := range aprvList {
			tmp := make(map[string]interface{})
			row, _ := strconv.Atoi(aprvData["rNum"])
			aprv, _ := strconv.Atoi(aprvData["aprvAmt"])
			cash, _ := strconv.Atoi(aprvData["cashAmt"])
			pca, _ := strconv.Atoi(aprvData["pcaAmt"])
			tot, _ := strconv.Atoi(aprvData["totAmt"])

			tmp["rNum"] = row
			tmp["trDt"] = aprvData["trDt"]
			tmp["aprvAmt"] = aprv
			tmp["cashAmt"] = cash
			tmp["pcaAmt"] = pca
			tmp["totAmt"] = tot
			tmp["diffColor"] = aprvData["diffColor"]
			tmp["dayColor"] = aprvData["dayColor"]

			dayList = append(dayList, tmp)
		}
		monthData := make(map[string]interface{})
		monthData["trMonth"] = sumData["trMonth"]
		monthData["dayList"] = dayList
		monthList = append(monthList, monthData)
	}
	data["monthList"] = monthList

	return data
}

func UpdateRedis(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	bizNum := params["biz_num"]

	for _, addr := range RedisAddr {
		con := redisConnect(addr)
		if con != nil {
			continue
		}
		defer con.Close()

		setRedisCalendar(con, bizNum)
	}

	return c.JSON(http.StatusOK, "ok")
}
