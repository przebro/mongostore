package collection

import (
	"context"
	"fmt"
	"testing"

	"github.com/przebro/databazaar/collection"
	tst "github.com/przebro/databazaar/collection/testing"
	"github.com/przebro/databazaar/selector"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var col collection.DataCollection

type InvalidDocument struct {
	ID    int    `json:"_id,omitempty"`
	REV   string `json:"_rev,omitempty"`
	Title string `json:"title"`
}

var (
	singleRecord        = tst.TestDocument{Title: "Blade Runner", Score: 8.1, Year: 1982, Oscars: false}
	singleRecordInvalid = InvalidDocument{ID: 12345, Title: "Invalid Title"}
	singleRecordWithID  tst.TestDocument
	testCollection      []tst.TestDocument
	cli                 *mongo.Client
)

const host = "127.0.0.1"
const port = 20017
const username = "admin"
const passwd = "notsecure"
const dbname = "testdb"

func init() {
	var err error = nil
	singleRecordWithID, testCollection = tst.GetSingleRecord("../data/testdata.json")

	uri := fmt.Sprintf("mongodb://%s:%d", host, port)

	credential := options.Credential{Username: username, Password: passwd}
	mongoOptions := options.Client().ApplyURI(uri).SetAuth(credential)

	cli, err = mongo.NewClient(mongoOptions)

	if err != nil {
		panic("")
	}
	err = cli.Connect(context.Background())

	db := cli.Database("testdb")

	col = Collection("databazaar", context.Background(), db)

}

func TestMain(m *testing.M) {

	m.Run()
	cli.Database("testdb").Drop(context.Background())
}

func TestInsertOne(t *testing.T) {

	r, err := col.Create(context.Background(), &singleRecord)

	if err == nil {
		t.Error(err)
	}

	if err != collection.ErrEmptyOrInvalidID {
		t.Error(err)
	}

	r, err = col.Create(context.Background(), &singleRecordWithID)
	if err != nil {
		t.Error(err)
	}

	if r.ID != singleRecordWithID.ID {
		t.Error("unexpected result, expected:", singleRecordWithID.ID, "actual:", r.ID)
	}

}

func TestGetOne(t *testing.T) {

	doc := tst.TestDocument{}
	err := col.Get(context.Background(), "single_record", &doc)
	if err != nil {
		t.Error(err)
	}

	if doc.ID != singleRecordWithID.ID {
		t.Error("unexpected result:", doc.ID)
	}
}

func TestInsertMany(t *testing.T) {

	tc := []interface{}{}

	for x := range testCollection {
		tc = append(tc, testCollection[x])
	}

	_, err := col.CreateMany(context.Background(), tc)
	if err != nil {
		t.Error(err)
	}
}

func TestSelectMany(t *testing.T) {

	result := []tst.TestDocument{}
	sel := selector.Gte("year", selector.Int(1986))
	qr, _ := col.AsQuerable()
	crsr, err := qr.Select(context.Background(), sel, nil)
	if err != nil {
		t.Error(err)
	}

	crsr.All(context.Background(), &result)

	if len(result) != 5 {
		t.Error("unexpected result: collection len expected:", 5, "actual:", len(result))
	}

	cnt := 0
	qr, _ = col.AsQuerable()
	crsr, err = qr.Select(context.Background(), sel, nil)
	for crsr.Next(context.Background()) {
		doc := tst.TestDocument{}
		crsr.Decode(&doc)
		cnt++
	}

	if len(result) != 5 {
		t.Error("unexpected result: collection len expected:", 5, "actual:", len(result))
	}

	err = crsr.Close()
	if err != nil {
		t.Error("Unexpected result:", err)
	}
}

func TestUpdate(t *testing.T) {

	singleRecord.Score = 7.3

	err := col.Update(context.Background(), &singleRecord)
	if err == nil {
		t.Error("unexpected result")
	}

	doc := tst.TestDocument{
		ID:     "movie_13",
		Oscars: true,
		Score:  7.9,
		Year:   1999,
		Title:  "The Matrix",
	}
	result, err := col.Create(context.Background(), &doc)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	doc.ID = result.ID
	doc.REV = result.Revision
	doc.Score = 2.3
	err = col.Update(context.Background(), &doc)
	if err != nil {
		t.Error("unexpected result:", err)
	}

}

func TestDelete(t *testing.T) {

	err := col.Delete(context.Background(), singleRecordWithID.ID)

	if err != nil {
		t.Error("unexpected result:", err)
	}

}
