import { TestBed } from '@angular/core/testing';

import { WebTemplateService } from './web-template.service';

describe('WebTemplateService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: WebTemplateService = TestBed.get(WebTemplateService);
    expect(service).toBeTruthy();
  });
});
