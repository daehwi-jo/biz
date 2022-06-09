package payment

var SelectPaidMngListCnt = `SELECT COUNT(*) as TOTAL_COUNT
						FROM dar_payment_report AS A
						INNER JOIN  dar_payment_hist AS B ON A.MOID = B.MOID
						INNER JOIN priv_rest_info AS C ON B.REST_ID = C.REST_ID
						WHERE B.REST_PAYMENT_ID ='X'
							AND B.PAYMENT_TY IN ('0', '3')
							AND B.PAY_CHANNEL  IN ('02', '03')    
							AND A.STATECD    = '0'  
							AND A.STATE      = '0000' 
							AND A.PAYMENT_DT < DATE_FORMAT(DATE_ADD(NOW(), INTERVAL -1 DAY), '%Y%m%d')
							AND C.REST_NM LIKE '%#{restNm}%'
							AND C.REST_TYPE='N'
                     `

var SelectPaidMngList = `SELECT B.REST_ID
				,C.REST_NM
				,DATE_FORMAT(A.PAYMENT_DT,'%Y-%m-%d') as PAYMENT_DT
				,B.CREDIT_AMT
				,replace(B.REST_PAY_AMT,'.0','') as REST_PAY_AMT
				,replace(B.TOT_FEE,'.0','') as TOT_FEE
				,PAY_PROXY_CD as payProxyCd
		FROM dar_payment_report AS A
		INNER JOIN  dar_payment_hist AS B ON A.MOID = B.MOID
		INNER JOIN priv_rest_info AS C ON B.REST_ID = C.REST_ID
		WHERE 
			B.REST_PAYMENT_ID = 'X'
			AND B.PAYMENT_TY IN ('0', '3')
			AND B.PAY_CHANNEL  IN ('02', '03')    
			AND A.STATECD    = '0'  
			AND A.STATE      = '0000'
			AND C.REST_NM LIKE '%#{restNm}%'
			AND C.REST_TYPE='N'
		ORDER BY A.PAYMENT_DT DESC`

var SelectCombineStoreListCnt = `SELECT COUNT(*) as TOTAL_COUNT
								FROM b_rest_combine AS A
								INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
                     `

var SelectCombineStoreList = `SELECT A.REST_ID
										,B.REST_NM
										,A.CHARGE_AMT
										,0 AS SETTLEMENT_ING_AMT
										,PAY_PROXY_CD as payProxyCd
								FROM b_rest_combine AS A
								INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
                     `

var SelectCombineSubStoreList = `SELECT B.REST_ID
									 ,B.REST_NM
								FROM b_rest_combine_sub AS A
								INNER JOIN priv_rest_info AS B ON A.SUB_REST_ID = B.REST_ID
								WHERE 
									A.REST_ID='#{restId}'
									
                     `

var SelectPaidListCnt = `SELECT COUNT(*) AS TOTAL_COUNT
				,ifnull(sum(A.PAYMENT_CNT),0) as ALL_PAYMENT_CNT
				,ifnull(sum(A.PAYMENT_AMT),0) as ALL_PAYMENT_AMT
				,ifnull(sum(A.REST_PAYMENT_AMT),0) as ALL_REST_PAYMENT_AMT
				,ifnull(sum(A.TOT_FEE),0) as ALL_TOT_FEE
				,ifnull(sum(A.CANCEL_CNT),0) as ALL_CANCEL_CNT
				,ifnull(sum(A.CANCEL_AMT),0) as ALL_CANCEL_AMT
       FROM DAR_REST_PAYMENT AS A
		INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
		WHERE
		(A.RESULT_CD NOT IN ('0000','9999','0001') OR A.RESULT_CD IS NULL)
		AND A.PAYMENT_DT >= '#{searchStartDt}' 
		AND A.PAYMENT_DT <= '#{searchEndDt}'
		AND B.REST_NM LIKE '%#{searchText}%'
		ORDER BY A.SEND_DATE DESC,A.PAYMENT_DT ASC
`

var SelectPaidList = `SELECT	
				 REST_PAYMENT_ID  						AS restPaymentId  		-- 가맹점 정산 아이디 
				,B.REST_NM 				 				AS restNm 				-- 가맹점명
				,B.BUSID				 				AS busId
				,DATE_FORMAT(A.PAYMENT_DT,'%Y-%m-%d')   AS paymentDt  			-- 결제일자 
				,DATE_FORMAT(A.SETTLMNT_DT,'%Y-%m-%d')	 	 				AS settlmntDt 	-- 정산 지급 요청 일자 
				,A.PAYMENT_CNT           				AS paymentCnt   		-- 결제 건수
				,A.PAYMENT_AMT           				AS paymentAmt			-- 결제 금액
				,A.REST_PAYMENT_AMT      				AS restPaymentAmt  	-- 지급 요청 금액
				,A.FIT_FEE               				AS fitFee				-- 총수수료
				,A.PG_FEE                				AS pgFee					-- FIT 수수료
				,A.TOT_FEE               				AS totFee				-- pg 수수료 
				,A.CANCEL_AMT            				AS cancelAmt			-- 취소 금액
				,A.CANCEL_CNT            				AS cancelCnt			-- 취소 건수 
				,IFNULL(A.RESULT_MSG,'')   				AS resultMsg      -- 결과 메세지
				,DATE_FORMAT(CONCAT(A.SEND_DATE,A.SEND_TIME),'%Y-%m-%d %h:%i:%s') as sendDate
		FROM DAR_REST_PAYMENT AS A
		INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
		WHERE
		(A.RESULT_CD NOT IN ('0000','9999','0001') OR A.RESULT_CD IS NULL)
		AND A.PAYMENT_DT >= '#{searchStartDt}' 
		AND A.PAYMENT_DT <= '#{searchEndDt}'
		AND B.REST_NM LIKE '%#{searchText}%'


		ORDER BY A.SEND_DATE DESC,A.PAYMENT_DT ASC
`

var SelectPaidExcelList = `SELECT D.REST_ID as EX01
								,D.BUSID as EX02
								,DATE_FORMAT(CONCAT(C.PAYMENT_DT,C.PAYMENT_TM),'%Y-%m-%d %h:%i') as EX03
								,DATE_FORMAT(A.SETTLMNT_DT,'%Y-%m-%d')	 	 				as EX04
								,PAYMETHOD as EX05
								,FNNAME as EX06
								,GOODSNAME as EX07
								,IF(A.PAYMENT_CNT = 1, A.PAYMENT_AMT, C.AMT)               as EX08
								,IF(A.PAYMENT_CNT = 1, A.REST_PAYMENT_AMT, C.REST_PAY_AMT) as EX09
								,IF(A.PAYMENT_CNT = 1, A.TOT_FEE, C.TOT_FEE)               as EX10
								,E.HP_NO as EX11
								,E.USER_NM as EX12
								,D.REST_NM AS EX13
						FROM dar_rest_payment AS A
						INNER JOIN dar_payment_hist AS B ON A.REST_PAYMENT_ID = B.REST_PAYMENT_ID
						INNER JOIN dar_payment_report AS C ON B.MOID = C.MOID
						INNER JOIN PRIV_REST_INFO AS D ON B.REST_ID = D.REST_ID
						INNER JOIN priv_user_info AS E ON B.USER_ID = E.USER_ID
		WHERE
		(A.RESULT_CD NOT IN ('0000','9999','0001') OR A.RESULT_CD IS NULL)
		AND C.PAYMENT_DT >= '#{searchStartDt}' 
		AND C.PAYMENT_DT <= '#{searchEndDt}'
		AND D.REST_NM LIKE '%#{searchText}%'
		ORDER BY A.SEND_DATE DESC,C.PAYMENT_DT ASC
`

var SelectPaidIngListCnt = `SELECT COUNT(*) AS TOTAL_COUNT
       	FROM DAR_REST_PAYMENT AS A
		INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
		WHERE  A.RESULT_CD  IN('0000','0001')

`

var SelectPaidIngList = `SELECT	
				 REST_PAYMENT_ID  						AS restPaymentId  		-- 가맹점 정산 아이디 
				,B.REST_NM 				 				AS restNm 				-- 가맹점명
				,B.BUSID				 				AS busId
				,DATE_FORMAT(A.PAYMENT_DT,'%Y-%m-%d')   AS paymentDt  			-- 결제일자 
				,DATE_FORMAT(A.SETTLMNT_DT,'%Y-%m-%d')	AS settlmntDt 	-- 정산 지급 요청 일자 
				,A.PAYMENT_CNT           				AS paymentCnt   		-- 결제 건수
				,A.PAYMENT_AMT           				AS paymentAmt			-- 결제 금액
				,A.REST_PAYMENT_AMT      				AS restPaymentAmt  	-- 지급 요청 금액
				,A.FIT_FEE               				AS fitFee				-- 총수수료
				,A.PG_FEE                				AS pgFee					-- FIT 수수료
				,A.TOT_FEE               				AS totFee				-- pg 수수료 
				,A.CANCEL_AMT            				AS cancelAmt			-- 취소 금액
				,A.CANCEL_CNT            				AS cancelCnt			-- 취소 건수 
				,A.RESULT_MSG            				AS resultMsg      -- 결과 메세지
				,DATE_FORMAT(CONCAT(A.SEND_DATE,A.SEND_TIME),'%Y-%m-%d %h:%i:%s') as sendDate
				,A.RESULT_CD
				,A.RESULT_SEQ
				,A.SETTLMNT_DT
		FROM DAR_REST_PAYMENT AS A
		INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
		WHERE  A.RESULT_CD  IN('0000','0001')
`

var SelectPaidOkListCnt = `SELECT COUNT(*) AS TOTAL_COUNT
       	FROM DAR_REST_PAYMENT AS A
		INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
		WHERE  A.RESULT_CD in('0000','0001')
		AND A.PAYMENT_DT >= '#{startDate}' 
		AND A.PAYMENT_DT <= '#{endDate}'
		AND A.RESULT_PAY_YN='Y'
		AND B.REST_NM LIKE CONCAT('%','#{searchRestNm}', '%')
		AND B.BUSID LIKE CONCAT('%','#{searchBusid}', '%')
		ORDER BY A.SEND_DATE DESC,A.PAYMENT_DT ASC
`

var SelectPaidOkList = `SELECT	
				 REST_PAYMENT_ID  						AS restPaymentId  		-- 가맹점 정산 아이디 
				,B.REST_NM 				 				AS restNm 				-- 가맹점명
				,B.BUSID				 				AS busId
				,DATE_FORMAT(A.PAYMENT_DT,'%Y-%m-%d')   AS paymentDt  			-- 결제일자 
				,DATE_FORMAT(A.SETTLMNT_DT,'%Y-%m-%d')	 	 				AS settlmntDt 	-- 정산 지급 요청 일자 
				,A.PAYMENT_CNT           				AS paymentCnt   		-- 결제 건수
				,A.PAYMENT_AMT           				AS paymentAmt			-- 결제 금액
				,A.REST_PAYMENT_AMT      				AS restPaymentAmt  	-- 지급 요청 금액
				,A.FIT_FEE               				AS fitFee				-- 총수수료
				,A.PG_FEE                				AS pgFee					-- FIT 수수료
				,A.TOT_FEE               				AS totFee				-- pg 수수료 
				,A.CANCEL_AMT            				AS cancelAmt			-- 취소 금액
				,A.CANCEL_CNT            				AS cancelCnt			-- 취소 건수 
				,A.RESULT_MSG            				AS resultMsg      -- 결과 메세지
				,DATE_FORMAT(CONCAT(A.SEND_DATE,A.SEND_TIME),'%Y-%m-%d %h:%i:%s') as sendDate
				,A.RESULT_CD
				,A.RESULT_SEQ
		FROM DAR_REST_PAYMENT AS A
		INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
		WHERE A.RESULT_CD in('0000','0001')
		AND A.PAYMENT_DT >= '#{startDate}' 
		AND A.PAYMENT_DT <= '#{endDate}'
		AND A.RESULT_PAY_YN='Y'
		AND B.REST_NM LIKE CONCAT('%','#{searchRestNm}', '%')
		AND B.BUSID LIKE CONCAT('%','#{searchBusid}', '%')
		ORDER BY A.SEND_DATE DESC,A.PAYMENT_DT ASC
`

var SelectPaidOkExcelList = `SELECT C.REST_NM as A01
							,C.BUSID  as A02
							,DATE_FORMAT(A.PAYMENT_DT,'%Y-%m-%d') as A03
							,B.CREDIT_AMT as A04
							,replace(B.REST_PAY_AMT,'.0','') as A05
							,replace(A.FIT_FEE,'.0','') as A06
                            ,replace(A.PG_FEE,'.0','') as A07
							,replace(A.TOT_FEE,'.0','') as A08
							,CASE WHEN PAY_INFO='1' THEN '카드'
										WHEN PAY_INFO='0' THEN 'BANK'
										WHEN PAY_INFO='2' THEN 'BANK'
										WHEN PAY_INFO='3' THEN '현금'
										END  as A09
							FROM dar_rest_payment AS AA
							INNER JOIN dar_payment_hist  AS B ON AA.REST_PAYMENT_ID = B.REST_PAYMENT_ID
							INNER JOIN dar_payment_report AS A  ON A.MOID = B.MOID
							INNER JOIN priv_rest_info AS C ON B.REST_ID = C.REST_ID
							WHERE
							AA.RESULT_CD in('0000','0001')
							AND AA.PAYMENT_DT >= '#{startDate}'
							AND AA.PAYMENT_DT <= '#{endDate}'
							AND C.REST_NM LIKE CONCAT('%', '#{searchRestNm}', '%')
							AND C.BUSID LIKE CONCAT('%','#{searchBusid}', '%')
							AND AA.RESULT_PAY_YN='Y'
							ORDER BY AA.SEND_DATE DESC,AA.PAYMENT_DT ASC
					`

var SelectPgInfo = `SELECT PG_MID
							,PG_UID
							,PG_PSWD
							,PG_MERCHANT_KEY
							,PAY_ACTION_URL
							,PAY_LOCAL_URL
							,RETURN_URL
					FROM sys_pg_info
					WHERE PG_CD='01' AND SERVER_CL ='03'
`

var SelectRestPaymentSendList = `SELECT REST_ID
									,SETTLMNT_DT
									,REST_PAYMENT_AMT
									,REST_PAYMENT_ID
							FROM DAR_REST_PAYMENT
							WHERE 
							REST_PAYMENT_ID IN(#{restPaymentIdArray}) 
`

var UpdateDarRestPaymentTpayRegResult = `UPDATE DAR_REST_PAYMENT SET
										          RESULT_SUB_ID 		= '#{sub_id}'
												 ,RESULT_SETTLMNT_DT 	= '#{settlmnt_dt}'
												 ,RESULT_SEQ 			= '#{seq}'
												 ,RESULT_MSG 	        = '#{result_msg}'
												 ,RESULT_CD 	        = '#{result_cd}'
												 ,SEND_DATE 	        = DATE_FORMAT(NOW(), '%Y%m%d')
												 ,SEND_TIME 	        = DATE_FORMAT(NOW(), '%H%i%s')
										WHERE   
												REST_PAYMENT_ID	 = '#{restPaymentId}'
											`

var UpdateDarRestPaymentTpayRegResultFail = `UPDATE DAR_REST_PAYMENT SET 
												      RESULT_MSG = '#{result_msg}'
													, RESULT_CD = '#{result_cd}'
												WHERE     
												 REST_PAYMENT_ID='#{restPaymentId}'
											`

var UpdateDarRestPaymentTpayResultPay = `UPDATE DAR_REST_PAYMENT SET
										          RESULT_PAY_YN 		= 'Y'
										WHERE   
												REST_PAYMENT_ID	 = '#{restPaymentId}'
											`

var UpdateDarRestPaymentSettlmentDt = `UPDATE DAR_REST_PAYMENT SET
										          SETTLMNT_DT = '#{settlmntDt}'
										WHERE   
												REST_PAYMENT_ID IN (#{restPaymentIdArray})
										`

var UpdateDarRestPayment = `UPDATE DAR_REST_PAYMENT SET
										          SETTLMNT_DT = '#{settlmntDt}'
												 ,REST_PAYMENT_AMT = '#{restPaymentAmt}'
										WHERE   
												REST_PAYMENT_ID = '#{restPaymentId}'
										`

var SelectPaidInfo = `		SELECT   
											REST_PAYMENT_ID		 					AS restPaymentId  		-- 가맹점 정산 아이디
											,B.REST_NM 				 					AS restNm 				-- 가맹점명
											,B.BUSID				 					AS busId
											,DATE_FORMAT(A.SETTLMNT_DT,'%Y-%m-%d')		AS settlmntDt
											,A.PAYMENT_AMT           					AS paymentAmt
											,A.REST_PAYMENT_AMT      					AS restPaymentAmt
											,DATE_FORMAT(A.PAYMENT_DT,'%Y-%m-%d')       AS paymentDt  			-- 결제일자
											FROM DAR_REST_PAYMENT AS A
											INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
											WHERE 
											REST_PAYMENT_ID= '#{restPaymentId}'
											`

var SelectRestFeesInfo = `SELECT
							REST_ID
							,REST_NM
							,PAYMETHOD
							,REST_FEES / 100 AS REST_FEES
							,USE_FEES_YN
						FROM DAR_REST_FEES
						WHERE
							REST_ID = '#{restId}'
						AND PAYMETHOD = '#{payMethod}'
						AND USE_FEES_YN = 'Y'
					`

var SelectPgFeesInfo = `SELECT
						ORG_CD
						,PAYMETHOD
						,FEE_CL
						,VAT_YN
						,FEE / 100 AS FEE
					FROM DAR_ORG_FEE
					WHERE
					 PAYMETHOD = '#{payMethod}'
					`

var InsertTpayPayment = `INSERT INTO DAR_PAYMENT_REPORT
							(
								  MOID
								,PAYMETHOD
								,TRANSTYPE
								,GOODSNAME
								,AMT
								,USERIP
								,TID
								,STATE
								,ADDAMT
								,STATECD
								,CARDNO
								,AUTHCODE
								,AUTHDATE
								,CARDQUOTA
								,FNCD
								,FNNAME
								,RESULTCD
								,RESULTMSG
								,PG_CD
								,PAYMENT_DT
								,PAYMENT_TM
							)
							VALUES
							(
								  '#{moid}'
								,'#{paymethod}'
								,'#{transtype}'
								,'#{goodsname}'
								,'#{amt}'
								,'#{userip}'
								,'#{tid}'
								,'#{state}'
								,'#{addAmt}'
								,'#{statecd}'
								,'#{cardno}'
								,'#{authcode}'
								,'#{authdate}'
								,'#{cardquota}'
								,'#{fncd}'
								,'#{fnname}'
								,'#{resultcd}'
								,'#{resultmsg}'
								,'#{pgCd}'
								,DATE_FORMAT(NOW(), '%Y%m%d')
								,DATE_FORMAT(NOW(), '%H%i%s')
							)
					`

var InsertTpayPaymentHistory = `INSERT INTO DAR_PAYMENT_HIST
								(
										  HIST_ID
										,REST_ID
										,GRP_ID
										,USER_ID
										,MOID
										,CREDIT_AMT
										,ADD_AMT
										,USER_TY
										,SEARCH_TY
										,PAYMENT_TY
										,PAY_INFO
										,REG_DATE
										,PG_CD
										,PAY_CHANNEL
								)
								VALUES
								(
										  '#{histId}'
										, '#{restId}'
										, '#{grpId}'
										, '#{userId}'
										, '#{moid}'
										, '#{creditAmt}'
										, '#{addAmt}'
										, '#{userTy}'
										, '#{searchTy}'
										, '#{paymentTy}'
										, '#{payInfo}'
										, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										, '#{pgCd}'
										, '#{payChannel}'          
								)
					`

var SelectRestUserInfo = `SELECT USER_NM
								,IFNULL(EMAIL,'') AS EMAIL
								,IFNULL(HP_NO,'') AS HP_NO
								,FN_GET_RESTNAME('#{restId}') AS REST_NM
						FROM priv_user_info
						WHERE 
						USER_ID= '#{mallUserId}'
					`

var SelectUnpaidPaymentInfo = `SELECT
								B.REST_NM
							  , (SELECT GRP_NM FROM PRIV_GRP_INFO WHERE GRP_ID = A.GRP_ID) AS GRP_NM
							  , IFNULL(SUM(CREDIT_AMT), 0) AS CREDIT_AMT
							  , DATE_FORMAT('#{selectedDate}', '%Y-%m-%d') AS AC_DATE
                              , B.REST_TYPE
                              , B.CEO_NM
							  FROM DAR_ORDER_INFO AS A
							  INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID	
							 WHERE 
							   A.REST_ID = '#{restId}'
                               AND A.order_ty IN ('1','2','3','5')
							   AND GRP_ID = '#{grpId}'
							   AND ORDER_STAT = '20'
							   AND PAID_YN = 'N'
							   AND PAY_TY = '1'
							   AND DATE_FORMAT(ORDER_DATE, '%Y%m%d') <= '#{selectedDate}'
							`

var UnpaidPaymentCreditNow = `UPDATE DAR_ORDER_INFO
										   SET PAID_YN = 'Y'
												 ,PAY_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
												 ,MOID = '#{moid}'
										 WHERE 
										   DATE_FORMAT(ORDER_DATE, '%Y%m%d') <= '#{selectedDate}'
										   AND REST_ID = '#{restId}'
										   AND order_ty IN ('1','2','3','5')
										   AND GRP_ID = '#{grpId}'
										   AND ORDER_STAT = '20'
										   AND PAID_YN = 'N'
										   AND PAY_TY = '1'
										`


var UnpaidPaymentOk = `UPDATE DAR_ORDER_INFO
										   SET PAID_YN = 'Y'
												 ,PAY_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
												 ,MOID = '#{moid}'
										 WHERE 
										   PAID_KEY  = '#{moid}'
										`

var SelectPayInfo = `SELECT GOODSNAME
							,AMT
							,FN_GET_USERNAME('#{mallUserId}') AS USER_NM
							,FN_GET_GRPNAME('#{sndGrpId}') AS GRP_NM
						FROM DAR_PAYMENT_REPORT
						WHERE
							MOID = '#{decMoid}'
						ORDER BY PAYMENT_DT DESC
						LIMIT 1
							`

var SelectTPayDupCheck = `SELECT COUNT(*) AS DUP_CNT
								FROM dar_payment_report
								WHERE 
									MOID='#{decMoid}'
									AND TID='#{tid}'
										`

var SelectPaidOrderAmtCheck = `SELECT IFNULL(SUM(TOTAL_AMT),0) AS TOTAL_AMT
								FROM DAR_ORDER_INFO
								WHERE 
								PAID_KEY='#{decMoid}'
										`


var SelectRestType = `SELECT REST_TYPE
							,IFNULL(CHARGE_AMT,0) AS CHARGE_AMT
								FROM priv_rest_info AS A
								LEFT OUTER JOIN b_rest_combine_sub AS B ON A.REST_ID = B.SUB_REST_ID
								LEFT OUTER JOIN b_rest_combine AS C ON B.REST_ID = C.REST_ID
								WHERE 
									A.rest_id='#{restId}'
										`


var SelectRestCombine = `SELECT CHARGE_AMT
						FROM b_rest_combine
						WHERE 
							rest_id='#{restId}'
										`

var UpdateAgrmPrepaidAmt = `UPDATE ORG_AGRM_INFO SET
							PREPAID_AMT = #{prepaidAmt},
							PREPAID_POINT = #{prepaidPoint},
							MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							WHERE 
 									AGRM_ID = '#{agrmId}'
										`

var SelectAgrmInfo = `SELECT 
						A.AGRM_ID
						, A.GRP_ID
						, A.REST_ID
						, A.REQ_STAT	
						, A.PAY_TY
						, IFNULL(A.PREPAID_AMT, 0) AS PREPAID_AMT
						, B.PAYMENT_USE_YN
						, IFNULL(A.PREPAID_POINT, 0) AS PREPAID_POINT
					FROM ORG_AGRM_INFO AS A
					left join priv_rest_info AS B on A.REST_ID = B.rest_id or (B.FRAN_YN = 'Y' and B.FRAN_ID = A.REST_ID)
					WHERE 
						A.GRP_ID = '#{grpId}'
						AND B.REST_ID = '#{restId}'
										`

var InsertPrepaid = `INSERT INTO DAR_PREPAID_INFO
									(
										  PREPAID_NO
										,GRP_ID
										,REST_ID
										,JOB_TY
										,PREPAID_AMT
										,REG_DATE
									)
									VALUES
									(
										  CONCAT(DATE_FORMAT(NOW(),'%Y%m%d%H%i%s'), '#{grpId}') 
										,'#{grpId}'
										,'#{restId}'
										,'#{jobTy}'
										,'#{prepaidAmt}'
										,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									)	
										`

var SelectPaidIngExcelList = `SELECT '01' AS A01
		,MAX(A.PAYMENT_DT) AS A02
		,'1198619035' AS A03
		,'' AS A04
		,'㈜에프아이티' AS A05
		,'조용준' AS A06
		,'서울시 구로구 디지털로31길 20, 409호(에이스테크노타워 5차)' AS A07
		,'서비스' AS A08
		,'소프트웨어 개발판매업' AS A09
		,'help@darayo.com' AS A10
		,C.BUSID AS A11
		,'' AS A12
		,C.REST_NM as A13
		,C.CEO_NM AS A14
		,CONCAT(ADDR,' ',ADDR2) AS A15
		,' ' AS A16
		,' ' AS A17
		,C.EMAIL AS A18
		,' ' AS A19
		,SUM(B.TOT_SUPLY_AMT) AS A20
		,SUM(B.TOT_VAT) AS A21
		,' ' AS A22
		,DATE_FORMAT(MAX(A.PAYMENT_DT),'%d') AS A23
		,'달아요 결제수수료' AS A24
		,' ' AS A25
		,' ' AS A26
		,' ' AS A27
		,SUM(B.TOT_SUPLY_AMT) AS A28
		,SUM(B.TOT_VAT) AS A29
		,' ' AS A30
		,' ' AS A31
		,' ' AS A32
		,' ' AS A33
		,' ' AS A34
		,' ' AS A35
		,' ' AS A36
		,' ' AS A37
		,' ' AS A38
		,' ' AS A39
		,' ' AS A40
		,' ' AS A41
		,' ' AS A42
		,' ' AS A43
		,' ' AS A44
		,' ' AS A45
		,' ' AS A46
		,' ' AS A47
		,' ' AS A48
		,' ' AS A49
		,' ' AS A50
		,' ' AS A51
		,' ' AS A52
		,' ' AS A53
		,' ' AS A54
		,' ' AS A55
		,' ' AS A56
		,' ' AS A57
		,' ' AS A58
		,'01' AS A59
 		FROM dar_rest_payment AS A
		INNER JOIN ( SELECT AA.REST_PAYMENT_ID
							,AA.USER_ID
							,AA.REST_ID
							,SUM(AA.TOT_SUPLY_AMT) AS TOT_SUPLY_AMT
							,SUM(AA.TOT_VAT) AS TOT_VAT
					FROM dar_payment_hist AS AA
					INNER JOIN dar_payment_report AS BB ON AA.MOID = BB.MOID
					WHERE 
						 BB.PAYMENT_DT >= '#{startDate}' 
					AND BB.PAYMENT_DT <= '#{endDate}'
					AND REST_PAYMENT_ID !='X'
					GROUP BY REST_PAYMENT_ID
							) AS B ON A.REST_PAYMENT_ID = B.REST_PAYMENT_ID		
		INNER JOIN priv_rest_info AS C ON B.REST_ID = C.REST_ID
		INNER JOIN priv_user_info AS D ON B.USER_ID = D.USER_ID
		WHERE 
        A.RESULT_CD in('0000','0001')
		AND A.PAYMENT_DT >= '#{startDate}' 
		AND A.PAYMENT_DT <= '#{endDate}'
		AND C.REST_NM LIKE CONCAT('%', #{searchKeyword}, '%') 
 		AND C.BUSID LIKE CONCAT('%', #{searchKeyword}, '%')
		AND A.RESULT_PAY_YN='N'
		GROUP BY A.REST_PAYMENT_ID
`

var UpdateTpayStoreRegResult = `UPDATE  PRIV_REST_INFO SET  
									 PAY_PROXY_CD = '#{resultCd}'
									,PAY_PROXY_REG_DT = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE 
									REST_ID='#{restId}'
											`

var SelectStorePayInfo = `SELECT REST_NM
									,BUSID
									,BANK_CD
									,ACCOUNT_NO
									,ACCOUNT_NM
									,IFNULL(ACCOUNT_CERT_YN,'N') AS ACCOUNT_CERT_YN
							FROM priv_rest_info
							WHERE 
								rest_id='#{restId}'
											`

var UpdateCombineChargeAmt = `
                            UPDATE b_rest_combine AS A
							INNER JOIN b_rest_combine_SUB AS B ON A.REST_ID = B.REST_ID
							SET A.CHARGE_AMT= #{chargeAmt}
								 ,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							WHERE 
								B.SUB_REST_ID='#{restId}'
										`

var UpdateCombinePaidChargeAmt = `
                            UPDATE b_rest_combine 
							SET CHARGE_AMT= #{chargeAmt}
								 ,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							WHERE 
								REST_ID='#{restId}'
										`

var SelectCombineOrderListCnt = `SELECT COUNT(*) as TOTAL_COUNT
								,SUM(CASE WHEN B.ORDER_STAT='20'  THEN B.CREDIT_AMT  ELSE 0 END) AS ORDER_AMT
								,SUM(CASE WHEN B.ORDER_STAT='21'  THEN B.CREDIT_AMT  ELSE 0 END) AS CANCEL_AMT
								,SUM(CASE WHEN CC.CP_STATUS='0' THEN B.CREDIT_AMT  ELSE 0 END) AS COUPON_0
								,SUM(CASE WHEN CC.CP_STATUS='1' THEN B.CREDIT_AMT  ELSE 0 END) AS COUPON_1
								,SUM(CASE WHEN CC.CP_STATUS='2' THEN B.CREDIT_AMT  ELSE 0 END) AS COUPON_2
								,SUM(CASE WHEN CC.CP_STATUS='0' THEN 1  ELSE 0 END) AS COUPON_0_CNT
								,SUM(CASE WHEN CC.CP_STATUS='1' THEN 1  ELSE 0 END) AS COUPON_1_CNT
								,SUM(CASE WHEN CC.CP_STATUS='2' THEN 1  ELSE 0 END) AS COUPON_2_CNT
								FROM b_rest_combine_sub AS a
								INNER JOIN dar_order_info AS b ON a.sub_rest_id = b.rest_id
								INNER JOIN dar_order_detail AS BB ON B.ORDER_NO = BB.ORDER_NO  
								INNER JOIN dar_order_coupon AS cc ON BB.ORDER_NO = cc.ORDER_NO  
								INNER JOIN priv_rest_info AS c ON b.REST_ID = c.rest_id
								INNER JOIN dar_sale_item_info AS D ON C.REST_ID = D.REST_ID AND BB.ITEM_NO = D.ITEM_NO
								WHERE 
									A.rest_id='#{restId}'
								AND (DATE_FORMAT(cc.cpno_exch_dt, '%Y-%m-%d') >=  '#{startDate}'
								AND DATE_FORMAT(cc.cpno_exch_dt, '%Y-%m-%d') <=  '#{endDate}'
								OR cc.cpno_exch_dt='' )
								AND B.REST_ID ='#{subRestId}' 
								ORDER BY cc.cpno_exch_dt DESC
											`

var SelectCombineOrderList = `SELECT DATE_FORMAT(cc.cpno_exch_dt, '%Y-%m-%d %H:%i:%s') AS ORDER_DATE	/* 주문일자 */
														, B.ORDER_NO			
														, B.REST_ID
														, C.REST_NM
														, B.CREDIT_AMT		
														, B.ORDER_STAT 
														, PAY_TY
												   		, D.ITEM_NM 
                                                    	, NOMAL_SALE_PRICE
														, NOMAL_SALE_VAT
														, SALE_PRICE
														, SALE_VAT
														, TOTAL_PRICE
								FROM b_rest_combine_sub AS a
								INNER JOIN dar_order_info AS b ON a.sub_rest_id = b.rest_id
								INNER JOIN dar_order_detail AS BB ON B.ORDER_NO = BB.ORDER_NO  
								INNER JOIN dar_order_coupon AS cc ON BB.ORDER_NO = cc.ORDER_NO
								INNER JOIN priv_rest_info AS c ON b.REST_ID = c.rest_id
								INNER JOIN dar_sale_item_info AS D ON C.REST_ID = D.REST_ID AND BB.ITEM_NO = D.ITEM_NO
								WHERE 
									A.rest_id='#{restId}'
								AND (DATE_FORMAT(cc.cpno_exch_dt, '%Y-%m-%d') >=  '#{startDate}'
								AND DATE_FORMAT(cc.cpno_exch_dt, '%Y-%m-%d') <=  '#{endDate}'
								OR cc.cpno_exch_dt='' )
								AND B.REST_ID ='#{subRestId}' 
								ORDER BY cc.cpno_exch_dt DESC
											`




var SelectCombineOrderWincubeListCnt = `SELECT COUNT(*) as TOTAL_COUNT
										,SUM(CASE WHEN B.ORDER_STAT='20'  THEN B.CREDIT_AMT  ELSE 0 END) AS ORDER_AMT
										,SUM(CASE WHEN B.ORDER_STAT='20'  THEN D.SALE_PRICE  ELSE 0 END) AS SALE_PRICE_AMT
										,SUM(CASE WHEN B.ORDER_STAT='21'  THEN B.CREDIT_AMT  ELSE 0 END) AS CANCEL_AMT
										,SUM(CASE WHEN B.ORDER_STAT='21'  THEN D.SALE_PRICE  ELSE 0 END) AS CANCEL_SALE_PRICE_AMT
										,SUM(CASE WHEN CC.CP_STATUS='0' THEN B.CREDIT_AMT  ELSE 0 END) AS COUPON_0
										,SUM(CASE WHEN CC.CP_STATUS='1' THEN B.CREDIT_AMT  ELSE 0 END) AS COUPON_1
										,SUM(CASE WHEN CC.CP_STATUS='2' THEN B.CREDIT_AMT  ELSE 0 END) AS COUPON_2
										,SUM(CASE WHEN CC.CP_STATUS='0' THEN 1  ELSE 0 END) AS COUPON_0_CNT
										,SUM(CASE WHEN CC.CP_STATUS='1' THEN 1  ELSE 0 END) AS COUPON_1_CNT
										,SUM(CASE WHEN CC.CP_STATUS='2' THEN 1  ELSE 0 END) AS COUPON_2_CNT
								FROM b_rest_combine_sub AS a
								INNER JOIN dar_order_info AS b ON a.sub_rest_id = b.rest_id
								INNER JOIN dar_order_detail AS BB ON B.ORDER_NO = BB.ORDER_NO  
								INNER JOIN dar_order_coupon AS cc ON BB.ORDER_NO = cc.ORDER_NO
								INNER JOIN priv_rest_info AS c ON b.REST_ID = c.rest_id
								INNER JOIN dar_sale_item_info AS D ON C.REST_ID = D.REST_ID AND BB.ITEM_NO = D.ITEM_NO
								WHERE 
									A.rest_id='#{restId}'
								AND B.ORDER_STAT='20'
								AND DATE_FORMAT(B.ORDER_DATE, '%Y-%m-%d') >=  '#{startDate}'
								AND DATE_FORMAT(B.ORDER_DATE, '%Y-%m-%d') <=  '#{endDate}'
								AND B.REST_ID ='#{subRestId}' 

								ORDER BY B.ORDER_DATE DESC
											`


var SelectCombineOrderWincubeList = `SELECT DATE_FORMAT(B.ORDER_DATE, '%Y-%m-%d %H:%i:%s') AS ORDER_DATE	/* 주문일자 */
														, B.ORDER_NO			
														, B.REST_ID
														, C.REST_NM
														, B.CREDIT_AMT		
														, B.ORDER_STAT 
														, PAY_TY
												   		, D.ITEM_NM 
                                          				, NOMAL_SALE_PRICE
														, NOMAL_SALE_VAT
														, SALE_PRICE
														, SALE_VAT
														, TOTAL_PRICE
														, CP_STATUS
								FROM b_rest_combine_sub AS a
								INNER JOIN dar_order_info AS b ON a.sub_rest_id = b.rest_id
								INNER JOIN dar_order_detail AS BB ON B.ORDER_NO = BB.ORDER_NO  
								INNER JOIN dar_order_coupon AS cc ON BB.ORDER_NO = cc.ORDER_NO
								INNER JOIN priv_rest_info AS c ON b.REST_ID = c.rest_id
								INNER JOIN dar_sale_item_info AS D ON C.REST_ID = D.REST_ID AND BB.ITEM_NO = D.ITEM_NO
								WHERE 
									A.rest_id='#{restId}'
								AND B.ORDER_STAT='20'
								AND DATE_FORMAT(B.ORDER_DATE, '%Y-%m-%d') >=  '#{startDate}'
								AND DATE_FORMAT(B.ORDER_DATE, '%Y-%m-%d') <=  '#{endDate}'
								AND B.REST_ID ='#{subRestId}' 

								ORDER BY B.ORDER_DATE DESC
											`

var SelectCombineStoreData = `	SELECT A.REST_ID
										,A.CHARGE_AMT
										,B.REST_NM
										,BANK_CD
										,ACCOUNT_NO
										,ACCOUNT_NM
										,ACCOUNT_CERT_YN
								FROM b_rest_combine AS A
								INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
								WHERE
									A.rest_id='#{restId}'
						`

var SelectCombinePaidList = `SELECT	
				 REST_PAYMENT_ID  						AS restPaymentId  		-- 가맹점 정산 아이디 
				,B.REST_NM 				 				AS restNm 				-- 가맹점명
				,B.BUSID				 				AS busId
				,DATE_FORMAT(A.PAYMENT_DT,'%Y-%m-%d')   AS paymentDt  			-- 결제일자 
				,DATE_FORMAT(A.SETTLMNT_DT,'%Y-%m-%d')	 	 				AS settlmntDt 	-- 정산 지급 요청 일자 
				,A.PAYMENT_CNT           				AS paymentCnt   		-- 결제 건수
				,A.PAYMENT_AMT           				AS paymentAmt			-- 결제 금액
				,A.REST_PAYMENT_AMT      				AS restPaymentAmt  	-- 지급 요청 금액
				,A.FIT_FEE               				AS fitFee				-- 총수수료
				,A.PG_FEE                				AS pgFee					-- FIT 수수료
				,A.TOT_FEE               				AS totFee				-- pg 수수료 
				,A.CANCEL_AMT            				AS cancelAmt			-- 취소 금액
				,A.CANCEL_CNT            				AS cancelCnt			-- 취소 건수 
				,IFNULL(A.RESULT_MSG,'')   				AS resultMsg      -- 결과 메세지
				,DATE_FORMAT(CONCAT(A.SEND_DATE,A.SEND_TIME),'%Y-%m-%d %h:%i:%s') as sendDate
		FROM DAR_REST_PAYMENT AS A
		INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
		WHERE
		A.rest_id='#{restId}'
		AND A.RESULT_PAY_YN='N'
		ORDER BY A.SEND_DATE DESC,A.PAYMENT_DT ASC
`

var SelectPaymentListCnt = `SELECT COUNT(*) AS TOTAL_COUNT
				FROM dar_payment_hist AS a
				INNER JOIN dar_payment_report AS AA ON AA.MOID = A.MOID
				INNER JOIN priv_rest_info AS b ON a.REST_ID = b.rest_id
				INNER JOIN priv_grp_info AS C ON A.GRP_ID = C.GRP_ID
				INNER JOIN priv_user_info AS D ON A.USER_ID = D.USER_ID
				WHERE 
					DATE_FORMAT(AA.PAYMENT_DT, '%Y-%m-%d') >=  '#{startDate}'
					AND DATE_FORMAT(AA.PAYMENT_DT, '%Y-%m-%d') <=  '#{endDate}'
					AND C.GRP_NM LIKE '%#{searchGrpNm}%'
					AND D.USER_NM LIKE '%#{searchUserNm}%'
					AND B.REST_NM LIKE '%#{searchRestNm}%'
					AND A.PAYMENT_TY = '#{searchStat}'
`

var SelectPaymentList = `SELECT A.MOID
						,B.REST_NM
						,B.REST_ID
						,A.GRP_ID
						,C.GRP_NM
						,A.USER_ID
						,IFNULL(D.USER_NM,'') AS USER_NM
						,A.USER_TY
						,A.PAY_CHANNEL
						,A.PAYMENT_TY
						,ACC_ST_DAY
						,AA.FNNAME
						,CASE WHEN PAY_INFO IN('1','3') THEN A.CREDIT_AMT ELSE AA.AMT END AS AMT
						,IFNULL(DATE_FORMAT(CONCAT(AA.PAYMENT_DT,AA.PAYMENT_TM),'%Y-%m-%d %H:%i:%s'), DATE_FORMAT(A.REG_DATE,'%Y-%m-%d %H:%i:%s')) as PAYMENT_DT
						,A.PAY_INFO
				FROM dar_payment_hist AS a
				LEFT OUTER JOIN dar_payment_report AS AA ON AA.MOID = A.MOID
				LEFT OUTER JOIN priv_rest_info AS b ON a.REST_ID = b.rest_id
				LEFT OUTER JOIN priv_grp_info AS C ON A.GRP_ID = C.GRP_ID
				LEFT OUTER JOIN priv_user_info AS D ON A.USER_ID = D.USER_ID
				WHERE 
					DATE_FORMAT(A.REG_DATE, '%Y-%m-%d') >=  '#{startDate}'
					AND DATE_FORMAT(A.REG_DATE, '%Y-%m-%d') <=  '#{endDate}'
					AND C.GRP_NM LIKE '%#{searchGrpNm}%'
					AND D.USER_NM LIKE '%#{searchUserNm}%'
					AND B.REST_NM LIKE '%#{searchRestNm}%'
					AND A.PAYMENT_TY = '#{searchStat}'
				ORDER BY A.REG_DATE DESC
`




var InsertTpayUnpaidPaymentHistory = `INSERT INTO DAR_PAYMENT_HIST
								(
										  HIST_ID
										,REST_ID
										,GRP_ID
										,USER_ID
										,MOID
										,CREDIT_AMT
										,ADD_AMT
										,USER_TY
										,SEARCH_TY
										,PAYMENT_TY
										,PAY_INFO
										,REG_DATE
										,PG_CD
										,PAY_CHANNEL
								)
								SELECT   CONCAT('#{histId}','_',REST_ID)
										,REST_ID
										,GRP_ID
										, '#{userId}'
										, '#{moid}'
										, SUM(TOTAL_AMT)
										, '#{addAmt}'
										, '#{userTy}'
										, '#{searchTy}'
										, '#{paymentTy}'
										, '#{payInfo}'
										, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										, '#{pgCd}'
										, '#{payChannel}'     
								FROM DAR_ORDER_INFO
								WHERE 
								PAID_KEY='#{moid}'
								GROUP BY REST_ID,GRP_ID
					`


var SelectTpayBillingCardInfo string = `SELECT A.USER_ID
										, SEQ
										, CARD_NAME
										, CARD_CODE
										, CARD_NUM
										, CARD_TOKEN
										, CARD_TYPE
										, A.USE_YN
										, EMAIL
										, HP_NO
										, USER_NM
										, FN_GET_RESTNAME('#{restId}') AS REST_NM
										, IFNULL(BILLING_PWD,'NONE') AS BILLING_PWD
									FROM b_tpay_billing_key AS A
									INNER JOIN priv_user_info AS b ON a.user_id= b.user_id
									WHERE
									A.USER_ID='#{userId}'
									AND A.SEQ='#{seq}'
										`
