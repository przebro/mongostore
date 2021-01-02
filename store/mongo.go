package store

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

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

//MongoStore - mongodb data store
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

	if capath := opt.Options[o.RootCACertOption]; capath != "" {

		pool := x509.NewCertPool()
		data, err := ioutil.ReadFile(capath)

		if err != nil {
			return nil, err
		}
		ok := pool.AppendCertsFromPEM(data)
		if !ok {
			return nil, errors.New("unable to read certificate")
		}

		var untrusted bool

		if trustopt := opt.Options[o.UntrustedOption]; trustopt != "" {

			if v, err := strconv.ParseBool(trustopt); err == nil {
				untrusted = v
			}
		}

		tlscfg := &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: untrusted,
			ServerName:         "local",
		}

		if ckeyf, ccertf := opt.Options[o.ClientKeyOption], opt.Options[o.ClientCertOption]; ckeyf != "" && ccertf != "" {
			ccert, err := tls.LoadX509KeyPair(ccertf, ckeyf)
			if err != nil {
				return nil, err
			}

			tlscfg.Certificates = []tls.Certificate{ccert}
		}

		mongoOptions.SetTLSConfig(tlscfg)

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

	s.client.ListDatabaseNames(context.Background(), nil)
	err := s.client.Ping(ctx, readpref.PrimaryPreferred())

	return "", err
}

//Close - closes connection
func (s *MongoStore) Close(ctx context.Context) {
	s.client.Disconnect(ctx)
}
