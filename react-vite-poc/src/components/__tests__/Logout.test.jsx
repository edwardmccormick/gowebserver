import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import Logout from '../logout';

// Mock fetch for API calls
global.fetch = jest.fn();

// Mock alert
global.alert = jest.fn();

// Mock Lucide React icon
jest.mock('lucide-react', () => ({
  LogOut: () => <span data-testid="logout-icon">LogoutIcon</span>
}));

describe('Logout Component', () => {
  const mockSetLoggedInUser = jest.fn();
  const mockSetJWT = jest.fn();
  const mockJwt = 'test-jwt-token';
  
  beforeEach(() => {
    // Reset mocks before each test
    jest.clearAllMocks();
    global.fetch.mockReset();
    global.alert.mockReset();
  });
  
  test('renders logout button correctly', () => {
    render(
      <Logout 
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        jwt={mockJwt}
      />
    );
    
    // Check that button is rendered
    const logoutButton = screen.getByRole('button', { name: /logout/i });
    expect(logoutButton).toBeInTheDocument();
    
    // Check that icon is rendered
    const logoutIcon = screen.getByTestId('logout-icon');
    expect(logoutIcon).toBeInTheDocument();
  });
  
  test('handles successful logout', async () => {
    // Mock successful logout response
    fetch.mockResolvedValueOnce({
      ok: true
    });
    
    render(
      <Logout 
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        jwt={mockJwt}
      />
    );
    
    // Click the logout button
    fireEvent.click(screen.getByRole('button', { name: /logout/i }));
    
    // Check that fetch was called with correct parameters
    expect(fetch).toHaveBeenCalledWith('http://localhost:8080/logout', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': mockJwt
      }
    });
    
    // Wait for state updates
    await waitFor(() => {
      // Verify state was reset
      expect(mockSetLoggedInUser).toHaveBeenCalledWith(null);
      expect(mockSetJWT).toHaveBeenCalledWith(null);
      
      // Verify success message was shown
      expect(alert).toHaveBeenCalledWith('Logout successful!');
    });
  });
  
  test('handles logout failure', async () => {
    // Mock failed logout response
    fetch.mockResolvedValueOnce({
      ok: false
    });
    
    render(
      <Logout 
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        jwt={mockJwt}
      />
    );
    
    // Click the logout button
    fireEvent.click(screen.getByRole('button', { name: /logout/i }));
    
    // Wait for error message
    await waitFor(() => {
      expect(alert).toHaveBeenCalledWith('Invalid email or password.');
    });
    
    // Verify state was not reset
    expect(mockSetLoggedInUser).not.toHaveBeenCalled();
    expect(mockSetJWT).not.toHaveBeenCalled();
  });
  
  test('handles network error', async () => {
    // Mock network error
    const errorMessage = 'Network error';
    fetch.mockRejectedValueOnce(new Error(errorMessage));
    
    // Mock console.error to prevent test output pollution
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    
    render(
      <Logout 
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        jwt={mockJwt}
      />
    );
    
    // Click the logout button
    fireEvent.click(screen.getByRole('button', { name: /logout/i }));
    
    // Wait for error handling
    await waitFor(() => {
      // Verify error was logged
      expect(consoleErrorSpy).toHaveBeenCalledWith('Error during logout:', expect.any(Error));
      
      // Verify error message was shown
      expect(alert).toHaveBeenCalledWith('An error occurred while logging out. Guess you`re stuck with us!');
    });
    
    // Verify state was not reset
    expect(mockSetLoggedInUser).not.toHaveBeenCalled();
    expect(mockSetJWT).not.toHaveBeenCalled();
    
    // Cleanup
    consoleErrorSpy.mockRestore();
  });
});