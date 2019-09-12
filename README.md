# Basic Blog
## Prototype Version 1
A bare-bones blog implementation with a Go REST API and no-nonsense raw JavaScript SPA front-end, with dummy authentication
## Run It
1. Install & run Docker
2. Clone this repo and `cd` into the root directory
3. Run `docker-compose up`
4. Go to `localhost:8080`
5. Log in with dummy credentials
   - username: `bobross`
   - password: `painter`

## Overview
### Front-End
After building various UIs recently in React, Angular, and Svelte/Sapper, I decided to tackle this simple UI in good ol' plain JavaScript. 'Gasp', right? At least it's ES6 (it won't run in older browsers) and it was a good reminder for me why JavaScript frameworks exist. Adding more features towards the end required more architectural changes and clunkier code that most frameworks handle very elegantly. And XSS? What's that? But it was very nice not to have to fight with WebPack, TypeScript, SSR, and the million other tools we now use to make life on the front-end 'easier'.
#### `basic-blog/index.html`
A single page with a single script tag that handles API fetch requests and DOM manipulation. Who needs MVC or component architecture when you can shove it all in one file?
### Back-End
For simplicity, I avoided using DDD or layered/clean architecture on the back-end, and built the blog functionality into a single main.go file. I tried hard, but couldn't avoid using dependency injection for the DB connection reference throughout the API handlers.
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
2. Rebuild and start the server after changes: `go build && ./basic-blog`