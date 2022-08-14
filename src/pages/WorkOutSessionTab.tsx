import { IonContent, IonSelect, IonSelectOption,IonHeader, IonPage, IonTitle, IonToolbar, IonCard, IonAlert, IonItem } from '@ionic/react';
import React, { useState, useEffect } from 'react';
import './WorkOutSessionTab.css';
import {ExerciseSession, ExerciseSessionData} from '../components/ExcerciseSession';

export interface WorkoutData
{
   DoneChanged:Function;
   RepsChanged:Function;
   WeightChanged: Function;
   exercisesSessionsGet: Function
}

const WorkOutSessionTab: React.FC<WorkoutData> = ({DoneChanged, RepsChanged, WeightChanged, exercisesSessionsGet}) => {

  const initial_data :ExerciseSessionData ={
    name : "",
    date : new Date(),
    exercises : [],
    commpleted : false
  }
  

  const [selectedExerciseSession, setExcerciseData] = useState<ExerciseSessionData>();

  const [exercisesSessions, setExercisesSessions ] = useState([initial_data])

  const setCurrentExerciseSession = (s : String)=>{
    for(let i = 0; i < exercisesSessions.length; i++){
      if(exercisesSessions[i].name ==s){
        setExcerciseData(exercisesSessions[i]);
        break;
      }
    }
   
  }

  useEffect( ()=>{
    exercisesSessionsGet( (data:ExerciseSessionData[] )=>{
      if(data.length != 0){
        setExercisesSessions(data)
      }
      
    })
  },[ exercisesSessionsGet])

  return (
    <IonPage>
      {selectedExerciseSession &&
      <IonCard> 
        <IonSelect placeholder="Select Exercise Session" value={selectedExerciseSession.name}   onIonChange={(e)=>{setCurrentExerciseSession(e.detail.value!)}}  >
            {exercisesSessions.map( (s,index)=><IonSelectOption key={index}   value={s.name}>{s.name}</IonSelectOption> )}
          

      </IonSelect> 
     </IonCard>}
     {selectedExerciseSession ?
      <IonContent fullscreen>
        
        <ExerciseSession exerciseSession={selectedExerciseSession}  RepsChanged={RepsChanged} WeightChanged={WeightChanged} DoneChanged={DoneChanged} />
      </IonContent> 
      :
      <IonCard>
      <IonItem color="danger">No incomplete exercises found</IonItem>
      </IonCard>
     }
    </IonPage>
  );
};

export default WorkOutSessionTab;
