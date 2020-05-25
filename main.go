package main

import (
	"BeardBar_on_GO/routes"
	"BeardBar_on_GO/session"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"html/template"
	"labix.org/v2/mgo"
)

//var posts map[string]*models.Post
//var postsCollection *mgo.Collection
//var inMemorySession *session.Session

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	//inMemorySession := session.NewSession()
	mongoSession, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	db := mongoSession.DB("blog")

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Map(db)
	m.Use(session.Middleware)
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		//Delims: render.Delims{"{[{", "}]}"}, 	// Sets delimiters to the specified strings.
		Charset:    "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,    // Output human readable JSON
		//IndentXML: true, 								// Output human readable XML
		//HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html"
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", routes.IndexHandler)
	m.Get("/login", routes.GetLoginHandler)
	m.Post("/login", routes.PostLoginHandler)
	m.Get("/logout", routes.LogoutHandler)
	m.Get("/write", routes.WriteHandler)
	m.Get("/edit/:id", routes.EditHandler)
	m.Get("/view/:id", routes.ViewHandler)
	m.Post("/SavePost", routes.SavePostHandler)
	m.Get("/delete/:id", routes.DeleteHandler)
	m.Post("/gethtml", routes.GetHtmlHandler)

	m.Run()
}
