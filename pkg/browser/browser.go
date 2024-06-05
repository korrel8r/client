// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

// package browser implements an HTML UI for web browsers.
package browser

import (
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/korrel8r/client/pkg/swagger/client"
)

var (
	//go:embed templates
	templates embed.FS
	//go:embed images
	images embed.FS
)

// Browser implements HTTP handlers for web browsers.
type Browser struct {
	client     *client.RESTAPI
	router     *gin.Engine
	images     http.FileSystem
	dir, files string
}

func New(restClient *client.RESTAPI, router *gin.Engine) (*Browser, error) {
	images, err := fs.Sub(images, "images")
	if err != nil {
		return nil, err
	}
	b := &Browser{
		client: restClient,
		router: router,
		images: http.FS(images),
	}
	if b.dir, err = os.MkdirTemp("", "korrel8r"); err == nil {
		b.files = filepath.Join(b.dir, "files")
		err = os.Mkdir(b.files, 0700)
	}
	if err != nil {
		return nil, err
	}
	log.Println("Using temporary directory: ", b.dir)
	c := &correlate{Browser: b}

	tmpl := template.Must(template.New("").ParseFS(templates, "templates/*.tmpl"))
	router.SetHTMLTemplate(tmpl)
	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/correlate") })
	router.GET("/correlate", c.HTML)
	router.Static("/files", b.files)
	router.StaticFS("/images", b.images)
	// Display errors converted into URLs.
	router.GET("/error", func(c *gin.Context) {
		httpError(c, errors.New(c.Request.URL.Query().Get("err")), http.StatusNotFound)
	})

	return b, nil
}

// Close should be called on shutdown to clean up external resources.
func (b *Browser) Close() { _ = os.RemoveAll(b.dir) }

func httpError(c *gin.Context, err error, code int) bool {
	if err != nil {
		_ = c.Error(err)
		c.HTML(code, "error.html.tmpl", c)
	}
	return err != nil
}
