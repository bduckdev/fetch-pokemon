// fetching data from pokeapi concurrently and non-concurrently.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const API_URL string = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

func main() {
	// non-concurrent
	start := time.Now()
	pokemon := []string{"vaporeon", "jolteon", "flareon", "espeon", "umbreon", "leafeon", "glaceon", "sylveon", "eevee"}
	for _, name := range pokemon {
		p := fetchPokemon(name)
		fmt.Printf("name: %s\nheight: %d\nweight: %d\n", p.Name, p.Height, p.Weight)
	}
	fmt.Println("Fetching without concurrency took ", time.Since(start))
	// concurrent
	start = time.Now()
	ch := make(chan Pokemon)
	var wg sync.WaitGroup
	for _, name := range pokemon {
		wg.Add(1)
		go fetchPokemonConcurrently(name, ch, &wg)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		fmt.Printf("name: %s\nheight: %d\nweight: %d\n", res.Name, res.Height, res.Weight)
	}
	fmt.Println("Fetching with concurrency took ", time.Since(start))
}

func fetchPokemon(name string) Pokemon {
	url := API_URL + name
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var p Pokemon
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		panic(err)
	}

	return p
}

func fetchPokemonConcurrently(name string, ch chan<- Pokemon, wg *sync.WaitGroup) {
	defer wg.Done()

	url := API_URL + name
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var p Pokemon
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		panic(err)
	}

	ch <- p
}
