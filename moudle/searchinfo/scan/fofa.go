package scan

import (
	"Road/moudle/common"
	"Road/moudle/sqlmoudle"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AutoGenerated struct {
	Mode    string     `json:"mode"`
	Error   bool       `json:"error"`
	Query   string     `json:"query"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
	Results [][]string `json:"results"`
}

func fofa_api(keyword string, email string, key string, page int, size int) string {
	input := []byte(keyword)
	encodeString := base64.StdEncoding.EncodeToString(input)
	api_request := fmt.Sprintf("https://fofa.info/api/v1/search/all?email=%s&page=%d&size=%d&key=%s&qbase64=%s&fields=ip,host,title,port,protocol", strings.Trim(email, " "), page, size, strings.Trim(key, " "), encodeString)
	return api_request
}

func fofahttp(url string, timeout string) *AutoGenerated {
	var itime, err = strconv.Atoi(timeout)
	if err != nil {
		log.Println("fofa超时参数错误: ", err)
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   time.Duration(itime) * time.Second,
		Transport: transport,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "*/*;q=0.8")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	res := &AutoGenerated{}
	json.Unmarshal(result, &res)
	return res
}

func Fafaall(fofaquery string) (urls []string) {
	Email := viper.Get("FofaEmail")
	Token := viper.Get("FofaToken")
	if Email == nil && Token == nil {
		log.Println("config file is empty,plz write your profile")
	}
	email := common.Strval(Email)
	Fofa_token := common.Strval(Token)
	Fofa_timeout := "10"
	for i := 1; i <= 20; i++ {
		url := fofa_api(fofaquery, email, Fofa_token, i, 500)
		res := fofahttp(url, Fofa_timeout)
		if len(res.Results) > 0 {
			for _, value := range res.Results {
				//fmt.Println(value[1])
				if strings.Contains(value[1], "http") {
					urls = append(urls, value[1])
				} else {
					urls = append(urls, "https://"+value[1])
				}
			}
		}
	}
	return urls
}

func FofaGet(Task_id string, qs string) {

	tempbanner := sqlmoudle.Bannerresult{
		Task_Id:     Task_id,
		Target:      "",
		Banner:      sql.NullString{String: "0", Valid: true},
		Server:      sql.NullString{String: "0", Valid: true},
		Status_Code: sql.NullString{},
		Title:       sql.NullString{String: "0", Valid: true},
		Last_time:   sql.NullTime{},
		Pocmatch:    sql.NullInt16{Int16: 0, Valid: true},
	}
	tempbanner.Banner.String = "0"
	tempbanner.Title.String = "0"
	//FofaRq(qs, FofaEmail, FofaToken, Url) // bug: 只請求了一次。可能shodan模塊也是這樣
	res := Fafaall(qs)
	for _, targets := range res {
		if len(targets) > 0 {
			tempbanner.Target = targets
			sqlmoudle.Insertbanner(tempbanner)
			//log.Println(tempbanner)
		}
	}
	log.Println("finfish")
}