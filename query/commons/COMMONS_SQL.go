package commons

var PagingQuery string = `
					LIMIT #{pageSize}
					OFFSET #{offSet}`

var SelectBoardList string = `SELECT
							BOARD_ID
							,TITLE
							,BOARD_TYPE
							,LINK_URL
							,DATE_FORMAT(REG_DATE, '%Y.%m.%d') AS regDate	
							FROM sys_boards
							WHERE START_DATE <= NOW() AND END_DATE >=NOW()
							`

var SelectCategoryList string = `SELECT CATEGORY_ID
									,CATEGORY_NM 
								FROM b_category
								WHERE 
								CATEGORY_GRP_CODE='#{grpCode}'
								AND USE_YN='Y'
							`
var SelectCodeist string = `SELECT CODE_ID
								,CODE_NM
								FROM b_code
								WHERE
								CATEGORY_ID='#{categoryId}'
								AND USE_YN='Y'
							`

// 버전
var SelectVersion string = `SELECT 
							  VERSION AS versionCode 
							, CASE WHEN AUTO_YN ='Y' THEN 'true' ELSE 'false' END AS  isRequireUpdate
							FROM sys_version_info
							WHERE 
							USE_YN = 'Y' 
							AND OS_TY = '#{osTy}'
							AND APP_TY = '#{appTy}'
							ORDER BY VERSION_ID DESC
							LIMIT 1
							`

var InsertPushLog string = `INSERT INTO sys_log_push
						(
						NAME
						, TITLE
						, BODY
						, APP_TY
						, REG_DATE
						, REG_ID
						)
						VALUES
						(
                          '#{name}'
						, '#{title}'
						, '#{body}'
						, '#{appTy}'
						, NOW()
						, '#{regId}'
						)
						`

var SelectCouponInfo string = `SELECT COUPON_NAME
											,USE_TYPE
											,COUPON_VAL
											,USE_SERVICE
											,ITEM_CODE
											,USE_YN
										FROM b_coupon AS a
										WHERE 
										START_DATE <= DATE_FORMAT(NOW(), '%Y%m%d')
										AND END_DATE  >= DATE_FORMAT(NOW(), '%Y%m%d')
										AND  COUPON_NO = '#{couponNo}'
                             		`

var SelectCouponChk string = `SELECT COUNT(*) as cnt
									FROM b_coupon_his
									WHERE 
									COUPON_NO =  '#{couponNo}'
									AND USER_KEY= '#{userKey}'
									AND USER_TYPE= '#{userType}'
                             		`

var SelectBillingChk string = `SELECT END_DATE
							,B_ID
							FROM e_billing AS A
							WHERE 
							STORE_ID='#{storeId}'
							`

var InsertBillingCouponUse string = `INSERT INTO e_billing
		(
			USER_ID
			,STORE_ID
			, NEXT_PAY_DAY	
			, REG_DATE
			, ITEM_CODE
			, START_DATE
			, END_DATE
			, PAY_YN
		)
		VALUES
		(
			'#{userId}'
			,'#{storeId}' 
			, DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')	 
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, '#{itemCode}'
			, DATE_FORMAT( now(), '%Y-%m-%d')
			, DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
			, '#{payYn}'
		)
		`

var UpdateBillingCouponUse string = `UPDATE e_billing SET
									MOD_DATE =DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									,NEXT_PAY_DAY = CASE WHEN END_DATE >= SYSDATE() THEN  DATE_FORMAT( date_add(END_DATE, interval '#{useMonth}' month) , '%Y-%m-%d')
														ELSE  DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,ITEM_CODE = '#{itemCode}'
									,START_DATE = CASE WHEN END_DATE >= SYSDATE() THEN  START_DATE
														ELSE   DATE_FORMAT( date_add(START_DATE, interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,END_DATE = CASE WHEN END_DATE >= SYSDATE() THEN  DATE_FORMAT( date_add(END_DATE, interval '#{useMonth}' month) , '%Y-%m-%d')
														ELSE  DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,PAY_YN = '#{payYn}'
								WHERE
									STORE_ID ='#{storeId}'
									AND USER_ID = '#{userId}'
								`

var InserCouponHistory string = `INSERT INTO b_coupon_his
								(COUPON_NO
								, USER_KEY
								, USER_TYPE
								, REG_DATE
								)
								VALUES 
								(
									'#{couponNo}'
									, '#{userKey}'
									, '#{userType}'
									, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								)
							`

var InserBillingPayment string = `INSERT INTO e_billing_payment
										(
										 PCD_PAY_OID
										  ,STORE_ID
										 ,REG_DATE
                                         ,PAY_TYPE
										 ,ETC
										)
										VALUES
										(
											concat('#{couponNo}'
												,'_','#{storeId}')
											, '#{storeId}'
											, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
											, 'C'
											, '#{etc}'
										)
									`

var InsertConfirmData string = `INSERT INTO biz_confirm_info 
									(
										CONFIRM_DIV
										,TEL_NUM 
										,EMAIL 
										,SND_DTM 
										,CONFIRM_NUM
										,CONFIRM_YN
										,CONFIRM_DTM
										,MOD_DTM
									)
									SELECT 
										IFNULL('#{confirmDiv}','')               
										,IFNULL('#{telNum}','') 
										,''                  
										,DATE_FORMAT(SYSDATE(), '%Y%m%d%H%i%s') 
										,IFNULL('#{confirmNum}','')               
										,'N'                                    
										, ''                                    
										, DATE_FORMAT(SYSDATE(), '%Y%m%d%H%i%s') 
									FROM DUAL`

var SelectConfirmCheck string = `SELECT 
									g.CONFIRM_DIV
									,g.TEL_NUM 
									,g.EMAIL 
									,g.SND_DTM 
									,g.CONFIRM_NUM 
									,DATE_FORMAT(DATE_ADD(SYSDATE(),INTERVAL -3 MINUTE), '%Y%m%d%H%i%s') AS NOW_DATE
									FROM biz_confirm_info g
									WHERE 
									g.CONFIRM_DIV = '#{confirmDiv}'
									AND  G.TEL_NUM='#{telNum}'
									`

var UpdateConfirmData string = `UPDATE biz_confirm_info SET CONFIRM_YN='Y'
															,CONFIRM_DTM=DATE_FORMAT(SYSDATE(), '%Y%m%d%H%i%s') 
								WHERE 
								CONFIRM_DIV = '#{confirmDiv}'
								AND TEL_NUM='#{telNum}'
								 `

var UpdateConfirmDataReset string = `UPDATE biz_confirm_info SET CONFIRM_YN='N'
								,SND_DTM=DATE_FORMAT(SYSDATE(), '%Y%m%d%H%i%s') 
								,CONFIRM_NUM ='#{confirmNum}'
								,MOD_DTM=DATE_FORMAT(SYSDATE(), '%Y%m%d%H%i%s') 
								,CONFIRM_DTM=''
								WHERE 
								CONFIRM_DIV = '#{confirmDiv}' 
								AND TEL_NUM='#{telNum}'
								`

var SelectCompAuthInfo string = `SELECT LN_AUTH_FAIL
										,LN_FAIL_DT
										,HOMETAX_AUTH_FAIL
										,HOMTAXT_FAIL_DT
										,TIMESTAMPDIFF(MINUTE,SYSDATE(),DATE_ADD(SYSDATE(), INTERVAL 30 MINUTE)) AS LN_FAIL_DT
										,TIMESTAMPDIFF(MINUTE,SYSDATE(),DATE_ADD(SYSDATE(), INTERVAL 30 MINUTE)) AS HOMTAXT_FAIL_DT
								FROM cc_comp_inf
								WHERE
									BIZ_NUM='#{bizNum}'

								`

var UpdateLnAuthFail string = `UPDATE cc_comp_inf SET LN_AUTH_FAIL= #{lnAuthFailCnt}
										,LN_FAIL_DT = DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s')
										,LN_ID ='#{loginId}'
										,LN_PSW ='#{password}'
										,LN_JOIN_STS_CD =3
										WHERE
										BIZ_NUM='#{bizNum}'
								`

var UpdateLnAuthSuccess string = `UPDATE cc_comp_inf SET LN_AUTH_FAIL= 0
										,LN_FAIL_DT = ''
										,LN_ID ='#{loginId}'
										,LN_PSW ='#{password}'
										,LN_JOIN_STS_CD =1
										WHERE
										BIZ_NUM='#{bizNum}'
								`

var UpdateHomeTaxAuthFail string = `UPDATE cc_comp_inf SET HOMETAX_AUTH_FAIL = #{homeTaxAuthFailCnt}
										,LN_FAIL_DT = DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s')
										,HOMETAX_ID= '#{loginId}'
										,HOMETAX_PSW='#{password}'
										,HOMETAX_JOIN_STS_CD=3
										WHERE
										BIZ_NUM='#{bizNum}'
								`

var UpdateHomeTaxAuthSuccess string = `UPDATE cc_comp_inf SET HOMETAX_JOIN_STS_CD=1
										,HOMETAX_ID= '#{loginId}'
										,HOMETAX_PSW='#{password}'
										,HOMETAX_AUTH_FAIL=0
										,HOMTAXT_FAIL_DT=''
										WHERE
<<<<<<< .mine
										BIZ_NUM='#{bizNum}'
								`

var SelectPushUser string = `SELECT
							A.USER_ID
							, A.REG_ID
							, A.OS_TY
							, A.REG_DATE
							, A.LOGIN_YN
							, B.PUSH_YN
							FROM
							SYS_REG_INFO A,
							PRIV_USER_INFO B
							WHERE
							A.USER_ID = B.USER_ID
							AND B.USE_YN = 'Y'
							AND A.USER_ID = '#{userId}'
							`

var SelectPushBizNum string = `SELECT C.USER_ID
							, C.REG_ID
							, C.OS_TY
							, C.LOGIN_YN
							, B.PUSH_YN
							FROM PRIV_REST_INFO AS A
							INNER JOIN PRIV_REST_USER_INFO AS B ON A.REST_ID = B.REST_ID
							INNER JOIN SYS_REG_INFO AS C ON B.USER_ID = C.USER_ID
							WHERE
							BUSID='#{bizNum}'
							AND B.REST_AUTH=0
							`

var SelectPushRest string = `SELECT C.USER_ID
							, C.REG_ID
							, C.OS_TY
							, C.LOGIN_YN
							, B.PUSH_YN
							FROM PRIV_REST_INFO AS A
							INNER JOIN PRIV_REST_USER_INFO AS B ON A.REST_ID = B.REST_ID
							INNER JOIN SYS_REG_INFO AS C ON B.USER_ID = C.USER_ID
							WHERE
							A.REST_ID='#{restId}'
							AND B.REST_AUTH=0
							`

var SelectPushGrp string = `SELECT
								A.USER_ID
								, A.REG_ID
								, A.OS_TY
								, A.LOGIN_YN
								, C.PUSH_YN
							FROM 
								SYS_REG_INFO A,
								PRIV_USER_INFO C
							WHERE 
								EXISTS (SELECT USER_ID FROM PRIV_GRP_USER_INFO B
										WHERE 
										A.USER_ID = B.USER_ID 
										AND B.GRP_ID = '#{grpId}'
										AND B.GRP_AUTH = '0'
								)
							AND A.LOGIN_YN = 'Y'
							AND A.USER_ID = C.USER_ID
							`

var RegDateOrderBy = `ORDER BY A.REG_DATE DESC`

var JoinDateOrderBy = `ORDER BY A.JOIN_DATE DESC`
