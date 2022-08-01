import React, { useState } from 'react';

import { IonContent, IonSelect, IonSelectOption, IonItem, IonInput,IonIcon, IonHeader, IonPage, IonTitle, IonToolbar, IonGrid, IonRow, IonCol, IonButton, IonLabel } from '@ionic/react';
import { add, settings, share, person, arrowForwardCircle, arrowBackCircle, arrowUpCircle, logoVimeo, logoFacebook, logoInstagram, logoTwitter } from 'ionicons/icons';

	

interface SetAddProps {
  NewSet: Function;
}

const SetAdd: React.FC<SetAddProps> = ({NewSet}) => {
    
	

   
  
  return (
    <IonItem>
        
      
        <IonButton onClick={(e)=>NewSet()}> Add Set</IonButton>
    </IonItem>
  );
};

export default SetAdd;