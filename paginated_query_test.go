package gopaq

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestFindWithPagination(t *testing.T) {
	database := newDatabase()
	database.migrateUp()
	defer database.migrateDown()
	getDb := database.getSessionFactory()

	t.Run("Page size is correct", func(t *testing.T) {
		db := getDb()
		defer db.Rollback()
		for i := 0; i < 3; i++ {
			panicResult(db.Create(&Plant{}))
		}

		query := db.Model(&Plant{})

		result, err := FindWithPagination(query, []*Plant{}, 1, 2)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(result.Items))

	})

	t.Run("Total is correct", func(t *testing.T) {
		db := getDb()
		defer db.Rollback()
		total := 5
		for i := 0; i < total; i++ {
			panicResult(db.Create(&Plant{}))
		}

		query := db.Model(&Plant{})

		result, err := FindWithPagination(query, []*Plant{}, 1, 1)

		assert.Nil(t, err)
		assert.Equal(t, uint(total), result.Total)

	})

	t.Run("Page is correct", func(t *testing.T) {
		db := getDb()
		defer db.Rollback()

		for i := 0; i < 10; i++ {
			panicResult(db.Create(&Plant{}))
		}

		plant1 := &Plant{Name: "Banana"}
		panicResult(db.Create(plant1))

		plant2 := &Plant{Name: "Apple"}
		panicResult(db.Create(plant2))

		query := db.Model(&Plant{}).Order("name DESC")

		result, err := FindWithPagination(query, []*Plant{}, 2, 2)

		assert.Nil(t, err)
		assert.Equal(t, "", result.Items[0].Name)
		assert.Equal(t, "", result.Items[1].Name)
	})

	t.Run("Accepts 0 for page and page size", func(t *testing.T) {
		db := getDb()
		defer db.Rollback()
		for i := 0; i < 100; i++ {
			panicResult(db.Create(&Plant{}))
		}

		query := db.Model(&Plant{})

		result, err := FindWithPagination(query, []*Plant{}, 0, 0)

		assert.Nil(t, err)
		assert.Equal(t, DefaultLimit, len(result.Items))
	})

	t.Run("Returns error if something unexpected happens", func(t *testing.T) {
		db := getDb()

		// table doesn't exists
		query := db.Model(&Animal{})

		_, err := FindWithPagination(query, []*Plant{}, 0, 0)

		assert.Error(t, err)
	})
}

type Plant struct {
	Id   uint
	Type string
	Name string
}

type Animal struct {
	Id   uint
	Type string
	Name string
}

type database struct {
	db *gorm.DB
}

func newDatabase() *database {
	return &database{
		db: createDb(),
	}
}

func (d *database) migrateUp() {
	panicResult(d.db.Exec(`
CREATE TABLE plant (
	id SERIAL PRIMARY KEY NOT NULL,
	type VARCHAR NOT NULL,
	name VARCHAR NOT NULL
);`))
}

func (d *database) migrateDown() {
	panicResult(d.db.Exec("DROP TABLE plant;"))
}

func (d *database) getSessionFactory() func() *gorm.DB {
	return func() *gorm.DB {
		return d.db.Begin()
	}
}

func createDb() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		"127.0.0.1",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	var db *gorm.DB = nil
	var dbErr error = nil

	for i := 0; i < 3; i++ {
		db, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			SkipDefaultTransaction: true,
		})

		if dbErr == nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}

	if dbErr != nil {
		panic(dbErr)
	}

	return db
}

func panicResult(result *gorm.DB) {
	if result.Error != nil {
		panic(result.Error)
	}
}
