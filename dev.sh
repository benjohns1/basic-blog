#!/bin/bash

# Start dev environment in Mac OSX
path=$PWD
osascript -e "tell app \"Terminal\"
    do script \"cd ${path} && docker-compose --file=docker-compose.dev.yml up & docker-compose --file=docker-compose.dev.yml down\"
    do script \"cd ${path}/app && npm start\"
    do script \"cd ${path}/services/api-gateway/cmd && go build && ./cmd\"
    do script \"cd ${path}/services/authentication/cmd && go build && ./cmd\"
    do script \"cd ${path}/services/post/cmd && go build && ./cmd\"
    do script \"cd ${path}/services/comment/cmd && go build && ./cmd\"
end tell"
