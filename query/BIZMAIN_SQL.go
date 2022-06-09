package query

// 전체사용자수
var AllUserCnt = `	SELECT	COUNT(*) AS TOT_CNT, 
							SUM(IF (DATE_FORMAT(JOIN_DATE, '%Y%m%d') BETWEEN DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 7 DAY), '%Y%m%d') AND DATE_FORMAT(NOW(), '%Y%m%d'), 1, 0)) AS WEEK_CNT
				FROM PRIV_USER_INFO
				where use_yn = 'y'	
				`

//전체 기업수
var AllCompanyCnt = `SELECT	COUNT(*) AS TOT_CNT, 
							SUM(IF (DATE_FORMAT(REG_DATE, '%Y%m%d') BETWEEN DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 7 DAY), '%Y%m%d') AND DATE_FORMAT(NOW(), '%Y%m%d'), 1, 0)) AS WEEK_CNT
						FROM PRIV_GRP_INFO
						where use_yn = 'y'						
`

//전체스토어수
var AllStoreCnt = `SELECT	COUNT(*) AS TOT_CNT, 
							SUM(IF (DATE_FORMAT(REG_DATE, '%Y%m%d') BETWEEN DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 7 DAY), '%Y%m%d') AND DATE_FORMAT(NOW(), '%Y%m%d'), 1, 0)) AS WEEK_CNT
					FROM PRIV_REST_INFO
					WHERE 
						TEST_YN='N'
						and use_yn='y'`

// 금월누적금액/건수, 금일사용금액/건수
var CurrentMonthAmtCnt = `SELECT 	DATE_FORMAT(NOW(), '%Y%m%d') AS ORDER_DAY, 
							IFNULL(SUM(CASE WHEN A.ORDER_DATE = DATE_FORMAT(NOW(), '%Y%m%d') THEN A.ORDER_AMT END ),0) AS DAY_AMT, 
							COUNT(CASE WHEN A.ORDER_DATE = DATE_FORMAT(NOW(), '%Y%m%d') THEN A.ORDER_CNT END ) AS DAY_CNT, 
							DATE_FORMAT(NOW(), '%Y%m') AS ORDER_MONTH, 
							IFNULL(SUM(CASE WHEN DATE_FORMAT(A.ORDER_DATE, '%Y%m') = DATE_FORMAT(NOW(), '%Y%m') THEN A.ORDER_AMT END ),0) AS MONTH_AMT, 
							COUNT(CASE WHEN DATE_FORMAT(A.ORDER_DATE, '%Y%m') = DATE_FORMAT(NOW(), '%Y%m') THEN A.ORDER_CNT END ) AS MONTH_CNT
					FROM
							(SELECT DATE_FORMAT(A.ORDER_DATE, '%Y%m%d') AS ORDER_DATE, 
									A.CREDIT_AMT AS ORDER_AMT,  
									A.USER_ID AS ORDER_CNT
								FROM
									DAR_ORDER_INFO A,
									PRIV_USER_INFO B,
									PRIV_GRP_USER_INFO C
								WHERE
									1 = 1
									AND A.USER_ID = B.USER_ID
									AND B.USER_ID = C.USER_ID
									AND A.GRP_ID = C.GRP_ID
									AND DATE_FORMAT(A.ORDER_DATE, '%Y%m') = DATE_FORMAT(NOW(), '%Y%m')
									AND A.ORDER_STAT = '20'
							) A`

// 전월누적금액/건수, 전일사용금액/건수
var PreviousMonthAmtCnt = `SELECT 	DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), '%Y%m%d') AS ORDER_DAY, 
									IFNULL(SUM(CASE WHEN A.ORDER_DATE = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), '%Y%m%d') THEN A.ORDER_AMT END ),0) AS DAY_AMT, 
									COUNT(CASE WHEN A.ORDER_DATE = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), '%Y%m%d') THEN A.ORDER_CNT END ) AS DAY_CNT, 
									DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH), '%Y%m') AS ORDER_MONTH, 
									IFNULL(SUM(CASE WHEN DATE_FORMAT(A.ORDER_DATE, '%Y%m') = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH), '%Y%m') THEN A.ORDER_AMT END ),0) AS MONTH_AMT, 
									COUNT(CASE WHEN DATE_FORMAT(A.ORDER_DATE, '%Y%m') = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH), '%Y%m') THEN A.ORDER_CNT END ) AS MONTH_CNT
							FROM 	
								(SELECT 	DATE_FORMAT(A.ORDER_DATE, '%Y%m%d') AS ORDER_DATE, 
											A.CREDIT_AMT AS ORDER_AMT,  
											A.USER_ID AS ORDER_CNT
									FROM
										DAR_ORDER_INFO A,
										PRIV_USER_INFO B,
										PRIV_GRP_USER_INFO C
									WHERE
										1 = 1
										AND A.USER_ID = B.USER_ID
										AND B.USER_ID = C.USER_ID
										AND A.GRP_ID = C.GRP_ID
										AND DATE_FORMAT(A.ORDER_DATE, '%Y%m') = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH), '%Y%m')
										AND A.ORDER_STAT = '20'
								) A`

//주석 되어있는 부분 본 쿼리 데이터 삽입
// 시간별 사용현황(금일)
var TimeUseToday = `SELECT 	A.ORDER_HOUR, 
						IFNULL(B.CNT, 0) AS HOUR_CNT, 
						IFNULL(B.AMT, 0) AS HOUR_AMT
				FROM
						(
							SELECT '01' AS ORDER_HOUR UNION ALL
							SELECT '02' UNION ALL
							SELECT '03' UNION ALL
							SELECT '04' UNION ALL
							SELECT '05' UNION ALL
							SELECT '06' UNION ALL
							SELECT '07' UNION ALL
							SELECT '08' UNION ALL
							SELECT '09' UNION ALL
							SELECT '10' UNION ALL
							SELECT '11' UNION ALL
							SELECT '12' UNION ALL
							SELECT '13' UNION ALL
							SELECT '14' UNION ALL
							SELECT '15' UNION ALL
							SELECT '16' UNION ALL
							SELECT '17' UNION ALL
							SELECT '18' UNION ALL
							SELECT '19' UNION ALL
							SELECT '20' UNION ALL
							SELECT '21' UNION ALL
							SELECT '22' UNION ALL
							SELECT '23' UNION ALL
							SELECT '24'
						) A LEFT OUTER JOIN
							(
								SELECT DATE_FORMAT(ORDER_DATE, '%H') AS ORDER_HOUR, COUNT(ORDER_NO) AS CNT, SUM(CREDIT_AMT) AS AMT
								FROM DAR_ORDER_INFO
								WHERE 1 = 1

								AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') = DATE_FORMAT(NOW(), '%Y%m%d')

								AND ORDER_STAT = '20'
								GROUP BY DATE_FORMAT(ORDER_DATE, '%H')
							) B
						ON A.ORDER_HOUR = B.ORDER_HOUR`

//AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') = DATE_FORMAT(NOW(), '%Y%m%d')
// 시간별 사용현황(전일)

var TimeUserYesterday = `SELECT A.ORDER_HOUR, IFNULL(B.CNT, 0) AS HOUR_CNT, IFNULL(B.AMT, 0) AS HOUR_AMT
		FROM
		(
			SELECT '01' AS ORDER_HOUR UNION ALL
			SELECT '02' UNION ALL
			SELECT '03' UNION ALL
			SELECT '04' UNION ALL
			SELECT '05' UNION ALL
			SELECT '06' UNION ALL
			SELECT '07' UNION ALL
			SELECT '08' UNION ALL
			SELECT '09' UNION ALL
			SELECT '10' UNION ALL
			SELECT '11' UNION ALL
			SELECT '12' UNION ALL
			SELECT '13' UNION ALL
			SELECT '14' UNION ALL
			SELECT '15' UNION ALL
			SELECT '16' UNION ALL
			SELECT '17' UNION ALL
			SELECT '18' UNION ALL
			SELECT '19' UNION ALL
			SELECT '20' UNION ALL
			SELECT '21' UNION ALL
			SELECT '22' UNION ALL
			SELECT '23' UNION ALL
			SELECT '24'
		) A LEFT OUTER JOIN
		(
			SELECT DATE_FORMAT(ORDER_DATE, '%H') AS ORDER_HOUR, COUNT(ORDER_NO) AS CNT, SUM(CREDIT_AMT) AS AMT
			FROM DAR_ORDER_INFO
			WHERE 1 = 1

			AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), '%Y%m%d')
							
			AND ORDER_STAT = '20'
			GROUP BY DATE_FORMAT(ORDER_DATE, '%H')
		) B
		ON A.ORDER_HOUR = B.ORDER_HOUR`

//AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), '%Y%m%d')

// 가맹점 매출액 순위
var StoreAmtRank = `SELECT REST_ID, FN_GET_RESTNAME(REST_ID) AS NM, SUM(CREDIT_AMT) AS ORDER_AMT
FROM DAR_ORDER_INFO
WHERE 1 = 1
AND ORDER_STAT = '20'
AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') = DATE_FORMAT(NOW(), '%Y%m%d')
GROUP BY REST_ID
ORDER BY ORDER_AMT DESC
LIMIT 5`

//AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') = DATE_FORMAT(NOW(), '%Y%m%d')

// 기업 구매액 순위
var CompanyAmtRank = `SELECT GRP_ID, FN_GET_GRPNAME(GRP_ID) AS NM, SUM(CREDIT_AMT) AS ORDER_AMT
		FROM DAR_ORDER_INFO
		WHERE 1 = 1
		AND ORDER_STAT = '20'
		AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') = DATE_FORMAT(NOW(), '%Y%m%d')
		GROUP BY GRP_ID
		ORDER BY ORDER_AMT DESC
		LIMIT 5`

//AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') = DATE_FORMAT(NOW(), '%Y%m%d')

func getKeyword(num string) []string {
	switch num {
	case "3": //오늘 INTERVAL 1 HOUR
		return []string{"%H", "HOUR", "8", "10"}
	case "2": //이번달 INTERVAL 1 DAY
		return []string{"%d", "DAY", "6", "8"}
	case "1": //올해 INTERVAL 1 MONTH
		return []string{"%m", "MONTH", "4", "6"}
	default:
		return []string{}
	}
}

func HomeChartData(key string) string {
	arr := getKeyword(key)
	query := `SELECT A.ORDER_S, IFNULL(B.CNT, 0) AS CNT, IFNULL(B.AMT, 0) AS AMT
		FROM
		(
			SELECT '01' AS ORDER_S UNION ALL
			SELECT '02' UNION ALL
			SELECT '03' UNION ALL
			SELECT '04' UNION ALL
			SELECT '05' UNION ALL
			SELECT '06' UNION ALL
			SELECT '07' UNION ALL
			SELECT '08' UNION ALL
			SELECT '09' UNION ALL
			SELECT '10' UNION ALL
			SELECT '11' UNION ALL
			SELECT '12' UNION ALL
			SELECT '13' UNION ALL
			SELECT '14' UNION ALL
			SELECT '15' UNION ALL
			SELECT '16' UNION ALL
			SELECT '17' UNION ALL
			SELECT '18' UNION ALL
			SELECT '19' UNION ALL
			SELECT '20' UNION ALL
			SELECT '21' UNION ALL
			SELECT '22' UNION ALL
			SELECT '23' UNION ALL
			SELECT '24' UNION ALL
			SELECT '25' UNION ALL
			SELECT '26' UNION ALL
			SELECT '27' UNION ALL
			SELECT '28' UNION ALL
			SELECT '29' UNION ALL
			SELECT '30' UNION ALL
			SELECT '31'
		) A LEFT OUTER JOIN
		(
			SELECT DATE_FORMAT(ORDER_DATE, '` + arr[0] + `') AS ORDER_S, COUNT(ORDER_NO) AS CNT, SUM(CREDIT_AMT) AS AMT
			FROM DAR_ORDER_INFO as c
			Inner join b_company_book as d on d.BOOK_ID = c.GRP_ID 
			WHERE 
			1 = 1
			AND c.ORDER_STAT = '20'
			and d.company_id = '#{companyId}'
			AND LEFT(DATE_FORMAT(c.ORDER_DATE, '%Y%m%d%H'),` + arr[2] + `) = LEFT(DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 ` + arr[1] + `), '%Y%m%d%H'),` + arr[2] + `) 
			GROUP BY LEFT(DATE_FORMAT(c.ORDER_DATE, '%Y%m%d%H'),` + arr[3] + `)
		) B
		ON A.ORDER_S = B.ORDER_S`

	return query
}

var StoreUseRank = `select a.REST_ID,
							(select REST_NM from priv_rest_info as c WHERE a.REST_ID = c.REST_ID) as REST_NM,
							SUM(a.CREDIT_AMT) as TOTAL_AMT
					from dar_order_info as a
					inner join b_company_book as b on b.BOOK_ID = a.GRP_ID
					WHERE 
						1=1
						and b.company_id = '#{companyId}'
						and a.ORDER_STAT = '20'
						and LEFT(a.ORDER_DATE,'#{YMD}') = LEFT(DATE_FORMAT(now(), '%Y%m%d'),'#{YMD}')
						GROUP BY REST_ID asc
						ORDER BY TOTAL_AMT desc
						limit 5`

// 년 = 4, 월 = 6, 오늘 = 8

var GrpCount = `SELECT COUNT(*) AS Count FROM priv_grp_info AS a
				INNER JOIN b_company_book AS b ON a.GRP_ID = b.BOOK_ID
				WHERE b.company_id = '#{companyId}'
				AND a.USE_YN = 'Y'`

var GrpUserCount = `SELECT COUNT(DISTINCT user_id) as Count
					FROM b_company_user 
					WHERE 
						company_id = '#{companyId}'
						and USE_YN = 'Y'`

var GrpJoinStore = `SELECT 
						COUNT(a.REST_ID) as Count
					FROM org_agrm_info as a 
					inner join b_company_book as b on a.GRP_ID = b.BOOK_ID 
					WHERE
						b.company_id = '#{companyId}'`

var UseAmt = `select 
	IFNULL(SUM(CASE WHEN LEFT(C.ORDER_DATE,6) = LEFT(DATE_FORMAT(NOW(), '%Y%m%d'),6) THEN C.CREDIT_AMT END ),0) AS monthAmt,
	IFNULL(SUM(CASE WHEN LEFT(C.ORDER_DATE,8) = LEFT(DATE_FORMAT(NOW(), '%Y%m%d'),8) THEN C.CREDIT_AMT END ),0) AS dayAmt
	from							
	(select *
		from b_company_book as a 
		inner join DAR_ORDER_INFO as b on a.BOOK_ID = b.GRP_ID 
		WHERE 
		a.company_id = '#{companyId}' 
		and b.ORDER_STAT = '20' 
	) C`
