import { Component, Input } from "@angular/core";

@Component({
  selector: "app-fact",
  templateUrl: "./fact.component.html",
  styleUrls: ["./fact.component.scss"],
})
export class FactComponent {
  @Input() id: string = "";
  @Input() text: string = "";
  @Input() source: string = "";
}
