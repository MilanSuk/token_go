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
	name  string
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

func Vocab_getPath(name string) string {
	return name + ".tiktoken"
}

func NewVocab(name string, enableDownload bool) (*Vocab, error) {
	vb := &Vocab{name: name}

	if enableDownload {
		err := VocabAddr_findOrDownload(name)
		if err != nil {
			return nil, err
		}
	}

	//open file
	fl, err := os.ReadFile(Vocab_getPath(name))
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

func (node *VocabNode) encode(runes []int, rune_pos int, next_rune_pos *int, ids *[]int) bool {
	if rune_pos == len(runes) {
		return false
	}

	next, found := node.chs[runes[rune_pos]]
	if found {
		found = next.encode(runes, rune_pos+1, next_rune_pos, ids)
		if !found {

			//If next.id < 0 return to parent node
			//This is for situation when there is "go" and "good" in vocabulary and it searches for "goo ".
			//Tree is "g" -> "go" -> "goo" -> "good". "goo" node has 'id'=-1(it's NOT in vocabulary), other nodes are 'id'>= 0(in vocabulary).
			if next.id >= 0 {
				*next_rune_pos = rune_pos + 1
				*ids = append(*ids, next.id)
				return true
			}
		}
	}

	return found
}

func (vb *Vocab) Encode(str string) []int {
	ids := make([]int, 0, len(str)+1) //pre-alloc
	if str == "" {
		return ids
	}

	runes := getRunes(str)
	rune_pos := 0
	for rune_pos < len(runes) {
		vb.items.encode(runes, rune_pos, &rune_pos, &ids)
	}

	return ids
}
