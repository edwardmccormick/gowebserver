import React from 'react';
import { render, waitFor } from '@testing-library/react';
import NotificationService from '../notificationservice';

// Mock EventSource since it's not available in the test environment
class MockEventSource {
  constructor(url) {
    this.url = url;
    this.onmessage = null;
    this.onerror = null;
    this.onopen = null;
    
    // Simulate connection
    setTimeout(() => {
      if (this.onopen) this.onopen();
    }, 0);
  }
  
  // Method to simulate receiving a message
  simulateMessage(data) {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(data) });
    }
  }
  
  close() {}
}

// Replace browser's EventSource with mock
global.EventSource = MockEventSource;

// Spy on console.log and console.error
const originalConsoleLog = console.log;
const originalConsoleError = console.error;
let consoleLogSpy;
let consoleErrorSpy;

beforeEach(() => {
  consoleLogSpy = jest.fn();
  consoleErrorSpy = jest.fn();
  console.log = consoleLogSpy;
  console.error = consoleErrorSpy;
});

afterEach(() => {
  console.log = originalConsoleLog;
  console.error = originalConsoleError;
});

describe('NotificationService', () => {
  test('establishes SSE connection with the correct URL', () => {
    const user = { id: 123 };
    const jwt = 'test-token';
    
    render(
      <NotificationService 
        user={user} 
        jwt={jwt}
        onNewMessage={jest.fn()}
        onNewMatch={jest.fn()}
      />
    );
    
    // Check that EventSource was constructed with the correct URL
    expect(consoleLogSpy).toHaveBeenCalledWith(
      expect.stringContaining(`Creating SSE connection to: http://localhost:8080/notifications/123?token=test-token`)
    );
  });
  
  test('calls onNewMessage when receiving a chat_message notification', async () => {
    const user = { id: 123 };
    const jwt = 'test-token';
    const onNewMessage = jest.fn();
    
    render(
      <NotificationService 
        user={user} 
        jwt={jwt}
        onNewMessage={onNewMessage}
        onNewMatch={jest.fn()}
      />
    );
    
    // Get the EventSource instance
    const eventSourceInstance = new EventSource();
    
    // Simulate receiving a chat message notification
    const testNotification = {
      type: 'chat_message',
      matchID: 456,
      count: 3,
      timestamp: new Date().toISOString()
    };
    
    eventSourceInstance.simulateMessage(testNotification);
    
    // Check that onNewMessage was called with the correct arguments
    await waitFor(() => {
      expect(onNewMessage).toHaveBeenCalledWith(456, 3);
    });
  });
  
  test('calls onNewMatch when receiving a match_update notification', async () => {
    const user = { id: 123 };
    const jwt = 'test-token';
    const onNewMatch = jest.fn();
    
    render(
      <NotificationService 
        user={user} 
        jwt={jwt}
        onNewMessage={jest.fn()}
        onNewMatch={onNewMatch}
      />
    );
    
    // Get the EventSource instance
    const eventSourceInstance = new EventSource();
    
    // Simulate receiving a match update notification
    const testNotification = {
      type: 'match_update',
      matchID: 789,
      timestamp: new Date().toISOString()
    };
    
    eventSourceInstance.simulateMessage(testNotification);
    
    // Check that onNewMatch was called with the correct arguments
    await waitFor(() => {
      expect(onNewMatch).toHaveBeenCalledWith(789);
    });
  });
  
  test('handles ping messages correctly', async () => {
    const user = { id: 123 };
    const jwt = 'test-token';
    const onNewMessage = jest.fn();
    const onNewMatch = jest.fn();
    
    render(
      <NotificationService 
        user={user} 
        jwt={jwt}
        onNewMessage={onNewMessage}
        onNewMatch={onNewMatch}
      />
    );
    
    // Get the EventSource instance
    const eventSourceInstance = new EventSource();
    
    // Simulate receiving a ping message
    eventSourceInstance.onmessage({ data: 'ping' });
    
    // Wait to ensure callbacks are not called
    await waitFor(() => {
      expect(consoleLogSpy).toHaveBeenCalledWith('SSE heartbeat received');
      expect(onNewMessage).not.toHaveBeenCalled();
      expect(onNewMatch).not.toHaveBeenCalled();
    });
  });
  
  test('does not establish connection without user or jwt', () => {
    // Render without user
    render(
      <NotificationService 
        jwt="test-token"
        onNewMessage={jest.fn()}
        onNewMatch={jest.fn()}
      />
    );
    
    // Check that no connection was attempted
    expect(consoleLogSpy).not.toHaveBeenCalledWith(
      expect.stringContaining('Creating SSE connection')
    );
    
    // Render without JWT
    render(
      <NotificationService 
        user={{ id: 123 }}
        onNewMessage={jest.fn()}
        onNewMatch={jest.fn()}
      />
    );
    
    // Check that no connection was attempted
    expect(consoleLogSpy).not.toHaveBeenCalledWith(
      expect.stringContaining('Creating SSE connection')
    );
  });
  
  test('handles JSON parsing errors gracefully', async () => {
    const user = { id: 123 };
    const jwt = 'test-token';
    const onNewMessage = jest.fn();
    
    render(
      <NotificationService 
        user={user} 
        jwt={jwt}
        onNewMessage={onNewMessage}
        onNewMatch={jest.fn()}
      />
    );
    
    // Get the EventSource instance
    const eventSourceInstance = new EventSource();
    
    // Simulate receiving an invalid JSON message
    eventSourceInstance.onmessage({ data: 'not-valid-json' });
    
    // Check that error is logged and callbacks are not called
    await waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalledWith(
        'Error processing SSE notification:',
        expect.any(Error)
      );
      expect(onNewMessage).not.toHaveBeenCalled();
    });
  });
  
  test('handles unknown notification types', async () => {
    const user = { id: 123 };
    const jwt = 'test-token';
    const onNewMessage = jest.fn();
    const onNewMatch = jest.fn();
    
    render(
      <NotificationService 
        user={user} 
        jwt={jwt}
        onNewMessage={onNewMessage}
        onNewMatch={onNewMatch}
      />
    );
    
    // Get the EventSource instance
    const eventSourceInstance = new EventSource();
    
    // Simulate receiving an unknown notification type
    const testNotification = {
      type: 'unknown_type',
      matchID: 456,
      timestamp: new Date().toISOString()
    };
    
    eventSourceInstance.simulateMessage(testNotification);
    
    // Check that warning is logged and callbacks are not called
    await waitFor(() => {
      expect(consoleLogSpy).toHaveBeenCalledWith('Unknown notification type: unknown_type');
      expect(onNewMessage).not.toHaveBeenCalled();
      expect(onNewMatch).not.toHaveBeenCalled();
    });
  });
});