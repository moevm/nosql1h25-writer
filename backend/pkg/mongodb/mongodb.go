package mongodb

import (
	"context"
	"time"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionTimeout = 7 * time.Second
const connectionAttempts = 3

func New(uri string) (mongoifc.Client, error) {
    var client mongoifc.Client
    var err error
    currentAttempt := 0

    for currentAttempt < connectionAttempts {
        client, err = connectToMongo(uri)
        if err == nil {
            return client, nil
        }
        currentAttempt++
    }
    return nil, err
}

func connectToMongo(uri string) (mongoifc.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
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