import { useEffect, useRef, useState } from 'react';
import Button from 'react-bootstrap/Button';
import NavLink from 'react-bootstrap/esm/NavLink';
import Form from 'react-bootstrap/Form';
import Modal from 'react-bootstrap/Modal';
import Stack from 'react-bootstrap/Stack';

function ChatModal() {
  const [show, setShow] = useState(false);

  const handleClose = () => setShow(false);
  const handleShow = () => setShow(true);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const ws = useRef(null);
  
  useEffect(() => {
    ws.current = new WebSocket('ws://localhost:8080/ws');

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
          setMessages((prev) => [...prev, { message: event.data, isme: false, id: Date.now() }]);
        }
      } catch (e) {
        // Not JSON, treat as a status message
        setMessages((prev) => [...prev, { message: event.data, isme: false, id: Date.now() }]);
      }
    };

    ws.current.onclose = () => {
      setMessages((prev) => [...prev, { message: 'Connection closed', isme: false, id: Date.now() }]);
    };

    ws.current.onerror = (err) => {
      setMessages((prev) => [...prev, { message: 'WebSocket error', isme: false, id: Date.now() }]);
    };

    return () => {
      ws.current.close();
    };
  }, []);
  
    const sendMessage = () => {
      if (ws.current && input) {
        ws.current.send(input);
        setInput('');
      }
    };

  return (
    <>
      <NavLink variant="primary" onClick={handleShow}>
        Chat
      </NavLink>

      <Modal className='h-50' show={show} onHide={handleClose}
            size="lg"
            aria-labelledby="contained-modal-title-vcenter"
            centered>
        <Modal.Header closeButton>
          <Modal.Title>Chat with The Other Person</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <div className='overflow-scroll h-75'>
            <Stack gap={1}>
            {messages.map((msg, idx) => (
              <div key={msg.id || idx}
              // {msg.isme ? "className='text-end bg-primary text-white'" : "className='text-start bg-secondary text-white'"}
              >
                {msg.isme ? "Me" : "Them"} - {msg.message}
              </div>
            ))}
            </Stack>
          </div>
          <Form>            
  
            <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
              <Form.Label>Type a message</Form.Label>
              <Form.Control
                type="text"
                placeholder="Type a message"
                autoFocus
              />
            </Form.Group>
            <Button variant='success' onClick={sendMessage}>Send</Button>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>
            Close
          </Button>
          <Button variant="primary" onClick={handleClose}>
            Save Changes
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
}

export default ChatModal;