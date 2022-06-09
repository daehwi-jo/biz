package _struct

type BookUser struct {
	CompanyId string
	GrpId     string
	UserList  []userList
}

type userList struct {
	UserNm string
	HpNo   string
	DeptId string
	DeptNm string
}
