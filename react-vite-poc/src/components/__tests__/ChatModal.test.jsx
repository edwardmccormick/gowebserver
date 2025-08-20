import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import ChatModal from '../chatmodal';

// Mock fetch for API calls
global.fetch = jest.fn();

// Mock WebSocket
class MockWebSocket {
  constructor(url) {
    this.url = url;
    this.onmessage = null;
    this.onclose = null;
    this.onopen = null;
    this.readyState = 1; // OPEN

    // Simulate connection
    setTimeout(() => {
      if (this.onopen) this.onopen({ target: this });
    }, 0);
  }

  send(data) {
    // Mock sending message
    this.lastMessage = data;
  }

  close() {
    if (this.onclose) {
      this.onclose({ target: this });
    }
  }

  // Helper method for tests to simulate incoming messages
  simulateMessage(data) {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(data) });
    }
  }
}

// Replace the browser's WebSocket with our mock
global.WebSocket = MockWebSocket;

describe('ChatModal Component', () => {
  const mockMatch = {
    ID: 1,
    Offered: 1,
    Accepted: 2,
    UnreadOffered: 3,
    OfferedProfile: { id: 1, name: 'Current User', profile: { url: 'user1.jpg' } },
    AcceptedProfile: { id: 2, name: 'Match User', profile: { url: 'user2.jpg' } }
  };

  const mockPerson = {
    id: 2,
    name: 'Match User',
    profile: { url: 'user2.jpg' }
  };

  const mockUser = {
    id: 1,
    name: 'Current User'
  };

  const mockJwt = 'test-token';
  const mockClearChatNotification = jest.fn();

  beforeEach(() => {
    fetch.mockClear();
    mockClearChatNotification.mockClear();
    
    // Mock fetch for chat history
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([
        { id: '1', message: 'Hello', who: 1, time: '2023-01-01T12:00:00Z' },
        { id: '2', message: 'Hi there!', who: 2, time: '2023-01-01T12:01:00Z' }
      ])
    });
  });

  test('renders button with unread count', () => {
    render(
      <ChatModal
        match={mockMatch}
        person={mockPerson}
        User={mockUser}
        jwt={mockJwt}
        unreadmessages={3}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    const button = screen.getByText(/chat with match user/i);
    expect(button).toBeInTheDocument();
    
    const badge = screen.getByText('3');
    expect(badge).toBeInTheDocument();
  });

  test('opens modal and loads chat history when clicked', async () => {
    render(
      <ChatModal
        match={mockMatch}
        person={mockPerson}
        User={mockUser}
        jwt={mockJwt}
        unreadmessages={3}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    // Click the chat button
    fireEvent.click(screen.getByText(/chat with match user/i));
    
    // Modal should be visible
    expect(screen.getByText(/chat with match user/i, { selector: '.modal-title' })).toBeInTheDocument();
    
    // Wait for chat history to load
    await waitFor(() => {
      expect(screen.getByText('Hello')).toBeInTheDocument();
      expect(screen.getByText('Hi there!')).toBeInTheDocument();
    });
    
    // Check that notification was cleared
    expect(mockClearChatNotification).toHaveBeenCalledWith(1);
    
    // Check that fetch was called correctly for chat history
    expect(fetch).toHaveBeenCalledWith(
      'http://localhost:8080/chatmessages/1',
      expect.objectContaining({
        headers: expect.objectContaining({
          'Authorization': 'test-token'
        })
      })
    );
  });

  test('sends message when submit button clicked', async () => {
    render(
      <ChatModal
        match={mockMatch}
        person={mockPerson}
        User={mockUser}
        jwt={mockJwt}
        unreadmessages={0}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    // Open the modal
    fireEvent.click(screen.getByText(/chat with match user/i));
    
    // Wait for chat history to load
    await waitFor(() => {
      expect(screen.getByText('Hello')).toBeInTheDocument();
    });
    
    // Type a message
    const input = screen.getByPlaceholderText(/type a message/i);
    fireEvent.change(input, { target: { value: 'New test message' } });
    
    // Submit the message
    const sendButton = screen.getByRole('button', { name: /send/i });
    fireEvent.click(sendButton);
    
    // Check WebSocket message was sent
    const ws = global.WebSocket.mock?.instances[0];
    await waitFor(() => {
      expect(JSON.parse(ws.lastMessage)).toEqual(
        expect.objectContaining({
          message: 'New test message',
          who: 1,
          matchID: 1
        })
      );
    });
    
    // Check input was cleared
    expect(input.value).toBe('');
  });

  test('receives and displays new messages', async () => {
    render(
      <ChatModal
        match={mockMatch}
        person={mockPerson}
        User={mockUser}
        jwt={mockJwt}
        unreadmessages={0}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    // Open the modal
    fireEvent.click(screen.getByText(/chat with match user/i));
    
    // Wait for chat history to load
    await waitFor(() => {
      expect(screen.getByText('Hello')).toBeInTheDocument();
    });
    
    // Get WebSocket instance and simulate incoming message
    const ws = global.WebSocket.mock?.instances[0];
    ws.simulateMessage({
      id: '3',
      message: 'This is a new message',
      who: 2,
      time: new Date().toISOString()
    });
    
    // Check that new message appears
    await waitFor(() => {
      expect(screen.getByText('This is a new message')).toBeInTheDocument();
    });
  });

  test('handles WebSocket errors', async () => {
    // Mock console.error to prevent test output noise
    const originalError = console.error;
    console.error = jest.fn();
    
    render(
      <ChatModal
        match={mockMatch}
        person={mockPerson}
        User={mockUser}
        jwt={mockJwt}
        unreadmessages={0}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    // Open the modal
    fireEvent.click(screen.getByText(/chat with match user/i));
    
    // Wait for chat history to load
    await waitFor(() => {
      expect(screen.getByText('Hello')).toBeInTheDocument();
    });
    
    // Get WebSocket instance and simulate error
    const ws = global.WebSocket.mock?.instances[0];
    ws.onerror({ error: new Error('WebSocket error') });
    
    // Verify error was logged
    expect(console.error).toHaveBeenCalled();
    
    // Restore console.error
    console.error = originalError;
  });

  test('closes WebSocket when modal is closed', async () => {
    render(
      <ChatModal
        match={mockMatch}
        person={mockPerson}
        User={mockUser}
        jwt={mockJwt}
        unreadmessages={0}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    // Open the modal
    fireEvent.click(screen.getByText(/chat with match user/i));
    
    // Wait for chat history to load
    await waitFor(() => {
      expect(screen.getByText('Hello')).toBeInTheDocument();
    });
    
    // Get WebSocket instance and spy on close method
    const ws = global.WebSocket.mock?.instances[0];
    const closeSpy = jest.spyOn(ws, 'close');
    
    // Close the modal
    fireEvent.click(screen.getByLabelText(/close/i));
    
    // Check that WebSocket was closed
    expect(closeSpy).toHaveBeenCalled();
  });
});