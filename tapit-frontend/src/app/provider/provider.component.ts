import { Component, OnInit } from '@angular/core';
import { ProviderService, TwilioProvider } from '../provider.service';


@Component({
  selector: 'app-provider',
  templateUrl: './provider.component.html',
  styleUrls: ['./provider.component.css']
})
export class ProviderComponent implements OnInit {

  currTwilioProvider: TwilioProvider = new TwilioProvider();

  submitProviders() {
    this.providerService.updateTwilioProvider(this.currTwilioProvider);
  }

  constructor(private providerService: ProviderService) { }

  ngOnInit() {
    this.providerService.getTwilioProviderObs().subscribe(currTwilio => {
      this.currTwilioProvider = currTwilio;
    });
  }

}
