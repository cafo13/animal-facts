import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { interval, Subscription } from 'rxjs'

@Component({
    selector: 'app-fact-controls',
    templateUrl: './fact-controls.component.html',
    styleUrls: ['./fact-controls.component.scss']
})
export class FactControlsComponent implements OnInit {
    @Input() value: number = 0
    @Output() timerEnd = new EventEmitter<boolean>()

    valueSub: Subscription | undefined = undefined

    ngOnInit(): void {
        this.valueSub = interval(1000).subscribe(() => {
            if (this.value <= 0) {
                this.timerEnd.emit(true)
            }
        })
    }

    ngOnDestroy(): void {
        this.valueSub?.unsubscribe()
    }
}
