import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';
import { NotificationService } from './notification.service';

export class Phonebook {
  id: number;
  name: string;
  size: number;
  createDate: Date;
  records: PhoneRecord[];
}

export class PhoneRecord {
  id: number;
  firstName: string;
  lastName: string;
  alias: string;
  phoneNumber: string;
}

export class PhonebookNotification {
  resultType: string;
  text: string;
  payload: Phonebook;
}

@Injectable({
  providedIn: 'root'
})
export class PhonebookService {

  phonebooks: Phonebook[] = [];

  phonebookUrl = '/api/phonebook';
  phonebookImportUrl = '/api/import-phonebook';

  httpOptions = {
    headers: new HttpHeaders({
      'Content-Type': 'application/json',
    }),
  };

  getPhonebooks() {
    this.http.get<Phonebook[]>(this.phonebookUrl).subscribe(phonebooks => {
      if (phonebooks === null) {
        this.phonebooks = [];
      } else {
        this.phonebooks = phonebooks;
      }
    });
  }

  getPhonebookObs(id: number): Observable<Phonebook> {
    return this.http.get<Phonebook>(this.phonebookUrl + '/' + id.toString());
  }

  addPhonebook(phonebook: Phonebook) {
    this.http.post<PhonebookNotification>(this.phonebookUrl, phonebook, this.httpOptions).subscribe(pbNotification => {
      this.notificationService.addNotification(pbNotification.resultType, pbNotification.text);
      this.phonebooks.push(pbNotification.payload);
      if (pbNotification.payload !== null) {
        this.router.navigate(['/phonebook']);
      }
    },
    err => {
      this.notificationService.addNotification('failure', 'Error in creating phonebook');
    });
  }

  editPhonebook(phonebook: Phonebook) {
    this.http.put<PhonebookNotification>(this.phonebookUrl + '/' + phonebook.id.toString(), phonebook, this.httpOptions)
      .subscribe(pbNotification => {
        this.notificationService.addNotification(pbNotification.resultType, pbNotification.text);
        if (pbNotification.payload !== null) {
          this.router.navigate(['/phonebook']);
        }
      },
      err => {
        this.notificationService.addNotification('failure', 'Error in editing phonebook');
      });
  }

  deletePhonebook(phonebook: Phonebook) {
    this.http.delete<PhonebookNotification>(this.phonebookUrl + '/' + phonebook.id.toString(), this.httpOptions)
      .subscribe(pbNotification => {
        this.notificationService.addNotification(pbNotification.resultType, pbNotification.text);
        this.router.navigate(['/phonebook']);
      },
      err => {
        this.notificationService.addNotification('failure', 'Error in deleting phonebook');
      });
    }

  uploadPhonebook(file: File): Observable<PhoneRecord[]> {
    const formData = new FormData();
    formData.append('phonebookFile', file);
    return this.http.post<PhoneRecord[]>(this.phonebookImportUrl, formData);
  }

  constructor(private http: HttpClient, private router: Router, private notificationService: NotificationService) {
    this.getPhonebooks();
  }
}
