package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Feature information
type Feature struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

var collection *mongo.Collection

// featureJSON := make(map[string][]Feature)
// err = json.Unmarshal([]byte(reqBody), &featureJSON)
// https://www.thepolyglotdeveloper.com/2019/02/developing-restful-api-golang-mongodb-nosql-database/
func httpServerFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		id := r.URL.Query().Get("id")
		fmt.Println("collection.FindOne: ", id)

		var feature Feature
		objectID, _ := primitive.ObjectIDFromHex(id)
		filter := bson.M{"_id": objectID}
		err := collection.FindOne(context.TODO(), filter).Decode(&feature)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Found a single document: %+v\n", feature)
		json.NewEncoder(w).Encode(feature)

	//see https://www.golangprograms.com/example-to-handle-get-and-post-request-in-golang.html
	case "POST":
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var feature Feature
		if err := json.Unmarshal([]byte(reqBody), &feature); err != nil {
			panic(err)
		}
		insertResult, err := collection.InsertOne(context.TODO(), feature)
		if err != nil {
			log.Fatal(err)
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)

		// For using forms:
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		// if err := r.ParseForm(); err != nil {
		// 	fmt.Fprintf(w, "ParseForm() err: %v", err)
		// 	return
		// }
		// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		// name := r.FormValue("name")
		// address := r.FormValue("address")
		// fmt.Fprintf(w, "Name = %s\n", name)
		// fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {

	fmt.Println("hello, world")
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/go")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection = client.Database("go").Collection("features")

	http.HandleFunc("/", httpServerFunc)
	err = http.ListenAndServe("localhost:8765", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
