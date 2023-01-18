import { HttpClient } from '@angular/common/http'
import { Injectable } from '@angular/core'

import { Observable } from 'rxjs'

export type Fact = {
    ID: string
    text: string
    category: string
    source: string
}

@Injectable({
    providedIn: 'root'
})
export class AnimalfactsService {
    apiBaseDomain = 'http://localhost:8080/api/v1'

    constructor(private http: HttpClient) {}

    getFact(id?: string): Observable<Fact> {
        if (id) {
            return this.http.get<Fact>(`${this.apiBaseDomain}/fact/${id}`)
        }
        return this.http.get<Fact>(`${this.apiBaseDomain}/fact`)
    }
}
