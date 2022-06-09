package main

import (
	ApiSrc "biz-web/src"
	Daily "biz-web/src/controller"
	"biz-web/src/controller/api/mongo"
	"biz-web/src/controller/cls"
	"net/http"
	"os"
	"strings"
	"time"

	echotemplate "biz-web/src/controller/cls/echotemplate"
	"github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm/module/apmechov4"
	//Daily "biz-web/src/controller"
	commons "biz-web/src/controller"
)

var lprintf func(int, string, ...interface{}) = cls.Lprintf

func main() {

	// test 12
	fname := cls.Cls_conf(os.Args)

	domains := cls.WebConf(fname)

	if cls.Db_conf(fname) < 0 {
		lprintf(4, "DataBase not setting \n")
		return
	}
	defer cls.DBc.Close()

	if cls.Db_conf2(fname) < 0 {
		lprintf(4, "DataBase2 not setting \n")
		return
	}

	if commons.Tpay_conf(fname) < 0 {
		lprintf(4, "Tpay_conf not setting \n")
	}

	defer cls.DBc2.Close()

	mongo.MongoDBConfig(fname)

	go cls.LogCollect("bizweb", fname)

	//Daily.CompayPriceDay()

	// 스케쥴 일정 등록
	s := gocron.NewScheduler(time.Local)

	// 기업 기념일 금액 충전
	s.Every(1).Day().At("15:00").Do(Daily.CompayPriceDay, false, "", "")

	// 윈큐브 상품 업데이트
	s.Every(1).Day().At("06:01").Do(Daily.WincubeItemUpdate)

	// 티페이 정산 지급 요청  데이터 생성
	s.Every(1).Day().At("01:01").Do(Daily.TpayMakePayStep1)
	// 매일 새벽 식대지원금 백업
	s.Every(1).Day().At("01:11").Do(Daily.BackupSupportBalance)
	// 지정된 일자까지 미사용된 선물 취소
	s.Every(1).Day().At("02:01").Do(Daily.UnusedGiftCancel)
	// 매월 초 지원금 초기화
	s.Every(1).Month(1).At("01:20").Do(Daily.MonthResetSupportBalance)
	// 스케줄러 결과 카카오 워크 전송
	s.Every(1).Day().At("09:30").Do(Daily.SendDailyJobResult)

	// 카카오 알림톡
	Daily.DailyConfig(fname)
	// 매일 성공 -> 매일 11:00 실행
	s.Every(1).Day().At("11:00").Do(Daily.YesterdayReport1)
	// 매일 실패 -> 매일 14:00 실행
	s.Every(1).Day().At("14:00").Do(Daily.YesterdayReport2)
	// 주간 -> 매주 월요일 15:30분 실행
	s.Every(1).Monday().At("15:30").Do(Daily.WeekReport)
	// 월간 -> 매월 2일 15:30분 실행
	s.Cron("30 15 2 * *").Do(Daily.MonthReport)

	//TODO 구글 캘린더 토큰 갱신 프로세스 제작해야함
	//s.Every(1).Day().At("12:00").Do(googleToken.RefreshGoogleToken)

	// redis 매출 데이터 insert/delete
	rSchd, r := cls.GetTokenValue("REDISSCHED", fname)
	if r != cls.CONF_ERR {
		redisSchd := strings.Split(rSchd, ",")
		for _, schedule := range redisSchd {
			s.Every(1).Day().At(schedule).Do(Daily.RedisInsert)
		}
	}

	s.StartAsync()
	defer s.Clear()

	// single domain
	if len(domains) == 1 {
		e := echo.New()
		e.Use(apmechov4.Middleware())

		e.Static("/public", "src/public")
		e.Static("/SharedStorage", "/app/SharedStorage")

		e.Renderer = echotemplate.New(echotemplate.TemplateConfig{
			Root:         "src/templates",
			Extension:    ".htm",
			Master:       "master",
			Partials:     []string{"/master"},
			DisableCache: true,
			Delims:       echotemplate.Delims{Left: "[[", Right: "]]"},
		}) //웹페이지 템플릿 로드 등

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			//	AllowOrigins: []string{"https://172.30.1.22", "https://172.30.1.22:8080"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}))

		// Web page Service
		e = ApiSrc.SvcSetting(e, fname)     //api 세팅
		e = ApiSrc.SvcSettingPage(e, fname) // 세팅

		// e = MocaApp.SvcSetting(e, fname)
		// start Web Server
		domains[0].EchoData = e
		cls.StartDomain(domains)
		return
	}

	for i := 0; i < len(domains); i++ {
		cls.Lprintf(4, "domains : %s\n", domains[i].Domain)
	}
	cls.StartDomain(domains)
}
