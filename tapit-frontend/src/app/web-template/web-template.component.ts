import { Component, OnInit } from '@angular/core';
import { WebTemplateService } from '../web-template.service';

@Component({
  selector: 'app-web-template',
  templateUrl: './web-template.component.html',
  styleUrls: ['./web-template.component.css']
})
export class WebTemplateComponent implements OnInit {

  constructor(private webTemplateService: WebTemplateService) { }

  ngOnInit() {
    this.webTemplateService.getWebTemplates();
  }

}
