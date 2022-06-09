package mngOrder

var SelectStoreMng = `SELECT *
					FROM ( SELECT A.REST_ID
							,A.REST_NM
							,C.GRP_NM
							,IFNULL(AA.ORDER_AMT,0) AS ORDER_AMT
							,IFNULL(AA.NOPAY_AMT,0) AS NOPAY_AMT
							,IFNULL(AA.ORDER_CNT,0) AS ORDER_CNT
							,DATE_FORMAT(IFNULL(AA.ORDER_DATE,''),'%Y-%m-%d') AS ORDER_DATE 
							,B.REQ_STAT
							,B.PAY_TY
							,C.GRP_ID
							,B.PREPAID_AMT
					FROM PRIV_REST_INFO AS A
					INNER JOIN ORG_AGRM_INFO AS B ON A.REST_ID = B.REST_ID AND B.REQ_STAT='1'
					INNER JOIN PRIV_GRP_INFO AS C ON B.GRP_ID = C.GRP_ID
					INNER JOIN B_COMPANY_BOOK AS D ON C.GRP_ID = D.BOOK_ID	
					LEFT OUTER JOIN (	SELECT   A.GRP_ID
											 , A.REST_ID
											 , IFNULL(SUM(A.TOTAL_AMT),0) AS ORDER_AMT
											 , IFNULL(SUM(CASE WHEN (A.PAID_YN='N' AND A.PAY_TY = '1') THEN A.TOTAL_AMT ELSE 0 END ),0) AS NOPAY_AMT 
											 , MAX(A.ORDER_DATE ) AS ORDER_DATE 
											 , MAX(A.PAY_DATE ) AS AUTH_DATE 
											 , IFNULL(COUNT(*),0) AS ORDER_CNT
										FROM DAR_ORDER_INFO A
										WHERE 1=1
										AND A.order_ty IN ('1','2','3','5')
										AND A.ORDER_STAT = '20'	
										AND DATE_FORMAT(A.ORDER_DATE,'%Y%m%d') >= '#{startDate}'
										AND DATE_FORMAT(A.ORDER_DATE,'%Y%m%d') <= '#{endDate}'
										GROUP BY A.GRP_ID, A.REST_ID
							) AS AA ON C.GRP_ID = AA.GRP_ID AND A.REST_ID = AA.REST_ID
					WHERE 
						D.COMPANY_ID='#{companyId}'
						AND c.grp_id = '#{searchGrpId}'
						AND A.REST_NM LIKE '%#{searchRestNm}%'						
						AND A.USE_YN = 'Y'
					  ) AS ZZ
					ORDER BY ORDER_DATE DESC

`

var SelectStoreMngCnt = ` SELECT COUNT(*) as TOTAL_COUNT
								 ,sum(AA.ORDER_CNT) as TOTAL_ORDER_COUNT
							FROM PRIV_REST_INFO AS A
							INNER JOIN ORG_AGRM_INFO AS B ON A.REST_ID = B.REST_ID AND B.REQ_STAT='1'
							INNER JOIN PRIV_GRP_INFO AS C ON B.GRP_ID = C.GRP_ID
							INNER JOIN B_COMPANY_BOOK AS D ON C.GRP_ID = D.BOOK_ID	
							LEFT OUTER JOIN (	SELECT   A.GRP_ID
											 , A.REST_ID
											 , IFNULL(SUM(A.TOTAL_AMT),0) AS ORDER_AMT
											 , IFNULL(SUM(CASE WHEN (A.PAID_YN='N' AND A.PAY_TY = '1') THEN A.TOTAL_AMT ELSE 0 END ),0) AS NOPAY_AMT 
											 , MAX(A.ORDER_DATE ) AS ORDER_DATE 
											 , MAX(A.PAY_DATE ) AS AUTH_DATE 
											 , IFNULL(COUNT(*),0) AS ORDER_CNT
										FROM DAR_ORDER_INFO A
										WHERE 1=1
										AND A.order_ty IN ('1','2','3','5')
										AND A.ORDER_STAT = '20'	
										AND DATE_FORMAT(A.ORDER_DATE,'%Y%m%d') >= '#{startDate}'
										AND DATE_FORMAT(A.ORDER_DATE,'%Y%m%d') <= '#{endDate}'
										GROUP BY A.GRP_ID, A.REST_ID
							) AS AA ON C.GRP_ID = AA.GRP_ID AND A.REST_ID = AA.REST_ID
							WHERE 
								D.COMPANY_ID = '#{companyId}'
								AND c.grp_id = '#{searchGrpId}'
          						AND A.REST_NM LIKE '%#{searchRestNm}%'
								AND A.USE_YN = 'Y'
`

var SelectStoreInfo = `SELECT A.REST_ID AS restId
										,A.REST_NM as restNm
										,A.ADDR as addr
										,IFNULL(A.ADDR2,'') AS addr2
										,IFNULL(A.TEL,'') AS tel
										,IFNULL(CONCAT(B.FILE_PATH,'/',B.SYS_FILE_NM)
								 			,CONCAT('/public/img/',CASE A.BUETY WHEN '00' THEN '한식' 
											  		  WHEN '01' THEN '중식' 
											  		  WHEN '02' THEN '일식' 
											  		  WHEN '03' THEN '양식' 
											  		  WHEN '04' THEN '카페' 
											  		  WHEN '05' THEN '분식' 
											  		  WHEN '06' THEN '부페' 
											  		  WHEN 'CA' THEN '부페' 
											  		  WHEN '07' THEN '기타' 
											  		  WHEN '08' THEN '유통' 
											  		  WHEN '09' THEN '뷰티' 
													  WHEN NULL THEN '기타' 
													  WHEN '' THEN '기타' 
													  ELSE cc.CODE_NM
											 END,'.png')) AS restImg
								FROM priv_rest_info AS A
								LEFT OUTER JOIN priv_rest_file AS B ON A.REST_ID = B.REST_ID AND B.FILE_TY='1'
								LEFT OUTER JOIN b_category AS bb ON a.category = bb.CATEGORY_ID AND A.USE_YN='Y'
							    LEFT OUTER JOIN b_code AS cc ON A.BUETY = cc.CODE_ID  AND A.USE_YN='Y'
								WHERE 
									A.REST_ID = '#{restId}'
					`

var SelectChargeAmtList string = `SELECT AMT
											,ADD_AMT
										,IFNULL(AMT+ADD_AMT,0) AS SUM_AMT
									FROM dar_prepayment_charge_info
									WHERE 
									REST_ID='#{restId}'
									AND USE_YN='Y'
									ORDER BY AMT asc
								`

var SelectUnpaidList string = `SELECT
									DATE_FORMAT(A.ORDER_DATE,'%Y-%m-%d')  AS ORDER_DATE
									,SUM(A.TOTAL_AMT) AS TOTAL_AMT
									,COUNT(*) AS ORDER_CNT
								FROM  dar_order_info AS A
								WHERE 
								A.REST_ID='#{restId}'
								AND A.order_ty IN ('1','2','3','5')
								AND A.GRP_ID ='#{grpId}'
								AND PAY_TY='1' 
								AND PAID_YN ='N'
								AND order_stat = '20'
								AND LEFT(A.ORDER_DATE,8) <='#{accStDay}'
								GROUP BY DATE_FORMAT(A.ORDER_DATE,'%Y.%m.%d')
								ORDER BY DATE_FORMAT(A.ORDER_DATE,'%Y.%m.%d')
								`

var SelectUnpaidListCount string = `SELECT COUNT(*) AS orderCnt
									, SUM(total_amt) AS TOTAL_AMT
								FROM  dar_order_info AS A
								INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID
								WHERE 
								A.REST_ID='#{restId}'
								AND A.order_ty IN ('1','2','3','5')
								AND A.GRP_ID ='#{grpId}'
								AND PAY_TY='1' 
								AND PAID_YN='N'
								AND order_stat = '20'
								AND LEFT(A.ORDER_DATE,8) <='#{accStDay}'
								`
