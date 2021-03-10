import { Component, OnInit } from '@angular/core';
import { TextTemplate, TextTemplateService } from '../text-template.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

@Component({
  selector: 'app-text-template-new',
  templateUrl: './text-template-new.component.html',
  styleUrls: ['./text-template-new.component.css']
})
export class TextTemplateNewComponent implements OnInit {

  newTextTemplate: TextTemplate = new TextTemplate();
  previewStr: string;
  id = 0;

  submitNewTextTemplate() {
    if (this.router.url === '/text-template/new') {
      this.textTemplateService.addTextTemplate(this.newTextTemplate);
    } else {
      this.editTextTemplate();
    }
  }

  updatePreview() {
    let tempStr = '';
    tempStr = this.newTextTemplate.templateStr;
    tempStr = tempStr.replace('{firstName}', 'John');
    tempStr = tempStr.replace('{lastName}', 'Smith');
    tempStr = tempStr.replace('{alias}', 'Johnny');
    tempStr = tempStr.replace('{phoneNumber}', '+6598765432');
    tempStr = tempStr.replace('{url}', 'https://www.example.com/eR2c1');

    this.previewStr = tempStr;
  }

  deleteTextTemplate() {
    this.textTemplateService.deleteTextTemplate(this.newTextTemplate);
  }

  editTextTemplate() {
    this.textTemplateService.editTextTemplate(this.newTextTemplate);
  }

  constructor(private textTemplateService: TextTemplateService, private router: Router, private route: ActivatedRoute) { }

  ngOnInit() {
    // if page is edit
    if (this.router.url !== '/text-template/new') {
      const idParam = 'id';
      this.route.params.subscribe( params => {
        this.id = parseInt(params[idParam], 10);
        this.textTemplateService.getTextTemplateObs(this.id).subscribe(currTT => {
          this.newTextTemplate = currTT;
          this.updatePreview();
        });
      });
    }
  }

}
