package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type pack struct {
	id      int
	dexid   string
	name    string
	pokemon []int
}

type pokemon struct {
	ID       int
	Dexid    string `json:"dexid"`
	Name     string `json:"name"`
	Health   int    `json:"health"`
	PType    string `json:"ptype"`
	Stage    int    `json:"stage"`
	Weakness string `json:"weakness"`
	Ablities []int  `json:"ablities"`
	Attacks  []int  `json:"attacks"`
	Rarity   string `json:"rarity"`
	Retreat  int    `json:"retreat"`
	Packs    []int  `json:"packs"`
}

type ablity struct {
	id          int
	name        string
	description string
}

//
// cost: "Grass", "Colorless", "Fire", "etc"
//

type attack struct {
	id          int
	name        string
	damage      int
	cost        []string
	description string
}

type user struct {
	id      int
	name    string
	duser   string
	pokemon []int
}

func dbCreate(db *sql.DB) {
	// create all of the tables that are in the database
	dbcall, err := db.Prepare("DROP TABLE IF EXISTS Pokemon")
	if err != nil {
		log.Print(err)
	}
	dbcall.Exec()
	// create the Pokemon table
	// this database contains all the pokemon in the game and their infomation
	dbcall, err = db.Prepare("CREATE TABLE IF NOT EXISTS Pokemon (id INTEGER PRIMARY KEY, dexid string, name STRING, health INTEGER, ptype STRING, stage INTEGER, weakness STRING, ablities STRING, attacks STRING, rarity STRING, retreat INTEGER, packs STRING) ")
	if err != nil {
		log.Print(err)
	}
	dbcall.Exec()
}

func importCardData(filename string) []pokemon {
	var cards []pokemon

	cardFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(cardFile)
	if err = jsonParser.Decode(&cards); err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	return cards

}

// func userPrint(db *sql.DB) {
// 	var data user
// 	dbcall, err := db.Query("SELECT * FROM User")
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	for dbcall.Next() {
// 		dbcall.Scan(&data.id, &data.name, &data.duser, &data.pokemon)
// 		log.Printf("ID:%d, Username:%s\n", data.id, data.name)
// 	}
// }

func locateDatabase() string {
	dbLocle := os.Getenv("SystemDrive") + "\\ProgramData\\.ptcgp"
	_, err := os.Stat(dbLocle)

	if os.IsNotExist(err) {
		err = os.Mkdir(dbLocle, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	return dbLocle
}

func main() {
	// access cli flags
	createFlag := flag.Bool("create", false, "This creates the database")

	// edit help message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of PTCGP Database Tool:\n")
		flag.PrintDefaults()
	}

	// parse flags
	flag.Parse()

	dbLocle := locateDatabase()
	fmt.Println()
	db, err := sql.Open("sqlite3", dbLocle+"\\ptcgpv1.db")

	if err != nil {
		log.Fatal(err)
	}

	// create flag
	if *createFlag {
		dbCreate(db)
	}

	a := importCardData(".\\cards.json")

	for _, card := range a {
		log.Println(card.Name)
	}

	defer db.Close()

}
