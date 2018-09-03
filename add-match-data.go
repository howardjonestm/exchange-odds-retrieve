//MOVE TO SDK ONCE COMPLETE
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
)

//ENV ENDTIME
//ENV MATCHID
//ENV TEAM

func main() {
	createTable()

	endTime, err := strconv.Atoi(os.Getenv("ENDTIME"))
	check(err)

	for endTime < int(time.Now().Unix()) {

		parseAddData(matchId)
		time.Sleep(10 * time.Second)
	}
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

var matchId string = os.Getenv("MATCHID")
var tableName string = os.Getenv("TEAM")

func createTable() error {

	db, err := sql.Open("mysql", "howardjones94:password@/football_price_tracker")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	createStatement := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
		time int not null, 
		home float null,
		home_liquidity float null,
		draw float null,
		draw_liquidity float null,
		away float null,
		away_liquidity float null
	)ENGINE=innodb;`, tableName)

	fmt.Println(createStatement)
	stmtIn, err := db.Prepare(createStatement)

	check(err)
	defer stmtIn.Close()

	_, err = stmtIn.Exec()
	check(err)
	return err
}

func addData(time int64, home, homeLiquidity, draw, drawLiquidity, away, awayLiquidity float64) error {
	db, err := sql.Open("mysql", "root:root@tcp(database:3306)/football_price_tracker")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	statement := fmt.Sprintf("INSERT INTO %s VALUES( ?, ?, ?, ?, ?, ?, ?, ?  )", tableName)
	stmtIn, err := db.Prepare(statement)
	check(err)
	defer stmtIn.Close()

	_, err = stmtIn.Exec(nil, time, home, homeLiquidity, draw, drawLiquidity, away, awayLiquidity)
	check(err)
	return err

}

func retrieveEventData(eventId string) string {
	requestUrl := "https://api.matchbook.com/edge/rest/events/" + eventId

	req, err := http.NewRequest("GET", requestUrl, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)

	fmt.Println("response status: ", resp.Status)
	fmt.Println("response headers: ", resp.Header)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	stringBody := string(body)
	return stringBody

}

func parseAddData(eventId string) error {
	content := retrieveEventData(eventId)

	dataJson := string(content)
	now := time.Now()
	secs := now.Unix()

	homeOdds := (gjson.Get(dataJson, `markets.0.runners.0.prices.0.odds`)).Float()
	homeLiquidity := (gjson.Get(dataJson, `markets.0.runners.0.prices.0.available-amount`)).Float()

	drawOdds := (gjson.Get(dataJson, `markets.0.runners.2.prices.0.odds`)).Float()
	drawLiquidity := (gjson.Get(dataJson, `markets.0.runners.2.prices.0.available-amount`)).Float()

	awayOdds := (gjson.Get(dataJson, `markets.0.runners.1.prices.0.odds`)).Float()
	awayLiquidity := (gjson.Get(dataJson, `markets.0.runners.1.prices.0.available-amount`)).Float()

	err := addData(secs, homeOdds, homeLiquidity, drawOdds, drawLiquidity, awayOdds, awayLiquidity)
	return err
}
