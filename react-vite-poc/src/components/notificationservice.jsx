import { useState, useEffect } from 'react';

/**
 * NotificationService component that establishes a Server-Sent Events (SSE) connection
 * for receiving real-time notifications from the server.
 * 
 * @param {Object} props
 * @param {string} props.jwt - JSON Web Token for authentication
 * @param {Object} props.user - Current user object
 * @param {Function} props.onNewMessage - Callback for new message notifications
 * @param {Function} props.onNewMatch - Callback for new match notifications
 * @returns {null} This component doesn't render anything visible
 */
function NotificationService({ jwt, user, onNewMessage, onNewMatch }) {
  const [eventSource, setEventSource] = useState(null);
  const [reconnectAttempts, setReconnectAttempts] = useState(0);
  const [reconnectTimer, setReconnectTimer] = useState(null);
  const MAX_RECONNECT_ATTEMPTS = 5;
  const BASE_RECONNECT_DELAY = 2000; // Start with 2 seconds

  useEffect(() => {
    // Only establish connection if we have a user and JWT
    if (!user || !jwt) {
      return;
    }

    // Clean up any existing connection
    if (eventSource) {
      eventSource.close();
    }

    // Create new EventSource connection
    const connectSSE = () => {
      try {
        console.log(`Establishing SSE connection for user ${user.id}...`);
        
        // Create the SSE connection with JWT auth as a query parameter
        // because EventSource doesn't support custom headers
        const sseUrl = `http://localhost:8080/notifications/${user.id}?token=${encodeURIComponent(jwt)}`;
        console.log(`Creating SSE connection to: ${sseUrl}`);
        const sse = new EventSource(sseUrl, {
          withCredentials: true
        });

        // Set up event listeners
        sse.onopen = () => {
          console.log('SSE connection established');
          setReconnectAttempts(0); // Reset reconnect attempts on successful connection
          
          // Clear any reconnect timer
          if (reconnectTimer) {
            clearTimeout(reconnectTimer);
            setReconnectTimer(null);
          }
        };

        sse.onmessage = (event) => {
          try {
            console.log('SSE event received:', event);
            
            if (event.data === 'ping') {
              console.log('SSE heartbeat received');
              return;
            }
            
            // Parse the notification
            const notification = JSON.parse(event.data);
            console.log('Received notification:', notification);

            // Handle different notification types
            switch (notification.type) {
              case 'chat_message':
                if (onNewMessage) {
                  onNewMessage(notification.matchID, notification.count);
                }
                break;
              case 'match_update':
                if (onNewMatch) {
                  onNewMatch(notification.matchID);
                }
                break;
              default:
                console.log(`Unknown notification type: ${notification.type}`);
            }
          } catch (error) {
            console.error('Error processing SSE notification:', error);
          }
        };

        sse.onerror = (error) => {
          console.error('SSE connection error:', error);
          sse.close();
          
          // Implement reconnection strategy
          if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
            const delay = BASE_RECONNECT_DELAY * Math.pow(2, reconnectAttempts);
            console.log(`Reconnecting in ${delay}ms... (attempt ${reconnectAttempts + 1})`);
            
            const timer = setTimeout(() => {
              setReconnectAttempts(prev => prev + 1);
              connectSSE();
            }, delay);
            
            setReconnectTimer(timer);
          } else {
            console.error('Max reconnection attempts reached');
          }
        };

        setEventSource(sse);
      } catch (error) {
        console.error('Failed to establish SSE connection:', error);
      }
    };

    connectSSE();

    // Clean up on component unmount
    return () => {
      if (eventSource) {
        console.log('Closing SSE connection');
        eventSource.close();
      }
      
      if (reconnectTimer) {
        clearTimeout(reconnectTimer);
      }
    };
  }, [user, jwt]);

  // This component doesn't render anything visible
  return null;
}

export default NotificationService;