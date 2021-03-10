import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TextTemplateComponent } from './text-template.component';

describe('TextTemplateComponent', () => {
  let component: TextTemplateComponent;
  let fixture: ComponentFixture<TextTemplateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TextTemplateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TextTemplateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
