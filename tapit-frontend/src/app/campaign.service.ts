import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { NotificationService } from './notification.service';

export class Campaign {
  id: number;
  name: string;
  fromNumber: string;
  size: number;
  currentStatus: string;
  createDate: Date;
  phonebookId: number;
  textTemplateId: number;
  webTemplateId: number;
  providerTag: string;
  jobs: Job[];
}

export class Job {
  id: number;
  currentStatus: string;
  webStatus: string;
  timeSent: Date;
  fromNum: string;
  toNum: string;
  fullUrl: string;
}

export class CampaignNotification {
  resultType: string;
  text: string;
  payload: Campaign;
}

@Injectable({
  providedIn: 'root'
})

export class CampaignService {

  campaigns: Campaign[] = [];

  campaignUrl = '/api/campaign';

  httpOptions = {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
  };

  getCampaigns() {
    this.http.get<Campaign[]>(this.campaignUrl).subscribe(campaigns => {
      if (campaigns === null) {
        this.campaigns = [];
      } else {
        this.campaigns = campaigns;
      }
    });
  }

  getCampaignObs(id: number): Observable<Campaign> {
    return this.http.get<Campaign>(this.campaignUrl + '/' + id.toString());
  }

  addCampaign(newCampaign: Campaign) {
    this.http.post<CampaignNotification>(this.campaignUrl, newCampaign, this.httpOptions).subscribe(campaignNotification => {
      this.notificationService.addNotification(campaignNotification.resultType, campaignNotification.text);
      this.campaigns.push(campaignNotification.payload);
      if (campaignNotification.payload !== null) {
        this.router.navigate(['/campaign']);
      }
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in creating template');
    });
  }

  addCampaignRun(newCampaign: Campaign) {
    this.http.post<CampaignNotification>(this.campaignUrl, newCampaign, this.httpOptions).subscribe(campaignNotification => {
      this.notificationService.addNotification(campaignNotification.resultType, campaignNotification.text);
      this.campaigns.push(campaignNotification.payload);
      if (campaignNotification.payload !== null) {
        this.startCampaign(campaignNotification.payload).subscribe();
        this.router.navigate(['/campaign']);
      }
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in creating template');
    });
  }

  deleteCampaign(campaign: Campaign) {
    this.http.delete<CampaignNotification>(this.campaignUrl + '/' + campaign.id.toString(), this.httpOptions)
      .subscribe(campaignNotification => {
        this.notificationService.addNotification(campaignNotification.resultType, campaignNotification.text);
        this.router.navigate(['/campaign']);
      },
      err => {
        this.notificationService.addNotification('failure', 'Error in deleting campaign');
      });
  }

  startCampaign(campaign: Campaign) {
    return this.http.get<CampaignNotification>(this.campaignUrl + '/' + campaign.id.toString() + '/' + 'start');
  }

  pauseCampaign(campaign: Campaign) {
    return this.http.get<CampaignNotification>(this.campaignUrl + '/' + campaign.id.toString() + '/' + 'pause');
  }

  constructor(private http: HttpClient, private router: Router, private notificationService: NotificationService) {
    this.campaigns = [];
    this.getCampaigns();
  }
}
