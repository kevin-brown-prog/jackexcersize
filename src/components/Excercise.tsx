import React from 'react';
import Set, {SetData} from './Set';
import {  IonList, IonItem, IonLabel, IonListHeader } from '@ionic/react';
	




export interface ExerciseData {
  name: string;
  sets : Array<SetData>;

}

 interface ExerciseProp{
  exerciseData: ExerciseData,
 
    DoneChanged : Function;
    RepsChanged : Function;
    WeightChanged: Function;
 }

const Exercise: React.FC<ExerciseProp> = ({  exerciseData, DoneChanged, RepsChanged, WeightChanged }) => {
	
  return (
      <IonItem >
       
       <IonList>
       <IonListHeader>
           <IonLabel>{exerciseData.name}</IonLabel>
        </IonListHeader>  
         {exerciseData.sets.map( (s,index)=><Set key={index} setData={s}  DoneChanged={DoneChanged} RepsChanged={RepsChanged} WeightChanged={WeightChanged}/>)}
       </IonList>
     </IonItem>
  
  );
};

export default Exercise;