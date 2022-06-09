package base

import (
	"biz-web/src/controller"
	"bytes"
	"biz-web/src/controller/cls"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

var Query = ``

func base(c echo.Context) error {

	//JSON 배열형태 및 단순 맵형이 아닐때
	bodyBytes, _ := ioutil.ReadAll(c.Request().Body) //바디에서 데이터를 받아 바이트 화 시킴
	var data ExamData                                //구조체 선언
	err := json.Unmarshal(bodyBytes, &data)          //선언한 구조체에 데이터 삽입
	if err != nil {                                  //에러체크
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	params := cls.GetParamJsonMap(c) //값을 map으로 변환

	//select
	resultData, err := cls.GetSelectTypeRequire(Query, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", err.Error()))
	}

	//insert , update
	selectQuery, err := cls.SetUpdateParam(Query, params) //쿼리문 작성
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(selectQuery) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{}) //두개 이상의 결과값을 리턴 할 경우
	result["resultData"] = resultData

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

type ExamData struct {
	ExamName string
	ExamData []ExamInData
}

type ExamInData struct {
	Title    string
	SubTitle int
}
