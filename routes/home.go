package routes

import (
	"BeardBar_on_GO/db/documents"
	"BeardBar_on_GO/models"
	"BeardBar_on_GO/session"
	"fmt"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
)

func IndexHandler(rnd render.Render, s *session.Session, db *mgo.Database) {

	fmt.Println(s.Username)

	postDocuments := []documents.PostDocument{}
	postsCollection := db.C("posts")
	postsCollection.Find(nil).All(&postDocuments)

	posts := []models.Post{}
	for _, doc := range postDocuments {
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkdown}
		posts = append(posts, post)
	}

	rnd.HTML(200, "index", posts)
}
