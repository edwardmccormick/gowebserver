import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Button from 'react-bootstrap/Button';
import humanity from '../../../humanity.json';
import DetailsSelections from './detailsselections';
import details from '../../../details.json';

function SignIn({ 
  setLoggedInUser , 
  setJWT}) {
    const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleSubmit = async () => {
    try {
      const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: formData.email,
          password: formData.password,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setLoggedInUser(data.person); // Update the logged-in user in App.jsx
        setJWT(data.token);
        setShowDropdown(true);
        alert('Login successful!');
      } else {
        alert('Invalid email or password.');
      }
    } catch (error) {
      console.error('Error during login:', error);
      alert('An error occurred while logging in.');
    }
  };

  return (
    <>     
    <div className='text-center'>
      <h4 className='text-center'>Login:</h4>
        <Form.Label htmlFor="name">Email</Form.Label>
          <Form.Control
            id="email"
            type='text'
            placeholder='yourmom@issofat.com'
            value={formData.email}
            onChange={handleChange}
            aria-label="Default"
            aria-describedby="email"
          />
    

        <Form.Label htmlFor="inputPassword" >Password</Form.Label>
        <Form.Control
          type="password"
          id="password"
          placeholder="Enter your password"
          value={formData.password}
          onChange={handleChange}
          aria-describedby="passwordHelpBlock"
        />
        <Form.Control.Feedback type="invalid" tooltip>
          Your password is 8-20 characters long, contains letters and numbers, and does not contain spaces, special characters, or emoji.
        </Form.Control.Feedback>
        <br />
        <Button variant="primary" onClick={handleSubmit}>Login</Button>
      </div>
    </>
  );
}

export function SignUpProfile() {
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
      id: parseInt(formData.id),
      name: formData.name,
      motto: formData.motto,
      lat: parseFloat(formData.latitude),
      long: parseFloat(formData.longitude),
      profile: formData.profile,
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
        alert('Person added successfully!');
      } else {
        alert('Failed to add person.');
      }
    } catch (error) {
      console.error('Error submitting form:', error);
      alert('An error occurred while submitting the form.');
    }
  };

 return (
    <>
      <div className="d-flex justify-content-center align-items-center flex-wrap">
        <div className="col-1 m-1 p-1">
          <Form.Label htmlFor="id">ID</Form.Label>
          <Form.Control
            id="id"
            type="number"
            placeholder="0"
            value={formData.id}
            onChange={handleChange}
          />
        </div>

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

        <div className="col-3 m-1 p-1">
          <Form.Label htmlFor="motto">Motto</Form.Label>
          <Form.Control
            id="motto"
            type="text"
            placeholder="Motto"
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
      </div>

      <DetailsSelections
        onChange={handleDetailsChange}
      />

      <Button variant="primary" onClick={handleSubmit}>
        Sign up!
      </Button>
    </>
  );
}

export default SignIn;
