package dbi

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv(dbhost)
	if !ok {
		panic("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv(dbport)
	if !ok {
		panic("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv(dbuser)
	if !ok {
		panic("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv(dbpass)
	if !ok {
		panic("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv(dbname)
	if !ok {
		panic("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

// DBI  database interface
type DBI struct {
	db *gorm.DB
}

// Radiologdata DB model structure
type Radiologdata struct {
	gorm.Model
	Label        string
	Description  string
	Address      uint
	Timestamp    time.Time
	Lqi          int
	Rssi         int
	Uptime       int
	Tempcpu      int
	Vrefcpu      int
	Ntc0         int
	Ntc1         int
	Photores     int
	Pressure     int
	Temppressure int
}

// Temperature of gived address
func (dbp *DBI) Temperature(address uint) {
	var d []Radiologdata
	dbp.db.Where(&Radiologdata{Address: address}).Find(&d)

	for i := 0; i < 10; i++ {
		fmt.Println(d[i].Ntc0)
	}
}

// Init init db module interface
func (dbp *DBI) Init() error {
	config := dbConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	var err error
	dbp.db, err = gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Unable to connect to DB: ", err)
		return err
	}
	//defer dbp.db.Close()
	log.Info("Successfully connected!")

	// Migrate the schema
	dbp.db.AutoMigrate(&Radiologdata{})
	return nil
}
