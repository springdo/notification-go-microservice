# Notification Microservice App
Test Receiver=["thirdpartybuilderz@mailinator.com"]

## What does this do?
It will expose on api that receives a JSON payload and either emails or sends item to Jira

## run and build
1. to install app dependencies `go get -v ./...`
1. to run the app `go run *.go`
1. to build & run the app `go build && ./notification-microservice` or `go build && ./$APP_NAME.exe`

## run on Mac OS
2. Setup `$GOPATH`
2. Create `redhat` repo space 
2. Run these instructions if you're Tom or Dónal
```
mkdir -p ~/Documents/go/src/redhat
export GOPATH=$HOME/Documents/go
ln -s /Absolute/path/to/notification-microservice $GOPATH/src/redhat/notification-microservice  
```


## other useful OC commands
1. `oc get pods --show-labels`
1. `oc rsh ${POD}`
1. `oc login https://console.xpertex.rht-labs.com -u ${OCP_USER} -p ${OCP_PASS}`