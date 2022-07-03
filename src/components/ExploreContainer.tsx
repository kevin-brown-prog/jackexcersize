import './ExploreContainer.css';
import React, { useState } from 'react';

import {  IonCheckbox, IonList, IonItem, IonLabel } from '@ionic/react';
	

interface ContainerProps {
  name: string;
}

const ExploreContainer: React.FC<ContainerProps> = ({ name }) => {
	const [checked, setChecked] = useState(false);
  return (
    <div className="container">
       <IonList>
      <IonItem>
         
        <IonCheckbox checked={checked} onIonChange={e => setChecked(e.detail.checked)} />
		<IonLabel>{name} </IonLabel>
      </IonItem>
      <IonItem>
        <IonLabel>Mega Man X</IonLabel>
      </IonItem>
      <IonItem>
        <IonLabel>The Legend of Zelda</IonLabel>
      </IonItem>
      <IonItem>
        <IonLabel>Pac-Man</IonLabel>
      </IonItem>
      <IonItem>
        <IonLabel>Super Mario World</IonLabel>
      </IonItem>
    </IonList>
    </div>
  );
};

export default ExploreContainer;
