import { IonContent, IonSelect, IonSelectOption,IonHeader, IonPage, IonTitle, IonToolbar, IonCard, IonAlert, IonItem } from '@ionic/react';
import React, { useState, useEffect } from 'react';
import './WorkOutSessionTab.css';
import {ExerciseSession, ExerciseSessionData} from '../components/ExcerciseSession';

export interface WorkoutData
{
   DoneChanged:Function;
   RepsChanged:Function;
   WeightChanged: Function;
   exercisesSessionsNotComplete: ExerciseSessionData[]
}

const WorkOutSessionTab: React.FC<WorkoutData> = ({DoneChanged, RepsChanged, WeightChanged, exercisesSessionsNotComplete}) => {

  

  const [selectedExerciseSession, setExcerciseData] = useState<ExerciseSessionData>(exercisesSessionsNotComplete[0]);

 

  const setCurrentExerciseSession = (s : String)=>{
    for(let i = 0; i < exercisesSessionsNotComplete.length; i++){
      if(exercisesSessionsNotComplete[i].name ==s){
        setExcerciseData(exercisesSessionsNotComplete[i]);
        break;
      }
    }
   
  }

  

  return (
    <IonPage class="limit-width">
      
      <IonCard> 
        <IonSelect placeholder="Select Exercise Session" value={selectedExerciseSession.name}   onIonChange={(e)=>{setCurrentExerciseSession(e.detail.value!)}}  >
            {exercisesSessionsNotComplete.map( (s,index)=><IonSelectOption key={index}   value={s.name}>{s.name}</IonSelectOption> )}
          

      </IonSelect> 
     </IonCard>
    
      <IonContent fullscreen>
        
        <ExerciseSession exerciseSession={selectedExerciseSession}  RepsChanged={RepsChanged} WeightChanged={WeightChanged} DoneChanged={DoneChanged} />
      </IonContent> 
    
    </IonPage>
  );
};

export default WorkOutSessionTab;
