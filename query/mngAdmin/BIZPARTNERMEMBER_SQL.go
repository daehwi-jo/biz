package mngAdmin

var SelectPartnerMemberListCnt = `
SELECT Count(*) as TOTAL_COUNT
FROM (
         select S1.STORE_ID
              , S1.USER_ID
              , S1.END_DATE
              , S2.REST_ID
              , S2.REST_NM
              , S2.BUSID
              , S2.USE_YN
         from e_billing as S1
                  inner JOIN priv_rest_info AS S2 ON S1.STORE_ID = S2.REST_ID
     ) AS A
         left OUTER JOIN cc_sync_inf AS B
                         ON A.BUSID = B.BIZ_NUM
                             AND B.SITE_CD = '1'
                             AND B.BS_DT = (select max(S3.BS_DT)
                                            from cc_sync_inf AS S3
                                            where S3.BIZ_NUM = A.BUSID
                                              and S3.SITE_CD = '1')

         LEFT OUTER JOIN cc_sync_inf AS C
                         ON A.BUSID = C.BIZ_NUM
                             AND C.SITE_CD = '2'
                             AND C.BS_DT = (select max(S4.BS_DT)
                                            from cc_sync_inf AS S4
                                            WHERE S4.BIZ_NUM = A.BUSID
                                              AND S4.SITE_CD = '2')
         LEFT OUTER JOIN sys_alimtalk_log AS D
                         ON A.USER_ID = D.USER_ID
                             AND D.SEQ =
                                 (select max(S5.SEQ)
                                  from sys_alimtalk_log AS S5
                                  WHERE S5.USER_ID = A.USER_ID)
WHERE 
	A.REST_NM LIKE '%#{restNm}%'
`

// y**** = 여신협회 , h**** = 홈택스
var SelectPartnerMemberList = `
SELECT A.STORE_ID                                                 AS restId
     , A.USER_ID                                                  AS userId
     , A.END_DATE                                                 AS endDate
     , A.REST_NM                                                  AS restNm
     , A.BUSID                                                    AS BizNum
     , IF(A.END_DATE >= DATE_FORMAT(now(), '%Y-%m-%d'), 'Y', 'N') AS partnerMemberUseYn
     , A.USE_YN                                                   AS storeUseYn
     , IFNULL(DATE_FORMAT(B.BS_DT,'%Y-%m-%d'), 'null')                                    AS yBsDt
     , IFNULL(B.STS_CD, 'null')                                   AS yStsCd
     , IFNULL(B.ERR_CD, 'null')                                   AS yErrCd
     , IFNULL(DATE_FORMAT(C.BS_DT,'%Y-%m-%d'), 'null')                                    AS hBsDt
     , IFNULL(C.STS_CD, 'null')                                   AS hStsCd
     , IFNULL(C.ERR_CD, 'null')                                   AS hErrCd
     , IFNULL(DATE_FORMAT(D.SEND_DATE,'%Y-%m-%d'), 'null') 		  AS sendDt
     , IFNULL(D.TEMPLATE_CODE, 'null')                            AS templateCode
     , IFNULL(D.TEMPLATE_NAME, 'null')                            AS templateName
FROM (
         select S1.STORE_ID
              , S1.USER_ID
              , S1.END_DATE
              , S2.REST_ID
              , S2.REST_NM
              , S2.BUSID
              , S2.USE_YN
         from e_billing as S1
         INNER JOIN priv_rest_info AS S2 ON S1.STORE_ID = S2.REST_ID
     ) AS A
         LEFT OUTER JOIN cc_sync_inf AS B
                         ON A.BUSID = B.BIZ_NUM
                             AND B.SITE_CD = '1'
                             AND B.BS_DT = (select max(S3.BS_DT)
                                            from cc_sync_inf AS S3
                                            where S3.BIZ_NUM = A.BUSID
                                              and S3.SITE_CD = '1')

         LEFT OUTER JOIN cc_sync_inf AS C
                         ON A.BUSID = C.BIZ_NUM
                             AND C.SITE_CD = '2'
                             AND C.BS_DT = (select max(S4.BS_DT)
                                            from cc_sync_inf AS S4
                                            WHERE S4.BIZ_NUM = A.BUSID
                                              AND S4.SITE_CD = '2')
         LEFT OUTER JOIN sys_alimtalk_log AS D
                         ON A.USER_ID = D.USER_ID
                             AND D.SEQ =
                                 (select max(S5.SEQ)
                                  from sys_alimtalk_log AS S5
                                  WHERE S5.USER_ID = A.USER_ID)
WHERE 
	A.REST_NM LIKE '%#{restNm}%'
`

var SelectPartnerMemberDate = `
SELECT STORE_ID, START_DATE, END_DATE, NEXT_PAY_DAY
FROM e_billing
WHERE 
	STORE_ID = '#{restId}'
`

var UpdatePartnerMemberDate = `
UPDATE e_billing
SET END_DATE = '#{endDate}'
    , NEXT_PAY_DAY = '#{endDate}'
    , MOD_DATE = DATE_FORMAT(now(),'%Y-%m-%d %H:%i:%s')
WHERE 
	STORE_ID = '#{restId}'
`

var SelectPartnerMemberInfoData = `
SELECT REST_ID as restId
     , B.REST_NM AS restNm
     , IF(A.END_DATE > DATE_FORMAT(now(), '%Y-%m-%d'), 'Y', 'N') AS partnerMemberYn
     , B.USE_YN AS useYn
     , A.START_DATE AS startDate
     , A.END_DATE AS endDate
     , A.NEXT_PAY_DAY AS nextPayDate
     , A.ITEM_CODE AS itemCd
     , B.BUSID AS bizNum
     , B.CEO_NM AS ceoNm
     , B.TEL AS tel
     , A.PAY_YN as payYn
FROM e_billing as A
         INNER JOIN priv_rest_info as B ON A.STORE_ID = B.REST_ID
WHERE 
	STORE_ID = '#{restId}'
`

var SelectPartnerMemberCollectList = `
SELECT DATE_FORMAT(C.BS_DT, '%Y-%m-%d') as BsDt
     , IFNULL(C.STS_CD, 'null')         as yStsCd
     , IFNULL(C.ERR_CD, 'null')         as yErrCd
     , IFNULL(D.STS_CD, 'null')         as hStsCd
     , IFNULL(D.ERR_CD, 'null')         as hErrCd
FROM priv_rest_info AS A
        LEFT OUTER JOIN cc_sync_inf AS C
                         ON A.BUSID = C.BIZ_NUM
                             AND C.SITE_CD = '1'
        LEFT OUTER JOIN cc_sync_inf AS D
                         ON A.BUSID = D.BIZ_NUM
                             AND D.SITE_CD = '2'
                            AND C.BS_DT = D.BS_DT
WHERE 
  A.REST_ID = '#{restId}'
  AND DATE_FORMAT(C.BS_DT, '%Y-%m-%d') >= '#{startDate}'
  AND DATE_FORMAT(C.BS_DT, '%Y-%m-%d') <= '#{endDate}'
ORDER BY C.BS_DT DESC
`

var SelectPartnerMemberCollectListCnt = `
SELECT COUNT(*) AS TOTAL_COUNT
FROM priv_rest_info AS A
        LEFT OUTER JOIN cc_sync_inf AS C
                         ON A.BUSID = C.BIZ_NUM
                             AND C.SITE_CD = '1'
        LEFT OUTER JOIN cc_sync_inf AS D
                         ON A.BUSID = D.BIZ_NUM
                             AND D.SITE_CD = '2'
                            AND C.BS_DT = D.BS_DT
WHERE 
  A.REST_ID = '#{restId}'
  AND DATE_FORMAT(C.BS_DT, '%Y-%m-%d') >= '#{startDate}'
  AND DATE_FORMAT(C.BS_DT, '%Y-%m-%d') <= '#{endDate}'
ORDER BY C.BS_DT DESC
`

var SelectPartnerMemberAlarmList = `
SELECT DATE_FORMAT(B.SEND_DATE,'%Y-%m-%d %H:%i:%s') as sendDate,
		IFNULL(B.TEMPLATE_CODE, 'null')  as templateCd,
		B.HP_NO as hpNo,
		C.USER_NM as userNm
FROM e_billing as A
         INNER JOIN sys_alimtalk_log as B ON A.USER_ID = B.USER_ID
        INNER JOIN priv_user_info as C ON B.USER_ID = C.USER_ID
WHERE
    A.STORE_ID = '#{restId}'
AND DATE_FORMAT(B.SEND_DATE,'%Y-%m-%d') >= '#{startDate}'
AND DATE_FORMAT(B.SEND_DATE,'%Y-%m-%d') <= '#{endDate}'
ORDER BY B.SEND_DATE DESC
`

var SelectPartnerMemberAlarmListCnt = `
SELECT COUNT(*) as TOTAL_COUNT
FROM e_billing as A
         INNER JOIN sys_alimtalk_log as B ON A.USER_ID = B.USER_ID
        INNER JOIN priv_user_info as C ON B.USER_ID = C.USER_ID
WHERE
    STORE_ID = '#{restId}'
AND DATE_FORMAT(B.SEND_DATE,'%Y-%m-%d') >= '#{startDate}'
AND DATE_FORMAT(B.SEND_DATE,'%Y-%m-%d') <= '#{endDate}'
`
