import React from 'react'

export interface IFact {
    Id: string
    Text: string
    Category: string
    Source: string
    Image: string
}

export const Fact = (fact: IFact) => {
    console.log('Fact', fact)
    return (
        <div className="fact">
            <p>{fact.Text}</p>
            <a href={fact.Source}>Source</a>
        </div>
    )
}
