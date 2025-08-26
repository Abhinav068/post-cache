package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var posts = make(map[int64]*Post)

func main() {
	http.HandleFunc("/post", createPostHandler)
	http.HandleFunc("/post/", updatePostHandler)
	http.HandleFunc("/posts", getAllPostsHandler)
	http.HandleFunc("/like/", AddLike)
	http.HandleFunc("/dislike/", Dislike)
	http.ListenAndServe(":8088", nil)
}

var mongoClient *mongo.Client
var postsCollection *mongo.Collection

func init() {
	var err error
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://abhinav:Abhinav@cluster0.9q0f4.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0&authMechanism=SCRAM-SHA-1"))
	if err != nil {
		panic(err)
	}
	postsCollection = mongoClient.Database("mydb").Collection("posts")
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postData Post
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id := time.Now().UnixMilli()
	posts[id] = &postData

	// Write to MongoDB
	_, err := postsCollection.InsertOne(context.Background(), postData)
	if err != nil {
		http.Error(w, "Failed to write to DB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Post created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	idParam := r.URL.Path[len("/post/"):]
	ID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var postData Post
	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	posts[ID] = &postData

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "Post updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func AddLike(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/like/"):]
	ID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, ok := posts[ID]
	if !ok {
		http.Error(w, "post not found", http.StatusBadRequest)
		return
	}
	post.Lock()
	defer post.Unlock()
	post.Likes++

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "likes updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}
func Dislike(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/dislike/"):]
	ID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, ok := posts[ID]
	if !ok {
		http.Error(w, "post not found", http.StatusBadRequest)
		return
	}
	post.Lock()
	defer post.Unlock()
	post.Likes--

	w.WriteHeader(http.StatusCreated)
	response := map[string]string{
		"message": "likes updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}
