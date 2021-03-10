import { Component, OnInit } from '@angular/core';
import { CampaignService, Campaign } from '../campaign.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { ProviderService } from '../provider.service';
import { PhonebookService } from '../phonebook.service';
import { TextTemplateService } from '../text-template.service';
import { WebTemplateService } from '../web-template.service';

@Component({
  selector: 'app-campaign-new',
  templateUrl: './campaign-new.component.html',
  styleUrls: ['./campaign-new.component.css']
})
export class CampaignNewComponent implements OnInit {

  constructor(
          private campaignService: CampaignService,
          private router: Router,
          private providerService: ProviderService,
          private phonebookService: PhonebookService,
          private textTemplateService: TextTemplateService,
          private webTemplateService: WebTemplateService,
  ) { }

  newCampaign: Campaign = new Campaign();

  templateStr = '';
  previewStr = '';

  submitNewCampaign() {
    this.campaignService.addCampaign(this.newCampaign);
  }

  submitNewCampaignRun() {
    this.campaignService.addCampaignRun(this.newCampaign);
  }

  updatePreviews() {
    if (this.newCampaign.textTemplateId !== 0 && this.newCampaign.phonebookId !== 0) {
      this.phonebookService.getPhonebookObs(this.newCampaign.phonebookId).subscribe(phonebook => {
        this.textTemplateService.getTextTemplateObs(this.newCampaign.textTemplateId).subscribe(textTemplate => {
          this.templateStr = textTemplate.templateStr;

          let tempStr = this.templateStr;
          tempStr = tempStr.replace('{firstName}', phonebook.records[0].firstName);
          tempStr = tempStr.replace('{lastName}', phonebook.records[0].lastName);
          tempStr = tempStr.replace('{alias}', phonebook.records[0].alias);
          tempStr = tempStr.replace('{phoneNumber}', phonebook.records[0].phoneNumber);
          tempStr = tempStr.replace('{url}', 'https://www.example.com/eR2c1');

          this.previewStr = tempStr;
        });
      });
    } else {
      this.templateStr = '';
      this.previewStr = '';
    }
  }

  ngOnInit() {
    this.newCampaign.textTemplateId = 0;
    this.newCampaign.webTemplateId = 0;
    this.newCampaign.phonebookId = 0;
  }

}
