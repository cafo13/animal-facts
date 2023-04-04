import { Component, OnDestroy, OnInit } from "@angular/core";
import { interval, Subscription } from "rxjs";

import { AnimalfactsService, Fact } from "src/app/services/animalfacts.service";

@Component({
  selector: "app-home",
  templateUrl: "./home.component.html",
  styleUrls: ["./home.component.scss"],
})
export class HomeComponent implements OnInit, OnDestroy {
  factRetentionSeconds = 30;
  progressSub: Subscription | undefined = undefined;
  progressValue: number = this.factRetentionSeconds;
  currentFact: Fact = { uuid: "", text: "", source: "", approved: false };

  constructor(private animalfactsService: AnimalfactsService) {}

  ngOnInit(): void {
    this.startFactTimer();
  }

  ngOnDestroy(): void {
    this.progressSub?.unsubscribe();
  }

  private loadFact() {
    this.animalfactsService.getFact().subscribe((data: Fact) => {
      this.currentFact = data;
    });
  }

  startFactTimer() {
    this.loadFact();
    this.progressSub = interval(1000).subscribe(
      () => (this.progressValue -= 1)
    );
  }

  resetFactTimer() {
    this.progressSub?.unsubscribe();
    this.progressValue = this.factRetentionSeconds;
    this.startFactTimer();
  }
}
