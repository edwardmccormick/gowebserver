import React, { useEffect, useState } from 'react';
import { fetchPersonById } from '../services/api'; // Use named import

function PersonDetail({ id }) {
  const [person, setPerson] = useState(null);

  useEffect(() => {
    fetchPersonById(id).then((data) => setPerson(data));
  }, [id]);

  if (!person) return <div>Loading...</div>;

  return (
    <div>
      <h1>{person.Name}</h1>
      <p>{person.Motto}</p>
    </div>
  );
}

export default PersonDetail;