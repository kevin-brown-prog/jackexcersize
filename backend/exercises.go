package main

import "time"

/*
var set1 = { set_id : 0,
	weight : 135,
	reps_or_duration : 5,
	done : false,

	};

	var set2 = { set_id : 1,
	  weight : 135,
	  reps_or_duration : 5,
	  done : false,

	 };

	 var set3 = { set_id : 2,
	  weight : 135,
	  reps_or_duration : 5,
	  done : false,

	  };

	  var Exercise1 = {
		name:"Bench Press",
		sets :[set1, set2,set3]
		is_time_based: false
	  }


	   set1 = { set_id : 3,
		weight : 135,
		reps_or_duration : 5,
		done : false,

		};

		 set2 = { set_id : 4,
		  weight : 135,
		  reps_or_duration : 5,
		  done : false,

		 };

		  set3 = { set_id : 5,
		  weight : 135,
		  reps_or_duration : 5,
		  done : false,

		  };



	  var Exercise2 = {
		name:"Squats",
		sets:[set1, set2, set3]
		is_time_based:false
	  }

	  let ExerciseSession ={
		name:"Monday session",
		date:new Date(),
		exercises : [Exercise1, Exercise2]
	  }


weight : 135,
	  reps_or_duration : 5,
	  done : false,


*/

/*
type Set struct {
	Weight             int       `firestore:"weight" json:"weight"`
	RepsOrDuration     int       `firestore:"reps_or_duration" json:"reps_or_duration"`
	Done               bool      `firestore:"done" json:"done"`
	IsTimeBased        bool      `firestore:"is_time_based" json:"is_time_based"`
	TimeStampAdded     time.Time `firestore:"TimeStampAdded" json:"TimeStampAdded"`
	TimeStampCompleted time.Time `firestore:"TimeStampCompleted" json:"TimeStampCompleted"`
}
type Excercise struct {
	Name string `firestore:"name" json:"name"`
	Sets Set    `firestore:"sets" json:"sets"`
}
type ExcerciseSession struct {
	Name       string    `firestore:"name" json:"name"`
	Date       time.Time `firestore:"date" json:"date"`
	Excercises Excercise `firestore:"exercises" json:"excercises"`
}
*/
type Set struct {
	Weight             int       `firestore:"weight" json:"weight"`
	RepsOrDuration     int       `firestore:"reps_or_duration" json:"reps_or_duration"`
	Done               bool      `firestore:"done" json:"done"`
	SetID              int       `firestore:"set_id" json:"set_id"`
	ExerciseSessionID  int       `firestore:"exercise_session_id" json:"exercise_session_id"`
	TimeStampAdded     time.Time `firestore:"time_stamp_added" json:"time_stamp_added"`
	TimeStampCompleted time.Time `firestore:"time_stamp_completed" json:"time_stamp_completed"`
}
type Exercise struct {
	Name        string `firestore:"name" json:"name"`
	Sets        []Set  `firestore:"sets" json:"sets"`
	IsTimeBased bool   `firestore:"is_time_based" json:"is_time_based"`
}
type ExerciseSession struct {
	Name      string     `firestore:"name" json:"name"`
	Date      time.Time  `firestore:"date" json:"date"`
	Exercises []Exercise `firestore:"exercises" json:"exercises"`
	Completed bool       `firestore:"completed" json:"completed"`
}
