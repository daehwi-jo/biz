package mngAdmin

var SelectStoreList = `SELECT A.REST_ID
					,A.REST_NM
					,A.BUSID 
					,IFNULL(A.CEO_NM,'') as USER_NM
					,A.TEL
					,A.USE_YN
					,DATE_FORMAT(IFNULL(A.REG_DATE,''),'%Y-%m-%d') AS REG_DATE
					,(SELECT COUNT(*) FROM org_agrm_info AS AA WHERE A.REST_ID = AA.REST_ID AND AA.REQ_STAT='1') AS LINK_CNT
					,A.ADDR
				FROM priv_rest_info AS A
				WHERE AUTH_STAT='1'
`

var SelectStoreInfo = `SELECT A.REST_ID
						,A.REST_NM
						,A.BUSID
						,A.CEO_NM
						,A.TEL
						,EMAIL
						,INTRO
						,OPEN_WEEK
						,OPEN_WEEKEND
						,CODE_NM AS BANK_NM
						,ACCOUNT_NO
						,ACCOUNT_NM
						,DATE_FORMAT(IFNULL(A.REG_DATE,''),'%Y-%m-%d') AS JOIN_DATE
						,IFNULL(CONCAT(C.FILE_PATH,'/',C.SYS_FILE_NM),'') AS REST_IMG
						,A.USE_YN
						,IFNULL(d.user_id ,'')as ceoId
						,A.REST_TYPE as storeType
						FROM priv_rest_info AS A
						left outer JOIN b_code AS B ON A.BANK_CD = B.CODE_ID AND B.CATEGORY_ID='BANK'
						LEFT OUTER JOIN priv_rest_file AS C ON A.REST_ID = C.REST_ID AND C.FILE_TY='1'
						left outer join priv_rest_user_info as d on a.rest_id = d.rest_id and rest_auth = 0
						WHERE 
						A.REST_ID = '#{restId}'`

var SelectStoreEtc = `
SELECT 	MEMO as memo
		,NOTICE as notice 
FROM priv_rest_etc 
WHERE 
	REST_ID = '#{restId}'
`

var SelectStoreLinkBookList = `SELECT B.GRP_ID
								,B.GRP_NM
								,IFNULL(A.PREPAID_AMT,0) AS PREPAID_AMT
								,DATE_FORMAT(IFNULL(A.AUTH_DATE,''),'%Y-%m-%d') AS AUTH_DATE
								,A.REQ_STAT
								,(SELECT COUNT(*) 
									FROM priv_grp_user_info AA 
									WHERE 
										A.GRP_ID = AA.GRP_ID 
										and AUTH_STAT = '1') AS USER_CNT
							    ,(SELECT IFNULL(SUM(AA.TOTAL_AMT),0) 
									FROM dar_order_info AS AA
									WHERE A.REST_ID = AA.REST_ID AND AA.GRP_ID = B.GRP_ID
									AND AA.order_ty IN ('1','2','3','5')
									AND AA.PAY_TY='1' 
									AND AA.PAID_YN ='N'
									AND AA.order_stat = '20') AS ORDER_AMT
							FROM org_agrm_info AS A
							INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID 
							WHERE REQ_STAT = '1' 
							AND A.REST_ID = '#{restId}'
							AND B.GRP_NM LIKE '%#{searchKeyword}%'
							ORDER BY A.REQ_DATE DESC`

var SelectStoreLinkBookListCnt = `
SELECT Count(*) as total
FROM org_agrm_info AS A
INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID
WHERE REQ_STAT = '1'
AND A.REST_ID = '#{restId}'
AND B.GRP_NM LIKE '%#{searchKeyword}%'
ORDER BY A.REQ_DATE DESC
`

var SelectStoreLinkBookInfo = `
select 
		B.GRP_ID as grpId
		,B.GRP_NM
		,(SELECT COUNT(*) 
			FROM priv_grp_user_info AA 
			WHERE B.GRP_ID = AA.GRP_ID 
			and AUTH_STAT = 1) AS USER_CNT
		, IFNULL(A.REQ_STAT, 4) AS REQ_STAT 
		, IFNULL(A.PREPAID_AMT, 0) AS PREPAID_AMT
		, IFNULL(A.PAY_TY,B.GRP_PAY_TY) as payTy
		, DATE_FORMAT(A.AUTH_DATE, '%Y-%m-%d') AS AUTH_DATE
		, DATE_FORMAT(A.MOD_DATE, '%Y-%m-%d') AS MOD_DATE
		, IFNULL(A.REQ_TY, '3') as reqTy
		, IFNULL(A.LINK_TY,'G') as linkTy
		, A.REJ_COMMENT as rejCom
FROM priv_grp_info as B 
LEFT OUTER JOIN org_agrm_info AS A ON A.GRP_ID = B.GRP_ID AND A.REST_ID = '#{restId}'
WHERE 
B.GRP_ID = '#{grpId}'
`

var SelectStoreBookSearch = `
select 
		B.GRP_ID as grpId
		,B.GRP_NM
		,(SELECT COUNT(*) 
			FROM priv_grp_user_info AA 
			WHERE B.GRP_ID = AA.GRP_ID 
			and AUTH_STAT = 1) AS USER_CNT
		, IFNULL(A.REQ_STAT, 4) AS REQ_STAT 
		, IFNULL(A.PREPAID_AMT, 0) AS PREPAID_AMT
		, B.GRP_PAY_TY
		, DATE_FORMAT(A.AUTH_DATE, '%Y-%m-%d') AS AUTH_DATE
		, D.user_nm as grpMaster
FROM priv_grp_info as B 
LEFT OUTER JOIN org_agrm_info AS A ON A.GRP_ID = B.GRP_ID AND A.REST_ID = '#{restId}'
INNER JOIN priv_grp_user_info AS C ON C.GRP_ID = B.GRP_ID AND C.GRP_AUTH = 0
INNER JOIN priv_user_info AS D ON C.USER_ID = D.USER_ID 
WHERE 
B.USE_YN = 'Y'
AND B.grp_nm LIKE '%#{searchKeyword}%'
`

var SelectStoreBookSearchCnt = `
select 
		Count(*) as total
FROM priv_grp_info as B 
LEFT OUTER JOIN org_agrm_info AS A ON A.GRP_ID = B.GRP_ID AND A.REST_ID = '#{restId}'
WHERE 
B.USE_YN = 'Y'
AND B.grp_nm LIKE '%#{searchKeyword}%'
`

var SelectStoreItemList = `SELECT ITEM_NO
									,A.ITEM_NM
									,A.ITEM_PRICE
									,B.CODE_NM
									,A.ITEM_MENU as CODE_ID
									,A.USE_YN
							FROM dar_sale_item_info AS A
							INNER JOIN dar_category_info AS B ON A.REST_ID = B.REST_ID AND A.ITEM_MENU = B.CODE_ID
							WHERE ITEM_STAT='1'
							AND A.USE_YN = 'Y'
							AND A.REST_ID = '#{restId}'`

var SelectStoreCategoryItemList = `SELECT ITEM_NO
									,A.ITEM_NM
									,A.ITEM_PRICE
									,B.CODE_NM
									,A.ITEM_MENU as CODE_ID
									,A.USE_YN
									,IFNULL(A.prod_id,'') as prodId
									,A.ITEM_IMG
									,A.PROD_ID as prodId
									,A.BEST_YN as bestYn
									,DATE_FORMAT(A.SALES_END_DATE,'%Y-%m-%d') as salesEndDate
							FROM dar_sale_item_info AS A
							INNER JOIN dar_category_info AS B ON A.REST_ID = B.REST_ID AND A.ITEM_MENU = B.CODE_ID
							WHERE A.ITEM_STAT='1'
							AND A.REST_ID = '#{restId}'
`

var SelectStoreCnt = `SELECT COUNT(*) as TOTAL_COUNT
       					FROM priv_rest_info AS A
						WHERE AUTH_STAT='1'
`

var SelectRestImgInfo = `SELECT FILE_NM
								,SYS_FILE_NM
						FROM priv_rest_file
						WHERE 
						REST_ID='#{restId}'
						AND FILE_TY='1'
						`

var SelectRestFileSeq = `SELECT
							IFNULL(LPAD(MAX(FILE_NO) + 1, 10, 0), '0000000001') as FILE_NO
							FROM
							PRIV_REST_FILE
						`

var InsertRestImg = `INSERT INTO PRIV_REST_FILE
								(
									FILE_NO					
									, REST_ID				
									, FILE_TY				                                                         
									, FILE_PATH			                                                                                                                  
									, FILE_NM			
									, SYS_FILE_NM	
									, FILE_SIZE			
									, REG_DATE			
									, MOD_DATE		                                                     
								) 
								VALUES 
								(
									'#{fileNo}'
									, '#{restId}'
									, '#{fileTy}'
									, '#{filePath}'
									, '#{fileName}'
									, '#{sysFileName}'
									, '#{fileSize}'
									, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								)
						`

var UpdateRestImg = `UPDATE PRIV_REST_FILE  SET
									 MOD_DATE=DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									, FILE_PATH   =		'#{filePath}'		                                                                                                                     
									, FILE_NM	  =		'#{fileName}'
									, SYS_FILE_NM =		'#{sysFileName}'	
									, FILE_SIZE   =		'#{fileSize}'
					WHERE
						REST_ID='#{restId}' 
						AND FILE_TY = '#{fileTy}'
						`

//계좌 정보 추가안됨
var UpdateStoreInfo = `UPDATE priv_rest_info 
		SET 
		 MOD_DATE=DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
		 ,rest_nm = '#{restNm}'
		 ,busid = '#{bizNum}'
		 ,tel = '#{tel}'
		 ,email = '#{email}'
		 ,ceo_nm = '#{ceo}'
		 ,intro = '#{intro}'
		 ,open_week = '#{open}'
		 ,BANK_CD = '#{bankCd}'
		 ,account_nm = '#{accountNm}'
		 ,account_no = '#{accountNo}'
		 ,use_yn = '#{useYn}'
		 ,ACCOUNT_CERT_YN = '#{accountCertYn}'
		WHERE
			rest_id = '#{restId}'`

var UpdateStoreEtc = `
UPDATE priv_rest_etc
SET
	MEMO = '#{memo}'
	,NOTICE = '#{notice}'
	,MOD_DATE = DATE_FORMAT(now(),'%Y%m%d%H%i%s')
WHERE
	REST_ID = '#{restId}'
`

var UpdateBookLink = `UPDATE org_agrm_info 
						SET 
							req_stat = '#{authCode}'
							, mod_date = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s') 
						WHERE 
							rest_id = '#{restId}' 
						and grp_id ='#{grpId}'`

var InsertBookLink = `
INSERT INTO org_agrm_info (
AGRM_ID
,GRP_ID
,REST_ID
,REQ_STAT
,REQ_TY
,AUTH_DATE 
,MOD_DATE
,PAY_TY
,LINK_TY
,REJ_COMMENT
)
VALUES (
'#{agrmId}'
,'#{grpId}'
,'#{restId}'
,'#{authCode}'
,'#{reqTy}'
,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
,'#{payTy}'
,'#{linkTy}'
,'#{rejCom}'
)
ON DUPLICATE KEY UPDATE 
REQ_STAT = '#{authCode}'
,REQ_TY = '#{reqTy}'
,AUTH_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s') 
,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
,PAY_TY = '#{payTy}'
,LINK_TY = '#{linkTy}'
,REJ_COMMENT = '#{rejCom}'
`

var SelectCategoryList = `
SELECT * 
FROM dar_category_info 
WHERE 
     rest_id = '#{restId}'
`

var SelectNewCategory = `
                SELECT 
                IFNULL(CATEGORY_ID,FN_GET_SEQUENCE('CATEGORY_ID')) AS categoryId
                ,(SELECT rest_nm from priv_rest_info 
					WHERE 
						REST_ID='#{restId}') AS categoryNm
                ,CONCAT(IFNULL(LPAD(MAX(SUBSTRING(CODE_ID, -5)) + 1, 5, 0), '00001')) as codeId
                FROM DAR_CATEGORY_INFO as A
                WHERE 
                A.REST_ID='#{restId}'
                AND A.CODE_ID <> '99999'
`

var InsertCategory = `
INSERT INTO dar_category_info (
category_id, category_nm, rest_id, code_id, code_nm, use_yn)

VALUES(
	'#{categoryId}'
	,'#{categoryNm}'
	,'#{restId}'
	,'#{codeId}'
	,'#{codeNm}'
	,'#{useYn}'
	
)
`

var UpdateCategory = `
UPDATE dar_category_info 
SET 
    code_nm = '#{codeNm}'
    ,use_yn = '#{useYn}'
WHERE 
    code_id = '#{codeId}' 
AND rest_id = '#{restId}'
`

var InsertMenu = `
INSERT INTO dar_sale_item_info
(
ITEM_NO
,ITEM_NM
,REST_ID
,ITEM_STAT
,ITEM_PRICE
,REG_DATE
,USE_YN
,ITEM_MENU
,PROD_ID
,SALES_END_DATE
,ITEM_IMG
) VALUES
(
'#{seqNo}'
,'#{menuNm}'
,'#{restId}'
,1
,'#{menuPrice}'
,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
,'#{useYn}'
,'#{codeId}'
,IF('#{productCd}' = '',NULL,'#{productCd}')
,IF('#{salesEndDate}' = '','20991231','#{salesEndDate}')
,IF('#{itemImage}' = '',NULL,'#{itemImage}')
)
`

var ItemNewSeq = `SELECT 
	LPAD((SELECT item_no from dar_sale_item_info ORDER BY item_no DESC LIMIT 1)+1,10,0) as seqNo`

var UpdateMenu = `
UPDATE dar_sale_item_info 
SET 
	ITEM_NM = '#{menuNm}'
	,ITEM_PRICE = '#{menuPrice}'
	,USE_YN = '#{useYn}'
	,ITEM_MENU = '#{codeId}'
	,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
	,PROD_ID = '#{prodId}'
	,BEST_YN = '#{bestYn}'
	,ITEM_IMG = IF('#{imagePath}' = '',NULL,'#{imagePath}')
	,SALES_END_DATE = IF('#{salesEndDate}' = '',NULL,'#{salesEndDate}')
WHERE 
	REST_ID = '#{restId}'
AND	ITEM_NO = '#{menuNo}'`

var UpdateMenuUseYn = `
UPDATE dar_sale_item_info 
SET 
	USE_YN = '#{useYn}'
WHERE 
	REST_ID = '#{restId}'
AND	ITEM_MENU = '#{codeId}'
`

var GetStoreService = `SELECT SERVICE_ID
										, SERVICE_NM
										, SERVICE_INFO
										, NOTICE_YN
										, USE_YN
										FROM priv_rest_service
										WHERE 
										REST_ID = '#{restId}'
										`

var InsertStoreBaseService = `INSERT INTO priv_rest_service(REST_ID, SERVICE_ID, SERVICE_NM, SERVICE_INFO, USE_YN)
									SELECT 
										'#{restId}'
										, SERVICE_ID
										, SERVICE_NM
										, SERVICE_INFO
										, USE_YN
									FROM priv_rest_service
									WHERE 
									REST_ID='R0000000000'
									`

var InsertStoreService = `
INSERT INTO priv_rest_service (REST_ID, SERVICE_ID, SERVICE_NM, SERVICE_INFO, USE_YN, NOTICE_YN, MOD_DATE)
SELECT 
		'#{restId}'
		,count(*)+1 
		,'#{serviceNm}'
		,'#{serviceInfo}'
		,'#{useYn}'
		,'#{noticeYn}'
		,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
FROM priv_rest_service 
where rest_id = '#{restId}'
`

var UpdateStoreService string = `UPDATE priv_rest_service SET
										SERVICE_INFO  = '#{serviceInfo}'
										, SERVICE_NM  = '#{serviceNm}'
										, USE_YN  = '#{useYn}'
										, NOTICE_YN  = '#{noticeYn}'
										, MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										WHERE 
											REST_ID='#{restId}'
											AND SERVICE_ID ='#{serviceId}'
											`

var SelectBankList = `
select CODE_ID, CODE_NM from b_code 
WHERE 
	category_id = 'BANK' 
	and USE_YN = 'Y'
`

var CheckPaymentUse = `
SELECT PAYMENT_USE_YN as allUseYn 
FROM priv_rest_info 
WHERE 
	rest_id = '#{restId}'
`

var UpdatePaymentUse = `
UPDATE priv_rest_info 
SET 
	PAYMENT_USE_YN = '#{allUseYn}'
WHERE 
	rest_id = '#{restId}'`

var SelectChargeList = `
SELECT 	SEQ_NO as seqNo
		,AMT as amt
		,ADD_AMT as addAmt
		,USE_YN as useYn
FROM dar_prepayment_charge_info 
WHERE 
	rest_id = '#{restId}'
`

var UpdateChargeItem = `
UPDATE dar_prepayment_charge_info 
SET 
	amt = '#{amt}'
	,add_amt = '#{addAmt}'
	,USE_YN = '#{useYn}'
WHERE 
	seq_no = '#{seqNo}'
AND	rest_id = '#{restId}'
`

var SelectChargeNewSeqNo = `
SELECT LPAD(max(seq_no)+1,10,0) as newSeqNo
FROM dar_prepayment_charge_info
`

var InsertChargeItem = `
INSERT INTO dar_prepayment_charge_info (
SEQ_NO
,REST_ID
,AMT
,ADD_AMT
,USE_YN
) VALUES (
'#{newSeqNo}'
,'#{restId}'
,'#{amt}'
,'#{addAmt}'
,'#{useYn}'
)
`

var SelectUnpaidListCount string = `SELECT COUNT(*) AS orderCnt
								, SUM(total_amt) AS TOTAL_AMT
								,FN_GET_GRPNAME(A.GRP_ID) AS BOOK_NM
								,(SELECT USER_ID FROM PRIV_GRP_USER_INFO AS AA WHERE AA.GRP_ID=A.GRP_ID AND AA.GRP_AUTH='0') AS USER_ID
								FROM  dar_order_info AS A
								INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID
								WHERE 
								A.REST_ID='#{restId}'
								AND A.order_ty IN ('1','2','3','5')
								AND A.GRP_ID ='#{grpId}'
								AND PAY_TY='1' 
								AND PAID_YN='N'
								AND order_stat = '20'
								AND DATE_FORMAT(A.ORDER_DATE,'%Y-%m-%d') <='#{accStDay}'
								`

var SelectUnpaidList string = `SELECT DATE_FORMAT(A.ORDER_DATE,'%Y-%m-%d %H:%i:%s')  AS ORDER_DATE
								,ORDER_TY
								,B.USER_NM AS orderer
								,CASE WHEN A.ORDER_TY ='1' THEN 'pay'
																			WHEN A.ORDER_TY ='2' THEN 'delivery'
																			WHEN A.ORDER_TY ='3' THEN 'takeout'
																			WHEN A.ORDER_TY ='5' THEN ''
																			END AS ORDER_TY
								,A.TOTAL_AMT
								FROM  dar_order_info AS A
								INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID
								WHERE 
								A.REST_ID='#{restId}'
								AND A.order_ty IN ('1','2','3','5')
								AND A.GRP_ID ='#{grpId}'
								AND PAY_TY='1' 
								AND PAID_YN ='N'
								AND order_stat = '20'
								AND DATE_FORMAT(A.ORDER_DATE,'%Y-%m-%d') <='#{accStDay}'
								ORDER BY ORDER_DATE DESC
								`

var SelectStorePaymentListCount string = `SELECT COUNT(*) AS orderCnt
								FROM dar_payment_hist AS a
								LEFT OUTER JOIN priv_rest_info AS b ON a.REST_ID = b.rest_id
								LEFT OUTER JOIN priv_grp_info AS C ON A.GRP_ID = C.GRP_ID
								WHERE 
									DATE_FORMAT(A.REG_DATE, '%Y-%m-%d') >=  '#{startDate}'
									AND DATE_FORMAT(A.REG_DATE, '%Y-%m-%d') <=  '#{endDate}'
									AND A.GRP_ID = '#{grpId}'
									AND A.REST_ID = '#{restId}'
									AND A.PAY_INFO IN ('0','3')
								ORDER BY A.REG_DATE DESC
								`

var SelectStorePaymentList string = `SELECT A.MOID
										,A.GRP_ID
										,C.GRP_NM
										,A.USER_TY
										,A.PAY_CHANNEL
										,A.PAYMENT_TY
										,ACC_ST_DAY
										,A.CREDIT_AMT 
										,DATE_FORMAT(A.REG_DATE,'%Y-%m-%d %H:%i:%s') as PAYMENT_DT
										,A.PAY_INFO
								FROM dar_payment_hist AS a
								LEFT OUTER JOIN priv_rest_info AS b ON a.REST_ID = b.rest_id
								LEFT OUTER JOIN priv_grp_info AS C ON A.GRP_ID = C.GRP_ID
								WHERE 
									DATE_FORMAT(A.REG_DATE, '%Y-%m-%d') >=  '#{startDate}'
									AND DATE_FORMAT(A.REG_DATE, '%Y-%m-%d') <=  '#{endDate}'
									AND A.GRP_ID = '#{grpId}'
									AND A.REST_ID = '#{restId}'
									AND A.PAY_INFO IN ('0','3')
								ORDER BY A.REG_DATE DESC
								`

var UpdateOrderPaid string = `UPDATE dar_order_info SET PAID_YN='Y'
												,moid = '#{moid}'
												,PAY_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE
								REST_ID='#{restId}'
								AND GRP_ID ='#{grpId}'
								AND PAY_TY='1' 
								AND PAID_YN='N'
								AND DATE_FORMAT(ORDER_DATE,'%Y-%m-%d') <='#{accStDay}'
								`

var InsertPaymentHistory string = `INSERT INTO dar_payment_hist
									(
										HIST_ID
										, REST_ID
										, GRP_ID
										, USER_ID
										, CREDIT_AMT
										, USER_TY
										, SEARCH_TY
										, PAYMENT_TY
										 ,PAY_INFO
										, REG_DATE
										, PAY_CHANNEL
										, ADD_AMT
										, MOID
										if #{accStDay} != '' then ,ACC_ST_DAY
									)
									VALUES
									(
										( SELECT FN_GET_SEQUENCE('PAYMENT_HIST_ID') AS TMP )
										, '#{restId}'
										, '#{grpId}'
										, '#{userId}'
										, '#{creditAmt}'
										, '#{userTy}'
										, '#{searchTy}'
										, '#{paymentTy}'
										, '#{payInfo}'
										, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										, '#{payChannel}'
										, '#{addAmt}'
										, '#{moid}'
										, '#{accStDay}'
									)`

var SelectStoreCancelCnt string = `SELECT COUNT(*) as CancelCnt
											FROM dar_payment_hist
											WHERE 
											MOID='#{moid}'
											AND PAYMENT_TY IN ('1','4')
											`

var SelectStoreChargeInfo string = `SELECT A.MOID
									,DATE_FORMAT(A.REG_DATE,'%Y.%m.%d %p%h:%i')  AS REG_DATE
									,B.GRP_NM
									,A.PAYMENT_TY
									,A.CREDIT_AMT
									,A.ADD_AMT
									,A.CREDIT_AMT + A.ADD_AMT AS TOTAL_AMT
									,IFNULL(DATE_ADD(C.SETTLMNT_DT, INTERVAL 10 DAY),'') AS expectInDate
									,A.USER_ID
									,A.PAY_INFO
									,A.GRP_ID AS BOOK_ID
									,ACC_ST_DAY
								FROM dar_payment_hist AS A
								INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID AND PAYMENT_TY in ('0','3')
								LEFT OUTER JOIN dar_rest_payment AS C ON A.REST_PAYMENT_ID = C.REST_PAYMENT_ID
								WHERE 
								A.MOID='#{moid}'
								AND A.REST_ID='#{restId}'
								`

var UpdateOrderPaidCancel string = `UPDATE dar_order_info SET PAID_YN='N'
								,moid = NULL
								,PAY_DATE = NULL
				WHERE
				moid= '#{moid}'
				`
