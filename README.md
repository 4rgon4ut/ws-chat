


Fiber based websocket nonblocking echo-chat server.

# Download
      
      git clone git@github.com:bestpilotingalaxy/ws-chat.git
      
      cd ws-chat

# Starting server
*change `.env` if necessary

        docker build  -t ws-chat-server .
        
        
        docker run  --network host  --env-file=.env  ws-chat-server:latest

# Starting client 
        
        docker run  --network host --env-file=.env  -it --entrypoint  /bin/bash  solsson/websocat
        
### Text messages

From container shell:

        websocat ws://0.0.0.0:$SERVER_PORT/ws

### JRPC calls

         websocat --jsonrpc -b ws://0.0.0.0:$SERVER_PORT/ws
         
message format: 
            
         BroadcastToAll "hello"
(now only `BroadcastToAll` method supported)



# Structure 
![image](https://user-images.githubusercontent.com/59182467/128648205-88b0217a-f0ff-4169-a102-dd12322a35ce.png)



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
