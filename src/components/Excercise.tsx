import React from 'react';
import Set, {SetData, SetEdit} from './Set';
import {  IonList, IonItem, IonLabel, IonListHeader } from '@ionic/react';

export const ExerciseEdit :React.FC<ExcersizeEditProp> = ({exerciseData, update}) => {
  return (
      <IonItem >
       <IonList>
       <IonListHeader>
           <IonLabel>here{exerciseData.name}</IonLabel>
        </IonListHeader>  
         {exerciseData.sets.map( (s,index)=><SetEdit key={index} index={index} setData={s} UpdateData={(i:number, new_set :SetData)=>{
          const ret = {...exerciseData}
          ret.sets[i] = new_set;
          update(ret);
         }}/>)}
       </IonList> 
     </IonItem>
  );
};



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
interface ExcersizeEditProp
{
  exerciseData : ExerciseData;
  update : Function
}

