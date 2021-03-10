import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  username: string;
  password: string;

  login() {
    this.authService.login(this.username, this.password);
  }

  routeRegister() {
    this.router.navigate(['/register']);
  }
  constructor(private authService: AuthService, private router: Router) { }

  ngOnInit() {
  }

}
