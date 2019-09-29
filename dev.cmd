REM Start dev environment in Windows
start "docker" cmd /k "docker-compose --file=docker-compose.dev.yml up & docker-compose --file=docker-compose.dev.yml down"
start "webapp" /D "app" cmd /k "npm start"
start "service/api-gateway" /D "services/api-gateway/cmd" cmd /k "go build && cmd"
start "service/authentication" /D "services/authentication/cmd" cmd /k "go build && cmd"
start "service/post" /D "services/post/cmd" cmd /k "go build && cmd"
start "service/comment" /D "services/comment/cmd" cmd /k "go build && cmd"