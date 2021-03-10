import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { CampaignService } from '../campaign.service';

@Component({
  selector: 'app-campaign',
  templateUrl: './campaign.component.html',
  styleUrls: ['./campaign.component.css']
})
export class CampaignComponent implements OnInit {

  constructor(private campaignService: CampaignService, private router: Router) { }

  ngOnInit() {
    this.campaignService.getCampaigns();
    const intervalId = setInterval(() => {
      this.campaignService.getCampaigns();
      if (!this.router.url.includes('/campaign')) {
        clearInterval(intervalId);
      }
    }, 2000);
  }

}
