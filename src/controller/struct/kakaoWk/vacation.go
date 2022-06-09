package _struct

type Vacation struct {
	Code         string `json:"code"`
	VacApplyList []struct {
		CreatedAt string `json:"createdAt"`
		CreatedBy int    `json:"createdBy"`
		UpdatedAt string `json:"updatedAt"`
		UpdatedBy int    `json:"updatedBy"`
		Seq       int    `json:"seq"`
		CmpID     int    `json:"cmpId"`
		Employee  struct {
			CreatedAt  string `json:"createdAt"`
			CreatedBy  int    `json:"createdBy"`
			UpdatedAt  string `json:"updatedAt"`
			UpdatedBy  int    `json:"updatedBy"`
			CmpID      int    `json:"cmpId"`
			EmpID      int    `json:"empId"`
			DeptID     int    `json:"deptId"`
			UserID     int    `json:"userId"`
			Department struct {
				CreatedAt          string      `json:"createdAt"`
				CreatedBy          int         `json:"createdBy"`
				UpdatedAt          string      `json:"updatedAt"`
				UpdatedBy          int         `json:"updatedBy"`
				CmpID              int         `json:"cmpId"`
				DeptID             int         `json:"deptId"`
				DeptCode           string      `json:"deptCode"`
				DeptName           string      `json:"deptName"`
				DeptOrd            int         `json:"deptOrd"`
				UpDeptID           int         `json:"upDeptId"`
				DeptPathID         string      `json:"deptPathId"`
				EffectiveStartDate interface{} `json:"effectiveStartDate"`
				EffectiveEndDate   interface{} `json:"effectiveEndDate"`
				UseYn              string      `json:"useYn"`
				Depth              int         `json:"depth"`
			} `json:"department"`
			Name          string      `json:"name"`
			AccountID     string      `json:"accountId"`
			WorkStartTime string      `json:"workStartTime"`
			WorkEndTime   string      `json:"workEndTime"`
			DutyCodeID    interface{} `json:"dutyCodeId"`
			TypeCodeID    int         `json:"typeCodeId"`
			TypeCode      string      `json:"typeCode"`
			EmployeeType  struct {
				CreatedAt      string      `json:"createdAt"`
				CreatedBy      int         `json:"createdBy"`
				UpdatedAt      interface{} `json:"updatedAt"`
				UpdatedBy      interface{} `json:"updatedBy"`
				EmployeeTypeID struct {
					TypeCode string `json:"typeCode"`
					CmpID    int    `json:"cmpId"`
				} `json:"employeeTypeId"`
				EmployeeTypeCode struct {
					CreatedAt        string      `json:"createdAt"`
					CreatedBy        int         `json:"createdBy"`
					UpdatedAt        interface{} `json:"updatedAt"`
					UpdatedBy        interface{} `json:"updatedBy"`
					CodeID           int         `json:"codeId"`
					SyncCodeID       interface{} `json:"syncCodeId"`
					CmpID            int         `json:"cmpId"`
					CodeTranslations []struct {
						CodeTranslationID struct {
							CmpID      int    `json:"cmpId"`
							CodeID     int    `json:"codeId"`
							LocaleCode string `json:"localeCode"`
						} `json:"codeTranslationId"`
						Name string `json:"name"`
					} `json:"codeTranslations"`
					CodeTranslation struct {
						CodeTranslationID struct {
							CmpID      int    `json:"cmpId"`
							CodeID     int    `json:"codeId"`
							LocaleCode string `json:"localeCode"`
						} `json:"codeTranslationId"`
						Name string `json:"name"`
					} `json:"codeTranslation"`
					CodeType     string `json:"codeType"`
					Code         string `json:"code"`
					CodeName     string `json:"codeName"`
					DisplayOrder int    `json:"displayOrder"`
					UseYn        string `json:"useYn"`
					SyncYn       string `json:"syncYn"`
				} `json:"employeeTypeCode"`
				CanDelYn                 string `json:"canDelYn"`
				DelYn                    string `json:"delYn"`
				DefaultOt                int    `json:"defaultOt"`
				HibernateLazyInitializer struct {
				} `json:"hibernateLazyInitializer"`
			} `json:"employeeType"`
			Email                 string      `json:"email"`
			MobileTelephoneNumber interface{} `json:"mobileTelephoneNumber"`
			EmpStatus             string      `json:"empStatus"`
			OfficeLocation        interface{} `json:"officeLocation"`
			GroupNo               int         `json:"groupNo"`
			ProfileImageURL       interface{} `json:"profileImageUrl"`
			GroupApplyDate        string      `json:"groupApplyDate"`
			GroupApplyEmpNo       interface{} `json:"groupApplyEmpNo"`
			SettingSeq            int         `json:"settingSeq"`
			JoinDate              string      `json:"joinDate"`
			RetirementDate        interface{} `json:"retirementDate"`
			LayoffDate            interface{} `json:"layoffDate"`
			ReinstatementDate     interface{} `json:"reinstatementDate"`
			LayOffReasonCodeID    interface{} `json:"layOffReasonCodeId"`
			LayoffFromDate        interface{} `json:"layoffFromDate"`
			LayoffToDate          interface{} `json:"layoffToDate"`
			LeaderIssuedDate      interface{} `json:"leaderIssuedDate"`
			LeaderReleaseDate     interface{} `json:"leaderReleaseDate"`
			TraineeEndDate        interface{} `json:"traineeEndDate"`
			UseServiceYn          string      `json:"useServiceYn"`
			LocaleCode            string      `json:"localeCode"`
			RenewalAt             string      `json:"renewalAt"`
			ConversationID        string      `json:"conversationId"`
			UpdateAccountID       interface{} `json:"updateAccountId"`
			UpdateName            interface{} `json:"updateName"`
			Delegate              interface{} `json:"delegate"`
			DelegateToMeList      interface{} `json:"delegateToMeList"`
			WorkType              interface{} `json:"workType"`
		} `json:"employee"`
		AgentEmployee interface{} `json:"agentEmployee"`
		AccountSeq    int         `json:"accountSeq"`
		ApprovalSeq   int         `json:"approvalSeq"`
		VacationCode  struct {
			CreatedAt        string      `json:"createdAt"`
			CreatedBy        int         `json:"createdBy"`
			UpdatedAt        interface{} `json:"updatedAt"`
			UpdatedBy        interface{} `json:"updatedBy"`
			CmpID            int         `json:"cmpId"`
			CodeID           int         `json:"codeId"`
			Depth            int         `json:"depth"`
			ParentCodeID     int         `json:"parentCodeId"`
			VacationTypeCode struct {
				CreatedAt        string      `json:"createdAt"`
				CreatedBy        int         `json:"createdBy"`
				UpdatedAt        interface{} `json:"updatedAt"`
				UpdatedBy        interface{} `json:"updatedBy"`
				CodeID           int         `json:"codeId"`
				SyncCodeID       interface{} `json:"syncCodeId"`
				CmpID            int         `json:"cmpId"`
				CodeTranslations []struct {
					CodeTranslationID struct {
						CmpID      int    `json:"cmpId"`
						CodeID     int    `json:"codeId"`
						LocaleCode string `json:"localeCode"`
					} `json:"codeTranslationId"`
					Name string `json:"name"`
				} `json:"codeTranslations"`
				CodeTranslation struct {
					CodeTranslationID struct {
						CmpID      int    `json:"cmpId"`
						CodeID     int    `json:"codeId"`
						LocaleCode string `json:"localeCode"`
					} `json:"codeTranslationId"`
					Name string `json:"name"`
				} `json:"codeTranslation"`
				CodeType     string `json:"codeType"`
				Code         string `json:"code"`
				CodeName     string `json:"codeName"`
				DisplayOrder int    `json:"displayOrder"`
				UseYn        string `json:"useYn"`
				SyncYn       string `json:"syncYn"`
			} `json:"vacationTypeCode"`
			VacationCodeTranslations []struct {
				VacationCodeTranslationID struct {
					CmpID      int    `json:"cmpId"`
					CodeID     int    `json:"codeId"`
					LocaleCode string `json:"localeCode"`
				} `json:"vacationCodeTranslationId"`
				Name string `json:"name"`
			} `json:"vacationCodeTranslations"`
			VacationCodeTranslation struct {
				VacationCodeTranslationID struct {
					CmpID      int    `json:"cmpId"`
					CodeID     int    `json:"codeId"`
					LocaleCode string `json:"localeCode"`
				} `json:"vacationCodeTranslationId"`
				Name string `json:"name"`
			} `json:"vacationCodeTranslation"`
			CodeName               string      `json:"codeName"`
			UseType                string      `json:"useType"`
			UseTerm                int         `json:"useTerm"`
			UseTime                int         `json:"useTime"`
			DeductYn               string      `json:"deductYn"`
			HolidayYn              string      `json:"holidayYn"`
			DivisionYn             string      `json:"divisionYn"`
			FileYn                 string      `json:"fileYn"`
			SalaryYn               string      `json:"salaryYn"`
			DescriptionYn          string      `json:"descriptionYn"`
			UseYn                  string      `json:"useYn"`
			DelYn                  string      `json:"delYn"`
			ApprovalYn             string      `json:"approvalYn"`
			VacationTypeCodeString string      `json:"vacationTypeCodeString"`
			ChangeTypeEnum         interface{} `json:"changeTypeEnum"`
			VacationFullName       string      `json:"vacationFullName"`
		} `json:"vacationCode"`
		StartDate          string      `json:"startDate"`
		EndDate            string      `json:"endDate"`
		StartTime          string      `json:"startTime"`
		EndTime            string      `json:"endTime"`
		UseMinute          int         `json:"useMinute"`
		Description        string      `json:"description"`
		FileID             interface{} `json:"fileId"`
		ApprovalAdminYn    string      `json:"approvalAdminYn"`
		VacationModifyType interface{} `json:"vacationModifyType"`
		WorkCloseYn        interface{} `json:"workCloseYn"`
		EarlyEnd           bool        `json:"earlyEnd"`
	} `json:"vacApplyList"`
	Message string `json:"message"`
}
