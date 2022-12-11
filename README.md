[![codecov](https://codecov.io/gh/DAtek/gopaq/branch/main/graph/badge.svg?token=WWY2L6G56Y)](https://codecov.io/gh/DAtek/gopaq)

# GOPAQ - Gorm Paginated Query

Limit-offset pagination for [GORM](https://gorm.io/) queries.


## Example
```go
package main

import (
	"fmt"
	"os"

	"github.com/DAtek/gopaq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)


type Plant struct {
	Id   uint
	Type string
	Name string
}

// You have to implement this function
func getSession() *gorm.DB {
    return nil
}

func main() {
	session := getSession()

	for i := 0; i < 10; i++ {
		session.Create(&Plant{})
	}

	plant1 := &Plant{Name: "Banana"}
	session.Create(plant1)

	plant2 := &Plant{Name: "Apple"}
	session.Create(plant2)

	query := session.Model(&Plant{}).Order("name DESC")
	result, _ := gopaq.FindWithPagination(query, []*Plant{}, 1, 2)

	fmt.Printf("plant1: %v\n", result.Items[0])
	fmt.Printf("plant2: %v\n", result.Items[1])
	fmt.Printf("returned items: %v\n", len(result.Items))
	fmt.Printf("total: %v\n", result.Total)
}

```
Output:
```
plant1: &{11  Banana}
plant2: &{12  Apple}
returned items: 2
total: 12
```