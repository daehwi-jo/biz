package mngAdmin

import (
	"biz-web/query/commons"
	darayoquery "biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetFeeRateList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	whereQuery := ""

	if params["restNm"] != "" {
		whereQuery += "AND A.REST_NM LIKE '%#{restNm}%'"
	}

	if params["useYn"] != "" {
		whereQuery += "AND A.USE_FEES_YN = '#{useYn}'"
	}

	listCnt, err := cls.GetSelectDataRequire(darayoquery.SelectFeeRateListCnt+whereQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	list, err := cls.GetSelectDataRequire(darayoquery.SelectFeeRateList+whereQuery+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	result := make(map[string]interface{})
	result["totalCount"] = listCnt[0]["TOTAL_COUNT"]
	if list != nil {
		result["feeRateList"] = list
	} else {
		result["feeRateList"] = []string{}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func PutFeeRate(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	//update
	feeRateUpdateQuery, err := cls.GetQueryJson(darayoquery.UpdateFeeRate, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	_, err = cls.QueryDB(feeRateUpdateQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "서버오류"))
	}

	//insert
	feeRateInsertQuery, err := cls.GetQueryJson(darayoquery.InsertFeeRate, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}
	_, err = cls.QueryDB(feeRateInsertQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "서버오류"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}
