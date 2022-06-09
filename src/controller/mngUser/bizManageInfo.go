package mngUser

import (
	bizMngQuery "biz-web/query/mngUser"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetManagerInfoMng(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	companyInfo, err := cls.GetSelectTypeRequire(bizMngQuery.SelectCompanyInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if companyInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	managerInfo, err := cls.GetSelectTypeRequire(bizMngQuery.SelectManagerInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if managerInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	result := make(map[string]interface{})
	result["companyInfo"] = companyInfo
	result["managerInfo"] = managerInfo

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ModifyCompanyInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	if params["authorCd"] == "BM" {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "회사 관리자만 변경할 수 있습니다. 회사 관리자에게 문의 해주세요."))
	}

	companyInfo, err := cls.GetSelectDataRequire(bizMngQuery.SelectCompanyInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if companyInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}
	if params["addr"] == "" {
		params["addr"] = companyInfo[0]["ADDR"]
	}
	if params["addr2"] == "" {
		params["addr2"] = companyInfo[0]["ADDR2"]
	}
	if params["ceoNm"] == "" {
		params["ceoNm"] = companyInfo[0]["CEO_NM"]
	}
	if params["companyNm"] == "" {
		params["companyNm"] = companyInfo[0]["COMPANY_NM"]
	}
	if params["homepage"] == "" {
		params["homepage"] = companyInfo[0]["HOMEPAGE"]
	}
	if params["tel"] == "" {
		params["tel"] = companyInfo[0]["TEL"]
	}

	query, err := cls.SetUpdateParam(bizMngQuery.UpdateCompanyInfo, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func ModifyManagerInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	if params["authorCd"] == "BM" {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "회사 관리자만 변경할 수 있습니다. 회사 관리자에게 문의 해주세요."))
	}

	managerInfo, err := cls.GetSelectDataRequire(bizMngQuery.SelectManagerInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if managerInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	if params["dept"] == "" {
		params["dept"] = managerInfo[0]["DEPT"]
	}
	if params["email"] == "" {
		params["email"] = managerInfo[0]["EMAIL"]
	}
	if params["class"] == "" {
		params["class"] = managerInfo[0]["COURSE"] //??? class??? course
	}
	if params["tel"] == "" {
		params["tel"] = managerInfo[0]["TEL"]
	}
	//이름 선택이 되는데 이것도 바뀌어야 하는거 아닌지

	query, err := cls.SetUpdateParam(bizMngQuery.UpdateManagerInfo, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}
