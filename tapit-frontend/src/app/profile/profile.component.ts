import { Component, OnInit } from '@angular/core';
import { AuthService, User } from '../auth.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {

  constructor(private authService: AuthService) { }
  currUser: User;

  updateUser() {
    this.authService.updateUser(this.currUser);
  }

  ngOnInit() {
    this.authService.getUserObs().subscribe(user => {
      this.currUser = JSON.parse(JSON.stringify(user));
    });
  }

}
