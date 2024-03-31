package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type DBEngine interface {
	GetDB() *mongo.Database
	Configure(...Option) DBEngine
	Close()
}
