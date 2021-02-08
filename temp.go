package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"
)

type Filmss struct {
	Results []Films `json : "result"`
}

type Films struct {
	Title         string   `json: "title"`
	Opening_crawl string   `json: "opening_crawl"`
	Director      string   `json: "director"`
	Producer      string   `json: "prodcuer"`
	Release_date  string   `json: "release_date"`
	Characters    []string `json: "characters"`
}

type Chara struct {
	Name   string `json: "name"`
	Gender string `json:"gender"`
	Height string `json:"height"`
}

type Person struct {
	Name  string
	Age   int
	Color string
}

type Comment struct {
	Post string
}

var responseObject Filmss
var responseObject2 Chara

func readForm(r *http.Request) *Comment {
	r.ParseForm()
	comment := new(Comment)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(comment, r.PostForm)
	if decodeErr != nil {

		log.Printf("error mapping parsed form data to struct : ",
			decodeErr)

	}
	return comment
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	// person := Person{
	// 	Name:  "tolu",
	// 	Age:   25,
	// 	Color: "black",
	// }
	if r.Method == "GET" {

		parsedTemplate, _ := template.ParseFiles("template/index.html")
		parsedTemplate.Execute(w, responseObject)

	} else {
		comment := readForm(r)
		fmt.Fprint(w, comment.Post)

	}

}

func main() {

	// database

	// db, err := sql.Open("mysql", "root:@tcp(localhost)/test")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// defer db.Close()

	// insert, err := db.Query("INSERT INTO test VALUES (8,'TooEiSTt')")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// defer insert.Close()

	response, err := http.Get("https://swapi.dev/api/films/")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(responseData))

	json.Unmarshal(responseData, &responseObject)

	fmt.Println(len(responseObject.Results))
	for i := 0; i < len(responseObject.Results); i++ {
		fmt.Println(responseObject.Results[i].Title, responseObject.Results[i].Director, responseObject.Results[i].Opening_crawl)
		for _, i := range responseObject.Results[i].Characters {
			response, err := http.Get(i)
			if err != nil {
				continue
			}
			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(string(responseData))

			json.Unmarshal(responseData, &responseObject2)

			fmt.Println(responseObject2.Gender, responseObject2.Height, responseObject2.Name)
		}

	}

	fmt.Println("it gets here")

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", renderTemplate)
	er := http.ListenAndServe(":8085", nil)
	if er != nil {
		log.Fatal("error starting http server : ", er)
		return
	}

}
