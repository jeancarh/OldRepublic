package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Satellite Satellites
type Satellite struct {
	Name     string   `json:"name"`
	Distance float64  `json:"distance"`
	Message  []string `json:"message"`
}

// Satellites Satellites
type Satellites struct {
	Satellites []Satellite `json:"satellites"`
}

// BodyResponse BodyResponse
type BodyResponse struct {
	Position Position `json:"position"`
	Message  string   `json:"message"`
}

// Position Position
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// BodySplitRequest BodySplitRequest
type BodySplitRequest struct {
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

// struct for use on trasliteration method
type point struct {
	x float64
	y float64
}

func norm(px, py float64) float64 {
	return math.Pow(math.Pow(px, 2)+math.Pow(py, 2), .5)
}

// this methos is being used to handle the POST request for /topsecret
func yourHandlerTs(w http.ResponseWriter, r *http.Request) {
	var satellite Satellites
	err := json.NewDecoder(r.Body).Decode(&satellite)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message, err := getMessage(satellite.Satellites[0].Message, satellite.Satellites[1].Message, satellite.Satellites[2].Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	if message != "" {
		locationX, locationY, err := getLocation(satellite.Satellites[0].Distance, satellite.Satellites[1].Distance, satellite.Satellites[2].Distance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		saveRedisData()
		response := buildResponse(message, locationX, locationY)
		fmt.Fprintf(w, "%v\n", response)
	}
}

// this methos is being used to handle the query param GET or POST request for /topsecret/{satellite}
func yourHandlerSplit(w http.ResponseWriter, r *http.Request) {
	var bodyRQ BodySplitRequest
	err := json.NewDecoder(r.Body).Decode(&bodyRQ)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.Method == "POST" {
		params := mux.Vars(r)
		sattelite := params["satellite"]
		fmt.Fprintf(w, "Data sent to %v\n", sattelite)
	}
	if r.Method == "GET" {
		_, _, err := getLocationSplit(bodyRQ.Distance)
		if err != nil {
			fmt.Fprintf(w, "there is not enough data")
			http.Error(w, err.Error(), http.StatusNotFound)
		}

	}
}

func main() {
	r := mux.NewRouter()
	// Routes r is the way to route path or handle funcions using gorilla/Mux
	r.Path("/topsecret").HandlerFunc(yourHandlerTs).Methods("POST")
	r.HandleFunc("/topsecret/{satellite}", yourHandlerSplit).Methods("POST", "GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}

// this method build the msg and return the error 404 where there is not data in the 3 arrays
func getMessage(messages ...[]string) (string, error) {
	kenobi := messages[0]
	skywalker := messages[1]
	sato := messages[2]
	completeMsg := []string{}
	for i := range kenobi {
		if kenobi[i] == "" && skywalker[i] == "" && sato[i] == "" {
			return "", errors.New("")
		}
		if kenobi[i] != "" {
			completeMsg = append(completeMsg, kenobi[i])
		}
		if skywalker[i] != "" {
			completeMsg = append(completeMsg, skywalker[i])
		}
		if sato[i] != "" {
			completeMsg = append(completeMsg, sato[i])
		}
	}
	result := strings.Join(unique(completeMsg), " ")
	return result, nil
}

// this method is used to delete repeated data on the final msg
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func saveRedisData() {

}

// this method use an math script to find the location
func getLocation(distances ...float64) (float64, float64, error) {
	var finalPose point
	kenobi := []float64{-500, -200}
	skywalker := []float64{100, -100}
	sato := []float64{500, 100}
	finalPose = trilateration(kenobi, skywalker, sato, distances[0], distances[1], distances[2])
	return finalPose.x, finalPose.y, nil
}

// this method return err because with only one location cannot be found the ship
func getLocationSplit(distances float32) (float32, float32, error) {
	return 0, 0, errors.New("")
}

//this method build the response for Post on top_secret query
func buildResponse(msg string, x, y float64) string {
	position := Position{
		X: x,
		Y: y,
	}
	bodyRS := BodyResponse{Message: msg, Position: position}
	b, err := json.Marshal(bodyRS)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
	return string(b)
}

func trilateration(point1 []float64, point2 []float64, point3 []float64, d1 float64, d2 float64, d3 float64) point {
	var resultPose point
	//unit vector in a direction from point1 to point 2
	p2p1Distance := math.Pow(math.Pow(point2[0]-point1[0], 2)+math.Pow(point2[1]-point1[1], 2), 0.5)
	ex1 := (point2[0] - point1[0]) / p2p1Distance
	ex2 := (point2[1] - point1[1]) / p2p1Distance
	aux1 := point3[0] - point1[0]
	aux2 := point3[1] - point1[1]
	//signed magnitude of the x component
	i := ex1*aux1 + ex2*aux2
	//the unit vector in the y direction.
	aux3 := point3[0] - point1[0] - i*ex1
	aux4 := point3[1] - point1[1] - i*ex2
	eyx := aux3 / norm(aux3, aux4)
	eyy := aux4 / norm(aux3, aux4)
	//the signed magnitude of the y component
	j := eyx*aux3 + eyy*aux4
	//coordinates
	x := (math.Pow(d1, 2) - math.Pow(d2, 2) + math.Pow(p2p1Distance, 2)) / (2 * p2p1Distance)
	y := (math.Pow(d1, 2)-math.Pow(d3, 2)+math.Pow(i, 2)+math.Pow(j, 2))/(2*j) - i*x/j
	//result coordinates
	finalX := point1[0] + x*ex1 + y*eyx
	finalY := point1[1] + x*ex2 + y*eyy
	resultPose.x = finalX
	resultPose.y = finalY
	return resultPose
}
