package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"sambragge/go-software-solutions/controllers"
	"sambragge/go-software-solutions/core"

	//"google.golang.org/appengine" // FOR PRODUCTION ONLY
	"net/http"
)

func startServer(p string) {

	port := ":" + p
	fmt.Println("=\n=\n=\n Server listening port ", port)
	log.Fatal(http.ListenAndServe(port, nil))

}

func main() {
	//appengine.Main() // FOR PRODUCTION ONLY

}

func init() {

	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer s.Close()

	theCore := core.NewCore(s.DB("go-software-solutions"))
	userController := controllers.NewUserController(theCore)
	postController := controllers.NewPostController(theCore)
	commentController := controllers.NewCommentController(theCore)
	viewController := controllers.NewViewController(theCore)

	log.Print("--- controllers created")

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "views/assets/images/gopher.png") })
	http.Handle("/api/users/", http.StripPrefix("/api/users", userController))
	http.Handle("/api/post/", http.StripPrefix("/api/post", postController))
	http.Handle("/api/comment/", http.StripPrefix("/api/comment", commentController))
	http.Handle("/a/", http.StripPrefix("/a", http.FileServer(http.Dir("views/assets"))))

	http.Handle("/", viewController)

	log.Print("--- routes handled")

	startServer("8080") // FOR DEVELOPMENT ONLY

}
