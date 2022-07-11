package main

//https://www.digitalocean.com/community/tutorials/debugging-go-code-with-visual-studio-code
import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type DB struct {
	firestore *firestore.Client
	app       *firebase.App
}

func (db *DB) New() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("jackworkout.json")
	var err error
	db.app, err = firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	db.firestore, err = db.app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func (db *DB) Delete() {
	if db.app != nil {
		//db.firestore
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Web!\n"))
}

func (db *DB) DoneChange(c *gin.Context) {
	id := c.Param("id")
	done := c.Param("done")
	db.firestore.Collection("Sets").Doc(id).Update(context.Background(),
		[]firestore.Update{
			{
				Path:  "done",
				Value: done,
			},
		})

}
func (db *DB) RepsChange(c *gin.Context) {
	id := c.Param("id")
	reps := c.Param("reps")
	db.firestore.Collection("Sets").Doc(id).Update(context.Background(),
		[]firestore.Update{
			{
				Path:  "reps_or_duration",
				Value: reps,
			},
		})
}
func (db *DB) WeightChanged(c *gin.Context) {
	id := c.Param("id")
	weight := c.Param("weight")
	db.firestore.Collection("Sets").Doc(id).Update(context.Background(),
		[]firestore.Update{
			{
				Path:  "weight",
				Value: weight,
			},
		})

}

func main() {
	fmt.Println("Hello, world")

	/*	iter := firestore.Collection("ExcerciseSessions").Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Printf("Failed to iterate: %v", err)
				break
			}
			fmt.Printf("%s\n", doc.CreateTime.String())

		}*/

	/*_, err := userReference.Set(context.Background(), models.User{
	        Jobs:               []*firestore.DocumentRef{client.Doc("/selfManagedEmployees/K4qhd5k1c...")}
	})*/
	/*
		excercise_sessions := firestore.Collection("ExcerciseSessions").ID
		sets := firestore.Collection("Sets")
		session.Set(ctx, Set{Weight: 135, RepsOrDuration: 10, Done: false, IsTimeBased: false, TimeStampAdded: time.Now()})
		session = firestore.Collection("ExcerciseSessions").NewDoc()
		session.Set(ctx, Set{Weight: 135, RepsOrDuration: 10, Done: false, IsTimeBased: false, TimeStampAdded: time.Now()})
		session = firestore.Collection("ExcerciseSessions").NewDoc()
		session.Set(ctx, Set{Weight: 135, RepsOrDuration: 10, Done: false, IsTimeBased: false, TimeStampAdded: time.Now()})
	*/
	/*mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloHandler)

	mux.HandleFunc("/done_changed", DoneChange)
	mux.HandleFunc("/reps_change", RepsChange)
	mux.HandleFunc("/weight_change", WeightChanged)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	addr := ":8084"
	http.ListenAndServe(addr, mux)*/
	db := DB{}
	db.New()
	defer db.Delete()
	router := gin.Default()
	router.PUT("/done_changed/:id/:done", db.DoneChange)
	router.PUT("/reps_change/:id/:reps", db.RepsChange)
	router.PUT("/weight_change/:id/:weight", db.WeightChanged)

	router.Run("localhost:8080")

}
