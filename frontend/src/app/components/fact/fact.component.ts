import { Component, Input } from "@angular/core";

@Component({
  selector: "app-fact",
  templateUrl: "./fact.component.html",
  styleUrls: ["./fact.component.scss"],
})
export class FactComponent {
  @Input() uuid: string = "";
  @Input() text: string = "";
  @Input() source: string = "";
  @Input() approved: boolean = false;
}
