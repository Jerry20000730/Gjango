package File

import (
	context "github.com/Jerry20000730/Gjango/web/Context"
	"github.com/Jerry20000730/Gjango/web/Utils"
	"net/http"
	"net/url"
)

type FileManager struct{}

func (f *FileManager) Download(ctx *context.Context, fileName string) {
	http.ServeFile(ctx.W, ctx.R, fileName)
}

func (f *FileManager) FileAttachment(ctx *context.Context, filePath, fileName string) {
	// Content-Disposition:
	// An opportunity to raise a "File Download" dialogue box for
	// a known MIME type with binary format or suggest a filename for dynamic content.
	// Quotes are necessary with special characters.
	if Utils.IsASCII(fileName) {
		ctx.W.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`"`)
	} else {
		ctx.W.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(fileName))
	}
	http.ServeFile(ctx.W, ctx.R, filePath)
}

func (f *FileManager) FileFromFileSystem(ctx *context.Context, filePath string, fs http.FileSystem) {
	// filePath is the relative path of the FileSystem
	defer func(old string) {
		ctx.R.URL.Path = old
	}(ctx.R.URL.Path)

	ctx.R.URL.Path = filePath

	http.FileServer(fs).ServeHTTP(ctx.W, ctx.R)
}
