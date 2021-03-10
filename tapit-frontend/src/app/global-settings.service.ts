import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { NotificationService } from './notification.service';
import { Observable, of } from 'rxjs';

export class GlobalSettings {
    secretRegistrationCode: string;
    threadsPerCampaign: number;
    bcryptCost: number;
    maxRequestRetries: number;
    waitBeforeRetry: number;
    webTemplatePrefix: string;
    webFrontPlaceholder: string;
}

export class GlobalSettingsNotification {
  resultType: string;
  text: string;
  payload: GlobalSettings;
}

@Injectable({
  providedIn: 'root'
})
export class GlobalSettingsService {
  globalSettings = new GlobalSettings();

  globalSettingsUrl = 'api/globalsettings';

  httpOptions = {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
  };

  getGlobalSettings() {
    this.http.get<GlobalSettings>(this.globalSettingsUrl, this.httpOptions).subscribe(globalSettings => {
      this.globalSettings = globalSettings;
    });
  }

  getGlobalSettingsObs(): Observable<GlobalSettings> {
    return this.http.get<GlobalSettings>(this.globalSettingsUrl, this.httpOptions);
  }

  updateGlobalSettings(globalSettings: GlobalSettings) {
    this.http.put<GlobalSettingsNotification>(this.globalSettingsUrl, globalSettings, this.httpOptions).subscribe(settingsMessage => {
      this.notificationService.addNotification(settingsMessage.resultType, settingsMessage.text);
      this.globalSettings = settingsMessage.payload;
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in updating settings');
    });
  }

  constructor(private http: HttpClient, private router: Router, private notificationService: NotificationService) {
    this.getGlobalSettings();
  }
}
