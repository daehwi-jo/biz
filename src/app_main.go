package biz

import (
	Daily "biz-web/src/controller"
	"biz-web/src/controller/api/mongo"
	apiPush "biz-web/src/controller/api/push"
	"biz-web/src/controller/mains"
	"biz-web/src/controller/mngAdmin"
	"biz-web/src/controller/mngBook"
	"biz-web/src/controller/mngOrder"
	"biz-web/src/controller/mngUser"
	"biz-web/src/controller/payment"
	"biz-web/src/controller/tpays"
	//공통 api
	"biz-web/src/controller/cls"

	"github.com/labstack/echo/v4"
)

var lprintf func(int, string, ...interface{}) = cls.Lprintf

func SvcSetting(e *echo.Echo, fname string) *echo.Echo {
	lprintf(4, "[INFO] sql start \n")

	api(e)

	cls.SetNotLoginUrl("/")
	lprintf(4, "[INFO] page start \n")

	return e
}

func api(e *echo.Echo) {
	//SELECT = GET	: GET~~~~ 메소드
	//UPDATA = PUT	: SET~~~~ 메소드
	//INSERT = POST	:
	//DELETE = DEL ****데이터는 거의 지우지 않음

	//CS 조회 및 저장
	//	searchType => 	M : 회원				searchKeyId =>	U~~~~~ : 회원
	//					G : 장부								B,G~~~~~ : 장부
	//					S : 가맹점							S~~~~~ : 스토어

	api := e.Group("/api") //api 그룹생성

	api.GET("/biz/health", mains.BaseUrl)
	api.GET("/PushTest", apiPush.PushTest) // 푸쉬 테스트

	api.POST("/login", mains.LoginOk) // 로그인 api

	api.GET("/SysHomeData", mains.SysHomeData)                 // system 로그인 시 회사 고정 데이터
	api.GET("/SysHomeContentsData", mains.SysHomeContentsData) // system 로그인 시 랭크 및 그래프 데이터

	api.GET("/HomeData", mains.HomeData)                 // 일반 로그인 시 회사 고정 데이터
	api.GET("/HomeContentsData", mains.HomeContentsData) // 일반 로그인 시 회사 랭크 및 그래프 데이터

	api.GET("/GetCompanyList", mngAdmin.GetCompanyIdList) // system 로그인 시 회사 목록 불러오기

	api.GET("/GetSysUserSettingMainData", mngAdmin.GetSysUserSettingMainData) //시스템 사용자 메인 데이터
	api.POST("/AddSysUser", mngAdmin.AddSysUser)                              //시스템 사용자 추가
	api.PUT("/ModSysUser", mngAdmin.ModSysUser)                               //시스템 사용자 수정
	api.PUT("/ModSysUserAuth", mngAdmin.ModSysUserAuth)                       //시스템 사용자 접속 설정 변경
	api.POST("/CheckSysUserPw", mngAdmin.CheckSysUserPw)                      //시스템 사용자 패스워드 체크

	//정산
	api.POST("/payment/GetCombineStoreList", payment.GetCombineStoreList)                //통합 정산 대상
	api.POST("/payment/GetCombineStoreData", payment.GetCombineStoreData)                //통합 정산 대상상세
	api.POST("/payment/GetCombineStoreData_wincube", payment.GetCombineStoreDataWincube) //통합 정산 대상상세 - 윈큐브
	api.POST("/payment/SetCombinePaidMake", payment.SetCombinePaidMake)                  //통합 정산 데이터 생성
	api.POST("/payment/GetCombineSubStoreList", payment.GetCombineSubStoreList)          //통합 정산 가맹점 불러오기
	api.POST("/payment/SetGiftconUpdate", payment.SetGiftconUpdate)                      //기프티콘 사용현황 업데이트

	api.POST("/payment/GetPaidMngList", payment.GetPaidMngList) //정산대상
	api.POST("/payment/SetStoreReg", payment.SetStoreReg)       //서브몰 등록

	api.POST("/payment/GetPaidList", payment.GetPaidList)   //지급요청대상
	api.POST("/payment/GetPaidExcel", payment.GetPaidExcel) //지급요청대상 엑셀

	api.POST("/payment/GetPaidIngList", payment.GetPaidIngList)   //지급처리중 데이터
	api.POST("/payment/GetPaidIngExcel", payment.GetPaidIngExcel) //지급처리중 엑셀데이터
	api.POST("/payment/GetPaidOkList", payment.GetPaidOkList)     //지급완료
	api.POST("/payment/GetPaidOkExcel", payment.GetPaidOkExcel)   //지급완료 엑셀

	api.POST("/payment/GetPaidInfo", payment.GetPaidInfo)         //지급상세
	api.POST("/payment/PaidReq", payment.SetPaidReq)              //지급요청
	api.POST("/payment/PaidOk", payment.SetPaidOk)                //지급완료 처리
	api.PUT("/payment/settlmntDtCh", payment.SetSettlmntDtChange) //지급 요청일 일괄 수정
	api.PUT("/payment/SetPaidInfo", payment.SetPaidInfo)          //지급상세 수정
	api.POST("/payment/SetMakePaid", payment.SetMakePaid)         //지급요청대상 생성
	api.POST("/payment/SetMakeFee", payment.SetMakeFee)           //수수료 생성

	api.POST("/payment/GetAccount", payment.GetAccount) //자금현황조회

	api.GET("/payment/manualTpay", payment.ManualTpay) //수수료 생성

	//Admin
	api.GET("/BizUserMng", mngAdmin.GetUserList)            //회원관리 메인
	api.GET("/UserInfo", mngAdmin.GetUserInfo)              //회원정보 불러오기
	api.PUT("/ResetPassword", mngAdmin.ResetPassword)       //비밀번호 초기화
	api.PUT("/UserInfoCommit", mngAdmin.ModifyUserInfo)     //회원수정 저장
	api.GET("/UserInfoCsMng", mngAdmin.GetBizUserCsInfo)    //유저 CS 메인
	api.GET("/CheckGifticonUsed", mngAdmin.CheckGifticon)   //기프티콘 확인
	api.GET("/ExtendGifticonUsed", mngAdmin.ExtendGifticon) //기프티콘 갱신

	api.GET("/BizBookMng", mngAdmin.GetBizBookList) //장부리스트 불러오기
	api.GET("/BookData", mngAdmin.GetBookData)
	api.PUT("/ModBookData", mngAdmin.ModBookData)
	api.GET("/BookInfoCsMng", mngAdmin.GetBizBookCsInfo)   //장부 CS 메인
	api.GET("/SearchStoreList", mngAdmin.GetBookStoreList) //가맹점목록 불러오기

	api.GET("/BizCompanyMng", mngAdmin.GetBizCompanyList)      //업체리스트 불러오기
	api.GET("/CompanyInfoCsMng", mngAdmin.GetBIZCompanyCsInfo) //기업 CS 메인
	api.POST("/AddCompany", mngAdmin.AddCompany)               //기업 추가
	api.GET("/ManagerList", mngAdmin.GetCompanyManager)        //메니저 리스트 불러오기
	api.PUT("/ModManagerAuthor", mngAdmin.ModCompanyManager)   //메니저 권한 CM수정
	api.GET("/CompanyUserList", mngAdmin.GetCompanyUsers)      //회사 유저 목록
	api.POST("/AddManager", mngAdmin.AddCompanyManager)        //회사 매니저 추가

	api.GET("/BizStoreMng", mngAdmin.GetBizStoreList)   //가맹점리스트 불러오기
	api.GET("/StoreInfoCs", mngAdmin.GetBizStoreCsInfo) //가맹점 CS 메인
	api.PUT("/ModStoreInfo", mngAdmin.ModStoreInfo)     //가맹점 상세정보 수정
	api.POST("/SetStoreImg", mngAdmin.SetStoreImg)      //가맹점 이미지 추가

	api.GET("/GetLinkBookList", mngAdmin.GetStoreLinkBookList)       //가맹점 연결장부 목록 불러오기
	api.GET("/SearchLinkStoreBook", mngAdmin.GetStoreLinkBookSearch) //가맹점에 연결할 장부 검색

	api.GET("/GetLinkBookInfo", mngAdmin.GetStoreLinkBookInfo) //가맹점 연결장부 상세정보
	api.POST("/ModBookLink", mngAdmin.ModStoreLinkBook)        //가맹점 장부 연결

	api.GET("/GetCategoryMenuList", mngAdmin.GetStoreCategoryMenu) //메뉴,카테고리 불러오기
	api.GET("/GetCategoryList", mngAdmin.GetStoreCategory)         //카테고리 불러오기
	api.POST("/AddCategory", mngAdmin.AddStoreCategory)            //카테고리 추가
	api.PUT("/ModCategory", mngAdmin.ModStoreCategory)             //카테고리 수정

	api.GET("/GetMenuList", mngAdmin.GetStoreMenu) //메뉴 불러오기
	api.POST("/AddMenu", mngAdmin.AddStoreMenu)    //메뉴 추가
	api.PUT("/ModMenu", mngAdmin.ModStoreMenu)     //메뉴 수정
	api.POST("/MenuImg", mngAdmin.AddMenuImage)    //메뉴 이미지 업로드
	api.GET("/GetBankList", mngAdmin.GetBankList)  //은행목록 불러오기

	api.GET("/GetFacility", mngAdmin.GetFacilityList)    //부대시설 리스트
	api.PUT("/UpdateFacility", mngAdmin.ModFacilityList) //부대시설 수정
	api.POST("/AddFacility", mngAdmin.AddFacility)       //부대시설 추가

	api.POST("/GetUnpaidList", mngAdmin.GetUnPaidList)       // 매장 정산 대상 조회
	api.POST("/SetpaidOk", mngAdmin.SetPaidOk)               // 매장 정산
	api.DELETE("/SetPaidCancel", mngAdmin.SetPaidCancel)     // 매장 정산 취소
	api.POST("/GetPaidOkLIST", mngAdmin.GetStorePaymentList) // 매장 결제 조회

	api.GET("/CheckAccount", mngAdmin.AcctNameSearch) //계좌 실명 조회

	api.GET("/ChargeList", mngAdmin.GetChargeList)     //선불충전관리 리스트
	api.PUT("/ModChargeList", mngAdmin.ModChargeList)  //선불충전관리 리스트 수정
	api.POST("/AddChargeItem", mngAdmin.AddChargeItem) //선불충전관리 아이템 추가

	api.GET("/GetBizCsList", mngAdmin.GetBizCsList)        //CS목록 불러오기
	api.POST("/AddBizCsContent", mngAdmin.AddBizCsContent) //CS목록 추가하기

	api.POST("/GetPaymentList", payment.GetPaymentList) //결제관리 리스트

	api.POST("/GetOrderList", mngAdmin.GetOrderList)           //주문 리스트
	api.POST("/GetOrderListExcel", mngAdmin.GetOrderListExcel) //주문 리스트 엑셀 다운로드
	api.POST("/GetOrderInfo", mngAdmin.GetOrderInfo)           //주문 상세
	api.POST("/SetOrderCancel", mngAdmin.SetOrderCancel)       //주문 취소

	api.POST("/GetTaskList", mngAdmin.GetTaskList) //주문 리스트

	//사용자관리
	api.GET("/ManagerInfoMng", mngUser.GetManagerInfoMng)   //관리자 정보관리 메인
	api.PUT("/MngInfoUpdateCo", mngUser.ModifyCompanyInfo)  //관리자 정보관리 회사정보 수정
	api.PUT("/MngInfoUpdateMng", mngUser.ModifyManagerInfo) //관리자 정보관리 관리자정보 수정

	//장부관리
	api.GET("/BookMngList", mngBook.GetBookList)                //장부리스트 메인
	api.POST("/AddGrpBook", mngBook.CreateCoBookGrp)            //장부추가
	api.PUT("/DelGrpBook", mngBook.DeleteGrp)                   //장부삭제
	api.GET("/BookInfo", mngBook.GetBookInfo)                   //장부 수정 - 불러오기
	api.PUT("/BookInfoUpdate", mngBook.ModifyBookInfo)          //장부 수정 - 저장
	api.GET("/BookUserChangeList", mngBook.GetGrpUserList)      //장부리스트 - 권한 전달 (장부관리자 변경)
	api.PUT("/BookMngChange", mngBook.UpdateGrpUser)            //장부리스트 - 장부관리자 변경 - 변경
	api.GET("/BookCompanyUserList", mngBook.GetCompanyUserList) //장부리스트 - 회사장부유저 불러오기
	api.POST("/AddBookManager", mngBook.AddBookManager)         //장부리스트 - 관리자 추가
	api.GET("/UnLinkBookList", mngBook.GetCompanyBookList)      //기업 장부 리스트 불러오기(검색)
	api.PUT("/ModAddBook", mngBook.ModifyBookCompanyId)         //기존 장부 연결

	api.GET("/BookUserMng", mngBook.GetBookUserList)                 //장부사용자관리 메인자
	api.GET("/BookUserInfo", mngBook.GetBookUserInfo)                //사용자 정보  - 불러오기
	api.GET("/GrpBookInfo", mngBook.GetGrpBookInfo)                  //장부 정보 - 불러오기
	api.PUT("/BookUserDel", mngBook.UpdateBookUserDel)               //사용자 정보 - 직원삭제
	api.PUT("/BookUserDisconnect", mngBook.UpdateBookUserDisconnect) //사용자 정보 - 연결해제
	api.PUT("/ModifyBookUserInfo", mngBook.ModifyBookUserInfo)       //사용자 정보 - 수정
	api.POST("/ModifyConnectUser", mngBook.InsertBookUserConnect)    //장부 연결
	api.POST("/AddUserInBook", mngBook.InsertBookUser)               //사용자 추가 - 장부사용자 추가 - 등록
	api.POST("/ParsingExcel", mngBook.ParsingExcel)                  //장부 사용자 추가 - excel 파싱
	api.PUT("/CreateInviteLink", mngBook.CreateInviteURL)            //초대링크 생성
	api.GET("/GetDeptCode", mngBook.GetDeptCode)                     //부서코드

	//주문관리
	api.GET("/GetCompanyBookList", mngOrder.GetBookList)         //회사 장부 목록
	api.GET("/OrderHistoryLookUp", mngOrder.GetOrderHistory)     //주문내역조회 메인
	api.GET("/OrderHistoryExcel", mngOrder.GetOrderHistoryExcel) //주문내역조회 엑셀 다운로드
	api.POST("/CalculateBookMng", mngOrder.GetCalculateBookList) //장부정산관리 메인
	api.GET("/PaymentHistoryLookUp", mngOrder.GetPaymentHistory) //결제내역조회 메인
	api.POST("/StoreInfoLookUp", mngOrder.GetStoreInfoList)      //스토어조회 메인

	api.POST("/GetStoreChargeInfo", mngOrder.GetStoreChargeInfo) //스토어 충전하기 정보
	api.POST("/GetStoreUnpaidList", mngOrder.GetStoreUnpaidInfo) //스토어 미정산 정보

	//검색
	api.GET("/Search", mngAdmin.GetSysUserSearch) //검색
	api.POST("/Mongo", mongo.SendMongoDB)         // mongodb query 검색

	admin := api.Group("/admin")

	admin.GET("/partnerMember", mngAdmin.GetBizPartnerMemberList)               //멤버십 리스트 불러오기
	admin.GET("/partnerMemberDate", mngAdmin.GetBizPartnerMemberDate)           //멤버십 날짜 불러오기
	admin.PUT("/partnerMemberDate", mngAdmin.ModBizPartnerMemberDate)           //멤버십 날짜 수정하기
	admin.GET("/partnerMemberInfo", mngAdmin.GetBizPartnerMemberInfoData)       //멤버십 상세정보 불러오기 (가입 정보 및 기간 정보)
	admin.GET("/partnerMemberCollect", mngAdmin.GetBizPartnerMemberCollectList) //멤버십 상세정보 불러오기 (수집기 수집 결과)
	admin.GET("/partnerMemberAlarm", mngAdmin.GetBizPartnerMemberAlarmList)     //멤버십 상세정보 불러오기 (알림톡 결과)

	admin.GET("/fee", mngAdmin.GetFeeRateList)
	admin.PUT("/fee", mngAdmin.PutFeeRate)
	admin.POST("/kakaoWork", mngAdmin.SetWorkTimeAndVacation)

	company := admin.Group("/company")
	company.POST("/biz", mngAdmin.TransBizUse)         //기업유저 비즈 사용으로 전환
	company.POST("/priceDay", mngAdmin.CompayPriceDay) // 기업고객 기념일 금액 충전

	contents := api.Group("/contents")
	contents.GET("/board", mngAdmin.GetBoardList)
	contents.POST("/AddBoard", mngAdmin.AddBoard)
	contents.GET("/boardInfo", mngAdmin.GetBoardInfo)
	contents.PUT("/boardInfo", mngAdmin.ModBoardInfo)

	contents.GET("/content", mngAdmin.GetContentList)
	contents.POST("/AddContent", mngAdmin.AddContent)
	contents.GET("/contentInfo", mngAdmin.GetContentInfo)
	contents.PUT("/contentInfo", mngAdmin.ModContentInfo)

	contents.GET("/banner", mngAdmin.GetBannerList)
	contents.POST("/AddBanner", mngAdmin.AddBanner)
	contents.POST("/AddBannerImg", mngAdmin.AddBannerImg)
	contents.GET("/bannerInfo", mngAdmin.GetBannerInfo)
	contents.PUT("/bannerInfo", mngAdmin.ModBannerInfo)

	//OCR CHECK
	ocr := api.Group("/ocr") //api 그룹생성
	ocr.GET("/textList", mngAdmin.GetOcrTextList)
	ocr.GET("/textData", mngAdmin.GetTextData)
	ocr.PUT("/textUpdate", mngAdmin.SetReceiptText)

	ocr.GET("/receiptList", mngAdmin.GetReceiptList)
	ocr.GET("/receiptData", mngAdmin.GetReceiptData)
	ocr.POST("/receiptUpdate", mngAdmin.SetReceiptData)

	//Tpay결제
	tpay := api.Group("/tpays")
	tpay.POST("/ready", tpays.TpayReady)               // 결제 시작
	tpay.GET("/ready", tpays.TpayReady)                // 결제 시작 -- 테스트용
	tpay.POST("/payResult", tpays.TpayResult)          // 결제 반영
	tpay.GET("/payCancelCall", tpays.TpayCancelPage)   // 결제 취소
	tpay.POST("/rePayResult", tpays.TpayResult)        // 결제 재반영
	tpay.POST("/readyWebData", tpays.TpayReadyWebData) // 웹결제 준비 데이터

	tpay.POST("/GetUnpaidReady", tpays.TpayUnpaidReady) // 일괄 정산 준비 데이터
	tpay.POST("/payResult_unpaid", tpays.TpayUnpaidNew) // 일괄 정산 반영

	//tpay.POST("/payResult_unpaid_new", tpays.TpayUnpaidNew)    // 일괄 정산 반영

	//TPAY 빌링
	//tpay.GET("/readyBilling", tpays.TpayReadyBilling)               // 빌링키 생성

	//tpay.POST("/genBillikey", tpays.TpayGenBillikey)            		// 빌링키 생성
	//tpay.POST("/delBillkey", tpays.TpayDelbillkey)            			// 빌링키 삭제
	//tpay.POST("/billingPay", tpays.TpayBillingPay)           			// 빌링키 결제
	//tpay.POST("/billingPayCancel", tpays.TpayBillingPayCancel)         // 빌링키 결제 취소

	//신규
	biz_v2 := e.Group("/api/biz/v2")
	tpayV2 := biz_v2.Group("/tpays")

	tpayV2.GET("/pgReady", tpays.TpayReadyNew)          // PG 결제시작
	tpayV2.POST("/pgReady", tpays.TpayReadyNew)         // PG 결제시작
	tpayV2.POST("/pgPay", tpays.TpayResultNew)          // 일반결제 - PG 카드 & 계좌이체
	tpayV2.POST("/simplePay", tpays.TpayBillingPayment) // 간편결제

	tpayV2.POST("/payNorderReady", tpays.PayNorderReady)   // 주문 즉시 결제 - 일반 준비창 오픈
	tpayV2.POST("/payNorder", tpays.PayNorder)             // 주문 즉시 결제 - 일반결과 처리
	tpayV2.POST("/simplePayNorder", tpays.SimplePayNorder) // 주문 즉시 결제 - 간편

	// redis
	redis := e.Group("/redis")
	redis.GET("/update", Daily.UpdateRedis)

}
