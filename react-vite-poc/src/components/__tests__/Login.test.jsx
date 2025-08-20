import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import SignIn from '../login';

// Mock fetch for API calls
global.fetch = jest.fn();

// Mock alert
global.alert = jest.fn();

describe('SignIn Component', () => {
  const mockSetLoggedInUser = jest.fn();
  const mockSetJWT = jest.fn();
  const mockSetUploadUrls = jest.fn();
  const mockSetUploadProfileUrls = jest.fn();
  const mockRefreshMatches = jest.fn();
  
  beforeEach(() => {
    fetch.mockClear();
    alert.mockClear();
    mockSetLoggedInUser.mockClear();
    mockSetJWT.mockClear();
    mockSetUploadUrls.mockClear();
    mockSetUploadProfileUrls.mockClear();
    mockRefreshMatches.mockClear();
  });

  test('renders login form correctly', () => {
    render(
      <SignIn
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
        refreshMatches={mockRefreshMatches}
      />
    );
    
    // Check form elements
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
  });

  test('handles input changes', () => {
    render(
      <SignIn
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
        refreshMatches={mockRefreshMatches}
      />
    );
    
    const emailInput = screen.getByLabelText(/email/i);
    const passwordInput = screen.getByLabelText(/password/i);
    
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    
    expect(emailInput.value).toBe('test@example.com');
    expect(passwordInput.value).toBe('password123');
  });

  test('submits the form with credentials and handles success', async () => {
    // Mock successful login response
    const mockResponse = {
      person: { id: 1, name: 'Test User' },
      token: 'test-token',
      upload_urls: ['url1', 'url2'],
      profile_upload_urls: ['profile-url1', 'profile-url2']
    };
    
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    });
    
    // Need to mock console.log for the profile_upload_urls logging
    const consoleSpy = jest.spyOn(console, 'log').mockImplementation();
    
    render(
      <SignIn
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
        refreshMatches={mockRefreshMatches}
      />
    );
    
    // Fill in the form
    fireEvent.change(screen.getByLabelText(/email/i), { 
      target: { value: 'test@example.com' } 
    });
    fireEvent.change(screen.getByLabelText(/password/i), { 
      target: { value: 'password123' } 
    });
    
    // Submit the form
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    
    // Check that fetch was called correctly
    expect(fetch).toHaveBeenCalledWith('http://localhost:8080/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        email: 'test@example.com', 
        password: 'password123' 
      })
    });
    
    // Wait for state updates
    await waitFor(() => {
      // Check that setter functions were called with correct values
      expect(mockSetLoggedInUser).toHaveBeenCalledWith(mockResponse.person);
      expect(mockSetJWT).toHaveBeenCalledWith(mockResponse.token);
      expect(mockSetUploadUrls).toHaveBeenCalledWith(mockResponse.upload_urls);
      expect(mockSetUploadProfileUrls).toHaveBeenCalledWith(mockResponse.profile_upload_urls);
      
      // Check that alert was called
      expect(alert).toHaveBeenCalledWith('Login successful!');
    });
    
    // Check that refreshMatches is called after timeout
    jest.advanceTimersByTime(500);
    expect(mockRefreshMatches).toHaveBeenCalled();
    
    // Clean up
    consoleSpy.mockRestore();
  });

  test('handles login error', async () => {
    // Mock failed login response
    fetch.mockResolvedValueOnce({
      ok: false
    });
    
    render(
      <SignIn
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
        refreshMatches={mockRefreshMatches}
      />
    );
    
    // Fill in the form
    fireEvent.change(screen.getByLabelText(/email/i), { 
      target: { value: 'wrong@example.com' } 
    });
    fireEvent.change(screen.getByLabelText(/password/i), { 
      target: { value: 'wrongpassword' } 
    });
    
    // Submit the form
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    
    // Check alert was called with error message
    await waitFor(() => {
      expect(alert).toHaveBeenCalledWith('Invalid email or password.');
    });
    
    // Check that state setter functions were not called
    expect(mockSetLoggedInUser).not.toHaveBeenCalled();
    expect(mockSetJWT).not.toHaveBeenCalled();
    expect(mockSetUploadUrls).not.toHaveBeenCalled();
    expect(mockSetUploadProfileUrls).not.toHaveBeenCalled();
    expect(mockRefreshMatches).not.toHaveBeenCalled();
  });

  test('handles network error', async () => {
    // Mock network error
    const errorMessage = 'Network error';
    fetch.mockRejectedValueOnce(new Error(errorMessage));
    
    // Mock console.error to prevent test output pollution
    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();
    
    render(
      <SignIn
        setLoggedInUser={mockSetLoggedInUser}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
        refreshMatches={mockRefreshMatches}
      />
    );
    
    // Fill in the form
    fireEvent.change(screen.getByLabelText(/email/i), { 
      target: { value: 'test@example.com' } 
    });
    fireEvent.change(screen.getByLabelText(/password/i), { 
      target: { value: 'password123' } 
    });
    
    // Submit the form
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    
    // Check console.error and alert were called
    await waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalledWith('Error during login:', expect.any(Error));
      expect(alert).toHaveBeenCalledWith('An error occurred while logging in.');
    });
    
    // Check that state setter functions were not called
    expect(mockSetLoggedInUser).not.toHaveBeenCalled();
    expect(mockSetJWT).not.toHaveBeenCalled();
    expect(mockSetUploadUrls).not.toHaveBeenCalled();
    expect(mockSetUploadProfileUrls).not.toHaveBeenCalled();
    
    // Cleanup
    consoleErrorSpy.mockRestore();
  });
  
  // Setup fake timers for testing timeouts
  beforeAll(() => {
    jest.useFakeTimers();
  });
  
  afterAll(() => {
    jest.useRealTimers();
  });
});