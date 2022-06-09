package _struct

type WorkTime struct {
	Result        string `json:"result"`
	Code          string `json:"code"`
	CommuteReport struct {
		CommuteSeq            int         `json:"commuteSeq"`
		DeptCode              string      `json:"deptCode"`
		DeptName              string      `json:"deptName"`
		EmpName               string      `json:"empName"`
		AccountID             string      `json:"accountId"`
		EmpID                 int         `json:"empId"`
		EmpStatus             string      `json:"empStatus"`
		DutyName              interface{} `json:"dutyName"`
		CommuteDay            string      `json:"commuteDay"`
		PlanStartTime         string      `json:"planStartTime"`
		PlanEndTime           string      `json:"planEndTime"`
		VacName               string      `json:"vacName"`
		VacTime               int         `json:"vacTime"`
		VacStartTime          interface{} `json:"vacStartTime"`
		WorkOnTime            string      `json:"workOnTime"`
		WorkOffTime           string      `json:"workOffTime"`
		WorkOnStatus          string      `json:"workOnStatus"`
		WorkOffStatus         string      `json:"workOffStatus"`
		CommuteUseStatus      interface{} `json:"commuteUseStatus"`
		WorkTimeMinute        int         `json:"workTimeMinute"`
		LateMinute            int         `json:"lateMinute"`
		ChangeDescriptionCode string      `json:"changeDescriptionCode"`
		ChangeDescription     string      `json:"changeDescription"`
		OnType                string      `json:"onType"`
		OffType               string      `json:"offType"`
		WorkOnIP              interface{} `json:"workOnIp"`
		WorkOnAddress         interface{} `json:"workOnAddress"`
		WorkOffIP             interface{} `json:"workOffIp"`
		WorkOffAddress        interface{} `json:"workOffAddress"`
		CanModify             bool        `json:"canModify"`
		Hyphen                string      `json:"hyphen"`
		TimeChangeSeq         int         `json:"timeChangeSeq"`
		TimeTypeCode          interface{} `json:"timeTypeCode"`
		WorkTypeYn            interface{} `json:"workTypeYn"`
		VacCodeID             int         `json:"vacCodeId"`
		OnLocationName        string      `json:"onLocationName"`
		OffLocationName       string      `json:"offLocationName"`
	} `json:"commuteReport"`
	Message string `json:"message"`
}
