package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Day struct {
	Date time.Time
	Open float64
	High float64
	Low float64
	Close float64
	Volume int
}

func GetDailyDataShort (symbol string, key string) (string, error) {
	req := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s",symbol,key)
	resp, err := http.Get(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if strings.Index(string(body), "Note") != -1{
		err = errors.New("exceedApiFrequency")
		return "", err
	}
	return string(body), err
}

func GetFloatValue (valueI interface{}) float64 {
	valueF, _ := strconv.ParseFloat(valueI.(string), 64)
	return valueF
}

func ScrapJsonBody (body string) ([]Day, error){
	byt := []byte(body)
	var days []Day
	var dat map[string]interface{}
	err := json.Unmarshal(byt, &dat)
	if err != nil {
		return nil, err
	}
	dailyTimeSeries := dat["Time Series (Daily)"].(map[string]interface{})
	for key := range dailyTimeSeries{
		date, err := time.Parse("2006-01-02", key)
		if err != nil {
			return nil, err
		}
		dayValues := dailyTimeSeries[key].(map[string]interface{})
		v, err := strconv.Atoi(dayValues["5. volume"].(string))
		if err != nil {
			return nil, err
		}
		day := Day{
			Date:   date,
			Open:   GetFloatValue(dayValues["1. open"]),
			High:   GetFloatValue(dayValues["2. high"]),
			Low:    GetFloatValue(dayValues["3. low"]),
			Close:  GetFloatValue(dayValues["4. close"]),
			Volume: v,
		}
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].Date.Before(days[j].Date )})
	return days, nil
	}

func requestHandlerPlot(r *http.Request, plotChannel chan []byte, errChannel chan error) {
	symbol := strings.Split(r.URL.Path[1:], ";")[0]
	apiKey := loadConfig().VentageKey
	body, err := GetDailyDataShort(symbol, apiKey)
	if err != nil {
		errChannel <- err
		return
	}
	plot, err := ScrapJsonBody(body)
	if err != nil {
		errChannel <- err
		return
	}
	jsonPlot, err := json.Marshal(plot)
	if err != nil {
		errChannel <- err
		return
	}
	plotChannel <- jsonPlot
}
