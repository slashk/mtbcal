package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/markbates/going/defaults"
	"github.com/slashk/mtbcal/actions"
)

func main() {
	port := defaults.String(os.Getenv("PORT"), "3000")
	log.Printf("Starting mtbcal on %s at port %s\n", actions.App().Host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), actions.App()))
}
