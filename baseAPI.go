package main

/*
	importing the packages
*/
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

/*
	Globle Variables,type,struct are defines
*/
type fruits map[string]int
type vegetables map[string]int

type payload struct {
	Stuff data
}
type data struct {
	Fruit   fruits
	Veggies vegetables
}

/*
MyServer - this is a struct defined for server
*/
type MyServer struct {
	r *mux.Router
}

/*
Root - This is the function which is called when base url is called
*/
func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server Running")
}

/*
GetJSONResponse - This is the function which is called when getjson url is called
					This is demo, How to return a JSON object
*/
func GetJSONResponse(w http.ResponseWriter, r *http.Request) {
	fruits := make(map[string]int)
	fruits["Apples"] = 25
	fruits["Oranges"] = 10

	vegetables := make(map[string]int)
	vegetables["Carrats"] = 10
	vegetables["Beets"] = 0

	d := data{fruits, vegetables}
	p := payload{d}
	data := p
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

/*
startServer - This function will starts the server.
				We define the routes in the function.
*/
func startServer(port int, url string) {
	/*
		Checks the port and url variables values
	*/
	if port <= 0 || len(url) == 0 {
		panic("invalid port or url")
	}

	/*
		Defines and prints the variable fullURL
	*/
	fullURL := fmt.Sprintf("%s:%d", url, port)
	fmt.Printf("starting server on %s\n", fullURL)

	/*
		Defines a Router
	*/
	rm := mux.NewRouter()

	/*
		initial splash screen
	*/

	/*
		Defines routes in the API
	*/
	rm.HandleFunc("/", Root).Methods("GET")
	rm.HandleFunc("/getjson", GetJSONResponse).Methods("GET")

	/*
		Starts the API
		Set the port
	*/
	http.ListenAndServe(fullURL, &MyServer{rm})

}

/*
ServeHTTP - This set the parameter and configures the server parameter
*/
func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

/*
main - This is main starting point
*/
func main() {
	/*
		Loads the env variables
	*/
	gotenv.Load()

	/*
		Get the port no from .env file.
		Convert string to int
		In case some error comes then process is stopped
	*/
	port, err := strconv.Atoi(os.Getenv("WEBSITE_PORT"))
	if err != nil {
		fmt.Println("port value is invalid")
		return
	}

	/*
		Gets the website ip from .env file.
	*/
	url := os.Getenv("WEBSITE_IP")

	/*
		calls the function and starts the server.
		start listeneing on port and url
	*/
	startServer(port, url)
}
