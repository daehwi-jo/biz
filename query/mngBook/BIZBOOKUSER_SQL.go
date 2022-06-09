package mngBook

var SelectGrpCode = `SELECT A.GRP_ID as value
			  				,A.GRP_NM as label
							,A.GRP_TYPE_CD
						FROM priv_grp_info as a
						inner join b_company_book as b on a.grp_id = b.book_id
						where 
							a.use_yn='Y' AND company_id = '#{companyId}'
						AND B.BOOK_ID IN (
											SELECT  GRP_ID
											FROM priv_grp_user_info as aa
											WHERE 
												aa.GRP_AUTH=0 and aa.user_id = '#{searchUid}' )`

var SelectGrpUserMng = `Select  A.U_ID
      			,B.GRP_ID 
			    ,B.GRP_NM 
				,A.USER_NM 
				,(SELECT LOGIN_ID FROM priv_user_info AS AA WHERE A.USER_NM = AA.USER_NM AND A.USER_ID = AA.USER_ID  AND USER_TY=0 LIMIT 1) AS  LOGIN_ID
				,C.AUTH_STAT 
				,IFNULL(C.SUPPORT_BALANCE,0) AS   SUPPORT_BALANCE
				,A.HP_NO 
				,(SELECT USER_ID FROM priv_user_info AS AA WHERE A.USER_NM = AA.USER_NM AND A.USER_ID = AA.USER_ID  AND USER_TY=0 LIMIT 1) AS  USER_ID
				,'' AS COURSE_ID
				,'' AS COURSE_NM
				,IFNULL(MON_AMT,0) AS MON_AMT
				,IFNULL(MON_CNT,0) AS MON_CNT
				,DATE_FORMAT(C.REG_DATE,'%Y-%m-%d') as REG_DATE
				,DATE_FORMAT(A.USER_BIRTH,'%Y-%m-%d') as USER_BIRTH
				,A.LUNAR_BIRTH_YN
			FROM b_company_user AS A
			LEFT OUTER JOIN priv_grp_info AS B ON A.book_id = B.grp_id AND B.USE_YN='Y'
			LEFT OUTER JOIN priv_grp_user_info as C on A.book_id = C.GRP_ID AND A.USER_ID = C.USER_ID
			LEFT OUTER JOIN (  SELECT IFNULL(COUNT(*),0) AS MON_CNT,
											  IFNULL(SUM(ORDER_AMT*ORDER_QTY),0) AS MON_AMT
											  ,BB.USER_ID,AA.GRP_ID
									FROM dar_order_info AS AA 
									INNER JOIN dar_order_detail AS BB ON AA.ORDER_NO = BB.ORDER_NO
									WHERE ORDER_STAT='20' AND LEFT(AA.ORDER_DATE,6)=date_format(now(),'%Y%m') 
									GROUP BY BB.USER_ID,AA.GRP_ID
									) AS D ON C.USER_ID = D.USER_ID AND B.GRP_ID = D.GRP_ID
			WHERE 
				A.COMPANY_ID = '#{companyId}' 
				AND A.USE_YN = 'Y'
				AND B.GRP_ID = '#{searchGrpId}'
				
				AND (('#{searchKey}' = 'userNm' AND A.USER_NM LIKE '%#{searchKeyword}%') 
				OR ('#{searchKey}' = 'userHp' AND A.HP_NO LIKE '%#{searchKeyword}%'))
				`

func OrderBySort(params map[string]string) string {
	switch params["sortKey"] {
	case "userNm":
		return `ORDER BY A.USER_NM ASC`
	case "authState":
		return `ORDER BY C.AUTH_STAT desc, A.USER_NM ASC`
	case "support":
		return `ORDER BY C.SUPPORT_BALANCE DESC`
	case "regDate":
		return `ORDER BY C.REG_DATE DESC`
	case "monAmt":
		return `ORDER BY D.MON_AMT DESC`
	case "":
		return `ORDER BY C.AUTH_STAT desc, A.USER_NM ASC`
	}
	return ""
}

var SelectGrpUserMngCnt = `SELECT COUNT(*) as TOTAL_COUNT
							FROM b_company_user AS A
							LEFT OUTER JOIN priv_grp_info AS B ON A.book_id = B.grp_id AND B.USE_YN='Y'
							LEFT OUTER JOIN priv_grp_user_info as C on A.book_id = C.GRP_ID AND A.USER_ID = C.USER_ID
							WHERE 
								A.COMPANY_ID = '#{companyId}' 
								AND A.USE_YN= 'Y'
								AND B.GRP_ID = '#{searchGrpId}'
								AND (('#{searchKey}' = 'userNm' AND A.USER_NM LIKE '%#{searchKeyword}%') 
								OR ('#{searchKey}' = 'userHp' AND A.HP_NO LIKE '%#{searchKeyword}%'))`

var SelectGrpCompanyUserInfo = `SELECT A.U_ID
				,A.USER_NM
				,A.HP_NO
				,C.GRP_NM
				,IFNULL(B.SUPPORT_BALANCE,0) AS SUPPORT_AMT
				,D.LOGIN_ID
				,A.DEPT
				,A.USER_ID
				,C.GRP_ID
				,A.USER_BIRTH
				,A.LUNAR_BIRTH_YN
		FROM b_company_user AS a
		LEFT OUTER JOIN priv_grp_user_info as B on A.book_id = B.GRP_ID AND A.USER_ID = B.USER_ID AND B.AUTH_STAT ='1'
		LEFT OUTER JOIN priv_grp_info AS C ON B.grp_id = C.grp_id
		LEFT OUTER JOIN priv_user_info as D on B.USER_ID = D.USER_ID
		WHERE 
			A.U_ID = '#{uId}'
`

var UpdateGroupUserInfo = `UPDATE PRIV_GRP_USER_INFO
		SET
			AUTH_STAT = '#{authStat}'
			,AUTH_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			,LEAVE_TY = '#{leaveTy}'
			,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
		WHERE 
			GRP_ID = '#{grpId}'
			AND USER_ID = '#{userId}' `

//빠진데이터 SET항목에 들어가야함
//GRP_AUTH = '#{grpAuth}'
//,JOIN_TY = '#{joinTy}'
//,DEPT_ID = '#{deptId}'
//,SUPPORT_BALANCE = '#{supportBalance}'

var UpdateCompanyUserDel = `UPDATE b_company_user SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								,USE_YN = 'N'
		WHERE 
			COMPANY_ID = '#{companyId}' 
			and U_ID = '#{uid}'`

var UpdateCompanyUserInfo = `
UPDATE b_company_user
SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
  , USER_NM  = '#{userNm}'
  , HP_NO    = '#{hpNo}'
  , DEPT     = '#{deptId}'
  , USER_ID  = '#{userId}'
  , USER_BIRTH = '#{userBirth}'
  , LUNAR_BIRTH_YN = '#{lunarBirthYn}'
WHERE 
	COMPANY_ID = '#{companyId}'
  and U_ID = '#{uid}'
`

var UpdateCompanyUserSupportAmt = `UPDATE priv_grp_user_info SET MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'), 
											SUPPORT_BALANCE = #{supportAmt}
									WHERE 
										GRP_ID = '#{grpId}' 
									and USER_ID = '#{userId}'`

var InsertConnectBook = `INSERT INTO PRIV_GRP_USER_INFO
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
			, '#{uId}'
			, '#{grpAuth}'
			, '#{authStat}'
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, '#{joinTy}'
			, ''
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
		)ON DUPLICATE KEY UPDATE  GRP_AUTH = '#{grpAuth}'
										,AUTH_DATE=DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										,JOIN_TY =  '#{joinTy}'
										,AUTH_STAT= '#{authStat}' 
										,MOD_DATE= DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')`

//Values안에 들어가야할항목
//deptId 항목 ''로 삽입함 - 받는 파라미터 없음

var UpdateConnectBook = `UPDATE PRIV_GRP_USER_INFO
		SET
			AUTH_STAT = '#{authStat}',
			AUTH_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'),
			GRP_AUTH = '#{grpAuth}',
			JOIN_TY = '#{joinTy}',
			if #{leaveTy} != '' then LEAVE_TY,
			MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
		WHERE 
			GRP_ID = '#{grpId}'
			AND USER_ID = '#{userId}' `

//빠진데이터 SET항목에 들어가야함
//,DEPT_ID = '#{deptId}'
//,SUPPORT_BALANCE = '#{supportBalance}'

var InsertCompanyUser = `INSERT INTO b_company_user(
									COMPANY_ID
									, USER_ID
									, BOOK_ID
									, USER_NM
									, HP_NO
									, DEPT
									, REG_DATE
									, USE_YN)
							VALUES (
									'#{companyId}'
									, IF('#{userId}'='<nil>',NULL,'#{userId}')
									, '#{grpId}'
									, '#{userNm}'
									, '#{hpNo}'
									, IF('#{deptId}'='<nil>',NULL,'#{deptId}')
									, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									, 'Y')`

var SelectGrpBookInfo = `select A.GRP_NM
								, A.GRP_TYPE_CD
								, CASE  WHEN B.CAT_CD IS NULL THEN 'N' ELSE 'Y' END AS CAT_YN
								, A.DETAIL_VIEW_YN
								, A.LIMIT_YN
								, IFNULL(A.LIMIT_AMT,0) AS LIMIT_AMT
								, IFNULL(A.LIMIT_DAY_AMT ,0) AS LIMIT_DAY_AMT
								, A.INVITE_LINK
								, IFNULL(A.SUPPORT_AMT ,0) AS SUPPORT_AMT
								, A.support_yn
								, A.intro
						FROM priv_grp_info AS A
						LEFT OUTER JOIN priv_catering_grp_info AS B ON A.GRP_ID = B.GRP_ID
						where 
							A.grp_id = '#{searchGrpId}'`

var UpdateGrpBookInvite = `UPDATE PRIV_GRP_INFO
								SET 
									MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'),
									INVITE_LINK = '#{inviteLink}',
							WHERE 
								GRP_ID = '#{grpId}'`

var SelectDeptCode = `SELECT CODE_ID AS VALUE
			  				,CODE_INFO  AS LABEL
						FROM b_company_code 
						WHERE 
							company_id = '#{companyId}' 
							AND CODE_TY = 0`

var SelectCheckAuthState = `select IFNULL(MAX(AUTH_STAT),NULL) as authState
								from PRIV_GRP_USER_INFO 
								WHERE 
									GRP_ID = '#{grpId}' 
									and USER_ID = '#{uId}'`

var SearchUserId = `select user_id as userId from priv_user_info
										where 1=1
										and user_nm = '#{userNm}'
										and hp_no = '#{hpNo}'`

var SearchCompanyUser = `select U_ID from b_company_user 
									where 
									book_id = '#{grpId}' 
									and USER_ID = '#{userId}'`

var UpdateCompanyUser = `UPDATE b_company_user 
						SET USE_YN = 'Y'
							,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
						where 
						user_id = '#{userId}' 
						and book_id = '#{grpId}'`

var SelectUserInfo = `select USER_NM as userNm, HP_NO as hpNo from priv_user_info where 
						user_id = '#{userId}' `

var SelectGrpUserInfo = `select Count(*) as state
							from priv_grp_user_info 
							where 
									user_id = '#{userId}' 
								and grp_id = '#{grpId}'`

var SelectCompanyDeptCode = `
SELECT code_id as codeId FROM b_company_code 
WHERE
company_id = '#{companyId}'
`

var InsertGrpUser = `
INSERT INTO priv_grp_user_info
	(grp_id
	,user_id
	,grp_auth
	,auth_stat
	,auth_date
	,join_ty
	,reg_date)
VALUES (
	'#{grpId}'
	,'#{userId}'
	,'1'
	,'1'
	,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
	,'0'
	,'#{deptId}'
	,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'))
`

var UpdateSupportBalance = `
UPDATE PRIV_GRP_USER_INFO AS A
INNER JOIN PRIV_GRP_INFO AS B ON A.GRP_ID = B.GRP_ID AND B.SUPPORT_YN = 'Y'
SET A.SUPPORT_BALANCE = B.FIRST_SUPPORT_AMT
WHERE 
A.GRP_ID = '#{grpId}'
AND A.USER_ID = '#{userId}'
`
