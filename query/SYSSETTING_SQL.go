package query

var SelectSysUserInfo = `
SELECT user_id as userId
	 , user_nm as userNm 
FROM sys_user_info 
WHERE 
	user_id = '#{loginId}'
`

var SelectSysUserInfoList = `
SELECT 
	user_id as userId, 
	user_nm as userNm, 
	added_by as added, 
	CONN_ALLOW_YN as useYn, 
	added_date as addDate 
FROM sys_user_info 
WHERE 
	author_cd = 'SYS'
`

var InsertSysUser = `
INSERT INTO sys_user_info 
(USER_ID
, USER_NM
, USER_PASS
, MSG_LANG_CD
, CONN_ALLOW_YN
, INIT_PASS_YN
, RETRY_CNT
, USER_MENU_AUTHOR_YN
, AUTHOR_CD
, OTP_YN
, ADDED_BY
, ADDED_DATE)
VALUES
	('#{userId}'
	, '#{userNm}'
	, '#{userPw}'
	, 'ko_KR'
	, 'Y'
	, 'N'
	, 0
	, '#{menuYn}'
	, '#{auth}'
	, 'N'
	, '#{lUserId}'
	, DATE_FORMAT(NOW(), '%Y-%m-%d %H:%i:%s')
	)
`

var UpdateSysUser = `
UPDATE sys_user_info
SET
	USER_NM = '#{userNm}'
	, MODIFY_BY = '#{userId}'
	, MODIFY_DATE = DATE_FORMAT(NOW(), '%Y-%m-%d %H:%i:%s')
	, USER_PASS = IF('#{userPw}' = '',(SELECT a.USER_PASS 
							  FROM (SELECT USER_PASS
							  		FROM sys_user_info 
							  		WHERE
										USER_ID = '#{userId}') as a
							  		) 
							,'#{userPw}')
	, LAST_PASS_UPD_DATE = IF('#{userPw}' = '',(SELECT a.LAST_PASS_UPD_DATE 
							  FROM (SELECT LAST_PASS_UPD_DATE
							  		FROM sys_user_info 
							  		WHERE
										USER_ID = '#{userId}') as a
							  		) 
							,DATE_FORMAT(NOW(), '%Y-%m-%d %H:%i:%s'))
WHERE
	USER_ID = '#{userId}'
`

var SelectSysUserPwCheck = `
SELECT IF(count(*) = 1,'1','0') as state 
FROM sys_user_info 
WHERE 
	user_id = '#{userId}' 
AND user_pass = '#{userPw}' 
AND conn_allow_yn = 'Y'
`

var UpdateSysUserUseYn = `
UPDATE sys_user_info 
SET
	CONN_ALLOW_YN = '#{useYn}'
WHERE 
	user_id = '#{userId}'
`
