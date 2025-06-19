import { LogOut } from 'lucide-react';
import Button from 'react-bootstrap/Button';

function Logout({ 
  setLoggedInUser, 
  setJWT,
  jwt
  }) {
    
  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleSubmit = async () => {
    try {
      const response = await fetch('http://localhost:8080/logout', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `${jwt}`, 
        },
      });

      if (response.ok) {
        setLoggedInUser(null); // Update the logged-in user in App.jsx
        setJWT(null);
        alert('Logout successful!');
      } else {
        alert('Invalid email or password.');
      }
    } catch (error) {
      console.error('Error during logout:', error);
      alert('An error occurred while logging out. Guess you\`re stuck with us!');
    }
  };

  return (
    <>     
    <div className='text-center p-1 m-1'>

        <Button variant="primary" onClick={handleSubmit}><LogOut />Logout</Button>
      </div>
    </>
  );
}

export default Logout;
