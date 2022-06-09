package mngAdmin

var SelectFeeRateListCnt = `
SELECT 
	count(*) as TOTAL_COUNT
FROM dar_rest_fees as A
INNER JOIN priv_rest_info as B ON A.REST_ID = B.REST_ID
WHERE 1=1
`

var SelectFeeRateList = `
SELECT 
	A.REST_ID as restId
    ,A.REST_NM as restNm
    ,B.CEO_NM as ceoNm
	,B.TEL as tel
    ,A.PAYMETHOD as paymentHod
    ,A.REST_FEES as restFee
    ,DATE_FORMAT(A.START_DATE,'%Y-%m-%d') as startDate
    ,DATE_FORMAT(A.END_DATE,'%Y-%m-%d') as endDate
    ,A.USE_FEES_YN as useYn
FROM dar_rest_fees as A
INNER JOIN priv_rest_info as B ON A.REST_ID = B.REST_ID
WHERE
      1=1
`

var UpdateFeeRate = `
UPDATE dar_rest_fees SET
	USE_FEES_YN = '#{useYn}' 
	,END_DATE = DATE_FORMAT(now(),'%Y%m%d%H%i%s')
WHERE
    REST_ID = '#{restId}'
    AND PAYMETHOD = '#{paymentHod}'
`

var InsertFeeRate = `
INSERT INTO dar_rest_fees (
	REST_ID
	,REST_NM
	,PAYMETHOD
	,START_DATE
	,END_DATE
	,USE_FEES_YN
	,REST_FEES
) SELECT 
	REST_ID
	, REST_NM
	, '#{paymentHod}' as PAYMETHOD
	, DATE_FORMAT(now(), '%Y%m%d%H%i%s') as START_DATE
	, '99991231235959' as START_DATE
	, 'Y' as USE_FEES_YN
	, '#{restFee}' as REST_FEES
FROM priv_rest_info 
WHERE
	REST_ID = '#{restId}'
`
