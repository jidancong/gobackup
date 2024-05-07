package dbs

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getPGDBs(host, username, password, port string) ([]string, error) {
	gormdb, err := pgClient(host, username, password, port, "")
	if err != nil {
		return nil, err
	}

	var databases []struct {
		DatabaseName string `gorm:"column:datname"`
	}

	if err := gormdb.Raw("SELECT datname FROM pg_database WHERE datistemplate = false").Scan(&databases).Error; err != nil {
		panic(err)
	}

	var dbNames []string
	for _, v := range databases {
		dbNames = append(dbNames, v.DatabaseName)
	}

	return dbNames, nil
}

func pgClient(host, username, password, port, database string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=template1 port=%s sslmode=disable TimeZone=Asia/Shanghai", host, username, password, port)
	if len(database) != 0 {
		dsn = fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", host, username, password, port, database)
	}
	gormdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormdb, nil
}

// 数据库客户端
func mysqlClient(username, password, host, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

// 获取所有数据库
func getDBs(username, password, host, port string) ([]string, error) {
	db, err := mysqlClient(username, password, host, port)
	if err != nil {
		return nil, err
	}

	var databases []struct {
		DatabaseName string `gorm:"column:Database"`
	}

	if err := db.Raw("SHOW DATABASES").Scan(&databases).Error; err != nil {
		return nil, err
	}

	var dbNames []string
	for _, db := range databases {
		dbNames = append(dbNames, db.DatabaseName)
	}

	return dbNames, err
}

// 获取数据库版本
func version(username, password, host, port string) (string, error) {
	db, err := mysqlClient(username, password, host, port)
	if err != nil {
		return "", err
	}

	var version string
	err = db.Raw("SELECT VERSION()").Scan(&version).Error
	return version, err
}

func getMongoClient(username, password, host, port string) (*mongo.Client, error) {
	applyUrl := fmt.Sprintf("mongodb://%s:%s", host, port)
	if username != "" {
		applyUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(applyUrl))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getMongoDBs(username, password, host, port string) ([]string, error) {
	client, err := getMongoClient(username, password, host, port)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())
	names, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	return names, err
}
