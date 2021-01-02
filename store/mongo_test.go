package store

import (
	"context"
	"testing"

	"github.com/przebro/databazaar/store"
)

type TestMongoStruct struct {
	ID    string `json:"_id,omitempty" bson:"_id,omitempty"`
	Value string `json:"value"`
	Rev   string `json:"_rev,omitempty" bson:"_rev,omitempty"`
}

type TestDocument struct {
	ID     string  `json:"_id" bson:"_id"`
	REV    string  `json:"_rev,omitempty" bson:"_rev,omitempty"`
	Title  string  `json:"title"`
	Score  float32 `json:"score"`
	Year   int     `json:"year"`
	Oscars bool    `json:"oscars"`
}

const mongoDbConn = "mongodb;127.0.0.1:20017/testdb?username=admin&password=notsecure"

func TestMongoStore(t *testing.T) {

	_, err := store.NewStore("mongodb;127.0.0.1:20017/testdb")

	if err != nil {
		t.Error("unexpected error:", err)
	}

	_, err = store.NewStore("mongodb;127.0.0.1:20017/?username=admin&password=notsecure")

	if err == nil {
		t.Error("unexpected result")
	}
}

func TestStatus(t *testing.T) {

	s, err := store.NewStore(mongoDbConn)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	_, err = s.Status(context.Background())

	if err != nil {
		t.Error("unexpected result:", err)
	}
}

func TestCreateCollection(t *testing.T) {

	store, err := store.NewStore(mongoDbConn)

	_, err = store.Collection(context.Background(), "databazaar")
	if err == nil {
		t.Error(err)
	}

	_, err = store.CreateCollection(context.Background(), "databazaar")
	if err != nil {
		t.Error(err)
	}

	_, err = store.Collection(context.Background(), "databazaar")
	if err != nil {
		t.Error(err)
	}
}

func TestClose(t *testing.T) {

	store, err := store.NewStore(mongoDbConn)
	if err != nil {
		t.Error(err)
	}
	store.Close(context.Background())

}