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

// var queryEngine *QueryEngine

// Init ...
func Init() {

	db_instance, err := ResolveDB()
	if err != nil {
		log.Fatal(err)
	}
	db_instance.AutoMigrate(&models.Article{})
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

// Pipeline for query: Filter -> Sort -> Paging -> Population -> Projection
func Query(tableName string, params QueryParams, result interface{}) *gorm.DB {
	queryEngine := new(QueryEngine)
	instance, _ := ResolveDB()
	queryEngine.Ref = instance.Table(tableName)
	//declare new object with type of modelType
	queryEngine.Ref = queryEngine.Filter(params).Find(result)
	return queryEngine.Ref
}

func (qe *QueryEngine) HandleAnd(params []WhereParams) *gorm.DB {
	if (len(params)) == 0 {
		return qe.Ref
	}
	for _, param := range params {
		qe.Ref = qe.Ref.Where(qe.HandleAnd(param.And)).Where(qe.HandleOr(param.Or)).Where(qe.HandleAndAttr(param.Attr))
	}
	return qe.Ref
}

func (qe *QueryEngine) HandleOr(params []WhereParams) *gorm.DB {
	if (len(params)) == 0 {
		return qe.Ref
	}
	for _, param := range params {
		qe.Ref = qe.Ref.Or(qe.HandleAnd(param.And)).Or(qe.HandleOr(param.Or)).Or(qe.HandleOrAttr(param.Attr))
	}
	return qe.Ref
}

func (qe *QueryEngine) HandleAndAttr(params map[string]AttributeParams) *gorm.DB {
	if (len(params)) == 0 {
		return qe.Ref
	}
	for key, param := range params {
		if param.Eq != nil {
			qe.Ref = qe.Ref.Where(key, param.Eq)
		}
		if param.Eqi != nil {
			qe.Ref = qe.Ref.Where(key+" = ?", param.Eqi)
		}
		if param.Ne != nil {
			qe.Ref = qe.Ref.Where(key+" <> ?", param.Ne)
		}
		if param.In != nil {
			qe.Ref = qe.Ref.Where(key+" IN ?", param.In)
		}
		if param.Nin != nil {
			qe.Ref = qe.Ref.Where(key+" NOT IN ?", param.Nin)
		}
		if param.Lt != "" {
			qe.Ref = qe.Ref.Where(key+" < ?", param.Lt)
		}
		if param.Lte != "" {
			qe.Ref = qe.Ref.Where(key+" <= ?", param.Lte)
		}
		if param.Gt != "" {
			qe.Ref = qe.Ref.Where(key+" > ?", param.Gt)
		}
		if param.Gte != "" {
			qe.Ref = qe.Ref.Where(key+" >= ?", param.Gte)
		}
		if param.Between != nil {
			qe.Ref = qe.Ref.Where(key+" BETWEEN ? AND ?", param.Between[0], param.Between[1])
		}
		if param.Contains != "" {
			qe.Ref = qe.Ref.Where(key+" LIKE ?", "%"+param.Contains+"%")
		}
		if param.NContains != "" {
			qe.Ref = qe.Ref.Where(key+" NOT LIKE ?", "%"+param.NContains+"%")
		}
		if param.Containsi != "" {
			qe.Ref = qe.Ref.Where(key+" ILIKE ?", "%"+param.Containsi+"%")
		}
		if param.NContainsi != "" {
			qe.Ref = qe.Ref.Where(key+" NOT ILIKE ?", "%"+param.NContainsi+"%")
		}
		if param.StartsWith != "" {
			qe.Ref = qe.Ref.Where(key+" LIKE ?", param.StartsWith+"%")
		}
		if param.EndsWith != "" {
			qe.Ref = qe.Ref.Where(key+" LIKE ?", "%"+param.EndsWith)
		}
		if param.NStartsWith != "" {
			qe.Ref = qe.Ref.Where(key+" NOT LIKE ?", param.NStartsWith+"%")
		}
		if param.Nil {
			qe.Ref = qe.Ref.Where(key + " IS NULL")
		}
		if param.NNil {
			qe.Ref = qe.Ref.Where(key + " IS NOT NULL")
		}
	}
	return qe.Ref
}

func (qe *QueryEngine) HandleOrAttr(params map[string]AttributeParams) *gorm.DB {

	if (len(params)) == 0 {
		return qe.Ref
	}
	for key, param := range params {
		if param.Eq != nil {
			qe.Ref = qe.Ref.Or(key, param.Eq)
		}
		if param.Eqi != nil {
			qe.Ref = qe.Ref.Or(key+" = ?", param.Eqi)
		}
		if param.Ne != nil {
			qe.Ref = qe.Ref.Or(key+" <> ?", param.Ne)
		}
		if param.In != nil {
			qe.Ref = qe.Ref.Or(key+" IN ?", param.In)
		}
		if param.Nin != nil {
			qe.Ref = qe.Ref.Or(key+" NOT IN ?", param.Nin)
		}
		if param.Lt != "" {
			qe.Ref = qe.Ref.Or(key+" < ?", param.Lt)
		}
		if param.Lte != "" {
			qe.Ref = qe.Ref.Or(key+" <= ?", param.Lte)
		}
		if param.Gt != "" {
			qe.Ref = qe.Ref.Or(key+" > ?", param.Gt)
		}
		if param.Gte != "" {
			qe.Ref = qe.Ref.Or(key+" >= ?", param.Gte)
		}
		if param.Between != nil {
			qe.Ref = qe.Ref.Or(key+" BETWEEN ? AND ?", param.Between[0], param.Between[1])
		}
		if param.Contains != "" {
			qe.Ref = qe.Ref.Or(key+" LIKE ?", "%"+param.Contains+"%")
		}
		if param.NContains != "" {
			qe.Ref = qe.Ref.Or(key+" NOT LIKE ?", "%"+param.NContains+"%")
		}
		if param.Containsi != "" {
			qe.Ref = qe.Ref.Or(key+" ILIKE ?", "%"+param.Containsi+"%")
		}
		if param.NContainsi != "" {
			qe.Ref = qe.Ref.Or(key+" NOT ILIKE ?", "%"+param.NContainsi+"%")
		}
		if param.StartsWith != "" {
			qe.Ref = qe.Ref.Or(key+" LIKE ?", param.StartsWith+"%")
		}
		if param.EndsWith != "" {
			qe.Ref = qe.Ref.Or(key+" LIKE ?", "%"+param.EndsWith)
		}
		if param.NStartsWith != "" {
			qe.Ref = qe.Ref.Or(key+" NOT LIKE ?", param.NStartsWith+"%")
		}
		if param.Nil {
			qe.Ref = qe.Ref.Or(key + " IS NULL")
		}
		if param.NNil {
			qe.Ref = qe.Ref.Or(key + " IS NOT NULL")
		}
	}
	return qe.Ref
}

func (qe *QueryEngine) Filter(params QueryParams) *gorm.DB {
	where := params.Where
	qe.Ref = qe.Ref.Where(qe.HandleAnd(where.And)).Where(qe.HandleOr(where.Or)).Where(qe.HandleAndAttr(where.Attr))
	return qe.Ref
}

// func (qe *QueryEngine) Create(modelName string, data interface{}) (interface{}, error) {
// 	return nil, nil
// }

// func (qe *QueryEngine) Update(modelName string, id uint, data interface{}) (interface{}, error) {
// 	return nil, nil
// }

// func (qe *QueryEngine) Delete(modelName string, id uint) (interface{}, error) {
// 	return nil, nil
// }

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
