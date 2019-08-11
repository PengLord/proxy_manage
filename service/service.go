package service

import (
	"time"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"log"
	"io/ioutil"
	"proxy_manage/config"
	"encoding/json"
	"proxy_manage/pojo"
	"proxy_manage/dao"
	"net/url"
)

type Proxies struct {
	Code    int      `json:"code"`
	Proxies []string `json:"proxies"`
}

func Callback(con config.Config) {
	//1通过第三方接口获取代理数据
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	result, err := httpClient.Get(con.ProxyUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer result.Body.Close()
	buff, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Print(err)
		return
	}
	proxies := Proxies{}
	err = json.Unmarshal(buff, &proxies)
	if err != nil {
		log.Println(err)
	}

	//2将代理数据存入道数据库
	for _, v := range proxies.Proxies {
		pro := pojo.ProxyPool{}
		pro.Proxy = v
		pro.CreatedAt = time.Now().Unix()

		h := sha1.New()
		h.Write([]byte(v))
		buff := h.Sum(nil)
		pro.Id = hex.EncodeToString(buff)
		dao.Save(pro)
	}
	log.Printf("----本次共录入%d代理IP-----\n", len(proxies.Proxies))
}

func Check(valid int) {
	pos := dao.FindAll()
	count := 0
	for _, v := range pos {
		checkurl := "https://baidu.com"
		if v.IsCn == 1 {
			checkurl = "https://google.com"
		}

		transport := &http.Transport{Proxy: func(request *http.Request) (*url.URL, error) {
			return url.Parse("http://" + v.Proxy)
		}}
		//1通过第三方接口获取代理数据
		httpClient := &http.Client{
			Timeout:   time.Duration(valid) * time.Second,
			Transport: transport,
		}
		_, err := httpClient.Get(checkurl)
		if err != nil {
			count++
			dao.DeleteById(v.Id)
		}
	}
	log.Printf("----本次清理无效IP数量是:%d---\n", count)
}
