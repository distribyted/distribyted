// +build ignore

package main

import (
	"log"

	"github.com/ajnavarro/distribyted"
	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(distribyted.Assets, vfsgen.Options{
		BuildTags:    "release",
		VariableName: "Assets",
		PackageName:  "distribyted",
	})
	if err != nil {
		log.Fatalln(err)
	}

	err := vfsgen.Generate(distribyted.Assets, vfsgen.Options{
		BuildTags:    "release",
		VariableName: "Templates",
		PackageName:  "distribyted",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
