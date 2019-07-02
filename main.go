package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	//_ "github.com/mattn/go-sqlite3"
)

type Recipe struct {
	Title       string `json:"title"`
	Href        string `json:"href"`
	Ingredients string `json:"ingredients"`
	Thumbnail   string `json:"thumbnail"`
}

type Recipes struct {
	Title   string   `json:"title"`
	Version string   `json:"version"`
	Href    string   `json:"href"`
	Recipes []Recipe `json:"results"`
}

//PORT port to be used
func main() {

	r := mux.NewRouter()
	http.Handle("/", r)
	r.Queries("i", "{i}", "q", "{q}", "p", "{p}")
	r.Handle("/api/", filterRecipe()).Methods("GET", "OPTIONS")

	port := ""

	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	} else {
		port = "8080"

	}

	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":" + port,
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func filterRecipe() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// params := mux.Vars(r)
		keys := r.URL.Query()

		keyI := keys.Get("i")
		keyQ := keys.Get("q")
		keyP := keys.Get("p")

		// Open our jsonFile
		jsonFile, err := os.Open("recipes.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var recipes Recipes
		json.Unmarshal(byteValue, &recipes)
		var recipesArray []Recipe
		var arrUtilizado []int
		for i := 0; i < len(recipes.Recipes); i++ {
			found := false
			if len(keyI) > 0 {
				recipeIngredients := strings.Split(recipes.Recipes[i].Ingredients, ",")
				paramIngredients := strings.Split(keyI, ",")

				for _, paramIngredient := range paramIngredients {
					for _, recipeIngredient := range recipeIngredients {
						recipeIngredient = strings.TrimSpace(recipeIngredient)
						if strings.Contains(recipeIngredient, paramIngredient) {
							recipesArray = append(recipesArray, recipes.Recipes[i])
							arrUtilizado = append(arrUtilizado, i)
							found = true
							break
						}
					}
					if found {
						break
					}
				}

			}
			if len(keyQ) > 0 {
				recipeTitles := strings.Split(recipes.Recipes[i].Title, " ")
				paramTitles := strings.Split(strings.TrimSpace(keyQ), " ")

				for _, paramTitle := range paramTitles {
					paramTitle = strings.ToLower(paramTitle)
					for _, recipeTitle := range recipeTitles {
						recipeTitle = strings.ToLower(recipeTitle)
						if strings.Contains(recipeTitle, paramTitle) {
							if !inArray(i, arrUtilizado) {
								recipesArray = append(recipesArray, recipes.Recipes[i])
								arrUtilizado = append(arrUtilizado, i)
								found = true
								break
							}
						}
					}
					if found {
						break
					}
				}
			}
		}

		

		page, _ := strconv.Atoi(keyP)
		maxPerPage := 4
		if page <= 0 {
			page = 1
		}
		maxToShow := page * maxPerPage
		tamArrRecipe := len(recipesArray)


		if len(keyI) > 0 || len(keykeyQI) > 0 {
			if tamArrRecipe == 0 {
				recipes.Recipes = recipesArray
			} else if maxToShow > tamArrRecipe {

				if tamArrRecipe-maxToShow <= 0 && page > 1 {
					recipes.Recipes = nil
				} else if tamArrRecipe-maxPerPage <= 0 {
					recipes.Recipes = recipesArray[0:tamArrRecipe]
				} else {
					recipes.Recipes = recipesArray[maxToShow-maxPerPage : tamArrRecipe]
				}
			} else {
				recipes.Recipes = recipesArray[maxToShow-maxPerPage : maxToShow]
			}
		}

		response, _ := json.Marshal(recipes)
		w.Write(response)
	})
}

func inArray(val int, array []int) (exists bool) {
	exists = false
	for _, v := range array {
		if val == v {
			exists = true
			break
		}
	}
	return exists
}
