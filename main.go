package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/envy"
	"github.com/slashk/mtbcal/actions"
)

func main() {
	port := envy.Get("PORT", "3000")
	log.Printf("Starting mtbcal on %s at port %s\n", actions.App().Host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), actions.App()))
}
