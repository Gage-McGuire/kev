#### Following adds support for passing files instead of prompting throght the terminal
```
// builds the kev binary and installs it so the kev cmd can be used
go build -o kev
go install

// in the dir your .kev file is located
kev run <fileName>.kev
```

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
    - [EVALUATOR](#evaluator)<br>

## INTRODUCTION
kev is a small and simple interpreted language I built as a way of not only expanding my knowlage in design patterns of interpreted languages but also to deepen my understanding of the go programming language it's self. 

Since this was a learning opportunity I tried to comment things rather... *verbosely*. So please don't be offended if you come across comments that seem elementary and describe processes rather granularly. I found this to be a necessary evil to ensure and help with the understanding of the source code. 

If you want to just test out kev and see what it can do, you can start with the [BASICS OF KEV](#basics-of-kev) section. However if you want a more indepth explanation of the code you can skip to the [DIVE INTO THE CODE](#dive-into-the-code) section.  

## BASICS OF KEV
*scripts to add kev to your path are in the works*
### Setup
For now if you want to start seeing what Kev can do you need to have Golang installed. Then you can `go run .` in the project directory. You should have something like the following show up.
```
       ,--.                       
   ,--/  /|    ,---,.             
,---,': / '  ,'  .' |       ,---. 
:   : '/ / ,---.'   |      /__./| 
|   '   ,  |   |   .' ,---.;  ; | 
'   |  /   :   :  |-,/___/ \  | | 
|   ;  ;   :   |  ;/|\   ;  \ ' | 
:   '   \  |   :   .' \   \  \: | 
|   |    ' |   |  |-,  ;   \  ' . 
'   : |.  \'   :  ;/|   \   \   ' 
|   | '_\.'|   |    \    \   `  ; 
'   : |    |   :   .'     :   \ | 
;   |,'    |   | ,'        '---"  
'---'      `----'                 
                                           
>> 
```
Now you can start typing commands!

### Variables
You're able to define a new variable with the `var` key word
```
var x = 5
var hello = "world"
var map_example = {"key":"value"}
```

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

<br>
<br>

### GENERAL CONTROL FLOW
When the code is first initalized `main.go` is ran which starts a new repl. This repl.go file creates a new enviroment and scans for the code input. Once the input is received a new lexer is created which is then used to initalize a parser. This parser will populate a `ast.Program` with all the given statments. These statments are parsed based on the current token type, more indepth explanation in the [PARSER](#parser) section. After the `ast.Program` is filled with the parsed statements the parser will return the filled `ast.Program` struct back to the repl which checks for parsing errors. If this check passes, the evaluator gets called with the `ast.Program` and the `env` we created earlier being passed to it. The evaluator will then check the type of the `ast.Node`, This can be anything like the program itself, an if statement, return statement, etc... based on the node type, the evaluator will process and evaluate the node. The main thing that needs attention is the `Eval(node ast.Node, env *object.Environment)` function is recursively called quite often. Finally the evaluator returns a `object.Object` which holds the type (represented as the `Type()` function in the `object.Object`) and the value (represented as the `Inspect()` function in the `object.Object`). This `object.Object` is stored in a variable named evaluated. To print the output to the console, the repl simply uses `evaluated.Inspect()` to get the object value and writes it to the console.

<br>
<br>

### REPL
The `repl.go` file first has a couple of constants defined. The first is the color codes that are used in `io.Writer`, which honestly does nothing except help make the errors look pretty. The next constant is the `PROMPT`. This is used as an indicator to the user when the repl is ready for the next prompt.

The most of the work happens in the `Start(in io.Reader, out io.Writer)` function. The `in` and `out` args are used to read user input and print the programs output respectively. We then start this function with `scanner` and `env` variables. The `scanner` is used to initalize a new `bufio.NewScanner()` which is used later in the function to scan for user's inputs. Then `env` calls the `NewEnviroment()` function from the object package which simply initalizes an empty enviroment object that will be used later to store our programs variables and what not.

Now we come across a for loop that has no termination statment so it will repeat the following code until the user decideds to terminate the interpreter. This loop starts by printing that `PROMPT` indicator we talked about earlier to specify its ready to accept instructions. We then simply scan for the user's input and get the line from the scanner. Now this is where we start to really see other parts of the program to come into play. First a new [LEXER](#lexer) is created which returns a pointer to the `Lexer{}` struct which is then used to create a new [PARSER](#parser) which gives us a new pointer to the `Parser{}` struct. It's ok if you don't know what the lexer and parser are doing, we'll cover that in later sections. The main thing you need to know is we are collecting the user's input and getting the parser ready to parse it into an `ast.Program`.

It's finally time to parse and evaluate. we simply start by calling `ParseProgram()` on the parser pointer. This give us the `*ast.Program` we've been needing and stores it into a variable named `program`. Now, before we start evaluating, we need to check to see if there was any parsing errors. This is done by simply calling the `Errors()` function in the parser, which returns a list errors that occured, and checking if that list is empty.  

Now lets say an error did occur. Our if statment would evaluate to true and the `printParserErrors(out io.writer, errors []string)` function gets called. This function doesn't need much explaining, most of it is just formating to make the errors look better. All that is really happing is we are indexing through the list of errors and printing them out.

Now lets say no errors occur. Well first, congratulations! We get to move onto the evaluator. We start by calling the `Eval()` function in the evaluator package, which is passed the `program` variable which is a `*ast.Program` and the `env` variable we made earlier in the function which is a `*object.Enviroment`. The evaluator gives us back a `object.Object`, which if you remember mentioned in the [GENERAL CONTROL FLOW](#general-control-flow) section, has a function called `Inspect()`. We use this `Inspect()` function on the object to get its value a print it to the console. This value could be nothing, a string, a integer, etc...

With that we finally come to the end of the repl. All that happens now is the for loop gets evaluated again, and keeps going through process again until the user terminates the program. 

<br>
<br>

### TOKEN
Let's take a look at the `token.go` file. This file really doesn't have much going on but is used quite a lot in other places of the program. First we see there is a whole bunch of constants. These are all characters or identifiers that are predefined for the language. When we look at these constants they are super simular if not the same as many other languages.

The main thing this code is revloving around is the `Token{}` struct. This struct holds the `TokenType` which is defined by the constants above and the `Literal` which simply put is the literal representation whatever the token is. If there's any confusion hopfully the following examples will clear things up.
```
The integer 10 would have the following Token struct:

Token {
    Type = INT
    Literal = 10
}


The string "Hello world" would have the following Token struct:

Token {
    Type = STRING
    Literal = Hello world
}


The Bool true would have the following Token struct:

Token {
    Type = TRUE
    Literal = true
}


The Identifier x would have the following Token struct:

Token {
    Type = IDENT
    Literal = x
}
```

All this might not mean much but in later sections we will see how tokens are populated and handled. For now we might be asking the question "How do you know if an identifier is actually a keyword like var, func, return ect...". This is done with the `LookupIdent(ident string)` function. This function takes the identifier it's passed when called and checks if it's in the keyword map. For now you can assume anytime a new identifier is made in the program the `LookupIdent()` function gets called. This function simply returns the corresponding token type if the identifier is found in the keyword map, and returns the `IDENT` token type if it's not found. Simply put if the identifier is not found in the keyword map, we can assume it's some sort of variable name.

That's pretty much it, not much else to explain. Like I said, you'll understand the actual use cases for tokens later on when we start looking into the [AST](#ast), [LEXER](#lexer), and [PARSER](#parser).

<br>
<br>

### OBJECT
Now we might be jumping the gun here since the object package doesn't get used/populated till later. However, sticking with the partern of introducing the different structs that are used we'll just go ahead a do it.

So in the `object.go` file you can start to see the different data types and what exactly is being stored as an object in kev. For instance, The `Hash{}` struct does not just hold a hash map. It holds a value named `Pairs` which is a map of `HashKeys` and `HashPairs` both of which have their own structs that hold values that are important to its own goals. The Hash data type might be one of the more advance objects so we'll come back to it once we see whats going on with the rest of the code in `object.go`.

Alrighty, starting from the top! The first thing we come across is the `ObjectType` and this is exactly what it says it is. It simply stores the data type of said object. We can assume this `ObjectType` will be one the following types that are defined in the constant below. This brings us to the `Object` interface. This was stated back in [GENERAL CONTROL FLOW](#general-control-flow) but this interface implements the `Type()` and `Inspect()` functions. The `Type()` function will return the data type of the object and the `Inspect()` function will return the value of the Object. This `Object` interface can be thought of as the base representation of an Object. So whenever you see a `object.Object` you know its this interface.

Next lets look at the objects that make use of Go's native data types:
```
type Integer struct {
	Value int64
}

type String struct {
	Value string
}

type Boolean struct {
	Value bool
}

type Null struct{}

type Error struct {
	Message string
}
```
All theses objects simply have the same data type as they do in go so there really isn't much "sauce" going on.

Now lets see the more unique/complex objects:
```
type ReturnValue struct {
	Value Object
}

type Array struct {
	Elements []Object
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

type Hash struct {
	Pairs map[HashKey]HashPair
}
```
Let's just go down the list starting with the `ReturnValue` struct. This struct represents the return statement. Something to keep in mind is this object holds the value of the return AFTER it is evaluated... here is an example
```
return(2+2) is represented as the following object

ReturnValue {
    Value = Object{ 
        Type() = INTEGER_OBJ
        Inspect() = 4
    }
}
``` 
The `Array` object works the same what but just holds all the objects in a slice. The `Function` object is where we start to run into things that we really havent seen before. Since we haven't discussed the [AST](#ast) yet don't worry to much about knowing exactly whats happening, I'll give you the gist. First we have a value called `Parameters` this is a slice of identifiers that are passed to the function. Next we have a `Body` which is a block statement, simply put its the code contained in the function. Lastly we have the `Env` which is the Enviroment object, we'll talk about this object in a second. Just know this is the outer enviroment that is surrounding the function. Just like the function object the `Hash` object might have some moving parts that dont make sense right now. The only thing you need to know is the key and value are stored in a `map[HashKey]HashPair`. The `Hash` object will come together and make sense in the [EVALUATOR](#evaluator).

The one thing that is consistent across all these objects is how they implement the `Type()` and `Inspect()` functions. Some of these functions have more going on than others but it's really just for string formatting and making the string representation of the object look pretty.  

BUT WAIT THERES MORE! If you look at the end of `object.go` you get to the section that holds the code for enviroments and built-in functions. I promise this is the last part of `object.go` and we can move onto the actual fun stuff.

When it comes to the enviroment, it all starts with the `NewEnvironment()` functon. This function makes a new store which is a `map[string]Object` and sets `outer = nil`. Outer will eventaually be set to a `*Enviroment`, This will represent the encompassing enviroment. The very next function is used to do exactly that. `NewEnclosedEnvironment(outer *Environment)` makes a new environment then takes the surrounding environment and sets it to outer. This way global/surrounding vairables will still be available in the enclosed enviroment but prohibits the action the other way around. The outer environment will never know the existence of the enclosed environment. Now let's see how we set and get the objects stored in environments. First we can set objects in the environment with the `Set(name string, val Object)` function. This function is called on a pointer to an environment so the name/object will get set in the corresponding store. Just like the `Set()` function, the `Get(name string)` function is called on a pointer to an environment. It first checks if the environment has the name in it's store. If it does, it returns the corresponding object. However, if doesn't it checks if there is a outer environment set, if so it then checks that environment.

Finally we get to built-in function objects. There isn't much going on since all thats really happening is the `Builtin` object stores the function that is to be used. Again this is another situation where you don't need to worry to much about whats going on and things pertaining to builtins will be explained better in the [EVALUATOR](#evaluator). 

With that we finished the object package! To recap, all that we are trying to doing in this package is make structs that can represent different aspects of objects that we use within kev. These aspects are super benfical with telling us the data types of the objects, giving us string representations of the values of the objects, and holding key value pairs for certain objects.    

<br>
<br>

### LEXER

<br>
<br>

### PARSER

<br>
<br>

### EVALUATOR
