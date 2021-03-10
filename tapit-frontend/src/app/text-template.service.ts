import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { NotificationService } from './notification.service';

export class TextTemplate {
  id: number;
  name: string;
  templateStr: string;
  createDate: Date;
}

export class TextTemplateNotification {
  resultType: string;
  text: string;
  payload: TextTemplate;
}

@Injectable({
  providedIn: 'root'
})

export class TextTemplateService {

  textTemplates: TextTemplate[] = [];

  templateUrl = '/api/text-template';

  httpOptions = {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
  };

  getTextTemplates() {
    this.http.get<TextTemplate[]>(this.templateUrl).subscribe(templates => {
      if (templates === null) {
        this.textTemplates = [];
      } else {
        this.textTemplates = templates;
      }
    });
  }

  getTextTemplateObs(id: number) {
    return this.http.get<TextTemplate>(this.templateUrl + '/' + id.toString());
  }

  addTextTemplate(newTextTemplate: TextTemplate) {
    this.http.post<TextTemplateNotification>(this.templateUrl, newTextTemplate, this.httpOptions).subscribe(templateNotification => {
      this.notificationService.addNotification(templateNotification.resultType, templateNotification.text);
      this.textTemplates.push(templateNotification.payload);
      if (templateNotification.payload !== null) {
        this.router.navigate(['/text-template']);
      }
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in creating template');
    });
  }

  deleteTextTemplate(textTemplate: TextTemplate) {
    this.http.delete<TextTemplateNotification>(this.templateUrl + '/' + textTemplate.id.toString(), this.httpOptions)
      .subscribe(templateNotification => {
        this.notificationService.addNotification(templateNotification.resultType, templateNotification.text);
        this.router.navigate(['/text-template']);
      },
      err => {
        this.notificationService.addNotification('failure', 'Error in deleting text template');
      });
  }

  editTextTemplate(textTemplate: TextTemplate) {
    this.http.put<TextTemplateNotification>(this.templateUrl + '/' + textTemplate.id.toString(), textTemplate, this.httpOptions)
      .subscribe(templateNotification => {
        this.notificationService.addNotification(templateNotification.resultType, templateNotification.text);
        if (templateNotification.payload !== null) {
          this.router.navigate(['/text-template']);
        }
      },
      err => {
        this.notificationService.addNotification('failure', 'Error in editing text template');
      });
  }

  constructor(private http: HttpClient, private router: Router, private notificationService: NotificationService) {
    this.getTextTemplates();
  }
}
