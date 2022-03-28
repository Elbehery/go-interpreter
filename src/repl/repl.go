package repl

import (
	"bufio"
	"fmt"
	"io"
	"playground/go-interpreter/src/lexer"
	"playground/go-interpreter/src/token"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Println(PROMPT)
		if !scanner.Scan() {
			return
		}

		txt := scanner.Text()
		l := lexer.New(txt)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
