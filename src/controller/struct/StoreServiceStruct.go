package _struct

type StoreServiceObj struct {
	RestId       string
	FacilityList []facilityList
}

type facilityList struct {
	ServiceId   string
	ServiceInfo string
	UseYn       string
	NoticeYn    string
	ServiceNm   string
}
