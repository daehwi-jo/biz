package mngAdmin

import (
	"biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	_struct "biz-web/src/controller/struct"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetOcrTextList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	result := make(map[string]interface{})

	textList, err := cls.GetSelectDataRequire(mngAdmin.SelectTextList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if textList != nil {
		result["textList"] = textList
	} else {
		result["textList"] = []string{}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetTextData(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	result := make(map[string]interface{})

	data, err := cls.GetSelectData(mngAdmin.SelectTextData, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result["receiptData"] = data[0]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func SetReceiptText(c echo.Context) error {
	var data _struct.ReceiptTextDataList

	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)

	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_ = c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	//TODO 로직 수정
	for i := 0; i < len(data.List); i++ {
		var p = make(map[string]string)
		p["restId"] = data.List[i].RestId
		p["bizNum"] = data.List[i].BizNum
		p["textId"] = data.List[i].TextId
		p["isMenu"] = data.List[i].IsMenu
		p["menuId"] = data.List[i].MenuId
		p["menuNm"] = data.List[i].MenuNm
		p["menuPrice"] = data.List[i].MenuPrice

		newMenuIds, err := cls.GetSelectData(mngAdmin.SelectNewMenuId, p, c) //새로운 메뉴 아이디 불러옴
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		p["newMenuId"] = newMenuIds[0]["newMenuId"]

		checker, err := cls.GetSelectData(mngAdmin.SelectMenuChecker, p, c) //메뉴 있는지 상태값 가져오기
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		var state string

		if len(state) != 0 { //1=금액 안 맞음, 2=일치
			state = checker[0]["state"]
			p["menuId"] = checker[0]["menuId"]
		} else { //code : 0=메뉴없음
			state = "0"
			p["menuId"] = p["newMenuId"]
		}

		tx, err := cls.DBc.Begin()
		if err != nil {
			//return "5100", errors.New("begin error")
		}

		//메뉴 판별
		if p["isMenu"] == "Y" {
			if state != "2" { //메뉴에 일치하지 않는경우
				insertMenu, err := cls.GetQueryJson(mngAdmin.InsertOCRNewMenu, p)
				if err != nil {
					return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리문 생성 오류"))
				}
				_, err = tx.Exec(insertMenu)
				if err != nil {
					return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
				}

				if state == "1" { //기존 메뉴 비활성
					updateMenu, err := cls.GetQueryJson(mngAdmin.UpdateOCRMenuYn, p)
					if err != nil {
						return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리문 생성 오류"))
					}
					_, err = tx.Exec(updateMenu)
					if err != nil {
						return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
					}
				}

				p["menuId"] = p["newMenuId"] //사용하는 메뉴가 바뀜으로 변경
			}
		}

		//TEXT데이터 갱신
		updateTextMenu, err := cls.GetQueryJson(mngAdmin.UpdateOCRTextMenuId, p)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리 생성오류2"))
		} else {
			_, err = cls.QueryDB(updateTextMenu)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "텍스트 업데이트 오류"))
			}
		}

		//에러 확인해서 롤백하질 커미할지 결정
		if err != nil {
			_ = tx.Rollback() //롤백
		} else {
			err = tx.Commit()
			if err != nil {
				lprintf(1, "[ERROR] SetMakeFee UpdateStoreInfoQuery error(%s) \n", err.Error())
			}
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetReceiptList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	result := make(map[string]interface{})

	receiptList, err := cls.GetSelectDataRequire(mngAdmin.SelectReceiptList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if receiptList != nil {
		result["receiptList"] = receiptList
	} else {
		result["receiptList"] = []string{}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetReceiptData(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	result := make(map[string]interface{})

	receiptData, err := cls.GetSelectDataRequire(mngAdmin.SelectReceiptData, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	receiptDataMenu, err := cls.GetSelectDataRequire(mngAdmin.SelectReceiptMenuData, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result["receiptData"] = receiptData[0]

	if receiptDataMenu == nil {
		result["menuData"] = []string{}
	} else {
		result["menuData"] = receiptDataMenu
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func SetReceiptData(c echo.Context) error {
	var params _struct.OcrReceiptDataList

	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)

	err := json.Unmarshal(bodyBytes, &params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	_ = c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var check []map[string]string
	var query string
	var newMenuId int

	p := make(map[string]string)
	p["restId"] = params.RestId
	p["bizNum"] = params.BizNum
	p["receiptId"] = params.ReceiptId
	p["state"] = params.State
	p["totalAmt"] = params.TotalAmt
	p["aprvDt"] = params.AprvDt
	p["aprvNo"] = params.AprvNo

	receiptMenuList, err := cls.SelectData(mngAdmin.SelectReceiptMenuList, p)
	receiptMenuLength := len(receiptMenuList)

	//  transaction start
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			_ = tx.Rollback() //롤백
		}
	}()

	for i, m := range params.List {
		p["receiptMenuId"] = m.ReceiptMenuId
		p["menuNm"] = m.MenuNm
		p["menuPrice"] = m.MenuPrice
		p["menuEa"] = m.MenuEa
		p["menuAmt"] = m.MenuAmt

		//메뉴확인
		check, err = cls.SelectData(mngAdmin.SelectMenuCheck, p) //code : 0=메뉴없음 ,1=금액 안 맞음,2=일치
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "check 오류"))
		}

		//메뉴 insert update
		if len(check) == 0 { //메뉴가 아예 없음
			p["menuId"] = "1"

			if newMenuId == 0 {
				newMenuId = 1
			}
			p["newMenuId"] = strconv.Itoa(newMenuId)
			newMenuId++

			query, err = cls.GetQueryJson(mngAdmin.InsertOCRNewMenu, p)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리문 생성 오류"))
			}
			_, err = tx.Exec(query)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			}
			p["menuId"] = p["newMenuId"]
			newMenuId++
		} else {
			state := check[0]["state"]
			p["menuId"] = check[0]["menuId"]

			if newMenuId == 0 {
				newMenuId, _ = strconv.Atoi(check[0]["newMenuId"])
			}
			p["newMenuId"] = strconv.Itoa(newMenuId)
			newMenuId++

			if state != "2" { //있는 메뉴와 완벽하게 일치하지 않을때
				query, err = cls.GetQueryJson(mngAdmin.InsertOCRNewMenu, p) //신규 메뉴 추가
				if err != nil {
					return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리문 생성 오류"))
				}
				_, err = tx.Exec(query)
				if err != nil {
					return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
				}

				if state == "1" { //기존 메뉴 비활성
					query, err = cls.GetQueryJson(mngAdmin.UpdateOCRMenuYn, p)
					if err != nil {
						return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리문 생성 오류"))
					}
					_, err = tx.Exec(query)
					if err != nil {
						return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
					}
				}

			} else { //기존의 메뉴와 완벽히 일치 할떄
				p["newMenuId"] = check[0]["menuId"]
			}
		}

		if i < receiptMenuLength { //receipt menu 에 없는 경우
			query, err = cls.GetQueryJson(mngAdmin.UpdateReceiptMenuData, p)
		} else { //receipt menu 에 있는 경우
			p["receiptMenuId"] = strconv.Itoa(i + 1)
			query, err = cls.GetQueryJson(mngAdmin.InsertReceiptMenuData, p)
		}
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리문 생성 오류"))
		}
		_, err = tx.Exec(query)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	//남은 항목 지우기
	for i := len(params.List); i < receiptMenuLength; i++ {
		p["receiptMenuId"] = receiptMenuList[i]["receiptMenuId"]
		Query, err := cls.GetQueryJson(mngAdmin.DeleteReceiptMenuData, p)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리문 생성 오류"))
		}
		_, err = tx.Exec(Query)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	//영수증 정보 업데이트
	Query, err := cls.GetQueryJson(mngAdmin.UpdateReceiptData, p)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리문 생성 오류"))
	}
	_, err = tx.Exec(Query)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		lprintf(1, "[ERROR] SetMakeFee UpdateStoreInfoQuery error(%s) \n", err.Error())
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}
