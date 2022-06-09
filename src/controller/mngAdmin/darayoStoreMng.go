package mngAdmin

import (
	"biz-web/query/commons"
	"biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"biz-web/src/controller/struct"
	"biz-web/src/controller/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var lprintf func(int, string, ...interface{}) = cls.Lprintf

func GetBizStoreList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	addQuery := ``
	if params["searchKey"] == "restNm" {
		addQuery = `AND A.REST_NM LIKE '%#{searchKeyword}%'`
	} else {

		addr := params["searchKeyword"]

		str := strings.Split(addr, " ")

		if len(str) >= 2 {
			for i := 0; i < len(str); i++ {
				params["addr"+strconv.Itoa(i)] = str[i][0:6]
				addQuery += `and A.ADDR LIKE '%#{addr` + strconv.Itoa(i) + `}%'`
			}
		} else {
			params["addr"] = str[0]
			addQuery += `and A.ADDR LIKE '%#{addr}%'`
		}
	}
	addQuery += `ORDER BY REG_DATE DESC`

	storeCnt, err := cls.GetSelectDataRequire(mngAdmin.SelectStoreCnt+addQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	storeList, err := cls.GetSelectTypeRequire(mngAdmin.SelectStoreList+addQuery+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = storeCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]
	result["totalPage"] = utils.GetTotalPage(storeCnt[0]["TOTAL_COUNT"])

	if storeList == nil {
		result["storeList"] = []string{}
	} else {
		result["storeList"] = storeList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetBizStoreCsInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	//e로 두번 못태움?
	storeInfo, err := cls.GetSelectType(mngAdmin.SelectStoreInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	bankList, err := cls.GetSelectTypeRequire(mngAdmin.SelectBankList, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	storeEtc, err := cls.GetSelectTypeRequire(mngAdmin.SelectStoreEtc, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	if storeInfo == nil {
		result["storeInfo"] = []string{}
	} else {
		result["storeInfo"] = storeInfo
	}
	if bankList == nil {
		result["bankList"] = []string{}
	} else {
		result["bankList"] = bankList
	}
	if storeEtc == nil {
		result["storeEtc"] = []string{}
	} else {
		result["storeEtc"] = storeEtc
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

//가맹점 정보 수정
func ModStoreInfo(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	accountCertYn := "N"
	if params["accountNm"] != "" {
		accountCertYn = "Y"
	}
	params["accountCertYn"] = accountCertYn

	params["memo"] = strings.Replace(params["memo"], `'`, `\'`, -1)

	params["notice"] = strings.Replace(params["notice"], `'`, `\'`, -1)

	Query, err := cls.SetUpdateParam(mngAdmin.UpdateStoreInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	Query, err = cls.GetQueryJson(mngAdmin.UpdateStoreEtc, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

//연결 장부 목록
func GetStoreLinkBookList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	listCnt, err := cls.GetSelectData(mngAdmin.SelectStoreLinkBookListCnt, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	list, err := cls.GetSelectTypeRequire(mngAdmin.SelectStoreLinkBookList+commons.PagingQuery, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData := make(map[string]interface{})
	resultData["linkBookList"] = list
	resultData["totalCnt"] = listCnt[0]["total"]
	resultData["pageSize"] = params["pageSize"]
	resultData["pageNo"] = params["pageNo"]
	resultData["totalPage"] = utils.GetTotalPage(listCnt[0]["total"])

	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

//연결장부 검색
func GetStoreLinkBookSearch(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	params["offSet"] = utils.GetOffset(params["pageSize"], params["pageNo"])

	listCnt, err := cls.GetSelectData(mngAdmin.SelectStoreBookSearchCnt, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	list, err := cls.GetSelectTypeRequire(mngAdmin.SelectStoreBookSearch+commons.PagingQuery, params, c) //null 체크 안함

	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData := make(map[string]interface{})
	resultData["linkBookList"] = list
	resultData["totalCnt"] = listCnt[0]["total"]
	resultData["pageSize"] = params["pageSize"]
	resultData["pageNo"] = params["pageNo"]
	resultData["totalPage"] = utils.GetTotalPage(listCnt[0]["total"])

	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

//연결장부 상세 정보
func GetStoreLinkBookInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	list, err := cls.GetSelectTypeRequire(mngAdmin.SelectStoreLinkBookInfo, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m["resultData"] = list

	return c.JSON(http.StatusOK, m)
}

//장부 연결 수정
func ModStoreLinkBook(c echo.Context) error {
	params := cls.GetParamJsonMap(c)
	Query, err := cls.GetQueryJson(mngAdmin.InsertBookLink, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

//카테고리록불러오기
func GetStoreCategory(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	list, err := cls.GetSelectTypeRequire(mngAdmin.SelectCategoryList, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if list == nil {
		m["resultData"] = []string{}
	} else {
		m["resultData"] = list
	}

	return c.JSON(http.StatusOK, m)
}

//카테고리 추가
func AddStoreCategory(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	news, err := cls.GetSelectData(mngAdmin.SelectNewCategory, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	params["categoryId"] = news[0]["categoryId"]
	params["categoryNm"] = news[0]["categoryNm"]
	params["codeId"] = news[0]["codeId"]

	Query, err := cls.SetUpdateParam(mngAdmin.InsertCategory, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func ModStoreCategory(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	Query, err := cls.SetUpdateParam(mngAdmin.UpdateCategory, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	if params["useYn"] == "N" {
		Query2, err := cls.SetUpdateParam(mngAdmin.UpdateMenuUseYn, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		} else {
			_, err = cls.QueryDB(Query2) // 쿼리 실행
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			}
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

//메뉴카테고리목록불러오기
func GetStoreCategoryMenu(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	menu, err := cls.GetSelectTypeRequire(mngAdmin.SelectStoreItemList, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	category, err := cls.GetSelectTypeRequire(mngAdmin.SelectCategoryList, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData := make(map[string]interface{})
	resultData["menuData"] = menu
	resultData["categoryData"] = category

	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

//메뉴 추가
func AddStoreMenu(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	seq, err := cls.GetSelectDataRequire(mngAdmin.ItemNewSeq, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	params["seqNo"] = seq[0]["seqNo"]

	Query, err := cls.GetQueryJson(mngAdmin.InsertMenu, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}
func AddMenuImage(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	if params["restId"] == "" {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류"))
	}

	if params["menuId"] == "" {
		seq, err := cls.GetSelectDataRequire(mngAdmin.ItemNewSeq, params, c) //null 체크 안함
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "파라미터 오류2"))
		}
		params["menuId"] = seq[0]["seqNo"]
	}

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
		src, err := file[0].Open()
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		defer src.Close()

		extend := strings.Split(file[0].Filename, ".")
		fileName = params["restId"] + "_" + params["menuId"] + "." + extend[len(extend)-1]

		filePath = "/app/SharedStorage/upload/REST_ITEM/" + params["restId"]

		finishPath := filePath + "/" + fileName

		////경로 생성 (0755-권한)
		if _, err := os.Stat(filePath); err != nil {
			if err = os.Mkdir(filePath, 0755); err != nil {
				cls.Lprintf(1, "error create directory  : %s\n", err.Error())
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "디렉토리 생성오류 "+err.Error()))
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

//메뉴 수정
func ModStoreMenu(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	Query, err := cls.SetUpdateParam(mngAdmin.UpdateMenu, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

//부대시설 목록
func GetFacilityList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	list, err := cls.GetSelectDataRequire(mngAdmin.GetStoreService, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	if list == nil {
		m["resultData"] = []string{}
	} else {
		m["resultData"] = list
	}

	return c.JSON(http.StatusOK, m)
}

//부대시설 수정
func AddFacility(c echo.Context) error {

	params := cls.GetParamJsonMap(c) //json을 map으로 컨버팅

	Query, err := cls.GetQueryJson(mngAdmin.InsertStoreService, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

//부대시설 수정
func ModFacilityList(c echo.Context) error {
	var data _struct.StoreServiceObj

	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)

	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	params := cls.GetParamJsonMap(c) //json을 map으로 컨버팅

	for i := 0; i < len(data.FacilityList); i++ {
		facilityData := make(map[string]string)
		facilityData["serviceInfo"] = data.FacilityList[i].ServiceInfo
		facilityData["useYn"] = data.FacilityList[i].UseYn
		facilityData["noticeYn"] = data.FacilityList[i].NoticeYn
		facilityData["serviceId"] = data.FacilityList[i].ServiceId
		facilityData["serviceNm"] = data.FacilityList[i].ServiceNm
		facilityData["restId"] = params["restId"]
		Query, err := cls.SetUpdateParam(mngAdmin.UpdateStoreService, facilityData)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		} else {
			_, err = cls.QueryDB(Query) // 쿼리 실행
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			}
		}

	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetBankList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	result := make(map[string]interface{})

	list, err := cls.GetSelectTypeRequire(mngAdmin.SelectBankList, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if list == nil {
		result["bankList"] = []string{}
	} else {
		result["bankList"] = list
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func SetStoreImg(c echo.Context) error {

	c.FormValue("file")

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	form.RemoveAll()

	params := cls.GetParamJsonMap(c)

	fileName := ""
	fileSize := ""
	filePath := ""
	sysFileName := ""

	restId := params["restId"]
	//파일명에 타임셋 추가
	now := time.Now().Format("_20060102150405")

	m := make(map[string]interface{})

	for _, file := range form.File {
		// Source
		src, err := file[0].Open()
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		defer src.Close()

		println("11")

		fileName = file[0].Filename
		fileSize = fmt.Sprintf("%d", file[0].Size)

		sysFileName = restId + now + params["extend"]

		now := time.Now().Format("2006-01")
		thisMonth := strings.Replace(now, "-", "", -1)

		filePath = "/SharedStorage/upload/CORP_CI/" + thisMonth
		finishPath := filePath + "/" + sysFileName

		os.Mkdir("/app"+filePath, 0755)

		dst, err := os.Create("/app" + finishPath)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	if fileSize == "" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "파일을 선택해주세요."))
	}

	params["fileName"] = fileName
	params["fileSize"] = fileSize
	params["filePath"] = filePath
	params["fileTy"] = "1"
	params["sysFileName"] = sysFileName

	exQuery := ""

	restImgInfo, err := cls.GetSelectDataRequire(mngAdmin.SelectRestImgInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if restImgInfo == nil {

		restFileSeq, err := cls.GetSelectData(mngAdmin.SelectRestFileSeq, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		params["fileNo"] = restFileSeq[0]["FILE_NO"]

		exQuery = mngAdmin.InsertRestImg
	} else {
		exQuery = mngAdmin.UpdateRestImg
	}

	query, err := cls.SetUpdateParam(exQuery, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	_, err = cls.QueryDB(query) // 쿼리 실행
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = params

	return c.JSON(http.StatusOK, m)
}

type TpayResult struct {
	Result_cd    string `json:"result_cd"`    //
	Result_msg   string `json:"result_msg"`   //
	Account_name string `json:"account_name"` //
}

// 계좌실명조회
func AcctNameSearch(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	bankCode := params["bankCode"]
	accountNo := params["accountNo"]
	//buyerAuthNum :=params["buyerAuthNum"]

	url := "https://webtx.tpay.co.kr/api/v1/acct_name_search?"
	api_key := "xG3E5I+uuUvo+3ui/PKAPhxhutmQteOf3UiZ3PYG/zpO6fHsJZdlY28GOAWP09Kp7ArmIQdFlG7elvpTf/AKqQ=="
	mid := "darayo001m"
	bank_code := bankCode
	account := accountNo
	buyer_auth_num := "2222"

	urlParameters := "api_key=" + api_key + "&mid=" + mid + "&bank_code=" + bank_code + "&account=" + account + "&buyer_auth_num=" + buyer_auth_num
	resp, err := http.Get(url + urlParameters)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result TpayResult
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	rAccountName := result.Account_name
	rcode := result.Result_cd
	if rcode != "000" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "올바른 '계좌번호'를 넣어주세요."))
	}

	msg := "예금주명 '" + rAccountName + "'(이)가 맞습니까?"

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["msg"] = msg
	m["account_name"] = rAccountName

	return c.JSON(http.StatusOK, m)
}

//메뉴 목록 불러오기
func GetStoreMenu(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	resultData := make(map[string]interface{})

	addQuery := ""
	if params["useVisible"] == "Y" {
		addQuery = `
					AND A.USE_YN = 'Y'`
	}
	if params["codeId"] != "all" {
		addQuery += `
					AND B.CODE_ID = '#{codeId}'`
	}

	menu, err := cls.GetSelectTypeRequire(mngAdmin.SelectStoreCategoryItemList+addQuery, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if menu == nil {
		resultData["menuData"] = []string{}
	} else {
		resultData["menuData"] = menu
	}

	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

func GetChargeList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	check, err := cls.GetSelectData(mngAdmin.CheckPaymentUse, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	chargeList, err := cls.GetSelectData(mngAdmin.SelectChargeList, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultData := make(map[string]interface{})
	resultData["useCheck"] = check[0]["allUseYn"]
	if chargeList != nil {
		resultData["chargeList"] = chargeList
	} else {
		resultData["chargeList"] = []string{}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)
}

func ModChargeList(c echo.Context) error {
	var data _struct.StoreChargeObj

	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)

	err := json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	p := make(map[string]string)
	p["allUseYn"] = data.AllUseYn
	p["restId"] = data.RestId

	Query, err := cls.SetUpdateParam(mngAdmin.UpdatePaymentUse, p)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리 맵핑 오류1"))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	for _, d := range data.ChargeList {
		p["seqNo"] = d.SeqNo
		p["amt"] = d.Amt
		p["addAmt"] = d.AddAmt
		p["useYn"] = d.UseYn

		Query2, err := cls.SetUpdateParam(mngAdmin.UpdateChargeItem, p)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리 맵핑 오류2"))
		} else {
			_, err = cls.QueryDB(Query2) // 쿼리 실행
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			}
		}

	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func AddChargeItem(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	newSeqNo, err := cls.GetSelectData(mngAdmin.SelectChargeNewSeqNo, params, c) //null 체크 안함
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	params["newSeqNo"] = newSeqNo[0]["newSeqNo"] + "_" + params["restId"]

	Query, err := cls.SetUpdateParam(mngAdmin.InsertChargeItem, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "쿼리 맵핑 오류"))
	} else {
		_, err = cls.QueryDB(Query) // 쿼리 실행
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func GetUnPaidList(c echo.Context) error {

	dprintf(4, c, "call GetUnPaidList\n")

	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectData(mngAdmin.SelectUnpaidListCount, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUnpaidListCount err(%s) \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	unPaidList, err := cls.GetSelectType(mngAdmin.SelectUnpaidList, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUnpaidList err(%s) \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if unPaidList == nil {
		LinkData := make(map[string]interface{})
		LinkData["totalAmt"] = 0
		LinkData["unPaidList"] = []string{}
		LinkData["bookNm"] = ""
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = LinkData
		return c.JSON(http.StatusOK, m)
	}

	totalAmt, _ := strconv.Atoi(resultData[0]["TOTAL_AMT"])
	orderCnt, _ := strconv.Atoi(resultData[0]["orderCnt"])

	result := make(map[string]interface{})
	result["totalCnt"] = orderCnt
	result["totalAmt"] = totalAmt
	result["unPaidList"] = unPaidList
	result["bookNm"] = resultData[0]["BOOK_NM"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)

}

func GetStorePaymentList(c echo.Context) error {

	dprintf(4, c, "call GetStorePaymentList\n")

	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectData(mngAdmin.SelectStorePaymentListCount, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUnpaidListCount err(%s) \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	unPaidList, err := cls.GetSelectType(mngAdmin.SelectStorePaymentList, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUnpaidList err(%s) \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if unPaidList == nil {
		LinkData := make(map[string]interface{})
		LinkData["totalAmt"] = 0
		LinkData["paidOkList"] = []string{}
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = LinkData
		return c.JSON(http.StatusOK, m)
	}

	orderCnt, _ := strconv.Atoi(resultData[0]["orderCnt"])

	result := make(map[string]interface{})
	result["totalCnt"] = orderCnt
	result["paidOkList"] = unPaidList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)

}

// 매장 정산 처리
func SetPaidOk(c echo.Context) error {

	dprintf(4, c, "call SetPaidOk\n")

	params := cls.GetParamJsonMap(c)

	totalAmt, _ := strconv.Atoi(params["totalAmt"])

	if totalAmt == 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "데이터가 부족합니다."))
	}

	unPaidInfo, err := cls.GetSelectData(mngAdmin.SelectUnpaidListCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if unPaidInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "요청하신 정산 기준일에 미처리 정산 내역이 없습니다."))
	}

	rtotalAmt, _ := strconv.Atoi(unPaidInfo[0]["TOTAL_AMT"])
	userId := unPaidInfo[0]["USER_ID"]

	if rtotalAmt != totalAmt {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "정산금액이 요청하신 내역과 다릅니다."))
	}

	params["payChannel"] = "01"
	params["userId"] = userId

	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -매장 충전(SetStoreCharging)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	// 충전 히스토리 insert
	params["creditAmt"] = strconv.Itoa(totalAmt)
	params["addAmt"] = "0"
	params["searchTy"] = "2"
	params["paymentTy"] = "3"
	params["userTy"] = "2"

	now := time.Now()
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	diff := now.Sub(then)

	moid := fmt.Sprintf("%d", diff.Milliseconds())
	moid = moid + params["userId"]
	params["moid"] = moid

	UpdateOrderPaidQuery, err := cls.SetUpdateParam(mngAdmin.UpdateOrderPaid, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateOrderPaidQuery parameter fail"))
	}

	_, err = tx.Exec(UpdateOrderPaidQuery)
	dprintf(4, c, "call set Query (%s)\n", UpdateOrderPaidQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateOrderPaidQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	InsertPaymentHistoryQuery, err := cls.SetUpdateParam(mngAdmin.InsertPaymentHistory, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertPaymentHistoryQuery parameter fail"))
	}

	_, err = tx.Exec(InsertPaymentHistoryQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertPaymentHistoryQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPaymentHistoryQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

// 매장 정산 취소
func SetPaidCancel(c echo.Context) error {

	dprintf(4, c, "call SetPaidCancel\n")

	params := cls.GetParamJsonMap(c)

	cancelInfo, err := cls.GetSelectData(mngAdmin.SelectStoreCancelCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	cancelCnt, _ := strconv.Atoi(cancelInfo[0]["CancelCnt"])

	if cancelCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "이미 취소된 결제 입니다."))
	}

	chargeInfo, err := cls.GetSelectData(mngAdmin.SelectStoreChargeInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if chargeInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "결제 내역이 없습니다."))
	}

	params["creditAmt"] = chargeInfo[0]["CREDIT_AMT"]
	params["addAmt"] = chargeInfo[0]["ADD_AMT"]
	params["userId"] = chargeInfo[0]["USER_ID"]
	params["payInfo"] = chargeInfo[0]["PAY_INFO"]
	params["grpId"] = chargeInfo[0]["BOOK_ID"]
	params["accStDay"] = chargeInfo[0]["ACC_ST_DAY"]

	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -매장 충전(SetStoreCharging)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	// 충전 히스토리 insert

	params["searchTy"] = "2"
	params["paymentTy"] = "4"
	params["userTy"] = "1"
	params["payChannel"] = "01"

	UpdateOrderPaidQuery, err := cls.SetUpdateParam(mngAdmin.UpdateOrderPaidCancel, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateOrderPaidCancel parameter fail"))
	}

	_, err = tx.Exec(UpdateOrderPaidQuery)
	dprintf(4, c, "call set Query (%s)\n", UpdateOrderPaidQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateOrderPaidQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	InsertPaymentHistoryQuery, err := cls.SetUpdateParam(mngAdmin.InsertPaymentHistory, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertPaymentHistoryQuery parameter fail"))
	}

	_, err = tx.Exec(InsertPaymentHistoryQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertPaymentHistoryQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPaymentHistoryQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}
