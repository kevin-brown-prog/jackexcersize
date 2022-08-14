import React from 'react';

import Exercise, {ExerciseData} from './Excercise'
import { IonList, IonItem, IonLabel } from '@ionic/react';
	




export interface ExerciseSessionData {
  name: string;
  date : Date
  exercises : Array<ExerciseData>;
  commpleted : boolean;
}

interface ExcerssizeSessionProp
{
    exerciseSession : ExerciseSessionData 
    DoneChanged : Function;
    RepsChanged : Function;
    WeightChanged: Function;
}

export const ExerciseSession: React.FC<ExcerssizeSessionProp> = ({ exerciseSession, DoneChanged, RepsChanged, WeightChanged }) => {
	
  return (
       <IonItem>
          
       <IonList>
       <IonItem>
           <IonLabel>{exerciseSession.name}</IonLabel>
        </IonItem>
         {exerciseSession.exercises.map( (s,index)=><Exercise key={index}  DoneChanged={DoneChanged} WeightChanged={WeightChanged} RepsChanged={RepsChanged} exerciseData={s} />)}
       </IonList>
       </IonItem>
   
   
  );
};

