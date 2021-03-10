import { Component, OnInit } from '@angular/core';
import { GlobalSettings, GlobalSettingsService } from '../global-settings.service';

@Component({
  selector: 'app-global-settings',
  templateUrl: './global-settings.component.html',
  styleUrls: ['./global-settings.component.css']
})
export class GlobalSettingsComponent implements OnInit {

  tempSettings = new GlobalSettings();
  displaySettings;

  updateGlobalSettings() {
    this.tempSettings.secretRegistrationCode = this.displaySettings.secretRegistrationCode;
    this.tempSettings.webTemplatePrefix = this.displaySettings.webTemplatePrefix;
    this.tempSettings.webFrontPlaceholder = this.displaySettings.webFrontPlaceholder;
    this.tempSettings.threadsPerCampaign = parseInt(this.displaySettings.threadsPerCampaign, 10) + 0;
    this.tempSettings.bcryptCost = parseInt(this.displaySettings.bcryptCost, 10);
    this.tempSettings.maxRequestRetries = parseInt(this.displaySettings.maxRequestRetries, 10);
    this.tempSettings.waitBeforeRetry = parseInt(this.displaySettings.waitBeforeRetry, 10);

    this.globalSettingsService.updateGlobalSettings(this.tempSettings);
  }

  constructor(private globalSettingsService: GlobalSettingsService) { }

  ngOnInit() {
    this.globalSettingsService.getGlobalSettingsObs().subscribe(settings => {
      console.log(settings);
      this.displaySettings = settings;
    });
  }

}
