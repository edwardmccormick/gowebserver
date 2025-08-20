import { useEffect, useRef, useState } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Modal from 'react-bootstrap/Modal';
import Badge from 'react-bootstrap/Badge';

function convertISODateToLocal(dateString) {
  const date = new Date(dateString);


  const localTime = date.toLocaleTimeString();
  return localTime; // Output: "09:03:56 PM"

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
      const response = await fetch(`http://localhost:8080/chat/markread/${match.ID}`, {
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
  
  // Track if the chat history has been loaded
  const [historyLoaded, setHistoryLoaded] = useState(false);
  // Track loading state to show a spinner while fetching
  const [isLoading, setIsLoading] = useState(false);
  // Track message IDs to prevent duplicates
  const [messageIds, setMessageIds] = useState(new Set());
  
    // Scroll to the bottom whenever messages change
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);


  // Fetch chat history from the server when needed
  const fetchChatHistory = async () => {
    if (!match || !User || historyLoaded) return;
    
    setIsLoading(true);
    try {

      const response = await fetch(`http://localhost:8080/chat/messages/${match.ID}`, {
        headers: {
          'Authorization': jwt,
          'Content-Type': 'application/json'
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        
        // Check if the response contains the messages array directly or within pagination
        const chatMessages = data.messages || data;
        
        // Process the messages to ensure proper display names
        const processedMessages = chatMessages.map(msg => {
          if (msg.who == person.id) {
            return {...msg, who: person.name};
          } else if (msg.who == User.id) {
            return {...msg, who: 'Me'};
          } else if (msg.who === 0) {
            return {...msg, who: 'AI Introduction'};
          }
          return msg;
        });
        
        // Track message IDs to prevent duplicates
        const newIds = new Set(messageIds);
        processedMessages.forEach(msg => {
          if (msg.id) newIds.add(msg.id);
        });
        setMessageIds(newIds);
        
        setMessages(processedMessages);
        setHistoryLoaded(true);
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
    if (showModal && !historyLoaded) {
      // Reset messageIds when starting a new chat session
      if (showModal) {
        setMessageIds(new Set());
      }
      fetchChatHistory();
      // Mark messages as read when opening the chat
      markMessagesAsRead();
    }
  }, [showModal, historyLoaded]);

  useEffect(() => {
    if (!match || !person || !User || !showModal) {
      return;
    }
    
    ws.current = new WebSocket(`ws://localhost:8080/ws?id=${match.ID}&user_id=${User.id}`);

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
          // Only append new messages if history hasn't been loaded yet
          if (!historyLoaded) {
            setHistoryLoaded(true);
            
            // Filter out any duplicate messages
            const uniqueMessages = data.filter(msg => !msg.id || !messageIds.has(msg.id));
            
            // Process message senders
            const processedMessages = uniqueMessages.map(msg => {
              if (msg.who == person.id) {
                return {...msg, who: person.name};
              } else if (msg.who == User.id) {
                return {...msg, who: 'Me'};
              } else if (msg.who === 0) {
                return {...msg, who: 'AI Introduction'};
              }
              return msg;
            });
            
            // Track new message IDs
            const newIds = new Set(messageIds);
            processedMessages.forEach(msg => {
              if (msg.id) newIds.add(msg.id);
            });
            setMessageIds(newIds);
            
            // Add unique messages to the chat
            if (processedMessages.length > 0) {
              setMessages((prev) => [...prev, ...processedMessages]);
            }
          }
        } else if (typeof data === 'object' && data !== null) {
          // Only add the message if it's not a duplicate
          if (!data.id || !messageIds.has(data.id)) {
            // Handle different types of senders
            if (data.who == person.id) {
              data.who = person.name;
            } else if (data.who == User.id) {
              data.who = 'Me';
            } else if (data.who === 0) {
              data.who = 'AI Introduction';
            }
            
            // Track the new message ID
            if (data.id) {
              setMessageIds(prev => new Set(prev).add(data.id));
            }
            
            setMessages((prev) => [...prev, data]);
          }
        } else {
          // If it's not JSON, treat as a status message
          const statusMsgId = Date.now();
          setMessages((prev) => [...prev, { message: event.data, who: 'System', id: statusMsgId, time: statusMsgId }]);
        }
      } catch (e) {
        // Not JSON, treat as a status message
        const errorMsgId = Date.now();
        setMessages((prev) => [...prev, { message: event.data, who: 'System', id: errorMsgId, time: errorMsgId }]);
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
      const msgId = Date.now();
      
      // Send message to the server
      ws.current.send(JSON.stringify({ 
        message: input, 
        who: User.id, 
        id: msgId, 
        time: date.toISOString() 
      }));
      
      // Add the message to our local state
      setMessages((prev) => [...prev, { 
        message: input, 
        who: 'Me', 
        id: msgId, 
        time: date.toLocaleTimeString() 
      }]);
      
      // Track the message ID to prevent duplicates
      setMessageIds(prev => new Set(prev).add(msgId));
      
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
      const response = await fetch(`http://localhost:8080/perfectdate/${match.ID}`, {
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
      const response = await fetch(`http://localhost:8080/vibechat/${match.ID}`, {
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
              <div className={
                      msg.who === 'System'
                      ? 'text-center bg-warning text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Alert'
                      ? 'text-center bg-danger text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Connection'
                      ? 'text-center bg-success text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Admin' || msg.who === 'Moderator'
                      ? 'align-self-center bg-info text-white rounded m-2 p-2 flex-item'
                      : msg.who === "AI Host" || msg.who === "AI Introduction"
                      ? 'align-self-center bg-purple gradient-bg text-white rounded m-2 p-2 flex-item'
                      : msg.who == User.id || msg.who === 'Me'
                      ? 'd-flex flex-row justify-content-end align-items-center'
                      : 'd-flex flex-row justify-content-start align-items-center'
                  }
                  key={`${person.id}${idx}`}>
                  {msg.who == User.id || msg.who === "Me" ?<span className="text-body-tertiary fs-6">{msg.time}</span> : null}
                  <span className={
                    msg.who === 'System'
                      ? 'text-center bg-warning text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Alert'
                      ? 'text-center bg-danger text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Connection'
                      ? 'text-center bg-success text-white rounded m-2 p-2 flex-item'
                      : msg.who === 'Admin' || msg.who === 'Moderator'
                      ? 'align-self-center bg-info text-white rounded m-2 p-2 flex-item'
                      : msg.who === "AI Host" || msg.who === "AI Introduction"
                      ? 'align-self-center bg-purple gradient-bg text-white rounded m-2 p-1 flex-item'
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