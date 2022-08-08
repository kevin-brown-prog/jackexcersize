import { IonContent, IonSelect, IonSelectOption,IonHeader, IonPage, IonTitle, IonToolbar } from '@ionic/react';
import React, { useState, useEffect } from 'react';
import './WorkOutSessionTab.css';
import {ExerciseSession, ExerciseSessionData} from '../components/ExcerciseSession';

export interface WorkoutData
{
   DoneChanged:Function;
   RepsChanged:Function;
   WeightChanged: Function;
   exercisesSessions:ExerciseSessionData[]
}

const WorkOutSessionTab: React.FC<WorkoutData> = ({DoneChanged, RepsChanged, WeightChanged, exercisesSessions}) => {

  const initial_data :ExerciseSessionData ={
    name : "",
    date : new Date(),
    exercises : []
  }


  const [selectedExerciseSession, setExcerciseData] = useState(exercisesSessions.length == 0? initial_data : exercisesSessions[0]);

  return (
    <IonPage>
      <IonSelect placeholder="Select Exercise Session" value={selectedExerciseSession}   onIonChange={(e)=>{setExcerciseData(e.detail.value!)}}  >
          {exercisesSessions.map( (s,index)=><IonSelectOption key={index}   value={s}>{s}</IonSelectOption> )}
         

     </IonSelect>
      <IonContent fullscreen>
        
        <ExerciseSession exerciseSession={selectedExerciseSession}  RepsChanged={RepsChanged} WeightChanged={WeightChanged} DoneChanged={DoneChanged} />
      </IonContent>
    </IonPage>
  );
};

export default WorkOutSessionTab;
