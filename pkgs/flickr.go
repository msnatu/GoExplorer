package flickr

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
)

const (
	endpoint = "https://api.flickr.com/services/rest/?"
	get_photos_method = "flickr.photos.search"
	api_key = "4ef2fe2affcdd6e13218f5ddd0e2500d"
	output_format = "json"
	results_per_page = "10"
	tags = "cute+puppies"
)

func test() {
	image_url := endpoint + "&method=" + get_photos_method + "&api_key=" + api_key + "&text="+ tags + "&format=" + output_format + "&nojsoncallback=1&per_page=" + results_per_page
	response, err := http.Get(image_url)
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
		fmt.Printf("%s\n", string(contents))
	}
}