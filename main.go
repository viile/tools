package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println(time.Now().Unix())
	fmt.Println("命令行参数数量:", len(os.Args))
	for k, v := range os.Args {
		fmt.Printf("args[%v]=[%v]\n", k, v)
	}
}
