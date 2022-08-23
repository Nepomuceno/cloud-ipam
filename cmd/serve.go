package cmd

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nepomuceno/cloud-ipam/api"
	"github.com/nepomuceno/cloud-ipam/ui"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a front end server and api for cloud IPAM",
	RunE:  createServerAndServe,
}

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "List server static files",
	RunE:  serveStaticFiles,
}

func serveStaticFiles(cmd *cobra.Command, args []string) error {
	files, error := fs.Glob(ui.StaticFiles, "**/*.*")
	if error != nil {
		return error
	}
	for _, file := range files {
		cmd.Printf("%s\n", file)
	}
	return nil
}

func createServerAndServe(cmd *cobra.Command, args []string) error {
	srv := gin.Default()
	err := router(srv)
	if err != nil {
		return err
	}
	return srv.Run(":8080")
}

func router(r *gin.Engine) error {
	api.NewApiClient(client).RegisterRoutes(r)
	subFS, err := fs.Sub(ui.StaticFiles, "dist")
	if err != nil {
		return err
	}
	staticFS, err := fs.Sub(subFS, "static")
	if err != nil {
		return err
	}
	r.StaticFS("/static", http.FS(staticFS))
	r.StaticFileFS("/", "main.html", http.FS(subFS))
	r.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(subFS))
	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.AddCommand(filesCmd)
}
