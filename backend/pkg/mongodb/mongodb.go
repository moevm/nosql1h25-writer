package mongodb

import (
	"context"
	"time"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectTimeoutSecond = 7

func New(uri string) (mongoifc.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*connectTimeoutSecond)
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
