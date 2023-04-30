import { NgModule } from "@angular/core";
import { BrowserModule } from "@angular/platform-browser";
import { HttpClientModule } from "@angular/common/http";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { AngularFireModule } from "@angular/fire/compat";
import { AngularFireAuthModule } from "@angular/fire/compat/auth";

import { MdbAccordionModule } from "mdb-angular-ui-kit/accordion";
import { MdbCarouselModule } from "mdb-angular-ui-kit/carousel";
import { MdbCheckboxModule } from "mdb-angular-ui-kit/checkbox";
import { MdbCollapseModule } from "mdb-angular-ui-kit/collapse";
import { MdbDropdownModule } from "mdb-angular-ui-kit/dropdown";
import { MdbFormsModule } from "mdb-angular-ui-kit/forms";
import { MdbModalModule } from "mdb-angular-ui-kit/modal";
import { MdbPopoverModule } from "mdb-angular-ui-kit/popover";
import { MdbRadioModule } from "mdb-angular-ui-kit/radio";
import { MdbRangeModule } from "mdb-angular-ui-kit/range";
import { MdbRippleModule } from "mdb-angular-ui-kit/ripple";
import { MdbScrollspyModule } from "mdb-angular-ui-kit/scrollspy";
import { MdbTabsModule } from "mdb-angular-ui-kit/tabs";
import { MdbTooltipModule } from "mdb-angular-ui-kit/tooltip";
import { MdbValidationModule } from "mdb-angular-ui-kit/validation";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { HomeComponent } from "./components/home/home.component";
import { FactComponent } from "./components/fact/fact.component";
import { FactControlsComponent } from "./components/fact-controls/fact-controls.component";

import { environment } from "../environments/environment";
import { SignInComponent } from "./components/sign-in/sign-in.component";
import { AuthService } from "./services/auth.service";
import { AdminAreaComponent } from "./components/admin-area/admin-area.component";
import { MenuBarComponent } from "./components/menu-bar/menu-bar.component";
import { FooterComponent } from "./components/footer/footer.component";
import { ImprintComponent } from "./components/imprint/imprint.component";
import { PrivacyPolicyComponent } from "./components/privacy-policy/privacy-policy.component";
import { AuthGuard } from "./guards/auth.guard";

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    FactComponent,
    FactControlsComponent,
    SignInComponent,
    AdminAreaComponent,
    MenuBarComponent,
    FooterComponent,
    ImprintComponent,
    PrivacyPolicyComponent,
  ],
  imports: [
    AngularFireModule.initializeApp(environment.firebaseConfig),
    AngularFireAuthModule,
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MdbAccordionModule,
    MdbCarouselModule,
    MdbCheckboxModule,
    MdbCollapseModule,
    MdbDropdownModule,
    MdbFormsModule,
    MdbModalModule,
    MdbPopoverModule,
    MdbRadioModule,
    MdbRangeModule,
    MdbRippleModule,
    MdbScrollspyModule,
    MdbTabsModule,
    MdbTooltipModule,
    MdbValidationModule,
  ],
  providers: [AuthService, AuthGuard],
  bootstrap: [AppComponent],
})
export class AppModule {}
