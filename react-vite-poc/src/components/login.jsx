import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
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
        <div className="m-1 p-1">
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
        </div>
        <div className="m-1 p-1">
          <Form.Label htmlFor="inputPassword" >Password</Form.Label>
          <Form.Control
            type="password"
            id="password"
            placeholder="Your password is bad and you should feel bad"
            value={formData.password}
            onChange={handleChange}
            
          />
        </div>
        
        <Button variant="danger" className='w-50' onClick={handleSubmit}><span className='text-black'>Login</span></Button>
      </div>
    </>
  );
}



export default SignIn;
