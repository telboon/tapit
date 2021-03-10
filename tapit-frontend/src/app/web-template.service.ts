import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { NotificationService } from './notification.service';

export class WebTemplate {
  id: number;
  name: string;
  templateType: string; // enum redirect, harvester
  redirectAgent: string;
  redirectNegAgent: string;
  redirectPlaceholderHtml: string;
  redirectUrl: string;
  harvesterBeforeHtml: string;
  harvesterAfterHtml: string;
  createDate: Date;
}

export class WebTemplateNotification {
  resultType: string;
  text: string;
  payload: WebTemplate;
}

@Injectable({
  providedIn: 'root'
})
export class WebTemplateService {

  wTemplateEnum = [
                    {name: 'Redirect Based On User Agent', tag: 'redirect'},
                    {name: 'Credentials Harvesting', tag: 'harvester'},
  ];

  templateUrl = '/api/web-template';
  httpOptions = {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
  };

  webTemplates: WebTemplate[] = [];

  getWebTemplates() {
    this.http.get<WebTemplate[]>(this.templateUrl).subscribe(templates => {
      if (templates === null) {
        this.webTemplates = [];
      } else {
        this.webTemplates = templates;
      }
    });
  }
  addWebTemplate(newWebTemplate: WebTemplate) {
    this.http.post<WebTemplateNotification>(this.templateUrl, newWebTemplate, this.httpOptions).subscribe(templateNotification => {
      this.notificationService.addNotification(templateNotification.resultType, templateNotification.text);
      this.webTemplates.push(templateNotification.payload);
      if (templateNotification.payload !== null) {
        this.router.navigate(['/web-template']);
      }
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in creating template');
    });
  }

  deleteWebTemplate(webTemplate: WebTemplate) {
    this.http.delete<WebTemplateNotification>(this.templateUrl + '/' + webTemplate.id.toString(), this.httpOptions)
      .subscribe(templateNotification => {
        this.notificationService.addNotification(templateNotification.resultType, templateNotification.text);
        this.router.navigate(['/web-template']);
      },
      err => {
        this.notificationService.addNotification('failure', 'Error in deleting web template');
      });
  }

  editWebTemplate(webTemplate: WebTemplate) {
    this.http.put<WebTemplateNotification>(this.templateUrl + '/' + webTemplate.id.toString(), webTemplate, this.httpOptions)
      .subscribe(templateNotification => {
        this.notificationService.addNotification(templateNotification.resultType, templateNotification.text);
        if (templateNotification.payload !== null) {
          this.router.navigate(['/web-template']);
        }
      },
      err => {
        this.notificationService.addNotification('failure', 'Error in editing web template');
      });
  }

  getWebTemplateObs(id: number) {
    return this.http.get<WebTemplate>(this.templateUrl + '/' + id.toString());
  }

  constructor(private http: HttpClient, private router: Router, private notificationService: NotificationService) {
    this.getWebTemplates();
  }
}
