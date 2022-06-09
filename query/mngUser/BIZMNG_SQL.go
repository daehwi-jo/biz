package mngUser

var SelectCompanyInfo = `SELECT COMPANY_ID
	          					,COMPANY_NM
	          					,BUSID
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

var SelectManagerInfo = `SELECT IFNULL(B.USER_NM,'') AS USER_NM
								,IFNULL(A.TEL,'') AS TEL
								,IFNULL(A.DEPT,'') AS DEPT
								,IFNULL(A.CLASS,'') AS COURSE
								,IFNULL(A.EMAIL,'') AS EMAIL
								,B.USER_ID
						FROM b_company_manager as a
						inner join priv_user_info as b on a.user_id = b.user_id
						where 
							company_id = '#{companyId}' and author_cd = 'CM'`

var UpdateCompanyInfo = `UPDATE b_company
		SET
			company_nm = '#{companyNm}'
			,CEO_NM = '#{ceoNm}'
			,TEL = '#{tel}'
			,ADDR = '#{addr}'
			,ADDR2 = '#{addr2}'
		    ,USE_YN = '#{useYn}'
			,homepage = '#{homepage}'
			,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
		WHERE 
			company_id = '#{companyId}'`

var UpdateManagerInfo = `UPDATE b_company_manager
		SET
			MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			,DEPT = '#{dept}'
			,CLASS = '#{class}'
			,TEL = '#{tel}'
			,EMAIL = '#{email}'
		WHERE 
			company_id = '#{companyId}'
			AND USER_ID = '#{userId}'
`
