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
        const new_data = {...exerciseData}
        let set =null;
        if(exerciseData.sets.length != 0)
        {
            let last_set = exerciseData.sets[exerciseData.sets.length-1];
            set ={
                set_id:"",
                weight:last_set.weight,
                reps_or_duration:last_set.reps_or_duration,
                done : false,
                is_time_based:false
      
              }
        }
        else{
        
        set ={
            set_id:"",
            weight:135,
            reps_or_duration:3,
            done : false,
            is_time_based:false
  
          }
        }
        new_data.sets.push(set)
        update(new_data);
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