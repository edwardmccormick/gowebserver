import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import App from '../../App';

// Mock react-router-dom
jest.mock('react-router-dom', () => ({
  BrowserRouter: ({ children }) => <div>{children}</div>,
  Routes: ({ children }) => <div>{children}</div>,
  Route: ({ element }) => element,
  useNavigate: () => jest.fn(),
  useLocation: () => ({ pathname: '/' }),
}));

// Mock components
jest.mock('../login', () => {
  return function MockLogin({ setLoggedInUser }) {
    return (
      <div data-testid="login-component">
        <button onClick={() => setLoggedInUser({ id: 1, name: 'Test User' })}>
          Mock Login
        </button>
      </div>
    );
  };
});

jest.mock('../logout', () => {
  return function MockLogout({ setLoggedInUser }) {
    return (
      <div data-testid="logout-component">
        <button onClick={() => setLoggedInUser(null)}>
          Mock Logout
        </button>
      </div>
    );
  };
});

jest.mock('../navbar', () => {
  return function MockNavbar({ loggedInUser }) {
    return (
      <div data-testid="navbar-component">
        {loggedInUser ? `Logged in as ${loggedInUser.name}` : 'Not logged in'}
      </div>
    );
  };
});

jest.mock('../notificationservice', () => {
  return function MockNotificationService({ user, jwt, onNewMessage }) {
    return (
      <div data-testid="notification-service">
        <button 
          onClick={() => onNewMessage(1, 3)}
          data-testid="trigger-notification"
        >
          Trigger Notification
        </button>
      </div>
    );
  };
});

jest.mock('../chatselect', () => {
  return function MockChatSelect({ User, matches, unreadNotifications, clearChatNotification }) {
    return (
      <div data-testid="chat-select">
        <span>User: {User?.name}</span>
        <span>Matches: {matches?.length || 0}</span>
        <span>Notifications: {JSON.stringify(unreadNotifications)}</span>
        <button 
          onClick={() => clearChatNotification(1)}
          data-testid="clear-notification"
        >
          Clear Notification
        </button>
      </div>
    );
  };
});

// Mock fetch
global.fetch = jest.fn();

describe('App Component', () => {
  beforeEach(() => {
    localStorage.clear();
    fetch.mockClear();
    
    // Mock matches response
    fetch.mockImplementation((url) => {
      if (url.includes('/matches')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve([
            { 
              ID: 1, 
              Offered: 1, 
              Accepted: 2, 
              UnreadOffered: 2,
              OfferedProfile: { id: 1, name: 'Test User' },
              AcceptedProfile: { id: 2, name: 'Match User' } 
            }
          ])
        });
      }
      
      // Default mock response
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve({})
      });
    });
  });

  test('renders without crashing', () => {
    render(<App />);
    expect(screen.getByTestId('navbar-component')).toBeInTheDocument();
  });

  test('starts with no logged in user', () => {
    render(<App />);
    expect(screen.getByText('Not logged in')).toBeInTheDocument();
  });

  test('handles user login', async () => {
    render(<App />);
    
    // Mock login button should be visible
    const loginButton = screen.getByText('Mock Login');
    expect(loginButton).toBeInTheDocument();
    
    // Simulate login
    fireEvent.click(loginButton);
    
    // Should update navbar
    await waitFor(() => {
      expect(screen.getByText('Logged in as Test User')).toBeInTheDocument();
    });
    
    // Should fetch matches after login
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/matches/1'),
      expect.anything()
    );
  });

  test('handles user logout', async () => {
    // Start with logged in user
    localStorage.setItem('jwt', 'test-token');
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'Test User' }));
    
    render(<App />);
    
    // Wait for logged in state
    await waitFor(() => {
      expect(screen.getByText('Logged in as Test User')).toBeInTheDocument();
    });
    
    // Logout component should be available
    const logoutButton = screen.getByText('Mock Logout');
    fireEvent.click(logoutButton);
    
    // Should update navbar
    await waitFor(() => {
      expect(screen.getByText('Not logged in')).toBeInTheDocument();
    });
    
    // Should clear localStorage
    expect(localStorage.getItem('jwt')).toBeNull();
    expect(localStorage.getItem('user')).toBeNull();
  });

  test('loads user from localStorage on mount', async () => {
    // Set user in localStorage before mount
    localStorage.setItem('jwt', 'test-token');
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'Test User' }));
    
    render(<App />);
    
    // Should show logged in state
    await waitFor(() => {
      expect(screen.getByText('Logged in as Test User')).toBeInTheDocument();
    });
    
    // Should fetch matches
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/matches/1'),
      expect.anything()
    );
  });

  test('handles notifications from NotificationService', async () => {
    // Start with logged in user
    localStorage.setItem('jwt', 'test-token');
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'Test User' }));
    
    render(<App />);
    
    // Wait for components to load
    await waitFor(() => {
      expect(screen.getByTestId('notification-service')).toBeInTheDocument();
      expect(screen.getByTestId('chat-select')).toBeInTheDocument();
    });
    
    // Trigger a notification
    fireEvent.click(screen.getByTestId('trigger-notification'));
    
    // ChatSelect should receive the notification
    await waitFor(() => {
      const notificationsText = screen.getByText(/Notifications:/);
      expect(notificationsText.textContent).toContain('{"1":3}');
    });
    
    // Clear notification
    fireEvent.click(screen.getByTestId('clear-notification'));
    
    // ChatSelect should no longer show the notification
    await waitFor(() => {
      const notificationsText = screen.getByText(/Notifications:/);
      expect(notificationsText.textContent).toContain('{}');
    });
  });

  test('initializes notifications from match data', async () => {
    // Start with logged in user
    localStorage.setItem('jwt', 'test-token');
    localStorage.setItem('user', JSON.stringify({ id: 1, name: 'Test User' }));
    
    render(<App />);
    
    // Wait for matches to load and notification to initialize
    await waitFor(() => {
      const notificationsText = screen.getByText(/Notifications:/);
      // Should initialize with the UnreadOffered count from the match
      expect(notificationsText.textContent).toContain('{"1":2}');
    });
  });
});