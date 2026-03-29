import { useEffect, useRef, useState } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Modal from 'react-bootstrap/Modal';
import Badge from 'react-bootstrap/Badge';
import { apiUrl, wsUrl } from '../config/api';
import ChatMessage from './chatmessage';

function normalizeMessage(message, person, user) {
  if (!message || typeof message !== 'object') {
    return null;
  }

  const normalized = { ...message };

  if (normalized.who == person.id) {
    normalized.who = person.name;
  } else if (normalized.who == user.id) {
    normalized.who = 'Me';
  } else if (normalized.who === 0) {
    normalized.who = 'AI Introduction';
  }

  return normalized;
}

function mergeMessages(existingMessages, incomingMessages) {
  const merged = [...existingMessages];
  const seenIds = new Set(existingMessages.filter((msg) => msg?.id).map((msg) => msg.id));

  incomingMessages.forEach((message) => {
    if (!message) {
      return;
    }

    if (message.id && seenIds.has(message.id)) {
      return;
    }

    if (message.id) {
      seenIds.add(message.id);
    }
    merged.push(message);
  });

  merged.sort((a, b) => new Date(a.time || 0) - new Date(b.time || 0));
  return merged;
}

export function ChatModal({match, person, User, unreadmessages, jwt, clearChatNotification}) {

  const [showModal, setShowModal] = useState(false); // Controls the ChatModal visibility
  const handleOpen = () => {
    setShowModal(true);
    // Clear notification when opening the chat
    if (clearChatNotification && match?.ID) {
      clearChatNotification(match.ID);
    }
  };
  
  const handleClose = () => {
    // Keep message history when closing modal to prevent reloading
    setShowModal(false);
  };
  
  // Function to mark messages as read when opening the chat
  const markMessagesAsRead = async () => {
    if (!match || !User || !jwt || !unreadmessages || unreadmessages === 0) return;
    
    try {
      const response = await fetch(apiUrl(`/chat/markread/${match.ID}`), {
        method: 'POST',
        headers: {
          'Authorization': jwt,
          'Content-Type': 'application/json'
        }
      });
      
      if (!response.ok) {
        console.error('Failed to mark messages as read:', await response.text());
      }
    } catch (error) {
      console.error('Error marking messages as read:', error);
    }
  };
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const ws = useRef(null);
  const messagesEndRef = useRef(null); // Ref for the last message
  // Track loading state to show a spinner while fetching
  const [isLoading, setIsLoading] = useState(false);
  
    // Scroll to the bottom whenever messages change
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);


  // Fetch chat history from the server when needed
  const fetchChatHistory = async () => {
    if (!match || !User || !jwt) return;
    
    setIsLoading(true);
    try {
      const response = await fetch(apiUrl(`/chatmessages/${match.ID}`), {
        headers: {
          'Authorization': jwt,
          'Content-Type': 'application/json'
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        
        // Check if the response contains the messages array directly or within pagination
        const chatMessages = data.messages || data;
        
        const processedMessages = Array.isArray(chatMessages)
          ? chatMessages
              .map((msg) => normalizeMessage(msg, person, User))
              .filter(Boolean)
          : [];

        setMessages(processedMessages);
        console.log(`Loaded ${processedMessages.length} messages for match ${match.ID}`);
      } else {
        console.error('Failed to fetch chat history:', await response.text());
      }
    } catch (error) {
      console.error('Error fetching chat history:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Load chat history when the modal is first opened
  useEffect(() => {
    if (showModal) {
      fetchChatHistory();
      // Mark messages as read when opening the chat
      markMessagesAsRead();
    }
  }, [showModal, match?.ID, jwt]);

  useEffect(() => {
    if (!match || !person || !User || !showModal) {
      return;
    }
    
    ws.current = new WebSocket(wsUrl(`/ws?id=${match.ID}&user_id=${User.id}`));

    ws.current.onopen = () => {
      console.log('WebSocket connection established');
    };

    ws.current.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (Array.isArray(data)) {
          const processedMessages = data
            .map((msg) => normalizeMessage(msg, person, User))
            .filter(Boolean);

          setMessages((prev) => mergeMessages(prev, processedMessages));
        } else if (typeof data === 'object' && data !== null) {
          const normalized = normalizeMessage(data, person, User);
          if (normalized) {
            setMessages((prev) => mergeMessages(prev, [normalized]));
          }
        }
      } catch (e) {
        console.error('Error parsing websocket payload:', e);
      }
    };

    ws.current.onclose = () => {
      console.log('WebSocket connection closed');
    };

    ws.current.onerror = (err) => {
      console.error('WebSocket error', err);
    };

    return () => {
      if (ws.current) {
        ws.current.close();
        ws.current = null;
      }
    };
  }, [showModal, match?.ID, person?.id, person?.name, User?.id]);
  
  const sendMessage = () => {
    const trimmedInput = input.trim();

    if (ws.current && trimmedInput) {
      let date = new Date();
      const msgId = Date.now();
      
      // Send message to the server
      ws.current.send(JSON.stringify({ 
        message: trimmedInput, 
        who: User.id, 
        id: msgId, 
        time: date.toISOString() 
      }));
      
      // Add the message to our local state
      setMessages((prev) => mergeMessages(prev, [{ 
        message: trimmedInput, 
        who: 'Me', 
        id: msgId, 
        time: date.toISOString() 
      }]));
      
      // Reset the input field
      setInput('');
    }
  };
  
  // Function to trigger the Perfect Date AI recommendation
  const triggerPerfectDate = async () => {
    if (!match || !User || !jwt) return;
    
    try {
      // Show a sending indicator in the chat
      setMessages(prev => [...prev, {
        message: "Generating the perfect date suggestion...",
        who: "System",
        id: Date.now(),
        time: new Date().toLocaleTimeString()
      }]);
      
      // Call the PerfectDate API using the jwt prop
      const response = await fetch(apiUrl(`/perfectdate/${match.ID}`), {
        method: 'POST',
        headers: {
          'Authorization': jwt,
          'Content-Type': 'application/json'
        }
      });
      
      if (!response.ok) {
        throw new Error(`API returned ${response.status}: ${await response.text()}`);
      }
      
      // The AI message will be automatically added to the chat via the WebSocket
      // No need to manually add it to the messages state
      
    } catch (error) {
      console.error('Error getting perfect date suggestion:', error);
      
      // Show error message in chat
      setMessages(prev => [...prev, {
        message: `Failed to generate date suggestion: ${error.message}`,
        who: "Alert",
        id: Date.now(),
        time: new Date().toLocaleTimeString()
      }]);
    }
  };

  // Function to trigger the VibeChat AI conversation starter
  const triggerVibeChat = async () => {
    if (!match || !User || !jwt) return;
    
    try {
      // Show a sending indicator in the chat
      setMessages(prev => [...prev, {
        message: "Generating a new conversation starter...",
        who: "System",
        id: Date.now(),
        time: new Date().toLocaleTimeString()
      }]);
      
      // Call the VibeChat API using the jwt prop passed from App.jsx
      const response = await fetch(apiUrl(`/vibechat/${match.ID}`), {
        method: 'POST',
        headers: {
          'Authorization': jwt,
          'Content-Type': 'application/json'
        }
      });
      
      if (!response.ok) {
        throw new Error(`API returned ${response.status}: ${await response.text()}`);
      }
      
      // The AI message will be automatically added to the chat via the WebSocket
      // No need to manually add it to the messages state
      
    } catch (error) {
      console.error('Error triggering VibeChat:', error);
      
      // Show error message in chat
      setMessages(prev => [...prev, {
        message: `Failed to generate conversation starter: ${error.message}`,
        who: "Alert",
        id: Date.now(),
        time: new Date().toLocaleTimeString()
      }]);
    }
  };

  return (!person || !User || !match 
    ? <div>Sign in to see your matches and chat.</div>
    :
  <>
<style jsx>{`
        .gradient-bg {
  background: linear-gradient(135deg, #9333ea, #4f46e5);
  border: none;
  position: relative;
}

.gradient-bg::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #4f46e5, #9333ea);
  opacity: 0;
  transition: opacity 0.5s ease;
  z-index: -1;
  border-radius: inherit;
}

.gradient-bg:hover::before {
  opacity: 1;
}

.bg-purple {
  background-color: #9333ea;
}

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
      key={`openchatwith${person.name}-${match.ID}`}
      variant="outline-success"
      className="p-2 fs-5 m-2 text-right d-flex flex-row justify-content-between align-items-center"
      onClick={handleOpen}
    >
      <img
        src={person.profile.url ? person.profile.url : '/profile.svg'}
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
      key={`chatwindowwith${person.name}-${match.ID}`}
      show={showModal}
      onHide={handleClose}
      size="lg"
      aria-labelledby="contained-modal-title-vcenter"
      centered
      scrollable={true}
    >

      <Modal.Header closeButton>
        <Modal.Title>Chat with {person.name} <img
              src={person.profile.url ? person.profile.url : '/profile.svg'}
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
            {isLoading ? (
              <div className="d-flex justify-content-center align-items-center p-5">
                <div className="spinner-border text-primary" role="status">
                  <span className="visually-hidden">Loading...</span>
                </div>
              </div>
            ) : messages.length === 0 ? (
              <div className="text-center text-muted p-4">
                No messages yet. Start the conversation!
              </div>
            ) : messages.map((msg, idx) => (
              <ChatMessage
                key={msg.id || `${person.id}-${idx}`}
                msg={msg}
                currentUserId={User.id}
                otherUserId={person.id}
                otherUserName={person.name}
              />
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
        <div className="d-flex justify-content-around w-100 mt-2">
          <Button
            variant="secondary"
            onClick={() => triggerVibeChat()}
            className="mx-2"
          >
            Vibe Chat ❤️
          </Button>
          <Button
            variant="info"
            onClick={() => triggerPerfectDate()}
            className="mx-2"
          >
            Perfect Date 🍷
          </Button>
          <Button
            variant="secondary"
            onClick={handleClose}
            className="mx-2"
          >
            Close
          </Button>
        </div>
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
        src={person.profile.url ? person.profile.url : '/profile.svg'}
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
