package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	PromethuesServerAddress string = "http://localhost:9090"

	//promethuesURL string = PromethuesServerAddress + "/api/v1/query?query="
)

type metric struct {
	Stage     string `json:"stage"`
	State     string `json:"state"`
	User      string `json:"user"`
	TimeRange string
}

type values []interface {
}
type metricValue struct {
	HitCount string
	Vtime    string
}

func (v values) String() (mValue metricValue) {

	for _, i := range v {
		switch items := i.(type) {
		case float64:
			timeT := time.Unix(int64(items), 0)

			mValue.Vtime = fmt.Sprint(timeT)
		case string:
			mValue.HitCount = items

		}

	}

	return
}

type result struct {
	Metrics metric   `json:"metric"`
	Value   []values `json:"values"`
}

type data struct {
	Results []result `json:"result"`
}
type queryResault struct {
	Datas data `json:"data"`
}

func (j queryResault) String() string {
	out := ""
	for rs := range j.Datas.Results {
		for va := range j.Datas.Results[rs].Value {
			out += fmt.Sprintf("%v    user: %v    stage: %v    %v    count: %v\n", j.Datas.Results[rs].Value[va].String().Vtime, j.Datas.Results[rs].Metrics.User, j.Datas.Results[rs].Metrics.Stage, strings.ToLower(j.Datas.Results[rs].Metrics.State), j.Datas.Results[rs].Value[va].String().HitCount)

		}
	}
	return out
}

func getRawData(url string) ([]byte, error) {
	promethuesURL := PromethuesServerAddress + "/api/v1/query?query="
	queryfqn := promethuesURL + url
	resp, getErr := http.Get(queryfqn)
	if getErr != nil {
		log.Println(getErr)
		return nil, errors.New("prometheus connection error!")
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
		return nil, readErr
	}
	return body, nil
}
func getQuery(query string) (queryResault, error) {
	body, readError := getRawData(query)
	if readError != nil {
		return queryResault{}, readError
	}

	logs := queryResault{}
	jsonError := json.Unmarshal(body, &logs)
	if jsonError != nil {
		log.Println(jsonError)
		return queryResault{}, jsonError

	}
	if len(logs.Datas.Results) == 0 {

		return queryResault{}, errors.New(" No Result\n you can increase Time Range\n e.g: 30m (30 minutes)\n 6h (6 houres)\n 1d (1 day)\n")
	}
	return logs, nil

}

func queryMakert(q metric) string {

	if q.TimeRange == "" {
		q.TimeRange = "1h"
	}
	if q.User == "" {
		q.User = ".*"
	}
	if q.Stage == "" {
		q.Stage = ".*"
	}
	if q.State == "" {
		q.State = ".*"
	}
	query := fmt.Sprintf(`radius_response{stage=~"%s",state=~"%s",user=~"%s"}[%s]`, q.Stage, q.State, q.User, q.TimeRange)
	return strings.ReplaceAll(query, " ", "%20")
}

func logs(w http.ResponseWriter, r *http.Request) {

	templ := templateHandler{filename: "logs.gohtml"}
	var q metric

	if r.Method == http.MethodPost {
		q.Stage = r.FormValue("stage")
		q.State = r.FormValue("state")
		q.User = r.FormValue("user")
		q.TimeRange = r.FormValue("timerange")
		//log.Println(q)

		qResulat, qErr := getQuery(queryMakert(q))
		if qErr != nil {
			templ.options = qErr
		} else {

			templ.options = qResulat
		}

	}

	templ.ServeHTTP(w, r)

}
