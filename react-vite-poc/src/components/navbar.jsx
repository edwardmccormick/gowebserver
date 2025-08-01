import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import NavDropdown from 'react-bootstrap/NavDropdown';
import Image from 'react-bootstrap/Image';
import SignIn from './login';
import Form from 'react-bootstrap/Form';
import ChatSelect from './chatselect';
import Logout from './logout';

function NavBar({
  User,
  setLoggedInUser, 
  setJWT, 
  jwt,
  refreshMatches, 
  matches, 
  pendings, 
  offereds, 
  setShowConfirmMatch,
  onSearchClick,
  onMeetClick,
  onFAQClick,
  onClickProfile,
  }) {

  
  return (
    <Navbar collapseOnSelect expand="sm" className="bg-body-tertiary">
      <Container>
        <Navbar.Brand href="#home"><Image src="./urmid.svg" height="50" alt="profile" roundedCircle/> urmid</Navbar.Brand>
        <Navbar.Toggle aria-controls="responsive-navbar-nav" />
        <Navbar.Collapse id="responsive-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link onClick={onSearchClick}>Search</Nav.Link>
            <Nav.Link onClick={onMeetClick}> Meet </Nav.Link>
            <Nav.Link onClick={onFAQClick}>FAQ</Nav.Link>
            <ChatSelect
              User={User}
              refreshMatches={refreshMatches}
              matches={matches}
              pendings={pendings}
              offereds={offereds}
              setShowConfirmMatch={setShowConfirmMatch}
            />
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
              <NavDropdown 
              align={'end'} 
              title={
                    <div>
                        <img className="focus-ring" 
                            src={ User ? User.profile.url : './profile.svg'}
                            style={{borderRadius: '50%'}}
                            height='50'
                            width='50' 
                            alt="user pic"
                        />
                    </div>
                }  id="collapsible-nav-dropdown" >
             
              {User == undefined ? (
                <>
                  
                  <NavDropdown.Item>It's free to look but....no touching, okay?</NavDropdown.Item>
                  <NavDropdown.Item>Have an account? Why don't you...</NavDropdown.Item>
                  <div className='m-2 p-2'>
                    <SignIn 
                    setLoggedInUser={setLoggedInUser} 
                    setJWT={setJWT}
                    />      
                  </div>
                </>
              ) : (
                <>
                  <NavDropdown.Item className='text-center'>    
                  
                  </NavDropdown.Item>
                  <NavDropdown.Item  className='text-center'>Currently logged in as: <br /> {User.name}</NavDropdown.Item>
                  <NavDropdown.Item
                  onClick={onClickProfile}
                  >
                    Update Profile
                  </NavDropdown.Item>
                  <NavDropdown.Item>
                    Settings
                  </NavDropdown.Item>
                  <NavDropdown.Item>
                    Account
                  </NavDropdown.Item>
                  <NavDropdown.Divider />
                  <NavDropdown.Item  className='text-center'>
                    <Logout 
                    setLoggedInUser={setLoggedInUser} 
                    setJWT={setJWT}
                    jwt={jwt}
                    />  
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