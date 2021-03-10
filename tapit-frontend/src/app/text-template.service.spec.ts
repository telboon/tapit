import { TestBed } from '@angular/core/testing';

import { TextTemplateService } from './text-template.service';

describe('TextTemplateService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: TextTemplateService = TestBed.get(TextTemplateService);
    expect(service).toBeTruthy();
  });
});
