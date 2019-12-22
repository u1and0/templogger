package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Datum : 1秒当たりのデータ
type Datum struct {
	Time  time.Time `json:"Time"`
	Temp  float64   `json:"Temperature"`
	Hum   float64   `json:"Humidity"`
	Atemp float64   `json:"TemperatureAcc"`
	Accx  []float64 `json:"AccelerationX"`
	Accy  []float64 `json:"AccelerationY"`
	Accz  []float64 `json:"AccelerationZ"`
}

// Data : 読み込んだファイル内のデータすべてを入れるスライス
type Data []*Datum

// Encoded : バイナリファイルから読みだした16進数
type Encoded struct {
	Bytes  []byte
	String string
}

func main() {
	data := Data{}
	flag.Parse()
	for _, file := range flag.Args() {
		var e Encoded
		fp, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer fp.Close()

		e.Bytes = make([]byte, 324) // 1秒あたり324Byte記録されている
		for {
			n, err := fp.Read(e.Bytes)
			if n == 0 {
				break
			}
			if err != nil {
				panic(err)
			}

			e.String = hex.EncodeToString(e.Bytes)
			/* 日時 */
			tm, err := e.TransTime()
			if err != nil {
				log.Fatalln(err)
			}
			/* 温度 */
			tmp, err := e.TransTemp()
			if err != nil {
				log.Fatalln(err)
			}
			/* 湿度 */
			hum, err := e.TransHum()
			if err != nil {
				log.Fatalln(err)
			}
			/* 加速度センサー内の温度 */
			atmp, err := e.TransAtemp()
			if err != nil {
				log.Fatalln(err)
			}
			/* 加速度 */
			accx, err := e.TransAcc("x")
			if err != nil {
				log.Fatalln(err)
			}
			accy, err := e.TransAcc("y")
			if err != nil {
				log.Fatalln(err)
			}
			accz, err := e.TransAcc("z")
			if err != nil {
				log.Fatalln(err)
			}
			/* 出力 */
			d := &Datum{
				Time:  tm,
				Temp:  tmp,
				Hum:   hum,
				Atemp: atmp,
				Accx:  accx,
				Accy:  accy,
				Accz:  accz,
			}
			data = data.Append(d)
		}
	}
	jdata, err := data.Jsonize("\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", string(jdata))
}

// TransTime : 一秒分324文字列から日付・時間の変換
func (e Encoded) TransTime() (time.Time, error) {
	s := e.String
	y, err := strconv.Atoi(s[2:4])
	m, err := strconv.Atoi(s[0:2])
	mm := time.Month(m) // 月のみMonth型をDate()関数に渡さなければならない
	d, err := strconv.Atoi(s[6:8])
	H, err := strconv.Atoi(s[4:6])
	M, err := strconv.Atoi(s[10:12])
	S, err := strconv.Atoi(s[8:10])
	tm := time.Date(2000+y, mm, d, H, M, S, 0, time.Local)
	return tm, err
}

// TransTemp : 一秒分324文字列から温度の変換
func (e Encoded) TransTemp() (float64, error) {
	s := e.String
	t, err := strconv.ParseInt(s[14:16]+s[12:14], 16, 0) // 16->10進数化
	tmp := -45 + 175*float64(t)/65535
	return tmp, err
}

// TransHum : 一秒分324文字列から湿度の変換
func (e Encoded) TransHum() (float64, error) {
	s := e.String
	t, err := strconv.ParseInt(s[18:20]+s[16:18], 16, 0) // 16->10進数化
	hum := 100 * float64(t) / 65535
	return hum, err
}

// TransAtemp : 一秒分324文字列から加速度付属の温度センサーの変換
func (e Encoded) TransAtemp() (float64, error) {
	s := e.String
	a, err := strconv.ParseUint(s[22:24]+s[20:22], 16, 0) // 16->10進数化
	if err != nil {
		fmt.Println(err)
	}
	atmp := float64(int16(a))/333.87 + 21
	return atmp, err
}

// TransAcc : 一秒分324文字列から加速度の変換
func (e Encoded) TransAcc(xyz string) ([]float64, error) {
	var (
		accu uint64 // バイトから読み取った文字列から変換したint加速度
		err  error
		acxl []float64 // 加速度配列
	)
	s := e.String
	axis := map[string]int{"x": 48, "y": 52, "z": 56}
	for i := axis[xyz]; i < len(s); i += 12 { // 初期バイトはxyzの方向による
		ss := s[i+2:i+4] + s[i:i+2]
		accu, err = strconv.ParseUint(ss, 16, 0) // 16->10進数化
		if err != nil {
			fmt.Println(err)
		}
		acci := int16(accu) // 一気にfloatに渡してはいけない
		// なぜなら、Uintの補数をintとして換算しないと小数点計算されてしまう
		acxl = append(acxl, 16000*float64(acci)/32768) // 換算加速度を加速度配列に格納
	}
	return acxl, err
}

// Append :append data slice
func (d Data) Append(a *Datum) Data {
	return append(d, a)
}

// Jsonize : json marshal
func (d Data) Jsonize(s string) ([]byte, error) {
	return json.MarshalIndent(d, "", s)
}
