import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { NotificationService } from './notification.service';
import { Observable, of } from 'rxjs';

export class User {
  username: string;
  password: string;
  name: string;
  email: string;
  secretCode: string;
}

export class UserNotification {
  resultType: string;
  text: string;
  payload: User;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  currUser = new User();
  loggedin = false;         // change this for testing
  loginUrl = 'api/login';
  logoutUrl = 'api/logout';
  registerUrl = 'api/register';
  myselfUrl = 'api/myself';

  httpOptions = {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
  };

  login(username: string, password: string) {
    this.currUser.username = username;
    this.currUser.password = password;
    this.http.post<UserNotification>(this.loginUrl, this.currUser, this.httpOptions).subscribe(usermessage => {
      if (usermessage.payload !== null) {
        this.loggedin = true;

        // update user
        this.currUser.username = usermessage.payload.username;
        this.currUser.email = usermessage.payload.email;
        this.currUser.name = usermessage.payload.name;

        this.notificationService.addNotification(usermessage.resultType, usermessage.text);
        this.router.navigate(['/campaign']);
      } else {
        this.notificationService.addNotification(usermessage.resultType, usermessage.text);
      }
    },
    err => {
        this.notificationService.addNotification('failure', 'Error in logging in');
    });
    this.currUser.password = '';
  }

  register(username: string, password: string, email: string, name: string, secretCode: string) {
    this.currUser.username = username;
    this.currUser.password = password;
    this.currUser.email = email;
    this.currUser.name = name;
    this.currUser.secretCode = secretCode;

    this.http.post<UserNotification>(this.registerUrl, this.currUser, this.httpOptions).subscribe(usermessage => {
      if (usermessage.payload !== null) {
        this.loggedin = true;
        this.notificationService.addNotification(usermessage.resultType, usermessage.text);
        this.router.navigate(['/campaign']);

        // update user
        this.currUser.username = usermessage.payload.username;
        this.currUser.email = usermessage.payload.email;
        this.currUser.name = usermessage.payload.name;
      } else {
        this.notificationService.addNotification(usermessage.resultType, usermessage.text);
      }
    });

    this.currUser.secretCode = '';
  }

  logout() {
    this.http.post<UserNotification>(this.logoutUrl, '', this.httpOptions).subscribe(usermessage => {
      this.notificationService.addNotification(usermessage.resultType, usermessage.text);
      this.loggedin = false;
      this.currUser = new User();
      this.router.navigate(['/']);
    });
  }

  getUser(): User {
    this.http.get<User>(this.myselfUrl, this.httpOptions).subscribe(thisUser => {
      this.currUser = thisUser;
      if (this.currUser.username !== '') {
        this.loggedin = true;
      } else {
        this.router.navigate(['/']);
      }
      // separate one to redirect main to campaign dashboard
      if (this.router.url === '/' || this.router.url === '') {
        this.router.navigate(['/campaign']);
      }
    },
    err => {
      this.router.navigate(['/']);
    });
    return this.currUser;
  }

  getUserObs(): Observable<User> {
    return this.http.get<User>(this.myselfUrl, this.httpOptions);
  }

  updateUser(user: User) {
    this.currUser = user;
    this.http.put<UserNotification>(this.myselfUrl, this.currUser, this.httpOptions).subscribe(usermessage => {
        this.notificationService.addNotification(usermessage.resultType, usermessage.text);
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in updating profile');
    });
    this.currUser.password = '';
  }

  constructor(private http: HttpClient, private router: Router, private notificationService: NotificationService) { }
}
