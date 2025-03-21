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
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/the-real-kodoninja/Amity/backend/models"
)

type App struct {
    Router *mux.Router
    DB     *mongo.Client
}

var jwtSecret = []byte("your-secret-key")

func (a *App) Initialize() {
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }

    a.DB = client
    a.Router = mux.NewRouter()
    a.initializeRoutes()
}

func (a *App) initializeRoutes() {
    fs := http.FileServer(http.Dir("../frontend/build"))
    a.Router.PathPrefix("/").Handler(http.StripPrefix("/", fs)).Methods("GET")

    a.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Amity Backend is running!")
    }).Methods("GET")
    a.Router.HandleFunc("/login", a.login).Methods("POST")
    a.Router.HandleFunc("/register", a.register).Methods("POST")

    api := a.Router.PathPrefix("/api").Subrouter()
    api.Use(a.jwtMiddleware)
    api.HandleFunc("/users/{username}", a.getUser).Methods("GET")
    api.HandleFunc("/users/{username}/update", a.updateUser).Methods("PUT")
    api.HandleFunc("/feed", a.getFeed).Methods("GET")
    api.HandleFunc("/explore", a.getExplore).Methods("GET")
    api.HandleFunc("/users/{username}/photos", a.getUserPhotos).Methods("GET")
    api.HandleFunc("/posts", a.createPost).Methods("POST")
    api.HandleFunc("/posts/{id}/like", a.likePost).Methods("POST")
    api.HandleFunc("/posts/{id}/share", a.sharePost).Methods("POST")
    api.HandleFunc("/posts/{id}/comment", a.addComment).Methods("POST")
    api.HandleFunc("/posts/{id}/mint", a.mintPost).Methods("POST")
    api.HandleFunc("/photos/{id}/mint", a.mintPhoto).Methods("POST")
    api.HandleFunc("/search/users", a.searchUsers).Methods("GET")
    api.HandleFunc("/groups", a.createGroup).Methods("POST")
    api.HandleFunc("/groups", a.getGroups).Methods("GET")
    api.HandleFunc("/groups/{id}", a.getGroup).Methods("GET")
    api.HandleFunc("/groups/{id}/join", a.joinGroup).Methods("POST")
    api.HandleFunc("/groups/{id}/leave", a.leaveGroup).Methods("POST")
    api.HandleFunc("/pages", a.createPage).Methods("POST")
    api.HandleFunc("/pages", a.getPages).Methods("GET")
    api.HandleFunc("/pages/{id}", a.getPage).Methods("GET")
    api.HandleFunc("/pages/{id}/follow", a.followPage).Methods("POST")
    api.HandleFunc("/pages/{id}/unfollow", a.unfollowPage).Methods("POST")
    api.HandleFunc("/friend-requests", a.sendFriendRequest).Methods("POST")
    api.HandleFunc("/friend-requests", a.getFriendRequests).Methods("GET")
    api.HandleFunc("/friend-requests/{id}/accept", a.acceptFriendRequest).Methods("POST")
    api.HandleFunc("/friend-requests/{id}/reject", a.rejectFriendRequest).Methods("POST")
    api.HandleFunc("/users/{username}/follow", a.followUser).Methods("POST")
    api.HandleFunc("/users/{username}/unfollow", a.unfollowUser).Methods("POST")
    api.HandleFunc("/users/{username}/block", a.blockUser).Methods("POST")
    api.HandleFunc("/users/{username}/unblock", a.unblockUser).Methods("POST")
    api.HandleFunc("/notifications", a.getNotifications).Methods("GET")
    api.HandleFunc("/notifications/{id}/read", a.markNotificationRead).Methods("POST")
    api.HandleFunc("/messages", a.sendMessage).Methods("POST")
    api.HandleFunc("/messages/{username}", a.getMessages).Methods("GET")
    api.HandleFunc("/lists", a.createList).Methods("POST")
    api.HandleFunc("/lists", a.getLists).Methods("GET")
    api.HandleFunc("/lists/{id}/add", a.addToList).Methods("POST")
    api.HandleFunc("/hangouts", a.createHangout).Methods("POST")
    api.HandleFunc("/hangouts", a.getHangouts).Methods("GET")
    api.HandleFunc("/hangouts/{id}/join", a.joinHangout).Methods("POST")
    api.HandleFunc("/hangouts/{id}/leave", a.leaveHangout).Methods("POST")

    a.Router.HandleFunc("/.well-known/webfinger", a.handleWebFinger).Methods("GET")
    a.Router.HandleFunc("/users/{username}/outbox", a.handleOutbox).Methods("GET")
    a.Router.HandleFunc("/users/{username}/inbox", a.handleInbox).Methods("POST")
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid credentials", http.StatusBadRequest)
        return
    }

    collection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := collection.FindOne(ctx, bson.M{"username": creds.Username}).Decode(&user)
    if err != nil || user.Password != creds.Password { // In production, use password hashing
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": creds.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid user data", http.StatusBadRequest)
        return
    }

    user.Followers = 0
    user.Following = []string{}
    user.Friends = []string{}
    user.BlockedUsers = []string{}
    user.ProfilePhoto = "https://via.placeholder.com/100"
    user.Banner = "https://via.placeholder.com/1500x500"
    user.Settings = models.UserSettings{
        NSFWEnabled: false,
        Privacy: struct {
            ProfileVisibility string `json:"profile_visibility" bson:"profile_visibility"`
            Messaging         string `json:"messaging" bson:"messaging"`
        }{
            ProfileVisibility: "public",
            Messaging:         "everyone",
        },
    }

    collection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := collection.InsertOne(ctx, user)
    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (a *App) jwtMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        if !strings.HasPrefix(authHeader, "Bearer ") {
            http.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if ok && token.Valid {
            ctx := context.WithValue(r.Context(), "username", claims["username"])
            r = r.WithContext(ctx)
        }

        next.ServeHTTP(w, r)
    })
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]

    collection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]

    var updates struct {
        ProfilePhoto string           `json:"profile_photo"`
        Banner       string           `json:"banner"`
        Settings     models.UserSettings `json:"settings"`
    }
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        http.Error(w, "Invalid update data", http.StatusBadRequest)
        return
    }

    collection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    update := bson.M{"$set": bson.M{}}
    if updates.ProfilePhoto != "" {
        update["$set"].(bson.M)["profile_photo"] = updates.ProfilePhoto
    }
    if updates.Banner != "" {
        update["$set"].(bson.M)["banner"] = updates.Banner
    }
    if updates.Settings != (models.UserSettings{}) {
        update["$set"].(bson.M)["settings"] = updates.Settings
    }

    result, err := collection.UpdateOne(ctx, bson.M{"username": username}, update)
    if err != nil {
        http.Error(w, "Error updating user", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("User updated"))
}

func (a *App) getFeed(w http.ResponseWriter, r *http.Request) {
    collection := a.DB.Database("amity").Collection("posts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Error fetching feed", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var posts []models.Post
    if err := cursor.All(ctx, &posts); err != nil {
        http.Error(w, "Error decoding posts", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func (a *App) getExplore(w http.ResponseWriter, r *http.Request) {
    // Fetch trending posts, groups, pages, and users (simplified)
    postsCollection := a.DB.Database("amity").Collection("posts")
    groupsCollection := a.DB.Database("amity").Collection("groups")
    pagesCollection := a.DB.Database("amity").Collection("pages")
    usersCollection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var posts []models.Post
    postsCursor, err := postsCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"likes": -1}).SetLimit(10))
    if err == nil {
        postsCursor.All(ctx, &posts)
        postsCursor.Close(ctx)
    }

    var groups []models.Group
    groupsCursor, err := groupsCollection.Find(ctx, bson.M{}, options.Find().SetLimit(5))
    if err == nil {
        groupsCursor.All(ctx, &groups)
        groupsCursor.Close(ctx)
    }

    var pages []models.Page
    pagesCursor, err := pagesCollection.Find(ctx, bson.M{}, options.Find().SetLimit(5))
    if err == nil {
        pagesCursor.All(ctx, &pages)
        pagesCursor.Close(ctx)
    }

    var users []models.User
    usersCursor, err := usersCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"followers": -1}).SetLimit(5))
    if err == nil {
        usersCursor.All(ctx, &users)
        usersCursor.Close(ctx)
    }

    response := map[string]interface{}{
        "posts":  posts,
        "groups": groups,
        "pages":  pages,
        "users":  users,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (a *App) getUserPhotos(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]

    collection := a.DB.Database("amity").Collection("photos")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"username": username})
    if err != nil {
        http.Error(w, "Error fetching photos", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var photos []models.Photo
    if err := cursor.All(ctx, &photos); err != nil {
        http.Error(w, "Error decoding photos", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(photos)
}

func (a *App) createPost(w http.ResponseWriter, r *http.Request) {
    var post models.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid post data", http.StatusBadRequest)
        return
    }

    post.ActivityID = fmt.Sprintf("https://amity.example.com/posts/%d", time.Now().UnixNano())
    post.ActivityType = "Create"
    post.Likes = 0
    post.Shares = 0
    post.Comments = []models.Comment{}

    collection := a.DB.Database("amity").Collection("posts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, post)
    if err != nil {
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    log.Printf("Created post with ID: %v", result.InsertedID)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

func (a *App) likePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]

    collection := a.DB.Database("amity").Collection("posts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var post models.Post
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$inc": bson.M{"likes": 1}},
    )
    if err != nil {
        http.Error(w, "Error liking post", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    // Create a notification
    username := r.Context().Value("username").(string)
    a.createNotification(post.Username, "like", username, fmt.Sprintf("%s liked your post", username), postID)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Post liked"))
}

func (a *App) sharePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]

    collection := a.DB.Database("amity").Collection("posts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var post models.Post
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$inc": bson.M{"shares": 1}},
    )
    if err != nil {
        http.Error(w, "Error sharing post", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    username := r.Context().Value("username").(string)
    a.createNotification(post.Username, "share", username, fmt.Sprintf("%s shared your post", username), postID)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Post shared"))
}

func (a *App) addComment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]

    var comment models.Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        http.Error(w, "Invalid comment data", http.StatusBadRequest)
        return
    }

    comment.Timestamp = time.Now().Format(time.RFC3339)

    collection := a.DB.Database("amity").Collection("posts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var post models.Post
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$push": bson.M{"comments": comment}},
    )
    if err != nil {
        http.Error(w, "Error adding comment", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    username := r.Context().Value("username").(string)
    a.createNotification(post.Username, "comment", username, fmt.Sprintf("%s commented on your post", username), postID)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Comment added"))
}

func (a *App) mintPost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID := vars["id"]

    // Simulate minting a post as an NFT (requires blockchain integration in production)
    collection := a.DB.Database("amity").Collection("posts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var post models.Post
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    // In a real app, integrate with a blockchain (e.g., Internet Computer) to mint the post as an NFT
    log.Printf("Minted post %s as NFT (simulated)", postID)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Post minted as NFT"))
}

func (a *App) mintPhoto(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    photoID := vars["id"]

    // Simulate minting a photo as an NFT
    collection := a.DB.Database("amity").Collection("photos")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(photoID)
    if err != nil {
        http.Error(w, "Invalid photo ID", http.StatusBadRequest)
        return
    }

    var photo models.Photo
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&photo)
    if err != nil {
        http.Error(w, "Photo not found", http.StatusNotFound)
        return
    }

    log.Printf("Minted photo %s as NFT (simulated)", photoID)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Photo minted as NFT"))
}

func (a *App) searchUsers(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    if query == "" {
        http.Error(w, "Missing query parameter", http.StatusBadRequest)
        return
    }

    collection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    regex := bson.M{"username": bson.M{"$regex": query, "$options": "i"}}
    cursor, err := collection.Find(ctx, regex)
    if err != nil {
        http.Error(w, "Error searching users", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var users []models.User
    if err := cursor.All(ctx, &users); err != nil {
        http.Error(w, "Error decoding users", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func (a *App) createGroup(w http.ResponseWriter, r *http.Request) {
    var group models.Group
    if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
        http.Error(w, "Invalid group data", http.StatusBadRequest)
        return
    }

    username := r.Context().Value("username").(string)
    group.Creator = username
    group.Members = []string{username}
    group.CreatedAt = time.Now().Format(time.RFC3339)

    collection := a.DB.Database("amity").Collection("groups")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, group)
    if err != nil {
        http.Error(w, "Error creating group", http.StatusInternalServerError)
        return
    }

    groupID := result.InsertedID.(primitive.ObjectID).Hex()
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(group)
}

func (a *App) getGroups(w http.ResponseWriter, r *http.Request) {
    collection := a.DB.Database("amity").Collection("groups")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Error fetching groups", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var groups []models.Group
    if err := cursor.All(ctx, &groups); err != nil {
        http.Error(w, "Error decoding groups", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(groups)
}

func (a *App) getGroup(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    groupID := vars["id"]

    collection := a.DB.Database("amity").Collection("groups")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(groupID)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    var group models.Group
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&group)
    if err != nil {
        http.Error(w, "Group not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(group)
}

func (a *App) joinGroup(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    groupID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("groups")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(groupID)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$addToSet": bson.M{"members": username}},
    )
    if err != nil {
        http.Error(w, "Error joining group", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Group not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Joined group"))
}

func (a *App) leaveGroup(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    groupID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("groups")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(groupID)
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$pull": bson.M{"members": username}},
    )
    if err != nil {
        http.Error(w, "Error leaving group", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Group not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Left group"))
}

func (a *App) createPage(w http.ResponseWriter, r *http.Request) {
    var page models.Page
    if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
        http.Error(w, "Invalid page data", http.StatusBadRequest)
        return
    }

    username := r.Context().Value("username").(string)
    page.Creator = username
    page.Followers = []string{username}
    page.CreatedAt = time.Now().Format(time.RFC3339)

    collection := a.DB.Database("amity").Collection("pages")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, page)
    if err != nil {
        http.Error(w, "Error creating page", http.StatusInternalServerError)
        return
    }

    pageID := result.InsertedID.(primitive.ObjectID).Hex()
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(page)
}

func (a *App) getPages(w http.ResponseWriter, r *http.Request) {
    collection := a.DB.Database("amity").Collection("pages")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Error fetching pages", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var pages []models.Page
    if err := cursor.All(ctx, &pages); err != nil {
        http.Error(w, "Error decoding pages", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(pages)
}

func (a *App) getPage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    pageID := vars["id"]

    collection := a.DB.Database("amity").Collection("pages")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(pageID)
    if err != nil {
        http.Error(w, "Invalid page ID", http.StatusBadRequest)
        return
    }

    var page models.Page
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&page)
    if err != nil {
        http.Error(w, "Page not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(page)
}

func (a *App) followPage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    pageID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("pages")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(pageID)
    if err != nil {
        http.Error(w, "Invalid page ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$addToSet": bson.M{"followers": username}},
    )
    if err != nil {
        http.Error(w, "Error following page", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Page not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Followed page"))
}

func (a *App) unfollowPage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    pageID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("pages")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(pageID)
    if err != nil {
        http.Error(w, "Invalid page ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$pull": bson.M{"followers": username}},
    )
    if err != nil {
        http.Error(w, "Error unfollowing page", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Page not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Unfollowed page"))
}

func (a *App) sendFriendRequest(w http.ResponseWriter, r *http.Request) {
    var request models.FriendRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid friend request data", http.StatusBadRequest)
        return
    }

    username := r.Context().Value("username").(string)
    request.From = username
    request.Status = "pending"
    request.CreatedAt = time.Now().Format(time.RFC3339)

    collection := a.DB.Database("amity").Collection("friend_requests")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, request)
    if err != nil {
        http.Error(w, "Error sending friend request", http.StatusInternalServerError)
        return
    }

    a.createNotification(request.To, "friend_request", username, fmt.Sprintf("%s sent you a friend request", username), result.InsertedID.(primitive.ObjectID).Hex())

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(request)
}

func (a *App) getFriendRequests(w http.ResponseWriter, r *http.Request) {
    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("friend_requests")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"to": username, "status": "pending"})
    if err != nil {
        http.Error(w, "Error fetching friend requests", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var requests []models.FriendRequest
    if err := cursor.All(ctx, &requests); err != nil {
        http.Error(w, "Error decoding friend requests", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(requests)
}

func (a *App) acceptFriendRequest(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    requestID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("friend_requests")
    usersCollection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(requestID)
    if err != nil {
        http.Error(w, "Invalid request ID", http.StatusBadRequest)
        return
    }

    var request models.FriendRequest
    err = collection.FindOne(ctx, bson.M{"_id": objectID, "to": username}).Decode(&request)
    if err != nil {
        http.Error(w, "Friend request not found", http.StatusNotFound)
        return
    }

    // Update the request status
    _, err = collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$set": bson.M{"status": "accepted"}},
    )
    if err != nil {
        http.Error(w, "Error accepting friend request", http.StatusInternalServerError)
        return
    }

    // Add to friends list for both users
    _, err = usersCollection.UpdateOne(
        ctx,
        bson.M{"username": username},
        bson.M{"$addToSet": bson.M{"friends": request.From}},
    )
    if err != nil {
        http.Error(w, "Error updating friends list", http.StatusInternalServerError)
        return
    }

    _, err = usersCollection.UpdateOne(
        ctx,
        bson.M{"username": request.From},
        bson.M{"$addToSet": bson.M{"friends": username}},
    )
    if err != nil {
        http.Error(w, "Error updating friends list", http.StatusInternalServerError)
        return
    }

    a.createNotification(request.From, "friend_accepted", username, fmt.Sprintf("%s accepted your friend request", username), requestID)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Friend request accepted"))
}

func (a *App) rejectFriendRequest(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    requestID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("friend_requests")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(requestID)
    if err != nil {
        http.Error(w, "Invalid request ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID, "to": username},
        bson.M{"$set": bson.M{"status": "rejected"}},
    )
    if err != nil {
        http.Error(w, "Error rejecting friend request", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Friend request not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Friend request rejected"))
}

func (a *App) followUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    usernameToFollow := vars["username"]

    username := r.Context().Value("username").(string)
    if username == usernameToFollow {
        http.Error(w, "Cannot follow yourself", http.StatusBadRequest)
        return
    }

    usersCollection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Add to following list
    _, err := usersCollection.UpdateOne(
        ctx,
        bson.M{"username": username},
        bson.M{"$addToSet": bson.M{"following": usernameToFollow}},
    )
    if err != nil {
        http.Error(w, "Error following user", http.StatusInternalServerError)
        return
    }

    // Increment followers count
    result, err := usersCollection.UpdateOne(
        ctx,
        bson.M{"username": usernameToFollow},
        bson.M{"$inc": bson.M{"followers": 1}},
    )
    if err != nil {
        http.Error(w, "Error updating followers", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    a.createNotification(usernameToFollow, "follow", username, fmt.Sprintf("%s followed you", username), "")

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Followed user"))
}

func (a *App) unfollowUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    usernameToUnfollow := vars["username"]

    username := r.Context().Value("username").(string)
    if username == usernameToUnfollow {
        http.Error(w, "Cannot unfollow yourself", http.StatusBadRequest)
        return
    }

    usersCollection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Remove from following list
    _, err := usersCollection.UpdateOne(
        ctx,
        bson.M{"username": username},
        bson.M{"$pull": bson.M{"following": usernameToUnfollow}},
    )
    if err != nil {
        http.Error(w, "Error unfollowing user", http.StatusInternalServerError)
        return
    }

    // Decrement followers count
    result, err := usersCollection.UpdateOne(
        ctx,
        bson.M{"username": usernameToUnfollow},
        bson.M{"$inc": bson.M{"followers": -1}},
    )
    if err != nil {
        http.Error(w, "Error updating followers", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Unfollowed user"))
}

func (a *App) blockUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    usernameToBlock := vars["username"]

    username := r.Context().Value("username").(string)
    if username == usernameToBlock {
        http.Error(w, "Cannot block yourself", http.StatusBadRequest)
        return
    }

    usersCollection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := usersCollection.UpdateOne(
        ctx,
        bson.M{"username": username},
        bson.M{"$addToSet": bson.M{"blocked_users": usernameToBlock}},
    )
    if err != nil {
        http.Error(w, "Error blocking user", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Blocked user"))
}

func (a *App) unblockUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    usernameToUnblock := vars["username"]

    username := r.Context().Value("username").(string)
    if username == usernameToUnblock {
        http.Error(w, "Cannot unblock yourself", http.StatusBadRequest)
        return
    }

    usersCollection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := usersCollection.UpdateOne(
        ctx,
        bson.M{"username": username},
        bson.M{"$pull": bson.M{"blocked_users": usernameToUnblock}},
    )
    if err != nil {
        http.Error(w, "Error unblocking user", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Unblocked user"))
}

func (a *App) createNotification(userID, notifType, from, message, relatedID string) {
    notification := models.Notification{
        UserID:    userID,
        Type:      notifType,
        From:      from,
        Message:   message,
        RelatedID: relatedID,
        IsRead:    false,
        CreatedAt: time.Now().Format(time.RFC3339),
    }

    collection := a.DB.Database("amity").Collection("notifications")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := collection.InsertOne(ctx, notification)
    if err != nil {
        log.Printf("Error creating notification: %v", err)
    }
}

func (a *App) getNotifications(w http.ResponseWriter, r *http.Request) {
    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("notifications")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"user_id": username})
    if err != nil {
        http.Error(w, "Error fetching notifications", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var notifications []models.Notification
    if err := cursor.All(ctx, &notifications); err != nil {
        http.Error(w, "Error decoding notifications", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(notifications)
}

func (a *App) markNotificationRead(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    notifID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("notifications")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(notifID)
    if err != nil {
        http.Error(w, "Invalid notification ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID, "user_id": username},
        bson.M{"$set": bson.M{"is_read": true}},
    )
    if err != nil {
        http.Error(w, "Error marking notification as read", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Notification not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Notification marked as read"))
}

func (a *App) sendMessage(w http.ResponseWriter, r *http.Request) {
    var message models.Message
    if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
        http.Error(w, "Invalid message data", http.StatusBadRequest)
        return
    }

    username := r.Context().Value("username").(string)
    message.From = username
    message.CreatedAt = time.Now().Format(time.RFC3339)

    // Check privacy settings
    usersCollection := a.DB.Database("amity").Collection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var recipient models.User
    err := usersCollection.FindOne(ctx, bson.M{"username": message.To}).Decode(&recipient)
    if err != nil {
        http.Error(w, "Recipient not found", http.StatusNotFound)
        return
    }

    if recipient.Settings.Privacy.Messaging == "friends" && !contains(recipient.Friends, username) {
        http.Error(w, "Recipient only accepts messages from friends", http.StatusForbidden)
        return
    }
    if recipient.Settings.Privacy.Messaging == "none" {
        http.Error(w, "Recipient has disabled messaging", http.StatusForbidden)
        return
    }

    // If is_ai is true, simulate an AI response
    if message.IsAI {
        message.Content = fmt.Sprintf("AI Response: I understand your message: '%s'. How can I assist you further?", message.Content)
        message.From = "AI Assistant"
    }

    collection := a.DB.Database("amity").Collection("messages")
    _, err = collection.InsertOne(ctx, message)
    if err != nil {
        http.Error(w, "Error sending message", http.StatusInternalServerError)
        return
    }

    a.createNotification(message.To, "message", username, fmt.Sprintf("%s sent you a message", username), "")

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(message)
}

func (a *App) getMessages(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    otherUser := vars["username"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("messages")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{
        "$or": []bson.M{
            {"from": username, "to": otherUser},
            {"from": otherUser, "to": username},
        },
    })
    if err != nil {
        http.Error(w, "Error fetching messages", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var messages []models.Message
    if err := cursor.All(ctx, &messages); err != nil {
        http.Error(w, "Error decoding messages", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(messages)
}

func (a *App) createList(w http.ResponseWriter, r *http.Request) {
    var list models.List
    if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
        http.Error(w, "Invalid list data", http.StatusBadRequest)
        return
    }

    username := r.Context().Value("username").(string)
    list.UserID = username
    list.CreatedAt = time.Now().Format(time.RFC3339)

    collection := a.DB.Database("amity").Collection("lists")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, list)
    if err != nil {
        http.Error(w, "Error creating list", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(list)
}

func (a *App) getLists(w http.ResponseWriter, r *http.Request) {
    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("lists")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"user_id": username})
    if err != nil {
        http.Error(w, "Error fetching lists", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var lists []models.List
    if err := cursor.All(ctx, &lists); err != nil {
        http.Error(w, "Error decoding lists", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(lists)
}

func (a *App) addToList(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    listID := vars["id"]

    var item models.ListItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        http.Error(w, "Invalid item data", http.StatusBadRequest)
        return
    }

    item.AddedAt = time.Now().Format(time.RFC3339)

    collection := a.DB.Database("amity").Collection("lists")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(listID)
    if err != nil {
        http.Error(w, "Invalid list ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$push": bson.M{"items": item}},
    )
    if err != nil {
        http.Error(w, "Error adding to list", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "List not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Added to list"))
}

func (a *App) createHangout(w http.ResponseWriter, r *http.Request) {
    var hangout models.Hangout
    if err := json.NewDecoder(r.Body).Decode(&hangout); err != nil {
        http.Error(w, "Invalid hangout data", http.StatusBadRequest)
        return
    }

    username := r.Context().Value("username").(string)
    hangout.Creator = username
    hangout.Participants = []string{username}
    hangout.CreatedAt = time.Now().Format(time.RFC3339)

    collection := a.DB.Database("amity").Collection("hangouts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, hangout)
    if err != nil {
        http.Error(w, "Error creating hangout", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(hangout)
}

func (a *App) getHangouts(w http.ResponseWriter, r *http.Request) {
    collection := a.DB.Database("amity").Collection("hangouts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        http.Error(w, "Error fetching hangouts", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var hangouts []models.Hangout
    if err := cursor.All(ctx, &hangouts); err != nil {
        http.Error(w, "Error decoding hangouts", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(hangouts)
}

func (a *App) joinHangout(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    hangoutID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("hangouts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(hangoutID)
    if err != nil {
        http.Error(w, "Invalid hangout ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$addToSet": bson.M{"participants": username}},
    )
    if err != nil {
        http.Error(w, "Error joining hangout", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Hangout not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Joined hangout"))
}

func (a *App) leaveHangout(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    hangoutID := vars["id"]

    username := r.Context().Value("username").(string)

    collection := a.DB.Database("amity").Collection("hangouts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(hangoutID)
    if err != nil {
        http.Error(w, "Invalid hangout ID", http.StatusBadRequest)
        return
    }

    result, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$pull": bson.M{"participants": username}},
    )
    if err != nil {
        http.Error(w, "Error leaving hangout", http.StatusInternalServerError)
        return
    }
    if result.MatchedCount == 0 {
        http.Error(w, "Hangout not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Left hangout"))
}

func (a *App) handleWebFinger(w http.ResponseWriter, r *http.Request) {
    resource := r.URL.Query().Get("resource")
    if resource == "" {
        http.Error(w, "Missing resource parameter", http.StatusBadRequest)
        return
    }

    response := map[string]interface{}{
        "subject": resource,
        "links": []map[string]string{
            {
                "rel":  "self",
                "type": "application/activity+json",
                "href": fmt.Sprintf("https://amity.example.com/users/%s", resource[5:]),
            },
        },
    }

    w.Header().Set("Content-Type", "application/jrd+json")
    json.NewEncoder(w).Encode(response)
}

func (a *App) handleOutbox(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]

    collection := a.DB.Database("amity").Collection("posts")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"username": username})
    if err != nil {
        http.Error(w, "Error fetching outbox", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(ctx)

    var posts []models.Post
    if err := cursor.All(ctx, &posts); err != nil {
        http.Error(w, "Error decoding posts", http.StatusInternalServerError)
        return
    }

    activities := map[string]interface{}{
        "@context": "https://www.w3.org/ns/activitystreams",
        "type":     "OrderedCollection",
        "items":    posts,
    }

    w.Header().Set("Content-Type", "application/activity+json")
    json.NewEncoder(w).Encode(activities)
}

func (a *App) handleInbox(w http.ResponseWriter, r *http.Request) {
    var activity map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
        http.Error(w, "Invalid activity data", http.StatusBadRequest)
        return
    }

    log.Printf("Received activity: %v", activity)
    w.WriteHeader(http.StatusAccepted)
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func (a *App) Run() {
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func main() {
    app := &App{}
    app.Initialize()
    app.Run()
}