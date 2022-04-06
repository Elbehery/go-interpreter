package repl

import (
	"bufio"
	"io"
	"playground/go-interpreter/src/lexer"
	"playground/go-interpreter/src/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errs []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, e := range errs {
		io.WriteString(out, "\t"+e+"\n")
	}
}
