package main

import (
	"github.com/himanshugarg165/Employee-Management/server"
)

func main() {
	e := server.New()
	e.Logger.Fatal(e.Start("127.0.0.1:9000"))
}
