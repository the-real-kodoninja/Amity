package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/the-real-kodoninja/Amity/backend/models"
	"github.com/the-real-kodoninja/Amity/backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var userCollection *mongo.Collection
var postCollection *mongo.Collection

func main() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	userCollection = client.Database("amity").Collection("users")
	postCollection = client.Database("amity").Collection("posts")

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/register", register).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/users/{username}", getUser).Methods("GET")
	r.HandleFunc("/users/{username}/update", updateUser).Methods("PUT")
	r.HandleFunc("/feed", getFeed).Methods("GET")
	r.HandleFunc("/explore", getExplore).Methods("GET")
	r.HandleFunc("/posts", createPost).Methods("POST")
	r.HandleFunc("/posts/{id}/like", likePost).Methods("POST")
	r.HandleFunc("/posts/{id}/react", reactToPost).Methods("POST")
	r.HandleFunc("/posts/{id}/share", sharePost).Methods("POST")
	r.HandleFunc("/posts/{id}/hide", hidePost).Methods("POST")
	r.HandleFunc("/posts/{id}/comment", addComment).Methods("POST")
	r.HandleFunc("/posts/{id}/comment/{commentId}/reply", addReply).Methods("POST")
	r.HandleFunc("/shorts", getShorts).Methods("GET")
	r.HandleFunc("/groups", getGroups).Methods("GET")
	r.HandleFunc("/groups", createGroup).Methods("POST")
	r.HandleFunc("/groups/{id}", getGroup).Methods("GET")
	r.HandleFunc("/groups/{id}/join", joinGroup).Methods("POST")
	r.HandleFunc("/groups/{id}/leave", leaveGroup).Methods("POST")
	r.HandleFunc("/pages", getPages).Methods("GET")
	r.HandleFunc("/pages", createPage).Methods("POST")
	r.HandleFunc("/pages/{id}", getPage).Methods("GET")
	r.HandleFunc("/pages/{id}/follow", followPage).Methods("POST")
	r.HandleFunc("/pages/{id}/unfollow", unfollowPage).Methods("POST")
	r.HandleFunc("/users/{username}/unfollow", unfollowUser).Methods("POST")
	r.HandleFunc("/users/{username}/verify", verifyUser).Methods("POST")
	r.HandleFunc("/users/{username}/block", blockUser).Methods("POST")
	r.HandleFunc("/friend-requests", getFriendRequests).Methods("GET")
	r.HandleFunc("/friend-requests", sendFriendRequest).Methods("POST")
	r.HandleFunc("/friend-requests/{id}/accept", acceptFriendRequest).Methods("POST")
	r.HandleFunc("/friend-requests/{id}/reject", rejectFriendRequest).Methods("POST")
	r.HandleFunc("/users/{username}/follow", followUser).Methods("POST")
	r.HandleFunc("/users/{username}/unfollow", unfollowUser).Methods("POST")
	r.HandleFunc("/users/{username}/block", blockUser).Methods("POST")
	r.HandleFunc("/users/{username}/unblock", unblockUser).Methods("POST")
	r.HandleFunc("/notifications", getNotifications).Methods("GET")
	r.HandleFunc("/notifications/{id}/read", markNotificationRead).Methods("POST")
	r.HandleFunc("/messages", sendMessage).Methods("POST")
	r.HandleFunc("/messages/{username}", getMessages).Methods("GET")
	r.HandleFunc("/lists", createList).Methods("POST")
	r.HandleFunc("/lists", getLists).Methods("GET")
	r.HandleFunc("/lists/{id}/add", addToList).Methods("POST")
	r.HandleFunc("/hangouts", createHangout).Methods("POST")
	r.HandleFunc("/hangouts", getHangouts).Methods("GET")
	r.HandleFunc("/hangouts/{id}/join", joinHangout).Methods("POST")
	r.HandleFunc("/hangouts/{id}/leave", leaveHangout).Methods("POST")
	r.HandleFunc("/search/users", searchUsers).Methods("GET")
	// admin functionality endpoints
	r.HandleFunc("/contact-admin", contactAdmin).Methods("POST")
	r.HandleFunc("/admin/messages", getAdminMessages).Methods("GET")
	r.HandleFunc("/users/{username}/ban", banUser).Methods("POST")
	r.HandleFunc("/users/{username}/unban", unbanUser).Methods("POST")
	r.HandleFunc("/posts/{id}/delete", deletePost).Methods("POST")
	r.HandleFunc("/admin/deleted-posts", getDeletedPosts).Methods("GET")
	r.HandleFunc("/admin/sponsored-posts", createSponsoredPost).Methods("POST")
	// File upload endpoint
	r.HandleFunc("/upload", uploadFile).Methods("POST")
	// added functionality
	r.HandleFunc("/users/{username}/pin-post", pinPost).Methods("POST")
	r.HandleFunc("/posts/live", createLivePost).Methods("POST")
	// monetization
	r.HandleFunc("/monetization", getMonetization).Methods("GET")
	r.HandleFunc("/monetization/update", updateMonetization).Methods("POST")

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = primitive.NewObjectID()
	user.Followers = 0
	user.Following = []string{}
	user.Friends = []string{}
	user.BlockedUsers = []string{}
	user.Notifications = []models.Notification{}
	user.Settings = models.UserSettings{
		NSFWEnabled:       false,
		ProfileVisibility: "public",
		Messaging:         "everyone",
	}
	user.Verified = false // Default to false
	user.IsAdmin = false  // Default to false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := generateToken(user.Username)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": creds.Username, "password": creds.Password}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := generateToken(user.Username)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	userJSON := struct {
		Username     string              `json:"username"`
		Location     string              `json:"location"`
		Followers    int                 `json:"followers"`
		Following    []string            `json:"following"`
		Friends      []string            `json:"friends"`
		BlockedUsers []string            `json:"blocked_users"`
		ProfilePhoto string              `json:"profile_photo"`
		BannerPhoto  string              `json:"banner_photo"`
		Settings     models.UserSettings `json:"settings"`
	}{
		Username:     user.Username,
		Location:     user.Location,
		Followers:    user.Followers,
		Following:    user.Following,
		Friends:      user.Friends,
		BlockedUsers: user.BlockedUsers,
		ProfilePhoto: user.ProfilePhoto,
		BannerPhoto:  user.BannerPhoto,
		Settings:     user.Settings,
	}
	json.NewEncoder(w).Encode(userJSON)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	var updates struct {
		Location     string              `json:"location"`
		ProfilePhoto string              `json:"profile_photo"`
		BannerPhoto  string              `json:"banner_photo"`
		Settings     models.UserSettings `json:"settings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateFields := bson.M{}
	if updates.Location != "" {
		updateFields["location"] = updates.Location
	}
	if updates.ProfilePhoto != "" {
		updateFields["profile_photo"] = updates.ProfilePhoto
	}
	if updates.BannerPhoto != "" {
		updateFields["banner_photo"] = updates.BannerPhoto
	}
	if updates.Settings != (models.UserSettings{}) {
		updateFields["settings"] = updates.Settings
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := userCollection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": updateFields})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getFeed(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := bson.M{}
	if filter != "" && filter != "all" {
		query["media.type"] = filter
	}
	cursor, err := postCollection.Find(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func getExplore(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := postCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(struct {
		Posts []models.Post `json:"posts"`
	}{Posts: posts})
}

func likePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	username := r.Header.Get("username") // Extract from JWT in a real app
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postCollection.UpdateOne(ctx, bson.M{"_id": mustObjectID(postID)}, bson.M{"$inc": bson.M{"likes": 1}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify post owner
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": mustObjectID(postID)}).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "like",
		From:      username,
		Content:   fmt.Sprintf("%s liked your post", username),
		RelatedID: postID,
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": post.Username},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func reactToPost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	var reaction struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postCollection.UpdateOne(ctx,
		bson.M{"_id": mustObjectID(postID)},
		bson.M{"$inc": bson.M{"reactions." + reaction.Type: 1}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify post owner
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": mustObjectID(postID)}).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "reaction",
		From:      username,
		Content:   fmt.Sprintf("%s reacted to your post with %s", username, reaction.Type),
		RelatedID: postID,
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": post.Username},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func sharePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postCollection.UpdateOne(ctx, bson.M{"_id": mustObjectID(postID)}, bson.M{"$inc": bson.M{"shares": 1}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify post owner
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": mustObjectID(postID)}).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "share",
		From:      username,
		Content:   fmt.Sprintf("%s shared your post", username),
		RelatedID: postID,
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": post.Username},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func hidePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postCollection.UpdateOne(ctx,
		bson.M{"_id": mustObjectID(postID)},
		bson.M{"$addToSet": bson.M{"hidden_by": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addComment(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment.ID = primitive.NewObjectID()
	comment.Timestamp = time.Now().Format(time.RFC3339)
	comment.Replies = []models.Comment{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postCollection.UpdateOne(ctx,
		bson.M{"_id": mustObjectID(postID)},
		bson.M{"$push": bson.M{"comments": comment}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify post owner
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": mustObjectID(postID)}).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "comment",
		From:      username,
		Content:   fmt.Sprintf("%s commented on your post", username),
		RelatedID: postID,
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": post.Username},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func addReply(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	commentID := mux.Vars(r)["commentId"]
	username := r.Header.Get("username")
	var reply models.Comment
	if err := json.NewDecoder(r.Body).Decode(&reply); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reply.ID = primitive.NewObjectID()
	reply.Timestamp = time.Now().Format(time.RFC3339)
	reply.Replies = []models.Comment{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postCollection.UpdateOne(ctx,
		bson.M{"_id": mustObjectID(postID), "comments._id": mustObjectID(commentID)},
		bson.M{"$push": bson.M{"comments.$.replies": reply}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify post owner and comment author
	var post models.Post
	err = postCollection.FindOne(ctx, bson.M{"_id": mustObjectID(postID)}).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	var comment models.Comment
	for _, c := range post.Comments {
		if c.ID.Hex() == commentID {
			comment = c
			break
		}
	}
	if comment.Username != "" {
		notification := models.Notification{
			ID:        primitive.NewObjectID(),
			Type:      "reply",
			From:      username,
			Content:   fmt.Sprintf("%s replied to your comment", username),
			RelatedID: postID,
			Timestamp: time.Now().Format(time.RFC3339),
			Read:      false,
		}
		_, err = userCollection.UpdateOne(ctx,
			bson.M{"username": comment.Username},
			bson.M{"$push": bson.M{"notifications": notification}},
		)
		if err != nil {
			log.Printf("Error sending notification: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func getShorts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := postCollection.Find(ctx, bson.M{"is_short": true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var shorts []models.Post
	if err := cursor.All(ctx, &shorts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(shorts)
}

func getGroups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := client.Database("amity").Collection("groups").Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var groups []models.Group
	if err := cursor.All(ctx, &groups); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(groups)
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	group.ID = primitive.NewObjectID()
	group.Members = []string{group.Creator}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("groups").InsertOne(ctx, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(group)
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var group models.Group
	err := client.Database("amity").Collection("groups").FindOne(ctx, bson.M{"_id": mustObjectID(groupID)}).Decode(&group)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(group)
}

func joinGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("groups").UpdateOne(ctx,
		bson.M{"_id": mustObjectID(groupID)},
		bson.M{"$addToSet": bson.M{"members": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func leaveGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("groups").UpdateOne(ctx,
		bson.M{"_id": mustObjectID(groupID)},
		bson.M{"$pull": bson.M{"members": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getPages(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := client.Database("amity").Collection("pages").Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var pages []models.Page
	if err := cursor.All(ctx, &pages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pages)
}

func createPage(w http.ResponseWriter, r *http.Request) {
	var page models.Page
	if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	page.ID = primitive.NewObjectID()
	page.Followers = []string{page.Creator}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("pages").InsertOne(ctx, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(page)
}

func getPage(w http.ResponseWriter, r *http.Request) {
	pageID := mux.Vars(r)["id"]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var page models.Page
	err := client.Database("amity").Collection("pages").FindOne(ctx, bson.M{"_id": mustObjectID(pageID)}).Decode(&page)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(page)
}

func followPage(w http.ResponseWriter, r *http.Request) {
	pageID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user is an admin
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Admins can follow any page
	_, err = client.Database("amity").Collection("pages").UpdateOne(ctx,
		bson.M{"_id": mustObjectID(pageID)},
		bson.M{"$addToSet": bson.M{"followers": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func unfollowPage(w http.ResponseWriter, r *http.Request) {
	pageID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user is an admin
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Admins can unfollow any page
	_, err = client.Database("amity").Collection("pages").UpdateOne(ctx,
		bson.M{"_id": mustObjectID(pageID)},
		bson.M{"$pull": bson.M{"followers": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getFriendRequests(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := client.Database("amity").Collection("friend_requests").Find(ctx, bson.M{"to": username})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var requests []models.FriendRequest
	if err := cursor.All(ctx, &requests); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(requests)
}

func sendFriendRequest(w http.ResponseWriter, r *http.Request) {
	var request models.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	request.ID = primitive.NewObjectID()
	request.Timestamp = time.Now().Format(time.RFC3339)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("friend_requests").InsertOne(ctx, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify recipient
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "friend_request",
		From:      request.From,
		Content:   fmt.Sprintf("%s sent you a friend request", request.From),
		RelatedID: request.ID.Hex(),
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": request.To},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func acceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	requestID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var request models.FriendRequest
	err := client.Database("amity").Collection("friend_requests").FindOne(ctx, bson.M{"_id": mustObjectID(requestID)}).Decode(&request)
	if err != nil {
		http.Error(w, "Friend request not found", http.StatusNotFound)
		return
	}

	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": request.From},
		bson.M{"$addToSet": bson.M{"friends": request.To}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": request.To},
		bson.M{"$addToSet": bson.M{"friends": request.From}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = client.Database("amity").Collection("friend_requests").DeleteOne(ctx, bson.M{"_id": mustObjectID(requestID)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify sender
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "friend_accept",
		From:      username,
		Content:   fmt.Sprintf("%s accepted your friend request", username),
		RelatedID: requestID,
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": request.From},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func rejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	requestID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var request models.FriendRequest
	err := client.Database("amity").Collection("friend_requests").FindOne(ctx, bson.M{"_id": mustObjectID(requestID)}).Decode(&request)
	if err != nil {
		http.Error(w, "Friend request not found", http.StatusNotFound)
		return
	}

	_, err = client.Database("amity").Collection("friend_requests").DeleteOne(ctx, bson.M{"_id": mustObjectID(requestID)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify sender
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "friend_reject",
		From:      username,
		Content:   fmt.Sprintf("%s rejected your friend request", username),
		RelatedID: requestID,
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": request.From},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func followUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	follower := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$inc": bson.M{"followers": 1}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": follower},
		bson.M{"$addToSet": bson.M{"following": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify user
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "follow",
		From:      follower,
		Content:   fmt.Sprintf("%s followed you", follower),
		RelatedID: follower,
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func unfollowUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	follower := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$inc": bson.M{"followers": -1}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": follower},
		bson.M{"$pull": bson.M{"following": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func blockUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	blocker := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(ctx,
		bson.M{"username": blocker},
		bson.M{"$addToSet": bson.M{"blocked_users": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func unblockUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	blocker := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(ctx,
		bson.M{"username": blocker},
		bson.M{"$pull": bson.M{"blocked_users": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getNotifications(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user.Notifications)
}

func markNotificationRead(w http.ResponseWriter, r *http.Request) {
	notifID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.UpdateOne(ctx,
		bson.M{"username": username, "notifications._id": mustObjectID(notifID)},
		bson.M{"$set": bson.M{"notifications.$.read": true}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message.ID = primitive.NewObjectID()
	message.Timestamp = time.Now().Format(time.RFC3339)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("messages").InsertOne(ctx, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify recipient
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "message",
		From:      message.From,
		Content:   fmt.Sprintf("%s sent you a message", message.From),
		RelatedID: message.ID.Hex(),
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": message.To},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	currentUser := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := client.Database("amity").Collection("messages").Find(ctx,
		bson.M{
			"$or": []bson.M{
				{"from": currentUser, "to": username},
				{"from": username, "to": currentUser},
			},
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}

func createList(w http.ResponseWriter, r *http.Request) {
	var list models.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list.ID = primitive.NewObjectID()
	list.Items = []models.ListItem{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("lists").InsertOne(ctx, list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func getLists(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := client.Database("amity").Collection("lists").Find(ctx, bson.M{"creator": username})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var lists []models.List
	if err := cursor.All(ctx, &lists); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(lists)
}

func addToList(w http.ResponseWriter, r *http.Request) {
	listID := mux.Vars(r)["id"]
	var item models.ListItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("lists").UpdateOne(ctx,
		bson.M{"_id": mustObjectID(listID)},
		bson.M{"$push": bson.M{"items": item}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createHangout(w http.ResponseWriter, r *http.Request) {
	var hangout models.Hangout
	if err := json.NewDecoder(r.Body).Decode(&hangout); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hangout.ID = primitive.NewObjectID()
	hangout.Participants = []string{hangout.Creator}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("hangouts").InsertOne(ctx, hangout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify friends
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"username": hangout.Creator}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	for _, friend := range user.Friends {
		notification := models.Notification{
			ID:        primitive.NewObjectID(),
			Type:      "hangout",
			From:      hangout.Creator,
			Content:   fmt.Sprintf("%s created a new hangout: %s", hangout.Creator, hangout.Name),
			RelatedID: hangout.ID.Hex(),
			Timestamp: time.Now().Format(time.RFC3339),
			Read:      false,
		}
		_, err := userCollection.UpdateOne(ctx,
			bson.M{"username": friend},
			bson.M{"$push": bson.M{"notifications": notification}},
		)
		if err != nil {
			log.Printf("Error sending notification to %s: %v", friend, err)
		}
	}

	json.NewEncoder(w).Encode(hangout)
}

func getHangouts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := client.Database("amity").Collection("hangouts").Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var hangouts []models.Hangout
	if err := cursor.All(ctx, &hangouts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(hangouts)
}

func joinHangout(w http.ResponseWriter, r *http.Request) {
	hangoutID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("hangouts").UpdateOne(ctx,
		bson.M{"_id": mustObjectID(hangoutID)},
		bson.M{"$addToSet": bson.M{"participants": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func leaveHangout(w http.ResponseWriter, r *http.Request) {
	hangoutID := mux.Vars(r)["id"]
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("hangouts").UpdateOne(ctx,
		bson.M{"_id": mustObjectID(hangoutID)},
		bson.M{"$pull": bson.M{"participants": username}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func searchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := userCollection.Find(ctx, bson.M{"username": bson.M{"$regex": query, "$options": "i"}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB limit
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileType := strings.Split(handler.Header.Get("Content-Type"), "/")[0]
	mediaType := "file"
	if fileType == "image" {
		mediaType = "photo"
	} else if fileType == "video" {
		mediaType = "video"
	}

	key := fmt.Sprintf("uploads/%s-%s", time.Now().Format("20060102150405"), handler.Filename)
	url, err := utils.UploadFileToS3("your-s3-bucket-name", key, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"type": mediaType,
		"url":  url,
		"size": fmt.Sprintf("%d", handler.Size),
	})
}

func generateToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}

func mustObjectID(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return objID
}

func contactAdmin(w http.ResponseWriter, r *http.Request) {
	var message models.AdminMessage
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message.ID = primitive.NewObjectID()
	message.Timestamp = time.Now().Format(time.RFC3339)
	message.Read = false
	message.From = r.Header.Get("username")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("admin_messages").InsertOne(ctx, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getAdminMessages(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")

	// Check if user is an admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
		return
	}

	cursor, err := client.Database("amity").Collection("admin_messages").Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var messages []models.AdminMessage
	if err := cursor.All(ctx, &messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(post.Content) > 280 {
		http.Error(w, "Post exceeds 280 characters", http.StatusBadRequest)
		return
	}

	post.ID = primitive.NewObjectID()
	post.Timestamp = time.Now().Format(time.RFC3339)
	post.Likes = 0
	post.Reactions = make(map[string]int)
	post.Shares = 0
	post.Comments = []models.Comment{}
	post.HiddenBy = []string{}
	post.Deleted = false
	post.Sponsored = false
	post.Live = false
	// Simplified NFT logic
	post.NFTAddress = fmt.Sprintf("0x%s", post.ID.Hex()) // Mock NFT address
	post.MintEarnings = 0.0                              // Initial earnings

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postCollection.InsertOne(ctx, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify followers
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"username": post.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	for _, follower := range user.Following {
		notification := models.Notification{
			ID:        primitive.NewObjectID(),
			Type:      "post",
			From:      post.Username,
			Content:   fmt.Sprintf("%s created a new post", post.Username),
			RelatedID: post.ID.Hex(),
			Timestamp: time.Now().Format(time.RFC3339),
			Read:      false,
		}
		_, err := userCollection.UpdateOne(ctx,
			bson.M{"username": follower},
			bson.M{"$push": bson.M{"notifications": notification}},
		)
		if err != nil {
			log.Printf("Error sending notification to %s: %v", follower, err)
		}
	}

	json.NewEncoder(w).Encode(post)
}

func banUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	requester := r.Header.Get("username")

	// Check if requester is an admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var requesterUser models.User
	err := userCollection.FindOne(ctx, bson.M{"username": requester}).Decode(&requesterUser)
	if err != nil {
		http.Error(w, "Requester not found", http.StatusNotFound)
		return
	}
	if !requesterUser.IsAdmin {
		http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
		return
	}

	// Ban the user
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"banned": true}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func unbanUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	requester := r.Header.Get("username")

	// Check if requester is an admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var requesterUser models.User
	err := userCollection.FindOne(ctx, bson.M{"username": requester}).Decode(&requesterUser)
	if err != nil {
		http.Error(w, "Requester not found", http.StatusNotFound)
		return
	}
	if !requesterUser.IsAdmin {
		http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
		return
	}

	// Unban the user
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"banned": false}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	requester := r.Header.Get("username")

	// Check if requester is an admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var requesterUser models.User
	err := userCollection.FindOne(ctx, bson.M{"username": requester}).Decode(&requesterUser)
	if err != nil {
		http.Error(w, "Requester not found", http.StatusNotFound)
		return
	}
	if !requesterUser.IsAdmin {
		http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
		return
	}

	// Mark post as deleted
	_, err = postCollection.UpdateOne(ctx,
		bson.M{"_id": mustObjectID(postID)},
		bson.M{"$set": bson.M{"deleted": true}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getDeletedPosts(w http.ResponseWriter, r *http.Request) {
	requester := r.Header.Get("username")

	// Check if requester is an admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var requesterUser models.User
	err := userCollection.FindOne(ctx, bson.M{"username": requester}).Decode(&requesterUser)
	if err != nil {
		http.Error(w, "Requester not found", http.StatusNotFound)
		return
	}
	if !requesterUser.IsAdmin {
		http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
		return
	}

	cursor, err := postCollection.Find(ctx, bson.M{"deleted": true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts)
}

func createSponsoredPost(w http.ResponseWriter, r *http.Request) {
	requester := r.Header.Get("username")

	// Check if requester is an admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var requesterUser models.User
	err := userCollection.FindOne(ctx, bson.M{"username": requester}).Decode(&requesterUser)
	if err != nil {
		http.Error(w, "Requester not found", http.StatusNotFound)
		return
	}
	if !requesterUser.IsAdmin {
		http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(post.Content) > 280 {
		http.Error(w, "Post exceeds 280 characters", http.StatusBadRequest)
		return
	}

	post.ID = primitive.NewObjectID()
	post.Timestamp = time.Now().Format(time.RFC3339)
	post.Likes = 0
	post.Reactions = make(map[string]int)
	post.Shares = 0
	post.Comments = []models.Comment{}
	post.HiddenBy = []string{}
	post.Deleted = false
	post.Sponsored = true
	// NFT logic
	// In createSponsoredPost
	post.NFTAddress = fmt.Sprintf("0x%s", post.ID.Hex())
	post.MintEarnings = 0.0

	_, err = postCollection.InsertOne(ctx, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func createLivePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post.ID = primitive.NewObjectID()
	post.Timestamp = time.Now().Format(time.RFC3339)
	post.Likes = 0
	post.Reactions = make(map[string]int)
	post.Shares = 0
	post.Comments = []models.Comment{}
	post.HiddenBy = []string{}
	post.Deleted = false
	post.Sponsored = false
	post.Live = true
	//NFT logic
	// In createLivePost
	post.NFTAddress = fmt.Sprintf("0x%s", post.ID.Hex())
	post.MintEarnings = 0.0

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := postCollection.InsertOne(ctx, post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify followers
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"username": post.Username}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	for _, follower := range user.Following {
		notification := models.Notification{
			ID:        primitive.NewObjectID(),
			Type:      "live",
			From:      post.Username,
			Content:   fmt.Sprintf("%s is live!", post.Username),
			RelatedID: post.ID.Hex(),
			Timestamp: time.Now().Format(time.RFC3339),
			Read:      false,
		}
		_, err := userCollection.UpdateOne(ctx,
			bson.M{"username": follower},
			bson.M{"$push": bson.M{"notifications": notification}},
		)
		if err != nil {
			log.Printf("Error sending notification to %s: %v", follower, err)
		}
	}

	json.NewEncoder(w).Encode(post)
}

func pinPost(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	requester := r.Header.Get("username")
	var body struct {
		PostID string `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if requester is the user
	if requester != username {
		http.Error(w, "Unauthorized: Can only pin your own posts", http.StatusUnauthorized)
		return
	}

	// Verify the post exists and belongs to the user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var post models.Post
	err := postCollection.FindOne(ctx, bson.M{"_id": mustObjectID(body.PostID), "username": username}).Decode(&post)
	if err != nil {
		http.Error(w, "Post not found or not owned by user", http.StatusNotFound)
		return
	}

	// Pin the post
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"pinned_post_id": body.PostID}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getMonetization(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var monetization models.Monetization
	err := client.Database("amity").Collection("monetization").FindOne(ctx, bson.M{"username": username}).Decode(&monetization)
	if err != nil {
		// If not found, return zero earnings
		monetization = models.Monetization{
			Username:      username,
			TotalEarnings: 0.0,
			AdEarnings:    0.0,
			NFTEarnings:   0.0,
		}
	}
	json.NewEncoder(w).Encode(monetization)
}

func updateMonetization(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	var update struct {
		AdEarnings  float64 `json:"ad_earnings"`
		NFTEarnings float64 `json:"nft_earnings"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("amity").Collection("monetization").UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{
			"$set": bson.M{
				"ad_earnings":    update.AdEarnings,
				"nft_earnings":   update.NFTEarnings,
				"total_earnings": update.AdEarnings + update.NFTEarnings,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update user's total NFT earnings
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"total_nft_earnings": update.NFTEarnings}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func verifyUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	requester := r.Header.Get("username") // Assumes username is set in the header via JWT

	// Check if requester is an admin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var requesterUser models.User
	err := userCollection.FindOne(ctx, bson.M{"username": requester}).Decode(&requesterUser)
	if err != nil {
		http.Error(w, "Requester not found", http.StatusNotFound)
		return
	}
	if !requesterUser.IsAdmin {
		http.Error(w, "Unauthorized: Admin access required", http.StatusUnauthorized)
		return
	}

	// Verify the user
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"verified": true}},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify the user
	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		Type:      "verified",
		From:      "system",
		Content:   "Your account has been verified!",
		RelatedID: username,
		Timestamp: time.Now().Format(time.RFC3339),
		Read:      false,
	}
	_, err = userCollection.UpdateOne(ctx,
		bson.M{"username": username},
		bson.M{"$push": bson.M{"notifications": notification}},
	)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	w.WriteHeader(http.StatusOK)
}
