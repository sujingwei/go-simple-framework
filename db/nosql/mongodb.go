/*
 * @Author: sujingwei 348149047@qq.com
 * @Date: 2024-05-19 16:04:58
 * @LastEditors: sujingwei 348149047@qq.com
 * @LastEditTime: 2024-05-20 15:08:53
 * @FilePath: \go-simple-framework\db\nosql\mongodb.go
 * @Description: mongodb 连接信息
 */
package nosql

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 单个mongodb的连接信息
type MongoConnectionConfig struct {
	Key         string `json:"key" yaml:"key"`                 // 连接的名称，默认为default
	Uri         string `json:"uri" yaml:"uri"`                 // 如：localhost:27017
	Username    string `json:"username" yaml:"username"`       // 用户名
	Password    string `json:"password" yaml:"password"`       // 密码
	AuthSource  string `json:"authSource" yaml:"authSource"`   // 连接的数据库
	MaxPoolSize uint64 `json:"maxPoolSize" yaml:"maxPoolSize"` // 连接池最大连接数
	MinPoolSize uint64 `json:"minPoolSize" yaml:"minPoolSize"` // 连接池最小连接数
}

// Mongodb连接池
type MongoDbConfig struct {
	Pool    []MongoConnectionConfig `json:"pool" yaml:"pool"`       // 连接池，连接池中默认没有出现default，就使用Default的配置，如果出现
	Default MongoConnectionConfig   `json:"default" yaml:"default"` // 单独配置mongodb
}

var (
	// MongoDB连接池, key => *mongo.Client
	mongodbPool              map[string]*mongo.Client = make(map[string]*mongo.Client) // 多个连接的连接池
	mongodbKeyToDatabaseName map[string]string        = make(map[string]string)        // key对应的数据库名称
	DefaultName              string                   = "default"                      // 默认连接客户端的名称
)

// 连接操作，返回 client
func (m *MongoConnectionConfig) connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if m.MaxPoolSize <= 0 {
		m.MaxPoolSize = 5
	}
	if m.MinPoolSize <= 0 {
		m.MinPoolSize = 1
	}
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username:   m.Username,
		Password:   m.Password,
		AuthSource: m.AuthSource,
	}).ApplyURI(m.Uri).SetMaxPoolSize(m.MaxPoolSize).SetMinPoolSize(m.MinPoolSize))
	return client, err
}

// 创建Mongodb连接池的集合
func createMongodbPool(mpc *MongoDbConfig) {
	var j int
	for i := 0; i < len(mpc.Pool); i++ {
		if mpc.Pool[i].Key == "default" || mpc.Pool[i].Key == "" {
			j++
		}
	}
	if j > 1 {
		panic("The same value is displayed in the MongoDB connection pool!")
	}

	for i := 0; i < len(mpc.Pool); i++ {
		var mcc MongoConnectionConfig = mpc.Pool[i]
		if mcc.AuthSource == "" {
			panic("mongodb AuthSource is Empty!")
		}
		if mcc.Uri == "" {
			panic("mongodb Uri is Empty!")
		}
		if c, err := mcc.connect(); err != nil {
			log.Fatal(err)
		} else {
			key := DefaultName // 如果为空
			if mcc.Key != "" {
				key = mcc.Key
			}
			// 保存到池中
			mongodbPool[key] = c
			// 对应的库名称
			mongodbKeyToDatabaseName[key] = mcc.AuthSource
			log.Printf("create mongodb pool, database: %s, minPoolSize: %d, maxPoolSize: %d\n", c.Database(mcc.AuthSource).Name(), mcc.MaxPoolSize, mcc.MinPoolSize)
		}
	}
}

// 项目中使用mongodb
func UseMongoDB(mongoDbConfig *MongoDbConfig) {
	if mongoDbConfig != nil {
		if mongoDbConfig.Default != nil {
			mongoDbConfig.Pool = append(mongoDbConfig.Pool, *mongoDbConfig.Default)
		}
	}
	if len(mongoDbConfig.Pool) > 0 {
		createMongodbPool(mongoDbConfig)
	}
}

// 获取Mongodb指定key的Client
func GetMongoDbClient(key string) *mongo.Client {
	return mongodbPool[key]
}

// 获取默认的Mongodb的Client
func GetMongoDbDefaultClient() *mongo.Client {
	return mongodbPool[DefaultName]
}

// 获取MongoDB的Database
func GetMongoDbDatabase(key string) *mongo.Database {
	return GetMongoDbClient(key).Database(mongodbKeyToDatabaseName[key])
}

// 获取MongoDB默认的Database
func GetMongoDbDefaultDatabase() *mongo.Database {
	return GetMongoDbDatabase(DefaultName)
}

// 获取key下的集合
func GetMongodbCollection(key string, collectionName string) *mongo.Collection {
	return GetMongoDbDatabase(key).Collection(collectionName)
}

// 获取默认key下的集合
func GetMongodbDefaultKeyCollection(collectionName string) *mongo.Collection {
	return GetMongoDbDefaultDatabase().Collection(collectionName)
}
