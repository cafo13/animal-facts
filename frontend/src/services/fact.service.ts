import axios from 'axios'

export class FactService {
    public async getRandomFact(): Promise<any> {
        const response = await axios.get('http://localhost:8080/fact')
        return response.data
    }

    public async getFactById(id: string) {
        const response = await axios.get(`http://localhost:8080/fact/${id}`)
        return response.data
    }
}
