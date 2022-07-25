import { Redirect, Route } from 'react-router-dom';
import {
  IonApp,
  IonIcon,
  IonLabel,
  IonRouterOutlet,
  IonTabBar,
  IonTabButton,
  IonTabs,
  setupIonicReact
} from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';



import { ellipse, square, triangle } from 'ionicons/icons';
import WorkOutSessionTab from './pages/WorkOutSessionTab';
import ExcerciseSessionCreate from './pages/ExcerciseSessionCreate';

import Tab3 from './pages/Tab3';

/* Core CSS required for Ionic components to work properly */
import '@ionic/react/css/core.css';

/* Basic CSS for apps built with Ionic */
import '@ionic/react/css/normalize.css';
import '@ionic/react/css/structure.css';
import '@ionic/react/css/typography.css';

/* Optional CSS utils that can be commented out */
import '@ionic/react/css/padding.css';
import '@ionic/react/css/float-elements.css';
import '@ionic/react/css/text-alignment.css';
import '@ionic/react/css/text-transformation.css';
import '@ionic/react/css/flex-utils.css';
import '@ionic/react/css/display.css';

/* Theme variables */
import './theme/variables.css';
import { ExerciseSessionData } from './components/ExcerciseSession';
import { serialize } from 'v8';
//https://www.pluralsight.com/guides/using-firebase-with-react-and-redux
setupIonicReact();

function DoneChange(id :string, is_done:boolean){
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: id, is_done: is_done })
};
fetch('/api/done_changed/'+encodeURIComponent(id)+"/"+encodeURIComponent(is_done), requestOptions)
    .then(response => response.json())
   .catch( (reason)=>{
 //   alert(reason);
   })




}
function RepsChange(id:string, reps_update:number){
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: id, reps: reps_update })
};
fetch('/api/reps_change/'+encodeURIComponent(id)+"/"+encodeURIComponent(reps_update), requestOptions)
    .then(response => response.json())
   .catch( (reason)=>{
 //   alert(reason);
   })
}
function WeightChanged(id:string, weight:number ){
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: id, weight: weight })
};
fetch('/api/weight_change/'+encodeURIComponent(id)+"/"+encodeURIComponent(weight), requestOptions)
    .then(response => response.json())
   .catch( (reason)=>{
 //   alert(reason);
   })
}

function AddExerciseSession(session : ExerciseSessionData, callback : Function ){
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(session)
};
fetch('/api/add_exercise_session', requestOptions)
    .then(response => response.json()).then(obj=>{

      callback(obj.ID)

    })
  
   .catch( (reason)=>{
 //   alert(reason);
   })
}


function DeleteExerciseSession(sessionID : string, callback : Function ){
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ID:sessionID})
};
fetch('/api/delete_exercise_session/' + encodeURIComponent(sessionID), requestOptions)
    .then(response => response.json()).then(obj=>{

      callback(obj.ID)

    })
  
   .catch( (reason)=>{
//    alert(reason);
   })
}

var set1 = { set_id : "0",
weight : 135,
reps_or_duration : 5,
done : false,
is_time_based: false
};

var set2 = { set_id : "1",
  weight : 135,
  reps_or_duration : 5,
  done : false,
  is_time_based: false,
 };

 var set3 = { set_id : "2",
  weight : 135,
  reps_or_duration : 5,
  done : false,
  is_time_based: false,
  };

  var Exercise1 = {
    name:"Bench Press",
    sets :[set1, set2,set3]

  }


   set1 = { set_id : "3",
    weight : 135,
    reps_or_duration : 5,
    done : false,
    is_time_based: false
    };
    
     set2 = { set_id : "4",
      weight : 135,
      reps_or_duration : 5,
      done : false,
      is_time_based: false,
     };
    
      set3 = { set_id : "5",
      weight : 135,
      reps_or_duration : 5,
      done : false,
      is_time_based: false,
      };



  var Exercise2 = {
    name:"Squats",
    sets:[set1, set2, set3]
  }

  let ExerciseSession ={
    name:"Monday session",
    date:new Date(),
    exercises : [Exercise1, Exercise2]
  }




const App: React.FC = () => (
  <IonApp>
    <IonReactRouter>
      <IonTabs>
        <IonRouterOutlet>
          <Route exact path="/tab1">
            <WorkOutSessionTab DoneChanged={DoneChange} RepsChanged={RepsChange} WeightChanged={WeightChanged} exercises={ExerciseSession} />
          </Route>
          <Route exact path="/tab2">
            <ExcerciseSessionCreate />
          </Route>
          <Route path="/tab3">
            <Tab3 />
          </Route>
          <Route exact path="/">
            <Redirect to="/tab1" />
          </Route>
        </IonRouterOutlet>
        <IonTabBar slot="bottom">
          <IonTabButton tab="tab1" href="/tab1">
            <IonIcon icon={triangle} />
            <IonLabel>Tab 1</IonLabel>
          </IonTabButton>
          <IonTabButton tab="tab2" href="/tab2">
            <IonIcon icon={ellipse} />
            <IonLabel>Tab 2</IonLabel>
          </IonTabButton>
          <IonTabButton tab="tab3" href="/tab3">
            <IonIcon icon={square} />
            <IonLabel>Tab 3</IonLabel>
          </IonTabButton>
        </IonTabBar>
      </IonTabs>
    </IonReactRouter>
  </IonApp>
);

export default App;
