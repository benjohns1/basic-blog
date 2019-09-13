# Basic Blog - Prototype Version 2
This basic blog implementation makes use of Material Angular on the front-end with a Go API back-end and Postgres for persistence.
## Run It
1. Install & run Docker
2. `git clone https://github.com/benjohns1/basic-blog`
3. `cd basic-blog`
5. Start the app with `docker-compose up`
6. Go to `localhost:8080`
7. Log in with dummy credentials
   - username: `bobross`
   - password: `painter`

## Overview
Building off of the feature-complete prototype 1, this version has a UI completely rewritten in Angular and has a slightly refactored Go API.
### Front-End
#### `./app`
The Angular app uses Angular Material for styling, components and services for separation of concerns.
 - `./src/app/app.component.ts`: Main app component and page layout
 - `./src/app/app-routing.module.ts`: App routing
 - `./src/app/blog.service.ts`: Wrapper service for blog REST API
 - `./src/app/authentication.service.ts`: Wrapper service for dummy authentication API and storing the authenication token client-side
 - `./src/app/post.ts`: Blog post model classes
 - `./src/app/comment.ts`: Blog comment model classes
 - `./src/app/comments`: Component handles the comment form and comment list
 - `./src/app/login-form`: Component handles the authentication form
 - `./src/app/page-not-found`: Page error component handler
 - `./src/app/post`: Component handles displaying and editing a single post
 - `./src/app/post-list`: Component displays lists of posts or deleted posts

### Back-End
#### `./services/api`
On startup, the simple Go service in `cmd/main.go` connects and sets up the DB persistence layer in `./internal/postgres/postgres.go` along with some dummy data, then injects it into the API server `./internal/api/api.go`'s Run() function.

## Infrastructure
Docker and Compose
 - Web App: node alpine image as builder -> targets a nginx image that hosts the Angular web application
 - API Server: golang image as builder -> targets a scratch image that hosts the Go API binary
 - DB: postgres image (with adminer image for DB dev/inspection)

## Development
Prereqs: Ensure you have the lastest Angular CLI and Docker installed.
1. Start the DB and Adminer: `docker-compose --file=docker-compose.dev.yml up`
2. In `services/api/cmd`, rebuild and start the API server after any changes: `go build -o blog && ./blog` (or `go build -o blog.exe && blog` on Windows)
3. In `app`, start the client dev server and make changes: `ng serve`
4. API URL: `localhost:3000`
5. App URL: `localhost:8080`
6. Adminer URL: `localhost:8081`

# Devlog
## 2019-09-12 (hours 8-10):
Rebuilt front-end in Material Angular for a much cleaner UX
What I was hoping to do but ran out of time:
 - Friendly user error and success messages (currenty only logged to console)
 - UI cleanup (show active route, favicon, custom color scheme with some more icons for flavor)
 - Real unit tests and some e2e tests
 - Refactor the Go API to use separate microservices for authentication, blog posts, and blog comments
 - Use protobufs for transport
 - Use event-driven architecture and event sourcing for back-end
 - Implement authentication for real users

## 2019-09-12 (hour 8):
Feature complete with raw UI, finished edit post functionality and updated readme for the initial prototype!  
`git checkout prototype-1` to view this working version (WARNING: it has a face only a mother could love)

## 2019-09-11 (hours 5-7):
Almost feature complete except for the edit post functionality in the UI

## 2019-09-10 (hours 3-5): 
Added dummy authentication and authenticated vs anonymous logic in the UI

## 2019-09-09: Design and Initial Features (hour 1-3)
Fleshed-out unauthenticated flows, UI, setup DB, and docker
### REST API Endpoints
 - /api/v1/authenticate (POST)
 - /api/v1/post (GET, POST)
 - /api/v1/post/{id} (GET, POST, DELETE)
 - /api/v1/post/{id}/comment (POST)
### DB Design
This is a target design for the final microservice bounded contexts, initially I'll implement the API within a single context, hard-code a dummy user, and only use the `post` and `comment` tables
#### Authentication Context
user
 - id
 - username
 - passwordhash
anon_user
 - id
 - ip
 - useragent
#### Blog Post Context
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
#### Comment Context
comment
 - id
 - body
 - post_id
 - commenter_id
Commenter
 - id
 - user_id (nullable)
 - displayname

## 2019-09-08: Initial Requirements and Approach (hours 0-1)
### Project Approach
I think it would be fun to try a simple project like this in stages since I believe I can get all the base requirements done in the time allotted with hopefully some to spare. I've been a practicing consultant for the past 8 years, so my estimates should only be off by a factor of 10,000% or thereabouts. (I can be an optimist and a pragmatist at the same time, right?) So, the goal is to get all requirements done with the bare minimum of effort and engineering, then refactor the architecture, and add a nice UI if there's time.
### Design Stages
#### 1: Prototype
Bare-bones fulfillment of requirements
 - Plain JavaScript front-end SPA
 - Simple Go REST API
#### 2: UI and Refactor
Style the UI and refactor into a cleaner architecture
 - Angular
 - Angular Material
 - Refactor API
#### 3: New Tech
Unlikely to get here in the time limit, but it's good to set high goals!
Refactor to use the following architectural patterns
 - Split API handlers into smaller bounded contexts:
   1. Authentication Service
   2. Blog Post Service
   3. Comment Service
 - gRPC for communication protocol
 - Event sourcing
 - CQRS
### Constraints
 - 10 hours max
 - Separate back-end API and front-end SPA
 - Relational DB
### Story Reference
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