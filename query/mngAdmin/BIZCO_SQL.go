package mngAdmin

var SelectCompanyList = `
select a.COMPANY_ID
     , a.COMPANY_NM
     , IFNULL(b.HP_NO,'') as USER_HP
     , IFNULL(b.USER_NM,'') as USER_NM
     , a.USE_YN
     , a.BUSID
     , IF(ISNULL(b.USER_NM), 'N', 'Y') as BIZ_YN
     , (SELECT COUNT(*)
        FROM b_company_book as AA
                 INNER JOIN priv_grp_info as AAA on AAA.grp_id = AA.book_id
        where AAA.USE_YN = 'Y'
          and A.company_id = AA.company_id) AS BOOK_CNT
     , (SELECT COUNT(*) AS store_cnt
        FROM (SELECT AA.company_id, bb.rest_id, COUNT(*) AS CNT
              FROM b_company_book AS aa
                       INNER JOIN org_agrm_info AS bb ON aa.BOOK_ID = bb.grp_id AND bb.REQ_STAT = '1'
              GROUP BY AA.company_id, bb.rest_id
             ) AS CC
        WHERE A.company_id = CC.company_id) AS STORE_CNT
     , (SELECT COUNT(*)
        FROM (SELECT AA.company_id, aa.user_id, COUNT(*) AS CNT
              FROM b_company_user AS aa
              WHERE aa.USER_ID IS NOT null
              GROUP BY AA.company_id, aa.user_id) AS DD
        WHERE A.company_id = DD.company_id) AS USER_CNT
FROM b_company as a
         left outer join (select ba.USER_ID, bb.USER_NM, ba.company_id, bb.HP_NO
                          FROM b_company_manager as ba
                                   INNER JOIN priv_user_info as bb
                                              on ba.USER_ID = bb.USER_ID and ba.AUTHOR_CD = 'CM') as b
                         on a.company_id = b.company_id
WHERE 1=1`

var SelectCompanyCnt = `
SELECT COUNT(*) as TOTAL_COUNT
FROM b_company as a
         LEFT OUTER JOIN (select ba.USER_ID, bb.USER_NM, ba.company_id, bb.HP_NO
                          FROM b_company_manager as ba
                                   INNER JOIN priv_user_info as bb
                                              on ba.USER_ID = bb.USER_ID and ba.AUTHOR_CD = 'CM') as b
                         on a.company_id = b.company_id
WHERE 1=1`

var SelectCompanyId = `
SELECT A.COMPANY_ID
     , A.COMPANY_NM
FROM b_company AS A
         left outer join b_company_manager as B on A.company_id = B.company_id
WHERE USE_YN = 'Y'
  AND A.COMPANY_NM LIKE '%#{companyNm}%'
  AND B.USER_ID IS NOT NULL
GROUP BY A.company_id
`

var SelectCompanyInfo = `
SELECT company_id as companyId,
       company_nm as companyNm,
       ADDR       as addr,
       ADDR2      as addr2,
       TEL        as tel,
       homepage   as homePage,
       BUSID      as bisId,
       CEO_NM     as ceoNm,
	   USE_YN     as useYn 
FROM b_company
WHERE 
	company_id = '#{companyId}'
`

var InsertCompany = `
INSERT INTO b_company (
                       company_nm
                      , busid
                      , reg_date)
VALUES (
        '#{companyNm}'
       , '#{bizNum}'
       , DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'))
`

//SelectDataForAddCompany priv_user_info, PRIV_GRP_INFO, b_company, b_company_code 데이터 추출
var SelectDataForAddCompany = `
SELECT USER_ID                                       as userId,
	   USER_NM                                       as userNm,
	   HP_NO                                         as hpNo,
       LOGIN_ID                                      as loginId, 
       LOGIN_PW                                      as loginPw, 
       EMAIL                                         as email,
       IFNULL(KAKAO_KEY, '')                         as kakaoKey, 
       IFNULL(APPLE_KEY, '')                         as appleKey, 
       IFNULL(NAVER_KEY, '')                         as naverKey, 
       (SELECT CONCAT('B', IFNULL(LPAD(MAX(SUBSTRING(GRP_ID, -10)) + 1, 10, 0), '0000000001')) 
        FROM PRIV_GRP_INFO)                          as grpId 
FROM priv_user_info 
WHERE 1=1  
	AND LOGIN_ID = '#{loginId}' 
`

// SelectDataForTransBiz 기업의 비즈 사용전환을 위한 데이터 추출
var SelectDataForTransBiz = `
SELECT a.company_id as companyId, 
       a.company_nm as companyNm, 
       b.BOOK_ID as grpId, 
       c.USER_ID as userId, 
       d.HP_NO as hpNo, 
       d.USER_NM as userNm, 
       d.LOGIN_ID as loginId, 
       IFNULL(d.KAKAO_KEY, '') as kakaoKey, 
       IFNULL(d.APPLE_KEY, '') as appleKey, 
       IFNULL(d.NAVER_KEY, '') as naverKey, 
       d.LOGIN_PW as loginPw 
FROM b_company as a 
         inner join b_company_book as b on a.company_id = b.company_id 
         left outer join priv_grp_user_info as c on c.GRP_ID = b.BOOK_ID and c.GRP_AUTH = '0' 
         left outer join priv_user_info as d on c.USER_ID = d.USER_ID 
WHERE 1 = 1 
  AND a.company_id = '#{companyId}' 
`

//SelectCompanyDept b_company_code 부서 데이터 추출
var SelectCompanyDept = `
SELECT * 
FROM b_company_code 
WHERE 1 = 1 
  AND COMPANY_ID = '#{companyId}'
  AND CODE_TY = '#{codeTy}'
`

//InsertCompanyDept b_company_code 회사 부서 추가
var InsertCompanyDept = `
INSERT INTO b_company_code (
	company_id, 
	code_ty, 
	code_info) 
VALUES (
	'#{companyId}', 
	'#{codeTy}', 
	'#{codeInfo}'
)`

//InsertSysUserInfo sys_user_info 데이터 삽입 **userId = loginId
var InsertSysUserInfo = `
INSERT INTO sys_user_info ( user_id
                          , user_nm
                          , user_pass
                          , MSG_LANG_CD
                          , conn_allow_yn
                          , last_conn_date
                          , user_menu_author_yn
                          , author_cd
                          , init_pass_yn
                          , OTP_YN
                          , ADDED_BY
                          , last_pass_upd_date
                          , added_date)
VALUES (
		'#{loginId}',
        '#{userNm}',
        '#{loginPw}',
        'ko_KR',
        'Y',
        DATE_FORMAT(NOW(), '%Y-%m-%d- %H:%i:%s'),
        'N',
        '#{authorCd}',
        'N',
        'N',
        IF('#{addedBy}' = '<null>', null, '#{addedBy}'),
        DATE_FORMAT(NOW(), '%Y-%m-%d- %H:%i:%s'),
        DATE_FORMAT(NOW(), '%Y-%m-%d- %H:%i:%s'))
`

//InsertCompanyBook b_company_book 데이터 삽입
var InsertCompanyBook = `
INSERT INTO b_company_book (
	company_id, 
	BOOK_ID)
VALUES (
	'#{companyId}', 
	'#{grpId}')
`

//SelectGrpUsers priv_grp_user_info 장부 유저 추출
var SelectGrpUsers = `
SELECT a.USER_ID as userId, b.USER_NM as userNm, b.HP_NO as hpNo, IF(ISNULL(a.LEAVE_TY),'Y','N') as useYn, a.GRP_AUTH as grpAuth, b.EMAIL as email
FROM priv_grp_user_info as a 
         INNER JOIN priv_user_info as b on a.USER_ID = b.USER_ID 
WHERE 1 = 1 
  AND GRP_ID = '#{grpId}' 
`

//InsertCompanyUser b_company_user 회사유저 삽입
var InsertCompanyUser = `
INSERT INTO b_company_user (
							COMPANY_ID,
                            USER_ID,
                            BOOK_ID,
                            USER_NM,
                            HP_NO,
                            DEPT,
                            REG_DATE,
                            USE_YN
) VALUES (
		'#{companyId}',
        '#{userId}',
        '#{grpId}',
        '#{userNm}',
        '#{hpNo}',
        '#{codeId}',
        DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'),
        '#{useYn}'
) 
`

//UpdateCompanyUseYn b_company useYn 갱신
var UpdateCompanyUseYn = `
UPDATE b_company
SET 
	USE_YN = '#{useYn}'
WHERE 
      company_id = '#{companyId}'
`

//InsertCompanyManager b_company_manager 데이터 삽입
var InsertCompanyManager = `
INSERT INTO b_company_manager (company_id, USER_ID, REG_DATE, AUTHOR_CD, DEPT, TEL, EMAIL)
VALUES (
        '#{companyId}',
        '#{userId}',
        DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'),
        '#{authorCd}',
        '#{codeInfo}',
        '#{hpNo}',
        '#{email}'
       )
`

//InsertGrp PRIV_GRP_INFO 데이터 삽입
var InsertGrp = `
INSERT INTO PRIV_GRP_INFO (
                           GRP_ID,
                           GRP_NM,
                           BUSID,
                           CEO_NM,
                           ADDR,
                           ADDR2,
                           TEL,
                           EMAIL,
                           AUTH_STAT,
                           REG_DATE,
                           AUTH_DATE,
                           INTRO,
                           LIMIT_YN,
                           LIMIT_AMT,
                           LIMIT_DAY_AMT,
                           DETAIL_VIEW_YN,
                           GRP_TYPE_CD,
                           SUPPORT_AMT,
                           GRP_PAY_TY
) VALUES (
        '#{grpId}',
        '#{grpNm}',
        '#{bizNum}',
        '#{userNm}',
        '',
        '',
        '#{hpNo}',
        '#{email}',
        '#{authStat}',
        DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'),
        DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'),
        '',
        '#{limitYn}',
        '#{limitAmt}',
        '#{limitDayAmt}',
        '#{detailViewYn}',
        '#{grpTypeCd}',
        '#{supportAmt}',
        '#{grpPayTy}'
)
`

var InsertGrpUser = `
INSERT INTO priv_grp_user_info(
     grp_id
	,user_id
	,grp_auth
	,auth_stat
	,auth_date
	,join_ty
	,reg_date
) VALUES (
	'#{grpId}'
	,'#{userId}'
	,'#{grpAuth}'
	,'#{authStat}'
	,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
	,'0'
	,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'))`

var SelectCompanyManagers = `
SELECT A.COMPANY_ID        as companyId,
       B.USER_ID           as userId,
       A.AUTHOR_CD         as AuthorCd,
       B.HP_NO             as HpNo,
       IFNULL(B.EMAIL, '') as Email,
       B.USER_NM           as userNm
FROM b_company_manager AS a
         INNER JOIN priv_user_info AS b ON a.user_id = b.user_id
WHERE 1 = 1
  AND a.company_id = '#{companyId}'
`

var SelectCompanyUsers = `
SELECT distinct a.USER_ID as user_id, B.USER_NM as userNm, A.HP_NO as HpNo FROM b_company_user AS a
INNER JOIN priv_user_info AS b ON a.USER_ID = b.USER_ID 
LEFT JOIN b_company_manager as c ON a.user_id = c.user_id and c.company_id = a.COMPANY_ID
WHERE 
a.COMPANY_ID = '#{companyId}'
AND a.USE_YN = 'Y'
AND c.USER_ID IS NULL
`

var SelectCompanyUsersCnt = `
SELECT count(distinct a.user_id) as total FROM b_company_user AS a
INNER JOIN priv_user_info AS b ON a.USER_ID = b.USER_ID 
LEFT JOIN b_company_manager as c ON a.user_id = c.user_id
WHERE 
a.COMPANY_ID = '#{companyId}'
AND a.USE_YN = 'Y'
AND c.USER_ID IS NULL
`

var InsertCompanyManagers = `
INSERT INTO b_company_manager (company_id
			,user_id
			,reg_date
			,mod_date
			,author_cd
			,dept,class
			,tel
			,email) 
SELECT
	'#{companyId}'
	,'#{userId}'
	,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
	,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
	,'#{authorCd}'
	,null
	,null
	,HP_NO
	,EMAIL
FROM priv_user_info
WHERE user_id = '#{userId}'
`

var UpdateCompanyManagerAuthorAllBM = `
UPDATE b_company_manager 
SET author_cd = 'BM' 
WHERE 
company_id = '#{companyId}' 
AND USER_ID != '#{userId}'`

var UpdateCompanyManagerAuthorCM = `
UPDATE b_company_manager 
SET author_cd = 'CM' 
WHERE 
company_id = '#{companyId}' 
AND USER_ID = '#{userId}'`

var SelectCheckLoginIdInSysInfo = `
select * from sys_user_info where 
USER_ID = '#{loginId}' 
`
