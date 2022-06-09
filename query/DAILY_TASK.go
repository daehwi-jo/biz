package query

var UpdateOrderCanCel string = ` UPDATE DAR_ORDER_INFO
							    SET ORDER_STAT = '#{orderStat}'
								WHERE 
								ORDER_NO = '#{orderNo}'
							 `

var UpdateGiftCancel string = ` UPDATE DAR_GIFT_INFO
							    SET GIFT_CAN_DTM = DATE_FORMAT(SYSDATE(), '%Y%m%d%H%i%s')
								, GIFT_STS_CD='#{giftStsCd}'
								WHERE 
								MOID = '#{giftMoid}'
								AND GIFT_STS_CD = '1'
								AND RCV_STS_CD IN ('0', '2')
							 `
var SelectAgrmInFo string = `SELECT AGRM_ID
							,PREPAID_AMT
							FROM ORG_AGRM_INFO 
							WHERE 
							GRP_ID = '#{grpId}'
							AND REST_ID ='#{restId}'
							`

var UpdateAgrm string = ` UPDATE ORG_AGRM_INFO
							    SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								,PREPAID_AMT = '#{prepaidAmt}'
								WHERE 
								AGRM_ID = '#{agrmId}'
							 `

var SelectGiftInFo string = `SELECT GIFT_STS_CD
									,ORDER_NO
									,RCV_STS_CD
							FROM dar_gift_info
							WHERE 
							MOID='#{giftMoid}'
							`
var SelectDailyResult string = `
							SELECT 
								task_succ_yn 
							FROM 
								b_dailytask 
							WHERE 
								left(reg_date,8) = DATE_FORMAT(NOW(), '%Y%m%d')
`

var SelectUnusedGiftList string = `SELECT  
										   G.SND_GRP_ID 
										 , G.REST_ID 
										 , G.SND_USER_ID 
										 , G.GIFT_AMT 
										 , G.MOID
								FROM       DAR_GIFT_INFO G 
								WHERE      1=1
								AND        GIFT_STS_CD = '1'
								AND        RCV_STS_CD  = '0'
								AND        DATE_FORMAT(NOW(), '%Y%m%d') > DATE_FORMAT(DATE_ADD(GIFT_SND_DTM, INTERVAL '#{cancelDur}' DAY), '%Y%m%d')
								`

var InsertBackupSupportBalance string = `INSERT INTO   PRIV_GRP_USER_BALANCE
											(
											GRP_ID
											, USER_ID
											, APLY_DATE
											, SUPPORT_BALANCE
											, REG_DATE
											, REG_ID
											)
											SELECT
											A.GRP_ID
											, A.USER_ID
											, DATE_FORMAT(NOW(), '%Y%m%d')
											, A.SUPPORT_BALANCE
											, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
											, 'SYSTEM'
											
											FROM      PRIV_GRP_USER_INFO  A
											, PRIV_GRP_INFO       B
											WHERE     1=1
											AND       A.GRP_ID      = B.GRP_ID
											AND       A.AUTH_STAT   = '1'
											AND       B.SUPPORT_YN = 'Y'
								`

var UpdateResetSupportBalance string = `UPDATE   PRIV_GRP_USER_INFO A
									,(
										SELECT  GRP_ID
										, GRP_NM
										, SUPPORT_AMT
										, SUPPORT_FORWARD_YN
										FROM    PRIV_GRP_INFO
										WHERE   1=1
										AND     SUPPORT_YN = 'Y'
										AND     AUTH_STAT = '1'
										) B
										SET        A.MOD_DATE        =  DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										,  A.SUPPORT_BALANCE =  CASE WHEN B.SUPPORT_FORWARD_YN = 'Y'
										THEN A.SUPPORT_BALANCE + B.SUPPORT_AMT
										ELSE B.SUPPORT_AMT
										END
										WHERE      1=1
										AND        A.GRP_ID = B.GRP_ID
										AND        A.AUTH_STAT = '1'
										`

var SelectDayBefore string = `SELECT   
								DATE_FORMAT(DATE_ADD(NOW(), INTERVAL -'#{dayBefore}' DAY), '%Y%m%d') as beforeDate`

var SelectBizDay string = `SELECT    
									   A.total_date as totalDate
							FROM      (
										SELECT     
												   total_date
												 , @ROWNUM:=@ROWNUM+1 As RN
										FROM       sys_week_date  
												 , (Select @ROWNUM:=0)ROWNUM
										WHERE      1=1
										AND        HOLIDAY = 'WD'
										AND        total_date > '#{paymentDt}'
										AND        total_date < DATE_FORMAT(DATE_ADD('#{paymentDt}', INTERVAL 30 DAY), '%Y%m%d')
										ORDER BY   total_date
									  ) A
							WHERE     A.RN = 2
							`

var UpdateMakePaymentId string = `UPDATE     DAR_PAYMENT_REPORT A
								, (
									SELECT    MOID
									, MAX(PAYMENT_TY)       AS PAYMENT_TY
									, MAX(REST_PAYMENT_ID)  AS REST_PAYMENT_ID
									FROM      DAR_PAYMENT_HIST
									WHERE     1=1
									AND       PAY_CHANNEL  IN ('02', '03')       
									GROUP BY  MOID
								) B
								, DAR_PAYMENT_HIST C
								, PRIV_REST_INFO D
								SET  C.REST_PAYMENT_ID = CONCAT(A.PAYMENT_DT, '_', C.REST_ID)
									,C.REST_PAYMENT_ID_DT = DATE_FORMAT(NOW(), '%Y%m%d') 
								WHERE      1=1
								AND        A.MOID       = B.MOID
								AND        A.MOID       = C.MOID
								AND 	   C.REST_ID    = D.REST_ID 
								AND D.PAY_PROXY_CD='0000'
								AND        A.STATECD    = '0'        	
								AND        A.STATE      = '0000'    	
								AND        C.PAYMENT_TY IN ('0', '3')   
								AND 	   C.REST_PAYMENT_ID ='X'
								AND D.REST_TYPE='N'
								AND A.PAYMENT_DT  < DATE_FORMAT(NOW(), '%Y%m%d')
							`

var UpdatePaymentFees string = `UPDATE dar_payment_report
								SET
									REST_PAY_AMT= '#{restPayAmt}',
									TOT_SUPLY_AMT= '#{totSuplyAmt}',
									TOT_VAT='#{totVat}',
									TOT_FEE='#{totFee}',
									PG_SUPLY_AMT='#{pgSuplyAmt}',
									PG_VAT='#{pgVat}',
									PG_FEE='#{pgFee}',
									FIT_SUPLY_AMT='#{fitSuplyAmt}',
									FIT_VAT='#{fitVat}',
									FIT_FEE='#{fitFee}'
								WHERE 
								moid = '#{moid}'

`

var UpdatePaymentHistFees string = `UPDATE dar_payment_hist
								SET
									REST_PAY_AMT= '#{restPayAmt}',
									TOT_SUPLY_AMT= '#{totSuplyAmt}',
									TOT_VAT='#{totVat}',
									TOT_FEE='#{totFee}',
									PG_SUPLY_AMT='#{pgSuplyAmt}',
									PG_VAT='#{pgVat}',
									PG_FEE='#{pgFee}',
									FIT_SUPLY_AMT='#{fitSuplyAmt}',
									FIT_VAT='#{fitVat}',
									FIT_FEE='#{fitFee}'
								WHERE 
								HIST_ID = '#{histId}'

`

var UpdateMakePayment string = `INSERT INTO DAR_REST_PAYMENT
								(
								REST_PAYMENT_ID   /* 가맹점 정산ID          	*/
								, REST_ID           /* 가맹점ID 	           	*/
								, REST_NM           /* 가맹점명 	           	*/
								, PAYMENT_DT        /* 지급대상 결제일자                */
								, PAYMENT_CNT       /* 지급대상 결제건수                */
								, PAYMENT_AMT       /* 지급대상 결제금액                */
								, REST_PAYMENT_AMT  /* 가맹점 실지급액(수수료 제외금액) */
								, TOT_FEE           /* 총수수료                         */
								, PG_FEE            /* PG사 수수료                      */
								, FIT_FEE           /* FIT 수수료                      	*/
								, CANCEL_CNT        /* 취소건수                        	*/
								, CANCEL_AMT        /* 취소금액                       	*/
								, SETTLMNT_DT       /* 정산지급 요청일자 YYYYMMDD 	*/
								, WORK_DATE         /* 작업일자                         */
								, WORK_TIME         /* 작업시간                         */
								)
								SELECT
								B.REST_PAYMENT_ID   			/* 가맹점 정산ID	*/
								, B.REST_ID           			/* 가맹점ID   		*/
								, C.REST_NM					/* 가맹점명   		*/
								, A.PAYMENT_DT					/* 결제일자   		*/
								, SUM(
								CASE WHEN A.STATECD = '0' THEN 1
								ELSE 0
								END
								)  AS PAYMENT_CNT  		/* 지급대상 결제건수 */
								, SUM(
								CASE WHEN A.STATECD = '0' THEN B.CREDIT_AMT
								ELSE 0
								END
								)  AS PAYMENT_AMT 		/* 지급대상 결제금액 */
								, SUM(
								CASE WHEN A.STATECD = '0' THEN B.REST_PAY_AMT
								ELSE 0
								END
								)  AS REST_PAYMENT_AMT	/* 가맹점 실지급액(수수료 제외금액) */
								, SUM(
								CASE WHEN A.STATECD = '0' THEN B.TOT_FEE
								ELSE 0
								END
								)  AS TOT_FEE 			/* 총수수료  	*/
								, SUM(
								CASE WHEN A.STATECD = '0' THEN B.PG_FEE
								ELSE 0
								END
								)  AS PG_FEE				/* PG사 수수료  */
								, SUM(
								CASE WHEN A.STATECD = '0' THEN B.FIT_FEE
								ELSE 0
								END
								)  AS FIT_FEE				/* FIT 수수료   */
								, SUM(
								CASE WHEN A.STATECD = '1' THEN 1
								ELSE 0
								END
								)  AS CANCEL_CNT			/* 취소건수     */
								, SUM(
								CASE WHEN A.STATECD = '1' THEN A.AMT
								ELSE 0
								END
								)  AS CANCEL_AMT			/* 취소금액     */
								, '#{settlmtDt}'
								, DATE_FORMAT(NOW(), '%Y%m%d')
								, DATE_FORMAT(NOW(), '%H%i%s')
								FROM       DAR_PAYMENT_REPORT A
								, DAR_PAYMENT_HIST   B
								, PRIV_REST_INFO     C
								WHERE      1=1
								AND        A.MOID       = B.MOID
								AND        B.REST_ID    = C.REST_ID
								AND        B.REST_PAYMENT_ID != 'X'
								AND        B.REST_PAYMENT_ID_DT = DATE_FORMAT(NOW(), '%Y%m%d')
								GROUP BY   B.REST_ID
								, B.REST_PAYMENT_ID
								, A.PAYMENT_DT
								`

var InsertDailyTask string = ` INSERT INTO b_dailytask
								(
								  TASK_NAME
								, TASK_SUCC_YN
								, CYCLE
								, REG_DATE
								)
								VALUES (
								'#{taskName}'
								, '#{taskSuccYn}'
								, '#{cycle}'
								, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								)
								`
var UpdateDailyTask string = `  UPDATE b_dailytask SET 
								TASK_SUCC_YN = '#{taskSuccYn}'
								WHERE 
								SEQ = '#{seq}'
								`

var SelectMakeList = `SELECT 
								 A.MOID
								,A.PAYMETHOD
								,B.CREDIT_AMT as AMT
								,B.REST_ID
								,PAYMENT_DT
								,B.HIST_ID
						FROM dar_payment_report AS A
						INNER JOIN dar_payment_hist AS B ON A.MOID = B.MOID AND  B.REST_PAYMENT_ID != 'X'  
						INNER JOIN priv_rest_info AS C ON B.REST_ID = C.REST_ID
						WHERE 
						B.REST_PAYMENT_ID_DT = DATE_FORMAT(NOW(), '%Y%m%d')
						AND   B.REST_PAYMENT_ID != 'X'
					`

var SelectFeeMakeList = `SELECT 
								 A.MOID
								,A.PAYMETHOD
								,AMT
								,B.REST_ID
								,PAYMENT_DT
								,B.HIST_ID
						FROM dar_payment_report AS A
						INNER JOIN dar_payment_hist AS B ON A.MOID = B.MOID
						WHERE
						A.REST_PAY_AMT IS null
						AND  A.STATECD ='0' 
						AND  A.STATE = '0000'
						AND  B.PAYMENT_TY IN ('0', '3')
					`

var SelectKakaoAlimYN string = `
							SELECT kakao_week, kakao_month, kakao_daily 
							FROM priv_rest_etc 
							WHERE 
								rest_id = '#{restId}'
`

var SelectSuccessDailyAlim string = `
							SELECT a.BIZ_NUM as bizNum, a.REST_ID as restId, b.USER_ID as userId 
							FROM cc_comp_inf a 
							LEFT JOIN priv_rest_user_info b ON a.REST_ID = b.REST_ID 
							WHERE b.REST_AUTH = 0 AND a.biz_num IN(
								SELECT biz_num 
								FROM cc_sync_inf 
								WHERE 
									bs_dt = '#{bsDt}' 
								AND err_cd = '0000' 
								AND biz_num IN(
									SELECT b.BIZ_NUM 
									FROM priv_rest_info a INNER JOIN cc_comp_inf b ON a.REST_ID = b.REST_ID 
									AND a.USE_YN ='Y'
							))
`

var SelectBillingInfo string = `
							SELECT CASE WHEN DATE_FORMAT(start_date,'%Y%m%d') <= DATE_FORMAT(NOW(),'%Y%m%d') THEN 'Y' ELSE 'N' END AS startYN, CASE WHEN DATE_FORMAT(end_date,'%Y%m%d') >= DATE_FORMAT(NOW(),'%Y%m%d') THEN 'Y' ELSE 'N' END AS endYN
							FROM e_billing 
							WHERE 
								user_id='#{userId}'
`

var SelectSuccessReportSendCheck string = `
							SELECT result 
							FROM sys_alimtalk_log  
							WHERE 
								template_code = '#{code}' 
							AND 
								left(send_date,8) = '#{today}' 
							AND 
								user_id='#{userId}'
`

var SelectSyncAmt string = `
							SELECT bs_dt, aprv_amt AS amt 
							FROM cc_sync_inf 
							WHERE 
								biz_num = '#{bizNum}' 
							AND (
								bs_dt = '#{bsDt}' 
							or 
								bs_dt = '#{bsDt2}'
							) GROUP BY bs_dt 
							ORDER BY bs_dt
`

var SelectSuccessWeekAlim string = `
							SELECT a.BIZ_NUM as bizNum, a.REST_ID as restId, b.USER_ID as userId 
							FROM cc_comp_inf a 
							LEFT JOIN priv_rest_user_info b ON a.REST_ID = b.REST_ID 
							WHERE b.REST_AUTH = 0 AND a.biz_num IN(
								SELECT biz_num 
								FROM cc_sync_inf 
								WHERE bs_dt 
									BETWEEN 
										'#{startDt}' 
									AND 
										'#{endDt}'
								AND err_cd != '0005' 
								AND biz_num IN(
									SELECT b.BIZ_NUM 
									FROM priv_rest_info a INNER JOIN cc_comp_inf b ON a.REST_ID = b.REST_ID 
									AND a.USE_YN ='Y'
							))
`

var SelectSuccessMonthAlim string = `
							SELECT a.BIZ_NUM as bizNum, a.REST_ID as restId, b.USER_ID as userId 
							FROM cc_comp_inf a 
							LEFT JOIN priv_rest_user_info b ON a.REST_ID = b.REST_ID 
							WHERE b.REST_AUTH = 0 AND a.biz_num IN(
								SELECT biz_num 
								FROM cc_sync_inf 
								WHERE 
									left(bs_dt,6) = '#{bsDt}' 
								AND err_cd != '0005' 
								AND biz_num IN(
									SELECT b.BIZ_NUM 
									FROM priv_rest_info a INNER JOIN cc_comp_inf b ON a.REST_ID = b.REST_ID 
									AND a.USE_YN ='Y'
							))
`

var SelectFailDailyAlim string = `
							SELECT a.BIZ_NUM as bizNum, a.REST_ID as restId, b.USER_ID as userId, a.COMP_NM as compNm 
							FROM cc_comp_inf a 
							LEFT JOIN priv_rest_user_info b ON a.REST_ID = b.REST_ID 
							WHERE b.REST_AUTH = 0 AND a.biz_num IN(
								SELECT biz_num 
								FROM cc_sync_inf 
								WHERE 
									bs_dt = '#{bsDt}' 
								AND err_cd != '0000' 
								AND err_cd != '0005' 
								AND biz_num IN(
									SELECT b.BIZ_NUM 
									FROM priv_rest_info a INNER JOIN cc_comp_inf b 
									ON a.REST_ID = b.REST_ID 
									AND a.USE_YN ='Y'
							))
`

var SelectFailDailyAlimCard string = `
							SELECT a.BIZ_NUM as bizNum, a.REST_ID as restId, b.USER_ID as userId, a.COMP_NM as compNm 
							FROM cc_comp_inf a 
							LEFT JOIN priv_rest_user_info b ON a.REST_ID = b.REST_ID 
							WHERE b.REST_AUTH = 0 AND a.biz_num IN(
								SELECT biz_num 
								FROM cc_sync_inf 
								WHERE 
									bs_dt = '#{bsDt}' 
								AND err_cd = '0005'
								AND biz_num IN(
									SELECT b.BIZ_NUM 
									FROM priv_rest_info a INNER JOIN cc_comp_inf b 
									ON a.REST_ID = b.REST_ID 
									AND a.USE_YN ='Y'
							))
`

var SelectSucessDailyAlimCard string = `
							SELECT a.BIZ_NUM as bizNum, a.REST_ID as restId, b.USER_ID as userId, a.COMP_NM as compNm 
							FROM cc_comp_inf a 
							LEFT JOIN priv_rest_user_info b ON a.REST_ID = b.REST_ID 
							WHERE b.REST_AUTH = 0
							AND
								a.REST_ID = '#{restId}'
							AND a.biz_num IN(
								SELECT biz_num 
								FROM cc_sync_inf 
								WHERE 
									bs_dt = '#{yesterDt}' 
								AND err_cd = '0000'
								AND biz_num IN(
									SELECT b.BIZ_NUM 
									FROM priv_rest_info a INNER JOIN cc_comp_inf b 
									ON a.REST_ID = b.REST_ID 
									AND a.USE_YN ='Y'
							))
`

// 입금캘린더 월별 합계리스트
var SelectPayCalendarSumList string = `
									SELECT 
										tr_month AS trMonth, 
										SUM(z.outp_expt_amt) AS outpExptAmt, 
										SUM(z.real_in_amt) AS realInAmt,
										SUM(z.real_in_amt) - SUM(Z.outp_expt_amt) AS diffAmt,
										CASE WHEN SUM(z.outp_expt_amt) = SUM(z.real_in_amt) THEN '0' 
											WHEN SUM(z.outp_expt_amt) > SUM(z.real_in_amt) THEN '1' ELSE '2' END diffColor 
									FROM (
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											IFNULL(SUM(pay_amt),0) AS outp_expt_amt,
											0 AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pca_dtl ON dt = outp_expt_dt 
										AND 
											biz_num = '#{bizNum}'
										WHERE 
											dt BETWEEN '#{startDt}' 
											AND '#{endDt}' 
										GROUP BY SUBSTR(dt, 1, 6)
										
										UNION ALL
										
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											0 AS outp_expt_amt,
											IFNULL(SUM(real_pay_amt),0) AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pay_dtl ON dt = pay_dt
										AND 
											biz_num = '#{bizNum}'
										WHERE 
											dt BETWEEN '#{startDt}' 
											AND '#{endDt}' 
										GROUP BY SUBSTR(dt, 1, 6)
									) z
									GROUP BY z.tr_month
									ORDER BY z.tr_month DESC
									`

// 입금 캘린더
var SelectPayCalendarList string = `
									SELECT 
										CAST(FORMAT(@RN := @RN + 1, 0) as unsigned) AS rNum,
										tr_dt AS trDt, 
										SUM(z.outp_expt_amt) AS outpExptAmt, 
										SUM(z.real_in_amt) AS realInAmt, 
										SUM(z.real_in_amt) - SUM(z.outp_expt_amt) AS diffAmt,
										CASE WHEN SUM(z.outp_expt_amt) = SUM(z.real_in_amt) THEN '0' 
											WHEN SUM(z.outp_expt_amt) > SUM(z.real_in_amt) THEN '1' ELSE '2' END diffColor,
										SUM(day_color) AS dayColor
									FROM (
										SELECT 
											a.dt AS tr_dt, 
											IF(a.dt >= DATE_FORMAT(NOW(), '%Y%m%d'), 0, IFNULL(SUM(b.pay_amt),0)) AS outp_expt_amt, 
											0 AS real_in_amt,
											0 AS day_color
										FROM 
											cc_date_info a 
											LEFT JOIN cc_pca_dtl b 
												ON a.dt = b.outp_expt_dt 
										AND 
											b.biz_num = '#{bizNum}'
										WHERE 
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt

										UNION ALL

										SELECT total_date as tr_dt,
										0 as outp_expt_amt,
										0 as real_in_amt,
										CASE WHEN DAY = 7 then 2
										WHEN datekind = 'W' then 1
										ELSE 3 END day_color 
										FROM 
											sys_week_date 
										WHERE 
											total_date 
										BETWEEN '#{startDt}' 
										AND '#{endDt}'
										
										UNION ALL
										
										SELECT 
											a.dt AS tr_dt, 
											0 AS outp_expt_amt,
											IFNULL(SUM(c.real_pay_amt),0) AS real_in_amt,
											0 AS day_color
										FROM 
											cc_date_info a 
											LEFT JOIN cc_pay_dtl c 
												ON a.dt = c.pay_dt 
										AND 
											c.biz_num = '#{bizNum}'
										WHERE  
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt
									) z INNER JOIN (SELECT @RN := 0) r
									GROUP BY z.tr_dt
									ORDER BY z.tr_dt
									`

// 매출캘린더 월별 합계리스트
var SelectAprvCalendarSumList string = `
									SELECT
										tr_month AS trMonth,
										SUM(z.aprv_amt) AS aprvAmt,
										SUM(z.cash_amt) AS cashAmt,
										SUM(z.pca_amt) AS pcaAmt,
										SUM(z.tot_amt) AS totAmt
									FROM (
										SELECT
											SUBSTR(a.TOTAL_DATE, 1, 6) AS tr_month,
                                            IFNULL(sum(b.TOT_AMT),0) AS aprv_amt,
                                            0 AS cash_amt,
                                            0 AS pca_amt,
                                            IFNULL(sum(b.TOT_AMT),0) AS tot_amt
                                        FROM
                                            sys_week_date a
                                            LEFT JOIN cc_aprv_sum b ON a.TOTAL_DATE = b.BS_DT 
										AND 
											b.biz_num = '#{bizNum}'
										WHERE
											a.TOTAL_DATE BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY SUBSTR(a.TOTAL_DATE, 1, 6)
									
										UNION ALL
									
										SELECT
											SUBSTR(dt, 1, 6) AS tr_month,
											0 AS aprv_amt,
											0 AS cash_amt,
											IFNULL(SUM(pca_amt),0) AS pca_amt,
											0 AS tot_amt
										FROM
											cc_date_info
											LEFT JOIN cc_pca_dtl ON dt = org_tr_dt 
										AND 
											biz_num = '#{bizNum}'
										WHERE
											dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)
									
										UNION ALL
									
										SELECT
											SUBSTR(DT, 1, 6) AS tr_month,
											0 AS aprv_amt,
											IFNULL(SUM(tot_amt),0) AS cash_amt,
											0 AS pca_amt,
											IFNULL(SUM(tot_amt),0) AS tot_amt
										FROM
											cc_date_info
											LEFT JOIN cc_cash_dtl ON dt = tr_dt 
										AND 
											biz_num = '#{bizNum}'
										WHERE
											dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)
									) z
									GROUP BY z.tr_month
									ORDER BY z.tr_month DESC
									`

var SelectAprvCalendarList string = `
									SELECT 
										CAST(FORMAT(@RN := @RN + 1, 0) as unsigned) AS rNum, 
										tr_dt AS trDt, 
										SUM(z.aprv_amt) AS aprvAmt, 
										SUM(z.cash_amt) AS cashAmt, 
										SUM(z.pca_amt) AS pcaAmt, 
										SUM(z.tot_amt) AS totAmt,
										CASE WHEN SUM(z.pca_amt) = 0 THEN '3' WHEN SUM(z.aprv_amt) > SUM(z.pca_amt) THEN '1' WHEN SUM(z.aprv_amt) < SUM(z.pca_amt) THEN '2' ELSE '0' END diffColor,
										SUM(day_color) AS dayColor
									FROM (
										SELECT 
											a.dt AS tr_dt, 
											IFNULL(SUM(b.tot_amt),0) AS aprv_amt, 
											0 AS cash_amt, 
											0 AS pca_amt, 
											IFNULL(SUM(b.tot_amt),0) AS tot_amt,
											0 as day_color
										FROM 
											tb_date_info a 
											LEFT JOIN cc_aprv_sum b 
												ON a.dt = b.bs_dt 
										AND 
											b.biz_num = '#{bizNum}'
										WHERE 
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt

										UNION ALL

										SELECT total_date as tr_dt,
										0 as aprv_amt,
										0 as cash_amt,
										0 as pca_amt,
										0 as tot_amt,
										CASE WHEN DAY = 7 then 2
										WHEN datekind = 'W' then 1
										ELSE 3 END day_color 
										FROM 
											sys_week_date 
										WHERE 
											total_date 
										BETWEEN '#{startDt}' 
										AND '#{endDt}'

										UNION ALL

										SELECT 
											a.dt AS tr_dt, 
											0 AS aprv_amt, 
											0 AS cash_amt, 
											IFNULL(SUM(b.pca_amt),0) AS pca_amt, 
											0 AS tot_amt,
											0 as day_color
										FROM 
											tb_date_info a 
											LEFT JOIN cc_pca_dtl b 
												ON a.dt = b.org_tr_dt 
										AND 
											b.biz_num = '#{bizNum}'
										WHERE  
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt
										
										UNION ALL
										
										SELECT 
											a.dt AS tr_dt, 
											0 AS aprv_amt, 
											IFNULL(SUM(c.tot_amt),0) AS cash_amt, 
											0 AS pca_amt, 
											IFNULL(SUM(c.tot_amt),0) AS tot_amt,
											0 as day_color
										FROM 
											tb_date_info a 
											LEFT JOIN cc_cash_dtl c 
												ON a.dt = c.tr_dt 
										AND 
											c.biz_num = '#{bizNum}'
										WHERE  
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt
										) z INNER JOIN (SELECT @RN := 0) R
										GROUP BY z.tr_dt
										ORDER BY z.tr_dt
									`

var SelectRedisComp string = `
							SELECT biz_num 
							FROM 
								cc_comp_inf 
							WHERE 
								comp_sts_cd = 1 
							AND 
								LENGTH(svc_open_dt) > 0 
							AND 
								ln_first_yn = 'Y'
							`

var InserCombineDarRestPayment string = `INSERT INTO DAR_REST_PAYMENT
								(
								REST_PAYMENT_ID   /* 가맹점 정산ID          	*/
								, REST_ID           /* 가맹점ID 	           	*/
								, REST_NM           /* 가맹점명 	           	*/
								, PAYMENT_DT        /* 지급대상 결제일자                */
								, PAYMENT_CNT       /* 지급대상 결제건수                */
								, PAYMENT_AMT       /* 지급대상 결제금액                */
								, REST_PAYMENT_AMT  /* 가맹점 실지급액(수수료 제외금액) */
								, TOT_FEE           /* 총수수료                         */
								, PG_FEE            /* PG사 수수료                      */
								, FIT_FEE           /* FIT 수수료                      	*/
								, CANCEL_CNT        /* 취소건수                        	*/
								, CANCEL_AMT        /* 취소금액                       	*/
								, SETTLMNT_DT       /* 정산지급 요청일자 YYYYMMDD 	*/
								, WORK_DATE         /* 작업일자                         */
								, WORK_TIME         /* 작업시간                         */
								)
						VALUE(
								 CONCAT(DATE_FORMAT(NOW(), '%Y%m%d'), '_','#{restId}')
								 ,'#{restId}'
								 ,'#{restNm}'
								 ,DATE_FORMAT(NOW(), '%Y%m%d')
								 ,1
								 ,'#{amt}'
								 ,'#{restPayAmt}'
								 ,'#{totFee}'
								 ,'#{pgFee}'
								 ,'#{fitFee}'
								 ,0
								 ,0
								 ,'#{settlmtDt}'
								 , DATE_FORMAT(NOW(), '%Y%m%d')
								 , DATE_FORMAT(NOW(), '%H%i%s')
						)
`

// 기업 기념일 정보
var SelectCompanyDayInfo string = `
						SELECT 
							company_id AS cId, 
							event_code AS eCode, 
							event_name AS eName,
							SHIP_DATE AS sDate, 
							PRICE as price
						FROM 
							b_company_day
						WHERE
							use_yn = 'Y'
`

var SelectCompanyDayDetailInfo string = `
						SELECT 
							company_id AS cId, 
							event_code AS eCode, 
							event_name AS eName,
							SHIP_DATE AS sDate, 
							PRICE as price
						FROM 
							b_company_day
						WHERE
							use_yn = 'Y'
						AND
							company_id = '#{cId}'
`

// 기업 기념일 생일 직원 정보
var SelectCompanyBirthdayUserInfo string = `
						SELECT 
							user_id AS uId, 
							book_id as bId,
							user_nm AS uName, 
							hp_no AS hp, 
							user_birth AS uBirth, 
							LUNAR_BIRTH_YN AS lunarYn 
						FROM 
							b_company_user 
						WHERE 
							company_id = '#{cId}' 
						AND (
								substr(USER_BIRTH, 5, 4) = '#{shipDate}' 
								OR
								( substr(USER_BIRTH, 5, 4) = '#{shipDateLunar}' AND lunar_birth_yn = 'Y' ) )
						AND 
							use_yn = 'Y'
`

var SelectCompanyBirthdayUserDetailInfo string = `
						SELECT 
							user_id AS uId, 
							book_id as bId,
							user_nm AS uName, 
							hp_no AS hp, 
							user_birth AS uBirth, 
							LUNAR_BIRTH_YN AS lunarYn 
						FROM 
							b_company_user 
						WHERE 
							company_id = '#{cId}'
						AND 
							use_yn = 'Y'
						AND
							user_id = '#{uId}'
`

// 기업 잔여 지원금 히스토리 조회
var SelectCompanyDayHis string = `
						SELECT 
							price, 
							price_date AS pDate,
							alarm_yn AS aYn
						FROM 
							b_company_day_his 
						WHERE 
							company_id=#{cId}
						AND 
							event_code='#{eCode}'
						AND 
							user_id='#{uId}' 
						AND 
							substr(PRICE_DATE, 1, 4)='#{thisYear}'
`

// 기업 잔여 지원금 정보 업데이트
var UpdateCompanyUserDetail string = ` 
						UPDATE priv_grp_user_info
						SET 
							support_balance = support_balance + '#{cPrice}'
						WHERE 
							grp_id='#{bId}' 
						AND 
							user_id='#{uId}'
						`

// 기업 잔여 지원금 히스토리
var InsertCompanyDayHis string = `INSERT INTO b_company_day_his
								(
								COMPANY_ID
								, EVENT_CODE         
								, USER_ID           
								, USER_NM        
								, HP_NO       
								, PRICE       
								, PRICE_DATE  
								, SHIP_DATE           
								)
						VALUE(
								 '#{cId}' 
								 ,'#{eCode}'
								 ,'#{uId}'
								 ,'#{uName}'
								 ,'#{hp}'
								 ,'#{price}'
								 ,'#{priceDate}'
								 ,'#{sDate}'
						)
`
