package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	_defaultConnAttempts = 3
	_defaultConnTimeout  = 5 * time.Second
)

type DBConnString string

type DBName string

type mongoDB struct {
	connAttempts int
	connTimeout  time.Duration

	client   *mongo.Client
	database *mongo.Database
}

var _ DBEngine = (*mongoDB)(nil)

func NewMongoDB(url DBConnString, dbName DBName) (DBEngine, error) {
	zap.S().Info("Connection string:", string(url))

	mongoDB := &mongoDB{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	var err error
	for mongoDB.connAttempts > 0 {
		clientOptions := options.Client().ApplyURI(string(url))
		mongoDB.client, err = mongo.Connect(context.Background(), clientOptions)
		if err == nil {
			break
		}

		mongoDB.connAttempts--
		zap.S().Warnf("MongoDB connection failed, attempts left: %d", mongoDB.connAttempts)
		time.Sleep(mongoDB.connTimeout)
	}

	if err != nil {
		zap.S().Error("Failed to connect to MongoDB")
		return nil, err
	}

	mongoDB.database = mongoDB.client.Database(string(dbName))
	zap.S().Info("Connected to MongoDB 🎉")

	return mongoDB, nil
}

func (m *mongoDB) Configure(opts ...Option) DBEngine {
	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *mongoDB) GetDB() *mongo.Database {
	return m.database
}

func (m *mongoDB) GetClient() *mongo.Client {
	return m.client
}

func (m *mongoDB) Close() {
	if m.client != nil {
		err := m.client.Disconnect(context.Background())
		if err != nil {
			zap.S().Error("Error disconnecting from MongoDB", zap.Error(err))
		}
	}
}
