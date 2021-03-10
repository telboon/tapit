import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TextTemplateNewComponent } from './text-template-new.component';

describe('TextTemplateNewComponent', () => {
  let component: TextTemplateNewComponent;
  let fixture: ComponentFixture<TextTemplateNewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TextTemplateNewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TextTemplateNewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
