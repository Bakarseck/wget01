package main

import (
	"fmt"
	"time"

	"github.com/Bakarseck/wget01/cmd"
)

func main() {
	fmt.Print("a")
	time.Sleep(time.Second*2)
	fmt.Println("\rb")
	cmd.Execute()
}
