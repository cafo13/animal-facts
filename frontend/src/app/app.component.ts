import { Component } from "@angular/core";
import { Router } from "@angular/router";
import { AuthService } from "./services/auth.service";

@Component({
  selector: "app-root",
  templateUrl: "./app.component.html",
  styleUrls: ["./app.component.scss"],
})
export class AppComponent {
  constructor(protected router: Router, protected authService: AuthService) {}

  navigateTo(componentName: string) {
    this.router.navigate([componentName], { skipLocationChange: true });
  }
}
