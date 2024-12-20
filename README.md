### ---- README IS STILL IN PROGRESS ----

# KEV
1. [INTRODUCTION](#introduction)<br>
2. [BASICS OF KEV](#basics-of-kev)<br>
3. [DIVE INTO THE CODE](#dive-into-the-code)<br>
    - [GENERAL CONTROL FLOW](#general-control-flow)<br>
    - [REPL](#repl)<br>
    - [TOKEN](#token)<br>
    - [OBJECT](#object)<br>
    - [AST](#ast)<br>
    - [LEXER](#lexer)<br>
    - [PARSER](#parser)<br>
    - [EVALUATOR](#evaluator)<br>

## INTRODUCTION
kev is a small and simple interpreted language I built as a way of not only expanding my knowlage in design patterns of interpreted languages but also to deepen my understanding of the go programming language it's self. 

Since this was a learning opportunity I tried to comment things rather... *verbosely*. So please don't be offended if you come across comments that seem elementary and describe processes rather granularly. I found this to be a necessary evil to ensure and help with the understanding of the source code. 

If you want to just test out kev and see what it can do, you can start with the [BASICS OF KEV](#basics-of-kev) section. However if you want a more indepth explanation of the code you can skip to the [DIVE INTO THE CODE](#dive-into-the-code) section.  

## BASICS OF KEV

## DIVE INTO THE CODE
This section was orginally going to go step by step through the processes and control flow of executing a given set of instructions. However, I found that approach to be sporadic when trying to explain each component of the code. So, This section is split into eight parts each of which has the following corresponding goal.<br>
<br>**[GENERAL CONTROL FLOW](#general-control-flow)**<br>
Give the reader an overview and basic understanding of the *flow* the source code takes to get an output given an input<br>
<br>**[REPL](#repl)**<br>
Explain the process behind how the repl (read-eval-print-loop) handles inputs and outputs<br>
<br>**[TOKEN](#token)**<br>
Give an overview of what tokens are and what they represent inside the kev programming language<br>
<br>**[OBJECT](#object)**<br>
Give an overview of what objects are and how they are used in kev<br>
<br>**[AST](#ast)**<br>
Explain how the ast (abstract syntax tree) is created and used<br>
<br>**[LEXER](#lexer)**<br>
Explain the process of how the lexer crawls through the input and how its able to read identifiers, numbers, and strings<br>
<br>**[PARSER](#parser)**<br>
Explain how Pratt Parsing is implemented in the parser and how the parser returns proper ast structures<br>
<br>**[EVALUATOR](#evaluator)**<br>
Explain the process behind how the evaluator takes in ast nodes and returns the correct objects which are used in getting the final output<br>
<br>***Keep in mind I commented the code to the best of my ability, so it wouldn't be bad practice to have the code as a reference as you read through each part.*** 

### GENERAL CONTROL FLOW
When the code is first initalized main.go is ran which starts a new repl. This repl.go file creates a new enviroment and scans for the code input. Once the input is received a new lexer is created which is used to initalize a parser. This parser will populate a `ast.Program` with all the given statments. These statments are parsed based on the current token type, more indepth explanation in the [PARSER](#parser) section. After the `ast.Program` is filled with the parsed statements the parser will return the filled `ast.Program` struct back to the repl which checks for parsing errors. If this check passes, the evaluator gets called with the `ast.Program` and the `env` we created earlier being passed to it. The evaluator will then check the type of the `ast.Node`, This can be anything like the program itself, an if statement, return statement, etc... based on the node type, the evaluator will process and evaluate the node. The main thing that needs attention is the `Eval(node ast.Node, env *object.Environment)` function is recursively called quite often. Finally the evaluator returns a `object.Object` which holds the type (represented as the `Type()` function in the `object.Object`) and the value (represented as the `Inspect()` function in the `object.Object`). This `object.Object` is stored in a variable named evaluated. To print the output to the console, the repl simply uses `evaluated.Inspect()` to get the object value and writes it to the console.          

### REPL

### TOKEN

### PARSER

### OBJECT

### LEXER

### EVALUATOR

### AST

