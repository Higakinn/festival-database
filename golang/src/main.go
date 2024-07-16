package main

import (
	"fmt"
	"os"
)


func main() {
	lang := os.Getenv("PROG_LANG")
	fmt.Println(fmt.Sprintf("%s batch start!!", lang))
	fmt.Println(fmt.Sprintf("%s batch end!!", lang))
}