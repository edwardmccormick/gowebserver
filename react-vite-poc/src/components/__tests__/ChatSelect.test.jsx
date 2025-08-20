import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import ChatSelect from '../chatselect';

// Mock ChatModal component to avoid testing its internals
jest.mock('../chatmodal', () => {
  return function MockChatModal(props) {
    return (
      <div data-testid="chat-modal" onClick={() => props.clearChatNotification?.(props.match?.ID)}>
        Chat with {props.person?.name}
        {props.unreadmessages > 0 && (
          <span data-testid="unread-badge">{props.unreadmessages}</span>
        )}
      </div>
    );
  };
});

// Mock the date conversion function to return a consistent value
jest.mock('react-bootstrap/esm/Button', () => {
  return function MockButton(props) {
    return <button {...props}>{props.children}</button>;
  };
});

jest.mock('react-bootstrap/Offcanvas', () => {
  return function MockOffcanvas(props) {
    return (
      <div data-testid="offcanvas" className={props.show ? 'show' : ''}>
        <div data-testid="offcanvas-header">
          <h5>{props.title}</h5>
          <button onClick={props.onHide}>Close</button>
        </div>
        <div data-testid="offcanvas-body">{props.children}</div>
      </div>
    );
  };
});

describe('ChatSelect Component', () => {
  const mockUser = { id: 1, name: 'Test User' };
  
  const mockMatches = [
    {
      ID: 1,
      Offered: 1,
      Accepted: 2,
      UnreadOffered: 3,
      UnreadAccepted: 0,
      OfferedProfile: { id: 1, name: 'Test User', profile: { url: 'test1.jpg' } },
      AcceptedProfile: { id: 2, name: 'Match 1', profile: { url: 'test2.jpg' } },
    },
    {
      ID: 2,
      Offered: 2,
      Accepted: 1,
      UnreadOffered: 0,
      UnreadAccepted: 5,
      OfferedProfile: { id: 2, name: 'Match 2', profile: { url: 'test3.jpg' } },
      AcceptedProfile: { id: 1, name: 'Test User', profile: { url: 'test1.jpg' } },
    },
  ];
  
  const mockPendings = [];
  const mockOffereds = [];
  const mockJwt = 'test-jwt-token';
  const mockUnreadNotifications = { 1: 3, 2: 5 };
  const mockClearChatNotification = jest.fn();
  
  test('renders chat toggle button when user and matches exist', () => {
    render(
      <ChatSelect
        User={mockUser}
        matches={mockMatches}
        pendings={mockPendings}
        offereds={mockOffereds}
        jwt={mockJwt}
        unreadNotifications={mockUnreadNotifications}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    const toggleButton = screen.getByText('Chat Selector');
    expect(toggleButton).toBeInTheDocument();
  });
  
  test('opens offcanvas when toggle button is clicked', () => {
    render(
      <ChatSelect
        User={mockUser}
        matches={mockMatches}
        pendings={mockPendings}
        offereds={mockOffereds}
        jwt={mockJwt}
        unreadNotifications={mockUnreadNotifications}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    const toggleButton = screen.getByText('Chat Selector');
    fireEvent.click(toggleButton);
    
    const offcanvas = screen.getByTestId('offcanvas');
    expect(offcanvas).toHaveClass('show');
  });
  
  test('displays correct unread message counts for matches', () => {
    render(
      <ChatSelect
        User={mockUser}
        matches={mockMatches}
        pendings={mockPendings}
        offereds={mockOffereds}
        jwt={mockJwt}
        unreadNotifications={mockUnreadNotifications}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    const toggleButton = screen.getByText('Chat Selector');
    fireEvent.click(toggleButton);
    
    // Find all chat modals
    const chatModals = screen.getAllByTestId('chat-modal');
    expect(chatModals).toHaveLength(2);
    
    // Check for unread badges
    const unreadBadges = screen.getAllByTestId('unread-badge');
    expect(unreadBadges).toHaveLength(2);
    expect(unreadBadges[0].textContent).toBe('3');
    expect(unreadBadges[1].textContent).toBe('5');
  });
  
  test('calls clearChatNotification when a chat modal is clicked', () => {
    render(
      <ChatSelect
        User={mockUser}
        matches={mockMatches}
        pendings={mockPendings}
        offereds={mockOffereds}
        jwt={mockJwt}
        unreadNotifications={mockUnreadNotifications}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    const toggleButton = screen.getByText('Chat Selector');
    fireEvent.click(toggleButton);
    
    // Find the first chat modal and click it
    const chatModals = screen.getAllByTestId('chat-modal');
    fireEvent.click(chatModals[0]);
    
    // Check that clearChatNotification was called with the correct match ID
    expect(mockClearChatNotification).toHaveBeenCalledWith(1);
  });
  
  test('returns null when user or matches are not provided', () => {
    const { container: container1 } = render(
      <ChatSelect
        User={null}
        matches={mockMatches}
        pendings={mockPendings}
        offereds={mockOffereds}
        jwt={mockJwt}
        unreadNotifications={mockUnreadNotifications}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    expect(container1.firstChild).toBeNull();
    
    const { container: container2 } = render(
      <ChatSelect
        User={mockUser}
        matches={null}
        pendings={mockPendings}
        offereds={mockOffereds}
        jwt={mockJwt}
        unreadNotifications={mockUnreadNotifications}
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    expect(container2.firstChild).toBeNull();
  });
  
  test('correctly calculates unread message counts using countUndeliveredMessages', () => {
    // This is testing the internal function but we can verify its behavior
    // by examining the rendered output
    render(
      <ChatSelect
        User={mockUser}
        matches={mockMatches}
        pendings={mockPendings}
        offereds={mockOffereds}
        jwt={mockJwt}
        unreadNotifications={{}} // Empty notifications to test internal counting
        clearChatNotification={mockClearChatNotification}
      />
    );
    
    const toggleButton = screen.getByText('Chat Selector');
    fireEvent.click(toggleButton);
    
    // Find all unread badges
    const unreadBadges = screen.getAllByTestId('unread-badge');
    expect(unreadBadges).toHaveLength(2);
    
    // First match: User is Offered, should use UnreadOffered (3)
    expect(unreadBadges[0].textContent).toBe('3');
    
    // Second match: User is Accepted, should use UnreadAccepted (5)
    expect(unreadBadges[1].textContent).toBe('5');
  });
});