package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Time struct {
	Data string
	Code int
	Msg  string
}

var time_api = "https://dailyhealth-api.sustech.edu.cn/api/util/getNowTime"
var save_api = "https://dailyhealth-api.sustech.edu.cn/api/form/save"

var sid = flag.String("id", "", "学号。")
var name = flag.String("name", "南大科", "姓名。")
var login_token = flag.String("token", "", "登录token")

func add_data() map[string]string {
	req_params := map[string]string{
		"jkmdz":       "",
		"xcmdz":       "",
		"sid":         *sid,
		"type":        "1",
		"xm":          *name,
		"deptName":    "本科2019级",
		"curCity":     "3",
		"dept":        "2019",
		"fanshenDate": "",
		"sfQz":        "",
		"sfZx":        "1",
		"sfQgwh":      "",
		"sfQghb":      "",
		"jcs":         "",
		"bxzz":        "",
		"tiwen":       "37",
		"local":       "美国",
		"formDate":    time.Now().Format("2000-01-01"),
		"nl":          "18",
		"jtzz":        "深圳市南方科技大学",
		"xb":          "男",
		"mobile":      "11111111111",
	}

	for i := 0; i < 120; i++ {
		req_params["ylzd"+strconv.Itoa(i)] = "0"
	}
	return req_params
}

func add_header(req *http.Request) {
	req.Header.Add("content-type", "application/json")
	req.Header.Add("userid", *sid)

	// req.Header.Add("accept", "application/json, text/plain, */*")
	// req.Header.Add("cache-control", "no-cache")
	// req.Header.Add("origin", "https://dailyhealth.sustech.edu.cn")
	// req.Header.Add("pragma", "no-cache")
	// req.Header.Add("referer", "https://dailyhealth.sustech.edu.cn/")
	// req.Header.Add("x-requested-with", "XMLHttpRequest")
}
func main() {

	flag.Parse()

	if *sid == "" {
		fmt.Println("请指定学号。")
		return
	}

	if *login_token == "" {
		fmt.Println("请指定健康申报系统token。")
		return
	}

	client := &http.Client{}

	req_params := add_data()
	time_rsp, _ := http.Get(time_api)
	var time_json Time
	err := json.NewDecoder(time_rsp.Body).Decode(&time_json)

	if err != nil {
		fmt.Println("parse err")
		return
	}

	req_token := *login_token + "-" + time_json.Data
	req_token_b64 := base64.StdEncoding.EncodeToString([]byte(req_token))
	req_params_byte, _ := json.Marshal(req_params)
	req, _ := http.NewRequest("POST", save_api, bytes.NewReader(req_params_byte))

	req.Header.Add("accesstoken", req_token_b64)

	add_header(req)

	rsp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(body))
}
