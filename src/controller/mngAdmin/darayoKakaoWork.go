package mngAdmin

import (
	"biz-web/src/controller"
	"biz-web/src/controller/cls"
	_struct "biz-web/src/controller/struct/kakaoWk"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func SetWorkTimeAndVacation(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	calendarId, _ := cls.GetTokenValue("GOOGLE_CALENDAR_ID", fname)
	kakaoApi, _ := cls.GetTokenValue("KAKAO_WORK_API", fname)

	//카카오워크 근무시간 업데이트
	setWorkTimeResult, err := setWorkTime(params, kakaoApi)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//카카오워크 휴가 정보 불러와서 구글 캘린더에 업데이트
	setVacationResult, err := setVacation(params["token"], kakaoApi, calendarId)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "근무시간 업데이트 완료 휴가배치 오류 :"+err.Error()))
	}

	result := make(map[string]interface{})
	result["kakaoWorkResult"] = setWorkTimeResult
	result["googleCalenderResult"] = setVacationResult

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)
}

func setVacation(token string, kApi string, calendarId string) ([]interface{}, error) {
	var result []interface{}

	start := time.Now().Format("2006-01-02")                //오늘날짜
	end := time.Now().AddDate(0, 0, 7).Format("2006-01-02") //오늘 +7 날짜

	req, _ := http.NewRequest("GET", kApi+"/vacApply?searchStartDate="+start+"&searchEndDate="+end, nil)
	setKakaoWorkRequestHeader(req, token)

	client := http.Client{}
	defer client.CloseIdleConnections()

	res, err := client.Do(req)

	body, err := ioutil.ReadAll(res.Body)

	var vacation _struct.Vacation
	err = json.Unmarshal(body, &vacation)
	if err != nil {
		return nil, errors.New("카카오워크 휴가 데이터 파싱 오류 " + err.Error())
	}

	//카카오워크 휴가 데이터 받아와서 데이터 재가공
	vacationList := make(map[int]interface{})
	for i, data := range vacation.VacApplyList {
		tmp := map[string]string{}
		tmp["title"] = data.Employee.Name + " " + data.VacationCode.VacationFullName
		tmp["startDate"] = data.StartDate
		tmp["startTime"] = data.StartTime
		tmp["endDate"] = data.EndDate
		tmp["endTime"] = data.EndTime
		tmp["description"] = data.Description
		vacationList[i] = tmp
	}

	ctx := context.Background()
	defer ctx.Done()

	b, err := ioutil.ReadFile("conf/google_client_secret_key.json")
	if err != nil {
		return nil, errors.New("구글 키 파일 로드 오류 " + err.Error())
	}

	configs, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, errors.New("구글 캘린더 설정 오류 " + err.Error())
	}

	gClient, err := getGoogleClient(configs)
	if err != nil {
		return nil, errors.New("구글 api 클라이언트 시작 오류 " + err.Error())
	}

	service, err := serviceStart(ctx, gClient)
	if err != nil {
		return nil, errors.New("구글 api 클라이언트 서비스 시작 오류" + err.Error())
	}
	defer gClient.CloseIdleConnections()

	start = time.Now().Format(time.RFC3339)
	end = time.Now().AddDate(0, 0, 7).Format(time.RFC3339)

	events, err := service.Events.List(calendarId).ShowDeleted(true).
		SingleEvents(true).TimeMin(start).TimeMax(end).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		return nil, errors.New("구글 캘린더 이벤트 로드 오류 " + err.Error())
	}

	//구글 캘린더 데이터 취소되지 않은것만 선별하여 데이터 재가공
	calendarEvent := make(map[int]interface{})
	for i, data := range events.Items {
		if data.Status == "confirmed" {
			startArr := strings.Split(data.Start.DateTime, "T")
			endArr := strings.Split(data.End.DateTime, "T")

			tmp := map[string]string{}
			tmp["title"] = data.Summary
			tmp["description"] = data.Description
			tmp["startDate"] = startArr[0]
			tmp["startTime"] = strings.Split(startArr[1], "+")[0]
			tmp["endDate"] = endArr[0]
			tmp["endTime"] = strings.Split(endArr[1], "+")[0]
			calendarEvent[i] = tmp
		}
	}

	//구글 캘린더 이벤트 삽입 메소드
	setEventFunc := func(service *calendar.Service, data map[string]string) error {
		start, _ := getRFC3339Time(data["startDate"], data["startTime"])
		end, _ := getRFC3339Time(data["endDate"], data["endTime"])

		//이벤트 정의
		event := &calendar.Event{
			Summary:     data["title"],
			Description: data["description"],
			Start: &calendar.EventDateTime{
				DateTime: start,
				TimeZone: "Asia/Seoul",
			},
			End: &calendar.EventDateTime{
				DateTime: end,
				TimeZone: "Asia/Seoul",
			},
		}

		//이벤트 삽입
		event, err := service.Events.Insert(calendarId, event).Do()
		if err != nil {
			return err
		}
		return nil
	}

	//카카오 휴가 리스트와 구글 캘린더 휴가 이벤트 비교 해서 작성되지 않은 휴가 삽입
	for _, data := range vacationList {
		kakaoVacation := data.(map[string]string)
		check := true
		for _, data := range calendarEvent {
			googleVacation := data.(map[string]string)
			if googleVacation["title"] == kakaoVacation["title"] &&
				googleVacation["endDate"] == kakaoVacation["endDate"] &&
				googleVacation["startDate"] == kakaoVacation["startDate"] {
				check = false
				break
			}
		}
		if check {
			err = setEventFunc(service, kakaoVacation)
			if err != nil {
				return nil, errors.New("구글 캘린더 이벤트 삽입 오류 " + err.Error())
			}
			result = append(result, kakaoVacation)
		}
	}

	return result, nil
}

func setWorkTime(params map[string]string, kApi string) ([]interface{}, error) {
	var result []interface{}

	req, _ := http.NewRequest("GET", kApi+"/commuteReport/?page=0&size=100&startDate="+params["startDate"]+"&endDate="+params["endDate"]+"&vacationCodes=&notEmpty=true&onEmpty=true&offEmpty=true", nil)
	setKakaoWorkRequestHeader(req, params["token"])

	client := http.Client{}
	defer client.CloseIdleConnections()

	res, err := client.Do(req)
	b, err := ioutil.ReadAll(res.Body)

	var workUser _struct.WorkUser
	err = json.Unmarshal(b, &workUser)
	if err != nil {
		return nil, err
	}

	//카카오워크 사원 근무시간 삽입 메소드
	workTimeInsert := func(token string, seqNumber string) (map[string]interface{}, error) {
		start := strings.Split(params["startTime"], ":")
		end := strings.Split(params["endTime"], ":")

		payload, _ := json.Marshal(map[string]interface{}{
			"onHour":      start[0],
			"onMinute":    start[1],
			"offHour":     end[0],
			"offMinute":   end[1],
			"description": "COMMUTE_COMMENT_004",
		})

		req, _ := http.NewRequest("PATCH", kApi+"/commuteReport/"+seqNumber, bytes.NewBuffer(payload))
		setKakaoWorkRequestHeader(req, params["token"])

		res, err := client.Do(req)
		b, err = ioutil.ReadAll(res.Body)

		var workTime _struct.WorkTime
		err = json.Unmarshal(b, &workTime)
		if err != nil {
			return nil, err
		}

		if workTime.Result != "success" {
			err = errors.New("not success " + seqNumber)
		}

		data := make(map[string]interface{})
		data["name"] = workTime.CommuteReport.EmpName
		data["seq"] = workTime.CommuteReport.CommuteSeq
		data["date"] = workTime.CommuteReport.CommuteDay
		data["startTime"] = workTime.CommuteReport.WorkOnTime
		data["endTime"] = workTime.CommuteReport.WorkOffTime
		data["result"] = workTime.Result

		return data, nil
	}

	//카카오워크 사원 근무시간 삽입
	for _, data := range workUser.Content {
		if data.CommuteSeq != 0 {
			insertResult, err := workTimeInsert(params["token"], fmt.Sprintf("%v", data.CommuteSeq))
			if err != nil {
				return nil, errors.New("카카오워크 사원 근무시간 삽입 오류 " + err.Error())
			}
			result = append(result, insertResult)
		}
	}

	return result, nil
}

func setKakaoWorkRequestHeader(req *http.Request, token string) {
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", token)
	req.Header.Set("Referer", "https://schedule-admin.we.kakaowork.com/manageTA/taRecord")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Host", "schedule-admin.we.kakaowork.com")
}

func getRFC3339Time(mDate string, mTime string) (string, error) {
	arras := func(strArr []string) []int {
		var intArr []int
		for _, s := range strArr {
			ia, _ := strconv.Atoi(s)
			intArr = append(intArr, ia)
		}
		return intArr
	}

	date := arras(strings.Split(mDate, "-"))
	times := arras(strings.Split(mTime, ":"))

	if len(date) != 3 {
		return "", errors.New("날짜 파라미터 오류")
	}
	if len(times) != 3 {
		return "", errors.New("타임 파라미터 오류")
	}
	location := time.Now().Location()
	result := time.Date(date[0], time.Month(date[1]), date[2], times[0], times[1], times[2], 00, location).Format(time.RFC3339)

	return result, nil
}

func getGoogleClient(config *oauth2.Config) (*http.Client, error) {
	tokFile := "conf/googleToken.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok), nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.New("에러")
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, errors.New("에러")
	}
	return tok, nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.New("에러")
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

func serviceStart(ctx context.Context, client *http.Client) (*calendar.Service, error) {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.New("에러")
	}
	return srv, err
}

func RefreshGoogleToken() error {
	b, err := ioutil.ReadFile("conf/google_client_secret_key.json")
	if err != nil {
		return errors.New("구글 키 파일 로드 오류 " + err.Error())
	}

	configs, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		fmt.Println("토큰 갱신오류", " 구글 캘린더 설정 오류 ", err.Error())
		return errors.New("구글 캘린더 설정 오류 " + err.Error())
	}

	_, err = getGoogleClient(configs)
	if err != nil {
		fmt.Println("토큰 갱신오류", " 구글 api 클라이언트 시작 오류 ", err.Error())
		return errors.New("구글 api 클라이언트 시작 오류 " + err.Error())
	}

	fmt.Println("토큰 갱신됨")

	return nil
}
