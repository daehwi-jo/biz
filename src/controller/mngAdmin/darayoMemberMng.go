package mngAdmin

import (
	"biz-web/query/commons"
	"biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	_struct "biz-web/src/controller/struct"
	"biz-web/src/controller/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// GetUserList 유저목록 //쿼리문에 로직 빼자 (만들다 말았음)
func GetUserList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	addQuery := ""
	switch params["sortKey"] {
	case "all":
		addQuery = ``
		break
	case "store":
		addQuery = `AND USER_TY = '1'`
		break
	case "member":
		addQuery = `AND USER_TY = '0'`
		break
	}

	userListCnt, err := cls.GetSelectDataRequire(mngAdmin.SelectUserListCnt+addQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	userList, err := cls.GetSelectTypeRequire(mngAdmin.SelectUserList+addQuery+commons.JoinDateOrderBy+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = userListCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(userListCnt[0]["TOTAL_COUNT"])

	if userList == nil {
		result["userList"] = []string{}
	} else {
		result["userList"] = userList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

// GetUserInfo 유저정보
// notice:
func GetUserInfo(c echo.Context) error {
	var params = cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectType(mngAdmin.SelectUserInfo, params, c)

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

// ResetPassword 비밀번호 초기화
// Updata - SetUpdateParam사용 -쿼리문 생성 -> ExecDB()사용
func ResetPassword(c echo.Context) error { //수정해야함
	var params = cls.GetParamJsonMap(c)

	//해당 데이터로 쿼리문 생성
	query, err := cls.SetUpdateParam(mngAdmin.UpdateUserPassword, params) //업데이트 하나로만
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//sysUser pass 초기화
	userData, err := cls.GetSelectData(mngAdmin.SelectUserLoginId, params, c)
	params["loginId"] = userData[0]["loginId"]

	query2, err := cls.SetUpdateParam(mngAdmin.UpdateSysUserPassword, params) //업데이트 하나로만
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(query2) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

// ModifyUserInfo 유저 정보 변경
// notice :
func ModifyUserInfo(c echo.Context) error { //수정해야함
	params := cls.GetParamJsonMap(c)

	userData, err := cls.GetSelectDataRequire(mngAdmin.SelectUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if userData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	if params["email"] == "" {
		params["email"] = userData[0]["email"]
	}
	if params["hpNo"] == "" {
		params["hpNo"] = userData[0]["hpNo"]
	}
	if params["useYn"] == "" {
		params["useYn"] = userData[0]["useYn"]
	}

	if params["useYn"] == "N" {
		params["authStat"] = "3"

		query, err := cls.SetUpdateParam(mngAdmin.UpdateUserDelGrp, params) //연결되어있는 장부에서 모두 떠나기
		_, err = cls.QueryDB(query)                                         // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

	}

	query2, err := cls.SetUpdateParam(mngAdmin.UpdateUserInfo, params) //업데이트 하나로만
	_, err = cls.QueryDB(query2)                                       // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

// GetBizUserCsInfo
// notice : (Java = BizUserMngDesc)
func GetBizUserCsInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	userInfo, err := cls.GetSelectType(mngAdmin.SelectUserInfo, params, c) //userInfo불러옴
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if userInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	userGrpData, err := cls.GetSelectTypeRequire(mngAdmin.SelectUserGrpList, params, c) //userGrpData불러옴
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	userGrpData2, err := cls.GetSelectTypeRequire(mngAdmin.SelectUserGrpList2, params, c) //userGrpData불러옴
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	unUsedCouponList, err := cls.GetSelectTypeRequire(mngAdmin.SelectUnUsedCouponList, params, c) //userGrpData불러옴
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	billingCardInfo, err := cls.GetSelectTypeRequire(mngAdmin.SelectTpayBillingCardList, params, c) //userGrpData불러옴
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["userInfo"] = userInfo[0]
	result["userGrpData"] = userGrpData
	result["userGrpData2"] = userGrpData2
	result["unUsedCouponList"] = unUsedCouponList
	result["billingCardInfo"] = billingCardInfo

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func CheckGifticon(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	url, _ := cls.GetTokenValue("GIFTICON_USED_CHECK", fname)

	if params["orderNo"] == "" {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "주문번호가 없습니다."))
	}

	resp, err := http.Get(url + "?orderNo=" + params["orderNo"])
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var gifticonState _struct.GifticonState
	err = json.Unmarshal(body, &gifticonState)
	if err != nil {
		panic(err)
	}

	result := make(map[string]interface{})
	if state := gifticonState.ResultData.CouponInfo.CpStatus; state != "" {
		result["state"] = "1"
	} else {
		result["state"] = "0"
	}

	m := make(map[string]interface{})
	m["resultCode"] = gifticonState.ResultCode
	m["resultMsg"] = gifticonState.ResultMsg
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ExtendGifticon(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	usedUpdateQuery, err := cls.SetUpdateParam(mngAdmin.UpdateCouponUsedDate, params) //쿼리 업데이트
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(usedUpdateQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	darUseDay := time.Now()
	result := make(map[string]interface{})
	result["darUseDay"] = darUseDay.Format("2006-01-02")

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}
