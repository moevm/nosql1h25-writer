package mongodb

import (
	"context"
	"time"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectTimeout = 7 * time.Second

func New(uri string) (mongoifc.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongoifc.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
