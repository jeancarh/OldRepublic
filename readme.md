# Old Republic

This is a golang project to identify location based on 3 known locations and 3 known distances from the object

## Installation

1.download the go file installation for your OS at https://golang.org/doc/install
2. check golang installation with $ go version```

## How to run the server locally
1. first open the project
2. Next Press F5 to run the server
3. on postman put this url ```localhost:8000/topsecret``` for to return the ship location with the next body on method POST, plase make sure that all the message arrays have the same lenghts.
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

## about libs

1. gorilla/mux for http server
2. math from go to math operations
3. json decoder
4. errors to manage errors on the app
5. logs to make easier the track on console

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](https://choosealicense.com/licenses/mit/)