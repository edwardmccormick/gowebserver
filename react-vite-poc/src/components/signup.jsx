import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';

export function SignUp() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleSubmit = async () => {
    const payload = {
      email: formData.email,
      password: formData.password
    };

    try {
      const response = await fetch('http://localhost:8080/signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        alert('User added successfully!');
      } else {
        alert('Failed to add User.');
      }
    } catch (error) {
      console.error('Error submitting form:', error);
      alert('An error occurred while submitting the form.');
    }
  };

 return (
    <>
      <div className="d-flex w-75 mx-auto justify-content-around align-items-center ">
        <div className="col-5 m-1 p-1">
          <Form.Label htmlFor="email">Your email address</Form.Label>
          <Form.Control
            id="email"
            type="text"
            placeholder="urmom@issofat.com"
            value={formData.email}
            onChange={handleChange}
          />
        </div>

        <div className="col-5 m-1 p-1">
          <Form.Label htmlFor="password">Password</Form.Label>
          <Form.Control
            id="password"
            type="password"
            placeholder="password123"
            value={formData.password}
            onChange={handleChange}
          />
        </div>
      </div>
      <div className='px-auto text-center'>
        <Button variant="outline-primary" className='w-25' onClick={handleSubmit}>
            Sign up!
        </Button>
      </div>
    </>
  );
}

export default SignUp;