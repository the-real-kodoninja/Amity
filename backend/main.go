package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
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

    // File upload endpoint
    r.HandleFunc("/upload", uploadFile).Methods("POST")

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
        Username     string          `json:"username"`
        Location     string          `json:"location"`
        Followers    int             `json:"followers"`
        Following    []string        `json:"following"`
        Friends      []string        `json:"friends"`
        BlockedUsers []string        `json:"blocked_users"`
        ProfilePhoto string          `json:"profile_photo"`
        BannerPhoto  string          `json:"banner_photo"`
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
        Location     string             `json:"location"`
        ProfilePhoto string             `json:"profile_photo"`
        BannerPhoto  string             `json:"banner_photo"`
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
    _, err := client.Database("amity").Collection("pages").UpdateOne(ctx,
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
    _, err := client.Database("amity").Collection("pages").UpdateOne(ctx,
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
        bson.M{"$inc