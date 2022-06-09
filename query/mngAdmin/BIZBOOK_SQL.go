package mngAdmin

var SelectBookList = `SELECT C.USER_NM
							,A.GRP_NM
							,A.GRP_ID
							,A.USE_YN
							,(SELECT COUNT(*) FROM org_agrm_info AS aa WHERE a.grp_id = aa.grp_id AND aa.REQ_STAT='1') AS storeCnt
							,IFNULL((SELECT sum(PREPAID_AMT) FROM org_agrm_info AS aa WHERE a.grp_id = aa.grp_id AND aa.REQ_STAT='1'),0) AS PREPAID_AMT
							,IFNULL((SELECT SUM(TOTAL_AMT) FROM dar_order_info AS AA WHERE A.GRP_ID = AA.GRP_ID AND ORDER_STAT='20' AND PAID_YN='N' AND PAY_TY=1),0) AS UNPAID_AMT
						FROM priv_grp_info AS A
						INNER JOIN priv_grp_user_info AS B ON A.GRP_ID = B.GRP_ID AND B.GRP_AUTH='0'
						INNER JOIN priv_user_info AS C ON B.USER_ID = C.USER_ID AND C.USER_TY='0' AND C.USE_YN='Y'
						WHERE 1=1
`

var SelectBookListCnt = `SELECT COUNT(*) as TOTAL_COUNT
       						FROM priv_grp_info AS A
							INNER JOIN priv_grp_user_info AS B ON A.GRP_ID = B.GRP_ID AND B.GRP_AUTH='0'
							INNER JOIN priv_user_info AS C ON B.USER_ID = C.USER_ID AND C.USER_TY='0' AND C.USE_YN='Y'
							WHERE 1=1
`

var SelectAddStoreList = `SELECT AA.REST_ID
									,AA.REST_NM
									,IFNULL(AA.BUSID,'') AS BUSID
									,AA.ADDR
							FROM priv_rest_info AS AA
							WHERE AA.USE_YN='Y'
								AND AA.REST_ID NOT IN (SELECT REST_ID
														FROM org_agrm_info AS A
														WHERE 
															GRP_ID = '#{searchGrpId}')
	    						AND AA.REST_NM LIKE '%#{searchKeyword}%'
`

//BookCs초기화면
var SelectBookInfo = `SELECT A.GRP_NM as grpNm
								,A.GRP_TYPE_CD as grpTypeCd
								,CASE  WHEN B.CAT_CD IS NULL THEN 'N' ELSE 'Y' END AS catTy 
								,A.DETAIL_VIEW_YN as detailYn
								,IFNULL(A.LIMIT_YN,'N') as limitYn
								,IFNULL(A.LIMIT_AMT,0) AS limitAmt
								,IFNULL(A.LIMIT_DAY_AMT ,0) AS limitDayAmt
								,IFNULL(A.INVITE_LINK,'') as inviteLink
								,IFNULL(A.SUPPORT_AMT ,0) AS supportAmt
								,A.support_yn as supportYn 
								,A.intro as intro
								,D.company_nm as companyNm 
								,C.BOOK_ID as bookId
							FROM priv_grp_info AS A
							LEFT OUTER JOIN priv_catering_grp_info AS B ON A.GRP_ID = B.GRP_ID
							LEFT OUTER JOIN b_company_book AS C ON A.GRP_ID = C.BOOK_ID
							LEFT OUTER JOIN b_company AS D ON C.company_id = D.company_id
							WHERE 
								A.grp_id = '#{searchGrpId}'
`

var SelectBookLinkStore = `SELECT A.REST_ID
									,B.REST_NM
									,DATE_FORMAT(IFNULL(A.REQ_DATE,''),'%Y-%m-%d') AS JOIN_DATE
									,IFNULL(A.PREPAID_AMT,0) AS PREPAID_AMT
								FROM org_agrm_info AS A
								INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID AND A.REQ_STAT='1'
								WHERE 
									A.GRP_ID = '#{searchGrpId}'
								ORDER BY A.REQ_DATE DESC`

var SelectBookLinkStoreCnt = `SELECT count(*) as total
								FROM org_agrm_info AS A
								INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID AND A.REQ_STAT='1'
								WHERE 
									A.GRP_ID = '#{searchGrpId}'
								ORDER BY A.REQ_DATE DESC`

var SelectCompanyGrpUserList = `
								SELECT B.USER_NM
								,B.USER_ID
								,IFNULL(A.SUPPORT_BALANCE,0) AS amt
								,DATE_FORMAT(IFNULL(A.AUTH_DATE,''),'%Y-%m-%d') AS JOIN_DATE
								,A.AUTH_STAT
							FROM priv_grp_user_info AS A
							INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID 
							WHERE 
								A.GRP_ID = '#{searchGrpId}'
							ORDER BY A.AUTH_STAT ASC, A.AUTH_DATE DESC`

var SelectPrivateGrpUserList = `
SELECT B.USER_NM
	,B.USER_ID
	,IFNULL(sum(C.ORDER_AMT),0) as amt
	,DATE_FORMAT(IFNULL(A.AUTH_DATE,''),'%Y-%m-%d') AS JOIN_DATE
	,A.AUTH_STAT
FROM priv_grp_user_info AS A
INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID 
INNER JOIN dar_order_detail AS C ON C.USER_ID = A.USER_ID
INNER JOIN dar_order_info AS D ON D.ORDER_STAT = '20' AND C.ORDER_NO = D.ORDER_NO AND D.GRP_ID = A.GRP_ID AND D.ORDER_TY !='4'
WHERE 
	A.GRP_ID = '#{searchGrpId}'
AND C.ORDER_DATE >= DATE_FORMAT(now(), '%Y%m')
GROUP BY A.USER_ID
ORDER BY A.AUTH_STAT ASC, A.AUTH_DATE DESC
`

var SelectBookUserListCnt = `SELECT count(*)
							FROM priv_grp_user_info AS A
							INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID 
							WHERE 
								A.GRP_ID = '#{searchGrpId}'
							ORDER BY A.AUTH_STAT ASC, A.AUTH_DATE DESC`

var SelectBookData = `SELECT A.GRP_NM
								, A.GRP_ID
								, A.GRP_TYPE_CD
								, A.USE_YN
								, A.INVITE_LINK
								, A.INTRO
								, A.SUPPORT_YN                                # 지원금사용여부
								, IFNULL(A.SUPPORT_AMT, 0)   AS SUPPORT_AMT   # 지원금 기준금액
								, A.SUPPORT_EXCEED_YN                         # 지원금 초과사용 가능여부
								, A.SUPPORT_FORWARD_YN                        # 지원금 이월여부
								, A.LIMIT_YN                                  # 사용제한설정
								, IFNULL(A.LIMIT_AMT, 0)     AS LIMIT_AMT     # 1회 사용금액 한도
								, IFNULL(A.LIMIT_DAY_AMT, 0) AS LIMIT_DAY_AMT # 1일 사용금액 한도
								, IFNULL(A.LIMIT_DAY_CNT, 0) AS LIMIT_DAY_CNT # 1일 사용횟수 한도
							FROM priv_grp_info AS A
							WHERE 
								A.grp_id = '#{grpId}'
`

var UpdateBookData = `
UPDATE priv_grp_info
	SET	
		GRP_NM = '#{grpNm}'
		,GRP_TYPE_CD = '#{grpTy}'
		,USE_YN= '#{useYn}'
		,SUPPORT_YN = '#{suppotYn}'
		,SUPPORT_AMT = '#{suppotAmt}'
		,SUPPORT_EXCEED_YN = '#{supportExceedYn}'
		,SUPPORT_FORWARD_YN = '#{supportForwardYn}'
		,LIMIT_YN = '#{limitYn}'
		,LIMIT_AMT = '#{limitAmt}'
		,LIMIT_DAY_AMT = '#{limitDayAmt}'
		,LIMIT_DAY_CNT = '#{limitDayCnt}'
		,intro = IF('#{intros}' = '',NULL, '#{intros}' )
	WHERE 
		grp_id = '#{grpId}'

`
