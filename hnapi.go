/* HNdb is intended to provide a stripped down interface customized to use specifically with Hacker News firbase.
It has been adapted from https://github.com/easyCZ/grpc-web-hacker-news/blob/master/server/hackernews/api.go.
*/
package hnapi

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"gopkg.in/zabawaba99/firego.v1"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const baseURL = "https://hacker-news.firebaseio.com"
const httpURL = "https://news.ycombinator.com/"
const version = "/v0"

var endPoint = map[string]string{
	"top":  "/topstories",
	"new":  "/newstories",
	"best": "/beststories",
	"ask":  "/askstories",
	"jobs": "/jobstories",
	"show": "/showstories",
}

// HNdb has an embedded struct for the firebase interface
type HNdb struct {
	*firego.Firebase
}

// NewHNdb establishes an API to Hacker New's Firebase.
func NewHNdb() *HNdb {
	hnURL := baseURL + version
	return &HNdb{
		firego.New(hnURL, nil),
	}
}

// GetItem retrieves the specified item and parses it.
func (db *HNdb) GetItem(id int) (*Item, error) {
	ref, err := db.Ref(fmt.Sprintf(version+"/item/%d", id))
	if err != nil {
		log.Fatalf("request story reference failed @ reference: %s", err)
	}

	var value Item
	if err := ref.Value(&value); err != nil {
		log.Fatalf("story #%d retrieval failed %s", id, err)
	}

	return &Item{
		ID:          value.ID,
		Deleted:     value.Deleted,
		Type:        value.Type,
		By:          value.By,
		Time:        value.Time,
		Text:        value.Text,
		Title:       value.Title,
		Dead:        value.Dead,
		Parent:      value.Parent,
		Poll:        value.Poll,
		Kids:        value.Kids,
		URL:         value.URL,
		Score:       value.Score,
		Parts:       value.Parts,
		Descendants: value.Descendants,
	}, nil
}

// GetPosts retrieves the specified type and number of posts.
func (db *HNdb) GetPosts(req *Request) (contentChan chan *Item) {
	contentChan = make(chan *Item)

	ref, err := db.Firebase.Ref(endPoint[req.PostType])
	if err != nil {
		log.Fatal("error firebase reference")
	}

	var ids []int

	if err := ref.Value(&ids); err != nil {
		log.Printf("%s stories request failed", req.PostType)
	}

	ids = ids[:req.NumPosts]

	for _, id := range ids {
		go func(id int) {
			item, _ := db.GetItem(id)
			contentChan <- item
		}(id)
	}
	return contentChan
}

// ClientWithAuth returns a client with the the authenticated session cookie
func ClientWithAuth(username, password string) (*http.Client, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}

	//u, _ := url.ParseRequestURI(apiUrl)
	//u.Path = resource
	//urlStr := u.String() // "https://news.ycombinator.com/login"

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.PostForm(httpURL + "login", url.Values{
		"acct": {username},
		"pw": {password},
		"goto": {`news`},
	})
	if err != nil {
		return nil, err
	}
	for _, cookie := range resp.Cookies() {
		fmt.Printf("cookie:  %+v", cookie)
	}

	return client, nil
}