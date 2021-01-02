package store

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/przebro/databazaar/collection"
	"github.com/przebro/databazaar/store"
	o "github.com/przebro/databazaar/store"
	mongodb "github.com/przebro/mongostore/collection"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	mongostore     = "mongodb"
	databaseOption = "database"
)

type MongoStore struct {
	name     string
	client   *mongo.Client
	database *mongo.Database
}

func init() {
	store.RegisterStoreFactory(mongostore, initMongoDB)
}

func initMongoDB(opt o.ConnectionOptions) (store.DataStore, error) {

	uri := fmt.Sprintf("%s://%s:%d", opt.Scheme, opt.Host, opt.Port)

	if opt.Path == "" {
		return nil, errors.New("database name required")

	}
	mongoOptions := options.Client().ApplyURI(uri)

	if opt.Options[o.UsernameOption] != "" && opt.Options[o.PasswordOption] != "" {
		credential := options.Credential{
			Username: opt.Options[o.UsernameOption],
			Password: opt.Options[o.PasswordOption],
		}
		mongoOptions.SetAuth(credential)
	}

	ctx := context.Background()

	client, err := mongo.NewClient(mongoOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	database := client.Database(opt.Path)

	return &MongoStore{name: opt.Path, client: client, database: database}, nil
}

func (s *MongoStore) Collection(ctx context.Context, name string) (collection.DataCollection, error) {

	names, err := s.database.ListCollectionNames(ctx, bson.D{{"name", name}})

	if err != nil || len(names) != 1 {
		return nil, fmt.Errorf("unable to find collection:%s", name)
	}

	return mongodb.Collection(name, ctx, s.database), nil
}

func (s *MongoStore) CreateCollection(ctx context.Context, name string) (collection.DataCollection, error) {

	err := s.database.CreateCollection(ctx, name)
	return mongodb.Collection(name, ctx, s.database), err
}

//Status - returns status of a connection
func (s *MongoStore) Status(ctx context.Context) (string, error) {

	err := s.client.Ping(ctx, readpref.PrimaryPreferred())

	return "", err
}

//Close - closes connection
func (s *MongoStore) Close(ctx context.Context) {
	s.client.Disconnect(ctx)
}
