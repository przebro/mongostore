package collection

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCursor struct {
	crsr *mongo.Cursor
}

func (c *mongoCursor) All(ctx context.Context, v interface{}) error {
	return c.crsr.All(ctx, v)
}
func (c *mongoCursor) Next(ctx context.Context) bool {
	return c.crsr.Next(ctx)
}
func (c *mongoCursor) Decode(v interface{}) error {
	return c.crsr.Decode(v)
}
func (c *mongoCursor) Close() error {
	return c.crsr.Close(context.Background())
}
