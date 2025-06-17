import React, { useState, useEffect } from 'react';
import Button from 'react-bootstrap/esm/Button';
import Offcanvas from 'react-bootstrap/Offcanvas';
import NavLink from 'react-bootstrap/esm/NavLink';
import ChatModal from './chatmodal';
import { ChatModalButton } from './chatmodal';

function ChatSelect({User}) {
  const [showOffcanvas, setShowOffcanvas] = useState(false); // Controls the Offcanvas visibility
  const [showModal, setShowModal] = useState(false); // Controls the ChatModal visibility
  const [selectedPerson, setSelectedPerson] = useState(null); // Tracks the currently selected person
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);

  const handleOffcanvasClose = () => setShowOffcanvas(false);
  const toggleOffcanvasShow = () => setShowOffcanvas((s) => !s);

  useEffect(() => {
    fetch(`http://localhost:8080/matches/${User.id}`)
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
        onClick={toggleOffcanvasShow} // Keep the toggle functionality
      >
        Chat Selector
      </NavLink>
      <Offcanvas show={showOffcanvas} onHide={handleOffcanvasClose}>
        <Offcanvas.Header closeButton>
          <Offcanvas.Title>Your current matches:</Offcanvas.Title>
        </Offcanvas.Header>
        <Offcanvas.Body>
          <h4 className='text-center'>Active matches:</h4>
            <div className='d-flex flex-column'>
                {matches.map((match) => (
                  {match.acceptedtime == undefined ? null : (
                  <ChatModalButton
                    key={match.person.id}
                    person={match.person}
                    setSelectedPerson={setSelectedPerson} // Pass the setter for the selected person
                    setShow={setShowModal} // Pass the setter for modal visibility
                    message={match.person.motto ? person.motto.length : 0}
                  />
                  )}
                ))}
            </div>

          <h4 className='text-center'>Pending matches:</h4>
            <div className='d-flex flex-column'>
                {matches.map((person) => (
                <ChatModalButton
                  key={person.id}
                  person={person}
                  setSelectedPerson={setSelectedPerson} // Pass the setter for the selected person
                  setShow={setShowModal} // Pass the setter for modal visibility
                  message={person.motto ? person.motto.length : 0}
              />
                ))}
            </div>
          <h4 className='text-center'>Closed matches:</h4>
            <div className='d-flex flex-column'>
                {matches.map((person) => (
                <ChatModalButton
                  key={person.id}
                  person={person}
                  setSelectedPerson={setSelectedPerson} // Pass the setter for the selected person
                  setShow={setShowModal} // Pass the setter for modal visibility
              />
                ))}
            </div>
          <ChatModal
            show={showModal}
            setShow={setShowModal}
            person={selectedPerson} // Pass the currently selected person to the modal
            User={User} // Pass the User prop to the ChatModal
          />
        </Offcanvas.Body>
      </Offcanvas>
    </>
  );
}

export default ChatSelect