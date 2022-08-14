package main

//https://www.digitalocean.com/community/tutorials/debugging-go-code-with-visual-studio-code
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	done, err := strconv.ParseBool(c.Param("done"))
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("done not bool %v", err)})
		return
	}

	set := db.firestore.Collection("Sets").Doc(id)
	setData, err := set.Get(context.Background())
	if err != nil {
		c.JSON(400, gin.H{"message": "Unable to find document"})
		return
	}
	exercise_session_id := setData.Data()["exercise_session_id"].(string)
	set.Update(context.Background(),
		[]firestore.Update{
			{
				Path:  "done",
				Value: done,
			},
		})

	query := db.firestore.Collection("Sets").Where("exercise_session_id", "==", exercise_session_id)
	iter := query.Documents(context.Background())
	complete := true
	for {
		setDoc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		setLocal := Set{}
		setDoc.DataTo(&setLocal)
		if !setLocal.Done {
			complete = false
			break
		}
	}
	if complete {
		db.firestore.Collection(("ExerciseSessions")).Doc(exercise_session_id).Update(context.Background(),
			[]firestore.Update{
				{
					Path:  "completed",
					Value: complete,
				},
			})

	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})

}
func (db *DB) RepsChange(c *gin.Context) {
	id := c.Param("id")
	reps, err := strconv.Atoi(c.Param("reps"))
	if err != nil {
		c.JSON(400, gin.H{"message": "reps not int"})
		return
	}
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

	weight, err := strconv.Atoi(c.Param("weight"))
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("weight not int %v", err)})
		return
	}
	db.firestore.Collection("Sets").Doc(id).Update(context.Background(),
		[]firestore.Update{
			{
				Path:  "weight",
				Value: weight,
			},
		})
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (db *DB) GetSessionsNotComplete(c *gin.Context) {
	query := db.firestore.Collection("ExerciseSessions").Where("completed", "!=", true)
	iter := query.Documents(context.Background())
	var docs []ExerciseSession = make([]ExerciseSession, 0, 10)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		} else if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return
		}

		excersizeSession := doc.Data()

		es := ExerciseSession{}
		es.Completed = excersizeSession["completed"].(bool)
		es.Date = excersizeSession["date"].(time.Time)
		es.Name = excersizeSession["name"].(string)
		exercisesIter := doc.Ref.Collection("Exercises").Documents(context.Background())
		es.Exercises = make([]Exercise, 0, 5)
		for {
			exercieseDoc, err := exercisesIter.Next()
			if err == iterator.Done {
				break
			}

			exerciseDocValue, err := db.firestore.Collection("Exercises").Doc(exercieseDoc.Data()["ID"].(string)).Get(context.Background())
			if err != nil {
				c.JSON(400, gin.H{"message": "Unable to find document"})
				return
			}
			exerciseMap := exerciseDocValue.Data()
			exercise := Exercise{}
			exercise.Name = exerciseMap["name"].(string)
			exercise.IsTimeBased = exerciseMap["is_time_based"].(bool)
			exercise.Sets = make([]Set, 0, 5)
			setsIter := exerciseDocValue.Ref.Collection("Sets").Documents(context.Background())

			for {
				setDocument, err := setsIter.Next()
				if err == iterator.Done {
					break
				}
				setsDocValue, err := db.firestore.Collection("Sets").Doc(setDocument.Data()["ID"].(string)).Get(context.Background())
				if err != nil {
					c.JSON(400, gin.H{"message": "Unable to find document"})
					return
				}
				set := Set{}
				err = setsDocValue.DataTo(&set)
				if err != nil {
					c.JSON(400, gin.H{"message": "Unable to find document"})
					return
				}

				exercise.Sets = append(exercise.Sets, set)
			}
			es.Exercises = append(es.Exercises, exercise)
		}

		docs = append(docs, es)

	}
	c.JSON(http.StatusOK, docs)
}

func (db *DB) AddExerciseSession(c *gin.Context) {

	var es ExerciseSession
	if err := c.ShouldBindJSON(&es); err != nil {
		c.JSON(400, gin.H{"message": "Bad data"})
		return
	}
	var exercise_session = db.firestore.Collection("ExerciseSessions").NewDoc()

	ns := struct {
		Name      string    `firestore:"name" json:"name"`
		Date      time.Time `firestore:"date" json:"date"`
		Completed bool      `firestore:"completed" json:"completed"`
	}{Name: es.Name, Date: es.Date, Completed: false}

	exercise_session.Set(context.Background(), ns)
	for _, e := range es.Exercises {
		var exercise = db.firestore.Collection("Exercises").NewDoc()
		exercise.Set(context.Background(), struct {
			Name        string `firestore:"name" json:"name"`
			IsTimeBased bool   `firestore:"is_time_based" json:"is_time_based"`
		}{Name: e.Name, IsTimeBased: e.IsTimeBased})
		exercise_session.Collection("Exercises").NewDoc().Set(context.Background(), struct{ ID string }{ID: exercise.ID})

		for _, s := range e.Sets {
			var set = db.firestore.Collection("Sets").NewDoc()
			exercise.Collection("Sets").Add(context.Background(), struct{ ID string }{ID: set.ID})
			s.SetID = set.ID
			s.ExerciseSessionID = exercise_session.ID
			s.TimeStampAdded = time.Now()
			s.TimeStampCompleted = time.Unix(0, 0)
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

	router.GET("/api/get_sessions_not_complete", db.GetSessionsNotComplete)

	router.Run("localhost:8080")

}
