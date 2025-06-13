import React, { useState, useEffect } from 'react';
import Button from 'react-bootstrap/esm/Button';
import Stack from 'react-bootstrap/Stack';
import Offcanvas from 'react-bootstrap/Offcanvas';
import NavLink from 'react-bootstrap/esm/NavLink';
import ChatModal from './chatmodal';


function ChatSelect() {
  const [show, setShow] = useState(false);

  const handleClose = () => setShow(false);
  const toggleShow = () => setShow((s) => !s);

    const [matches, setMatches] = useState([]);
    const [loading, setLoading] = useState(true);
  
    useEffect(() => {
      fetch('http://localhost:8080/people')
        .then((res) => res.json())
        .then((data) => {
          setMatches(data);
          setLoading(false);
        })
        .catch(() => setLoading(false));
    }, []);
  
    console.log(matches);
    console.log(loading);

  return (
    <>
      <NavLink
        to="/chat" // Specify the route you want to navigate to
        onClick={toggleShow} // Keep the toggle functionality
      >
        Chat Selector
      </NavLink>
      <Offcanvas show={show} onHide={handleClose}>
        <Offcanvas.Header closeButton>
          <Offcanvas.Title>Your current matches:</Offcanvas.Title>
        </Offcanvas.Header>
        <Offcanvas.Body>
          <h4 className='text-center'>Active matches:</h4>
            <div className='d-flex flex-column'>
                {matches.map((person) => (
                <ChatModal
                person={person}
                >
                {/* // props=variant='outline-success fs-5 m-2 text-right d-flex flex-row justify-content-around' className="p-2"> */}

                </ChatModal>
                ))}
            </div>
          <h4 className='text-center'>Pending matches:</h4>
            <Stack gap={1}>
                {matches.map((person) => (
                <div className="p-2">
                    <img
                        src={person.profile ? person.profile : '/profile.svg'}
                        style={{ borderRadius: '50%' }}
                        className="m-1 p-1"
                        height="50"
                        width="50"
                        alt={`${person.name}'s profile`}
                    />
                    {person.name}
                </div>
                ))}
            </Stack>
          <h4 className='text-center'>Closed matches:</h4>
            <Stack gap={1}>
                {matches.map((person) => (
                <div className="p-2">
                    <img
                        src={person.profile ? person.profile : '/profile.svg'}
                        style={{ borderRadius: '50%' }}
                        className="m-1 p-1"
                        height="50"
                        width="50"
                        alt={`${person.name}'s profile`}
                    />
                    {person.name}
                </div>
                ))}
            </Stack>
        </Offcanvas.Body>
      </Offcanvas>
    </>
  );
}

export default ChatSelect