package main

import (
	"fmt"
	"testing"
	"time"
)

// 先頭から12バイトを時間に変換
func TestTransTime(t *testing.T) {
	e := Encoded{String: "1119132336336e68f8d8101e9ffd000"}
	actual, err := e.TransTime()
	expected := time.Date(2019, 11, 23, 13, 33, 36, 0, time.Local)
	/*
	   1119 =2019年11月
	   1323 = 23日13時
	   3633 = 33分36秒
	*/
	if err != nil {
		fmt.Println(err)
	}
	if actual != expected {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
}

// 12バイト目から4バイトを温度に変換
func TestTransTemp(t *testing.T) {
	e := Encoded{String: "1119132336336e68f8d8101e9ffd000"}
	actual, err := e.TransTemp()
	expected := 26.38857099259937
	/* 0x686e = 26.38℃ */
	if err != nil {
		fmt.Println(err)
	}
	if actual != expected {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
}

func TestAdd(t *testing.T) {
	a1 := &Datum{Temperature: 10}
	a2 := &Datum{Temperature: 20}
	d := Data{a1}
	actual := d.Add(a2)
	expected := []Datum{Datum{
		Temperature: 10},
		Datum{Temperature: 20},
	}
	for i, e := range expected {
		if actual[i].Temperature != e.Temperature {
			t.Fatalf("got: %v want: %v", actual[i].Temperature, e.Temperature)
		}
	}
}
