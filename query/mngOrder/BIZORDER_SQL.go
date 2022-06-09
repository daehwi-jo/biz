package mngOrder

// SelectGrpCode
// notice : bizBookMngQuery와 동일
var SelectGrpCode = `SELECT A.GRP_ID as value
		,A.GRP_NM as label
		,A.GRP_TYPE_CD
	FROM priv_grp_info as a
	inner join b_company_book as b on a.grp_id = b.book_id
	INNER JOIN priv_grp_user_info as c on a.grp_id = c.grp_id and c.GRP_AUTH = 0
	where a.use_yn='Y'
		AND B.BOOK_ID IN (
							SELECT  GRP_ID
							FROM priv_grp_user_info as aa
							WHERE aa.GRP_AUTH=0 )
		AND company_id = '#{companyId}'
		AND c.user_id LIKE '%#{userId}%'
`

var SelectOrderList = `SELECT A.REST_NM
								,(SELECT COUNT(*) FROM DAR_ORDER_DETAIL AS AA WHERE E.ORDER_NO = AA.ORDER_NO) AS MENU_CNT
								,E.ORDER_NO
								,E.TOTAL_AMT
								,DATE_FORMAT(E.ORDER_DATE,'%Y-%m-%d') AS ORDER_DATE
								,DATE_FORMAT(E.ORDER_DATE,'%H:%i:%S') AS ORDER_TIME
								,C.GRP_NM
								,IFNULL(E.ORDER_TY,0) AS ORDER_TYPE
								,F.USER_NM
								,G.COURSE_NM
						FROM PRIV_REST_INFO AS A
						INNER JOIN ORG_AGRM_INFO AS B ON A.REST_ID = B.REST_ID AND B.REQ_STAT='1'
						INNER JOIN PRIV_GRP_INFO AS C ON B.GRP_ID = C.GRP_ID
						INNER JOIN b_company_book AS D ON C.GRP_ID = D.BOOK_id
						INNER JOIN DAR_ORDER_INFO AS E ON A.REST_ID = E.REST_ID AND C.GRP_ID = E.GRP_ID
						INNER JOIN PRIV_USER_INFO AS F ON E.USER_ID = F.USER_ID
						LEFT OUTER JOIN (SELECT AA.GRP_ID , AA.COURSE_NM , BB.USER_ID , AA.COURSE_ID
											FROM priv_course_info AS AA
											INNER JOIN priv_course_user_info AS BB ON AA.COURSE_ID = BB.COURSE_ID)  G ON G.USER_ID = E.USER_ID AND G.GRP_ID = E.GRP_ID and G.course_id=E.TABLE_NO
						WHERE 
							D.COMPANY_ID = '#{companyId}' 
							AND E.ORDER_STAT=20
							AND DATE_FORMAT(E.ORDER_DATE,'%Y-%m-%d') 
							BETWEEN DATE_FORMAT('#{searchStartDt}','%Y-%m-%d') 
							AND DATE_FORMAT('#{searchEndDt}','%Y-%m-%d')
							AND B.GRP_ID = '#{searchGrpId}'
				 			AND G.COURSE_ID = '#{searchCourseId}'
							AND F.USER_NM LIKE '%#{searchUserNm}%'
						ORDER BY DATE_FORMAT(E.ORDER_DATE,'%Y-%m-%d') DESC,DATE_FORMAT(E.ORDER_DATE,'%H:%i:%S')  DESC`

var SelectOrderListCount = `SELECT COUNT(*) as TOTAL_COUNT, sum(E.TOTAL_AMT) as ALL_AMT
        					FROM PRIV_REST_INFO AS A
							INNER JOIN ORG_AGRM_INFO AS B ON A.REST_ID = B.REST_ID AND B.REQ_STAT='1'
							INNER JOIN PRIV_GRP_INFO AS C ON B.GRP_ID = C.GRP_ID
							INNER JOIN b_company_book AS D ON C.GRP_ID = D.BOOK_id
							INNER JOIN DAR_ORDER_INFO AS E ON A.REST_ID = E.REST_ID AND C.GRP_ID = E.GRP_ID
							INNER JOIN PRIV_USER_INFO AS F ON E.USER_ID = F.USER_ID
							WHERE 
								D.COMPANY_ID='#{companyId}'
								AND E.ORDER_STAT=20
								AND DATE_FORMAT(E.ORDER_DATE,'%Y-%m-%d') BETWEEN DATE_FORMAT('#{searchStartDt}','%Y-%m-%d') 
								AND DATE_FORMAT('#{searchEndDt}','%Y-%m-%d')
								AND G.COURSE_ID = '#{searchCourseId}'
								AND B.GRP_ID = '#{searchGrpId}'
`

func SelectOrderListCountWhere(params map[string]string) string {
	query := ``
	if params["searchGrpId"] != "" {
		query += `AND B.GRP_ID = '#{searchGrpId}'`
	}
	if params["searchUserNm"] != "" {
		query += `AND F.USER_NM LIKE '%#{searchUserNm}%'`
	}
	return query
}

//AND G.COURSE_ID = '%#{searchCourseId}%' 제외함
//필수조건으로 필요한 값들 searchStartDt searchEndDt searchGrpId searchCourseId searchUserNm

var SelectOrderExcelList = `SELECT DATE_FORMAT(E.ORDER_DATE,'%Y-%m-%d') AS EXCEL1
					,DATE_FORMAT(E.ORDER_DATE,'%H:%i:%S') AS EXCEL2
					,IFNULL(F.USER_NM,'') as EXCEL3
					,IFNULL(C.GRP_NM,'') as EXCEL4
					,IFNULL(A.REST_NM,'') as EXCEL5
					,IFNULL(G.ITEM_NM,'') as EXCEL6
					,(SELECT COUNT(*) FROM DAR_ORDER_DETAIL AS AA WHERE E.ORDER_NO = AA.ORDER_NO) AS EXCEL7
					,IFNULL(E.TOTAL_AMT,0) as EXCEL8
					,CASE WHEN E.PAY_TY='0' THEN  '선불' ELSE '후불' END AS EXCEL9
					,CASE WHEN E.PAY_TY='1' THEN  PAID_YN ELSE '' END AS EXCEL10
					,'' AS  EXCEL11
			FROM PRIV_REST_INFO AS A
			INNER JOIN ORG_AGRM_INFO AS B ON A.REST_ID = B.REST_ID AND B.REQ_STAT='1'
			INNER JOIN PRIV_GRP_INFO AS C ON B.GRP_ID = C.GRP_ID
			INNER JOIN b_company_book AS D ON C.GRP_ID = D.BOOK_id
			INNER JOIN DAR_ORDER_INFO AS E ON A.REST_ID = E.REST_ID AND C.GRP_ID = E.GRP_ID
			INNER JOIN PRIV_USER_INFO AS F ON E.USER_ID = F.USER_ID
			INNER JOIN ( SELECT  AA.ORDER_NO
										,MAX(BB.ITEM_NM) AS ITEM_NM
							FROM DAR_ORDER_DETAIL AS AA
							INNER JOIN DAR_SALE_ITEM_INFO AS BB ON AA.ITEM_NO = BB.ITEM_NO
							GROUP BY AA.ORDER_NO ) AS G ON E.ORDER_NO = G.ORDER_NO
			WHERE 
				D.COMPANY_ID='#{companyId}' 
				AND E.ORDER_STAT=20
				AND DATE_FORMAT(E.ORDER_DATE,'%Y-%m-%d') 
					BETWEEN DATE_FORMAT('#{searchStartDt}','%Y-%m-%d') 
				AND DATE_FORMAT('#{searchEndDt}','%Y-%m-%d')
				AND B.GRP_ID = '#{searchGrpId}'
				AND F.USER_NM LIKE '%#{searchUserNm}%' 
			ORDER BY DATE_FORMAT(E.ORDER_DATE,'%Y-%m-%d') DESC,DATE_FORMAT(E.ORDER_DATE,'%H:%i:%S')  DESC`
