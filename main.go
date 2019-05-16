package main

import (
	"fmt"
	"grustosoftware/controllers"
	"grustosoftware/core"
	"log"

	"gopkg.in/mgo.v2"

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

	theCore := core.NewCore(s.DB("grusto"))
	userController := controllers.NewUserController(theCore)
	boardController := controllers.NewBoardController(theCore)
	postController := controllers.NewPostController(theCore)
	commentController := controllers.NewCommentController(theCore)
	voteController := controllers.NewVoteController(theCore)
	viewController := controllers.NewViewController(theCore)
	projectController := controllers.NewProjectController(theCore)

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "views/assets/images/gopher.png") })
	http.Handle("/api/users/", http.StripPrefix("/api/users", userController))
	http.Handle("/api/boards/", http.StripPrefix("/api/boards", boardController))
	http.Handle("/api/posts/", http.StripPrefix("/api/posts", postController))
	http.Handle("/api/comments/", http.StripPrefix("/api/comments", commentController))
	http.Handle("/api/votes/", http.StripPrefix("/api/votes", voteController))
	http.Handle("/api/projects/", http.StripPrefix("/api/projects", projectController))
	http.Handle("/a/", http.StripPrefix("/a", http.FileServer(http.Dir("views/assets/"))))

	http.Handle("/", viewController)

	startServer("8080") // FOR DEVELOPMENT ONLY

}

func serveView(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "rview/build/index.html")
}
