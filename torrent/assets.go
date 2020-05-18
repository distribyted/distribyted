// +build !release

package distribyted

//go:generate go run ./assets_generate.go
import (
	"net/http"
)

var Assets http.FileSystem = http.Dir("assets")
var Templates http.FileSystem = http.Dir("templates")
