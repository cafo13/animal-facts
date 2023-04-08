import { Component, OnInit } from "@angular/core";
import { AuthService } from "src/app/services/auth.service";

@Component({
  selector: "app-menu-bar",
  templateUrl: "./menu-bar.component.html",
  styleUrls: ["./menu-bar.component.scss"],
})
export class MenuBarComponent implements OnInit {
  constructor(protected auth: AuthService) {}

  ngOnInit() {}
}
