// fetching data from pokeapi concurrently and non-concurrently.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
)

const API_URL string = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name    string `json:"name"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

func main() {
	start := time.Now()
	pokemon := []string{"vaporeon", "jolteon", "flareon", "espeon", "umbreon", "leafeon", "glaceon", "sylveon", "eevee"}
	// non-concurrent
	for _, name := range pokemon {
		p := fetchPokemon(name)
		printResult(p)
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

	for p := range ch {
		printResult(p)
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

func printResult(p Pokemon) {
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{25, 25}
	asciiArt, err := aic_package.Convert(p.Sprites.FrontDefault, flags)
	if err != nil {
		panic(err)
	}

	fmt.Printf("front sprite: %s\n", asciiArt)
	fmt.Printf("name: %s\nheight: %d\nweight: %d\n", p.Name, p.Height, p.Weight)
}
