package query

func SearchMemberSql(key string) string {
	query := `select 
					USER_NM as uNm, 
					USER_ID as uId, 
					HP_NO as HpNo 
				from priv_user_info 
				where
				`

	if key == "userNm" {
		query += `USER_NM LIKE '%#{searchKeyword}%'`
	} else {
		query += `HP_NO LIKE '%#{searchKeyword}%'`
	}
	return query
}

func SearchBookSql(key string) string {
	query := `select 
					a.USER_NM as uNm , 
					a.USER_ID as uId , 
					b.GRP_ID as grpId, 
					c.GRP_NM as grpNm  
				from priv_user_info as a 
				inner join priv_grp_user_info as b on a.USER_ID = b.USER_ID and b.GRP_AUTH='0'
				inner join priv_grp_info as c on c.GRP_ID = b.GRP_ID
				where 
				`

	if key == "userNm" {
		query += `a.USER_NM LIKE '%#{searchKeyword}%'`
	} else {
		query += `a.HP_NO LIKE '%#{searchKeyword}%'`
	}
	return query
}

var SearchStoreSql = `select 
					REST_ID as restId, 
					REST_NM as restNm from priv_rest_info 
				where
				REST_NM LIKE '%#{searchKeyword}%'

				`
