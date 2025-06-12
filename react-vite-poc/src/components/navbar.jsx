import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import NavDropdown from 'react-bootstrap/NavDropdown';
import Image from 'react-bootstrap/Image';
import FormTextExample from './login';
import Button from 'react-bootstrap/esm/Button';
import ChatModal from './chatmodal';
import Form from 'react-bootstrap/Form';
import Badge from 'react-bootstrap/Badge';

function NavBar({profile, username}) {
  return (
    <Navbar collapseOnSelect expand="sm" className="bg-body-tertiary">
      <Container>
        <Navbar.Brand href="#home"><Image src="./urmid.svg" height="50" alt="profile" roundedCircle/> urmid</Navbar.Brand>
        <Navbar.Toggle aria-controls="responsive-navbar-nav" />
        <Navbar.Collapse id="responsive-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link href="#features">Match</Nav.Link>
            <Nav.Item> <ChatModal /> </Nav.Item>

          </Nav>
          <Nav className="d.flex justify-content-around align-items-center">
            <Nav.Item>
              <Form>
                <Form.Check // prettier-ignore
                  type="switch"
                  id="light-dark-switch"
                />
              </Form>
            </Nav.Item>
            {/* <Nav.Link eventKey={2} href="#memes" /> */}
              <NavDropdown align={'end'} title={
                    <div>
                        <img className="focus-ring" 
                            src={profile}
                            style={{borderRadius: '50%'}}
                            height='50'
                            width='50' 
                            alt="user pic"
                        />
                    </div>
                }  id="collapsible-nav-dropdown" >
              <NavDropdown.Item href="#action/3.1" className='text-center'>Currently logged in as: <br /> {username}</NavDropdown.Item>
              {username == undefined && profile == undefined ? (
                <>
                  
                  <NavDropdown.Item>It's free to look but....no touching, okay?</NavDropdown.Item>
                  <NavDropdown.Item>Have an account? Why don't you...</NavDropdown.Item>
                  <div className='m-2 p-2'>
                    <FormTextExample />      
                  </div>
                </>
              ) : (
                <>
                  <NavDropdown.Item className='text-center'>    
                    <Button variant="primary" className='w-100'>
                        Chat
                      <Badge bg="secondary">9</Badge>
                      <span className="visually-hidden">unread messages</span>
                    </Button>
                  </NavDropdown.Item>
                  <NavDropdown.Item>
                    Update Profile
                  </NavDropdown.Item>
                  <NavDropdown.Item>
                    Settings
                  </NavDropdown.Item>
                  <NavDropdown.Item>
                    Account
                  </NavDropdown.Item>
                  <NavDropdown.Divider />
                  <NavDropdown.Item href="#action/3.2" className='text-center'>
                    <Button>Logout</Button>
                  </NavDropdown.Item>
                </>
              )}
              
              
            </NavDropdown>
            {/* </Nav.Link> */}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default NavBar;