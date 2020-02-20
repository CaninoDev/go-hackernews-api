package hnapi

import (
	"gopkg.in/zabawaba99/firego.v1"
	"net/url"
	"reflect"
	"testing"
)

func TestClientWithAuth(t *testing.T) {


	client, err := ClientWithAuth(username, password)
	if err != nil {
		t.Errorf("Got error #{err}")
		return
	}

	hnURL, err := url.Parse("https://news.ycombinator.com")
	if err != nil {
		t.Errorf("Got error #{err}")
	}

	for _, cookie := range client.Jar.Cookies(hnURL) {
		if (cookie.Name == "user") && (cookie.Value != "") {
			return
		}
	}
	t.Error("http.Client did not get authentication cookie")
}

func TestHNdb_GetItem(t *testing.T) {
	type fields struct {
		Firebase *firego.Firebase
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Item
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &HNdb{
				Firebase: tt.fields.Firebase,
			}
			got, err := db.GetItem(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHNdb_GetPosts(t *testing.T) {
	type fields struct {
		Firebase *firego.Firebase
	}
	type args struct {
		req *Request
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantContentChan chan *Item
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &HNdb{
				Firebase: tt.fields.Firebase,
			}
			if gotContentChan := db.GetPosts(tt.args.req); !reflect.DeepEqual(gotContentChan, tt.wantContentChan) {
				t.Errorf("GetPosts() = %v, want %v", gotContentChan, tt.wantContentChan)
			}
		})
	}
}

func TestNewHNdb(t *testing.T) {
	tests := []struct {
		name string
		want *HNdb
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHNdb(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHNdb() = %v, want %v", got, tt.want)
			}
		})
	}
}