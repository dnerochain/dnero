package main

import (
	"log"
	"strings"

	"github.com/spf13/cobra/doc"
	dnero "github.com/dnerochain/dnero/cmd/dnero/cmd"
	dnerocli "github.com/dnerochain/dnero/cmd/dnerocli/cmd"
)

func generateDneroCLIDoc(filePrepender, linkHandler func(string) string) {
	var all = dnerocli.RootCmd
	err := doc.GenMarkdownTreeCustom(all, "./wallet/", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func generateDneroDoc(filePrepender, linkHandler func(string) string) {
	var all = dnero.RootCmd
	err := doc.GenMarkdownTreeCustom(all, "./ledger/", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	filePrepender := func(filename string) string {
		return ""
	}

	linkHandler := func(name string) string {
		return strings.ToLower(name)
	}

	generateDneroCLIDoc(filePrepender, linkHandler)
	generateDneroDoc(filePrepender, linkHandler)
	Walk()
}
