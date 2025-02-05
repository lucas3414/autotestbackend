package main

import "go-gin-demo/cmd"

func main() {

	defer cmd.Clean()

	cmd.Start()

}
