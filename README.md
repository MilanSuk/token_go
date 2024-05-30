## Token_go
Simple & fast Encoder/Decoder for tiktoken vocabulary.
Implemented from scratch(no regex library). Tokenizer is in vocab.go which has ~110 lines of code.



## Performance
p50k_base.tiktoken:
- Encoder: 4.458M toks/sec, 18.451 MB/sec
- Decoder: 37.817M toks/sec, 156.516 MB/sec

cl100k_base.tiktoken:
- Encoded 3.979M toks/sec, 16.875 MB/sec
- Decoded 35.825M toks/sec, 151.952 MB/sec

*note: put megabytes of text into data.txt.*



## Examples
Encode/Decode:
<pre><code>vb, err := NewVocab("p50k_base.tiktoken")

toks := vb.Encode("Hi there!")
fmt.Println(toks)

str := vb.Decode(toks)
fmt.Println(str)
</code></pre>

Client/Server:
<pre><code>go NewServer()   //run server in extra thread

toks, err := Client_encode([]byte("Hi there!"), "cl100k_base")
fmt.Println(toks)

text, err := Client_decode([]int{13347, 1070, 0}, "cl100k_base")
fmt.Println(text)
</code></pre>



## Build
Written in Go language(https://go.dev/doc/install). No dependencies.

<pre><code>git clone https://github.com/milansuk/token_go
cd token_go
go build
./token_go
</code></pre>



## Author
Milan Suk

Email: milan@skyalt.com

Twitter: https://twitter.com/milansuk/

**Sponsor**: https://github.com/sponsors/MilanSuk

*Feel free to follow or contact me with any idea, question or problem.*



## Contributing
Your feedback and code are welcome!

For bug report or question, please use [GitHub's Issues](https://github.com/skyaltlabs/skyalt/issues)

SkyAlt is licensed under **Apache v2.0** license. This repository includes 100% of the code.
