package mngAdmin

/*
	최초 작성일자 : 21.06.29
	작상자 : 김형곤
*/

// SelectUserList Select
var SelectUserList = `SELECT A.USER_ID as uId
							,A.USER_NM
							,A.LOGIN_ID
							,A.USER_TY
							,IFNULL(C.rest_nm,'') as restNm
							,A.HP_NO
							,A.USE_YN
							,CASE 	WHEN A.KAKAO_KEY IS NOT NULL THEN 'KAKAO' 
									WHEN A.APPLE_KEY IS NOT NULL THEN 'APPLE' 
									WHEN A.NAVER_KEY IS NOT NULL THEN 'NAVER' 
									ELSE 'ID' END AS LOGIN_TY 
							,IFNULL(A.EMAIL,'') AS EMAIL
							,DATE_FORMAT(A.JOIN_DATE,'%Y-%m-%d') AS JOIN_DATE
						FROM priv_user_info AS A
						LEFT OUTER JOIN priv_rest_user_info as B ON A.user_id = B.user_id 
						LEFT OUTER JOIN priv_rest_info as C ON B.rest_id = C.rest_id
						WHERE 1=1
						AND (('#{searchKey}' = 'userNm' AND A.USER_NM LIKE '%#{searchKeyword}%') 
						OR ('#{searchKey}' = 'userHp' AND A.HP_NO LIKE '%#{searchKeyword}%'))
						`

// SelectUserListCnt
var SelectUserListCnt = `SELECT COUNT(*) as TOTAL_COUNT
       					FROM priv_user_info AS A
       					WHERE 1=1
						AND (('#{searchKey}' = 'userNm' AND A.USER_NM LIKE '%#{searchKeyword}%') 
						OR ('#{searchKey}' = 'userHp' AND A.HP_NO LIKE '%#{searchKeyword}%'))
						`

// SelectUserInfo searchKeyword string
var SelectUserInfo = `SELECT USER_ID as uId
							,USER_NM
							,HP_NO
							,EMAIL
							,USE_YN
							,LOGIN_ID
							,DATE_FORMAT(IFNULL(A.JOIN_DATE,''),'%Y-%m-%d') AS JOIN_DATE
							,(SELECT MAX(ADDED_DATE) FROM sys_log_access AS AA WHERE A.LOGIN_ID = AA.USER_ID AND LOG_IN_OUT='I') AS LAST_JOIN_DATE
							,CASE WHEN APPLE_KEY IS NOT NULL THEN 'APPLE' 
								WHEN NAVER_KEY IS NOT NULL THEN 'NAVER' 
								WHEN KAKAO_KEY IS NOT NULL THEN 'KAKAO' 
							ELSE 'ID' END AS JOIN_TYPE 
						FROM PRIV_USER_INFO AS A
							WHERE 
							USER_ID = '#{userId}' `

//#{} , ${} 으로 들어오는 데이터는 varchar 형으로 지정해주어야한다 (varchar = 가변길이 문자열 자료형)
//주의 웨어 다음에 #{} 또는 ${} 들어가면 오류 반드시 한줄 내려서 적어 줄것

// DefaultUpdateUserInfo updateUserInfo 기본 (byJava)
//var DefaultUpdateUserInfo = `UPDATE priv_user_info SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
//								,LOGIN_PW = '#{ passwd }'
//								,USER_NM = '#{ userNm }'
//								,EMAIL = '#{ email }'
//								,HP_NO = '#{ hpNo }'
//								,USE_YN = '#{ useYn }'
//			WHERE
//			USER_ID = '#{ userId }'`

// UpdateUserInfo Update
var UpdateUserInfo = `UPDATE priv_user_info 
						SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								,EMAIL = '#{email}'
								,HP_NO = '#{hpNo}'
								,USE_YN = '#{useYn}'
			WHERE 
			USER_ID = '#{userId}'`

// UpdateUserPassword Update
// 비밀번호 초기화
var UpdateUserPassword = `UPDATE priv_user_info 
							SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								,LOGIN_PW = '#{passwd}'
			WHERE 
			USER_ID = '#{userId}'`

// SelectUserGrpList
// notice :
var SelectUserGrpList = `SELECT B.GRP_NM
								,B.GRP_ID
								,A.GRP_AUTH
								,IFNULL(A.SUPPORT_BALANCE,0) AS PREPAID_AMT
								,DATE_FORMAT(IFNULL(A.REG_DATE,''),'%Y-%m-%d') AS REG_DATE
								,(select count(*) from priv_grp_user_info as C where C.grp_id = B.GRP_ID and C.LEAVE_TY IS NULL) as cnt
							FROM priv_grp_user_info AS A
							INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID
							WHERE 
							A.USER_ID = '#{userId}'
							AND A.GRP_AUTH = '0'
							ORDER BY A.REG_DATE DESC `

var SelectUserGrpList2 = `SELECT B.GRP_NM
								,B.GRP_ID
								,A.GRP_AUTH
								,IFNULL(A.SUPPORT_BALANCE,0) AS PREPAID_AMT
								,DATE_FORMAT(IFNULL(A.REG_DATE,''),'%Y-%m-%d') AS REG_DATE
							FROM priv_grp_user_info AS A
							INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID
							WHERE 
							A.USER_ID = '#{userId}'
							AND A.GRP_AUTH = '1'
							AND A.LEAVE_TY IS NULL
							ORDER BY A.REG_DATE DESC `

var SelectUnUsedCouponList = `
SELECT 	
		A.ORDER_NO AS orderNo
		, D.REST_NM as restNm
		, C.ITEM_NM as menuNm
		, C.ITEM_PRICE as price
		, A.CPNO as cpNo
		, DATE_FORMAT(A.EXCH_TO_DY, '%Y-%m-%d') as useDay
		, DATE_FORMAT(A.EXPIRE_DATE, '%Y-%m-%d') as darUseDay
FROM dar_order_coupon as A
INNER JOIN dar_order_info AS B ON A.ORDER_NO = B.ORDER_NO
INNER JOIN dar_sale_item_info AS C ON A.PROD_ID = C.PROD_ID
INNER JOIN priv_rest_info AS D ON C.REST_ID = D.REST_ID
WHERE A.CP_STATUS = '0'
AND B.USER_ID = '#{userId}'`

var UpdateCouponUsedDate = `
UPDATE dar_order_coupon
SET 
	EXPIRE_DATE = DATE_FORMAT(now(), '%Y%m%d')
WHERE 
	ORDER_NO = '#{orderNo}'
`

var UpdateUserDelCompany = `UPDATE b_company_user 
					SET 
					USE_YN = '#{useYn}' 
					,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
					WHERE 
					USER_ID = '#{userId}'`

var UpdateUserDelGrp = `update priv_grp_user_info 
					SET 
					AUTH_STAT = '#{authStat}'
					,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
					WHERE 
					USER_ID = '#{userId}'`

var SelectUserLoginId = `
SELECT login_id as loginId
FROM priv_user_info 
WHERE
user_id = '#{userId}'
`

var UpdateSysUserPassword = `
UPDATE sys_user_info 
SET 
  user_pass = '#{passwd}'
, modify_date = DATE_FORMAT(NOW(), '%Y-%m-%d %H:%i:%s') 
WHERE 
user_id = '#{loginId}'
`



var SelectTpayBillingCardList string = `SELECT USER_ID
										, SEQ
										, CARD_NAME
										, CARD_CODE
										, CARD_NUM
										, CARD_TYPE
										, DATE_FORMAT(IFNULL(REG_DATE,''),'%Y-%m-%d') AS REG_DATE
									FROM b_tpay_billing_key
									WHERE
									USER_ID='#{userId}'
									AND USE_YN='Y'
										`
