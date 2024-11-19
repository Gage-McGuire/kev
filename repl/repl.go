package repl

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/kev/evaluator"
	"github.com/kev/lexer"
	"github.com/kev/object"
	"github.com/kev/parser"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect()+"\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, Red+"****** PARSING ERROR ******\n"+Reset)
	io.WriteString(out, Yellow+"THE FOLLOWING ERRORS OCCURED:\n"+Reset)
	for idx, msg := range errors {
		io.WriteString(out, Gray+strconv.Itoa(idx+1)+": "+msg+Reset+"\n")
	}
}
