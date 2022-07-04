package common

import (
	"Road/moudle/finger"
	"Road/moudle/queue"
	"Road/moudle/sqlmoudle"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"plugin"
	"strconv"
	"strings"
	"time"
)

type configini struct {
	Fingerpath string
	Httpproxy  string
	Sqlstr     string
	Fileloader string
}

var (
	configjson configini
	Finpx      *finger.Packjson
)

func TaskInit(path string) {
	dir_list, err := ioutil.ReadDir(path)

	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range dir_list {
		p, err := plugin.Open(path + v.Name())
		if err != nil {
			log.Println(err)
			continue
		}
		if p == nil {
			continue
		}
	}
	//fmt.Println(comm.Attackplugins)
}

type Attackplugin interface {
	Attack(target string, filepath string) bool
}

type PluginInfo struct {
	Name        string
	Description string
	Query       string
}

var Attackplugins = make(map[string][]Attackplugin)

func Regist(target string, plugin Attackplugin) {
	Attackplugins[target] = append(Attackplugins[target], plugin)
}

func Attackmatch(target string, banner string) bool {
	flag := false
	for name, list := range Attackplugins {
		if strings.Contains(name, banner) {
			for _, plugin := range list {
				flag = plugin.Attack(target, configjson.Fingerpath)
				return flag
			}
		}
	}
	return flag
}

//随机UA头
func rndua() string {
	ua := []string{"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.1 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2226.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1",
		"Mozilla/5.0 (Windows NT 6.3; rv:36.0) Gecko/20100101 Firefox/36.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10; rv:33.0) Gecko/20100101 Firefox/33.0",
		"Mozilla/5.0 (X11; Linux i586; rv:31.0) Gecko/20100101 Firefox/31.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20130401 Firefox/31.0",
		"Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
		"Mozilla/5.0 (compatible, MSIE 11, Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows; Intel Windows) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67"}
	n := rand.Intn(13) + 1
	return ua[n]
}

func Request(target string, client *http.Client, banner sqlmoudle.Bannerresult, requestype string) {
	req, err := http.NewRequest("GET", target, nil)
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", rndua())
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Referer", "https://www.google.com/")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[-]: ", err)
		log.Println("[-] DELETE id:", banner.Id)
		sqlmoudle.DeleteTask(banner.Id)
		return
	}
	defer resp.Body.Close()
	out := finger.Checkbanner(target, resp, Finpx, banner)
	sqlmoudle.Upadtetask(sqlmoudle.Taskjob{
		Task_Id:  banner.Task_Id,
		Progress: 2,
	})
	tempbanner := sqlmoudle.Bannerresult{Id: banner.Id, Task_Id: banner.Task_Id, Target: target, Banner: sql.NullString{"NULL", true}, Server: sql.NullString{out.Server, true}, Status_Code: sql.NullString{strconv.Itoa(out.Statuscode), true}, Title: sql.NullString{out.Title, true}, Pocmatch: sql.NullInt16{
		Int16: 0,
		Valid: true,
	}}
	if requestype == "update" {
		if len(out.Cms) > 50 {
			//banner := sqlmoudle.Bannerresult{Id: banner.Id, Target: target, Banner: sql.NullString{"NULL", true}, Server: sql.NullString{out.Server, true}, Status_Code: sql.NullString{strconv.Itoa(out.Statuscode), true}, Title: sql.NullString{out.Title, true}, Pocmatch: sql.NullInt16{
			//	Int16: 0,
			//	Valid: true,
			//}}
			tempbanner.Banner.String = "Honeypot"
			tempbanner.Title.String = "Honeypot"
			sqlmoudle.UpdateTask(tempbanner)
			return
		}
		if out.Cms != "" {
			//banner := sqlmoudle.Bannerresult{Id: banner.Id, Task_Id: banner.Task_Id, Target: target, Banner: sql.NullString{out.Cms, true}, Server: sql.NullString{out.Server, true}, Status_Code: sql.NullString{strconv.Itoa(out.Statuscode), true}, Title: sql.NullString{out.Title, true}, Pocmatch: sql.NullInt16{
			//	Int16: 0,
			//	Valid: true,
			//}}
			tempbanner.Banner.String = out.Cms
			if Attackmatch(target, out.Cms) {
				tempbanner.Pocmatch.Int16 = 1
			} else {
				tempbanner.Pocmatch.Int16 = 0
			}
			sqlmoudle.UpdateTask(tempbanner)
		} else {
			//banner := sqlmoudle.Bannerresult{Id: banner.Id, Task_Id: banner.Task_Id, Target: target, Banner: sql.NullString{"NULL", true}, Server: sql.NullString{out.Server, true}, Status_Code: sql.NullString{strconv.Itoa(out.Statuscode), true}, Title: sql.NullString{out.Title, true}, Pocmatch: sql.NullInt16{
			//	Int16: 0,
			//	Valid: true,
			//}}

			sqlmoudle.UpdateTask(tempbanner)
		}
		if len(out.Jsurl) > 0 {
			requestype := "insert"
			Request(out.Jsurl[0], client, tempbanner, requestype)
		}
	} else if requestype == "insert" {
		if len(out.Cms) > 30 {
			//banner := sqlmoudle.Bannerresult{Task_Id: banner.Task_Id, Target: target, Banner: sql.NullString{"NULL", true}, Server: sql.NullString{out.Server, true}, Status_Code: sql.NullString{strconv.Itoa(out.Statuscode), true}, Title: sql.NullString{out.Title, true}, Pocmatch: sql.NullInt16{
			//	Int16: 0,
			//	Valid: true,
			//}}
			tempbanner.Banner.String = "Honeypot"
			tempbanner.Title.String = "Honeypot"
			sqlmoudle.Insertbanner(tempbanner)
			return
		}
		if out.Cms != "" {
			//banner := sqlmoudle.Bannerresult{Id: banner.Id, Task_Id: banner.Task_Id, Target: target, Banner: sql.NullString{out.Cms, true}, Server: sql.NullString{out.Server, true}, Status_Code: sql.NullString{strconv.Itoa(out.Statuscode), true}, Title: sql.NullString{out.Title, true}, Pocmatch: sql.NullInt16{
			//	Int16: 0,
			//	Valid: true,
			//}}
			tempbanner.Banner.String = out.Cms
			if Attackmatch(target, out.Cms) {
				tempbanner.Pocmatch.Int16 = 1
			} else {
				tempbanner.Pocmatch.Int16 = 0
			}
			sqlmoudle.Insertbanner(tempbanner)
		} else {
			sqlmoudle.Insertbanner(tempbanner)
		}
	}
}

func Readconfig() {
	path := "moudle/file/config.json"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &configjson)
	if err != nil {
		return
	}
	return
}

func Makeclient() *http.Client {
	if configjson.Httpproxy != "" {
		proxy, _ := url.Parse(configjson.Httpproxy)
		tr := &http.Transport{
			//关闭证书验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			//设置超时
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
			//设置代理
			Proxy: http.ProxyURL(proxy),
		}
		client := &http.Client{
			Transport: tr,
		}
		return client
	} else {
		tr := &http.Transport{
			//关闭证书验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			//设置超时
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
		}
		client := &http.Client{
			Transport: tr,
		}
		return client
	}
}

func Getfinger(Taskid string) {
	fingerpath := configjson.Fingerpath
	errs := finger.LoadWebfingerprint(fingerpath)
	if errs != nil {
		log.Println("[-] Webfingerprint file error!!!")
		return
	}
	Finpx = finger.GetWebfingerprint()
	target := sqlmoudle.Queryinfo(configjson.Sqlstr + " AND Task_Id=\"" + Taskid + "\"")
	client := Makeclient()
	//定义协程池，设置最大数量
	pool := queue.New(100)
	for i := 0; i < len(target); i++ {
		pool.Add(1)
		go func(i int) {
			Request(target[i].Target, client, target[i], "update")
			pool.Done()
		}(i)
	}
	pool.Wait()
}
