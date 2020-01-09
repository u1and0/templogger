package main

import (
	"fmt"
	"testing"
	"time"
)

// 先頭から12ビットを時間に変換
func TestTransTime(t *testing.T) {
	e := Encoded{String: "1119132336330000000000000000000"}
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

// 12ビット目から8ビットを温度に変換
func TestTransTemp(t *testing.T) {
	e := Encoded{String: "0000000000006e68000000000000000"}
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

// 16ビット目から8ビットを温度に変換
func TestTransHum(t *testing.T) {
	e := Encoded{String: "0000000000000000b78c00000000000"}
	actual, err := e.TransHum()
	expected := 54.96757457846952
	/* 0x686e = 54.9676% */
	if err != nil {
		fmt.Println(err)
	}
	if actual != expected {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
}

// 20ビット目から8ビットを温度に変換
func TestTransAtemp(t *testing.T) {
	e := Encoded{String: "0000000000000000000c10100000000"}
	actual, err := e.TransAtemp()
	expected := 33.31617096474676
	/* 0x01c1 = 33.3162℃ */
	if err != nil {
		fmt.Println(err)
	}
	if actual != expected {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
}

// 24ビット目から8ビットをジャイロに変換
func TestTransGyro(t *testing.T) {
	e := Encoded{String: "00000000000000000000000ff770"}
	actual, err := e.TransGyro("x")
	expected := 220.63446044921875
	/* 0x01c1 = 220.63*/
	if err != nil {
		fmt.Println(err)
	}
	if actual != expected {
		t.Fatalf("got: %v want: %v", actual, expected)
	}
}

func TestAppend(t *testing.T) {
	a1 := &Datum{Temp: 10}
	a2 := &Datum{Temp: 20}
	d := Data{a1}
	actual := d.Append(a2)
	expected := []Datum{Datum{
		Temp: 10},
		Datum{Temp: 20},
	}
	for i, e := range expected {
		if actual[i].Temp != e.Temp {
			t.Fatalf("got: %v want: %v", actual[i].Temp, e.Temp)
		}
	}
}
