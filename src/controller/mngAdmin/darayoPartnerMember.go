package mngAdmin

import (
	"biz-web/query/commons"
	"biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetBizPartnerMemberList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	whereQuery := ""

	if params["partnerMemberYn"] == "Y" {
		whereQuery += "AND A.END_DATE >= DATE_FORMAT(now(), '%Y-%m-%d')"
	} else if params["partnerMemberYn"] == "N" {
		whereQuery += "AND A.END_DATE < DATE_FORMAT(now(), '%Y-%m-%d')"
	}

	orderQuery := "ORDER BY A.END_DATE DESC"

	partnerMemberCnt, err := cls.GetSelectDataRequire(mngAdmin.SelectPartnerMemberListCnt+whereQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//wnw
	partnerMemberList, err := cls.GetSelectTypeRequire(mngAdmin.SelectPartnerMemberList+whereQuery+orderQuery+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = partnerMemberCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(partnerMemberCnt[0]["TOTAL_COUNT"])

	if partnerMemberList == nil {
		result["partnerMemberList"] = []string{}
	} else {
		result["partnerMemberList"] = partnerMemberList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetBizPartnerMemberDate(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	partnerMemberDate, err := cls.GetSelectDataRequire(mngAdmin.SelectPartnerMemberDate, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["restId"] = partnerMemberDate[0]["STORE_ID"]
	result["startDate"] = partnerMemberDate[0]["START_DATE"]
	result["endDate"] = partnerMemberDate[0]["END_DATE"]
	result["nextPayDate"] = partnerMemberDate[0]["NEXT_PAY_DAY"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ModBizPartnerMemberDate(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	partnerMemberDateUpdateQuery, err := cls.GetQueryJson(mngAdmin.UpdatePartnerMemberDate, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(partnerMemberDateUpdateQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "Update 오류"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetBizPartnerMemberInfoData(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	partnerMemberInfoDate, err := cls.GetSelectTypeRequire(mngAdmin.SelectPartnerMemberInfoData, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	if partnerMemberInfoDate != nil {
		result["partnerMemberInfoDate"] = partnerMemberInfoDate[0]
	} else {
		result["partnerMemberInfoDate"] = []string{}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetBizPartnerMemberCollectList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	collectListCnt, err := cls.GetSelectDataRequire(mngAdmin.SelectPartnerMemberCollectListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	collectList, err := cls.GetSelectTypeRequire(mngAdmin.SelectPartnerMemberCollectList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = collectListCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(collectListCnt[0]["TOTAL_COUNT"])

	if collectList != nil {
		result["collectList"] = collectList
	} else {
		result["collectList"] = []string{}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetBizPartnerMemberAlarmList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	alarmListCnt, err := cls.GetSelectDataRequire(mngAdmin.SelectPartnerMemberAlarmListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	alarmList, err := cls.GetSelectTypeRequire(mngAdmin.SelectPartnerMemberAlarmList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = alarmListCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(alarmListCnt[0]["TOTAL_COUNT"])

	if alarmList != nil {
		result["alarmList"] = alarmList
	} else {
		result["alarmList"] = []string{}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}
