import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import NavDropdown from 'react-bootstrap/NavDropdown';
import Image from 'react-bootstrap/Image';
import FormTextExample from './login';
import Button from 'react-bootstrap/esm/Button';
import ChatModal from './chatmodal';

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
            <Nav.Link href="#deets">More deets</Nav.Link>
            <Nav.Link eventKey={2} href="#memes">
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
                <NavDropdown.Item>Bro you're not even signed in</NavDropdown.Item>
              ) : (
                <NavDropdown.Item>Look at you all logged in and shit</NavDropdown.Item>
              )}
              
              <NavDropdown.Item href="#action/3.2" className='text-center'>
                <Button>Logout</Button>
              </NavDropdown.Item>
                <FormTextExample />      
              
              <NavDropdown.Item href="#action/3.3">Something</NavDropdown.Item>
              <NavDropdown.Divider />
              <NavDropdown.Item href="#action/3.4">
                Separated link
              </NavDropdown.Item>
            </NavDropdown>
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default NavBar;