import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { WebTemplateNewComponent } from './web-template-new.component';

describe('WebTemplateNewComponent', () => {
  let component: WebTemplateNewComponent;
  let fixture: ComponentFixture<WebTemplateNewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ WebTemplateNewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WebTemplateNewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
