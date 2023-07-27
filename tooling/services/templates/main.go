package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gitamped/fertilize/parser"
)

func main() {

	pattern := flag.String("pattern", "github.com/kjvonly/service/services/users", "package file path starting with github.com/kjvonly/")
	directory := flag.String("directory", "../../../services/user", "package directory")
	excluded := flag.String("excluded", "UserRpcService", "excluded interfaces")
	flag.Parse()
	err := os.Chdir(filepath.Join(*directory))
	if err != nil {
		log.Fatalf("error changing directory")
	}
	patterns := []string{*pattern}
	p := parser.New(patterns...)
	p.ExcludeInterfaces = []string{*excluded}
	p.Verbose = false
	log.Println(os.Getwd())
	def, err := p.Parse()
	if err != nil {
		panic(fmt.Sprintf("err parsing: %s", err))
	}
	b, err := json.Marshal(def)
	t, _ := ioutil.ReadFile("./handlers.tmpl")
	var data map[string]parser.Definition
	json.Unmarshal(b, &data)

	tmpl, _ := template.New("test").Parse(string(t))

	for k, v := range data {
		p := strings.Replace(k, *pattern, "", -1)
		os.Truncate(filepath.Join(p, "handlers.go"), 0)
		f, err := os.OpenFile(filepath.Join(p, "handlers.go"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		tmpl.Execute(f, v)
		if err != nil {
			log.Fatal(err)
		}
	}
}
