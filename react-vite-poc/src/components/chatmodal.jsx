import { useEffect, useRef, useState } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Modal from 'react-bootstrap/Modal';
import Badge from 'react-bootstrap/Badge';

function convertISODateToLocal(dateString) {
  const date = new Date(dateString);


  const localTime = date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: true,
  });
  return localTime; // Output: "09:03:56 PM"

}

export function ChatModal({match, person, User, unreadmessages}) {

  const [showModal, setShowModal] = useState(false); // Controls the ChatModal visibility
  const handleClose = () => setShowModal(false);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const ws = useRef(null);
  const messagesEndRef = useRef(null); // Ref for the last message
  
    // Scroll to the bottom whenever messages change
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);


  useEffect(() => {
    if (!match || !person || !User || !showModal) {
      return;
    }
    ws.current = new WebSocket(`ws://localhost:8080/ws?id=${match.id}&user_id=${User.id}`);

    ws.current.onopen = () => {
      setMessages((prev) => [...prev, { message: ' opened.', who: 'Connection', id: Date.now(), time: Date.now() }]);
      console.log('WebSocket connection established');
    };

    ws.current.onmessage = (event) => {
      try {
        let date = new Date()
        const data = JSON.parse(event.data);
        // If the server sends an array of messages
        if (Array.isArray(data)) {
          setMessages((prev) => [...prev, ...data]);
        } else if (typeof data === 'object' && data !== null) {
          data.who == person.id ? data.who = person.name: null
          setMessages((prev) => [...prev, data]);
        } else {
          // If it's not JSON, treat as a status message
          setMessages((prev) => [...prev, { message: event.data, who: 'System', id: Date.now(), time: Date.now()}]);
        }
      } catch (e) {
        // Not JSON, treat as a status message
        setMessages((prev) => [...prev, { message: event.data, who: 'System', id: Date.now(), time: Date.now() }]);
      }
    };

    ws.current.onclose = () => {
      setMessages((prev) => [...prev, { message: 'Connection closed', who: 'Alert', id: Date.now(), time: Date.now() }]);
    };

    ws.current.onerror = (err) => {
      setMessages((prev) => [...prev, { message: 'WebSocket error', who: 'System', id: Date.now(), time: Date.now() }]);
    };

    return () => {
    if (ws.current) {
      ws.current.close();
    }
  };
  }, [showModal]);
  
  const sendMessage = () => {
    if (ws.current && input) {
      let date = new Date();
      ws.current.send(JSON.stringify({ message: input, who: User.id, id: Date.now() })); // Send the message as a JSON object
      setMessages((prev) => [...prev, { message: input, who: 'Me', id: Date.now(), time: date.toLocaleTimeString() }]); // Add the message to the chat
      setInput(''); // Reset the input field to an empty string
    }
  };

  return (!person || !User || !match 
    ? <div>Sign in to see your matches and chat.</div>
    :
  <>
<style jsx>{`
        .vibe-chat-btn {
  position: relative;
  background: linear-gradient(45deg, #ff4500, #ff8c00, #ffd700);
  color: white;
  font-weight: bold;
  text-transform: uppercase;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  box-shadow: 0 0 15px rgba(255, 69, 0, 0.8), 0 0 30px rgba(255, 140, 0, 0.6);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  overflow: hidden;
}

.vibe-chat-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 0 25px rgba(255, 69, 0, 1), 0 0 50px rgba(255, 140, 0, 0.8);
}

.vibe-chat-btn::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255, 69, 0, 0.6), rgba(255, 140, 0, 0.4), transparent);
  animation: flame 1.5s infinite ease-in-out;
  z-index: -1;
}

@keyframes flame {
  0% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 1;
  }
  50% {
    transform: translate(-50%, -50%) scale(1.2);
    opacity: 0.8;
  }
  100% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 1;
  }
}
      `}</style>

    <>
    <Button
      key={`openchatwith${person.name}-${match.id}`}
      variant="outline-success"
      className="p-2 fs-5 m-2 text-right d-flex flex-row justify-content-between align-items-center"
      onClick={() => setShowModal(true)}
    >
      <img
        src={person.profile ? person.profile : '/profile.svg'}
        style={{ borderRadius: '50%' }}
        className="m-1 p-1"
        height="50"
        width="50"
        alt={`${person.name}'s profile`}
      />
      {person.name ? `Chat with ${person.name}` : 'Chat'}
      {unreadmessages ?
      <>
      <Badge bg="secondary" className='align-self-start'>{unreadmessages}</Badge>
        <span className="visually-hidden">unread messages</span>
      </>
      : null}
    </Button>
    </>
    <>
    <Modal
      key={`chatwindowwith${person.name}-${match.id}`}
      show={showModal}
      onHide={handleClose}
      size="lg"
      aria-labelledby="contained-modal-title-vcenter"
      centered
      scrollable={true}
    >

      <Modal.Header closeButton>
        <Modal.Title>Chat with {person.name} <img
              src={person.profile ? person.profile : '/profile.svg'}
              style={{ borderRadius: '50%' }}
              className="m-1 p-1"
              height="50"
              width="50"
              alt={`${person.name}'s profile`}
            /></Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <div className="overflow-scroll">
          <div className="m-1 p-1 d-flex flex-column">
            {messages.map((msg, idx) => (
              <div className={
                      msg.who === 'System'
                      ? 'text-center bg-warning text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Alert'
                      ? 'text-center bg-danger text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Connection'
                      ? 'text-center bg-success text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Admin' || msg.who === 'Moderator' || msg.who === "AI Host"
                      ? 'align-self-center bg-info text-white rounded m-2 p-2 flex-item'
                      : msg.who === msg.who == `${User.id}` || msg.who === 'Me'
                      ? 'd-flex flex-row justify-content-end align-items-center'
                      : 'd-flex flex-row justify-content-start align-items-center'
                  }
                  key={`${person.id}${idx}`}>
                  {msg.who == `${User.id}` || msg.who === "Me" ?<span className="text-body-tertiary fs-6">{msg.time}</span> : null}
                  <span className={
                    msg.who === 'System'
                      ? 'text-center bg-warning text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Alert'
                      ? 'text-center bg-danger text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Connection'
                      ? 'text-center bg-success text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Admin' || msg.who === 'Moderator' || msg.who === "AI Host"
                      ? 'align-self-center bg-info text-white rounded m-2 p-2 flex-item'
                      : msg.who == `${User.id}` || msg.who === 'Me'
                      ? 'bg-primary text-white rounded m-2 p-2 flex-item'
                      : 'bg-light text-black rounded m-2 p-2 flex-item'
                  }
                >
                  {msg.who} - {msg.message} 
                  </span>
                  {msg.who == `${person.id}` || msg.who == "Them" || msg.who == `${person.name}` || msg.who !=="Me" 
                  ?<span className="text-body-tertiary fs-6">  {
                    convertISODateToLocal(msg.time)                    
                    }</span> : null}
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
        </div>
      </Modal.Body>
      <Modal.Footer className="d-flex flex-row justify-content-around mx-3 flex-wrap">
        <div className='flex-grow-1'>
          <InputGroup>
            <Form.Control
              type="text"
              placeholder="Type a message"
              autoFocus
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  e.preventDefault(); // Prevent the default behavior of the Enter key
                  sendMessage(); // Call the sendMessage function
                }
              }}
            />
            <Button variant="success" onClick={sendMessage}>
              Send
            </Button>
          </InputGroup>
        </div>
        <Button
          variant="secondary"
          onClick={handleClose}
          className="col-3 mx-5 vibe-chat-btn"
        >
          Vibe Chat ❤️
        </Button>
        <Button
          variant="secondary"
          onClick={handleClose}
          className="col-3 mx-5"
        >
          Close
        </Button>
      </Modal.Footer>
    </Modal>
    </>
  </>    
  );
}

export function ChatModalButton({ match, person, setSelectedPerson, setSelectedMatch, setShow, message }) {
  const handleShow = () => {
    setSelectedMatch(match)
    setSelectedPerson(person); // Set the selected person
    setShow(true); // Open the modal
  };

  return (
    <>
    <Button
      variant="outline-success"
      className="p-2 fs-5 m-2 text-right d-flex flex-row justify-content-between align-items-center"
      onClick={handleShow}
    >
      <img
        src={person.profile ? person.profile : '/profile.svg'}
        style={{ borderRadius: '50%' }}
        className="m-1 p-1"
        height="50"
        width="50"
        alt={`${person.name}'s profile`}
      />
      {person.name ? `Chat with ${person.name}` : 'Chat'}
      {message ?
      <>
      <Badge bg="secondary" className='align-self-start'>{message}</Badge>
        <span className="visually-hidden">unread messages</span>
      </>
      : null}
    </Button>
    
    </>
  );
}

export default ChatModal;