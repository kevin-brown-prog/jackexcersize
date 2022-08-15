package main

//https://www.digitalocean.com/community/tutorials/debugging-go-code-with-visual-studio-code
import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type DB struct {
	firestore *firestore.Client
	app       *firebase.App
	db        *sql.DB
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
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

	db.db, err = sql.Open("sqlite3", "./foo.db")
	checkErr(err)

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

	id_int, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}

	stmt, err := db.db.Prepare("SELECT exerciseid FROM sets where id=?")
	checkErr(err)

	rows, err := stmt.Query(id_int)
	checkErr(err)
	var exerciseId int
	if rows.Next() {
		err = rows.Scan(&exerciseId)
		if err != nil {
			checkErr(err)
		}
	} else {
		c.JSON(400, gin.H{"message": "Unable to find set"})
		return
	}

	stmt, err = db.db.Prepare("SELECT exerciseSessionId FROM Exercises where id=?")
	checkErr(err)

	rows, err = stmt.Query(exerciseId)
	checkErr(err)
	var exerciseSessionId int
	if rows.Next() {
		err = rows.Scan(&exerciseSessionId)
		if err != nil {
			checkErr(err)
		}
	} else {
		c.JSON(400, gin.H{"message": "Unable to find set"})
		return
	}

	stmt, err = db.db.Prepare("SELECT User from ExerciseSessions where id=?")
	checkErr(err)
	rows, err = stmt.Query(exerciseSessionId)
	var user string
	if rows.Next() {
		err = rows.Scan(&user)
		if err != nil {
			checkErr(err)
		}
	} else {
		c.JSON(400, gin.H{"message": "Unable to find ExerciseSession"})
		return
	}
	fmt.Println(user)
	stmt, err = db.db.Prepare("update sets set done=?, TimestampCompleted=? where id=?")
	checkErr(err)

	var is_done int
	var timeStamp int64
	if done {
		is_done = 1
		timeStamp = time.Now().Unix()
	} else {
		is_done = 0
		timeStamp = 0
	}
	res, err := stmt.Exec(is_done, timeStamp, id_int)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("not able to update set")
	}

	stmt, err = db.db.Prepare("SELECT done from sets where exerciseid=?")
	checkErr(err)
	rows, err = stmt.Query(exerciseId)
	complete := 1
	for rows.Next() {
		var done bool
		err = rows.Scan(&done)
		if err != nil {
			checkErr(err)
		}
		if done {
			complete = 0
			break
		}
	}

	stmt, err = db.db.Prepare("Update ExerciseSessions set complete=? where id=?")
	checkErr(err)
	res, err = stmt.Exec(complete, exerciseSessionId)
	affect, err = res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		panic("not able to update set")
	}

	/*

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

		}*/

	c.JSON(http.StatusOK, gin.H{"message": "OK"})

}
func (db *DB) RepsChange(c *gin.Context) {
	id := c.Param("id")
	reps, err := strconv.Atoi(c.Param("reps"))
	if err != nil {
		c.JSON(400, gin.H{"message": "reps not int"})
		return
	}
	id_int, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}
	stmt, err := db.db.Prepare("update sets set repsorduration=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(reps, id_int)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		c.JSON(400, gin.H{"message": fmt.Sprintf("Unable to find set")})
		return
	}
	/*
		db.firestore.Collection("Sets").Doc(id).Update(context.Background(),
			[]firestore.Update{
				{
					Path:  "reps_or_duration",
					Value: reps,
				},
			})*/
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
func (db *DB) WeightChanged(c *gin.Context) {
	id := c.Param("id")

	weight, err := strconv.Atoi(c.Param("weight"))
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("weight not int %v", err)})
		return
	}

	id_int, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}
	stmt, err := db.db.Prepare("update sets set weight=? where id=?")
	checkErr(err)

	res, err := stmt.Exec(weight, id_int)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	if affect != 1 {
		c.JSON(400, gin.H{"message": fmt.Sprintf("Unable to find set")})
		return
	}
	/*
		db.firestore.Collection("Sets").Doc(id).Update(context.Background(),
			[]firestore.Update{
				{
					Path:  "weight",
					Value: weight,
				},
			})*/
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (db *DB) GetSessionsNotComplete(c *gin.Context) {

	stmt, err := db.db.Prepare("SELECT ID, Complete, DateComplete,Name from ExerciseSessions where complete=0")
	checkErr(err)
	rows, err := stmt.Query()
	checkErr(err)
	sessions := make([]ExerciseSession, 0, 5)
	for rows.Next() {

		var complete int
		var dateComplete int64
		var name string
		var exerciseSessionId int
		err = rows.Scan(&exerciseSessionId, &complete, &dateComplete, &name)
		if err != nil {
			checkErr(err)
		}
		session := ExerciseSession{Name: name,
			Date:      time.Unix(dateComplete, 0),
			Completed: complete != 0,
			Exercises: make([]Exercise, 0, 5),
		}

		stmt, err := db.db.Prepare("SELECT ID, Name,IsTimeBased from Exercises where ExerciseSessionID=?")
		checkErr(err)
		rowsExercises, err := stmt.Query(exerciseSessionId)
		checkErr(err)
		for rowsExercises.Next() {
			/*		"ID"	INTEGER NOT NULL UNIQUE,
					"Name"	TEXT NOT NULL,
					"IsTimeBased"	INTEGER NOT NULL,
					"ExerciseSessionID"	INTEGER NOT NULL,*/
			var exerciseId int
			var name string
			var isTimeBased int
			err = rowsExercises.Scan(&exerciseId, &name, &isTimeBased)
			checkErr(err)
			exercise := Exercise{Name: name, IsTimeBased: isTimeBased != 0, Sets: make([]Set, 0, 5)}

			/*
										"ID"	INTEGER NOT NULL UNIQUE,
								"Weight"	INTEGER NOT NULL,
								"RepsOrDuration"	INTEGER NOT NULL,
								"Done"	INTEGER NOT NULL,
								"ExerciseID"	INTEGER NOT NULL,
								"TimestampAdded"	INTEGER NOT NULL,
				            	"TimestampCompleted"	INTEGER NOT NULL,


			*/

			stmt, err := db.db.Prepare("SELECT ID, Weight,RepsOrDuration,Done,TimestampAdded,TimestampCompleted from Sets where ExerciseID=?")
			checkErr(err)
			rowsSets, err := stmt.Query(exerciseId)
			checkErr(err)
			for rowsSets.Next() {
				var set_id int
				var weight int
				var repsOrDuration int
				var done int
				var timeStampCompleted int64
				var timeStampAdded int64
				err = rowsSets.Scan(&set_id, &weight, &repsOrDuration, &done, &timeStampCompleted, &timeStampAdded)
				checkErr(err)
				set := Set{SetID: set_id, RepsOrDuration: repsOrDuration, Done: done != 0, ExerciseSessionID: exerciseSessionId,
					TimeStampAdded:     time.Unix(timeStampCompleted, 0),
					TimeStampCompleted: time.Unix(timeStampCompleted, 0)}
				exercise.Sets = append(exercise.Sets, set)
			}
			session.Exercises = append(session.Exercises, exercise)

		}

		sessions = append(sessions, session)

	}

	c.JSON(http.StatusOK, sessions)
}

func (db *DB) AddExerciseSession(c *gin.Context) {

	var es ExerciseSession
	if err := c.ShouldBindJSON(&es); err != nil {
		c.JSON(400, gin.H{"message": "Bad data"})
		return
	}
	/*
		CREATE TABLE "ExerciseSessions" (
			"ID"	INTEGER NOT NULL UNIQUE,
			"Complete"	INTEGER NOT NULL,
			"DateComplete"	INTEGER NOT NULL,
			"User"	TEXT NOT NULL,
			"Name"	TEXT NOT NULL,
			PRIMARY KEY("ID" AUTOINCREMENT)
		);*/

	stmt, err := db.db.Prepare("INSERT INTO ExerciseSessions(Complete, DateComplete, User, Name) values(?,?,?,?)")
	checkErr(err)
	res, err := stmt.Exec(es.Completed, 0, "bob", es.Name)
	checkErr(err)
	exerciseSessionId, err := res.LastInsertId()
	checkErr(err)

	for _, e := range es.Exercises {

		stmt, err := db.db.Prepare("INSERT INTO Exercises(Name, IsTimeBased, ExerciseSessionID) values(?,?,?)")
		checkErr(err)
		res, err := stmt.Exec(e.Name, e.IsTimeBased, exerciseSessionId)
		checkErr(err)
		exerciseId, err := res.LastInsertId()
		checkErr(err)
		for _, s := range e.Sets {

			/*
						"ID"	INTEGER NOT NULL UNIQUE,
				"Weight"	INTEGER NOT NULL,
				"RepsOrDuration"	INTEGER NOT NULL,
				"Done"	INTEGER NOT NULL,
				"ExerciseID"	INTEGER NOT NULL,
				, TimestampAdded, TimestampComplete*/

			stmt, err := db.db.Prepare("INSERT INTO Sets(Weight, RepsOrDuration, Done,ExerciseID, TimestampAdded, TimestampCompleted) values(?,?,?,?,?,?)")
			checkErr(err)
			_, err = stmt.Exec(s.Weight, s.RepsOrDuration, s.Done, exerciseId, time.Now().Unix(), 0)
			checkErr(err)
		}

	}

	stmt, err := db.db.Prepare("SELECT ID, Complete, DateComplete,Name from ExerciseSessions where complete=0")
	checkErr(err)
	rows, err := stmt.Query()
	checkErr(err)

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
