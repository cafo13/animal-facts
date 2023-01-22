import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss']
})
export class AppComponent {
    constructor(protected router: Router) {}

    navigateTo(componentName: string) {
        this.router.navigate([componentName], { skipLocationChange: true })
    }
}
