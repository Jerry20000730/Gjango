package File

import (
	context "github.com/Jerry20000730/Gjango/web/Context" // Custom context package for handling web contexts.
	"github.com/Jerry20000730/Gjango/web/Utils"           // Utility package for miscellaneous helper functions.
	"net/http"
	"net/url"
)

// FileManager provides methods for file handling operations such as downloading,
// attaching files for download, and serving files from a file system.
type FileManager struct{}

// Download sends a file to the client for download. It uses http.ServeFile
// to serve the file located at fileName.
//
// Parameters:
// - ctx: The web context containing the HTTP request and response writer.
// - fileName: The path to the file to be downloaded.
func (f *FileManager) Download(ctx *context.Context, fileName string) {
	http.ServeFile(ctx.W, ctx.R, fileName)
}

// FileAttachment sets the Content-Disposition header to attach a file for download,
// suggesting a filename for the client to use. It handles ASCII and non-ASCII filenames
// by properly encoding them.
//
// Parameters:
// - ctx: The web context containing the HTTP request and response writer.
// - filePath: The path to the file to be attached for download.
// - fileName: The suggested filename for the download. Special characters are handled.
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

// FileFromFileSystem serves a file from a specified file system. It temporarily modifies
// the request URL path to the path of the file to be served, ensuring the file server
// serves the correct file.
//
// Parameters:
// - ctx: The web context containing the HTTP request and response writer.
// - filePath: The relative path within the file system to the file to be served.
// - fs: The file system from which the file will be served.
func (f *FileManager) FileFromFileSystem(ctx *context.Context, filePath string, fs http.FileSystem) {
	// filePath is the relative path of the FileSystem
	defer func(old string) {
		ctx.R.URL.Path = old // Restore the original URL path after serving the file.
	}(ctx.R.URL.Path)

	ctx.R.URL.Path = filePath // Temporarily set the URL path to the file path.

	http.FileServer(fs).ServeHTTP(ctx.W, ctx.R) // Serve the file.
}
