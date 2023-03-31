import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FactControlsComponent } from './fact-controls.component';

describe('FactControlsComponent', () => {
  let component: FactControlsComponent;
  let fixture: ComponentFixture<FactControlsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ FactControlsComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(FactControlsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
