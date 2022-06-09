package mngAdmin

import (
	"biz-web/query/commons"
	"biz-web/query/mngAdmin"
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf

func GetOrderList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	resultCnt, err := cls.GetSelectData(mngAdmin.SelectOrderListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultList, err := cls.GetSelectType(mngAdmin.SelectOrderList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = resultCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]

	if resultList == nil {
		result["resultList"] = []string{}
	} else {
		result["resultList"] = resultList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetOrderInfo(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	orderInfo, err := cls.GetSelectData(mngAdmin.SelectOrderInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult(err.Error(), "DB fail"))
	}
	if orderInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "잘못된 주문 정보 입니다."))
	}

	qrOrderTy := orderInfo[0]["QR_ORDER_TYPE"]

	totalMenu, err := cls.GetSelectType(mngAdmin.SelectOrderDetail, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult(err.Error(), "DB fail"))
	}

	//	userList := make([]map[string]string)

	if qrOrderTy == "2" {

		userSplitList, err := cls.GetSelectType(mngAdmin.SelectOrderUserSplitAmt, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult(err.Error(), "DB fail"))
		}

		order := make(map[string]interface{})
		order["restType"] = orderInfo[0]["REST_TYPE"]
		order["orderNo"] = orderInfo[0]["ORDER_NO"]
		order["restNm"] = orderInfo[0]["REST_NM"]
		order["grpNm"] = orderInfo[0]["GRP_NM"]
		totalAmt, _ := strconv.Atoi(orderInfo[0]["TOTAL_AMT"])
		order["totalAmt"] = totalAmt
		order["orderStat"] = orderInfo[0]["ORDER_STAT"]
		order["orderDate"] = orderInfo[0]["ORDER_DATE"]
		order["orderCancelDate"] = orderInfo[0]["ORDER_CANCEL_DATE"]
		order["totalMenu"] = totalMenu
		order["usersList"] = userSplitList

		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = order

		return c.JSON(http.StatusOK, m)

	} else {

		userDetail, err := cls.GetSelectData(mngAdmin.SelectOrderUserDetail, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult(err.Error(), "DB fail"))
		}

		userList := make([]map[string]interface{}, len(userDetail))
		for i := range userDetail {

			params["userId"] = userDetail[i]["USER_ID"]
			userMenu, err := cls.GetSelectType(mngAdmin.SelectOrderUserMenu, params, c)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			}

			pOrderAmt, _ := strconv.Atoi(userDetail[i]["ORDER_AMT"])
			order2 := make(map[string]interface{})
			order2["userNm"] = userDetail[i]["USER_NM"]
			order2["orderAmt"] = pOrderAmt
			order2["menus"] = userMenu
			order2["memo"] = userDetail[i]["MEMO"]
			userList[i] = order2

		}

		order := make(map[string]interface{})
		order["restType"] = orderInfo[0]["REST_TYPE"]
		order["orderNo"] = orderInfo[0]["ORDER_NO"]
		order["restNm"] = orderInfo[0]["REST_NM"]
		order["grpNm"] = orderInfo[0]["GRP_NM"]
		totalAmt, _ := strconv.Atoi(orderInfo[0]["TOTAL_AMT"])
		order["totalAmt"] = totalAmt
		order["orderStat"] = orderInfo[0]["ORDER_STAT"]
		order["orderDate"] = orderInfo[0]["ORDER_DATE"]
		order["orderCancelDate"] = orderInfo[0]["ORDER_CANCEL_DATE"]
		order["totalMenu"] = totalMenu
		order["usersList"] = userList

		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = order

		return c.JSON(http.StatusOK, m)

	}

}

// 매장 주문취소
func SetOrderCancel(c echo.Context) error {

	dprintf(4, c, "call SetOrderCancel\n")

	params := cls.GetParamJsonMap(c)

	orderInfo, err := cls.GetSelectData(mngAdmin.SelectOrder, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if orderInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "주문 내용이 없습니다."))
	}

	payTy := orderInfo[0]["PAY_TY"]
	orderStat := orderInfo[0]["ORDER_STAT"]
	totalAmt, _ := strconv.Atoi(orderInfo[0]["TOTAL_AMT"])
	bookId := orderInfo[0]["BOOK_ID"]
	storeId := orderInfo[0]["STORE_ID"]
	userId := orderInfo[0]["USER_ID"]
	pointUse, _ := strconv.Atoi(orderInfo[0]["POINT_USE"])

	params["bookId"] = bookId
	params["storeId"] = storeId
	params["userId"] = userId

	if orderStat != "20" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "취소가 불가능한 주문입니다."))
	}

	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -주문 취소(SetOrderCancel)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	// 선불일 경우 금액 환불
	if payTy == "0" {

		linkInfo, err := cls.GetSelectData(mngAdmin.SelectLinkInfo, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		if linkInfo == nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("99", "협약이 내용이 없습니다."))
		}
		prepaidAmt, _ := strconv.Atoi(linkInfo[0]["PREPAID_AMT"])
		linkId := linkInfo[0]["LINK_ID"]

		// 포인트 화불
		params["linkId"] = linkId
		params["prepaidAmt"] = strconv.Itoa(prepaidAmt + totalAmt)
		prepaidPoint, _ := strconv.Atoi(linkInfo[0]["PREPAID_POINT"])
		params["prepaidPoint"] = strconv.Itoa(prepaidPoint + pointUse)

		UpdateLinkQuery, err := cls.SetUpdateParam(mngAdmin.UpdateLink, params)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", err.Error()))
		}

		_, err = tx.Exec(UpdateLinkQuery)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

	}

	//지원금 환불

	params["orderAmt"] = strconv.Itoa(totalAmt)
	UserSupportBalanceUpdateQuery, err := cls.SetUpdateParam(mngAdmin.UpdateBookUserSupportBalance, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", err.Error()))
	}
	_, err = tx.Exec(UserSupportBalanceUpdateQuery)
	if err != nil {

		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//주문 취소
	OrderCancelQuery, err := cls.SetUpdateParam(mngAdmin.UpdateOrderCancel, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "OrderCancel parameter fail"))
	}
	_, err = tx.Exec(OrderCancelQuery)
	if err != nil {

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

func GetOrderListExcel(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	resultList, err := cls.GetSelectType(mngAdmin.SelectOrderListExcel, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	if resultList == nil {
		result["resultList"] = []string{}
	} else {
		result["resultList"] = resultList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func GetTaskList(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	pageSize, _ := strconv.Atoi(params["pageSize"])
	pageNo, _ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	resultCnt, err := cls.GetSelectData(mngAdmin.SelectTaskListCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultList, err := cls.GetSelectType(mngAdmin.SelectTaskList+commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	result := make(map[string]interface{})
	result["totalCount"] = resultCnt[0]["TOTAL_COUNT"]
	result["pageSize"] = params["pageSize"]
	result["pageNo"] = params["pageNo"]

	if resultList == nil {
		result["resultList"] = []string{}
	} else {
		result["resultList"] = resultList
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}
