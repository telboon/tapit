import { Component, OnInit } from '@angular/core';
import { TextTemplateService } from '../text-template.service';

@Component({
  selector: 'app-text-template',
  templateUrl: './text-template.component.html',
  styleUrls: ['./text-template.component.css']
})
export class TextTemplateComponent implements OnInit {

  constructor(private textTemplateService: TextTemplateService) { }

  ngOnInit() {
    this.textTemplateService.getTextTemplates();
  }

}
