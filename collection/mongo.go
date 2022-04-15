package collection

import (
	"context"
	"encoding/json"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/przebro/databazaar/collection"
	"github.com/przebro/databazaar/result"
	"github.com/przebro/databazaar/selector"
	query "github.com/przebro/mongostore/query"
)

type mongoDbCollection struct {
	name       string
	database   *mongo.Database
	collection *mongo.Collection
}

//Collection - returns an instance of Collection
func Collection(name string, ctx context.Context, database *mongo.Database) collection.DataCollection {

	col := database.Collection(name)

	return &mongoDbCollection{name: name, database: database, collection: col}
}

func (d *mongoDbCollection) Create(ctx context.Context, document interface{}) (*result.BazaarResult, error) {

	var res *result.BazaarResult = &result.BazaarResult{}

	id, _, err := collection.RequiredFields(document)
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, collection.ErrEmptyOrInvalidID
	}

	doc, err := bson.Marshal(document)

	if err != nil {
		return nil, err
	}

	result, err := d.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	res.ID = result.InsertedID.(string)

	reflect.TypeOf(result.InsertedID)

	return res, nil
}
func (d *mongoDbCollection) CreateMany(ctx context.Context, docs []interface{}) ([]result.BazaarResult, error) {

	res, err := d.collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	rset := []result.BazaarResult{}
	for _, x := range res.InsertedIDs {

		rset = append(rset, result.BazaarResult{ID: x.(string)})
	}

	return rset, nil
}
func (d *mongoDbCollection) Get(ctx context.Context, id string, result interface{}) error {

	filter := bson.D{{"_id", id}}

	_, err := collection.IsStruct(result)
	if err != nil {
		return err
	}

	res := d.collection.FindOne(ctx, filter)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return collection.ErrNoDocuments
		}
		return res.Err()
	}

	err = res.Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (d *mongoDbCollection) All(ctx context.Context) (collection.BazaarCursor, error) {

	filter := bson.D{}
	crsr, err := d.collection.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	return &mongoCursor{crsr}, nil
}

func (d *mongoDbCollection) Select(ctx context.Context, s selector.Expr, fld selector.Fields) (collection.BazaarCursor, error) {

	var err error

	builder := query.NewBuilder()
	sxpr := builder.Build(s)
	query := bson.M{}
	json.Unmarshal([]byte(sxpr), &query)

	crsr, err := d.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	return &mongoCursor{crsr}, nil
}
func (d *mongoDbCollection) Update(ctx context.Context, doc interface{}) error {

	id, _, err := collection.RequiredFields(doc)
	if err != nil {
		return err
	}
	if id == "" {
		return collection.ErrEmptyOrInvalidID
	}

	filter := bson.D{{"_id", id}}

	withUpsert := true
	_, err = d.collection.ReplaceOne(ctx, filter, doc, &options.ReplaceOptions{Upsert: &withUpsert})

	return err
}

func (d *mongoDbCollection) BulkUpdate(ctx context.Context, docs []interface{}) error {

	models := []mongo.WriteModel{}
	var id string
	var err error

	for _, doc := range docs {

		if id, _, err = collection.RequiredFields(doc); err != nil {
			return err
		}
		m := mongo.NewReplaceOneModel()
		m.Filter = bson.D{{"_id", id}}
		m.SetUpsert(true)
		m.Replacement = doc

		models = append(models, m)
	}

	_, err = d.collection.BulkWrite(ctx, models)

	return err
}

func (d *mongoDbCollection) Delete(ctx context.Context, id string) error {

	if id == "" {
		return collection.ErrEmptyOrInvalidID
	}

	filter := bson.D{{"_id", id}}
	_, err := d.collection.DeleteOne(ctx, filter)
	return err

}

func (d *mongoDbCollection) Count(ctx context.Context) (int64, error) {

	return d.collection.EstimatedDocumentCount(ctx)
}

func (d *mongoDbCollection) AsQuerable() (collection.QuerableCollection, error) {
	return d, nil
}
