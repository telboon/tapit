import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PhonebookComponent } from './phonebook.component';

describe('PhonebookComponent', () => {
  let component: PhonebookComponent;
  let fixture: ComponentFixture<PhonebookComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PhonebookComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PhonebookComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
