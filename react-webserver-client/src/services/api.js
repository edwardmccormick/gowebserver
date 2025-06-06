import axios from 'axios';

const API_URL = 'http://localhost:8080/people';

export const fetchPeople = async () => {
    try {
        const response = await axios.get(API_URL);
        return response.data;
    } catch (error) {
        console.error('Error fetching people:', error);
        throw error;
    }
};

export const fetchPersonById = async (id) => {
    try {
        const response = await axios.get(`${API_URL}/${id}`);
        return response.data;
    } catch (error) {
        console.error(`Error fetching person with ID ${id}:`, error);
        throw error;
    }
};