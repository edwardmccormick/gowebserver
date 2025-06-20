import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';

export function SignUp({ 
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
    const payload = {
      email: formData.signupEmail,
      password: formData.signupPassword
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
        const data = await response.json();
        setLoggedInUser(data.person); // Update the logged-in user in App.jsx
        console.log(data.person);
        setJWT(data.token);
        alert('Signup successful');
      } else {
        alert('We already have a user with that email address');
      }
    } catch (error) {
      console.error('Error submitting form:', error);
      alert('An error occurred while submitting the form.');
    }
  };

 return (
    <>
      <h4 className="text-center">Sign up</h4>
      <div className="text-center ">
        <div className=" m-1 p-1 ">
          <Form.Label htmlFor="signupEmail">Your email address</Form.Label>
          <Form.Control
            id="signupEmail"
            type="text"
            placeholder="urmom@iskindofaslut.com"
            value={formData.signupEmail}
            onChange={handleChange}
          />
        </div>

        <div className="m-1 p-1">
          <Form.Label htmlFor="password">Password</Form.Label>
          <Form.Control
            id="signupPassword"
            type="password"
            placeholder="password123...is a 𝘵𝘦𝘳𝘳𝘪𝘣𝘭𝘦 password"
            value={formData.signupPassword}
            onChange={handleChange}
          />
        </div>
      </div>
      <div className='px-auto text-center'>
        <Button variant="danger" className='w-50 ' onClick={handleSubmit}>
            <span className='text-black'>Sign up</span>
        </Button>
      </div>
    </>
  );
}

export default SignUp;