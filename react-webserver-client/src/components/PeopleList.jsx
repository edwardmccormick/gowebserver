import React, { useEffect, useState } from 'react';
import axios from 'axios';

const PeopleList = () => {
    const [people, setPeople] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchPeople = async () => {
            try {
                const response = await axios.get('http://localhost:8080/people');
                setPeople(response.data);
            } catch (err) {
                setError(err);
            } finally {
                setLoading(false);
            }
        };

        fetchPeople();
    }, []);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error fetching people: {error.message}</p>;

    return (
        <div>
            <h1>People List</h1>
            <ul>
                {people.map(person => (
                    <li key={person.id}>
                        {person.name} - {person.motto}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default PeopleList;