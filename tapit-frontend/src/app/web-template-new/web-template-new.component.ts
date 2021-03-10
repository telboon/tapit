import { Component, OnInit } from '@angular/core';
import { WebTemplateService, WebTemplate } from '../web-template.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

@Component({
  selector: 'app-web-template-new',
  templateUrl: './web-template-new.component.html',
  styleUrls: ['./web-template-new.component.css']
})
export class WebTemplateNewComponent implements OnInit {

  newWebTemplate: WebTemplate = new WebTemplate();
  id = 0;

  submitNewWebTemplate() {
    if (this.router.url === '/web-template/new') {
      this.webTemplateService.addWebTemplate(this.newWebTemplate);
    } else {
      this.editWebTemplate();
    }
  }

  deleteWebTemplate() {
    this.webTemplateService.deleteWebTemplate(this.newWebTemplate);
  }

  editWebTemplate() {
    this.webTemplateService.editWebTemplate(this.newWebTemplate);
  }

  constructor(private webTemplateService: WebTemplateService, private router: Router, private route: ActivatedRoute) { }

  ngOnInit() {
    // if page is edit
    if (this.router.url !== '/web-template/new') {
      const idParam = 'id';
      this.route.params.subscribe( params => {
        this.id = parseInt(params[idParam], 10);
        this.webTemplateService.getWebTemplateObs(this.id).subscribe(currWT => {
          this.newWebTemplate = currWT;
        });
      });
    }
  }

}
