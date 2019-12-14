package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

const filename = "data/11211002.DAT"

func main() {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	buf := make([]byte, 324) // 1秒あたり324Byte記録されている
	for {
		n, err := fp.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}

		enco := hex.EncodeToString(buf)
		fmt.Printf("%s\n", enco)
	}
}
