package mngOrder

var SelectPaymentList = `SELECT  DATE_FORMAT(A.REG_DATE ,'%Y-%m-%d') AS REG_DATE 
				,DATE_FORMAT(A.REG_DATE,'%H:%i:%S') AS REG_TIME
				,GRP_NM
				,REST_NM
				,A.SEARCH_TY
				,A.PAYMENT_TY
				,A.PAY_INFO
				,A.CREDIT_AMT
				,IFNULL((SELECT DATE_FORMAT(CANCELDATE ,'%Y-%m-%d')  FROM dar_payment_report AS  AA WHERE A.MOID= AA.MOID LIMIT 1),'') AS CANCEL_DATE
				,IFNULL((SELECT AMT  FROM dar_payment_report AS  AA WHERE A.MOID= AA.MOID LIMIT 1),'') AS P_AMT
				FROM dar_payment_hist AS A
				INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
				INNER JOIN PRIV_GRP_INFO AS C ON A.GRP_ID = C.GRP_ID
				INNER JOIN b_company_book AS D ON C.grp_id = D.book_id
				WHERE 
					D.company_id='#{companyId}'
				AND C.GRP_ID = '#{searchGrpId}'
				AND B.REST_NM LIKE CONCAT('%', '#{searchRestNm}', '%')
				AND DATE_FORMAT(A.REG_DATE,'%Y-%m-%d') BETWEEN DATE_FORMAT('#{searchStartDt}','%Y-%m-%d') AND DATE_FORMAT('#{searchEndDt}','%Y-%m-%d')
				ORDER BY A.REG_DATE DESC`

var SelectPaymentListCnt = `SELECT COUNT(*) as TOTAL_COUNT, sum(A.CREDIT_AMT) as ALL_AMT
        					FROM dar_payment_hist AS A
							INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
							INNER JOIN PRIV_GRP_INFO AS C ON A.GRP_ID = C.GRP_ID
							INNER JOIN b_company_book AS D ON C.grp_id = D.book_id
							WHERE 
								D.company_id = '#{companyId}'
								AND C.GRP_ID = '#{searchGrpId}'
								AND B.REST_NM LIKE '%#{searchRestNm}%'
								AND DATE_FORMAT(A.REG_DATE,'%Y-%m-%d') BETWEEN DATE_FORMAT('#{searchStartDt}','%Y-%m-%d') 
								AND DATE_FORMAT('#{searchEndDt}','%Y-%m-%d')`
