package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"

	firebase "firebase.google.com/go"
)

type Fb struct {
	client *firestore.Client
	ctx    context.Context
}

func hello(name string) string {
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}

func Home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/home.gohtml")
	t.Execute(w, nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/create.gohtml")
	t.Execute(w, nil)
}

func (fb *Fb) CreateSubmit(w http.ResponseWriter, r *http.Request) {

	_, _, err := fb.client.Collection("users").Add(fb.ctx, map[string]interface{}{
		"pin":  "Ada",
		"last": "Lovelace",
		"born": 1815,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}

func Join(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/join.gohtml")
	t.Execute(w, nil)
}

func main() {

	// Use the application default credentials
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "dmzoo-b6d89"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	fb := &Fb{client: client, ctx: ctx}

	r := mux.NewRouter()
	r.HandleFunc("/create", fb.CreateSubmit).Methods("POST")
	r.HandleFunc("/create", Create)
	r.HandleFunc("/join", Join)
	r.HandleFunc("/", Home)

	log.Fatal(http.ListenAndServe(":80", r))

}
