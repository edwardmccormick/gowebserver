import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { SignUp } from '../signup';

// Mock fetch for API calls
global.fetch = jest.fn();

// Mock alert
global.alert = jest.fn();

// Mock console.error and console.log
const originalConsoleError = console.error;
const originalConsoleLog = console.log;

describe('Signup Component', () => {
  // Mock props
  const mockSetPendingID = jest.fn();
  const mockSetJWT = jest.fn();
  const mockSetUploadUrls = jest.fn();
  const mockSetUploadProfileUrls = jest.fn();
  
  beforeEach(() => {
    fetch.mockClear();
    alert.mockClear();
    mockSetPendingID.mockClear();
    mockSetJWT.mockClear();
    mockSetUploadUrls.mockClear();
    mockSetUploadProfileUrls.mockClear();
    
    // Mock console methods
    console.error = jest.fn();
    console.log = jest.fn();
  });
  
  afterAll(() => {
    // Restore console methods
    console.error = originalConsoleError;
    console.log = originalConsoleLog;
  });

  test('renders signup form correctly', () => {
    render(
      <SignUp
        setPendingID={mockSetPendingID}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
      />
    );
    
    // Check form elements
    expect(screen.getByLabelText(/your email address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /sign up/i })).toBeInTheDocument();
  });

  test('handles input changes', () => {
    render(
      <SignUp
        setPendingID={mockSetPendingID}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
      />
    );
    
    const emailInput = screen.getByLabelText(/your email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    
    expect(emailInput.value).toBe('test@example.com');
    expect(passwordInput.value).toBe('password123');
  });

  // Note: The current component doesn't validate passwords matching or email format

  test('submits the form with valid data', async () => {
    // Mock successful signup response
    const mockResponse = {
      id: 1,
      token: 'test-token',
      upload_urls: ['url1', 'url2'],
      profile_upload_urls: ['profile-url1', 'profile-url2']
    };
    
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    });
    
    render(
      <SignUp
        setPendingID={mockSetPendingID}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
      />
    );
    
    const emailInput = screen.getByLabelText(/your email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const signupButton = screen.getByRole('button', { name: /sign up/i });
    
    // Fill in the form with valid data
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    
    // Submit the form
    fireEvent.click(signupButton);
    
    // Check that fetch was called correctly
    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/signup', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          email: 'test@example.com', 
          password: 'password123' 
        })
      });
      
      // Check that state setter functions were called with correct values
      expect(mockSetPendingID).toHaveBeenCalledWith(mockResponse.id);
      expect(mockSetJWT).toHaveBeenCalledWith(mockResponse.token);
      expect(mockSetUploadUrls).toHaveBeenCalledWith(mockResponse.upload_urls);
      expect(mockSetUploadProfileUrls).toHaveBeenCalledWith(mockResponse.profile_upload_urls);
      
      // Check that success alert was shown
      expect(alert).toHaveBeenCalledWith('Signup successful');
    });
  });

  test('handles signup error', async () => {
    // Mock failed signup response
    fetch.mockResolvedValueOnce({
      ok: false,
      status: 400
    });
    
    render(
      <SignUp
        setPendingID={mockSetPendingID}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
      />
    );
    
    const emailInput = screen.getByLabelText(/your email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const signupButton = screen.getByRole('button', { name: /sign up/i });
    
    // Fill in the form
    fireEvent.change(emailInput, { target: { value: 'existing@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    
    // Submit the form
    fireEvent.click(signupButton);
    
    // Check error alert appears
    await waitFor(() => {
      expect(alert).toHaveBeenCalledWith('We already have a user with that email address');
    });
    
    // Verify state setter functions were not called
    expect(mockSetPendingID).not.toHaveBeenCalled();
    expect(mockSetJWT).not.toHaveBeenCalled();
    expect(mockSetUploadUrls).not.toHaveBeenCalled();
    expect(mockSetUploadProfileUrls).not.toHaveBeenCalled();
  });

  test('handles network error', async () => {
    // Mock network error
    const errorMessage = 'Network error';
    fetch.mockRejectedValueOnce(new Error(errorMessage));
    
    render(
      <SignUp
        setPendingID={mockSetPendingID}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
      />
    );
    
    const emailInput = screen.getByLabelText(/your email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const signupButton = screen.getByRole('button', { name: /sign up/i });
    
    // Fill in the form
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    
    // Submit the form
    fireEvent.click(signupButton);
    
    // Check error handling
    await waitFor(() => {
      expect(console.error).toHaveBeenCalledWith('Error submitting form:', expect.any(Error));
      expect(alert).toHaveBeenCalledWith('An error occurred while submitting the form.');
    });
    
    // Verify state setter functions were not called
    expect(mockSetPendingID).not.toHaveBeenCalled();
    expect(mockSetJWT).not.toHaveBeenCalled();
    expect(mockSetUploadUrls).not.toHaveBeenCalled();
    expect(mockSetUploadProfileUrls).not.toHaveBeenCalled();
  });
  
  test('handles payload correctly', async () => {
    // This test specifically verifies the payload structure
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({})
    });
    
    render(
      <SignUp
        setPendingID={mockSetPendingID}
        setJWT={mockSetJWT}
        setUploadUrls={mockSetUploadUrls}
        setUploadProfileUrls={mockSetUploadProfileUrls}
      />
    );
    
    const emailInput = screen.getByLabelText(/your email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    
    // Fill in the form with test data
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'securepassword' } });
    
    // Submit the form
    fireEvent.click(screen.getByRole('button', { name: /sign up/i }));
    
    // Check that the correct payload was used
    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/signup',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            email: 'test@example.com',
            password: 'securepassword'
          })
        })
      );
    });
  });
});