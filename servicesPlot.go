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

//Реализация интерфейса PlotManager, имеет параметр apiKey являющийся клюом к API Alpha Ventage
type PlotManagerAlphaVantage struct {
	apiKey string
}

//Метод структуры PlotManagerAlphaVantage, принимает символ финансового актива, возвращает список экземпляров структуры Candle
func (plotManager PlotManagerAlphaVantage) GetPlot(symbol string) ([]Candle, error) {
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

//Метод принимающий символ финансового актива и ключ API, производит запроса на Alpha Ventage и возвращает тело ответа
func GetPlotJson(symbol string, apiKey string) (string, error) {
	req := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", symbol, apiKey)
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

//Метод принимающий в себя тело ответа из функции GetPlotJson и превращающий его в список структуры Candle
func ScrapJsonBody(body string) ([]Candle, error) {
	byt := []byte(body)
	var days []Candle
	var dat map[string]interface{}
	err := json.Unmarshal(byt, &dat)
	if err != nil {
		return nil, err
	}
	dailyTimeSeries := dat["Time Series (Daily)"].(map[string]interface{})
	for dateKey := range dailyTimeSeries {
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			return nil, err
		}
		dayValues := dailyTimeSeries[dateKey].(map[string]interface{})
		value, err := strconv.Atoi(dayValues["5. volume"].(string))
		if err != nil {
			return nil, err
		}
		prices, err := GetFloatPrices(dayValues["1. open"].(string), dayValues["2. high"].(string),
			dayValues["3. low"].(string), dayValues["4. close"].(string))
		if err != nil {
			return nil, err
		}
		day := Candle{
			Date:   date,
			Open:   prices[0],
			High:   prices[1],
			Low:    prices[2],
			Close:  prices[3],
			Volume: value,
		}
		days = append(days, day)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].Date.Before(days[j].Date) })
	return days, nil
}

//Вспомогательный метод превращающий цены в формате String в цены в формате Float64
func GetFloatPrices(openS string, highS string, lowS string, closeS string) ([]float64, error) {
	openF, err := strconv.ParseFloat(openS, 64)
	if err != nil {
		return nil, err
	}
	highF, err := strconv.ParseFloat(highS, 64)
	if err != nil {
		return nil, err
	}
	lowF, err := strconv.ParseFloat(lowS, 64)
	if err != nil {
		return nil, err
	}
	closeF, err := strconv.ParseFloat(closeS, 64)
	if err != nil {
		return nil, err
	}
	preces := []float64{openF, highF, lowF, closeF}
	return preces, nil
}
