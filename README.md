## Token_go
Simple & fast Encoder/Decoder for tiktoken vocabulary.
Implemented from scratch(no regex library). Tokenizer is in vocab.go which has ~120 lines of code.



## Performance
p50k_base.tiktoken:
- Encoder: 4.625M toks/sec, 19.143 MB/sec, 1 thread
- Decoder: 37.817M toks/sec, 156.516 MB/sec, 1 thread

cl100k_base.tiktoken:
- Encoded 3.949M toks/sec, 16.748 MB/sec, 1 thread
- Decoded 35.825M toks/sec, 151.952 MB/sec, 1 thread

Server(p50k_base)
- 8x clients calls 100K times Encode("Hi there!" + index).
- 800K total requests in 26.7sec => 30K req/sec.



## Examples
Encode/Decode:
<pre><code>vb, err := NewVocab("p50k_base.tiktoken", true)

toks := vb.Encode("Hi there!")
fmt.Println(toks)

str := vb.Decode(toks)
fmt.Println(str)
</code></pre>

Client/Server:
<pre><code>go NewServer("8090", true)   //run server in extra thread

client := NewClient("localhost:8090", "p50k_base")

toks, err := client.Encode([]byte("Hi there!"))
fmt.Println(toks)

text, err := client.Decode([]int{17250, 612, 0})
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
