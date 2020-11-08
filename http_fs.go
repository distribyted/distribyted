// +build !release

package distribyted

//go:generate go run ./build_tools/assets_generate/main.go
import (
	"net/http"

	"github.com/shurcooL/httpfs/union"
)

var HttpFS = union.New(map[string]http.FileSystem{
	"/assets":    http.Dir("assets"),
	"/templates": http.Dir("templates"),
})
