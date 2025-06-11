import { useEffect, useRef, useState } from 'react';

export default function WebSocketChat() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const ws = useRef(null);

  useEffect(() => {
    ws.current = new WebSocket('ws://localhost:8080/ws');

    ws.current.onmessage = (event) => {
      setMessages((prev) => [...prev, event.data]);
    };

    ws.current.onclose = () => {
      setMessages((prev) => [...prev, 'Connection closed']);
    };

    ws.current.onerror = (err) => {
      setMessages((prev) => [...prev, 'WebSocket error']);
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
    <div className="border p-3 m-3">
      <h3>Chat POC</h3>
      <div className='overflow-scroll h-50' style={{ minHeight: 100, background: '#f8f9fa', marginBottom: 10, padding: 10 }}>
        {messages.map((msg, idx) => (
          <div key={idx}>{msg}</div>
        ))}
      </div>
      <input
        value={input}
        onChange={e => setInput(e.target.value)}
        onKeyDown={e => e.key === 'Enter' && sendMessage()}
        placeholder="Type a message..."
        className="form-control mb-2"
      />
      <button onClick={sendMessage} className="btn btn-primary">Send</button>
    </div>
  );
}