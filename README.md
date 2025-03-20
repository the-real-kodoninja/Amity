# Amity - A Decentralized Social Networking Platform on Mastodon and Motoko

![Amity 2015](https://github.com/the-real-kodoninja/Amity/blob/main/references/11148567_955037097860610_1543569738534476929_o.jpg?raw=true) 

Amity is a modern reimagining of the GreenHeartPT social networking platform, originally created in 2015 by [the-real-kodoninja](https://github.com/the-real-kodoninja). The original platform was built using PHP and jQuery, but Amity brings the same design and functionality into the present with a fresh tech stack: **Go**, **Gorilla Mux**, **Material-UI**, **React**, and **MongoDB**, now integrated with the **Mastodon network** and deployed on the **Motoko chain** (Internet Computer blockchain). Amity aims to provide a fun, simple, and secure social networking experience where users can connect, share, and engage with others while maintaining the nostalgic design of the original GreenHeartPT, all within a decentralized ecosystem.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Original GreenHeartPT (2015)](#original-greenheartpt-2015)
- [Tech Stack](#tech-stack)
- [Mastodon Network Integration](#mastodon-network-integration)
- [Motoko Chain Deployment](#motoko-chain-deployment)
- [Does This Commit MongoDB?](#does-this-commit-mongodb)
- [Getting Started](#getting-started)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

Amity is a decentralized social networking platform designed to foster connections, creativity, and community. It allows users to create profiles, share posts and photos, connect with others, and engage in a variety of social interactions. The platform supports anonymity, media sharing, and a clean, grid-based layout that emphasizes user content. Amity retains the exact design of the original GreenHeartPT platform, which was in beta (version 1.0019) in 2015, but updates the underlying technology for better performance, scalability, and decentralization.

By integrating with the **Mastodon network**, Amity leverages the Fediverse to enable cross-platform interactions, allowing users to follow and interact with Mastodon users across different instances. The platform is deployed on the **Motoko chain** (part of the Internet Computer blockchain), ensuring decentralized data storage and user sovereignty. Amity uses **MongoDB** as its primary database for efficient data management, with considerations for how this interacts with the blockchain.

The platform is built with a focus on:
- **User Engagement:** Connect with friends, follow others, and share posts, photos, and media.
- **Anonymity:** Post anonymously to encourage open expression.
- **Simplicity:** A clean, intuitive interface that’s easy to navigate.
- **Decentralization:** Leverage the Mastodon network and Motoko chain for a decentralized social experience.
- **Community:** Discover trends, find connections, and engage with like-minded users across the Fediverse.

---

## Features

Amity includes a rich set of features, all of which are carried over from the original GreenHeartPT platform, with modern implementations and enhancements for decentralization. Below is a detailed breakdown of the platform’s functionality:

### 1. **User Profiles**
- **Profile Header:** Each user has a profile with a cover photo, profile picture, username, location (e.g., "United States of America"), and last active status (e.g., "2 weeks ago").
- **Stats:** Displays counts for "arena," "photos," "connections," and "followers" (e.g., "photos 11," "connections 6," "followers 1").
- **Connections Grid:** A grid of profile pictures for the user’s connections, linking to their profiles.
- **Following Section:** Shows the number of users the current user is following (e.g., "following 0").
- **Photos Section:** A grid of user-uploaded photos with timestamps (e.g., "1 month ago") and interaction options (like, share).
- **Actions:** Users can connect with or follow others, with buttons to toggle these states (e.g., "connected," "following").

### 2. **Feed**
- **Global Feed:** Displays posts from users the current user is connected to or following, including Mastodon users via ActivityPub.
- **Posts:** Each post includes the username, timestamp, content (text, images), and hashtags (e.g., "I so need this #bike").
- **Interactions:** Users can like, share, and comment on posts, with counts displayed (e.g., "4 likes," "1 share").
- **Post Creation:** A "post" button allows users to create new posts with text, images, and hashtags, which can be federated to Mastodon.
- **Post Actions:** Options to hide, report, or delete posts via a dropdown menu.

### 3. **Photos Page**
- **Photo Grid:** Displays a grid of photos uploaded by a user, with timestamps and interaction options (like, share).
- **Media Support:** Supports images of various themes (e.g., comic characters, anime, personal photos).

### 4. **Search**
- **User Search:** A search bar at the top allows users to search for other users by username (e.g., searching for "b" shows "batfan," "batman," "bunny"), including Mastodon users.
- **Results:** Displays a list of matching users with their usernames, locations, and profile pictures.

### 5. **Notifications**
- **Notification Dropdown:** Shows recent activity, such as "larry posted on his page" or "birdman posted on his page," including Mastodon interactions.
- **Connection Requests:** Displays a "No connection requests" message when there are none.

### 6. **Anonymity**
- **Anonymous Posting:** Users can post anonymously, with their identity hidden (e.g., labeled as "anonymous").
- **Toggle Option:** A setting to enable/disable anonymity for posts.

### 7. **Login/Signup**
- **Login Form:** A simple form with fields for email and password, a "Login" button, and links for "Create Account" and "Forgot Something?"
- **Sidebar Features:** Highlights platform features like "Stay connected with everyone," "Promote your post," "Connect with what you like," "Become anonymous," and "Create topics & trends."
- **Stats:** Shows "Total Users" and "Total Posts" (e.g., "Total Users: 13," "Total Post: 115") with thumbnails of recent posts.
- **Decentralized Identity:** Users can log in using their Internet Computer identity (via Motoko chain).

### 8. **Developer Page**
- **Static Content:** A page with sections like "Goals," "How to Join," "Software," "Browser," "Programming," "Graphics," and "Updates."
- **Development Tools:** Mentions the use of Firefox Developer Edition, Notepad++, Sublime Text, Gedit, FileZilla, and GIMP.
- **Feedback Encouragement:** Encourages users to provide feedback on how they can help improve the platform.

### 9. **Trends & Connections**
- **Trends Section:** Displays popular hashtags (e.g., "#BETA released") and trending content, including Mastodon trends.
- **Find Connections:** A list of suggested users to connect with, showing their usernames, locations, and profile pictures, including Mastodon users.

### 10. **Beta Testing and Feedback**
- **Beta Status:** The platform is in beta (version 1.0019 in the original), with a "Want to help emoorephp?" link for feedback.
- **Feedback Form:** A form or modal for users to submit feedback during the beta phase.

### 11. **General Layout**
- **Header:** Includes the logo, search bar, user profile icon (with status indicator), and navigation tabs (arena, feed, photos, connections, followers).
- **Sidebar (Left):** Shows user info, connections grid, following section, help prompt, and footer links (about, policy, terms, ads, promote, developer, groups, language, social).
- **Main Content Area:** Displays the feed, photos, or other user activity in a scrollable area with cards for each item.
- **Right Column:** Contains tabs for "trends," "news," and "browse," with a "Popular & Trending" section and a "Find Connections" list.

---

## Original GreenHeartPT (2015)

The original GreenHeartPT platform was created in 2015 by [the-real-kodoninja](https://github.com/the-real-kodoninja) using PHP and jQuery. It was a beta project (version 1.0019) under the "Aviyon" brand, focusing on social networking with a clean, grid-based design. Below are screenshots from the original platform, showcasing its design and features. These images are stored in the `references` folder of this repository.

### Screenshots of GreenHeartPT

1. **Login Page**
   - **Description:** The login page features a simple form for email and password, with a sidebar highlighting features like "Stay connected with everyone" and "Become anonymous." It also shows stats like "Total Users: 13" and "Total Post: 115."
     
      ![references/login.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11169579_955037271193926_4355236387694178411_o.jpg?raw=true)

2. **User Profile (batgirl)**
   - **Description:** The profile page for user "batgirl" shows a cover photo with *The Powerpuff Girls* characters, a profile picture of Batgirl, and stats (e.g., "photos 11," "connections 6"). The photos section displays comic-related images like Harley Quinn.
     
     ![references/batgirl.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11079998_955037524527234_2253325221238154678_o.jpg?raw=true)
     ![references/batgirl.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11154576_955037581193895_8025093307947206791_o.jpg?raw=true)
     ![references/batgirl.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11164635_955037657860554_7614564749226175066_o.jpg?raw=true)

3. **User Profile (emoorephp)**
   - **Description:** The profile page for user "emoorephp" features a scenic cover photo, a profile picture, and stats (e.g., "feed 45," "connections 9"). The user is a developer at "Aviyon" and a "fitness enthusiast."
     
     ![references/emoorephp.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11148567_955037097860610_1543569738534476929_o.jpg?raw=true)
     ![references/emoorephp.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11161703_955037167860603_3753640108829535325_o.jpg?raw=true)
     ![references/emoorephp.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11169554_955038464527140_1781779550120457726_o.jpg?raw=true)
     ![references/emoorephp.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11174302_955037351193918_9146954543274840299_o.jpg?raw=true)

4. **Feed Page**
     **Description:** The global feed shows posts from connected users, such as "larry" posting "#BETA" and "emoorephp" posting "I so need this #bike" with an image of a futuristic bike.
     
     ![references/feed.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/1980501_955038101193843_5565818349684851117_o.jpg?raw=true)

5. **Developer Page**
   - **Description:** The developer page provides info about the platform’s goals, software (e.g., Firefox Developer Edition, Notepad++), and encourages feedback.
     
     ![references/developer.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11174728_955038264527160_6465011376463576850_o.jpg?raw=true)

6. **Search Results**
   - **Description:** Searching for "b" shows users like "batfan," "batman," and "bunny," with their usernames, locations, and profile pictures.
     
     ![references/search.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11157523_955038047860515_7463366602084980064_o.jpg?raw=true)

7. **Photos Page (batgirl)**
   - **Description:** The photos page for "batgirl" displays a grid of comic-related images, such as Harley Quinn and Batgirl, with timestamps.
     
     ![references/photos.png](https://github.com/the-real-kodoninja/Amity/blob/main/references/11011939_955037934527193_2760126084887791562_o.jpg?raw=true)

### Back to Amity
[⬅ Back to Amity README](#amity---a-decentralized-social-networking-platform-on-mastodon-and-motoko)

---

## Tech Stack

Amity is built with a modern stack to ensure performance, scalability, and decentralization while preserving the original design of GreenHeartPT. The technologies used are:

- **Backend:**
  - **Go:** A fast, statically typed language for building the server-side application.
  - **Gorilla Mux:** A powerful URL router and dispatcher for Go, used for handling RESTful API routes.
  - **MongoDB:** A NoSQL database for storing user profiles, posts, photos, notifications, and feedback.

- **Frontend:**
  - **React:** A JavaScript library for building the user interface.
  - **Material-UI:** A React component library for implementing the design of GreenHeartPT with modern, responsive components.

- **Decentralized Infrastructure:**
  - **Mastodon Network (Fediverse):** Integrates with Mastodon via the ActivityPub protocol for federated social networking.
  - **Motoko Chain (Internet Computer):** Deploys the application on the Internet Computer blockchain using Motoko for decentralized data storage and identity management.

- **Development Tools:**
  - **Firefox Developer Edition:** For testing and debugging the web application.
  - **VS Code:** Recommended code editor for development.
  - **Docker:** For containerizing the application (optional).
  - **Git:** For version control, hosted on GitHub.
  - **DFX (Internet Computer SDK):** For deploying and managing canisters on the Motoko chain.

---

## Mastodon Network Integration

Amity integrates with the **Mastodon network**, a decentralized social network part of the Fediverse, using the **ActivityPub** protocol. This allows Amity users to:
- Follow and interact with users on other Mastodon instances (e.g., mozilla.social).
- Share posts that federate to the Mastodon network, making them visible to users on other instances.
- Receive notifications and messages from Mastodon users.
- Search for and connect with Mastodon users across the Fediverse.

### Implementation Details:
- **ActivityPub Protocol:** Amity implements the ActivityPub server-to-server (S2S) protocol to federate with Mastodon instances. This involves handling activities like `Create`, `Follow`, `Like`, and `Announce` (reblog).
- **User Identity:** Each Amity user gets a federated handle (e.g., `@username@amity.social`), which is compatible with Mastodon.
- **API Endpoints:** Extend the Go backend to handle ActivityPub requests (e.g., `/users/{username}/outbox`, `/users/{username}/inbox`) using Gorilla Mux.
- **MongoDB Storage:** Store ActivityPub activities (e.g., posts, follows) in MongoDB, ensuring they are synced with the Mastodon network.

This integration makes Amity a part of the broader Fediverse, allowing for cross-platform social interactions while maintaining its unique design and features.

---

## Motoko Chain Deployment

Amity is deployed on the **Motoko chain**, which is part of the **Internet Computer (IC)** blockchain developed by DFINITY. Motoko is a programming language designed for the IC, and deploying Amity on this chain ensures decentralized data storage, user identity management, and application hosting.

### Implementation Details:
- **Canisters:** The Amity backend is deployed as a set of canisters (smart contracts) on the Internet Computer. Each canister handles specific functionality (e.g., user profiles, posts, ActivityPub integration).
- **Motoko Backend:** Rewrite parts of the Go backend in Motoko to handle decentralized data storage and identity management. For example, user authentication can use the Internet Computer’s identity system (Internet Identity).
- **MongoDB Integration:** MongoDB is used as an off-chain database for performance, but critical data (e.g., user identities, post metadata) is mirrored on the IC blockchain for decentralization.
- **Frontend Hosting:** The React frontend (with Material-UI) is hosted on the IC as a canister, ensuring the entire application is decentralized.
- **DFX SDK:** Use the DFX command-line tool to deploy and manage canisters on the IC.

### Benefits:
- **Decentralization:** Users own their data, and the platform is resistant to censorship.
- **Scalability:** The IC provides scalable compute and storage for social networking workloads.
- **Security:** Blockchain-based identity ensures secure user authentication.

---

## Does This Commit MongoDB?

The phrase "does this commit MongoDB" in the context of integrating with the Mastodon network and deploying on the Motoko chain likely refers to whether MongoDB is "committed" to the blockchain (i.e., whether MongoDB data is stored on-chain or if MongoDB itself is replaced by the blockchain). Let’s break this down:

### **MongoDB’s Role in Amity**
- MongoDB is used as the primary database for Amity to store user profiles, posts, photos, notifications, and ActivityPub activities. It provides fast, scalable, and flexible data storage, which is ideal for a social networking platform.
- MongoDB is an off-chain database, meaning it runs on a traditional server (e.g., MongoDB Atlas or a local instance) rather than on the Motoko chain.

### **Motoko Chain Integration**
- The Motoko chain (Internet Computer) uses canisters to store data on-chain. These canisters can store user identities, post metadata, and other critical data in a decentralized manner.
- However, storing all data (e.g., large media files, full post content) on the blockchain is impractical due to cost and performance constraints. The Internet Computer charges cycles (a form of gas) for storage and computation, and large datasets like images or videos would be expensive to store on-chain.

### **Does MongoDB Get Committed to the Blockchain?**
- **No, MongoDB itself is not committed to the blockchain.** MongoDB remains an off-chain database, running on a traditional server or cloud service. It handles the bulk of the data storage for Amity, such as user profiles, posts, and media files.
- **Some data is mirrored on-chain:** To ensure decentralization, critical data (e.g., user identities, post metadata, ActivityPub activities) is stored in canisters on the Motoko chain. For example:
  - User identities are managed using Internet Identity, which is stored on the IC.
  - Post metadata (e.g., post ID, author, timestamp, hashtags) is stored on-chain to ensure immutability and federation with Mastodon.
  - Media files (e.g., images, videos) are stored in MongoDB (or a service like MongoDB GridFS) and referenced on-chain via URLs or hashes.
- **Hybrid Approach:** Amity uses a hybrid data storage model:
  - **On-chain (Motoko Chain):** User identities, post metadata, and ActivityPub activities for decentralization and federation.
  - **Off-chain (MongoDB):** Full post content, media files, and non-critical data for performance and cost efficiency.

### **Mastodon Network Implications**
- The Mastodon network integration via ActivityPub does not directly affect MongoDB’s role. ActivityPub activities (e.g., posts, follows) are stored in MongoDB for quick access and processing, but their metadata is also mirrored on the Motoko chain to ensure federation and decentralization.
- For example, when a user posts on Amity, the post is stored in MongoDB, but its metadata (e.g., post ID, author, timestamp) is recorded on the Motoko chain and federated to Mastodon via ActivityPub.

### **Conclusion**
- MongoDB is not "committed" to the blockchain in the sense of being replaced or fully stored on-chain. Instead, Amity uses a hybrid approach where MongoDB handles off-chain storage for performance, while the Motoko chain handles critical data for decentralization. This ensures Amity can scale efficiently while maintaining the benefits of a decentralized architecture.

---

## Getting Started

To get started with Amity, follow these steps:

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/the-real-kodoninja/Amity.git
   cd Amity
Set Up the Backend (Go):
Ensure Go is installed (go version to check).
Install dependencies:
bash

go get github.com/gorilla/mux
go get go.mongodb.org/mongo-driver/mongo
Set up MongoDB (locally or via a cloud service like MongoDB Atlas).
Update the MongoDB connection string in the backend configuration.
Set Up the Frontend (React):
Navigate to the frontend directory:
bash

cd frontend
Install dependencies:
bash


npm install
Start the React development server:
bash

npm start
Set Up the Motoko Chain (Internet Computer):
Install the DFX SDK for the Internet Computer:
bash

sh -ci "$(curl -fsSL https://internetcomputer.org/install.sh)"
Navigate to the canisters directory:
bash


cd canisters
Deploy the canisters:
bash

dfx deploy
Run the Backend (Go):
Navigate to the backend directory:
bash

cd backend
Start the Go server:
bash

go run main.go
Access the Application:
Open your browser and navigate to http://localhost:3000 for the frontend (or the IC canister URL after deployment).
The backend API will be available at http://localhost:8080 (or the port you configure).
Contributing
Contributions are welcome! If you’d like to contribute to Amity, please follow these steps:

Fork the repository.
Create a new branch (git checkout -b feature/your-feature).
Make your changes and commit them (git commit -m "Add your feature").
Push to your branch (git push origin feature/your-feature).
Open a pull request.
Please ensure your code follows the project’s style guidelines and includes tests where applicable.

License
This project is licensed under the MIT License. See the LICENSE file for details.
---

### **Feasibility and Implications**

Let’s discuss the feasibility of integrating Amity with the Mastodon network and deploying it on the Motoko chain, as well as the implications for MongoDB.

#### **Mastodon Network Integration**
- **ActivityPub Protocol:** Mastodon uses the ActivityPub protocol for federation, which Amity can implement in Go. Libraries like `go-fed/activity` can help handle ActivityPub activities (e.g., `Create`, `Follow`, `Like`). The Go backend (using Gorilla Mux) can expose endpoints like `/users/{username}/outbox` and `/users/{username}/inbox` to handle federation.
- **Challenges:** Implementing ActivityPub is complex, as it requires handling signatures (HTTP Signatures), WebFinger for user discovery, and ensuring compatibility with Mastodon’s implementation. However, this is feasible and aligns with the Fediverse’s goal of decentralized social networking.
- **MongoDB Role:** MongoDB will store ActivityPub activities for quick access, but metadata (e.g., post IDs, user handles) will also be mirrored on the Motoko chain to ensure decentralization.

#### **Motoko Chain Deployment**
- **Internet Computer (IC):** The Motoko chain refers to the Internet Computer blockchain, where Motoko is the primary programming language for writing canisters (smart contracts). Deploying Amity on the IC is feasible, as the IC supports hosting full-stack applications (backend and frontend).
- **Go and Motoko:** The Go backend can be partially rewritten in Motoko for on-chain logic (e.g., user identity, post metadata). However, Go can still handle off-chain tasks (e.g., ActivityPub federation, MongoDB interactions) and communicate with Motoko canisters via HTTP or Candid (IC’s interface description language).
- **Frontend:** The React frontend (with Material-UI) can be deployed as a canister on the IC, ensuring the entire application is decentralized.
- **Challenges:** The IC has limitations on storage and computation (e.g., cycles cost for storage). Storing large media files on-chain is impractical, so MongoDB remains essential for off-chain storage.

#### **Does This Commit MongoDB?**
As discussed in the README, MongoDB is not fully committed to the blockchain. Instead, Amity uses a hybrid approach:
- **Off-chain (MongoDB):** Stores full post content, media files, and non-critical data for performance and cost efficiency.
- **On-chain (Motoko Chain):** Stores user identities, post metadata, and ActivityPub activities for decentralization and federation.
This hybrid model ensures Amity can scale while maintaining the benefits of decentralization.

#### **Does This Conflict with MongoDB?**
- **No Conflict:** MongoDB works well in this setup as an off-chain database. It complements the Motoko chain by handling large-scale data storage (e.g., media files) that would be expensive to store on-chain.
- **Data Synchronization:** The Go backend will need to sync data between MongoDB and the Motoko chain. For example, when a user posts, the post is stored in MongoDB, and its metadata is recorded on the IC. This requires careful design to ensure consistency (e.g., using eventual consistency or a message queue).

---

### **Next Steps**
1. **Project Structure:** Set up the repository with directories for the Go backend (`backend`), React frontend (`frontend`), and Motoko canisters (`canisters`).
2. **Go Backend:** Implement the REST API with Gorilla Mux, including ActivityPub endpoints for Mastodon integration.
3. **MongoDB Setup:** Configure MongoDB to store user data, posts, and ActivityPub activities.
4. **Motoko Canisters:** Write Motoko code for on-chain logic (e.g., user identity, post metadata) and deploy using DFX.
5. **Frontend:** Build the React frontend with Material-UI, ensuring it matches the original GreenHeartPT design, and deploy it as a canister on the IC.
6. **Testing:** Test the integration with Mastodon (e.g., federate posts to a Mastodon instance like mozilla.social) and ensure data sync between MongoDB and the Motoko chain.
