import React from 'react'
import { useState, useEffect } from 'react'

import './App.css'
import { Header } from './components/Header'
import { Fact } from './components/Fact'
import { FactService } from './services/fact.service'

const App = () => {
    const factService = new FactService()

    const [fact, setFact] = useState({
        Id: '',
        Text: '',
        Category: '',
        Source: '',
        Image: '',
    })
    useEffect(() => {
        async function getFact() {
            const response = await factService.getRandomFact()
            console.log(response)
            setFact(response.fact)
        }
        getFact()
    }, [])

    return (
        <div className="App">
            <Header />
            {fact && (
                <Fact
                    Id={fact.Id}
                    Text={fact.Text}
                    Category={fact.Category}
                    Source={fact.Source}
                    Image={fact.Image}
                />
            )}
        </div>
    )
}

export default App
