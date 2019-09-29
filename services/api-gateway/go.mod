module github.com/benjohns1/basic-blog-ms/services/api-gateway

go 1.13

replace (
	github.com/benjohns1/basic-blog-ms/services/authentication => ../authentication
	github.com/benjohns1/basic-blog-ms/services/comment => ../comment
	github.com/benjohns1/basic-blog-ms/services/post => ../post
)

require (
	github.com/benjohns1/basic-blog-ms/services/authentication v0.0.0-00010101000000-000000000000
	github.com/benjohns1/basic-blog-ms/services/comment v0.0.0-00010101000000-000000000000
	github.com/benjohns1/basic-blog-ms/services/post v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.23.1
)
