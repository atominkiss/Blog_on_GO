package main

import (
	"BeardBar_on_GO/db/documents"
	"BeardBar_on_GO/models"
	"BeardBar_on_GO/session"
	"BeardBar_on_GO/utils"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
)

const COOKIE_NAME = "sessionId"

//var posts map[string]*models.Post
var postsCollection *mgo.Collection
var inMemorySession *session.Session

func getLoginHandler(rnd render.Render) {
	rnd.HTML(200, "login", nil)
}

func postLoginHandler(rnd render.Render, r *http.Request, w http.ResponseWriter) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username)
	fmt.Println(password)

	sessionId := inMemorySession.Init(username)
	cookie := &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   sessionId,
		Expires: time.Now().Add(5 * time.Minute),
	}
	http.SetCookie(w, cookie)
	rnd.Redirect("/")
}

func indexHandler(rnd render.Render, r *http.Request) {
	cookie, _ := r.Cookie(COOKIE_NAME)
	inMemorySession.Get(cookie.Value)
	postDocuments := []documents.PostDocument{}
	postsCollection.Find(nil).All(&postDocuments)

	posts := []models.Post{}
	for _, doc := range postDocuments {
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkdown}
		posts = append(posts, post)
	}

	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	post := models.Post{}
	rnd.HTML(200, "write", post)
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)
	if err != nil {
		rnd.Redirect("/")
		return
	}

	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMarkdown}

	rnd.HTML(200, "write", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := utils.ConvertMarkdownToHtml(contentMarkdown)

	postDocument := documents.PostDocument{id, title, contentHtml, contentMarkdown}
	if id != "" {
		postsCollection.UpdateId(id, postDocument)
	} else {
		id = utils.GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
	}

	rnd.Redirect("/")
}

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]

	if id == "" {
		rnd.Redirect("/")
		return
	}

	postsCollection.RemoveId(id)
	rnd.Redirect("/")
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	inMemorySession := session.NewSession()
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	postsCollection = session.DB("blog").C("posts")

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

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

	m.Get("/", indexHandler)
	m.Get("/login", getLoginHandler)
	m.Post("/login", postLoginHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Post("/SavePost", savePostHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/gethtml", getHtmlHandler)

	m.Run()
}
