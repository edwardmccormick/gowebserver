import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Button from 'react-bootstrap/Button';


function FormTextExample() {
  return (
    <>      
    <h4 className='text-center'>Login:</h4>
      <InputGroup className="mb-3">
        <InputGroup.Text id="email">
          Email Address:
        </InputGroup.Text>
        <Form.Control
          aria-label="Default"
          aria-describedby="email"
        />
      </InputGroup>
        <p className='text-center'>OR</p>
      <InputGroup className="mb-3">
        <InputGroup.Text id="username">
          Username:
        </InputGroup.Text>
        <Form.Control
          aria-label="Default"
          aria-describedby="username"
        />
      </InputGroup>

      <Form.Label htmlFor="inputPassword">Password</Form.Label>
      <Form.Control
        type="password"
        id="inputPassword"
        aria-describedby="passwordHelpBlock"
      />
      <Form.Text id="passwordHelpBlock" muted>
        Your password is 8-20 characters long, contains letters and numbers,
        and does not contain spaces, special characters, or emoji.
      </Form.Text>
      <br />
      <Button variant="primary">Login</Button>
    </>
  );
}

function BasicExample() {
  return (
    <>
      <div className='d-flex justify-content-center align-items-center flex-wrap'>
        <div className='col-3 m-1 p-1'>
          <Form.Label htmlFor="createName">Name</Form.Label>
          <Form.Control
            id='createName'
            type='text'
            placeholder="Name"
            aria-label="Name"
            aria-describedby="Name"
          />
        </div>
      
      <div className='col-3 m-1 p-1'>
        <Form.Label htmlFor="createMotto">Motto</Form.Label>
        <Form.Control id='createMotto' type='text' aria-label="Motto or Title for the person" />
      </div>

      <div className='col-1 m-1 p-1'>
        <Form.Label htmlFor="createLatitude">Latitude</Form.Label>
        <Form.Control
          placeholder="Recipient's username"
          aria-label="Recipient's username"
          aria-describedby="basic-addon2"
        />
      </div>

      <div className='col-1 m-1 p-1'>
        <Form.Label htmlFor="createLongitude">Longitude</Form.Label>
        <Form.Control
          placeholder="Longitude"
          aria-label="Recipient's Lattitude"
          aria-describedby="basic-addon3"
        />
        <Form.Control id="basic-url" aria-describedby="basic-addon3" />
      </div>

      <div className='col-3 m-1 p-1'>
            <Form.Select aria-label="Default select example">
              <option>Open this select menu</option>
              <option value="0">Zero</option>
              <option value="1">One</option>
              <option value="2">Two</option>
              <option value="3">Three</option>
              <option value="4">Four</option>
              <option value="5">Five</option>
              <option value="6">Six</option>
              <option value="7">Seven</option>
              <option value="8">Eight</option>
              <option value="9">Nine</option>
              <option value="10">Ten</option>
            </Form.Select>
      </div>


      <Button variant="primary">Add User</Button>
      </div>
    </>
  );
}

export default FormTextExample;
export { BasicExample };
