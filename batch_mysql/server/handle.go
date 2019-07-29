package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func (s *srv) AdHandle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("request parameters is wrong!"))
		return
	}
	err = s.AdCache.AddContent(body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("add content is wrong! %s", err)))
		return
	}
	w.Write([]byte("It's success."))
}

func (s *srv) DebrisTypeHandle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("request parameters is wrong!"))
		return
	}
	err = s.DebrisTypeCache.AddContent(body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("add content is wrong! %s", err)))
		return
	}
	w.Write([]byte("It's success."))
}

func (s *srv) TestHandle(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
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
	goroutinCount, _ := strconv.Atoi(r.FormValue("gcount"))
	requestCount, _ := strconv.Atoi(r.FormValue("rcount"))
	if goroutinCount == 0 {
		goroutinCount = 100
	}
	if requestCount == 0 {
		requestCount = 1000
	}
	wait := sync.WaitGroup{}
	wait.Add(goroutinCount * requestCount)
	for j := 0; j < goroutinCount; j++ {
		go func(j int) {
			for i := 0; i < requestCount; i++ {
				err := s.AdCache.AddContent([]byte(fmt.Sprintf(content, j, i)))
				if err != nil {
					panic(err)
				}

				//err = s.DebrisTypeCache.AddContent([]byte(fmt.Sprintf(dt, i)))
				//if err != nil{
				//	panic(err)
				//}
				wait.Done()
			}
		}(j)
	}
	wait.Wait()
	w.Write([]byte(fmt.Sprintf("It's success, cost time: %s", time.Since(now))))
}
