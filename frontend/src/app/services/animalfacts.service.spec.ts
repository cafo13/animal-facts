import { HttpClientModule } from "@angular/common/http";
import { TestBed } from "@angular/core/testing";

import { AnimalfactsService } from "./animalfacts.service";

describe("AnimalfactsService", () => {
  let service: AnimalfactsService;

  beforeEach(() => {
    TestBed.configureTestingModule({ imports: [HttpClientModule] });
    service = TestBed.inject(AnimalfactsService);
  });

  it("should be created", () => {
    expect(service).toBeTruthy();
  });
});
