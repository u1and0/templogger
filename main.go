package main

import (
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	// VERSION : version
	VERSION = "0.2.1r"
)

var (
	// バージョンフラグ
	showVersion bool
	// ヘルプフラグ
	showHelp bool
	// 出力ファイル形式
	dumpFormat string
	// インデントフラグ
	indent bool
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
	templogger -f json data/12161037.DAT
複数のファイルをJSON化
	templogger --format json data/12161037.DAT data/12161237.DAT
すべてのDATファイルをJSON化
	templogger --format json data/*.DAT
-tオプションで読みやすいようにインデントを入れます
	templogger --format json -t data/*.DAT

-f, -format		dump format "csv" or "json"
-h, -help		show help message
-t, -indent		indent to format output (must use with "--format json")
-v, -version		show version
`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func main() {
	flag.Usage = flagUsage
	flag.StringVar(&dumpFormat, "f", "csv", "dump format csv or json")
	flag.StringVar(&dumpFormat, "format", "csv", "dump format csv or json")
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

	// File walk thru & Change format binary to text
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
			// Compile data
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
			data = data.Append(d)
			/* 行ごとではなく一度にdumpする理由
			JSON Array形式で出力するので、
			Arrayの接続の,とか終わりの]とか出力するのが難しいので、
			goのオブジェクト上でSliceにしてそれをJson Marshalかけるのが楽。
			いずれにせよcsv出力するときはappendしてSliceオブジェクトにするので、
			1行ずつ出力することにパフォーマンスの改善もない。
			*/
		}
	}

	// Output
	switch dumpFormat {
	case "csv": // dump to a file
		if err := data.ToCSV(); err != nil {
			log.Fatalf("%s", err)
		}
	case "json": // dump to stdout
		out, err := data.ToJSON(indent)
		if err != nil {
			log.Fatalf("%s", err)
		}
		fmt.Printf("%v\n", string(out))
	default:
		err := errors.New("error unknown file type")
		log.Fatalf("%s", err)
	}
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

// TrimExtension : trim file extensiton
func TrimExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

// Append :append data slice
func (d Data) Append(a *Datum) Data {
	return append(d, a)
}

// ToJSON : convert data slice as JSON format
func (d Data) ToJSON(indent bool) (b []byte, err error) {
	if indent {
		return json.MarshalIndent(d, "", "\t")
	}
	return json.Marshal(d)
}

// ToCSV : convert data slice as CSV format
// create a csv file
// file name is a first argument of dat file
func (d Data) ToCSV() error {
	filename := TrimExtension(flag.Args()[0]) + ".csv"
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	// w := csv.NewWriter(f) // as UTF-8
	// Create japanese sjis csv file
	w := csv.NewWriter(transform.NewWriter(f, japanese.ShiftJIS.NewEncoder())) // as SJIS

	layout := "2006-01-02 15:04:05" // time format
	// csv header
	w.Write([]string{
		"時間", "温度", "湿度",
		"ジャイロX", "ジャイロY", "ジャイロZ",
		"コンパスX", "コンパスY", "コンパスZ",
		"加速度X", "加速度Y", "加速度Z",
	})
	for _, datum := range d { // csv items
		var record []string
		record = append(record,
			datum.Time.Format(layout),
			fmt.Sprintf("%0.4f", datum.Temp),
			fmt.Sprintf("%0.4f", datum.Hum),
			fmt.Sprintf("%0.4f", datum.Gyrox),
			fmt.Sprintf("%0.4f", datum.Gyroy),
			fmt.Sprintf("%0.4f", datum.Gyroz),
			fmt.Sprintf("%0.4f", datum.Compx),
			fmt.Sprintf("%0.4f", datum.Compy),
			fmt.Sprintf("%0.4f", datum.Compz),
			fmt.Sprintf("%0.4f", datum.Accx[len(datum.Accx)-1]),
			fmt.Sprintf("%0.4f", datum.Accy[len(datum.Accy)-1]),
			fmt.Sprintf("%0.4f", datum.Accz[len(datum.Accz)-1]),
		)
		w.Write(record)
	}
	w.Flush()
	return err
}
