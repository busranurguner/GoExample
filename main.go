package main

import (
	"context"
	"encoding/json"
	"fmt"

	//"io/ioutil"
	"log"
	"net/http" //HTTP istemci ve sunucu uygulamaları sağlar -> get/put/post

	_ "api/docs" // docs is generated by Swag CLI, you have to import it.

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Article struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id"`
	Title   string             `json:"Title" bson:"Title"`
	Desc    string             `json:"desc" bson:"desc"`
	Content string             `json:"content" bson:"content"`
}

//var Articles []Article
var client *mongo.Client
var ctx = context.TODO()

func handleRequests() {

	//connection
	clientOptions := options.Client().ApplyURI("YOUR_MONGO_URI_HERE") 
	client, _ = mongo.Connect(ctx, clientOptions)

	// creates a new instance of a mux router
	myRouter := mux.NewRouter()

	//endpoints
	myRouter.HandleFunc("/article", getAllArticles).Methods("GET")
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")

	//swagger
	myRouter.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(":5050", myRouter))
}

// GetArticle godoc
// @Summary Get details of all article
// @Description Get details of all articles
// @Tags articles
// @Accept  json
// @Produce  json
// @Success 200 {array} Article
// @Router /article [get]
func getAllArticles(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var art []Article
	collection := client.Database("Example").Collection("Articles")
	cursor, _ := collection.Find(ctx, bson.M{})

	for cursor.Next(ctx) {
		var article Article
		cursor.Decode(&article)
		art = append(art, article)
	}

	json.NewEncoder(w).Encode(art)
}

// CreateOrder godoc
// @Summary Create a new article
// @Description Create a new article with the input paylod
// @Tags articles
// @Accept  json
// @Produce  json
// @Param article body Article true "Create article"
// @Success 200 {object} Article
// @Router /article [post]
func createNewArticle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var article Article
	_ = json.NewDecoder(r.Body).Decode(&article)
	collection := client.Database("Example").Collection("Articles")
	result, _ := collection.InsertOne(ctx, article)
	json.NewEncoder(w).Encode(result) // JSON dizesine kodlama ve ardından yanıtımızın bir parçası olarak yazma

}

// @title Article API
// @version 1.0
// @description
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email hello@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5050
// @BasePath /
func main() {

	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
