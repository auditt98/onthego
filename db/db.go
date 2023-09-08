package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/auditt98/onthego/models"
	deepgorm "github.com/survivorbat/gorm-deep-filtering"
	gormqonvert "github.com/survivorbat/gorm-query-convert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init ...
func Init() {

	db_instance, err := ResolveDB()
	db_instance.Use(deepgorm.New())
	config := gormqonvert.CharacterConfig{
		GreaterThanPrefix:      ">",
		GreaterOrEqualToPrefix: ">=",
		LessThanPrefix:         "<",
		LessOrEqualToPrefix:    "<=",
		NotEqualToPrefix:       "!=",
		LikePrefix:             "~",
		NotLikePrefix:          "!~",
	}
	db_instance.Use(gormqonvert.New(config))
	DB = db_instance
	if err == nil {
		fmt.Println("Running migrations...")
		db_instance.AutoMigrate(&models.User{}, &models.Album{}, &models.Photo{}, &models.Comment{}, &models.Like{})

	}
}

func ResolveDB() (*gorm.DB, error) {
	var dbInstance *gorm.DB
	var err error
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")) //user, pass, host, port, dbname
		dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})

	case "postgres":
		dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
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

func Query(searchParams *SearchParams, additionalFilter map[string]any, model interface{}, count *int64) *gorm.DB {
	dbLoader := DB.Debug().Where(searchParams.Filters)
	if additionalFilter != nil {
		dbLoader = dbLoader.Where(additionalFilter)
	}
	for _, populate := range searchParams.Populate {
		if populate == "*" {
			dbLoader = dbLoader.Preload(clause.Associations)
			continue
		}
		dbLoader = dbLoader.Preload(populate)
	}
	for _, sort := range searchParams.Sort {
		parsedSort := sort
		parsedSort = strings.Replace(parsedSort, "-", " desc", 1)
		parsedSort = strings.Replace(parsedSort, "+", "", 1)
		dbLoader = dbLoader.Order(parsedSort)
	}
	page := searchParams.Page
	perPage := searchParams.PerPage
	dbLoader.Scopes(Paginate(&page, &perPage)).Find(model).Count(count)
	searchParams.Page = page
	searchParams.PerPage = perPage
	return dbLoader
}

func QueryOne(searchParams *SearchParams, additionalFilter map[string]any, model interface{}) *gorm.DB {
	dbLoader := DB.Debug().Where(searchParams.Filters)
	if additionalFilter != nil {
		dbLoader = dbLoader.Where(additionalFilter)
	}
	for _, populate := range searchParams.Populate {
		if populate == "*" {
			dbLoader = dbLoader.Preload(clause.Associations)
			continue
		}
		dbLoader = dbLoader.Preload(populate)
	}
	for _, sort := range searchParams.Sort {
		parsedSort := sort
		parsedSort = strings.Replace(parsedSort, "-", " desc", 1)
		parsedSort = strings.Replace(parsedSort, "+", "", 1)
		dbLoader = dbLoader.Order(parsedSort)
	}
	dbLoader.First(model)
	return dbLoader
}
