package mngBook

var SelectBookGrpList = `SELECT C.GRP_ID
					,C.GRP_NM
					,DATE_FORMAT(C.REG_DATE, '%Y-%m-%d') as REG_DATE
					, D.USER_ID  
			   	, E.USER_NM  
			   	, GRP_TYPE_CD  
			   	,(select count(*) from priv_grp_user_info as  aa where B.BOOK_ID = aa.grp_id and aa.AUTH_STAT='1') as GRP_USER_CNT 
			   	,(select count(*) from org_agrm_info as  bb where B.BOOK_ID= bb.grp_id and bb.REQ_STAT='1') as GRP_REST_CNT 
			   	,c.intro
			FROM b_company AS A
			INNER JOIN b_company_book AS B ON A.COMPANY_ID = B.company_id
			INNER JOIN priv_grp_info AS C ON B.BOOK_ID = C.GRP_ID AND C.USE_YN='Y'
			inner join priv_grp_user_info as D on  C.grp_id = D.grp_id and D.GRP_AUTH=0
			inner join priv_user_info as E on D.USER_ID = E.user_id
			where 
				A.company_id = '#{companyId}'
			AND D.USER_ID LIKE '%#{userId}%'
			AND B.BOOK_ID IN (
						SELECT  GRP_ID
						FROM priv_grp_user_info as aa
						WHERE 
							aa.GRP_AUTH=0
			)`

//and aa.user_id = '#{searchUid}' 마지막줄에 (aa.GRP_AUTH=0 뒤에) 추가되면 값이 안나옴

var SelectBookGrpListCnt = ` SELECT COUNT(*) as TOTAL_COUNT
        						FROM b_company AS A
								INNER JOIN b_company_book AS B ON A.COMPANY_ID = B.company_id
								INNER JOIN priv_grp_info AS C ON B.BOOK_ID = C.GRP_ID AND C.USE_YN='Y'
								inner join priv_grp_user_info as D on  C.grp_id = D.grp_id and D.GRP_AUTH=0
								inner join priv_user_info as E on D.USER_ID = E.user_id
								where 
									A.company_id = '#{companyId}'
								AND D.USER_ID LIKE '%#{userId}%'
`

/*
AND B.BOOK_ID IN (SELECT GRP_ID
					FROM priv_grp_user_info as aa
					WHERE
					aa.GRP_AUTH=0
					and aa.user_id = '#{searchUid}')` 이부분 날림
*/

var SelectBookInfo = `select A.GRP_NM
				, A.GRP_TYPE_CD
				, CASE  WHEN B.CAT_CD IS NULL THEN 'N' ELSE 'Y' END AS CAT_YN
				, A.SUPPORT_FORWARD_YN
				, A.LIMIT_YN
				, IFNULL(A.LIMIT_AMT,0) AS LIMIT_AMT
				, IFNULL(A.LIMIT_DAY_AMT ,0) AS LIMIT_DAY_AMT
				, IFNULL(A.LIMIT_DAY_CNT ,0) AS LIMIT_DAY_CNT
				, A.INVITE_LINK
				, IFNULL(A.SUPPORT_AMT ,0) AS SUPPORT_AMT
				, A.support_yn
				, A.intro
		FROM priv_grp_info AS A
		LEFT OUTER JOIN priv_catering_grp_info AS B ON A.GRP_ID = B.GRP_ID
		where 
			A.grp_id = '#{searchGrpId}'`

// SelectGrpUserList 권한전달 모달 검색부분 사용시 쿼리문 추가
var SelectGrpUserList = `select B.USER_NM
        					  	,B.LOGIN_ID
        	  					,B.USER_ID 
						from b_company_manager AS A
						INNER JOIN PRIV_USER_INFO AS B ON A.USER_ID = B.USER_ID
						where 
							a.company_id = '#{companyId}'
							AND (('#{searchKey}' = 'userNm' AND B.USER_NM LIKE '%#{searchKeyword}%') 
							OR ('#{searchKey}' = 'loginId' AND B.LOGIN_ID LIKE '%#{searchKeyword}%')
							OR ('#{searchKey}' = 'userHp' AND B.HP_NO LIKE '%#{searchKeyword}%'))`

//이거 수정중
var UpdateGrpUser = `UPDATE PRIV_GRP_USER_INFO
		SET
			GRP_AUTH = '#{grpAuth}'
			,AUTH_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
		WHERE 
			GRP_ID = '#{grpId}'
			AND USER_ID = '#{userId}'`

//받아오는 값 없어서 빼어 놓음 SET 안에 들어가야함
//,AUTH_STAT = '#{authStat}'
//,JOIN_TY = '#{joinTy}'
//,DEPT_ID = '#{deptId}'
//,LEAVE_TY = '#{leaveTy}'
//,SUPPORT_BALANCE = '#{supportBalance}'

var DeleteGrp = `UPDATE PRIV_GRP_INFO
			SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
				,USE_YN ='N'
			WHERE 
				GRP_ID = '#{grpId}'`

var UpdateGrp = `UPDATE PRIV_GRP_INFO
		SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			,GRP_NM = '#{grpNm}'
			,GRP_TYPE_CD = '#{grpTypeCd}'
			,LIMIT_YN = '#{limitYn}'
			,LIMIT_AMT = '#{limitAmt}'
			,INVITE_LINK = '#{inviteLink}'
			,SUPPORT_AMT = '#{supportAmt}'
			,SUPPORT_FORWARD_YN = '#{supportFwYn}'
			,SUPPORT_YN = '#{supportYn}'
			,LIMIT_DAY_AMT = '#{limitDayAmt}'
			,LIMIT_DAY_CNT = '#{limitDayCnt}'
			,INTRO = '#{intro}'
		WHERE 
			GRP_ID = '#{grpId}'`

// SelectCreateGrpSeq 1단계 장부번호 생성
var SelectCreateGrpSeq = `SELECT
			CONCAT('B',IFNULL(LPAD(MAX(SUBSTRING(GRP_ID, -10)) + 1, 10, 0), '0000000001')) as GRP_ID
		FROM
			PRIV_GRP_INFO`

// SelectCompanyInfo 2단계 2-1 회사 정보 가져옴
var SelectCompanyInfo = `SELECT COMPANY_ID
	          					,COMPANY_NM
	          					,BUSID as bizNum
	          					,CEO_NM
	          					,ADDR,ADDR2
	          					,TEL
	          					,EMAIL
	          					,HOMEPAGE
	          					,IFNULL(LAT,'') AS LAT
	          					,IFNULL(LNG,'') AS LNG
	          					,IFNULL(INTRO,'') AS INTRO
						FROM b_company
						WHERE 
							company_id = '#{companyId}'`

// InsertCreateGroup 장부추가 2단계 2-2받아온 정보로 장부 생성
var InsertCreateGroup string = `INSERT INTO PRIV_GRP_INFO
		(
			GRP_ID,
			GRP_NM,
			if #{busid} != '' then BUSID,
			if #{ceoNm} != '' then CEO_NM,
			if #{addr} != '' then ADDR,
			if #{addr2} != '' then ADDR2,
			if #{tel} != '' then TEL,
			if #{email} != '' then EMAIL,
			if #{lat} != '' then LAT,
			if #{lng} != '' then LNG,
			if #{supportFwYn} != '' then SUPPORT_FORWARD_YN, 
			if #{paymentId} != '' then PAYMENT_ID,
			AUTH_STAT,
			if #{authComment} != '' then AUTH_COMMENT,
			REG_DATE,
			MOD_DATE,
			AUTH_DATE,
			if #{intro} != '' then INTRO,
			if #{maxPersonNum} != '' then MAX_PERSON_NUM,
			USE_YN,
			LIMIT_YN,
			LIMIT_AMT, 
			DETAIL_VIEW_YN,
			GRP_TYPE_CD,
			SUPPORT_AMT,
            MORE_ORDER_YN,
            SUPPORT_YN,
			LIMIT_DAY_AMT,
			LIMIT_DAY_CNT,
			if #{grpPayTy} != '' then GRP_PAY_TY
		)
		VALUES
		(
  			  '#{grpId}'
			, '#{grpNm}'
			, '#{bizNum}'
			, '#{ceoNm}'
			, '#{addr}'
			, '#{addr2}'
			, '#{tel}'
			, '#{email}'
			, '#{lat}'
			, '#{lng}'
			, '#{supportFwYn}'
			, '#{paymentId}'
			, '#{authStat}'
			, '#{authComment}'
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, NULL
			, '#{intro}'
			, '#{maxPersonNum}'
			, 'Y'
			, '#{supportExceedYn}'
			, '#{limitYn}'
			, '#{limitAmt}'
			, '#{detailViewYn}'
			, '#{grpTypeCd}'
			, '#{supportAmt}'
			, 'N'
            , IFNULL('#{supportYn}', 'N')
            , '#{limitDayAmt}'
            , '#{limitDayCnt}'
			, IFNULL('#{grpPayTy}', 1)
		)`

// SelectUserInfo 3단계 3-1 유저아이디로 유저 정보를 불러옴
var GetDefaultGrpBookMng = `SELECT 
			A.USER_ID
			, LOGIN_ID 
			, USER_NM
			, B.COMPANY_ID 
			, DEPT
			, TEL
		FROM  PRIV_USER_INFO AS A
		inner JOIN b_company_manager AS B ON A.USER_ID = B.USER_ID
		WHERE 
			b.company_id= "#{companyId}"
			AND b.AUTHOR_CD ='CM'
			LIMIT 1`

// InsertCreateGroupUser 3단계 3-2 장부 유저를 추가함
var InsertCreateGroupUser = `INSERT INTO PRIV_GRP_USER_INFO
		(
			GRP_ID
			, USER_ID
			, GRP_AUTH
			, AUTH_STAT
			, AUTH_DATE
			, JOIN_TY
			, DEPT_ID
			, REG_DATE
			, MOD_DATE
		)
		VALUES
		(
			'#{grpId}'
			, '#{userId}'
			, '#{grpAuth}'
			, '#{authStat}'
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, '#{joinTy}'
			, '#{deptId}'
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
		)ON DUPLICATE KEY UPDATE 
							GRP_AUTH = '#{grpAuth}'
										,AUTH_DATE=DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										,JOIN_TY =  '#{joinTy}'
										,AUTH_STAT= '#{authStat}' 
										,MOD_DATE= DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')`

// InsertCompanyBook 3단계 3-3 회사장부를 추가함
var InsertCompanyBook = `INSERT INTO b_company_book(
									company_id
									, BOOK_ID
									, REG_DATE
									)
						VALUES(
								'#{companyId}'
								, '#{grpId}'
								, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s' ))`

var SelectCompanyUser = `select DISTINCT 
										b.user_nm as userNm
										, b.user_id as userId
										, b.hp_no as hpNo
										, DATE_FORMAT(b.join_date, '%Y-%m-%d') as date
										, b.email as email
							from b_company_user as a 
							left join priv_user_info as b on a.user_id = b.user_id 
							left outer join b_company_manager as c on a.user_id = c.user_id and c.company_id = a.company_id
							where 1=1 
							and c.user_id is null
							and a.company_id = '#{companyId}'
							AND (('#{searchKey}' = 'userNm' AND B.USER_NM LIKE '%#{searchKeyword}%')
							OR ('#{searchKey}' = 'userHp' AND B.HP_NO LIKE '%#{searchKeyword}%'))
							order by b.join_date asc`

var InsertGrpManager = `INSERT INTO b_company_manager(
									company_id
									,user_id
									,reg_date
									,mod_date
									,author_cd
									,dept
									,class
									,TEL
									,email
						)
							VALUES(
									 '#{companyId}'
									, '#{userId}'
									, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									, 'BM'
									, null
									, null
									, '#{hpNo}'
									, '#{email}'
							)`

var SelectNotBizCompanyBook = `
SELECT b.grp_nm AS grpNm, (SELECT count(*) FROM priv_grp_user_info AS e WHERE a.book_id = e.GRP_ID) AS grpCnt, a.book_id AS grpId, e.user_nm as grpMaster
FROM b_company_book AS a
INNER JOIN priv_grp_info AS b ON a.book_id = b.grp_id
INNER JOIN priv_grp_user_info AS c ON a.book_id = c.grp_id AND c.GRP_AUTH = 0
LEFT OUTER JOIN b_company_manager AS d ON d.company_id = a.company_id
INNER JOIN priv_user_info AS e ON c.user_id = e.user_id 
WHERE 
d.company_id IS NULL
AND a.company_id != '#{companyId}'
AND b.grp_nm LIKE '%#{searchKeyword}%'
`

var SelectNotBizCompanyBookCnt = `
SELECT count(*) as total
FROM b_company_book AS a
INNER JOIN priv_grp_info AS b ON a.book_id = b.grp_id
INNER JOIN priv_grp_user_info AS c ON a.book_id = c.grp_id AND c.GRP_AUTH = 0
LEFT OUTER JOIN b_company_manager AS d ON d.company_id = a.company_id
INNER JOIN priv_user_info AS e ON c.user_id = e.user_id 
WHERE 
d.company_id IS NULL
AND a.company_id != '#{companyId}'
AND b.grp_nm LIKE '%#{searchKeyword}%'
`

var UpdateBookCompanyId = `
UPDATE b_company_book
		SET
			company_id = '#{companyId}'
		WHERE 
			BOOK_ID = '#{grpId}'`

var InsertCompanyGrpUser = `
`

var SelectDeptCodeFirst = `
SELECT min(code_id) as deptCode FROM b_company_code 
WHERE 
company_id = '#{companyId}'
`

var SelectCompanyGrpUser = `
SELECT a.user_id as userId 
		,a.grp_id as grpId
		,b.user_nm as userNm
		,b.hp_no as hpNo
		,b.use_yn as useYn 
FROM priv_grp_user_info as a
INNER JOIN priv_user_info as b on a.user_id = b.user_id
WHERE 
a.grp_id = '#{grpId}'
`

var SelectLastUId = `
select max(u_id)+1 as UId from b_company_user`

var InsertCompanyBookUser = `
INSERT INTO b_company_user 
(COMPANY_ID
, USER_ID
, BOOK_ID
, USER_NM
, HP_NO
, DEPT
, REG_DATE
, USE_YN)
VALUES
(
'#{companyId}'
,'#{userId}' 
,'#{grpId}' 
,'#{userNm}' 
,'#{hpNo}' 
,'#{dept}' 
, DATE_FORMAT(now(), '%Y%m%d%H%i%s') 
,'#{useYn}'
)
`

var SelectUserId = `
SELECT USER_ID 
FROM priv_user_info 
WHERE 
	USER_NM = '#{userNm}' 
AND HP_NO = '#{hpNo}'`

var SelectCheckGrpUser = `
SELECT A.USER_ID AS userId, IFNULL(B.GRP_ID, '') AS grpId 
FROM priv_user_info AS A 
         LEFT OUTER JOIN (SELECT * FROM priv_grp_user_info WHERE 
															GRP_ID = '#{grpId}') as B
                         ON A.USER_ID = B.USER_ID
WHERE 
      A.USER_NM = '#{userNm}' 
  AND A.HP_NO = '#{hpNo}' 
`
