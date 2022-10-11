package main

import (
	"fmt"
	"log"
	"portfolio_golang/src/rest"
	"time"
)

func main() {
	log.Printf(fmt.Sprintf("[GIN|MAIN] %s | INIT MAIN APP", time.Now().Format("2006/01/02 - 15:04:05")))

	r := rest.RunGinWebApp()
	r.Run("127.0.0.1:8082")
}
