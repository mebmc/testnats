package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/go-nats-streaming/pb"
)

type Entry struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Guess     string `json:"guess"`
}

type msgHandler struct {
	db *sql.DB
}

func main() {
	var clusterID = "test-cluster"
	var clientID = "ingestor01"
	var URL = "nats://docker.local:4222"
	var qgroup = "quiz"
	var durable = "quiz"
	var unsubscribe bool
	var subj = "1"

	db, err := sql.Open("mysql", "root:@tcp(docker.local:3306)/quiz?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	mh := msgHandler{db: db}

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", URL, clusterID, clientID)

	startOpt := stan.StartAt(pb.StartPosition_NewOnly)

	sub, err := sc.QueueSubscribe(subj, qgroup, mh.receiveMsg, startOpt, stan.DurableName(durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", subj, clientID, qgroup, durable)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if durable == "" || unsubscribe {
				sub.Unsubscribe()
			}
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func (z msgHandler) receiveMsg(msg *stan.Msg) {
	e := Entry{}
	json.Unmarshal(msg.Data, &e)
	log.Printf("Entry! FName: %v, LName: %v, Email: %v, Guess: %v", e.FirstName, e.LastName, e.Email, e.Guess)
	z.insertEntry(e)
}

func (z msgHandler) insertEntry(e Entry) {
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	query := fmt.Sprintf("INSERT INTO quiz_1 VALUES ( 0, '%v', '%v', '%v', '%v','%v', '%v' )", e.ID, e.FirstName, e.LastName, e.Email, e.Guess, ts)
	insert, err := z.db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	insert.Close()
}
