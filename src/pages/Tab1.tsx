import { IonContent, IonHeader, IonPage, IonTitle, IonToolbar } from '@ionic/react';

import './Tab1.css';
import {ExerciseSession, ExerciseSessionData} from '../components/ExcerciseSession';

export interface WorkoutData
{
   DoneChanged:Function;
   RepsChanged:Function;
   WeightChanged: Function;
   exercises:ExerciseSessionData
}

const Tab1: React.FC<WorkoutData> = ({DoneChanged, RepsChanged, WeightChanged, exercises}) => {
  return (
    <IonPage>
      <IonHeader>
          <IonTitle>Todays Exercise</IonTitle>
      </IonHeader>
      <IonContent fullscreen>
        
        <ExerciseSession exerciseSession={exercises}  RepsChanged={RepsChanged} WeightChanged={WeightChanged} DoneChanged={DoneChanged} />
      </IonContent>
    </IonPage>
  );
};

export default Tab1;
