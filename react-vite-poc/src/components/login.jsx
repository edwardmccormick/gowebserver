import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Button from 'react-bootstrap/Button';
import nonlinear from '../../../nonlinear.json';
import VerbiageSelections from './verbiageselections';
import verbiage from '../../../verbiage.json';

function FormTextExample() {
  return (
    <>     
    <div className='text-center'>
      <h4 className='text-center'>Login:</h4>
        <Form.Label htmlFor="name">Email</Form.Label>
          <Form.Control
            id="email"
            type='text'
            placeholder='yourmom@issofat.com'
            aria-label="Default"
            aria-describedby="email"
          />
    

        <Form.Label htmlFor="inputPassword" >Password</Form.Label>
        <Form.Control
          type="password"
          id="inputPassword"
          aria-describedby="passwordHelpBlock"
        />
        <Form.Control.Feedback type="invalid" tooltip>
          Your password is 8-20 characters long, contains letters and numbers, and does not contain spaces, special characters, or emoji.
        </Form.Control.Feedback>
        <br />
        <Button variant="primary">Login</Button>
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
    verbiage: Object.keys(verbiage).reduce((acc, key) => ({ ...acc, [key]: 0 }), {}),
  });

  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleVerbiageChange = (key, value) => {
    setFormData((prev) => ({
      ...prev,
      verbiage: { ...prev.verbiage, [key]: parseInt(value, 10) },
    }));
  };

  const handleSubmit = async () => {
    const payload = {
      name: formData.name,
      motto: formData.motto,
      lat: parseFloat(formData.latitude),
      long: parseFloat(formData.longitude),
      profile: formData.profile,
      verbiage: formData.verbiage,
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

      <VerbiageSelections
        onChange={(key, value) => handleVerbiageChange(key, value)}
      />

      <Button variant="primary" onClick={handleSubmit}>
        Sign up!
      </Button>
    </>
  );
}

export default FormTextExample;
