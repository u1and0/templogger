package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"
)

const filename = "data/11211002.DAT"

// type Row struct {
// 	Key uint16
// 	Val uint16
// }

// func (r Row) String() string {
// 	return fmt.Sprintf("(%s: %v)", string(r.Key), r.Val)
// }

func main() {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	buf := make([]byte, 324) // 1秒あたり324Byte記録されている
	// row := Row{}

	for {
		n, err := fp.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}

		// var i uint32

		// binary.LittleEndian.PutUint32(buf, i)
		// fmt.Printf("%x", binary.BigEndian.Uint32(buf)) // タイムスタンプしか出力されない
		// binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &i)

		enco := hex.EncodeToString(buf)

		var y, m, d, H, M, S int
		y, err = strconv.Atoi(enco[2:4])
		m, err = strconv.Atoi(enco[0:2])
		mm := time.Month(m) // 月のみMonth型をDate()関数に渡さなければならない
		d, err = strconv.Atoi(enco[6:8])
		H, err = strconv.Atoi(enco[4:6])
		M, err = strconv.Atoi(enco[10:12])
		S, err = strconv.Atoi(enco[8:10])
		if err != nil {
			fmt.Println(err)
		}
		tm := time.Date(2000+y, mm, d, H, M, S, 0, time.Local)
		// fmt.Printf("%d%d%d\n", y, m, d)

		// fmt.Printf("%s\n", enco)
		fmt.Printf("%v\n", tm)
	}
}
