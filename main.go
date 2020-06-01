package main

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"task/controllers"
	"task/helpers"
	"task/models"
	"time"
)

func main() {
	//Establishing MongoDb connection
	Collections := helpers.MongoConnection()
	//Feeding defaultUser data to DB
	defaultUser := os.Getenv("DEFAULT_USER")
	if defaultUser != "" {
		count, _ := Collections["user"].Find(bson.M{"email": defaultUser}).Count()
		if count == 0 {
			u := models.UserModel{
				ID:        bson.NewObjectId(),
				Email:     defaultUser,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := Collections["user"].Insert(u)
			if err != nil {
				fmt.Println("Log -- cannot create default user account", err)
			}
			fmt.Println("Log -- created default user account", err)
		}
	}

	//Creating chiRouter object
	chiRouter := chi.NewRouter()

	//Solving cors issue
	chiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//Defining routes
	chiRouter.Get("/", controllers.HomePage)
	chiRouter.Get("/login", controllers.GoogleLogin)
	chiRouter.Get("/callback", controllers.HandleCallBack)
	chiRouter.Get("/googleLogin", controllers.HandleCallBack)
	chiRouter.With(JwtMiddleware).Get("/getUserList", controllers.GetUserList)
	chiRouter.With(JwtMiddleware).Post("/addUser", controllers.AddUser)
	chiRouter.With(JwtMiddleware).Post("/deleteUser", controllers.DeleteUser)

	chiRouter.With(JwtMiddleware).Get("/logout", controllers.Logout)

	//Static files
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(chiRouter, "/page", filesDir)

	//Port listening
	http.ListenAndServe(":3000", chiRouter)
}

// FileServer is serving static files
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
