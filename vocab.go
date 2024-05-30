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
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type VocabNode struct {
	parent *VocabNode
	id     int
	chs    map[int]*VocabNode
}

func (node *VocabNode) Add(str []int, id int) {
	if node.chs == nil {
		node.chs = make(map[int]*VocabNode)
	}

	c, found := node.chs[str[0]]
	if !found {
		c = &VocabNode{id: -1, parent: node}
		node.chs[str[0]] = c
	}

	if len(str) == 1 {
		c.id = id
	} else {
		c.Add(str[1:], id)
	}
}

type Vocab struct {
	path  string
	words map[string]int //word-id pair
	ids   []string       //word[id]
	items VocabNode      //unicode(rune) tree
}

func getRunes(str string) []int {
	ids := make([]int, 0, len(str)) //pre-alloc
	for _, ch := range str {
		ids = append(ids, int(ch))
	}
	return ids //list of unicodes
}

func NewVocab(path string) (*Vocab, error) {
	vb := &Vocab{path: path}

	//open file
	fl, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	//decode file
	vb.words = make(map[string]int)
	lines := strings.Split(string(fl), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		token, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, err
		}
		rank, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		vb.words[string(token)] = rank
	}

	//create ids
	{
		max_id := 0
		for _, id := range vb.words {
			if id > max_id {
				max_id = id
			}
		}
		vb.ids = make([]string, max_id+1)
		for word, id := range vb.words {
			vb.ids[id] = word
		}
	}

	//create tree
	{
		vb.items.id = -1
		for word, id := range vb.words {

			str := getRunes(word)
			vb.items.Add(str, id)
		}
	}

	return vb, nil
}

func (vb *Vocab) Decode(ids []int) string {
	var buffer bytes.Buffer
	for _, id := range ids {
		buffer.WriteString(vb.ids[id]) //faster than +=
	}
	return buffer.String()
}

func (vb *Vocab) Encode(str string) []int {
	ids := make([]int, 0, len(str)) //pre-alloc

	runes := getRunes(str)

	items := &vb.items
	for i := 0; i < len(runes); i++ {
		next_item, found := items.chs[runes[i]]
		if found {
			items = next_item
		} else {

			//Go back in tree.
			//This is for situation when there is "go" and "good" in vocabulary and it search for "goo ".
			//Tree is "g" -> "go" -> "goo" -> "good". "goo" tree leaf has 'id'=-1(it's NOT in vocabulary), others are >= 0.
			for items.id < 0 && items.parent != nil {
				items = items.parent
				i--
			}

			if items.id >= 0 {
				i--
			} else {
				fmt.Println("warning: tree ID not exist")
			}

			ids = append(ids, items.id) //add

			items = &vb.items //reset
		}
	}

	//add last
	if items != &vb.items {
		ids = append(ids, items.id) //add
	}

	return ids
}
