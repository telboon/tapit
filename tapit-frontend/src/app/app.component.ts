import { Component } from '@angular/core';
import { RouterModule, Routes, Router } from '@angular/router';
import { AuthService } from './auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'tapit-frontend';
  navlinks = [
    {
      link: '/campaign',
      name: 'Campaigns',
      loginOnly: true,
    },
    {
      link: '/phonebook',
      name: 'Phonebook',
      loginOnly: true,
    },
    {
      link: '/text-template',
      name: 'Text Templates',
      loginOnly: true,
    },
    {
      link: '/web-template',
      name: 'Web Templates',
      loginOnly: true,
    },
  ];
  constructor( private router: Router, private authService: AuthService) {
    authService.getUser();
  }
}
