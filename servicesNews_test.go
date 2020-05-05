package main

import (
	"testing"
)

func TestGetNewsYahoo(t *testing.T)  {
	news, err := getNewsYahoo("IBM")
	if err != nil{
		t.Error(err)
	}
	if len(news) == 0 {
		t.Error("empty news")
	}
	news, err = getNewsYahoo("testShit")
	if err == nil{
		t.Error(err)
	}
}
