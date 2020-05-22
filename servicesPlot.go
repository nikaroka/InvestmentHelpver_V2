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
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

type PlotManagerAlphaVantage struct {
	apiKey string
}

func (plotManager PlotManagerAlphaVantage) GetPlot(symbol string) ([]Day, error) {
	apiKey := plotManager.apiKey
	body, err := GetPlotJson(symbol, apiKey)
	if err != nil {
		return nil, err
	}
	plot, err := ScrapJsonBody(body)
	if err != nil {
		return nil, err
	}
	return plot, nil
}

func GetPlotJson(symbol string, key string) (string, error) {
	req := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", symbol, key)
	resp, err := http.Get(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if strings.Index(string(body), "Note") != -1 {
		err = errors.New("exceedApiFrequency")
		return "", err
	}

	if strings.Index(string(body), "Invalid API call") != -1 {
		err = errors.New("wrongSymbolApiCall")
		return "", err
	}
	return string(body), err
}

func ScrapJsonBody(body string) ([]Day, error) {
	byt := []byte(body)
	var days []Day
	var dat map[string]interface{}
	err := json.Unmarshal(byt, &dat)
	if err != nil {
		return nil, err
	}
	dailyTimeSeries := dat["Time Series (Daily)"].(map[string]interface{})
	for key := range dailyTimeSeries {
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
	sort.Slice(days, func(i, j int) bool { return days[i].Date.Before(days[j].Date) })
	return days, nil
}

func GetFloatValue(valueI interface{}) float64 {
	valueF, _ := strconv.ParseFloat(valueI.(string), 64)
	return valueF
}
