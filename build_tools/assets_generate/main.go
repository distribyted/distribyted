// +build ignore

package main

import (
	"github.com/distribyted/distribyted"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(distribyted.HttpFS, vfsgen.Options{
		BuildTags:    "release",
		VariableName: "HttpFS",
		PackageName:  "distribyted",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("problem generating static files")
	}
}
