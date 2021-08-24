package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// Pokemon struct (PascalCase)
type Pokemon struct {
	Id      string   `json:"id"`
	Name    string   `json:"pokeName"`
	Element string   `json:"pokeElement"`
	Trainer *Trainer `json:"trainer"`
}

// Trainer struct (PascalCase)
type Trainer struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// creating pokemons slice -> An array has a fixed size.A slice, on the other hand, is a dynamically-size
var pokemons []Pokemon

// CRUD functions

func getPokemons(w http.ResponseWriter, _ *http.Request) {
	// response from the server the client can accept
	w.Header().Set("Content Type", "application/json")

	// we are sending our request with 'w' and we are getting pokemons as result
	json.NewEncoder(w).Encode(pokemons)
}

func getPokemonById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conten Type", "application/json")

	// we need parameter because we want to return pokemon with given id
	params := mux.Vars(r)

	// we need to search our given params (id) in pokemons
	for _, item := range pokemons {
		if item.Id == params["id"] {
			// we will see as result only one pokemon not pokemons
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content Type", "application/json")

	// Create a pokemon with Pokemon struct
	var pokemon Pokemon

	// we are sending our request as body cuz we are creating a pokemon
	// then we are putting the result in our pokemon variable -> '&' == variable's memory address.
	_ = json.NewDecoder(r.Body).Decode(&pokemon)

	// we are creating our pokemons id with random method
	// so when we are sending our body we don't need to give an id
	pokemon.Id = strconv.Itoa(rand.Intn(99999999999))

	// we are adding our new pokemon to our pokemon slice
	pokemons = append(pokemons, pokemon)

	// we will see our created pokemon as result
	json.NewEncoder(w).Encode(pokemon)

}

func deletePokemonById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content Type", "application/json")

	// we need a parameter
	params := mux.Vars(r)
	for index, item := range pokemons {

		if item.Id == params["id"] {

			/* for example;
			if our id is 3 then our index will be 2
			0,1 will remain and then 2 removed and our other slice will start with index 3
			final result will be as indexes 0,1,3,4,5,6,7...*/
			pokemons = append(pokemons[:index], pokemons[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(pokemons)
}

func updatePokemonById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content Type", "application/json")

	// for our body we need a parameter
	params := mux.Vars(r)

	for index, item := range pokemons {
		if item.Id == params["id"] {

			// like delete func we are doing the same and removing our given id from our slice
			pokemons = append(pokemons[:index], pokemons[index+1:]...)

			// then we are creating new one
			var pokemon Pokemon
			_ = json.NewDecoder(r.Body).Decode(&pokemon)

			// we are using the same id
			pokemon.Id = params["id"]

			// and then we are append it to our slice
			pokemons = append(pokemons, pokemon)

			// updated pokemon will be our result
			json.NewEncoder(w).Encode(pokemon)
			return
		}
	}
}

func main() {

	// Router
	r := mux.NewRouter()

	// Dump datas
	pokemons = append(pokemons, Pokemon{
		"1",
		"Pikachu",
		"Electricity",
		&Trainer{"Ash", "Ketchum"},
	})

	pokemons = append(pokemons, Pokemon{
		"2",
		"Charmender",
		"Flame",
		&Trainer{"Ash", "Ketchum"},
	})

	pokemons = append(pokemons, Pokemon{
		"3",
		"Psyduck",
		"Water",
		&Trainer{"Misty", "Williams"},
	})

	pokemons = append(pokemons, Pokemon{Id: "4", Name: "Onix", Element: "Rock",
		Trainer: &Trainer{FirstName: "Brock", LastName: "Harrison"}})

	// CRUD opertaions
	r.HandleFunc("/api/pokemons", getPokemons).Methods("GET")
	r.HandleFunc("/api/pokemon/{id}", getPokemonById).Methods("GET")
	r.HandleFunc("/api/pokemons", createPokemon).Methods("POST")
	r.HandleFunc("/api/pokemon/{id}", updatePokemonById).Methods("PUT")
	r.HandleFunc("/api/pokemon/{id}", deletePokemonById).Methods("DELETE")

	// For avoid No 'Access-Control-Allow-Origin' error.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	// after go build we will see this message
	fmt.Println("Server started at port 8000")

	// automaticly will catch the error if there is one
	log.Fatal(http.ListenAndServe(":8000", handler))
}
