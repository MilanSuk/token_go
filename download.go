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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type VocabAddr struct {
	Name string
	Url  string
}

var g_vocab_addrs = []VocabAddr{
	{"r50k_base", "https://openaipublic.blob.core.windows.net/encodings/r50k_base.tiktoken"},
	{"p50k_base", "https://openaipublic.blob.core.windows.net/encodings/p50k_base.tiktoken"},
	{"p50k_edit", "https://openaipublic.blob.core.windows.net/encodings/p50k_base.tiktoken"},
	{"cl100k_base", "https://openaipublic.blob.core.windows.net/encodings/cl100k_base.tiktoken"},
	{"o200k_base", "https://openaipublic.blob.core.windows.net/encodings/o200k_base.tiktoken"},
}

func VocabAddr_findOrDownload(name string) error {
	//is file exist?
	filePath := Vocab_getPath(name)
	{
		info, err := os.Stat(filePath)
		if !os.IsNotExist(err) && !info.IsDir() {
			return nil //ok, file exist
		}
	}

	//find url
	var url string
	{
		for _, it := range g_vocab_addrs {
			if it.Name == name {
				url = it.Url
			}
		}
		if url == "" {
			return errors.New(name + " not found in addresses")
		}
	}

	//download
	{
		fmt.Printf("Downloading: '%s' from '%s' ...", name, url)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		os.WriteFile(filePath, data, 0644)
		fmt.Println("done")
	}

	return nil //ok
}
