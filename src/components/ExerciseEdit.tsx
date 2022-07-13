import React from 'react';
import Set, {SetData, SetEdit} from './Set';
import {  IonList, IonItem, IonLabel, IonListHeader } from '@ionic/react';




import {ExerciseData} from '../components/Excercise'
import SetAdd from '../components/SetAdd'

interface ExcersizeEditProp
{
  exerciseData : ExerciseData;
  update : Function
}



 const ExerciseEdit :React.FC<ExcersizeEditProp> = ({exerciseData, update}) => {
    const NewSet = ()=>{

    }
    return (
        <IonItem >
         <IonList>
            <IonListHeader>
             <IonLabel>{exerciseData.name}</IonLabel>
          </IonListHeader>  
           {exerciseData.sets.map( (s,index)=><SetEdit key={index} index={index} setData={s} UpdateData={(i:number, new_set :SetData)=>{
            const ret = {...exerciseData}
            ret.sets[i] = new_set;
            update(ret);
           }}/>)}
           <SetAdd NewSet={NewSet} />
           </IonList>
       </IonItem>
    );

  };

  export default ExerciseEdit;