package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vonhraban/secret-server/secret"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoCollectionName = "secret"

type mongoVault struct {
	clock      secret.Clock
	collection *mongo.Collection
}

type mongoSecret struct {
	Hash           string     `bson:"_id"`
	SecretText     string     `bson:"text"`
	RemainingViews int        `bson:"remaining_views"`
	CreatedAt      time.Time  `bson:"created_at"`
	ExpiresAt      *time.Time `bson:"expires_at"`
}

func (s *mongoSecret) toDomainSecret() *secret.Secret {
	var expiresAt time.Time
	if s.ExpiresAt != nil {
		expiresAt = *s.ExpiresAt
	}
	return &secret.Secret{
		Hash:           s.Hash,
		SecretText:     s.SecretText,
		RemainingViews: s.RemainingViews,
		CreatedAt:      s.CreatedAt,
		ExpiresAt:      expiresAt,
	}
}

func mongoSecretFromDomainSecret(d *secret.Secret) *mongoSecret {
	expiresAt := &d.ExpiresAt
	if d.ExpiresAt.IsZero() {
		expiresAt = nil
	}

	return &mongoSecret{
		Hash:           d.Hash,
		SecretText:     d.SecretText,
		RemainingViews: d.RemainingViews,
		CreatedAt:      d.CreatedAt,
		ExpiresAt:      expiresAt,
	}
}

func NewMongoVault(clock secret.Clock, host string, port int, databaseName string, username string, password string) *mongoVault {
	connectionURL := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	clientOptions := options.Client().ApplyURI(connectionURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	collection := client.Database(databaseName).Collection(mongoCollectionName)

	return &mongoVault{
		collection: collection,
		clock:      clock,
	}
}

func (v *mongoVault) Store(secret *secret.Secret) error {
	_, err := v.collection.InsertOne(context.TODO(), mongoSecretFromDomainSecret(secret))
	if err != nil {
		return err
	}

	return nil
}

func (v *mongoVault) Retrieve(hash string) (*secret.Secret, error) {
	var result *mongoSecret

	filter := bson.D{{"_id", hash}}
	if err := v.collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, secret.SecretNotFoundError
		}
		return nil, err
	}

	return result.toDomainSecret(), nil
}

func (v *mongoVault) DecreaseRemainingViews(hash string) error {
	filter := bson.D{{"_id", hash}}
	update := bson.D{
		{"$inc", bson.D{
			{"remaining_views", -1},
		}},
	}

	updateResult, err := v.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return errors.New("Not found")
	}

	return nil
}
