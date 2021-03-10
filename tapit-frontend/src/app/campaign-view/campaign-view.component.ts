import { Component, OnInit } from '@angular/core';
import { CampaignService, Campaign, Job, CampaignNotification } from '../campaign.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { NotificationService } from '../notification.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-campaign-view',
  templateUrl: './campaign-view.component.html',
  styleUrls: ['./campaign-view.component.css']
})
export class CampaignViewComponent implements OnInit {

  currCampaign: Campaign = new Campaign();

  id = 0;

  constructor(
        private campaignService: CampaignService,
        private router: Router,
        private route: ActivatedRoute,
        private notificationService: NotificationService,
        private http: HttpClient
  ) { }

  downloadVisits(id: string) {
    fetch('/api/jobs/' + id + '/visits')
      .then(resp => resp.blob())
      .then(blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        // the filename you want
        a.download = 'visits-' + id + '.csv';
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        this.notificationService.addNotification('success', 'Successfully retrived web visits');
      })
      .catch(() => this.notificationService.addNotification('failure', 'Failed to download web visits'));
  }

  startCampaign() {
    this.campaignService.startCampaign(this.currCampaign).subscribe(campaignNotification => {
      this.notificationService.addNotification(campaignNotification.resultType, campaignNotification.text);
      this.campaignService.getCampaignObs(this.id).subscribe(campaign => {
        this.currCampaign = campaign;
      });
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in starting campaign');
    });
  }

  pauseCampaign() {
    this.campaignService.pauseCampaign(this.currCampaign).subscribe(campaignNotification => {
      this.notificationService.addNotification(campaignNotification.resultType, campaignNotification.text);
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in pausing campaign');
    });
  }

  deleteCampaign() {
    this.campaignService.deleteCampaign(this.currCampaign);
  }

  updateThisCampaign() {
      this.campaignService.getCampaignObs(this.id).subscribe(campaign => {
        this.currCampaign = campaign;
      });
  }

  ngOnInit() {
      const idParam = 'id';
      this.route.params.subscribe( params => {
        this.id = parseInt(params[idParam], 10);
      });
      this.updateThisCampaign();
      const intervalId = setInterval(() => {
        this.updateThisCampaign();
        if (!this.router.url.includes('/campaign')) {
          clearInterval(intervalId);
        }
      }, 2000);
  }

}
