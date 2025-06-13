import { useEffect, useRef, useState } from 'react';
import Button from 'react-bootstrap/Button';
import NavLink from 'react-bootstrap/esm/NavLink';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Modal from 'react-bootstrap/Modal';

function ChatModal({person}) {
  console.log(person);
  const [show, setShow] = useState(false);

  const handleClose = () => setShow(false);
  const handleShow = () => setShow(true);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const ws = useRef(null);
  
  useEffect(() => {
    ws.current = new WebSocket('ws://localhost:8080/ws');

    ws.current.onopen = () => {
      setMessages((prev) => [...prev, { message: ' opened.', who: 'Connection', id: Date.now() }]);
      console.log('WebSocket connection established');
    };

    ws.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        // If the server sends an array of messages
        if (Array.isArray(data)) {
          setMessages((prev) => [...prev, ...data]);
        } else if (typeof data === 'object' && data !== null) {
          setMessages((prev) => [...prev, data]);
        } else {
          // If it's not JSON, treat as a status message
          setMessages((prev) => [...prev, { message: event.data, who: 'System', id: Date.now() }]);
        }
      } catch (e) {
        // Not JSON, treat as a status message
        setMessages((prev) => [...prev, { message: event.data, who: 'System', id: Date.now() }]);
      }
    };

    ws.current.onclose = () => {
      setMessages((prev) => [...prev, { message: 'Connection closed', who: 'Alert', id: Date.now() }]);
    };

    ws.current.onerror = (err) => {
      setMessages((prev) => [...prev, { message: 'WebSocket error', who: 'System', id: Date.now() }]);
    };

    return () => {
    if (ws.current) {
      ws.current.close();
    }
  };
  }, []);
  
    const sendMessage = () => {
      if (ws.current && input) {
        ws.current.send(input);
        setInput( {message: 'WebSocket error', who: 'System', id: Date.now() });
      }
    };

  return (
  <>
    {!person ? (
      <div>
        <span>Kind of empty in here, isn't it?</span>
      </div>
    ) : (
      <Button
        variant="outline-success"
        className="p-2 fs-5 m-2 text-right d-flex flex-row justify-content-around"
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
      </Button>
    )}

    <Modal
      show={show}
      onHide={handleClose}
      size="lg"
      aria-labelledby="contained-modal-title-vcenter"
      centered
      scrollable={true}
    >
      <Modal.Header closeButton>
        <Modal.Title>Chat with The Other Person</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <div className="overflow-scroll">
          <div className="m-1 p-1 d-flex flex-column">
            {messages.map((msg, idx) => (
              <div
                key={msg.id || idx}
                className={
                  msg.who === 'System'
                    ? 'text-center bg-warning text-white rounded m-2 p-2 flex-item'
                    : msg.who === 'Alert'
                    ? 'text-center bg-danger text-white rounded m-2 p-2 flex-item'
                    : msg.who === 'Connection'
                    ? 'text-center bg-success text-white rounded m-2 p-2 flex-item'
                    : msg.who === 'Admin'
                    ? 'align-self-center bg-info text-white rounded m-2 p-2 flex-item'
                    : msg.who === 'Them'
                    ? 'align-self-start bg-light text-black rounded m-2 p-2 flex-item'
                    : 'align-self-end bg-primary text-white rounded m-2 p-2 flex-item'
                }
              >
                {msg.who} - {msg.message}
              </div>
            ))}
          </div>
        </div>
      </Modal.Body>
      <Modal.Footer className="d-flex flex-row justify-content-around mx-3">
        <Form className="flex-grow-1">
          <InputGroup>
            <Form.Control
              type="text"
              placeholder="Type a message"
              autoFocus
              value={input}
              onChange={(e) => setInput(e.target.value)}
            />
            <Button variant="success" onClick={sendMessage}>
              Send
            </Button>
          </InputGroup>
        </Form>

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
);
}

export default ChatModal;