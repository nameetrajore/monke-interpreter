# Monke Interpreter

This is an interpreter for the Monke programming language, written in Go. The Monke language is a simple, toy language that is designed to be easy to learn and implement. This interpreter is based on the book &ldquo;Writing an Interpreter in Go&rdquo; by Thorsten Ball.


<a id="org683af5b"></a>

## Features

1.  The interpreter can parse and execute Monkey source code.
2.  The interpreter has a REPL (Read-Eval-Print Loop) that allows you to interactively enter and execute Monke code.
3.  The interpreter supports a myriad of features like
    - Functions
    - Variables (Integer and String)
    - Statements
    - Expressions
    - Arrays
    - Hashmaps


<a id="org5f23474"></a>

## Getting Started

To get started, you will need to have Go installed on your computer. Once you have Go installed, you can clone this repository and run the interpreter in REPL mode with the following command: 
```
go run main.go
```

**OR** if you want to feed in a program you can write it in a file with the extension .grr and give it to the interpreter in the following manner:
```
go run main.go test.grr
```
