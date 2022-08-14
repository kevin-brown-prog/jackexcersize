import React, { useState, useEffect } from 'react';

import { IonInput, IonCheckbox, IonListHeader, IonItem, IonLabel} from '@ionic/react';
	

export interface SetData{
    set_id : string;
    weight : number;
    reps_or_duration : number;
    done : boolean;
  
}

interface SetProperty
{
    setData : SetData;
    WeightChanged : Function;
    DoneChanged : Function;
    RepsChanged : Function;
    is_time_based : boolean
}




const Set: React.FC<SetProperty> = ({  setData, DoneChanged, RepsChanged, WeightChanged, is_time_based }) => {
	const [isDone, setIsDone] = useState(setData.done);
    const [weightComplete, setWeightComplete] = useState(setData.weight)
    const [repsComplete, setRepsComplete] = useState(setData.reps_or_duration);
  /*  useEffect(() => {
        DoneChanged(setData.set_id,isDone);
      }, [isDone, DoneChanged, setData.set_id]); // Only re-run the effect if count changes

    useEffect(() => {
        
        RepsChanged(setData.set_id,repsComplete);
      }, [repsComplete, RepsChanged, setData.set_id]); // Only re-run the effect if count changes


    useEffect(()=>{
            WeightChanged(setData.set_id, weightComplete);
    }, [weightComplete, WeightChanged, setData.set_id]);*/

  const OnChangeReps = (val : string)=>{
   
        const re = /^[0-9\b]+$/;
    
        // if value is not blank, then test the regex
    
        if (val === '' || re.test(val)) {
            let new_val = parseInt(val);
            setRepsComplete(new_val);
            RepsChanged(setData.set_id, new_val)
        }
        else{
            setRepsComplete(repsComplete);
        }
  }
  const OnDoneChanged = (val : boolean)=>{
    setIsDone(val);
    DoneChanged(setData.set_id, val)
  }

  const OnChangeWeight = (val : string)=>{
   
    const re = /^[0-9\b]+$/;

    // if value is not blank, then test the regex

    if (val === '' || re.test(val)) {
        let new_weight = parseInt(val);
        setWeightComplete(new_weight);
        WeightChanged(new_weight)

    }
    else{
        setWeightComplete(weightComplete);
    }
}

  return (
    
         
    <IonListHeader>
      <IonItem>
         
        <IonCheckbox checked={isDone} onIonChange={e => OnDoneChanged(e.detail.checked)} />
		
      </IonItem>
      <IonItem>
        <IonLabel>Weight</IonLabel>
      </IonItem>
      <IonItem>
      <IonInput inputmode="numeric" value={weightComplete} onIonInput={(e: any)=>OnChangeWeight(e.target.value )}></IonInput>
      </IonItem>
      <IonItem>
        <IonLabel>{is_time_based? "Duration seconds" : "Reps"}</IonLabel>
      </IonItem>
      <IonItem>
      <IonInput inputmode="numeric" pattern="[0-9]" min={0}  value={repsComplete}  onIonInput={(e: any)=>OnChangeReps(e.target.value )}></IonInput>
      </IonItem>
    
    </IonListHeader>
  );
};

export default Set;


interface SetInputData{
  index : number;
  setData : SetData
  UpdateData : Function
  is_time_based: boolean
}

export const SetEdit: React.FC<SetInputData> = ({index, setData, UpdateData, is_time_based}) => {
	
    const [weightComplete, setWeightComplete] = useState(setData.weight)
    const [repsComplete, setRepsComplete] = useState(setData.reps_or_duration);
    
    useEffect(() => {
      UpdateData(index, {...setData, weight : weightComplete, reps_or_duration: repsComplete });
      }, [weightComplete,index, weightComplete, repsComplete ]); // Only re-run the effect if count changes

    
  

  

  const OnChangeReps = (val : string)=>{
   
        const re = /^[0-9\b]+$/;
    
        // if value is not blank, then test the regex
    
        if ( re.test(val)) {
            setRepsComplete(parseInt(val));
        }
        else{
            setRepsComplete(repsComplete);
        }
  }

  const OnChangeWeight = (val : string)=>{
   
    const re = /^[0-9\b]+$/;

    // if value is not blank, then test the regex
    
    if ( re.test(val)) {
        setWeightComplete(parseInt(val));
    }
    else{
        setWeightComplete(weightComplete);
    }
}

  return (
    
         
    <IonListHeader>
      <IonItem>
        <IonLabel>Weight</IonLabel>
      
      <IonInput inputmode="decimal" value={weightComplete} onIonInput={(e: any)=>OnChangeWeight(e.target.value )}></IonInput>
     
       <IonLabel>{is_time_based? "Duration seconds" : "Reps"}</IonLabel>
  
      <IonInput inputmode="numeric" pattern="[0-9]" min={0}  value={repsComplete}  onIonInput={(e: any)=>OnChangeReps(e.target.value )}></IonInput>
     
       
      </IonItem>
    </IonListHeader>
  );
};



