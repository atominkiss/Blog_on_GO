package main

import (
	"BeardBar_on_GO/models"
	"fmt"
	"html/template"
	"net/http"
)

var posts map[string]*models.Post

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	//fmt.Println(posts)

	t.ExecuteTemplate(w, "index", posts)
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "write", nil)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	id := r.FormValue("id")
	post, found := posts[id]
	if !found {
		http.NotFound(w, r)
	}

	t.ExecuteTemplate(w, "write", post)
}

func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id = GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	http.Redirect(w, r, "/", 302)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		http.NotFound(w, r)
	}

	delete(posts, id)

	http.Redirect(w, r, "/", 302)
}

func main() {
	fmt.Println("Listening on port :8080")
	posts = make(map[string]*models.Post, 0)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/SavePost", savePostHandler)
	http.HandleFunc("/delete", deleteHandler)

	http.ListenAndServe(":8080", nil)
}
