import { Component, OnDestroy, OnInit } from '@angular/core'
import { interval, Subscription } from 'rxjs'

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
export class HomeComponent implements OnInit, OnDestroy {
    factRetentionSeconds = 30
    progressSub: Subscription | undefined = undefined
    progressValue: number = 0
    currentFact: Fact = { Id: '', Text: '', Category: '', Source: '' }

    constructor(private animalfactsService: AnimalfactsService) {}

    ngOnInit(): void {
        this.startFactTimer()
    }

    ngOnDestroy(): void {
        this.progressSub?.unsubscribe()
    }

    private loadFact() {
        this.animalfactsService.getFact().subscribe((data: FactApiResponse) => {
            this.currentFact = data.Fact
        })
    }

    startFactTimer() {
        this.loadFact()
        this.progressSub = interval(100).subscribe(() => {
            this.progressValue += 10 / this.factRetentionSeconds
        })
    }

    resetFactTimer() {
        this.progressSub?.unsubscribe()
        this.progressValue = 0
        this.startFactTimer()
    }
}
