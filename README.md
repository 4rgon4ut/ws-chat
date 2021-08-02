# Installation
## Download
      
      git clone git@github.com:bestpilotingalaxy/ws-chat.git
      
      cd ws-chat

## Starting server
*change `.env` if necessary

        docker build  -t ws-chat-server .
        
        
        docker run  --network host  --env-file=.env  ws-chat-server:latest

## Starting client
        
        docker run  --network host --env-file=.env  -it --entrypoint  /bin/bash  solsson/websocat
        
From container shell:

        websocat ws://0.0.0.0:$SERVER_PORT/ws

