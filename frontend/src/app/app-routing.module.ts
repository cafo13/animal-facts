import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { HomeComponent } from "./components/home/home.component";
import { SignInComponent } from "./components/sign-in/sign-in.component";
import { AuthGuard } from "./guards/auth.guard";
import { AdminAreaComponent } from "./components/admin-area/admin-area.component";
import { ImprintComponent } from "./components/imprint/imprint.component";
import { PrivacyPolicyComponent } from "./components/privacy-policy/privacy-policy.component";

const routes: Routes = [
  { path: "home", component: HomeComponent },
  { path: "sign-in", component: SignInComponent },
  { path: "imprint", component: ImprintComponent },
  { path: "privacy-policy", component: PrivacyPolicyComponent },
  {
    path: "admin",
    component: AdminAreaComponent,
    canActivate: [AuthGuard],
  },
  { path: "", redirectTo: "/home", pathMatch: "full" },
];

@NgModule({
  imports: [CommonModule, RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
