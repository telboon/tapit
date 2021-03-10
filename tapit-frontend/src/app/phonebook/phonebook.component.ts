import { Component, OnInit } from '@angular/core';
import { PhonebookService } from '../phonebook.service';

@Component({
  selector: 'app-phonebook',
  templateUrl: './phonebook.component.html',
  styleUrls: ['./phonebook.component.css']
})
export class PhonebookComponent implements OnInit {

  constructor(private phonebookService: PhonebookService) { }

  ngOnInit() {
    this.phonebookService.getPhonebooks();
  }

}
