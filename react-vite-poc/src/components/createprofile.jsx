import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import DetailsSelections from './detailsselections';
import details from '../../../details.json';

export function CreateProfile({setLoggedInUser, loggedInUser}) {
  const [formData, setFormData] = useState({
    name: '',
    motto: '',
    latitude: '',
    longitude: '',
    profile: '',
    details: Object.keys(details).reduce((acc, key) => ({ ...acc, [key]: 0 }), {}),
  });

  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleDetailsChange = (key, value) => {
    setFormData((prev) => ({
      ...prev,
      details: { ...prev.details, [key]: parseInt(value, 10) },
    }));
  };

  const handleSubmit = async () => {
    const payload = {
      id: loggedInUser,
      age: parseInt(formData.age),
      name: formData.name,
      motto: formData.motto,
      lat: parseFloat(formData.latitude),
      long: parseFloat(formData.longitude),
      profile: formData.profile,
      description: formData.description,
      details: formData.details,
    };

    try {
      const response = await fetch('http://localhost:8080/people', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        const data = await response.json();
        setLoggedInUser(data); // Update the logged-in user in App.jsx
        console.log('Profile created successfully:', data);
        alert('Completed your profile! Nice');
      } else {
        alert('Ruh roh Shaggy');
      }
      } catch (error) {
        console.error('Error during profile creation:', error);
        alert('An error occurred during profile creation.');
      }
  };

 return (
    <>
      <div className="d-flex justify-content-center align-items-center flex-wrap">


        <div className="col-3 m-1 p-1">
          <Form.Label htmlFor="name">Name</Form.Label>
          <Form.Control
            id="name"
            type="text"
            placeholder="Name"
            value={formData.name}
            onChange={handleChange}
          />
        </div>

        <div className="col-1 m-1 p-1">
          <Form.Label htmlFor="age">Age</Form.Label>
          <Form.Control
            id="age"
            type="text"
            placeholder="0"
            value={formData.age}
            onChange={handleChange}
          />
        </div>

        <div className="col-3 m-1 p-1">
          <Form.Label htmlFor="motto">Title</Form.Label>
          <Form.Control
            id="motto"
            type="text"
            placeholder="Besides your name and picture, this is the first chance you have to turn people off."
            value={formData.motto}
            onChange={handleChange}
          />
        </div>

        <div className="col-3 m-1 p-1">
          <Form.Label htmlFor="profile">Profile Picture</Form.Label>
          <Form.Control
            id="profile"
            type="text"
            placeholder="https://wwww.example.com/profile.jpg"
            value={formData.profile}
            onChange={handleChange}
          />
        </div>

        <div className="col-1 m-1 p-1">
          <Form.Label htmlFor="latitude">Latitude</Form.Label>
          <Form.Control
            id="latitude"
            type="text"
            placeholder="Latitude"
            value={formData.latitude}
            onChange={handleChange}
          />
        </div>

        <div className="col-1 m-1 p-1">
          <Form.Label htmlFor="longitude">Longitude</Form.Label>
          <Form.Control
            id="longitude"
            type="text"
            placeholder="Longitude"
            value={formData.longitude}
            onChange={handleChange}
          />
        </div>

        <div className="col-1 m-1 p-1">
          <Form.Label htmlFor="longitude">Describe yourself</Form.Label>
          <Form.Control
            id="description"
            type="textarea"
            placeholder="This is the part where you pretend to have something interesting to say."
            value={formData.description}
            onChange={handleChange}
          />
        </div>
      </div>

      <DetailsSelections
        onChange={handleDetailsChange}
      />

      <Button variant="primary"  onClick={handleSubmit}>
        Finish your profile
      </Button>
    </>
  );
}

export default CreateProfile;