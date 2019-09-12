# Basic Blog - Prototype Version 1
# Devlog 2019-09-11: Prototype Version 1
A bare-bones blog implementation with a Go REST API and no-nonsense plain JavaScript SPA front-end, with dummy authentication
## Run It
1. Install & run Docker
2. `git clone https://github.com/benjohns1/basic-blog`
3. `cd basic-blog`
4. `git checkout prototype-1`
5. Start the app with `docker-compose up`
6. Go to `localhost:8080`
7. Log in with dummy credentials
   - username: `bobross`
   - password: `painter`

## Overview
### Front-End
After building various UIs recently in React, Angular, and Svelte/Sapper, I decided to tackle this simple UI in good ol' plain JavaScript. 'Gasp', right? At least it's ES6 (it won't run in older browsers) and it was a good reminder for me why JavaScript frameworks exist. Adding more features towards the end required more architectural changes and clunkier code that most frameworks handle very elegantly. And XSS? What's that? But it was very nice not to have to fight with WebPack, TypeScript, SSR, and the million other tools we now use to make life on the front-end 'easier'.
#### `basic-blog/index.html`
A single page with a single script tag that handles API fetch requests and DOM manipulation. Who needs MVC or component architecture when you can shove it all in one file?
### Back-End
For simplicity, I avoided using DDD or layered/clean architecture on the back-end, and built the blog functionality into a single main.go file. No matter how hard I tried, I couldn't avoid using dependency injection for the DB connection reference throughout the API handlers ;-)
#### `basic-blog/main.go`
 - API: A simple JSON REST API to handle CRUD for blog posts and comments, and a dummy authentication endpoint that returns an "auth token" for subsequent API requests (and I use the term "auth token" very loosely ;-P)
 - Persistence: PostgreSQL database with hard-coded credentials stored in plaintext the repo (that's best practice, right?!?) "I was on a short deadline with this project and just needed to get it done." Famous last words before a massive data breach? Probably.
 - App Server: Statically serving `index.html` in a theater near you (actually only on `localhost:8080`)

### Infrastructure
Docker and Compose
 - DB: Official postgres image (with adminer image for DB dev/inspection)
 - Server: golang:1.13.0 image as builder image targeting a scratch image that hosts the application binary

## Development
1. Start the DB and Adminer: `docker-compose --file=docker-compose.dev.yml up`
2. Rebuild and start the server after changes: `go build -o blog && ./blog` (or `go build -o blog.exe && blog` on Windows)

# Devlog 2019-09-12 (hour 8):
Feature complete, finished edit post functionality and updated readme for the initial prototype!

# Devlog 2019-09-11 (hours 5-7):
Almost feature complete except for the edit post functionality in the UI

# Devlog 2019-09-10 (hours 3-5): 
Added dummy authentication and authenticated vs anonymous logic in the UI

# Devlog 2019-09-09: Design and Initial Features (hour 1-3)
Fleshed-out unauthenticated flows, UI, setup DB, and docker
## DB Design
This is a target design for the final microservice bounded contexts, initially I'll implement the API within a single context, hard-code a dummy user, and only use the `post` and `comment` tables
### Authentication Context
user
 - id
 - username
 - passwordhash
anon_user
 - id
 - ip
 - useragent
### Blog Post Context
post
 - id
 - title
 - created_time
 - body
 - author_id
author
 - id
 - user_id
 - displayname
### Comment Context
comment
 - id
 - body
 - post_id
 - commenter_id
Commenter
 - id
 - user_id (nullable)
 - displayname

## REST API Endpoints
 - /api/v1/authenticate (POST)
 - /api/v1/post (GET, POST)
 - /api/v1/post/{id} (GET, POST, DELETE)
 - /api/v1/post/{id}/comment (POST)

# Devlog 2019-09-08: Initial Requirements and Approach (hours 0-1)
## Project Approach
I think it would be fun to try a simple project like this in stages since I believe I can get all the base requirements done in the time allotted with hopefully some to spare. I've been a practicing consultant for the past 8 years, so my estimates should only be off by a factor of 10,000% or thereabouts. (I can be an optimist and a pragmatist at the same time, right?) So, the goal is to get all requirements done with the bare minimum of effort and engineering, then refactor the architecture, and add a nice UI if there's time.
## Design Stages
### 1: Prototype
Bare-bones fulfillment of requirements
 - Plain JavaScript front-end SPA
 - Simple Go REST API
### 2: UI and Refactor
Style the UI and refactor into a cleaner architecture
 - Angular
 - Angular Material
 - Refactor API
### 3: New Tech
Unlikely to get here in the time limit, but it's good to set high goals!
Refactor to use the following architectural patterns
 - Split API handlers into smaller bounded contexts:
   1. Authentication Service
   2. Blog Post Service
   3. Comment Service
 - gRPC for communication protocol
 - Event sourcing
 - CQRS
## Constraints
 - 10 hours max
 - Separate back-end API and front-end SPA
 - Relational DB
## Story Reference
1. As an authenticated user, I can post a new blog post. The post has a title, created date,
and body.
a. NOTE: do not worry about Rich text or WSIWYG editors for the post body. Plain
text is fine.
2. As any user, I can publicly view a list of blog posts, sorted by created date in descending
order.
3. As any user, I can publicly view an individual blog post.
4. As any user, I can comment on a blog post and view all comments of the blog post.
5. As a user, I can authenticate by providing a username and password.
a. NOTE: "fake" authentication is sufficient.
b. Do not worry about user creation screens or endpoints for authentication.
c. Authenticate a user against a hardcoded username/password on the server.
6. As an authenticated user, I can delete, undelete, or edit an existing blog post.
7. As another developer on your team I can read your clear, thorough documentation and
run your project with a single command.