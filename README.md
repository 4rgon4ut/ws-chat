


# Installation
### Download
      
      git clone git@github.com:bestpilotingalaxy/ws-chat.git
      
      cd ws-chat

### Starting server
*change `.env` if necessary

        docker build  -t ws-chat-server .
        
        
        docker run  --network host  --env-file=.env  ws-chat-server:latest

### Starting client
        
        docker run  --network host --env-file=.env  -it --entrypoint  /bin/bash  solsson/websocat
        
From container shell:

        websocat ws://0.0.0.0:$SERVER_PORT/ws


# Structure 
![image](https://user-images.githubusercontent.com/59182467/127800515-bc5ed38d-ceda-40b4-8063-5caeb63b8eb8.png)



# Debug

I'm personally use VSCode with dlv-dap debugger
* https://github.com/golang/vscode-go/blob/master/docs/dlv-dap.md


### launch.json
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "debugAdapter": "dlv-dap",
            "trace": "verbose",
            "program": "${workspaceFolder}/cmd/server/",
            "envFile": "${workspaceFolder}/.env"
        }
    ]
}
```