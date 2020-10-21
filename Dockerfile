FROM golang

WORKDIR /
COPY . /

RUN go get "github.com/gorilla/mux"

EXPOSE 3000

CMD go run main.go