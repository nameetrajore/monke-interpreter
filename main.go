package main
import (
	"fmt"
	"log"
	"os"
	"os/user"
	"monke/repl"
	"path/filepath"
)

func main(){
	user, err := user.Current()
	if err != nil{
		panic(err)
	}

	if len(os.Args) < 2 {

	fmt.Printf("Hello %s! This is the Monke Programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands \n")
	repl.Start(os.Stdin, os.Stdout)

	} else {

	fileName := os.Args[1]
	file, err := os.Open(fileName)

	if filepath.Ext(fileName) != ".grr" && filepath.Ext(fileName) != ".brr" && filepath.Ext(fileName) != ".hoot" && filepath.Ext(fileName) != "coo" {
			fmt.Println("Invalid file extension. Expected .grr | .brr | .coo | .hoot .")
			return
		}

	if err != nil {
		log.Fatal(err)
	}

	repl.Interpret(file, os.Stdout)
	defer file.Close()
	}
}
