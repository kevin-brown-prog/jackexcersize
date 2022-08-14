
import React, { useState } from 'react';

import { IonContent, IonSelect, IonSelectOption, IonItem, IonInput,IonIcon, IonHeader, IonPage, IonTitle, IonToolbar, IonGrid, IonRow, IonCol, IonButton, IonLabel } from '@ionic/react';
import { add, settings, share, person, arrowForwardCircle, arrowBackCircle, arrowUpCircle, logoVimeo, logoFacebook, logoInstagram, logoTwitter } from 'ionicons/icons';

	

interface ExcersiseAddProps {
  NewExercise: Function;
}

const ExcersiseAdd: React.FC<ExcersiseAddProps> = ({NewExercise}) => {
    
	

    const excercises = ['Bench Press', 'Saftey Bar Squat', 
  'Wide Grip Pullups', 'Narrow Grip Pullups', 'Trap Bar Dead Lifts']
   const [currentExercise, setCurrentExercise] = useState(excercises[0]);
   let currentSel = excercises[0]
  return (
         <IonRow>
          
          <IonCol>
          <IonSelect placeholder="Select Excersice" value={currentExercise}   onIonChange={(e)=>{setCurrentExercise(e.detail.value!)}}  >
          {excercises.map( (s,index)=><IonSelectOption key={index}   value={s}>{s}</IonSelectOption> )}
         

          </IonSelect>


          </IonCol>
          <IonCol> <IonButton onClick={(e)=>NewExercise(currentExercise)}> <IonIcon icon={add} /></IonButton></IonCol>
          
          </IonRow>
  );
};

export default ExcersiseAdd;