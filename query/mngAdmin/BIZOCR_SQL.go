package mngAdmin

var SelectTextList = `
SELECT a.rest_id as restId
		,a.biz_num as bizNum
		,a.text_id as textId
		,a.receipt_id as receiptId
		,a.text as texts
		,a.is_menu as isMenu
		,a.menu_id as menuId
		,IFNULL(b.menu, '') as menuNm
		,IFNULL(b.price, '') as price
FROM cc_ocr_text as a 
LEFT OUTER JOIN cc_ocr_menu as b on a.menu_id = b.menu_id and a.biz_num = b.biz_num
WHERE 
is_menu = '#{isMenu}'
AND a.biz_num LIKE '%#{bizNum}%'
ORDER BY a.BS_DT, a.BS_TIME ASC
`

var SelectTextData = `
SELECT 	b.biz_num as bizNum
		,IFNULL(CONCAT('/app/SharedStorage/receipt','/',a.receipt_image_name),'') AS receiptImg 
		,b.text_id as textId
		,b.text as texts 
FROM cc_ocr_text as b 
INNER JOIN cc_ocr_receipt as a on b.receipt_id = a.receipt_id and b.biz_num = a.biz_num
WHERE 
		b.text_id = '#{textId}'
AND		b.biz_num = '#{bizNum}'
AND		b.RECEIPT_ID = '#{receiptId}'
`

var UpdateOCRTextMenuId = `
UPDATE cc_ocr_text 
SET is_menu = '#{isMenu}'
	,menu_id = IF('#{isMenu}'='Y','#{menuId}',0) 
WHERE 
	biz_num = '#{bizNum}' 
AND text_id = '#{textId}'
`

var SelectNewMenuId = `
SELECT IFNULL(MAX(MENU_ID),0) + 1 as newMenuId
		FROM cc_ocr_menu
WHERE 
	BIZ_NUM = '#{bizNum}' 
AND REST_ID = '#{restId}'
`

var SelectMenuChecker = `
SELECT 
CASE  
WHEN menu = '#{menuNM}' AND price = '#{menuPrice}' THEN 2
WHEN menu = '#{menuNM}' THEN 1
ELSE 0
END  as state
FROM cc_ocr_menu
WHERE 	
	biz_num = '#{bizNum}' 
AND use_yn = 'Y'
GROUP BY state 
ORDER BY state DESC LIMIT 1`

var SelectMenuCheck = `
SELECT 	IFNULL(a.stat,0) as state
		, a.menuId
		, (SELECT IFNULL(max(menu_id)+1,1) as newMenuId 
			FROM cc_ocr_menu as b 
			WHERE 
				biz_num = '#{bizNum}') as newMenuId 
FROM (SELECT IF(menu = '#{menuNm}'
			,IF(price = '#{menuPrice}',2,1),0) AS stat
			, menu_id AS menuId  
		FROM cc_ocr_menu
WHERE 
	biz_num = '#{bizNum}' 
AND use_yn = 'Y'
GROUP BY stat) AS a
ORDER BY stat DESC LIMIT 1
`

var InsertOCRNewMenu = `
INSERT INTO cc_ocr_menu 
		(rest_id
		, biz_num
		, menu_id
		, menu
		, price
		, use_yn
		, bs_dt
		, bs_time)
VALUES (
		'#{restId}'
		,'#{bizNum}'
		,'#{newMenuId}'
		,'#{menuNm}'
		,'#{menuPrice}'
		,'Y'
		,DATE_FORMAT(now(),'%Y%m%d')
		,DATE_FORMAT(now(),'%H%i%s'))
`

var UpdateOCRMenuYn = `
UPDATE cc_ocr_menu 
SET use_yn = 'N'
WHERE 
	menu_id = '#{menuId}' 
AND biz_num = '#{bizNum}'
`

var SelectReceiptList = `
SELECT 	biz_num as bizNum
		,receipt_id as receiptId
		,response_id as responseId
		,STATUS as state
		,DATE_FORMAT(scan_dt,'%Y-%m-%d') as scanDt
		,aprv_no as aprvNo
FROM cc_ocr_receipt
WHERE 
	STATUS = '#{status}'
AND biz_num LIKE '%#{bizNum}%'
`

var SelectReceiptData = `
SELECT 	a.biz_num as bizNum
		,a.rest_id as restId
		,a.status as state
		,a.STATUS_DETAIL as stateDetail
		,DATE_FORMAT(a.scan_dt,'%Y-%m-%d') as scanDt
		,DATE_FORMAT(a.aprv_dt,'%Y-%m-%d') as aprvDt
		,a.aprv_no as aprvNo
		,a.tot_amt as totalAmt
		,IFNULL(CONCAT('/app/SharedStorage/receipt','/',a.receipt_image_name),'') AS receiptImg 
FROM cc_ocr_receipt as a 
WHERE 
		a.receipt_id = '#{receiptId}'
AND		a.biz_num = '#{bizNum}'
`

var SelectReceiptMenuData = `
SELECT 	IFNULL(a.receipt_menu_id,'') as receiptMenuId
		, a.rest_id as restId
		, a.menu_id as menuId
		, a.receipt_id as receiptId
		, a.quantity as quantity
		, a.biz_num as bizNum
		, a.amt as totalPrice
		, b.price as menuPrice
		, b.menu as menuNm
FROM cc_ocr_receipt_menu as a
INNER JOIN cc_ocr_menu as b on b.biz_num = a.biz_num
							and a.menu_id = b.menu_id
WHERE 
	a.receipt_id = '#{receiptId}'
and	a.biz_num = '#{bizNum}'
`

var InsertReceiptMenuData = `
INSERT INTO cc_ocr_receipt_menu (REST_ID, BIZ_NUM, RECEIPT_MENU_ID, MENU_ID, RECEIPT_ID, QUANTITY, AMT) 
VALUES (
'#{restId}'
,'#{bizNum}'
,'#{receiptMenuId}'
,'#{newMenuId}'
,'#{receiptId}'
,'#{menuEa}'
,'#{menuAmt}'
)
`

var UpdateReceiptMenuData = `
UPDATE cc_ocr_receipt_menu 
SET 
	MENU_ID = '#{newMenuId}'
	, QUANTITY = '#{menuEa}'
	, AMT = '#{menuAmt}'
WHERE 
	REST_ID = '#{restId}'
AND BIZ_NUM = '#{bizNum}'
AND RECEIPT_MENU_ID = '#{receiptMenuId}'
AND RECEIPT_ID = '#{receiptId}'
`

var SelectNewReceiptMenuId = `
SELECT IFNULL(max(RECEIPT_MENU_ID)+1,1) as newReceiptMenuId 
FROM cc_ocr_receipt_menu 
WHERE 
	REST_ID = '#{restId}' 
AND BIZ_NUM = '#{bizNum}' 
AND RECEIPT_ID = '#{receiptId}'
`

var UpdateReceiptData = `
UPDATE cc_ocr_receipt 
SET
	STATUS = '#{state}'
	,TOT_AMT = '#{totalAmt}'
	,APRV_NO = '#{aprvNo}'
	,APRV_DT = '#{aprvDt}'
WHERE 
	rest_id = '#{restId}'
AND BIZ_NUM = '#{bizNum}'
AND RECEIPT_ID = '#{receiptId}'
`

var SelectReceiptMenuList = `
SELECT receipt_menu_id as receiptMenuId 
FROM cc_ocr_receipt_menu 
WHERE 
	rest_id = '#{restId}' 
and biz_num = '#{bizNum}' 
and receipt_id = '#{receiptId}'
`

var DeleteReceiptMenuData = `
DELETE FROM cc_ocr_receipt_menu 
WHERE 
	rest_id = '#{restId}' 
and biz_num = '#{bizNum}' 
and receipt_id = '#{receiptId}'
and receipt_menu_id = '#{receiptMenuId}'`
