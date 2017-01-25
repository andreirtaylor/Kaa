package kaa

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/appengine"
	//"google.golang.org/appengine/log"
	"log"
	"net/http"
	"os"
)

type JSON map[string]string

func handler(w http.ResponseWriter, r *http.Request) {

}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}

func respond(res http.ResponseWriter, obj JSON) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(obj)
}

func handleInfo(w http.ResponseWriter, req *http.Request) {
	connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user := mustGetenv("CLOUDSQL_USER")
	password := mustGetenv("CLOUDSQL_PASSWORD")

	w.Header().Set("Content-Type", "text/plain")

	sqlStr := fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName)

	if appengine.IsDevAppServer() {
		sqlStr = "root@/" // dev server has no password baby
	}

	db, err := sql.Open(
		"mysql",
		sqlStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open db: %v", err), 500)
		return
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not query db: %v", err), 500)
		return
	}
	defer rows.Close()

	buf := bytes.NewBufferString("Databases:\n")
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			http.Error(w, fmt.Sprintf("Could not scan result: %v", err), 500)
			return
		}
		fmt.Fprintf(buf, "- %s\n", dbName)
	}
	w.Write(buf.Bytes())

	//respond(res, JSON{
	//	"color": "#ff0000",
	//	"head":  "https://golang.org/doc/gopher/gopherbw.png",
	//})
}

func handleStart(res http.ResponseWriter, req *http.Request) {
	respond(res, JSON{
		"taunt": "battlesnake-go!",
	})
}

func handleMove(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	//client, err := storage.NewClient(ctx)

	//if err != nil {
	//	log.Errorf(ctx, "Cannot initialize data storage")
	//}
	ctx.Done()
	data, err := NewMoveRequest(req)
	if err != nil {
		respond(res, JSON{"move": "north", "taunt": "can't parse this!"})
		return
	}

	snake := data.GetSnake(SNAKE_ID)
	if snake == nil {
		respond(res, JSON{"move": "north", "taunt": "snake not found!"})
		return
	}

	board := NewBoard(data)
	location := snake.Head()

	var chosenDirection string

	if tile := board.GetTile(location.X, location.Y-1); tile == EMPTY {
		chosenDirection = "north"
	} else if tile := board.GetTile(location.X, location.Y+1); tile == EMPTY {
		chosenDirection = "south"
	} else if tile := board.GetTile(location.X-1, location.Y); tile == EMPTY {
		chosenDirection = "west"
	} else if tile := board.GetTile(location.X+1, location.Y); tile == EMPTY {
		chosenDirection = "east"
	} else {
		chosenDirection = "north"
	}

	respond(res, JSON{
		"move":  chosenDirection,
		"taunt": "go go go!",
	})
}

func handleEnd(res http.ResponseWriter, req *http.Request) {
	respond(res, JSON{})
}
