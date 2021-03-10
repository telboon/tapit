import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PhonebookNewComponent } from './phonebook-new.component';

describe('PhonebookNewComponent', () => {
  let component: PhonebookNewComponent;
  let fixture: ComponentFixture<PhonebookNewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PhonebookNewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PhonebookNewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
