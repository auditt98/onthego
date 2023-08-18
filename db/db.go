package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/auditt98/onthego/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Init ...
func Init() {

	db_instance, err := ResolveDB()
	if err == nil {
		fmt.Println("Running migrations...")
		db_instance.AutoMigrate(&models.Article{})
		db_instance.AutoMigrate(&models.User{})
	}
}

func ResolveDB() (*gorm.DB, error) {
	var dbInstance *gorm.DB
	var err error
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")) //user, pass, host, port, dbname
		dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case "postgres":
		dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "mssql":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		dbInstance, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("DB_DRIVER not found")
	}
	// check if db exists
	err = dbInstance.Exec(fmt.Sprintf("USE %s;", os.Getenv("DB_NAME"))).Error
	if err != nil && strings.Contains((err.Error()), "ERROR: database \""+os.Getenv("DB_NAME")+"\" does not exist") {
		fmt.Println("Database ", os.Getenv("DB_NAME"), " not found. Creating...")
		err = dbInstance.Exec("CREATE DATABASE " + os.Getenv("DB_NAME")).Error
	}
	return dbInstance, err
}

// Pipeline for query: Population -> Filter -> Sort -> Paging -> Projection
func Query(tableName string, params QueryParams, result interface{}) *gorm.DB {
	queryEngine := new(QueryEngine)
	instance, _ := ResolveDB()
	queryEngine.TableName = tableName
	queryEngine.Ref = instance.Table(tableName)
	queryEngine.Ref = queryEngine.Populate(params).Filter(params).Sort(params).Paginate(params).Projection(params).ToGorm().Find(result)
	return queryEngine.Ref
}

func (qe *QueryEngine) ToGorm() *gorm.DB {
	return qe.Ref
}
