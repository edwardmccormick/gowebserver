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
      <div className='d-flex align-items-center'>
        <div className='col-3'>
          <Form.Label htmlFor="createName">Name</Form.Label>
          <Form.Control
            id='createName'
            type='text'
            placeholder="Name"
            aria-label="Name"
            aria-describedby="Name"
          />
        </div>
      
      <div className='col-3'>
        <Form.Label htmlFor="createMotto">Motto</Form.Label>
        <Form.Control id='createMotto' as="textarea" aria-label="With textarea" />
      </div>

      <InputGroup className="mb-3">
        <Form.Control
          placeholder="Recipient's username"
          aria-label="Recipient's username"
          aria-describedby="basic-addon2"
        />
        <InputGroup.Text id="basic-addon2">@example.com</InputGroup.Text>
      </InputGroup>

      <Form.Label htmlFor="basic-url">Your vanity URL</Form.Label>
      <InputGroup className="mb-3">
        <InputGroup.Text id="basic-addon3">
          https://example.com/users/
        </InputGroup.Text>
        <Form.Control id="basic-url" aria-describedby="basic-addon3" />
      </InputGroup>

      <InputGroup className="mb-3">
        <InputGroup.Text>$</InputGroup.Text>
        <Form.Control aria-label="Amount (to the nearest dollar)" />
        <InputGroup.Text>.00</InputGroup.Text>
      </InputGroup>


      <Button variant="primary">Login</Button>
      </div>
    </>
  );
}

export default FormTextExample;
export { BasicExample };
