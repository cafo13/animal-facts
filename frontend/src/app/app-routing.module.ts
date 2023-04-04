import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { HomeComponent } from "./components/home/home.component";
import { SignInComponent } from "./components/sign-in/sign-in.component";

const routes: Routes = [
  { path: "home", component: HomeComponent },
  { path: "sign-in", component: SignInComponent },
  { path: "", redirectTo: "/home", pathMatch: "full" },
];

@NgModule({
  imports: [CommonModule, RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
