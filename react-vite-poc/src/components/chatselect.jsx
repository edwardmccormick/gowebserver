import React, { useState, useEffect, use } from 'react';
import Button from 'react-bootstrap/esm/Button';
import Offcanvas from 'react-bootstrap/Offcanvas';
import NavLink from 'react-bootstrap/esm/NavLink';
import ChatModal from './chatmodal';
import { ChatModalButton } from './chatmodal';

function convertISODateToLocal(dateString) {
  const date = new Date(dateString);


  const localTime = date.toLocaleTimeString('en-US', {
    month: 'long',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
  });
  return localTime; // Output: "09:03:56 PM"

}

function ChatSelect({User}) {
  const [showOffcanvas, setShowOffcanvas] = useState(false); // Controls the Offcanvas visibility
  const [showModal, setShowModal] = useState(false); // Controls the ChatModal visibility
  const [selectedPerson, setSelectedPerson] = useState(null); // Tracks the currently selected person
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);
  const [pendings, setPendings] = useState([]);
  const [offereds, setOffereds] = useState([]);

  const handleOffcanvasClose = () => setShowOffcanvas(false);
  const toggleOffcanvasShow = () => setShowOffcanvas((s) => !s);

  useEffect(() => {
      if (User == null || User == undefined) return; 
    fetch(`http://localhost:8080/matches/${User.id}`)
      .then((res) => res.json())
      .then((data) => {
        const offered = data.filter((match) => match.accepted_time == "0001-01-01T00:00:00Z"); // This is what a null date looks like in Go
        const accepted = data.filter((match) => match.accepted_time != "0001-01-01T00:00:00Z"); // A non-null date
        const pending = offered.filter((match) => match.offered != User.id)
        setMatches(accepted);
        setPendings(pending);
        setOffereds(offered);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [User]); // Run the effect only when User changes

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
          <h4 className='text-center'>Chat with your active matches:</h4>
            <div className='d-flex flex-column'>
                {matches.map((match) => (
                  // {match.accepted_time == undefined ? null : (
                  <ChatModalButton
                    key={match.person.id}
                    person={match.person}
                    setSelectedPerson={setSelectedPerson} // Pass the setter for the selected person
                    setShow={setShowModal} // Pass the setter for modal visibility
                    message={match.person.motto ? match.person.motto.length : 69}
                  />
                  // )}
                ))}
            </div>

          <h4 className='text-center'>Matches to review:</h4>
            <div className='d-flex flex-column'>
                {pendings.map((pending) => (
                <Button       
                  variant="outline-warning"
                  className="p-2 fs-5 m-2 text-right d-flex flex-row justify-content-between align-items-center"
                >
                      <img
                        src={pending.person.profile ? pending.person.profile : '/profile.svg'}
                        style={{ borderRadius: '50%' }}
                        className="m-1 p-1"
                        height="50"
                        width="50"
                        alt={`${pending.person.name}'s profile`}
                      />
                      {pending.person.name} liked you at {convertISODateToLocal(pending.offered_time)}
                  
                </Button>
                ))}
            </div>
          <h4 className='text-center'>Matches waiting on a response:</h4>
            <div className='d-flex flex-column'>
                {offereds.map((pending) => (
                  <Button 
                   variant="outline-warning"
                    className="p-2 fs-5 m-2 text-right d-flex flex-row justify-content-between align-items-center"
                  >
                    <img
                      src={pending.person.profile ? pending.person.profile : '/profile.svg'}
                      style={{ borderRadius: '50%' }}
                      className="m-1 p-1"
                      height="50"
                      width="50"
                      alt={`${pending.person.name}'s profile`}
                    />
                      You liked {pending.person.name} at {convertISODateToLocal(pending.offered_time)}
                   
                  </Button>
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