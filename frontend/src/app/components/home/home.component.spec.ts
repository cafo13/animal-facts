import { HttpClientModule } from '@angular/common/http'
import { ComponentFixture, TestBed } from '@angular/core/testing'
import { FactControlsComponent } from '../fact-controls/fact-controls.component'
import { FactComponent } from '../fact/fact.component'

import { HomeComponent } from './home.component'

describe('HomeComponent', () => {
    let component: HomeComponent
    let fixture: ComponentFixture<HomeComponent>

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [HttpClientModule],
            declarations: [HomeComponent, FactComponent, FactControlsComponent]
        }).compileComponents()

        fixture = TestBed.createComponent(HomeComponent)
        component = fixture.componentInstance
        fixture.detectChanges()
    })

    it('should create', () => {
        expect(component).toBeTruthy()
    })
})
