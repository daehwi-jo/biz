package query

var SelectLoginIdDupCheck string = `SELECT count(*) as loginIdCnt
								FROM priv_user_info
								WHERE 
								LOGIN_ID ='#{bizNum}'`

var SelectBizNumDupCheck string = `SELECT count(*) as bizCnt
								FROM priv_rest_info
								WHERE 
								BUSID ='#{bizNum}'`

var SelectUserLoginCheck string = `SELECT USER_ID
											,USER_NM
											,AUTHOR_CD
											,CONN_ALLOW_YN
									FROM sys_user_info
									WHERE 
									USER_ID='#{loginId}'
									AND USER_PASS='#{password}'
									`

var SelectCompanyData string = `SELECT B.COMPANY_ID
									,B.COMPANY_NM
									,C.USER_ID
									FROM b_company_manager AS A
									INNER JOIN b_company AS B ON A.COMPANY_ID = B.COMPANY_ID
									INNER JOIN priv_user_info AS C ON A.USER_ID = C.USER_ID
									WHERE 
									C.LOGIN_ID='#{loginId}'
									`

var SelectJoinCheck string = `SELECT A.USER_ID
									,A.USER_NM
									,A.HP_NO
									,A.LOGIN_ID
									,IFNULL(B.REST_ID,'NONE') AS STORE_ID
									,IFNULL(C.REST_NM,'NONE') AS REST_NM
									,IFNULL(C.CEO_NM,'NONE') AS CEO_NM
								FROM priv_user_info AS A
								LEFT OUTER JOIN priv_rest_user_info AS b ON a.user_id = b.user_id
								LEFT OUTER JOIN priv_rest_info AS c ON b.rest_id = c.rest_id
								WHERE 
								USER_TY='1'
								AND A.USER_ID='#{userId}'
								`

var SelectUserStoreInfo string = `SELECT A.REST_ID AS STORE_ID
											,B.REST_NM AS STORE_NM
											,B.BUSID  AS BIZ_NUM
									FROM priv_rest_user_info AS A
									INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
									WHERE 
									A.USER_ID ='#{userId}'
									`

var InsertLoginAccess string = `	INSERT INTO SYS_LOG_ACCESS (
											USER_ID
											,ADDED_DATE
											,LOG_IN_OUT
											,IP
											,SUCC_YN 
											,SERVICE 
											,TYPE
								) VALUES (
									'#{loginId}'
									,SYSDATE()
									,'#{logInOut}'
									,'#{ip}'
									,'#{succYn}'
									,'#{osTy}'
									,'#{type}'
								)
								`

var SelectIntroMsg string = `SELECT CODE_NM as introMsg
								FROM b_code
								WHERE CATEGORY_ID='M001'
								mngOrder by RAND()
								LIMIT 1
								`

var SelectStoreService string = `SELECT IFNULL(B.B_ID,'N') AS billingYn
										,IFNULL(C.LN_ID,'N') AS cardSalesYn
										,A.ADDR
										,IFNULL(C.HOMETAX_ID,'N') AS homeTaxYn
										,A.REST_NM
										,AA.REST_AUTH
								FROM priv_rest_info AS A
								INNER JOIN priv_rest_user_info AS AA ON A.REST_ID = AA.REST_ID  
								LEFT OUTER JOIN e_billing AS B ON A.REST_ID = B.STORE_ID  AND B.END_DATE >= SYSDATE()
								LEFT OUTER JOIN cc_comp_inf AS C ON A.REST_ID = C.REST_ID
								WHERE 
								A.REST_ID= '#{storeId}'
								AND AA.USER_ID= '#{userId}'
								`

var SelectCreatUserSeq string = `SELECT CONCAT('U',IFNULL(LPAD(MAX(SUBSTRING(USER_ID, -10)) + 1, 10, 0), '0000000001')) as newUserId
								 FROM priv_user_info`

var InserCreateUser string = `INSERT INTO priv_user_info
									(
										USER_ID,
										USER_NM,
										LOGIN_ID,
										LOGIN_PW,
										USER_TY,
										HP_NO,
										ATLOGIN_YN,
										GEOLOC_YN,
										PUSH_YN,
										if #{kakaoPw} != '' then KAKAO_PW,
										if #{kakaoKey} != '' then KAKAO_KEY,
										if #{applePw} != '' then APPLE_PW,
										if #{appleKey} != '' then APPLE_KEY,
										if #{naverPw} != '' then NAVER_PW,
										if #{naverKey} != '' then NAVER_KEY,
										if #{recomCode} != '' then RECOM_CODE,
										if #{channelCode} != '' then CHANNEL_CODE,
										if #{userBirth} != '' then USER_BIRTH,
										USE_YN,
										JOIN_DATE
									)
									VALUES
									(
										'#{userId}'
										, '#{userNm}'
										, '#{loginId}'
										, '#{loginPw}'
										, '#{userTy}'
										, '#{userTel}'
										, '#{atLoginYn}'
										, 'Y'
										, '#{pushYn}'
										, '#{kakaoPw}'
										, '#{kakaoKey}'
										, '#{applePw}'
										, '#{appleKey}'
										, '#{naverPw}'
										, '#{naverKey}'
										, '#{recomCode}'
										, '#{channelCode}'
										, '#{userBirth}'
										, 'S'
										, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										)`

var InsertTermsUser string = `INSERT INTO b_user_terms
							(USER_ID, 
							TERMS_OF_SERVICE, 
							TERMS_OF_PERSONAL, 
							TERMS_OF_PAYMENT, 
							TERMS_OF_BENEFIT, 
							REG_DATE
							)
							VALUES (
							'#{userId}'
							,'#{termsOfService}'
							,'#{termsOfPersonal}'
							,'#{termsOfPayment}'
							,'#{termsOfBenefit}'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')						
							)
							`

var SelectStoreSeq string = `SELECT CONCAT('S',IFNULL(LPAD(MAX(SUBSTRING(REST_ID, -10)) + 1, 10, 0), '0000000001')) as storeSeq
							FROM priv_rest_info
							`

var SelectUserInfo string = `SELECT USER_NM AS userName
							, HP_NO as userTel
							, ifnull(USER_BIRTH,'') as birthday
							, EMAIL AS email
							, use_yn
							, RECOM_CODE as recomCode
							FROM priv_user_info
							WHERE 
							USER_ID ='#{userId}'  
							`

var UpdateUserInfo string = `UPDATE priv_user_info SET 
												 MOD_DATE= DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')	
                                                , USER_BIRTH = '#{birthday}'
												, USER_NM='#{userName}'
												, USER_TEL='#{userTel}'
												, USE_YN='#{useYn}'
												, LOGIN_PW='#{loginPw}'
							WHERE 
							USER_ID ='#{userId}'`

var InsertStore string = `INSERT INTO priv_rest_info
							(
							  REST_ID
							, REST_NM
							, BUSID
							, CATEGORY
							, BUETY
							, ADDR
							, ADDR2
							, LAT
							, LNG
							, AUTH_STAT
							, USE_YN
							, CEO_BIRTHDAY
							, CEO_TTI
							, CEO_NM
							, TEL
							, EMAIL
							, H_CODE
							, REG_DATE
							)
							VALUES (
							'#{storeId}'
							,'#{storeNm}'
							,'#{bizNum}'
							,'#{category}'
							,'#{kind}'
							,'#{addr}'
							,'#{addr2}'
							,'#{lat}'
							,'#{lng}'
							,'1'
							,'Y'
							,'#{ceoBirthday}'
							,'#{ceoTti}'
							,'#{ceoName}'
							,'#{storeTel}'
							,''
							,'#{hCode}'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')						
							)
							`

var InsertStoreUser string = `INSERT INTO priv_rest_user_info
								(REST_ID
								, USER_ID
								, REST_AUTH
								, PRINT_CON_YN
								, DAYSUM_YN
								, MONSUM_YN
								, PAYHIST_YN
								, GRPHIST_YN
								, PREPAID_YN
								, UNPAID_YN
								, ORDER_YN
								, MENY_YN
								, AGRM_YN
								, EVENT_YN
								, PUSH_YN
								, USE_YN
								, REG_DATE
								)
								VALUES (
										'#{storeId}'
										,'#{userId}'
										,'0'
										,'N'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								)
								`
var SelectUserIdSearch string = `SELECT LOGIN_ID
								,DATE_FORMAT(JOIN_DATE, '%Y.%m.%d') AS JOIN_DATE
								,CASE 	WHEN kakao_key IS NOT NULL THEN 'KAKAO'
										WHEN apple_key IS NOT NULL THEN 'APPLE'
										WHEN naver_key IS NOT NULL THEN 'NAVER'
									ELSE 'ID' END AS LOGIN_TYPE
								FROM priv_user_info
								WHERE 
								USER_NM='#{userNm}'
								AND HP_NO='#{userTel}'
								AND USER_TY='1'
								`

var SelectUserPwSearch string = `SELECT COUNT(*) AS CNT
								,USER_ID
								FROM priv_user_info
								WHERE 
								USER_NM='#{userNm}'
								AND LOGIN_ID = '#{loginId}'
								AND HP_NO='#{userTel}'
								AND USER_TY='1'
								`
