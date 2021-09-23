package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"log"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func main() {
	InitEnv()

	key := os.Getenv("SECRET") // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30       // 30 days
	isProd := false            // Set to true when serving over https
	CLIENT_ID := os.Getenv("CLIENT_ID")
	SECRET := os.Getenv("SECRET")

	log.Println("CLIENT ID: ", CLIENT_ID)
	log.Println("SECRET:", SECRET)

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(CLIENT_ID, SECRET, "https://40ac-2800-300-6421-79c0-00-4.ngrok.io/auth/google/callback", "email", "profile"),
	)

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			log.Println("RESPONSE: ", fmt.Sprint(res), fmt.Sprint(err))
			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		log.Println("USER: ", fmt.Sprint(user))
		t.Execute(res, user)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})
	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}

func InitEnv() {
	env := os.Getenv("GO_ENVIRONMENT")
	if "" == env {
		env = "develop"
	}

	log.Print("Environment: " + env)
	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatal("Error loading .env." + env)
	}
}
