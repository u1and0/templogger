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

// Datum : JSON要素
type Datum struct {
	Time        time.Time `json:"Time"`
	Temperature float64   `json:"Temperature"`
	Accx        []float64 `json:"AccelerationX"`
	Accy        []float64 `json:"AccelerationY"`
	Accz        []float64 `json:"AccelerationZ"`
}

type (
	encoded []byte
)

var data []*Datum

func main() {
	flag.Parse()
	for _, file := range flag.Args() {
		fp, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer fp.Close()

		buf := make(encoded, 324) // 1秒あたり324Byte記録されている
		for {
			n, err := fp.Read(buf)
			if n == 0 {
				break
			}
			if err != nil {
				panic(err)
			}

			enco := hex.EncodeToString(buf)

			// /* 日時変換 */
			tm, err := TransTime(enco)
			if err != nil {
				log.Fatalln(err)
			}

			/* 温度変換 */
			tmp, err := TransTemp(enco)
			if err != nil {
				log.Fatalln(err)
			}

			/* 加速度X */
			accx, err := TransAcc(enco, "x")
			if err != nil {
				log.Fatalln(err)
			}

			/* 加速度Y */
			accy, err := TransAcc(enco, "y")
			if err != nil {
				log.Fatalln(err)
			}

			/* 加速度Z */
			accz, err := TransAcc(enco, "z")
			if err != nil {
				log.Fatalln(err)
			}

			/* 出力 */
			d := &Datum{
				Time:        tm,
				Temperature: tmp,
				Accx:        accx,
				Accy:        accy,
				Accz:        accz}
			data = append(data, d)
		}
	}
	jdata, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", string(jdata))
}

// TransAcc : 一秒分324文字列から加速度Xの変換
func TransAcc(s, xyz string) ([]float64, error) {
	var (
		accu uint64 // バイトから読み取った文字列から変換したint加速度
		err  error
		acxl []float64 // 加速度配列
	)
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

// TransTime : 一秒分324文字列から日付・時間の変換
func TransTime(s string) (time.Time, error) {
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
func TransTemp(s string) (float64, error) {
	t, err := strconv.ParseInt(s[14:16]+s[12:14], 16, 0) // 16->10進数化
	tmp := -45 + 175*float64(t)/65535
	return tmp, err
}
