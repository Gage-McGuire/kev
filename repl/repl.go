package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kev/lexer"
	token "github.com/kev/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		for next_token := l.NextToken(); next_token.Type != token.EOF; next_token = l.NextToken() {
			fmt.Printf("%+v\n", next_token)
		}
	}
}
