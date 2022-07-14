import { IonContent , IonList, IonItem, IonInput,IonIcon, IonHeader, IonPage, IonTitle, IonToolbar, IonGrid, IonRow, IonCol, IonButton, IonLabel } from '@ionic/react';
import { add, settings, share, person, arrowForwardCircle, arrowBackCircle, arrowUpCircle, logoVimeo, logoFacebook, logoInstagram, logoTwitter } from 'ionicons/icons';

import ExerciseAdd from '../components/ExcerciseAdd'
import ExerciseEdit from '../components/ExerciseEdit'
import {ExerciseData} from '../components/Excercise'
import './ExcerciseSessionCreate.css';
import React, { useState, useEffect } from 'react';
import {ExerciseSessionData} from '../components/ExcerciseSession'
import { SetData } from '../components/Set';





const ExcerciseSessionCreate: React.FC = () => {
  const initial_data :ExerciseSessionData ={
    name : "",
    date : new Date(),
    exercises : []
  }
  
  const [exerciseData, setExcerciseData] = useState(initial_data as ExerciseSessionData);



  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonTitle>Create Excersise Session</IonTitle>
        </IonToolbar>
      </IonHeader>
      <IonContent fullscreen>
     
      <IonItem>
            <IonInput value={exerciseData.name} placeholder="Session Name" 
              onIonChange={e => setExcerciseData( ()=>{
                 const new_data: ExerciseSessionData ={...exerciseData};
                 new_data.name = e.detail.value!;
                 return new_data;
              }
                )}></IonInput>
      </IonItem>
      <IonList>
     
              {exerciseData.exercises.map((data, index)=>
                  <ExerciseEdit key={index} exerciseData={data} update={(newData : ExerciseData)=>{
                    const new_data: ExerciseSessionData ={...exerciseData};
                    new_data.exercises[index] = newData;
                      setExcerciseData(new_data);
                  }} />
              )}
      
       </IonList>
      <IonGrid>
       <ExerciseAdd NewExercise={(e_name:string)=>{
        const new_data = {...exerciseData}
        const new_exercise : ExerciseData={
          name : e_name,
         sets : []
        }
        const set : SetData ={
          set_id:"",
          weight:135,
          reps_or_duration:3,
          done : false,
          is_time_based:false

        }
      new_exercise.sets.push(set)
        new_data.exercises.push(new_exercise)
        setExcerciseData(new_data);
       }}/>

       </IonGrid>
    </IonContent>
    </IonPage>
  );
};

export default ExcerciseSessionCreate;
