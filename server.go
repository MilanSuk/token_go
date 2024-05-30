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
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type Server struct {
	mu     sync.Mutex
	vocabs []*Vocab

	//collect stats ...
}

func (s *Server) findVocab(path string) *Vocab {
	for _, vb := range s.vocabs {
		if vb.path == path {
			return vb
		}
	}
	return nil
}

func (s *Server) GetVocab(path string) *Vocab {
	s.mu.Lock()
	defer s.mu.Unlock()

	vb := s.findVocab(path)
	if vb == nil {
		//add
		var err error
		vb, err = NewVocab(path)
		if err == nil {
			s.vocabs = append(s.vocabs, vb)
		}
	}
	return vb
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	var vocabName string
	var encode, decode bool

	//Encoder
	vocabName, encode = strings.CutPrefix(r.URL.Path, "/encode/")
	if encode {
		vocab := s.GetVocab(vocabName + ".tiktoken")
		if vocab == nil {
			http.Error(w, vocabName+" vocab not found", http.StatusBadRequest)
			return
		}

		//get string
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "body read failed", http.StatusBadRequest)
			return
		}

		ids := vocab.Encode(string(body))
		//fmt.Println("Encode:", ids)
		wids := ulit_integers_to_bytes(ids)
		w.Write(wids)
		return
	}

	//Decoder
	vocabName, decode = strings.CutPrefix(r.URL.Path, "/decode/")
	if decode {
		vocab := s.GetVocab(vocabName + ".tiktoken")
		if vocab == nil {
			http.Error(w, vocabName+" vocab not found", http.StatusBadRequest)
			return
		}

		//get string
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "body read failed", http.StatusBadRequest)
			return
		}

		ids := ulit_bytes_to_integers(body)
		str := vocab.Decode(ids)
		//fmt.Println("Decode:", str)
		w.Write([]byte(str))
		return
	}

	http.Error(w, r.URL.Path+" parsing failed", http.StatusBadRequest)

	//http.ServeFile(w, r, "index.html")
}

func NewServer() (*Server, error) {
	server := &Server{}

	http.HandleFunc("/", server.Handle)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func Client_encode(str []byte, vocab string) ([]int, error) {

	client := http.DefaultClient
	res, err := client.Post("http://localhost:8090/encode/"+vocab, "text/plain", bytes.NewBuffer(str))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if len(body)%4 != 0 {
		return nil, fmt.Errorf("invalid size of return array")
	}

	ids := ulit_bytes_to_integers(body)
	return ids, nil //ok
}

func Client_decode(ids []int, vocab string) ([]byte, error) {

	wids := ulit_integers_to_bytes(ids)

	client := http.DefaultClient
	res, err := client.Post("http://localhost:8090/decode/"+vocab, "text/plain", bytes.NewBuffer(wids))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ulit_integers_to_bytes(ids []int) []byte {
	data := make([]byte, len(ids)*4)
	for i, id := range ids {
		binary.LittleEndian.PutUint32(data[i*4:(i*4)+4], uint32(id))
	}
	return data
}
func ulit_bytes_to_integers(data []byte) []int {
	ids := make([]int, len(data)/4)
	for i := 0; i < len(data); i += 4 {
		ids[i/4] = int(binary.LittleEndian.Uint32(data[i : i+4]))
	}
	return ids
}

func TestServer() {

	//run server
	go NewServer()

	//constants
	vocab := "p50k_base"
	str := "Hi there!"

	//decode
	ids, err := Client_encode([]byte(str), vocab)
	if err != nil {
		fmt.Println(err)
		return
	}

	//encode
	str2, err := Client_decode(ids, vocab)
	if err != nil {
		fmt.Println(err)
		return
	}

	//print results
	fmt.Println(str)
	fmt.Println(ids)
	fmt.Println(string(str2))
}
