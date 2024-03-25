package bin

import (
	"bufio"
	_ "database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"os"
)

var DB *sqlx.DB
var data []string
var Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func ConnectDB() {
	file, err := os.Open("config.txt")
	CheckErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
		CheckErr(scanner.Err())
	}
	info := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", data[0], data[1], data[2])
	DB, err = sqlx.Connect("postgres", info)
	CheckErr(err)
}
