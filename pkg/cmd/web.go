// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korrel8r/client/pkg/browser"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web [FLAGS]",
	Short: "Connect to remote korrel8r, show graphs in local HTTP server.",
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, args []string) error {
		gin.DefaultWriter = log.Writer()
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		router := gin.New()
		router.Use(gin.Recovery())
		router.Use(gin.Logger())
		c := newClient()
		b, err := browser.New(c, router)
		if err != nil {
			return err
		}
		defer b.Close()
		s := http.Server{
			Addr:    *addr,
			Handler: router,
		}
		log.Println("Listening on ", *addr, " connected to ", *korrel8rURL)
		return s.ListenAndServe()
	},
}

var (
	addr *string
)

func init() {
	rootCmd.AddCommand(webCmd)
	addr = webCmd.Flags().StringP("addr", "a", ":8081", "Listening address for web server")
}
