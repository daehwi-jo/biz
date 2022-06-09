package biz

import (
	mains "biz-web/src/controller/mains"
	"biz-web/src/controller/tpays"

	//공통 api
	"biz-web/src/controller/cls"

	"github.com/labstack/echo/v4"
)

func SvcSettingPage(e *echo.Echo, fname string) *echo.Echo {
	lprintf(4, "[INFO] sql start \n")

	page(e)

	cls.SetNotLoginUrl("/")
	lprintf(4, "[INFO] page start \n")

	return e
}

func page(e *echo.Echo) {

	//홈페이지 라우터
	page := e.Group("")

	//메인
	page.GET("/", mains.Home)                  //로그인 화면
	page.GET("/login", mains.Login)            //로그인 화면
	page.GET("/home", mains.Home2)             //홈 화면
	page.GET("/search", mains.Search)          //검색화면
	page.POST("/sysSetting", mains.SysSetting) //검색화면

	//달아요 관리
	page.GET("/admin/combine", mains.Combine)                        //통합 정산 대상
	page.GET("/admin/combineDesc", mains.CombineDesc)                //통합 정산 대상 상세
	page.GET("/admin/combineDesc_wincube", mains.CombineDescWincube) //통합 정산 대상 상세 -윈큐브

	page.GET("/admin/paidMng", mains.PaidMng)   //정산대상
	page.GET("/admin/paidList", mains.PaidList) //지급요청대상
	page.GET("/admin/paidIng", mains.PaidIng)   //지급요청처리중
	page.GET("/admin/paidOk", mains.PaidOk)     //지급요청완료
	page.GET("/admin/account", mains.Account)   //자금현황 조회

	//회원관리
	page.GET("/admin/member", mains.AdminMember)         //회원 관리 	화면
	page.GET("/admin/memberInfo", mains.AdminMemberInfo) //회원 관리 상세 페이지

	//장부관리
	page.GET("/admin/book", mains.AdminBook)         //장부 관리 	화면
	page.GET("/admin/bookInfo", mains.AdminBookInfo) //장부 관리 상세 페이지

	//기업관리
	page.GET("/admin/company", mains.AdminCompany)         //기업 관리 	화면
	page.GET("/admin/companyInfo", mains.AdminCompanyInfo) //기업 관리 상세 페이지

	//가맹점 관리
	page.GET("/admin/store", mains.AdminStore)         //가맹점 관리 	화면
	page.GET("/admin/storeInfo", mains.AdminStoreInfo) //가맹점 관리 상세 페이지

	//멤버쉽 관리
	page.GET("/admin/partnerMember", mains.AdminPartnerMember)         //멤버쉽 관리 	화면
	page.GET("/admin/partnerMemberInfo", mains.AdminPartnerMemberInfo) //멤버쉽 관리 상세 화면

	//주문관리
	page.GET("/admin/order", mains.AdminOrder)         //주문 관리 	화면
	page.GET("/admin/orderInfo", mains.AdminOrderInfo) //주문 관리 상세 페이지

	page.GET("/admin/payment", mains.AdminPayment)         //결제 관리
	page.GET("/admin/paymentInfo", mains.AdminPaymentInfo) //결제 관리 상세 페이지

	//영수증 검증
	page.GET("/admin/ocrText", mains.AdminOCRText)       //영수증 검증 화명
	page.GET("/admin/ocrReceipt", mains.AdminOCRReceipt) //영수증 검증 화명

	//기타관리
	page.GET("/admin/task", mains.AdminTask) //스케줄 관리 	화면
	page.GET("/admin/log", mains.AdminLog)   //로그 관리 	화면
	page.GET("/admin/kakaoWork", mains.AdminKakaoWork)

	//컨텐츠관리
	page.GET("/admin/content", mains.AdminContent)
	page.GET("/admin/contentInfo", mains.AdminContentInfo)
	page.GET("/admin/board", mains.AdminBoard)
	page.GET("/admin/boardInfo", mains.AdminBoardInfo)
	page.GET("/admin/banner", mains.AdminBanner)
	page.GET("/admin/bannerInfo", mains.AdminBannerInfo)

	page.GET("/admin/feeRate", mains.AdminFeeRate)

	// 사용자 관리
	page.GET("/user/ManagerInfo", mains.MngManagerInfo)     //관리자 정보 관리
	page.GET("/user/AnniversaryInfo", mains.MngAnniversary) //기념일 관리

	//장부 관리
	page.GET("/book/GrpBookList", mains.MngGrpList)    //장부 리스트
	page.GET("/book/GrpBookUserMng", mains.MngGrpUser) //장부 사용자 관리

	//주문 관리
	page.GET("/order/OrderList", mains.MngOrderList)       //주문 내역 조회
	page.GET("/order/CalculateMng", mains.MngCalculateMng) //장부 정산관리
	page.GET("/order/PaymentList", mains.MngPaymentList)   //결제 내역 조회
	page.GET("/order/StoreList", mains.MngGrpRestList)     //가맹점 조회
	page.GET("/order/payCharge", tpays.MngPayCharge)       //가맹점 조회 - 충전하기
	page.GET("/order/payCalculate", mains.MngPayCalculate) //가맹점 조회 - 정산하기

	//기타
	page.GET("/address", mains.Address) //우편번호 찾기
}
