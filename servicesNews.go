package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

//Реализация интерфейса NewsManager, работает с сайтом https://finance.yahoo.com
type NewsManagerYahoo struct {
}

//Метод структуры NewsManagerYahoo, принимает символ финансового актива, возвращает список экземпляров структуры News
func (newsManager NewsManagerYahoo) GetNews(symbol string) ([]News, error) {
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/news?p=%s", symbol, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	news := []News{}
	html.Find("a").Each(func(i int, selection *goquery.Selection) {
		isBox := selection.Children().HasClass("StretchedBox")
		if isBox == true {
			headline := selection.Text()
			link, _ := selection.Attr("href")
			if strings.Index(link, "http") == -1 {
				link = "https://finance.yahoo.com" + link
			}
			news = append(news, News{headline, link})
		}
	})
	if len(news) == 0 {
		err = errors.New("emptyNews")
		return nil, err
	}
	return news, err
}
