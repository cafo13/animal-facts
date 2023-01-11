import { HttpClient } from '@angular/common/http'
import { Injectable } from '@angular/core'

import { Observable } from 'rxjs'

export type Fact = {
    Id: string
    Text: string
    Category: string
    Source: string
}

export type FactApiResponse = {
    Fact: Fact
}

@Injectable({
    providedIn: 'root'
})
export class AnimalfactsService {
    apiBaseDomain = 'http://localhost:8080/api/v1'

    constructor(private http: HttpClient) {}

    getFact(id?: string): Observable<FactApiResponse> {
        if (id) {
            return this.http.get<FactApiResponse>(
                `${this.apiBaseDomain}/fact/${id}`
            )
        }
        return this.http.get<FactApiResponse>(`${this.apiBaseDomain}/fact`)
    }
}
