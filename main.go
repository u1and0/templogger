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

const (
	// VERSION : version
	VERSION = "0.2.1"
)

var (
	// バージョンフラグ
	showVersion bool
	// ヘルプフラグ
	showHelp bool
	// インデントフラグ
	indent bool
	// 最終的にPrintするJSONデータ
	jdata []byte
	// エラー
	err error
)

// Datum : 1秒当たりのデータ
type Datum struct {
	Time  time.Time `json:"Time"`
	Temp  float64   `json:"Temperature"`
	Hum   float64   `json:"Humidity"`
	Atemp float64   `json:"TemperatureAcc"`
	Gyrox float64   `json:"GyroX"`
	Gyroy float64   `json:"GyroY"`
	Gyroz float64   `json:"GyroZ"`
	Compx float64   `json:"CompassX"`
	Compy float64   `json:"CompassY"`
	Compz float64   `json:"CompassZ"`
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

func flagUsage() {
	usageText := `SDカードにためたバイナリデータをテキスト(JSON形式)にして標準出力にdumpします。

Usage:
単一のファイルをJSON化
	templogger data/12161037.DAT
複数のファイルをJSON化
	templogger data/12161037.DAT data/12161237.DAT
すべてのDATファイルをJSON化
	templogger data/*.DAT
-tオプションで読みやすいようにインデントを入れます
	templogger -t data/*.DAT

-h, -help		show help message
-t, -indent		indent to format output
-v, -version	show version
`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func main() {
	flag.Usage = flagUsage
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&indent, "t", false, "indent to format output")
	flag.BoolVar(&indent, "indent", false, "indent to format output")
	data := Data{}
	flag.Parse()

	if showVersion {
		fmt.Println("templogger version:", VERSION)
		return // versionを表示して終了
	}

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
			/* ジャイロ */
			gyrox, err := e.TransGyro("x")
			if err != nil {
				log.Fatalln(err)
			}
			gyroy, err := e.TransGyro("y")
			if err != nil {
				log.Fatalln(err)
			}
			gyroz, err := e.TransGyro("z")
			if err != nil {
				log.Fatalln(err)
			}
			/* コンパス */
			compx, err := e.TransCompass("x")
			if err != nil {
				log.Fatalln(err)
			}
			compy, err := e.TransCompass("y")
			if err != nil {
				log.Fatalln(err)
			}
			compz, err := e.TransCompass("z")
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
				Gyrox: gyrox,
				Gyroy: gyroy,
				Gyroz: gyroz,
				Compx: compx,
				Compy: compy,
				Compz: compz,
				Accx:  accx,
				Accy:  accy,
				Accz:  accz,
			}
			data = append(data, d)
		}
	}
	if indent {
		jdata, err = json.MarshalIndent(data, "", "\t")
	} else {
		jdata, err = json.Marshal(data)
	}
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

// TransGyro : 一秒分324文字列からジャイロセンサーの変換
func (e Encoded) TransGyro(xyz string) (float64, error) {
	s := e.String
	i := map[string]int{"x": 24, "y": 28, "z": 32}[xyz]
	a, err := strconv.ParseUint(s[i+2:i+4]+s[i:i+2], 16, 0) // 16->10進数化
	if err != nil {
		fmt.Println(err)
	}
	g := 250 * float64(int16(a)) / 32768
	return g, err
}

// TransCompass : 一秒分324文字列からコンパスセンサーの変換
func (e Encoded) TransCompass(xyz string) (float64, error) {
	s := e.String
	i := map[string]int{"x": 36, "y": 40, "z": 44}[xyz]
	a, err := strconv.ParseUint(s[i+2:i+4]+s[i:i+2], 16, 0) // 16->10進数化
	if err != nil {
		fmt.Println(err)
	}
	g := 4800 * float64(int16(a)) / 32768
	return g, err
}

// TransAcc : 一秒分324文字列から加速度の変換
func (e Encoded) TransAcc(xyz string) ([]float64, error) {
	s := e.String
	x := map[string]int{"x": 48, "y": 52, "z": 56}
	var (
		err  error
		acxl []float64 // 加速度配列
	)
	for i := x[xyz]; i < len(s); i += 12 { // 初期バイトはxyzの方向による
		a, err := strconv.ParseUint(s[i+2:i+4]+s[i:i+2], 16, 0) // 16->10進数化
		if err != nil {
			fmt.Println(err)
		}
		// 一気にfloatに渡してはいけない
		// なぜなら、Uintの補数をintとして換算しないと小数点計算されてしまう
		acxl = append(acxl, 16000*float64(int16(a))/32768) // 換算加速度を加速度配列に格納
	}
	return acxl, err
}
