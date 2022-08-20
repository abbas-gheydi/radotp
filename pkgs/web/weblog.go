package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	promethuesServerAddress string = "http://localhost:9090"

	promethuesURL string = promethuesServerAddress + "/api/v1/query?query="
)

type metric struct {
	Stage     string `json:"stage"`
	State     string `json:"state"`
	User      string `json:"user"`
	TimeRange string
}

type values []interface {
}

func (v values) String() string {
	out := ""
	for _, i := range v {
		switch items := i.(type) {
		case float64:
			timeT := time.Unix(int64(items), 0)

			out = fmt.Sprint(timeT)

		}
	}

	return out
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
			out += fmt.Sprintf("%v user %v on %v stage, %v\n", j.Datas.Results[rs].Value[va], j.Datas.Results[rs].Metrics.User, j.Datas.Results[rs].Metrics.Stage, j.Datas.Results[rs].Metrics.State)

		}
	}
	return out
}

func getRawData(url string) ([]byte, error) {
	queryfqn := promethuesURL + url
	resp, getErr := http.Get(queryfqn)
	if getErr != nil {
		log.Println(getErr)
		return nil, getErr
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
	return logs, nil

}

func queryMakert(q metric) string {

	if q.TimeRange == "" {
		q.TimeRange = "1h"
	}
	if q.User == "" {
		q.User = ".*"
	}
	return fmt.Sprintf(`radius_response{stage=~"%s",state=~"%s",user=~"%s"}[%s]`, q.Stage, q.State, q.User, q.TimeRange)
}

func logs(w http.ResponseWriter, r *http.Request) {

	templ := templateHandler{filename: "logs.gohtml"}
	var q metric

	if r.Method == http.MethodPost {
		q.Stage = r.FormValue("stage")
		q.State = r.FormValue("state")
		q.User = r.FormValue("user")
		q.TimeRange = r.FormValue("timerange")
		log.Println(q)

		qResulat, qErr := getQuery(queryMakert(q))
		if qErr != nil {
			templ.options = qErr
		} else {

			templ.options = qResulat
		}

	}

	templ.ServeHTTP(w, r)

}
