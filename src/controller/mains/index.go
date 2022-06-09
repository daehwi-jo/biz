package mains

import (
	"biz-web/query"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func PaidList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/payment/paidList.htm", m)
}

func PaidIng(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/payment/paidIng.htm", m)
}

func PaidOk(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/payment/paidOk.htm", m)
}

func PaidMng(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/payment/paidMng.htm", m)
}

func Combine(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/payment/combine.htm", m)
}

func CombineDesc(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["restId"] = params["restId"]

	return c.Render(http.StatusOK, "mngAdmin/payment/combineDesc.htm", m)
}

func CombineDescWincube(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["restId"] = params["restId"]

	return c.Render(http.StatusOK, "mngAdmin/payment/combineDesc_wincube.htm", m)
}

func Account(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/payment/account.htm", m)
}

func Login(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "login.htm", m)
}

func Home(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "index.htm", m)
}

func Home2(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "index2.htm", m)
}

func SysHomeData(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	userData, err := cls.GetSelectDataRequire(query.AllUserCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	companyData, err := cls.GetSelectDataRequire(query.AllCompanyCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	storeData, err := cls.GetSelectDataRequire(query.AllStoreCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	currentMonthData, err := cls.GetSelectDataRequire(query.CurrentMonthAmtCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	previousMonthData, err := cls.GetSelectDataRequire(query.PreviousMonthAmtCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["userData"] = userData[0]
	result["companyData"] = companyData[0]
	result["storeData"] = storeData[0]
	result["currentMonthData"] = currentMonthData[0]
	result["previousMonthData"] = previousMonthData[0]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func SysHomeContentsData(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	today, err := cls.GetSelectDataRequire(query.TimeUseToday, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	yesterday, err := cls.GetSelectData(query.TimeUserYesterday, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//날짜 값 지정 하면 데이터 출력되지 않음
	storeAmtRank, err := cls.GetSelectData(query.StoreAmtRank, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	companyAmtRank, err := cls.GetSelectData(query.CompanyAmtRank, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["today"] = today
	result["yesterday"] = yesterday

	if storeAmtRank == nil {
		result["storeAmtRank"] = []string{}
	} else {
		result["storeAmtRank"] = storeAmtRank
	}

	if companyAmtRank == nil {
		result["companyAmtRank"] = []string{}
	} else {
		result["companyAmtRank"] = companyAmtRank
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func HomeData(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	grpCount, err := cls.GetSelectDataRequire(query.GrpCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	grpUserCount, err := cls.GetSelectDataRequire(query.GrpUserCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	joinStoreCount, err := cls.GetSelectDataRequire(query.GrpJoinStore, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	useAmt, err := cls.GetSelectDataRequire(query.UseAmt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["grpCount"] = grpCount[0]
	result["grpUserCount"] = grpUserCount[0]
	result["joinStoreCount"] = joinStoreCount[0]
	result["useAmt"] = useAmt[0]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func HomeContentsData(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	chartData, err := cls.GetSelectDataRequire(query.HomeChartData(params["chartKey"]), params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	storeRank, err := cls.GetSelectDataRequire(query.StoreUseRank, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["chartData"] = chartData

	if storeRank == nil {
		result["storeRank"] = []string{}
	} else {
		result["storeRank"] = storeRank
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func Search(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]
	m["saKeyword"] = params["saKeyword"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "search.htm", m)
}

func SysSetting(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["params"] = params

	check, err := cls.GetSelectDataRequire(query.SelectSysUserPwCheck, params, c)
	if err != nil {
		return c.Render(http.StatusOK, "index.htm", m)
	}

	if check[0]["state"] == "1" {
		return c.Render(http.StatusOK, "sysSetting.htm", m)
	}
	return c.Render(http.StatusOK, "index.htm", m)
}

func AdminMember(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/member/memberMain.htm", m)
}

func AdminMemberInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/member/memberInfo.htm", m)
}

func AdminBook(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/book/bookMain.htm", m)
}
func AdminBookInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/book/bookInfo.htm", m)
}

func AdminCompany(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/company/companyMain.htm", m)
}

func AdminCompanyInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/company/companyInfo.htm", m)
}

func AdminStore(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/store/storeMain.htm", m)
}
func AdminStoreInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/store/storeInfo.htm", m)
}

func MngManagerInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngUser/ManagerInfoMng.htm", m)
}

func MngAnniversary(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngUser/AnniversaryMng.htm", m)
}

func MngGrpList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngBook/GrpBookList.htm", m)
}
func MngGrpUser(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]
	m["grpId"] = params["searchGrpId"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngBook/GrpBookUserMng.htm", m)
}
func MngOrderList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngOrder/orderList.htm", m)
}
func MngCalculateMng(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngOrder/calculateMng.htm", m)
}

func MngPaymentList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngOrder/orderPaymentList.htm", m)
}

func MngGrpRestList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngOrder/orderStoreList.htm", m)
}

func MngPayCharge(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["grpId"] = params["grpId"]
	m["restId"] = params["restId"]

	return c.Render(http.StatusOK, "mngOrder/payCharge.htm", m)
}

func MngPayCalculate(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["grpId"] = params["grpId"]
	m["restId"] = params["restId"]
	m["searchTy"] = params["searchTy"]

	return c.Render(http.StatusOK, "mngOrder/payCalculate.htm", m)
}

func BaseUrl(c echo.Context) error {

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func LoginOk(c echo.Context) error { //로그인 완료시

	dprintf(4, c, "call LoginOk\n")
	ip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectDataRequire(query.SelectUserLoginCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		// 접속 로그
		params["logInOut"] = "I"
		params["ip"] = ip
		params["succYn"] = "N"
		params["type"] = "id"
		LoginAcceesLog(c, params)
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "아이디 또는 비밀번호가 잘못되었습니다."))
	}

	loginId := resultData[0]["USER_ID"]
	userNm := resultData[0]["USER_NM"]
	authorCd := resultData[0]["AUTHOR_CD"]
	connYn := resultData[0]["CONN_ALLOW_YN"]
	companyId := ""
	companyNm := ""
	userId := ""

	if connYn == "N" {
		params["logInOut"] = "I"
		params["ip"] = ip
		params["succYn"] = "N"
		params["type"] = "id"
		LoginAcceesLog(c, params)
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "접속이 거부된 아이디 입니다."))
	}

	//fmt.Println(authorCd)

	if authorCd != "SYS" {

		compData, err := cls.GetSelectDataRequire(query.SelectCompanyData, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		companyId = compData[0]["COMPANY_ID"]
		companyNm = compData[0]["COMPANY_NM"]
		userId = compData[0]["USER_ID"]

	}

	//토큰 발행

	c = cls.SetLoginJWT(c, loginId)

	// 접속 로그
	params["logInOut"] = "I"
	params["ip"] = ip
	params["succYn"] = "Y"
	params["type"] = "id"
	LoginAcceesLog(c, params)

	userData := make(map[string]interface{})
	userData["loginId"] = loginId
	userData["userNm"] = userNm
	userData["authorCd"] = authorCd
	userData["companyId"] = companyId
	userData["companyNm"] = companyNm
	userData["userId"] = userId

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)
}

func LoginAcceesLog(c echo.Context, params map[string]string) {

	dprintf(4, c, "call LoginAcceesLog\n")
	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.GetQueryJson(query.InsertLoginAccess, params)
	if err != nil {
		dprintf(4, c, "LoginAcceesLog query parameter fail\n")
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		dprintf(4, c, "LoginAcceesLog DB fail\n")
	}

}

func AdminPayment(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/payment/payment.htm", m)
}

func AdminPaymentInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/payment/paymentInfo.htm", m)
}

func AdminOrder(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/order/order.htm", m)
}
func AdminOrderInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["orderNo"] = params["orderNo"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/order/orderInfo.htm", m)
}

func AdminTask(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/task/task.htm", m)
}

func AdminLog(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "mngAdmin/task/log.htm", m)
}

func AdminOCRText(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/ocrReceipt/ocrText.htm", m)
}

func AdminOCRReceipt(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//if params["authorCd"] == "SYS" {
	//	m["SYSTEM"] = true
	//} else {
	//	m["SYSTEM"] = false
	//}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "mngAdmin/ocrReceipt/ocrReceipt.htm", m)
}

func AdminContent(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/content/content.htm", m)
}

func AdminContentInfo(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/content/contentInfo.htm", m)
}

func AdminBoard(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/content/board.htm", m)
}

func AdminBoardInfo(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/content/boardInfo.htm", m)
}

func AdminBanner(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/content/banner.htm", m)
}

func AdminBannerInfo(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/content/bannerInfo.htm", m)
}

func AdminFeeRate(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/feeRate/feeRate.htm", m)
}

func AdminKakaoWork(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/task/kakaoWork.htm", m)
}

func AdminPartnerMember(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/store/partnerMemberMain.htm", m)
}

func AdminPartnerMemberInfo(c echo.Context) error {

	m := make(map[string]interface{})

	return c.Render(http.StatusOK, "mngAdmin/store/partnerMemberInfo.htm", m)
}

func Address(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["os"] = params["os"]

	return c.Render(http.StatusOK, "etc/address.htm", m)
}
