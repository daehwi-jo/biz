package mongo

import (
	"biz-web/src/controller/cls"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"strings"
)

/* log format */
// 로그 레벨(5~1:INFO, DEBUG, GUIDE, WARN, ERROR), 1인 경우 DB 롤백 필요하며, 에러 테이블에 저장
// darayo printf(로그레벨, 요청 컨텍스트, format, arg) => 무엇을(서비스, 요청), 어떻게(input), 왜(원인,조치)
var lprintf func(int, string, ...interface{}) = cls.Lprintf

type mongoStr struct{
	Name 		string  //`json:"name"`
	Query 		string 	//`json:"Query"`
	Kst     	string
	ResultCode 	string
	ResultMsg  	string
}

var mClient *mongo.Client

func MongoDBConfig(fname string){

	var useName, passWord, mIp, mPort string

	v, r := cls.GetTokenValue("MONGO_ID", fname)
	if r == cls.CONF_OK{
		useName = v
	}

	v, r = cls.GetTokenValue("MONGO_PASSWD", fname)
	if r == cls.CONF_OK{
		passWord = v
	}

	v, r = cls.GetTokenValue("MONGO_IP", fname)
	if r == cls.CONF_OK{
		mIp = v
	}

	v, r = cls.GetTokenValue("MONGO_PORT", fname)
	if r == cls.CONF_OK{
		mPort = v
	}

	credential := options.Credential{
		Username:   useName,
		Password:   passWord,
	}

	dbAddr := fmt.Sprintf("mongodb://%s:%s", mIp, mPort)
	clientOption := options.Client().ApplyURI(dbAddr).SetAuth(credential)
	//clientOption := options.Client().ApplyURI(dbAddr)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		cls.Lprintf(1, "[ERROR] mongo connect err(%s)\n", err.Error())
		return
	}

	cls.Lprintf(4, "[INFO] mongo connect ok\n")
	mClient = client
}

func SendMongoDB(c echo.Context) error {

	if mClient == nil{
		return c.JSON(http.StatusOK, nil)
	}

	params := cls.GetParamJsonMap(c)

	var appName, resultCode []string
	var currentPage, count int

	startDt := params["startDt"]
	if len(startDt) == 0{
		startDt = "2020-11-20 16:40:03"
	}

	endDt := params["endDt"]
	if len(endDt) == 0{
		endDt = "2025-12-29 16:40:06"
	}


	// cashapi,mocaapi,bizweb ...
	appNames := params["appName"]
	if len(appNames) == 0{
		appName = []string{"cashapi", "mocaapi", "bizweb", "partnerweb"}
	}else{
		as := strings.Split(appNames, ",")
		for _,v := range as{
			appName = append(appName, v)
		}
	}

	// log level, result code level...
	resultCodes := params["resultCode"]
	if len(resultCodes) == 0{
		resultCode = []string{"0", "1", "2", "3", "4", "5", "99"}
	}else{
		rs := strings.Split(resultCodes, ",")
		for _,v := range rs{
			resultCode = append(resultCode, v)
		}
	}

	currentPages := params["page"]
	if len(currentPages) == 0{
		currentPage = 1
	}else{
		currentPage,_= strconv.Atoi(currentPages)
	}

	limits := params["limit"]
	if len(limits) == 0{
		count = 200
	}else{
		count,_ = strconv.Atoi(limits)
	}

	var match bson.M
	search := params["search"]
	if len(search) == 0{
		match = bson.M{
			"$match": bson.M{
				// 기간 검색
				"kst":	bson.M{
					"$gte": startDt,
					"$lte": endDt,
				},
				// 단어 완전 일치해야 함
				"name": bson.M{
					//"$in": []string{"cashapi", "bizweb"},
					"$in": appName,
				},
				"resultCode": bson.M{
					//"$in": []string{"0", "1", "99"},
					"$in": resultCode,
				},
			},
		}
	}else{
		match = bson.M{
			"$match": bson.M{
				// 기간 검색
				"kst":	bson.M{
					"$gte": startDt,
					"$lte": endDt,
				},
				// 단어 완전 일치해야 함
				"name": bson.M{
					//"$in": []string{"cashapi", "bizweb"},
					"$in": appName,
				},
				"resultCode": bson.M{
					//"$in": []string{"0", "1", "99"},
					"$in": resultCode,
				},
				// 부분 검색
				"$text": bson.M{
					"$search": search,
				},
			},
		}
	}

	project := bson.M{
		"$project": bson.M{"_id": 0},
	}

	sort := bson.M{
		"$sort": bson.M{"kst": -1},
	}

	skip := bson.M{
		"$skip": (currentPage - 1) * count,
	}

	limit := bson.M{
		"$limit": count,
	}

	pipeline := []bson.M{match, project, sort, skip, limit}

	coll := mClient.Database("darayo").Collection("query")

	cur, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		cls.Lprintf(1, "[ERROR] mongo aggregate err(%s)\n", err.Error())
		return c.JSON(http.StatusOK, nil)
	}
	defer cur.Close(context.TODO())

	var m []mongoStr
	if err = cur.All(context.TODO(), &m); err != nil {
		cls.Lprintf(1, "[ERROR] mongo todo err(%s)\n", err.Error())
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, m)
}