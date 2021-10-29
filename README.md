# Websockets Manager

Allows clients to subscribe to notifications by creating a websocket connection
with the server. The server also hosts an HTTP endpoint that will accept
a message and notify all connected clients of this message.

## Start Server

`go run ws/hub/main.go`

## Connect some Clients

`go run ws/client/main.go -userid=foo`
`go run ws/client/main.go -userid=bar`

## Send a Message

`curl -v localhost:8080/send -H 'X-Compass-WS-User:foo' -d '{"message": "hey foo"}'`
`curl -v localhost:8080/send -H 'X-Compass-WS-User:bar' -d '{"message": "hey bar"}'`

## Send a listing event
`curl -v localhost:8080/mls -H 'X-Compass-WS-User:foo' -d '{"city": "MyCity", "Zipcode" : "0982-2232", "address": "My address \n street name", "email": "user@compass.com"}'`

## Deploy and Run

### Deploy binary
```
env GOOS=linux GOARCH=amd64 go build ws/hub/main.go
scp -i ~/.ssh/websockets.pem main ec2-user@3.237.18.207:~
ssh -i .ssh/websockets.pem ec2-user@3.237.18.207
./main -addr 0.0.0.0:8080 &
```
### Deploy source
```
HACKATHON_SRC=hackathon-websockets-2020
scp -r -i ~/.ssh/websockets.pem <hackathon-src> ec-user@3.237.18.207:~/$HACKATHON_SRC
ssh -i .ssh/websockets.pem ec2-user@3.237.18.207
cd $HACKATHON_SRC; go run ws/hub/main.go -addr 0.0.0.0:8080 &
```

## Run Agent Home locally
```
cd development/uc-frontend/
git checkout alice.yoon/hackathon2020-websocket
cd apps/agent-home
./pnpm install
./pnpm run dev
```
