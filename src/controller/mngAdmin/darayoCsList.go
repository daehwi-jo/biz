package mngAdmin

import (
	"biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetBizCsList
func GetBizCsList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectTypeRequire(mngAdmin.SelectCsList, params, c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	if resultData == nil {
		m["resultData"] = []string{}
	} else {
		m["resultData"] = resultData
	}

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	return c.JSON(http.StatusOK, m)
}

// AddBizCsContent
// notice : Insert
func AddBizCsContent(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	qeury, err := cls.SetUpdateParam(mngAdmin.InsertCsContent, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(qeury) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}
