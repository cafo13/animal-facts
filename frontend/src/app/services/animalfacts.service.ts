import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";

import { Observable } from "rxjs";

import { environment } from "src/environments/environment";

export type Fact = {
  uuid: string;
  text: string;
  source: string;
  approved: boolean;
};

@Injectable({
  providedIn: "root",
})
export class AnimalfactsService {
  apiBaseDomain = environment.apiEndpoint;

  constructor(private http: HttpClient) {}

  getFact(id?: string): Observable<Fact> {
    if (id) {
      return this.http.get<Fact>(`${this.apiBaseDomain}/fact/${id}`);
    }
    return this.http.get<Fact>(`${this.apiBaseDomain}/fact`);
  }
}
