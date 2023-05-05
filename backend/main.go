package main

//https://www.digitalocean.com/community/tutorials/debugging-go-code-with-visual-studio-code
import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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

	var err error
	db.db, err = sql.Open("sqlite3", "./workout.db")
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

func Transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
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
	Transact(db.db, func(tx *sql.Tx) error {
		return DoneMethod(tx, id_int, c, done)

	})

}

func DoneMethod(tx *sql.Tx, id_int int, c *gin.Context, done bool) error {
	stmt, err := tx.Prepare("SELECT exerciseid FROM sets where id=?")
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}

	rows, err := stmt.Query(id_int)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	var exerciseId int
	if rows.Next() {
		err = rows.Scan(&exerciseId)
		if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return err
		}
	} else {
		c.JSON(400, gin.H{"message": "Unable to find set"})
		return fmt.Errorf("No rows in sets found for %v", id_int)
	}

	stmt, err = tx.Prepare("SELECT exerciseSessionId FROM Exercises where id=?")
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}

	rows, err = stmt.Query(exerciseId)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	var exerciseSessionId int
	if rows.Next() {
		err = rows.Scan(&exerciseSessionId)
		if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return err
		}
	} else {
		c.JSON(400, gin.H{"message": "Unable to find exercise"})
		return fmt.Errorf("No rows in Exercises found for %v", exerciseId)
	}

	stmt, err = tx.Prepare("SELECT User from ExerciseSessions where id=?")
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	rows, err = stmt.Query(exerciseSessionId)
	var user string
	if rows.Next() {
		err = rows.Scan(&user)
		if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return err
		}
	} else {
		c.JSON(400, gin.H{"message": "Unable to find ExerciseSession"})
		return fmt.Errorf("No rows in ExerciseSessions found for %v", exerciseSessionId)
	}
	fmt.Println(user)
	stmt, err = tx.Prepare("update sets set done=?, TimestampCompleted=? where id=?")
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}

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
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	if affect != 1 {
		panic("not able to update set")
	}

	stmt, err = tx.Prepare("SELECT done from sets where exerciseid=?")
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	rows, err = stmt.Query(exerciseId)
	complete := 1
	for rows.Next() {
		var done bool
		err = rows.Scan(&done)
		if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return err
		}
		if done {
			complete = 0
			break
		}
	}

	stmt, err = tx.Prepare("Update ExerciseSessions set complete=? where id=?")
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	res, err = stmt.Exec(complete, exerciseSessionId)
	affect, err = res.RowsAffected()
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	if affect != 1 {
		panic("not able to update set")
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
	return nil
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
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}

	res, err := stmt.Exec(reps, id_int)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}
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
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}

	res, err := stmt.Exec(weight, id_int)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}
	if affect != 1 {
		c.JSON(400, gin.H{"message": fmt.Sprintf("Unable to find set")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func (db *DB) GetSessionsNotComplete(c *gin.Context) {

	/*		"ID"	INTEGER NOT NULL UNIQUE,
			"Name"	TEXT NOT NULL,
			"IsTimeBased"	INTEGER NOT NULL,
			"ExerciseSessionID"	INTEGER NOT NULL,*/
	/*
								"ID"	INTEGER NOT NULL UNIQUE,
						"Weight"	INTEGER NOT NULL,
						"RepsOrDuration"	INTEGER NOT NULL,
						"Done"	INTEGER NOT NULL,
						"ExerciseID"	INTEGER NOT NULL,
						"TimestampAdded"	INTEGER NOT NULL,
		            	"TimestampCompleted"	INTEGER NOT NULL,


	*/
	Transact(db.db, func(tx *sql.Tx) error {
		return GetSessionsNotCompleteImp(tx, c)
		//return DoneMethod(tx, id_int, c, done)

	})

}

func GetSessionsNotCompleteImp(tx *sql.Tx, c *gin.Context) error {
	stmt, err := tx.Prepare("SELECT ID, Complete, DateComplete,Name from ExerciseSessions where complete=0")
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	rows, err := stmt.Query()
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return err
	}
	sessions := make([]ExerciseSession, 0, 5)
	for rows.Next() {

		var complete int
		var dateComplete int64
		var name string
		var exerciseSessionId int
		err = rows.Scan(&exerciseSessionId, &complete, &dateComplete, &name)
		if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return err
		}
		session := ExerciseSession{Name: name,
			Date:      time.Unix(dateComplete, 0),
			Completed: complete != 0,
			Exercises: make([]Exercise, 0, 5),
		}

		stmt, err := tx.Prepare("SELECT ID, Name,IsTimeBased from Exercises where ExerciseSessionID=?")
		checkErr(err)
		rowsExercises, err := stmt.Query(exerciseSessionId)
		checkErr(err)
		for rowsExercises.Next() {

			var exerciseId int
			var name string
			var isTimeBased int
			err = rowsExercises.Scan(&exerciseId, &name, &isTimeBased)
			if err != nil {
				c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
				return err
			}
			exercise := Exercise{Name: name, IsTimeBased: isTimeBased != 0, Sets: make([]Set, 0, 5)}

			stmt, err := tx.Prepare("SELECT ID, Weight,RepsOrDuration,Done,TimestampAdded,TimestampCompleted from Sets where ExerciseID=?")
			if err != nil {
				c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
				return err
			}
			rowsSets, err := stmt.Query(exerciseId)
			if err != nil {
				c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
				return err
			}
			for rowsSets.Next() {
				var set_id int
				var weight int
				var repsOrDuration int
				var done int
				var timeStampCompleted int64
				var timeStampAdded int64
				err = rowsSets.Scan(&set_id, &weight, &repsOrDuration, &done, &timeStampCompleted, &timeStampAdded)
				if err != nil {
					c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
					return err
				}
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
	return nil
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

	/*
				"ID"	INTEGER NOT NULL UNIQUE,
		"Weight"	INTEGER NOT NULL,
		"RepsOrDuration"	INTEGER NOT NULL,
		"Done"	INTEGER NOT NULL,
		"ExerciseID"	INTEGER NOT NULL,
		, TimestampAdded, TimestampComplete*/
	var exerciseSessionId int
	Transact(db.db, func(tx *sql.Tx) error {
		var err2 error
		exerciseSessionId, err2 = AddExerciseSessionTransaction(tx, c, es)
		return err2
	})

	c.JSON(http.StatusOK, gin.H{"ID": exerciseSessionId})
}

func AddExerciseSessionTransaction(tx *sql.Tx, c *gin.Context, es ExerciseSession) (int, error) {
	stmt, err := tx.Prepare("INSERT INTO ExerciseSessions(Complete, DateComplete, User, Name) values(?,?,?,?)")
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return -1, err
	}
	res, err := stmt.Exec(es.Completed, 0, "bob", es.Name)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return -1, err
	}
	exerciseSessionId, err := res.LastInsertId()
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return -1, err
	}

	for _, e := range es.Exercises {

		stmt, err := tx.Prepare("INSERT INTO Exercises(Name, IsTimeBased, ExerciseSessionID) values(?,?,?)")
		if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return -1, err
		}
		res, err := stmt.Exec(e.Name, e.IsTimeBased, exerciseSessionId)
		if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return -1, err
		}
		exerciseId, err := res.LastInsertId()
		if err != nil {
			c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
			return -1, err
		}
		for _, s := range e.Sets {

			stmt, err := tx.Prepare("INSERT INTO Sets(Weight, RepsOrDuration, Done,ExerciseID, TimestampAdded, TimestampCompleted) values(?,?,?,?,?,?)")
			if err != nil {
				c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
				return -1, err
			}
			_, err = stmt.Exec(s.Weight, s.RepsOrDuration, s.Done, exerciseId, time.Now().Unix(), 0)
			if err != nil {
				c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
				return -1, err
			}
		}

	}
	return int(exerciseSessionId), nil
}

func (db *DB) DeleteExerciseSession(c *gin.Context) {

	id := c.Param("id")

	id_int, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"message": fmt.Sprintf("%v", err)})
		return
	}

	stmt, err := db.db.Prepare("SELECT ID from Exercise where ExerciseSessionID=?")
	checkErr(err)
	rows, err := stmt.Query(id_int)
	exerciseIds := make([]int, 0, 5)
	if rows.Next() {
		var exerciseId int
		err = rows.Scan(&exerciseId)
		if err != nil {
			checkErr(err)
		}
		exerciseIds = append(exerciseIds, exerciseId)
		deleteStmt, err := db.db.Prepare("Delete from Sets where ExerciseID=?")
		checkErr(err)
		_, err = deleteStmt.Exec(id_int)
		checkErr(err)

	} else {
		c.JSON(400, gin.H{"message": "Unable to find ExerciseSession"})
		return
	}
	for _, id := range exerciseIds {
		deleteStmt, err := db.db.Prepare("Delete from Exercises where ID=?")
		checkErr(err)
		_, err = deleteStmt.Exec(id)
		checkErr(err)
	}
	deleteStmt, err := db.db.Prepare("Delete from ExerciseSessions where ID=?")
	checkErr(err)
	_, err = deleteStmt.Exec(id_int)
	checkErr(err)

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
