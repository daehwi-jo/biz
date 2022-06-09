package mngAdmin

import (
	"biz-web/query"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetSysUserSearch(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	uIdList, err := cls.GetSelectDataRequire(query.SearchMemberSql(params["searchKey"]), params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	grpIdList, err := cls.GetSelectDataRequire(query.SearchBookSql(params["searchKey"]), params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	restIdList, err := cls.GetSelectDataRequire(query.SearchStoreSql, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})

	if uIdList == nil {
		result["uIdList"] = []string{}
	} else {
		result["uIdList"] = uIdList
	}

	if grpIdList == nil {
		result["grpIdList"] = []string{}
	} else {
		result["grpIdList"] = grpIdList
	}

	if restIdList == nil {
		result["restIdList"] = []string{}
	} else {
		result["restIdList"] = restIdList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetSysUserSettingMainData(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	sysUserInfo, err := cls.GetSelectDataRequire(query.SelectSysUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	sysUserInfoList, err := cls.GetSelectDataRequire(query.SelectSysUserInfoList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	if sysUserInfo != nil {
		result["sysUserInfo"] = sysUserInfo[0]
	} else {
		result["sysUserInfo"] = []string{}
	}
	if sysUserInfoList != nil {
		result["sysUserInfoList"] = sysUserInfoList
	} else {
		result["sysUserInfoList"] = []string{}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func AddSysUser(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	Query, err := cls.SetUpdateParam(query.InsertSysUser, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_, err = cls.QueryDB(Query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func ModSysUser(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	Query, err := cls.GetQueryJson(query.UpdateSysUser, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_, err = cls.QueryDB(Query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func ModSysUserAuth(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	Query, err := cls.GetQueryJson(query.UpdateSysUserUseYn, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_, err = cls.QueryDB(Query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func CheckSysUserPw(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	status, err := cls.GetSelectDataRequire(query.SelectSysUserPwCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if status[0]["state"] == "1" {
		return c.JSON(http.StatusOK, controller.SetErrResult("00", "응답 성공"))
	} else {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "비밀번호 오류"))
	}
}
