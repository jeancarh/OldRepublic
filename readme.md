# Instalation
to start using OldRepublic locally you need to follow this steps

1. download the go file installation for your OS at https://golang.org/doc/install
2. check golang installation with $ go version```
3. install the http server library ```go get -u github.com/gorilla/mux ```


## How to run the project

1. first open the project
2. Next Press F5 to run the server
3. on postman put this url localhost:8000/topsecret for to return the ship location with the next body on method POST, plase make sure that all the message arrays have the same lenghts.
```
{
"satellites": [
{
"name": "kenobi",
"distance": 100.0,
"message": ["", "", "", "mensaje", ""]
},
{
"name": "skywalker",
"distance": 115.5,
"message": ["", "es", "", "", "secreto"]
},
{
"name": "sato",
"distance": 142.7,
"message": ["este", "", "un", "", ""]
}
]
}
```
4. to send data for only one known location or satellite sent POST or GET to this url ```localhost:8000/topsecret/{satellite_name}```

```
{
"distance": 100.0,
"message": ["este", "", "", "mensaje", ""]
}
```

5. another way to start the server is with ```go run main.go``` the project is going to start to hear on  ```localhost:8000```

## How to use the productive project on aws

1. follow the steps of the body in the "How to run the project" steps, the only thing that change is the url of EC2 instance on aws

```
http://ec2-3-133-20-139.us-east-2.compute.amazonaws.com:3000/topsecret/kenobi
```

## about libs
1. gorilla/mux for http server
2. math from go to math operations
3. json decoder
4. errors to manage errors on the app
5. logs to make easier the track on console
