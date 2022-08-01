package main

//https://www.digitalocean.com/community/tutorials/debugging-go-code-with-visual-studio-code
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
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

	//excercise_sessions := db.firestore.Collection("ExcerciseSessions").ID
	/*sets := db.firestore.Collection("Sets")
	session := sets.NewDoc()
	session.Set(ctx, Set{Weight: 135, RepsOrDuration: 10, Done: false, IsTimeBased: false, TimeStampAdded: time.Now()})
	session = sets.NewDoc()
	session.Set(ctx, Set{Weight: 135, RepsOrDuration: 10, Done: false, IsTimeBased: false, TimeStampAdded: time.Now()})
	session = sets.NewDoc()
	session.Set(ctx, Set{Weight: 135, RepsOrDuration: 10, Done: false, IsTimeBased: false, TimeStampAdded: time.Now()})*/

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

	c.JSON(http.StatusOK, gin.H{"message": "OK"})

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
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
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
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (db *DB) AddExerciseSession(c *gin.Context) {

	var es ExerciseSession
	if err := c.ShouldBindJSON(&es); err != nil {
		c.JSON(400, gin.H{"message": "Bad data"})
		return
	}
	var exercise_session = db.firestore.Collection("ExerciseSessions").NewDoc()
	exercise_session.Set(context.Background(), struct {
		Name string
		Date time.Time
	}{Name: es.Name, Date: es.Date})
	for _, e := range es.Exercises {
		var exercise = db.firestore.Collection("Exercises").NewDoc()
		exercise.Set(context.Background(), struct {
			Name        string
			IsTimeBased bool
		}{Name: e.Name, IsTimeBased: e.IsTimeBased})
		exercise_session.Collection("Exercises").NewDoc().Set(context.Background(), struct{ ID string }{ID: exercise.ID})

		for _, s := range e.Sets {
			var set = db.firestore.Collection("Sets").NewDoc()
			exercise.Collection("Sets").Add(context.Background(), struct{ ID string }{ID: set.ID})
			set.Set(context.Background(), s)

		}
	}

	c.JSON(http.StatusOK, gin.H{"ID": exercise_session.ID})
}

func (db *DB) DeleteExerciseSession(c *gin.Context) {

	id := c.Param("id")
	ref := db.firestore.Collection("ExerciseSessions").Doc(id)
	if ref == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "id not found when deleting"})
		return
	}
	iter := ref.Collection("Exercises").Documents(context.Background())

	excercise_delete_batch := db.firestore.Batch()

	for {
		exercise, err := iter.Next()

		if err == iterator.Done {
			break
		}
		sets_iter := exercise.Ref.Collection("Sets").Documents(context.Background())
		for {
			set, err := sets_iter.Next()
			if err == iterator.Done {
				break
			}
			doc_map := set.Data()
			set_id, ok := doc_map["ID"].(string)
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"message": "id not found in sets when deleting"})
				return
			}

			excercise_delete_batch.Delete(db.firestore.Collection("Sets").Doc(set_id))

		}

		excercise_delete_batch.Delete(exercise.Ref)

	}
	excercise_delete_batch.Commit(context.Background())
	//mapbody:=make(map[string]string)

	//  json.Unmarshal(c.Request.Body,&mapbody)

	c.JSON(http.StatusOK, gin.H{"ID": id})
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
	router.POST("/api/done_changed/:id/:done", db.DoneChange)
	router.POST("/api/reps_change/:id/:reps", db.RepsChange)
	router.POST("/api/weight_change/:id/:weight", db.WeightChanged)
	router.POST("/api/add_exercise_session", db.AddExerciseSession)
	router.POST("/api/delete_exercise_session", db.DeleteExerciseSession)

	router.Run("localhost:8080")

}
