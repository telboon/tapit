import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {
  username = '';
  password = '';
  email = '';
  name = '';
  secretCode = '';

  register() {
    this.authService.register(this.username, this.password, this.email, this.name, this.secretCode);
  }
  constructor(private authService: AuthService) { }

  ngOnInit() {
  }

}
