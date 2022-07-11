import { IonContent, IonHeader, IonPage, IonTitle, IonToolbar } from '@ionic/react';

import './WorkOutSessionTab.css';
import {ExerciseSession, ExerciseSessionData} from '../components/ExcerciseSession';

export interface WorkoutData
{
   DoneChanged:Function;
   RepsChanged:Function;
   WeightChanged: Function;
   exercises:ExerciseSessionData
}

const WorkOutSessionTab: React.FC<WorkoutData> = ({DoneChanged, RepsChanged, WeightChanged, exercises}) => {
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

export default WorkOutSessionTab;
