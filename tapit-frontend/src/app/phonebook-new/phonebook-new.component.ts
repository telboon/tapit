import { Component, OnInit } from '@angular/core';
import { PhonebookService, Phonebook, PhoneRecord } from '../phonebook.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

@Component({
  selector: 'app-phonebook-new',
  templateUrl: './phonebook-new.component.html',
  styleUrls: ['./phonebook-new.component.css']
})
export class PhonebookNewComponent implements OnInit {

  constructor(private phonebookService: PhonebookService, private router: Router, private route: ActivatedRoute) { }

  id = 0;

  newPhonebook: Phonebook = new Phonebook();
  newPhoneRecords: PhoneRecord[] = [];
  additionalRecord: PhoneRecord = new PhoneRecord();

  insertAdditionalRecord() {
    this.newPhoneRecords = this.newPhoneRecords.concat(this.additionalRecord);
    this.additionalRecord = new PhoneRecord();
    this.additionalRecord.phoneNumber = '';
  }

  importPhoneRecords(files: FileList) {
    this.phonebookService.uploadPhonebook(files.item(0)).subscribe(data => {
      this.newPhoneRecords = this.newPhoneRecords.concat(data);
    });
  }

  submitNewPhonebook() {
    if (this.router.url === '/phonebook/new') {
      if (this.additionalRecord.phoneNumber !== '') {
        this.insertAdditionalRecord();
      }
      this.newPhonebook.records = this.newPhoneRecords;
      this.phonebookService.addPhonebook(this.newPhonebook);
    } else {
      this.editPhonebook();
    }
  }

  deletePhonebook() {
    this.phonebookService.deletePhonebook(this.newPhonebook);
  }

  editPhonebook() {
    this.newPhonebook.records = this.newPhoneRecords;
    this.phonebookService.editPhonebook(this.newPhonebook);
  }

  ngOnInit() {
    this.additionalRecord = new PhoneRecord();
    this.additionalRecord.phoneNumber = '';

    // if page is edit
    if (this.router.url !== '/phonebook/new') {
      const idParam = 'id';
      this.route.params.subscribe( params => {
        this.id = parseInt(params[idParam], 10);
        this.phonebookService.getPhonebookObs(this.id).subscribe(currPb => {
          this.newPhonebook = currPb;
          this.newPhoneRecords = this.newPhonebook.records;
        });
      });
    }
  }

}
