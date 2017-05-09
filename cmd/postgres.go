package cmd

import (
	"database/sql"
	"fmt"
	"html/template"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sql.DB
var tpl *template.Template

func init() {
	var (
		err      error
		host     string
		database string
		username string
		password string
	)

	v := viper.New()

	v.AutomaticEnv()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.go-ops-controller/")
	v.SetEnvPrefix("ops")
	err = v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			fmt.Println("Parser error")
		}
		fmt.Println("No configuration file found")

	}

	host = v.GetString("host")
	database = v.GetString("dbname")
	username = v.GetString("username")
	password = v.GetString("password")
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", username, password, database, host)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to the database.")
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
