package main

import (
	"strings"
	"testing"
)

func TestGetDailyDataShort (t *testing.T)  {
	apiKey := loadConfig().VentageKey
	body, err := GetDailyDataShort("IBM", apiKey)
	if strings.Index(body, "Invalid API call") != -1 {
		t.Error("zero len body")
	}

	if err != nil {
		t.Error(err)
	}

	body, err = GetDailyDataShort("testKey", apiKey)
	if strings.Index(body, "Invalid API call") == -1 {
		t.Error("find testKey body")
	}
}

func TestScrapJsonBody (t *testing.T)  {
	apiKey := loadConfig().VentageKey
	body, _ := GetDailyDataShort("IBM", apiKey)
	days, err := ScrapJsonBody(body)
	if err != nil {
		t.Error(err)
	}

	if days == nil {
		t.Error("nil days")
	}
}


