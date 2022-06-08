package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
)

type Event struct {
	User     string  `json:"user"`
	URLPath  string  `json:"resource"`
	Duration float64 `json:"duration"`
	Method   string  `json:"method"`
}

type Http struct {
	IpAddress string `json:"clientHost"`
}

type Record struct {
	Event `json:"json"`
	Http  `json:"http"`
	Time  int64
}

func (r *Record) regexPath() string {
	_, err := regexp.MatchString("^\\/\\w+", r.URLPath)
	if err != nil {
		return ""
	}

	re, _ := regexp.Compile("^\\/\\w+")

	return re.FindString(r.URLPath) + "-" + r.Method
}

var records []Record
var users []string
var ips []string
var paths []string

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: Please type json file path")
		return
	}

	path := args[1]
	err := fetchData(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	initData()
	fmt.Println("group by users")
	for _, v := range users {
		analysis(v)
	}
	fmt.Println("group by ips")
	for _, v := range ips {
		analysis(v)
	}
	fmt.Println("group by urls")
	for _, v := range paths {
		analysis(v)
	}
}

func fetchData(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	json.Unmarshal(byteData, &result)

	data, err := json.Marshal(result["events"])
	if err != nil {
		return err
	}

	var events []map[string]interface{}
	json.Unmarshal(data, &events)

	for _, v := range events {
		data, _ := json.Marshal(v["event"])
		var record Record
		json.Unmarshal(data, &record)
		record.Time = int64(v["timestamp"].(float64))
		records = append(records, record)
	}

	return nil
}

func initData() {
	tempUsers := make(map[string]bool)
	tempIPs := make(map[string]bool)
	tempPaths := make(map[string]bool)
	for _, v := range records {
		_, ok := tempUsers[v.User]
		if !ok {
			tempUsers[v.User] = true
			users = append(users, v.User)
		}

		_, ok = tempIPs[v.IpAddress]
		if !ok {
			tempIPs[v.IpAddress] = true
			ips = append(ips, v.IpAddress)
		}

		_, ok = tempPaths[v.regexPath()]
		if !ok {
			tempPaths[v.regexPath()] = true
			paths = append(paths, v.regexPath())
		}

	}
}

func analysis(key string) {
	var sum float64
	var count int
	var max float64
	var min float64
	var median float64
	temp := make([]float64, 0)
	for _, v := range records {

		if v.User == key || v.IpAddress == key || v.regexPath() == key {
			sum += v.Duration
			count++
			if v.Duration > max {
				max = v.Duration
			}

			if v.Duration < min || min == 0 {
				min = v.Duration
			}
			temp = append(temp, v.Duration)
		}
	}
	sort.Float64s(temp)
	if len(temp)%2 == 0 {
		median = (temp[(len(temp)/2)-1] + temp[(len(temp)/2)]) / 2
	} else {
		median = temp[(len(temp) / 2)]
	}

	fmt.Printf("%v  response time { avg: %.2f median: %.2f max: %.2f min %.2f } api count: %d\n", key, sum/float64(count), median, max, min, count)
}
