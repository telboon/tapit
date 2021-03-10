import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { NotificationService } from './notification.service';
import { Observable, of } from 'rxjs';

export class TwilioProvider {
  accountSID: string;
  authToken: string;
}

export class TwilioProviderNotification {
  resultType: string;
  text: string;
  payload: TwilioProvider;
}

@Injectable({
  providedIn: 'root'
})
export class ProviderService {

  twilioProviderSettings: TwilioProvider = new TwilioProvider();
  twilioUrl = '/api/provider/twilio';

  httpOptions = {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
  };

  providerEnums = [
                    {name: 'Twilio', tag: 'twilio'},
  ];

  getTwilioProvider() {
    this.http.get<TwilioProvider>(this.twilioUrl, this.httpOptions).subscribe(thisTwilio => {
      this.twilioProviderSettings = thisTwilio;
    });
  }

  getTwilioProviderObs(): Observable<TwilioProvider> {
    return this.http.get<TwilioProvider>(this.twilioUrl, this.httpOptions);
  }

  updateTwilioProvider(tProvider: TwilioProvider) {
    this.http.post<TwilioProviderNotification>(this.twilioUrl, tProvider, this.httpOptions).subscribe(usermessage => {
        this.notificationService.addNotification(usermessage.resultType, usermessage.text);
        this.twilioProviderSettings = usermessage.payload;
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in updating Twilio provider');
    });
  }

  constructor(private http: HttpClient, private notificationService: NotificationService) {
    this.getTwilioProvider();
  }
}
