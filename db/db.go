package db

import (
	"fmt"
	"log"
	"os"

	"github.com/auditt98/onthego/models"
	_redis "github.com/go-redis/redis/v7"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// DB ...

var db *gorm.DB

// var queryEngine *QueryEngine

// Init ...
func Init() {

	db_instance, err := ResolveDB()
	if err != nil {
		log.Fatal(err)
	}
	db = db_instance
	db.AutoMigrate(&models.Article{})
}

func ResolveDB() (*gorm.DB, error) {

	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")) //user, pass, host, port, dbname
		db_instance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		return db_instance, err
	case "postgres":
		dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		db_instance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		return db_instance, err
	case "mssql":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		db_instance, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		return db_instance, err
	default:
		return nil, fmt.Errorf("DB_DRIVER not found")
	}
}

func query(modelName string) *QueryEngine {
	queryEngine := new(QueryEngine)
	modelType, ok := models.ModelMap[modelName]
	if !ok {
		log.Fatal("Model not found")
	}
	queryEngine.Ref = db.Model(modelType)
	return queryEngine
}

// Pipeline for query: Filter -> Sort -> Paging -> Population -> Projection

func (qe *QueryEngine) FindOne(params QueryParams) (interface{}, error) {

	return nil, nil
}

func (qe *QueryEngine) FindMany(query interface{}, args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (qe *QueryEngine) FindWithCount(query interface{}, args ...interface{}) ([]interface{}, int64, error) {
	return nil, 0, nil
}

func (qe *QueryEngine) Create(modelName string, data interface{}) (interface{}, error) {
	return nil, nil
}

func (qe *QueryEngine) Update(modelName string, id uint, data interface{}) (interface{}, error) {
	return nil, nil
}

func (qe *QueryEngine) Delete(modelName string, id uint) (interface{}, error) {
	return nil, nil
}

func (qe *QueryEngine) Filter() *gorm.DB {
	// SELECT * FROM ... where (everything inside and) AND (everything inside or) AND (everything inside not)
	//

	// Handle WhereParams
	// type WhereParams struct {
	// 	And  []WhereParams
	// 	Or   []WhereParams
	// 	Not  []interface{}
	// 	Attr map[string]AttributeParams
	// }

	// var a = QueryParams{
	// 	Where: WhereParams{
	// 		Or: []WhereParams{
	// 			{
	// 				Attr: map[string]AttributeParams{
	// 					"email": {
	// 						Eq: "test",
	// 					},
	// 				},
	// 			},
	// 			{
	// 				Attr: map[string]AttributeParams{
	// 					"email": {
	// 						Eq: "test",
	// 					},
	// 				},
	// 			},
	// 		},
	// 		Attr: map[string]AttributeParams{
	// 			"email": {
	// 				Eq: "test",
	// 			},
	// 		},
	// 	},
	// }

	// applies implicit AND to all attributes

	return db
}

func (qe *QueryEngine) Sort() *gorm.DB {
	return db
}

func (qe *QueryEngine) Paging() *gorm.DB {
	return db
}

func (qe *QueryEngine) Population() *gorm.DB {
	return db
}

func (qe *QueryEngine) Projection() *gorm.DB {
	// db.whe
	return db
}

// RedisClient ...
var RedisClient *_redis.Client

// InitRedis ...
func InitRedis(selectDB ...int) {

	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")
	RedisClient = _redis.NewClient(&_redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       selectDB[0],
		// DialTimeout:        10 * time.Second,
		// ReadTimeout:        30 * time.Second,
		// WriteTimeout:       30 * time.Second,
		// PoolSize:           10,
		// PoolTimeout:        30 * time.Second,
		// IdleTimeout:        500 * time.Millisecond,
		// IdleCheckFrequency: 500 * time.Millisecond,
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },
	})

}

// GetRedis ...
func GetRedis() *_redis.Client {
	return RedisClient
}
