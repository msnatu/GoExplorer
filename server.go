package main

import (
	"fmt"
	"net/http"
	"html/template"
	"io/ioutil"
	"os"
//	"encoding/json"
	"encoding/xml"
//	"reflect"
)


const (
	endpoint = "https://api.flickr.com/services/rest/?"
	get_photos_method = "flickr.photos.search"
	api_key = "4ef2fe2affcdd6e13218f5ddd0e2500d"
	output_format = "json"
	results_per_page = "10"
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

var page_templates = template.Must(template.ParseFiles(
	"./tpl/head.html",
	"./tpl/page_body.html"))


/********************************************************/
func main() {
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("www"))))
	http.HandleFunc("/", loadHomePage)
	fmt.Printf("tested!!");
	http.ListenAndServe(":8080", nil)
}

func loadHomePage(w http.ResponseWriter, r *http.Request) {
	p := &HomePage{Title: "Title", Body: "Body"}
	renderTpl(w, "head", p)
	getPuppies(w, p)
}

func renderTpl(w http.ResponseWriter, tmpl string, p *HomePage) {
	err := page_templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPuppies(w http.ResponseWriter, p *HomePage) {
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

		// Test data
//		var testXML string = `
//							<rsp>
//								<photos>
//									<photo id="99448848" owner="Landis"/>
//									<photo id="43453534" owner="John"/>
//								</photos>
//							</rsp>
//							`
		data := &FlickrResponse{}
		xml.Unmarshal(contents, data)
		if(len(data.Photos) > 0) {
			photo_collection := data.Photos[0]
			for i := 0; i < len(photo_collection.Photo); i++ {
				img_data := photo_collection.Photo[i];
				img_url := "https://farm" + img_data.FarmId + ".staticflickr.com/" + img_data.ServerId + "/" + img_data.Id + "_" + img_data.Secret + "_q.jpg"
				fmt.Println(img_url)
				img_elem := "<img src='" + img_url + "'>"
				fmt.Fprint(w, img_elem)
			}
		}
	}
}