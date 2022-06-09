package _struct

type WorkUser struct {
	Content []struct {
		CommuteSeq            int         `json:"commuteSeq"`
		DeptCode              string      `json:"deptCode"`
		DeptName              string      `json:"deptName"`
		EmpName               string      `json:"empName"`
		AccountID             string      `json:"accountId"`
		EmpID                 int         `json:"empId"`
		EmpStatus             interface{} `json:"empStatus"`
		DutyName              interface{} `json:"dutyName"`
		CommuteDay            string      `json:"commuteDay"`
		PlanStartTime         string      `json:"planStartTime"`
		PlanEndTime           string      `json:"planEndTime"`
		VacName               string      `json:"vacName"`
		VacTime               int         `json:"vacTime"`
		VacStartTime          interface{} `json:"vacStartTime"`
		WorkOnTime            interface{} `json:"workOnTime"`
		WorkOffTime           interface{} `json:"workOffTime"`
		WorkOnStatus          string      `json:"workOnStatus"`
		WorkOffStatus         string      `json:"workOffStatus"`
		CommuteUseStatus      string      `json:"commuteUseStatus"`
		WorkTimeMinute        int         `json:"workTimeMinute"`
		LateMinute            int         `json:"lateMinute"`
		ChangeDescriptionCode interface{} `json:"changeDescriptionCode"`
		ChangeDescription     interface{} `json:"changeDescription"`
		OnType                interface{} `json:"onType"`
		OffType               interface{} `json:"offType"`
		WorkOnIP              interface{} `json:"workOnIp"`
		WorkOnAddress         interface{} `json:"workOnAddress"`
		WorkOffIP             interface{} `json:"workOffIp"`
		WorkOffAddress        interface{} `json:"workOffAddress"`
		CanModify             bool        `json:"canModify"`
		Hyphen                string      `json:"hyphen"`
		TimeChangeSeq         int         `json:"timeChangeSeq"`
		TimeTypeCode          string      `json:"timeTypeCode"`
		WorkTypeYn            interface{} `json:"workTypeYn"`
		VacCodeID             interface{} `json:"vacCodeId"`
		OnLocationName        interface{} `json:"onLocationName"`
		OffLocationName       interface{} `json:"offLocationName"`
	} `json:"content"`
	Pageable struct {
		Sort struct {
			Unsorted bool `json:"unsorted"`
			Sorted   bool `json:"sorted"`
			Empty    bool `json:"empty"`
		} `json:"sort"`
		PageSize   int  `json:"pageSize"`
		PageNumber int  `json:"pageNumber"`
		Offset     int  `json:"offset"`
		Paged      bool `json:"paged"`
		Unpaged    bool `json:"unpaged"`
	} `json:"pageable"`
	Last          bool `json:"last"`
	TotalElements int  `json:"totalElements"`
	TotalPages    int  `json:"totalPages"`
	First         bool `json:"first"`
	Number        int  `json:"number"`
	Sort          struct {
		Unsorted bool `json:"unsorted"`
		Sorted   bool `json:"sorted"`
		Empty    bool `json:"empty"`
	} `json:"sort"`
	NumberOfElements int  `json:"numberOfElements"`
	Size             int  `json:"size"`
	Empty            bool `json:"empty"`
}
