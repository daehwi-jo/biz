package mngAdmin

// SelectCsList CsListLoad
// notice : select
var SelectCsList = `SELECT A.CONTENTS
							,substr(A.CONTENTS,1,30) AS shortContents
							,DATE_FORMAT(IFNULL(A.REG_DATE,''),'%Y-%m-%d') AS REG_DATE
							,A.SEQ
					FROM b_cs AS A
					WHERE 1=1
					AND A.KEY_ID = '#{searchKeyId}'
					AND A.type = '#{searchType}'
					order by a.reg_date desc`

// InsertCsContent CsListAdd
var InsertCsContent = `INSERT INTO b_cs
						(
								KEY_ID
								,type
								,contents
								,REG_DATE
						)
						VALUES
						(
								'#{keyId}'
								,'#{type}'
								,'#{contents}'
								,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
						)`

var SelectTaskListCnt string = `SELECT 	COUNT(*) AS TOTAL_COUNT
											FROM b_dailytask
											WHERE 
											DATE_FORMAT(REG_DATE, '%Y-%m-%d') >=  '#{startDate}'
											AND DATE_FORMAT(REG_DATE, '%Y-%m-%d') <=  '#{endDate}'
											ORDER BY REG_DATE DESC
											`

var SelectTaskList string = `SELECT TASK_NAME
													,TASK_SUCC_YN
													,CYCLE
													,DATE_FORMAT(REG_DATE, '%Y-%m-%d') as REG_DATE
											FROM b_dailytask
											WHERE 
											DATE_FORMAT(REG_DATE, '%Y-%m-%d') >=  '#{startDate}'
											AND DATE_FORMAT(REG_DATE, '%Y-%m-%d') <=  '#{endDate}'
											ORDER BY REG_DATE DESC
											`
