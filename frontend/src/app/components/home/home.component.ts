import { Component } from '@angular/core'

import {
    AnimalfactsService,
    Fact,
    FactApiResponse
} from 'src/app/services/animalfacts.service'

@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.scss']
})
export class HomeComponent {
    fact: string = ''

    constructor(animalfactsService: AnimalfactsService) {
        animalfactsService.getFact().subscribe((data: FactApiResponse) => {
            this.fact = data.Fact.Text
        })
    }
}
