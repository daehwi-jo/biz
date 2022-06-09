package mngAdmin

var SelectOrderListCnt = `	
			    	SELECT 	COUNT(*) AS TOTAL_COUNT
			    	FROM dar_order_info AS A
			    	INNER JOIN (SELECT AA.ORDER_NO ,ITEM_NO, COUNT(*) AS ITEM_CNT
			    					FROM dar_order_detail AS AA
			    					WHERE 
									DATE_FORMAT(AA.ORDER_DATE, '%Y-%m-%d') >=  '#{startDate}'
									AND DATE_FORMAT(AA.ORDER_DATE, '%Y-%m-%d') <=  '#{endDate}'
									GROUP BY AA.ORDER_NO ) AS E ON A.ORDER_NO  = E.ORDER_NO
			    	INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID AND A.ORDER_TY !='4'
			    	INNER JOIN priv_rest_info AS C ON A.REST_ID = C.REST_ID
			    	INNER JOIN priv_grp_info AS D ON A.GRP_ID = D.GRP_ID
			    	LEFT OUTER JOIN dar_sale_item_info AS F ON A.REST_ID = F.REST_ID AND E.ITEM_NO = F.ITEM_NO
					WHERE 1=1
					AND GRP_NM LIKE '%#{searchGrpNm}%'
					AND USER_NM LIKE '%#{searchUserNm}%'
					AND REST_NM LIKE '%#{searchRestNm}%'
					AND ORDER_STAT = '#{searchStat}'
			`

var SelectOrderList = `	SELECT *
		    	FROM 
		    		(
			    	SELECT 	DATE_FORMAT(A.ORDER_DATE, '%Y-%m-%d %H:%i:%s') AS ORDER_DATE	/* 주문일자 */
						, A.ORDER_NO			/* 주문등록번호 */
						, A.USER_ID		/* 주문자ID */
						, B.USER_NM	  /* 주문자명 */
						, A.REST_ID
						, C.REST_NM
						, D.GRP_NM
						, (SELECT SUM(ORDER_QTY) FROM DAR_ORDER_DETAIL WHERE ORDER_NO = A.ORDER_NO) AS ORDER_CNT	/* 주문건수 */
						, A.CREDIT_AMT			/* 외상금액 */
						, A.ORDER_STAT /* 주문상태 */
						, A.ORDER_COMMENT
						, A.PAID_YN
						, DATE_FORMAT(A.PAY_DATE, '%Y-%m-%d %H:%i:%s') AS PAY_DATE	/* 결제일자 */
						, PAY_TY
	               , ITEM_CNT
	               , IFNULL(CASE WHEN E.ITEM_CNT > 1 THEN CONCAT(F.ITEM_NM, ' 외') 
	               		  ELSE F.ITEM_NM 
	               		  END ,'금액권')AS ITEM_NM

			    	FROM dar_order_info AS A
			    	INNER JOIN (SELECT AA.ORDER_NO ,ITEM_NO, COUNT(*) AS ITEM_CNT
			    					FROM dar_order_detail AS AA
			    					WHERE 
									DATE_FORMAT(AA.ORDER_DATE, '%Y-%m-%d') >=  '#{startDate}'
									AND DATE_FORMAT(AA.ORDER_DATE, '%Y-%m-%d') <=  '#{endDate}'
									GROUP BY AA.ORDER_NO ) AS E ON A.ORDER_NO  = E.ORDER_NO
			    	INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID AND A.ORDER_TY !='4'
			    	INNER JOIN priv_rest_info AS C ON A.REST_ID = C.REST_ID
			    	INNER JOIN priv_grp_info AS D ON A.GRP_ID = D.GRP_ID
			    	LEFT OUTER JOIN dar_sale_item_info AS F ON A.REST_ID = F.REST_ID AND E.ITEM_NO = F.ITEM_NO
			    	ORDER BY ORDER_DATE DESC
			    	) AS T
					WHERE 1=1
					AND GRP_NM LIKE '%#{searchGrpNm}%'
					AND USER_NM LIKE '%#{searchUserNm}%'
					AND REST_NM LIKE '%#{searchRestNm}%'
                    AND ORDER_STAT = '#{searchStat}'
					ORDER BY ORDER_DATE DESC
		`

var SelectOrderListExcel = `
					SELECT 
                          DATE_FORMAT(A.ORDER_DATE, '%Y-%m-%d %H:%i:%s') AS a01	
						, D.GRP_NM as a02
						, B.USER_NM	 as a03 
						, C.REST_NM as a04
						, IFNULL(CASE WHEN E.ITEM_CNT > 1 THEN CONCAT(F.ITEM_NM, ' 외') 
	               		  ELSE F.ITEM_NM 
	               		  END ,'금액권') AS a05
						, (SELECT SUM(ORDER_QTY) FROM DAR_ORDER_DETAIL WHERE ORDER_NO = A.ORDER_NO) AS a06
						, A.CREDIT_AMT	as a07
					    , CASE WHEN PAY_TY='0' THEN '선불' ELSE '후불' END AS a08
					    , CASE WHEN PAID_YN='Y' THEN  DATE_FORMAT(A.PAY_DATE, '%Y-%m-%d %H:%i:%s') ELSE PAID_YN END AS a09
						, CASE WHEN A.ORDER_STAT='21' THEN '주문취소' ELSE '주문완료' END AS a10
			    	FROM dar_order_info AS A
			    	INNER JOIN (SELECT AA.ORDER_NO ,ITEM_NO, COUNT(*) AS ITEM_CNT
			    					FROM dar_order_detail AS AA
			    					WHERE 
									DATE_FORMAT(AA.ORDER_DATE, '%Y-%m-%d') >=  '#{startDate}'
									AND DATE_FORMAT(AA.ORDER_DATE, '%Y-%m-%d') <=  '#{endDate}'
									GROUP BY AA.ORDER_NO ) AS E ON A.ORDER_NO  = E.ORDER_NO
			    	INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID AND A.ORDER_TY !='4'
			    	INNER JOIN priv_rest_info AS C ON A.REST_ID = C.REST_ID
			    	INNER JOIN priv_grp_info AS D ON A.GRP_ID = D.GRP_ID
			    	LEFT OUTER JOIN dar_sale_item_info AS F ON A.REST_ID = F.REST_ID AND E.ITEM_NO = F.ITEM_NO
					WHERE 1=1
					AND GRP_NM LIKE '%#{searchGrpNm}%'
					AND USER_NM LIKE '%#{searchUserNm}%'
					AND REST_NM LIKE '%#{searchRestNm}%'
					AND ORDER_STAT = '#{searchStat}'
					ORDER BY ORDER_DATE DESC
		`

var SelectOrderInfo string = `SELECT A.ORDER_NO
											,B.REST_NM
											,C.GRP_NM
											,A.TOTAL_AMT
											,A.ORDER_STAT
											,DATE_FORMAT(A.ORDER_DATE,'%Y.%m.%d %p%h:%i') AS ORDER_DATE
											,ifnull(A.ORDER_COMMENT,'') AS ORDER_COMMENT
											,A.QR_ORDER_TYPE
											,DATE_FORMAT(A.ORDER_CANCEL_DATE,'%Y.%m.%d %p%h:%i') AS ORDER_CANCEL_DATE
											,B.REST_TYPE
									FROM dar_order_info AS A
									INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
									INNER JOIN priv_grp_info AS C ON A.GRP_ID = C.GRP_ID
									WHERE 
									A.ORDER_NO = '#{orderNo}'
									`

var SelectOrderDetail string = `SELECT CASE WHEN  A.ITEM_NO='9999999999'  THEN '금액권'  ELSE B.ITEM_NM END  as menuNm
										 ,SUM(ORDER_AMT* ORDER_QTY) as menuPrice
										 ,SUM(ORDER_QTY) as menuQty
										 ,IFNULL(CPNO,'') AS CPNO
										 ,IFNULL(ORD_NO,'') AS ORD_NO
										 ,IFNULL(EXCH_FR_DY,'') AS EXCH_FR_DY
										 ,IFNULL(EXCH_TO_DY,'') AS EXCH_TO_DY
										 ,IFNULL(EXPIRE_DATE,'') AS EXPIRE_DATE
										 ,IFNULL(CP_STATUS,'') AS CP_STATUS
										 ,IFNULL(EXCHPLC,'') AS EXCHPLC
										 ,IFNULL(EXCHCO_NM,'') AS EXCHCO_NM
										 ,IFNULL(CPNO_EXCH_DT,'') AS CPNO_EXCH_DT
										 ,IFNULL(CPNO_STATUS,'') AS CPNO_STATUS
										 ,IFNULL(BALANCE,'') AS BALANCE
								FROM dar_order_detail AS A 
								LEFT OUTER JOIN dar_sale_item_info AS B ON A.ITEM_NO = B.ITEM_NO
								LEFT OUTER JOIN dar_order_coupon AS c ON a.order_no = c.order_no
								WHERE 
								A.ORDER_NO= '#{orderNo}'
                                GROUP BY A.ITEM_NO 
								`

var SelectOrderUserSplitAmt string = `SELECT C.USER_NM
											,SPLIT_AMT AS ORDER_AMT
											,C.USER_ID
											,IFNULL(D.MEMO,'') AS MEMO
										FROM dar_order_split AS A
										INNER JOIN priv_user_info AS C ON A.USER_ID = C.USER_ID
										LEFT OUTER JOIN dar_order_memo AS D ON A.ORDER_NO=D.ORDER_NO AND C.USER_ID = D.USER_ID
										WHERE 
											A.ORDER_NO= '#{orderNo}'
								      `

var SelectOrderUserDetail string = `SELECT C.USER_NM
										,SUM(A.ORDER_QTY * A.ORDER_AMT) AS ORDER_AMT
										,C.USER_ID
										,IFNULL(D.MEMO,'') AS MEMO
								FROM dar_order_detail AS A
								INNER JOIN priv_user_info AS C ON A.USER_ID = C.USER_ID
								LEFT OUTER JOIN dar_order_memo AS D ON A.ORDER_NO=D.ORDER_NO AND C.USER_ID = D.USER_ID
								WHERE 
								A.ORDER_NO= '#{orderNo}'
								GROUP BY  C.USER_NM,C.USER_ID
								
								`

var SelectOrderUserMenu string = `SELECT A.ORDER_QTY 
										,A.ORDER_AMT
										,CASE A.ITEM_NO WHEN '9999999999' THEN '금액권' ELSE B.ITEM_NM END AS ITEM_NM
								FROM dar_order_detail AS A
								LEFT OUTER JOIN dar_sale_item_info AS B ON A.ITEM_NO = B.ITEM_NO
								WHERE 
								A.ORDER_NO= '#{orderNo}'
								AND A.USER_ID = '#{userId}'
								`

var SelectOrder string = `SELECT A.ORDER_NO
							,A.TOTAL_AMT
							,A.ORDER_STAT
							,A.PAY_TY
							,A.GRP_ID AS BOOK_ID
							,A.REST_ID AS STORE_ID
							,A.USER_ID
							,A.POINT_USE
					FROM dar_order_info AS A
					WHERE 
					A.ORDER_NO = '#{orderNo}'
					`

var SelectLinkInfo string = `SELECT 
							A.AGRM_ID AS LINK_ID
							, A.GRP_ID
							, FN_GET_GRPNAME(A.GRP_ID) AS BOOK_NM
							, A.REST_ID
							, FN_GET_RESTNAME(A.REST_ID) AS STORE_NM
							, A.REQ_STAT
							, FN_GET_CODENAME('AGRM_STAT', A.REQ_STAT) AS REQ_STAT_NM
							, A.REQ_TY
							, DATE_FORMAT(A.REQ_DATE, '%Y-%m-%d') AS REQ_DATE	
							, DATE_FORMAT(A.AUTH_DATE, '%Y-%m-%d') AS AUTH_DATE
							, A.PAY_TY
							, IFNULL(A.PREPAID_AMT, 0) AS PREPAID_AMT
							, IFNULL(A.PREPAID_POINT, 0) AS PREPAID_POINT
							FROM org_agrm_info  AS a
							INNER join PRIV_REST_INFO AS B on A.REST_ID = B.REST_ID or (B.FRAN_YN = 'Y' and B.FRAN_ID = A.REST_ID)
							WHERE
							B.REST_ID ='#{storeId}'
							AND A.GRP_ID='#{bookId}'
							`

var UpdateLink string = `UPDATE org_agrm_info SET 
							PREPAID_AMT = '#{prepaidAmt}'
							,PREPAID_POINT = #{prepaidPoint}
							,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							WHERE 
							AGRM_ID ='#{linkId}'
						`

var UpdateBookUserSupportBalance string = `UPDATE priv_grp_user_info SET  
								SUPPORT_BALANCE = SUPPORT_BALANCE + #{orderAmt},
								MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE 
								GRP_ID = '#{bookId}'
								AND USER_ID ='#{userId}'`

var UpdateOrderCancel string = ` UPDATE dar_order_info SET ORDER_STAT='21'
												, ORDER_CANCEL_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE
								ORDER_NO = '#{orderNo}'
									`
