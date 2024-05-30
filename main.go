/*
Copyright 2024 Milan Suk

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this db except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

func main() {
	//Profiling
	//Os_StartProfile("cpu.prof")	//run "sh perf" to show result
	//defer Os_StopProfile()

	TestServer() //run HTTP server and make request to encode/decode

	vb, err := NewVocab("p50k_base.tiktoken")
	//vb, err := NewVocab("cl100k_base.tiktoken")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//direct en/decode  test
	//ids := vb.Encode("kjl")
	//fmt.Println(ids)
	//str := vb.Decode(ids)
	//fmt.Println(str)

	fl, err := os.ReadFile("data.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	//encode
	st := ulit_getTime()
	tks := vb.Encode(string(fl))
	ulit_printStat("Encoded", st, len(tks), len(fl))

	//decode
	st = ulit_getTime()
	fl2 := vb.Decode(tks)
	ulit_printStat("Decoded", st, len(tks), len(fl))

	//compare decoded with original file
	if !bytes.Equal(fl, []byte(fl2)) {
		fmt.Println("Error")
		os.Exit(-1)
	}

	//TODO: test.txt ...
	//TODO: multi-threading ...
}

func ulit_getTime() int64 {
	return time.Now().UnixMicro() //micro seconds
}
func ulit_printStat(tp string, startTime_microsec int64, num_tokens, file_size int) {
	dt := float64(time.Now().UnixMicro()-startTime_microsec) / 1000000

	fmt.Printf("%s %d toks: %.3fM toks/sec, %.3f MB/sec\n", tp, num_tokens, float64(num_tokens)/dt/1000000, float64(file_size)/dt/(1000000))
}

/*
//Load .parquet file
type RowType struct{ Text string }
rows, err := parquet.ReadFile[RowType]("000_00000.parquet")
if err != nil {
	fmt.Println(err)
	return
}
vtokens := make([][]int, len(rows))
for _, r := range rows {
	tks := vb.Encode(r.Text)
	vtokens = append(vtokens, tks)
}
fmt.Println(len(vtokens))
*/
