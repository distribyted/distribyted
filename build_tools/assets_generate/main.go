// +build ignore

package main

import (
	"github.com/distribyted/distribyted"
	"github.com/shurcooL/vfsgen"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := vfsgen.Generate(distribyted.HttpFS, vfsgen.Options{
		BuildTags:    "release",
		VariableName: "HttpFS",
		PackageName:  "distribyted",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
