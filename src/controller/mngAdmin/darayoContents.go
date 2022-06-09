package mngAdmin

import (
	"biz-web/query/commons"
	darayoquery "biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/utils"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetBoardList(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	whereQuery := ""

	if params["bKind"] != "" {
		whereQuery += `
						AND B_KIND = '#{bKind}'`
	}
	if params["boardType"] != "" {
		whereQuery += `
						AND BOARD_TYPE = '#{boardType}'`
	}
	if params["useYn"] != "" {
		whereQuery += `
						AND USE_YN = '#{useYn}'`
	}

	orderQuery := `
					ORDER BY BOARD_ID DESC`

	boardCnt, err := cls.GetSelectDataRequire(darayoquery.SelectBoardCnt+whereQuery+orderQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	boardList, err := cls.GetSelectTypeRequire(darayoquery.SelectBoardList+whereQuery+orderQuery+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = boardCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(boardCnt[0]["TOTAL_COUNT"])

	if boardList == nil {
		result["boardList"] = []string{}
	} else {
		result["boardList"] = boardList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func AddBoard(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	addBoardQuery, err := cls.GetQueryJson(darayoquery.InsertBoard, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	_, err = cls.QueryDB(addBoardQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "서버오류"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetBoardInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	boardInfo, err := cls.GetSelectTypeRequire(darayoquery.SelectBoardInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})

	if boardInfo == nil {
		result["boardList"] = []string{}
	} else {
		result["boardInfo"] = boardInfo[0]
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ModBoardInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	boardUpdateQuery, err := cls.GetQueryJson(darayoquery.UpdateBoardInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}
	_, err = cls.QueryDB(boardUpdateQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "서버오류"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetContentList(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	whereQuery := ""
	if params["useYn"] != "" {
		whereQuery += `AND USE_YN = '` + params["useYn"] + `'`
	}

	orderQuery := `
					ORDER BY CONTENT_ID DESC`

	contentCnt, err := cls.GetSelectDataRequire(darayoquery.SelectContentCnt+whereQuery+orderQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	contentList, err := cls.GetSelectTypeRequire(darayoquery.SelectContentList+whereQuery+orderQuery+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = contentCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(contentCnt[0]["TOTAL_COUNT"])

	if contentList == nil {
		result["contentList"] = []string{}
	} else {
		result["contentList"] = contentList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func AddContent(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	link := params["link"]

	fname := cls.Cls_conf(os.Args)
	url, _ := cls.GetTokenValue("CONTENT_CRAWLER_URL", fname)

	result := make(map[string]interface{})
	result["url"] = url
	result["link"] = link

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetContentInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	contentInfo, err := cls.GetSelectTypeRequire(darayoquery.SelectContentInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	result := make(map[string]interface{})

	if contentInfo == nil {
		result["contentInfo"] = []string{}
	} else {
		result["contentInfo"] = contentInfo[0]
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ModContentInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	contentUpdateQuery, err := cls.GetQueryJson(darayoquery.UpdateContentInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	_, err = cls.QueryDB(contentUpdateQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "서버오류"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetBannerList(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	whereQuery := ""

	if params["bKind"] != "" {
		whereQuery += `
						AND B_KIND = '#{bKind}'`
	}
	if params["bannerType"] != "" {
		whereQuery += `
						AND BOARD_TYPE = '#{bannerType}'`
	}
	if params["useYn"] != "" {
		whereQuery += `
						AND USE_YN = '#{useYn}'`
	}

	orderQuery := `
					ORDER BY BANNER_ID DESC`

	bannerCnt, err := cls.GetSelectDataRequire(darayoquery.SelectBannerCnt+whereQuery+orderQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	bannerList, err := cls.GetSelectTypeRequire(darayoquery.SelectBannerList+whereQuery+orderQuery+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = bannerCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(bannerCnt[0]["TOTAL_COUNT"])

	if bannerList == nil {
		result["bannerList"] = []string{}
	} else {
		result["bannerList"] = bannerList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func AddBanner(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	addBannerQuery, err := cls.GetQueryJson(darayoquery.InsertBanner, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	_, err = cls.QueryDB(addBannerQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "서버오류"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetBannerInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	bannerInfo, err := cls.GetSelectTypeRequire(darayoquery.SelectBannerInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	result := make(map[string]interface{})

	if bannerInfo == nil {
		result["bannerInfo"] = []string{}
	} else {
		result["bannerInfo"] = bannerInfo[0]
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func ModBannerInfo(c echo.Context) error { //복붙해서 사용하자
	params := cls.GetParamJsonMap(c)

	bannerUpdateQuery, err := cls.GetQueryJson(darayoquery.UpdateBannerInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	_, err = cls.QueryDB(bannerUpdateQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "서버오류"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func AddBannerImg(c echo.Context) error {

	c.FormValue("file")

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	form.RemoveAll()

	fileName := ""
	filePath := ""

	m := make(map[string]interface{})

	for _, file := range form.File {
		// Source
		src, err := file[0].Open()
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		defer src.Close()

		fileName = file[0].Filename

		now := time.Now().Format("2006-01-02")
		thisMonth := strings.Replace(now, "-", "", -1)

		filePath = "/app/SharedStorage/upload/BANNER/" + thisMonth
		finishPath := filePath + "/" + fileName

		//경로 생성 (0755-권한)
		if _, err := os.Stat(filePath); err != nil {
			if err = os.Mkdir(filePath, 0755); os.IsNotExist(err) {
				cls.Lprintf(1, "error create directory  : %s\n", err.Error())
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "디렉토리 생성오류"+err.Error()))
			}
		}

		//출력 파일 생성
		dst, err1 := os.Create(finishPath)
		if err1 != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "출력파일 생성 오류"+err.Error()))
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "복사 오류"+err.Error()))
		}
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultPath"] = strings.Replace(filePath+"/"+fileName, "/app", "", -1)

	return c.JSON(http.StatusOK, m)
}
