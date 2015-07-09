package main

import (
	"fmt"
	"net/http"
	"html/template"
	"io/ioutil"
	"os"
	"encoding/xml"
	"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


const (
	endpoint = "https://api.flickr.com/services/rest/?"
	get_photos_method = "flickr.photos.search"
	api_key = "4ef2fe2affcdd6e13218f5ddd0e2500d"
	results_per_page = "200"
	tags = "cute+puppies"
)

type Photo struct {
	Id  string    `xml:"id,attr"`
	Owner string `xml:"owner,attr"`
	FarmId string `xml:"farm,attr"`
	ServerId string `xml:"server,attr"`
	Secret string `xml:"secret,attr"`
	ImgTitle string `xml:"title,attr"`
}

type Photos struct {
	Photo   []Photo  `xml:"photo"`
}

type FlickrResponse struct {
	XMLName xml.Name `xml:"rsp"`
	Photos   []Photos  `xml:"photos"`
}

type HomePage struct {
	Title string
	Body  string
}

var page_templates = template.Must(template.ParseFiles("./tpl/head.html", "./tpl/page_body.html"))

/********************************************************/

func main() {
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("www"))))
	http.HandleFunc("/", loadHomePage)
	http.HandleFunc("/vote", handleVote)
	fmt.Printf("tested!!");
	http.ListenAndServe(":8080", nil)
}

func loadHomePage(w http.ResponseWriter, r *http.Request) {
	p := &HomePage{Title: "Title", Body: "Body"}
	renderTpl(w, "head", p)
	puppies := getPuppies(w, p)
	fmt.Fprintf(w, puppies)
	renderTpl(w, "page_body", p)
}

func renderTpl(w http.ResponseWriter, tmpl string, p *HomePage) {
	err := page_templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPuppies(w http.ResponseWriter, p *HomePage) string {
	img_elem := ""
	rest_url := endpoint + "&method=" + get_photos_method + "&api_key=" + api_key + "&text="+ tags + "&nojsoncallback=1&per_page=" + results_per_page
	response, err := http.Get(rest_url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		data := &FlickrResponse{}
		xml.Unmarshal(contents, data)
		if(len(data.Photos) > 0) {
			photo_collection := data.Photos[0]
			for i := 0; i < len(photo_collection.Photo); i++ {
				img_data := photo_collection.Photo[i]
				img_url := "https://farm" + img_data.FarmId + ".staticflickr.com/" + img_data.ServerId + "/" + img_data.Id + "_" + img_data.Secret + "_q.jpg"
				existing_votes := getImageVotes(img_data.Id)
				up_vote := "0"
				down_vote := "0"
				if existing_votes[0] != "" {
					up_vote = existing_votes[0]
				}
				if existing_votes[1] != "" {
					down_vote = existing_votes[1]
				}

				img_elem += "<div class='puppy-image-container' img_id='"+ img_data.Id +"'>"
				img_elem += "<img src='" + img_url + "' class='puppy-box'>"
				img_elem += "<div class='puppy-image-upvote' style='display:none;'></div>"
				img_elem += "<div class='puppy-image-upvote-info' style='display:none;'>" + up_vote +"</div>"
				img_elem += "<div class='puppy-image-downvote' style='display:none;'></div>"
				img_elem += "<div class='puppy-image-downvote-info' style='display:none;'>" + down_vote +"</div>"
				img_elem += "</div>"
			}
		}
	}
	return img_elem
}

func getImageVotes(img_id string) []string {
	db, err := sql.Open("mysql", "root:p@ssw0rd@tcp(127.0.0.1:3306)/puppies")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var (
		id int
		image_id int
		up_vote string
		down_vote string
	)
	//Select
	rows, _ := db.Query("SELECT * FROM votes WHERE image_id = ?", img_id)
	defer rows.Close()

	total_rows := 0
	for rows.Next() {
		err = rows.Scan(&id, &image_id, &up_vote, &down_vote)
		if err != nil {
			fmt.Println(err)
		}
		total_rows++
	}

	//Insert
	if total_rows == 0 {
		stmt, _ := db.Prepare("INSERT INTO votes(image_id) VALUES(?)")
		stmt.Exec(img_id)
	}

	return []string{up_vote, down_vote}
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	img_id := r.FormValue("img_id")
	is_upvote := r.FormValue("is_upvote")

	db, err := sql.Open("mysql", "root:p@ssw0rd@tcp(127.0.0.1:3306)/puppies")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var (
		up_vote string
		down_vote string
		resp_vote int
	)

	db.QueryRow("SELECT up_vote, down_vote FROM votes WHERE image_id = ?", img_id).Scan(&up_vote, &down_vote)
	if(is_upvote == "true") {
		var new_up_vote int
		up_vote, _ := strconv.Atoi(up_vote)
		new_up_vote = up_vote + 1
		stmt, _ := db.Prepare("UPDATE votes SET up_vote = ? WHERE image_id = ?")
		stmt.Exec(new_up_vote, img_id)
		resp_vote = new_up_vote
	} else {
		var new_down_vote int
		down_vote, _ := strconv.Atoi(down_vote)
		new_down_vote = down_vote + 1
		stmt, _ := db.Prepare("UPDATE votes SET down_vote = ? WHERE image_id = ?")
		stmt.Exec(new_down_vote, img_id)
		resp_vote = new_down_vote
	}

	type JsonResponse map[string]interface{}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, resp_vote)
}