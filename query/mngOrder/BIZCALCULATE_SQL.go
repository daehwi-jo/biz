package mngOrder

var SelectGrpCalculateList = `SELECT E.REST_ID, E.REST_NM
					,SUM(D.TOTAL_AMT) AS TOTAL_AMT
					,B.GRP_NM
					,B.GRP_ID
				FROM  b_company_book AS A
				inner join priv_grp_info as B on A.book_id = B.grp_id
				inner join org_agrm_info as C on B.GRP_ID = C.GRP_ID AND C.REQ_STAT='1'
				inner join DAR_ORDER_INFO as D on C.GRP_ID = D.GRP_ID AND C.REST_ID = D.REST_ID AND D.ORDER_STAT='20' AND D.PAY_TY='1'
				inner join priv_rest_info as E on D.REST_ID = E.REST_ID	
				WHERE D.PAID_YN='N' 
					AND C.PAY_TY=1 
					AND A.company_id='#{companyId}'
					AND C.GRP_ID = '#{searchGrpId}'
					AND DATE_FORMAT(D.ORDER_DATE,'%Y-%m-%d') <= '#{searchEndDt}'
				GROUP BY E.REST_ID, E.REST_NM,B.GRP_NM,B.GRP_ID`

var SelectGrpCalculateCnt = `SELECT COUNT(*) as TOTAL_COUNT
									FROM( 
										SELECT E.REST_ID, E.REST_NM
												,SUM(D.TOTAL_AMT) AS TOTAL_AMT
												FROM  b_company_book AS A
												inner join priv_grp_info as B on A.book_id = B.grp_id
												inner join org_agrm_info as C on B.GRP_ID = C.GRP_ID AND C.REQ_STAT='1'
												inner join DAR_ORDER_INFO as D on C.GRP_ID = D.GRP_ID AND C.REST_ID = D.REST_ID AND D.ORDER_STAT='20' AND D.PAY_TY='1'
												inner join priv_rest_info as E on D.REST_ID = E.REST_ID	
												WHERE D.PAID_YN='N' 
													AND C.PAY_TY=1 
													AND A.company_id= '#{companyId}'
				 									AND C.GRP_ID = '#{searchGrpId}'
													AND DATE_FORMAT(D.ORDER_DATE,'%Y-%m-%d') <= '#{searchEndDt}'
												GROUP BY E.REST_ID, E.REST_NM) A`



var SelectTpayUnpaidList = `SELECT E.REST_ID, E.REST_NM
					,SUM(D.TOTAL_AMT) AS TOTAL_AMT
				FROM  b_company_book AS A
				inner join priv_grp_info as B on A.book_id = B.grp_id
				inner join org_agrm_info as C on B.GRP_ID = C.GRP_ID AND C.REQ_STAT='1'
				inner join DAR_ORDER_INFO as D on C.GRP_ID = D.GRP_ID AND C.REST_ID = D.REST_ID AND D.ORDER_STAT='20' AND D.PAY_TY='1'
				inner join priv_rest_info as E on D.REST_ID = E.REST_ID	
				WHERE D.PAID_YN='N' 
					AND C.PAY_TY=1 
					AND A.company_id='#{companyId}'
					AND C.GRP_ID = '#{grpId}'
					AND D.REST_ID IN(#{restIdArray}) 
					AND DATE_FORMAT(D.ORDER_DATE,'%Y%m%d') <= '#{selectedDate}'
				GROUP BY E.REST_ID, E.REST_NM,B.GRP_NM,B.GRP_ID`

var UpdateUnpaidReadyData =`UPDATE  dar_order_info SET PAID_KEY='#{moid}'
							WHERE PAID_YN='N' 
							AND PAY_TY='1' 
							AND ORDER_STAT='20'
							AND GRP_ID = '#{grpId}'
							AND REST_ID IN(#{restIdArray}) 
							AND DATE_FORMAT(ORDER_DATE,'%Y%m%d') <= '#{selectedDate}'

`