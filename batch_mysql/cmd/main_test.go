package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"
	"wiwieo/batch_mysql/constant"
	"wiwieo/batch_mysql/server"
)

const ad = `{
	"title":"baidu",
	"as_id":%d,
	"pic":"www.baidu.com",
	"url":"www.url.com",
	"sort":1,
	"open":1,
	"content":"this is content.%d"
}`

const dt = `{
	"title":"baidu",
	"sort":%d
}`

const AdUrl = `http://192.168.244.128:9999/add-ad`

//const AdUrl = `http://localhost:9999/add-ad`
const DtUrl = `http://localhost:9999/add-dt`

const (
	GoroutinCount = 1000
	RequestCount  = 2000
)

func BenchmarkRequestAd(b *testing.B) {
	cli := &http.Client{
		CheckRedirect: nil,
	}
	for i := 0; i < b.N; i++ {
		request, _ := http.NewRequest("POST", AdUrl, bytes.NewReader([]byte(fmt.Sprintf(ad, i, i))))
		request.Close = true
		_, err := cli.Do(request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestRequestAd(t *testing.T) {
	cli := &http.Client{
		CheckRedirect: nil,
	}
	for i := 0; i < GoroutinCount; i++ {
		go func(i int) {
			for {
				postAd(cli, i, i)
				time.Sleep(10 * time.Second)
			}
		}(i)
	}
}

func postAd(cli *http.Client, i, j int) {
	request, _ := http.NewRequest("POST", AdUrl, bytes.NewReader([]byte(fmt.Sprintf(ad, i, j))))
	request.Close = true
	_, err := cli.Do(request)
	if err != nil {
		panic(err)
	}
}

func postDt(cli *http.Client, i int) {
	request, _ := http.NewRequest("POST", DtUrl, bytes.NewReader([]byte(fmt.Sprintf(dt, i))))
	request.Close = true
	_, err := cli.Do(request)
	if err != nil {
		panic(err)
	}
}

func BenchmarkAddContent(b *testing.B) {
	s := server.NewSrv()
	content := `
{
	"title":"baidu",
	"as_id":%d,
	"pic":"www.baidu.com",
	"url":"www.url.com",
	"sort":1,
	"open":1,
	"content":"this is content.%d"
}
`
	for i := 0; i < b.N; i++ {
		println(i)
		err := s.AdCache.AddContent([]byte(fmt.Sprintf(content, i, i)))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestAddContent(t *testing.T) {
	s := server.NewSrv()
	content := `
{
	"title":"baidu",
	"as_id":%d,
	"pic":"www.baidu.com",
	"url":"www.url.com",
	"sort":1,
	"open":1,
	"content":"this is content.%d"
}
`
	count := 1000000
	for i := 0; i < count; i++ {
		err := s.AdCache.AddContent([]byte(fmt.Sprintf(content, i, i)))
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(time.Duration(constant.Config.MaxTriggerTime+10)*time.Second)
}
