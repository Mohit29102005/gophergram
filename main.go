package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"

	"github.com/kanozec/gophergram/controllers"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kanozec/gophergram/model"
)

type config struct {
	ServerPort string `env:"SERVERPORT" envDefault:"8080"`
}

func tmpInitialize(templatesDir string) multitemplate.Renderer {
	render := multitemplate.NewRenderer()

	frontLayouts, err := filepath.Glob(templatesDir + "/layouts/*.gohtml")
	if err != nil {
		panic(err.Error)
	}

	includes, err := filepath.Glob(templatesDir + "/includes/*.gohtml")
	if err != nil {
		panic(err.Error)
	}
	statics, err := filepath.Glob(templatesDir + "/*.html")

	for _, include := range includes {
		layoutClone := make([]string, len(includes))
		copy(layoutClone, frontLayouts)
		files := append(layoutClone, include)
		render.AddFromFiles(strings.Replace(filepath.Base(include), ".gohtml", "", 1), files...)
		// render.AddFromFiles(filepath.Base(include), files...)
	}
	for _, s := range statics {
		render.AddFromFiles(strings.Replace(filepath.Base(s), ".html", "", 1), s)
	}
	return render
}

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("File .env not found, using the default value")
	}

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		fmt.Println("Failed to parse ENV")
	}

	db, err := model.InitializeDB()
	if err != nil {
		fmt.Println("Initialize db error")
	}
	defer db.Close()
	router := gin.Default()
	// router.LoadHTMLGlob("templates/*.html")

	//TODO
	// var templatePage = []string{
	// 	"templates/layouts/frontend.gohtml",
	// 	"templates/layouts/navbar.gohtml",
	// 	"templates/layouts/footer.gohtml",
	// }
	// templatePage = append(templatePage, "templates/includes/contact.gohtml")
	// router.SetHTMLTemplate(template.Must(template.ParseFiles(templatePage...)))
	router.HTMLRender = tmpInitialize("./templates")
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404", nil)
	})
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})
	router.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact", nil)
	})
	staticC := controllers.NewStatic()
	router.GET("/type", staticC.Contact)
	router.Run(":" + cfg.ServerPort)
}
