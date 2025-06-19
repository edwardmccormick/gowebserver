import React, { useState, useEffect } from 'react';
import Button from 'react-bootstrap/esm/Button';
import Offcanvas from 'react-bootstrap/Offcanvas';
import NavLink from 'react-bootstrap/esm/NavLink';
import ChatModal from './chatmodal';

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

function countUndeliveredMessages(User, match) {
      let accepted_unread
      let offered_unread 
      (match.accepted_chat == undefined || match.accepted_chat == null || match.accepted_chat == {}) ? accepted_unread = 0 : accepted_unread = match.accepted_chat.length;
      (match.offered_chat == null || match.offered_chat == undefined || match.offered_chat == {}) ? offered_unread = 0 : offered_unread = match.accepted_chat.length;
      if (match.offered == User.id) {
        return offered_unread 
       } 
       else {
        return accepted_unread;
      }
}

function ChatSelect({User, matches, pendings, offereds, setShowConfirmMatch}) {
  const [showOffcanvas, setShowOffcanvas] = useState(false); // Controls the Offcanvas visibility
  // const [showModal, setShowModal] = useState(false); // Controls the ChatModal visibility
  // const [selectedPerson, setSelectedPerson] = useState(null); // Tracks the currently selected person
  // const [selectedMatch, setSelectedMatch] = useState(null);
  const handleOffcanvasClose = () => setShowOffcanvas(false);
  const toggleOffcanvasShow = () => setShowOffcanvas((s) => !s);
  
  // console.log(matchesObject);
  // console.log(pendingsObject);
  // console.log(offeredsObject);


  // let matches 
  // { (matchesObject !== null || matchesObject !== undefined) ? matches= matchesObject?.matches : null; }
  // let pendings 
  // if (pendingsObject !== null || pendingsObject !== undefined) {pendings = pendingsObject?.pendings;}
  // let offereds
  // if (offeredsObject !== null || offeredsObject !== undefined) {offereds = offeredsObject?.offereds;}

  useEffect(() => {
    console.log("Matches updated:", matches);
  }, [matches]);

  return (!User || !matches 
    ? null
    :
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
              {(matches == undefined || matches == null) 
                ? <div className="d-flex align-items-center">
                    <strong role="status">Loading...</strong>
                    <div className="spinner-border ms-auto" aria-hidden="true"></div>
                  </div>
                  
                : matches.map((match) => (
                  // {match.accepted_time == undefined ? null : (
                  <div key={`${match.person.name}-${match.id}-${match.person.id}`}>
                  <ChatModal
                    key={`chatwith${match.person.name}-${match.id}`}
                    person={match.person} // Pass the currently selected person to the modal
                    match={match}
                    User={User} // Pass the User prop to the ChatModal
                    unreadmessages={countUndeliveredMessages(User, match)}
                  />
                  {/* <ChatModalButton
                    key={match.person.id}
                    match={match}
                    person={match.person}
                    setSelectedPerson={setSelectedPerson} // Pass the setter for the selected person
                    setShow={setShowModal} // Pass the setter for modal visibility
                    setSelectedMatch={setSelectedMatch}
                    message={match.person.motto ? match.person.motto.length : 69}
                  /> */}
                  </div>
                  // )}
                ))
              }
            </div>
          { pendings == undefined || pendings == null || offereds == undefined || offereds == null
            ? null
            : <>
            <h4 className='text-center'>Matches to review:</h4>
              <div className='d-flex flex-column'>
                {pendings.map((pending) => (
                <Button
                  key={`pendingmatchfrom${pending.person.name}-${pending.id}`}     
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
                    key={`pendingmatchfrom${pending.person.name}-${pending.id}`}
                    variant="outline-warning"
                    className="p-2 fs-5 m-2 text-right d-flex flex-row justify-content-between align-items-center"
                    onClick={() => {
                      setShowOffcanvas(false);
                      setShowConfirmMatch(true)
                    }}
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
            </>
          }
          {/* <ChatModal
            show={showModal}
            setShow={setShowModal}
            person={selectedPerson} // Pass the currently selected person to the modal
            match={selectedMatch}
            User={User} // Pass the User prop to the ChatModal
          /> */}
        </Offcanvas.Body>
      </Offcanvas>
    </>
  );
}

export default ChatSelect