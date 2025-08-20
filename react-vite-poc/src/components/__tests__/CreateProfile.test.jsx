import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { CreateProfile } from '../createprofile';

// Mock fetch for API calls
global.fetch = jest.fn();

// Mock alert
global.alert = jest.fn();

// Mock console methods to reduce noise
global.console.log = jest.fn();
global.console.error = jest.fn();

// Mock child components
jest.mock('../detailsselections', () => ({
  __esModule: true,
  default: ({ onChange, selectedValues }) => (
    <div data-testid="details-selections">
      <button 
        onClick={() => onChange('detail1', '1')} 
        data-testid="mock-detail-button"
      >
        Set Detail
      </button>
    </div>
  ),
}));

jest.mock('../editor', () => ({
  useQuillLoader: () => true,
  QuillEditor: ({ onContentChange, initialContent }) => (
    <div data-testid="quill-editor">
      <textarea 
        data-testid="mock-editor"
        onChange={(e) => onContentChange(e.target.value)}
        defaultValue={initialContent ? JSON.stringify(initialContent) : ''}
      />
    </div>
  ),
}));

jest.mock('./photomanager', () => ({
  __esModule: true,
  default: ({ photos, onPhotoUpdate, onProfilePhotoSelect }) => (
    <div data-testid="photo-manager">
      <button 
        onClick={() => onPhotoUpdate([...photos, { url: 'new-photo.jpg', s3key: 's3key-123' }])} 
        data-testid="add-photo-button"
      >
        Add Photo
      </button>
      <button 
        onClick={() => onProfilePhotoSelect({ S3Key: 'profile-s3key' })} 
        data-testid="set-profile-button"
      >
        Set Profile Photo
      </button>
    </div>
  ),
}));

jest.mock('./carousel', () => ({
  __esModule: true,
  default: () => <div data-testid="carousel">Carousel Mock</div>,
}));

jest.mock('./matchlist', () => ({
  __esModule: true,
  default: ({ people }) => (
    <div data-testid="match-list">
      {people.map(person => (
        <div key={person.id} data-testid="preview-person">
          {person.name}, {person.age}
        </div>
      ))}
    </div>
  ),
}));

// Mock geolocation
const mockGeolocation = {
  getCurrentPosition: jest.fn().mockImplementation((success) => {
    success({ coords: { latitude: 37.7749, longitude: -122.4194 } });
  }),
};
global.navigator.geolocation = mockGeolocation;

// Mock window.confirm
global.confirm = jest.fn().mockImplementation(() => true);

describe('CreateProfile Component', () => {
  // Mock props
  const mockSetLoggedInUser = jest.fn();
  const mockSetPendingID = jest.fn();
  const mockSetUploadUrls = jest.fn();
  const mockSetUploadProfileUrls = jest.fn();
  
  const defaultProps = {
    setLoggedInUser: mockSetLoggedInUser,
    pendingID: 123,
    setPendingID: mockSetPendingID,
    loggedInUser: null,
    uploadUrls: ['url1', 'url2'],
    uploadProfileUrls: ['profile-url1'],
    setUploadUrls: mockSetUploadUrls,
    setUploadProfileUrls: mockSetUploadProfileUrls,
    jwt: 'mock-jwt-token'
  };
  
  beforeEach(() => {
    // Reset mocks before each test
    jest.clearAllMocks();
    global.fetch.mockReset();
    global.alert.mockReset();
    global.confirm.mockClear();
  });
  
  test('renders the form fields correctly', () => {
    render(<CreateProfile {...defaultProps} />);
    
    // Check that form elements are rendered
    expect(screen.getByLabelText(/name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/age/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/title/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/latitude/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/longitude/i)).toBeInTheDocument();
    expect(screen.getByText(/use my location/i)).toBeInTheDocument();
    expect(screen.getByTestId('quill-editor')).toBeInTheDocument();
    expect(screen.getByTestId('photo-manager')).toBeInTheDocument();
    expect(screen.getByTestId('details-selections')).toBeInTheDocument();
    expect(screen.getByText(/preview profile/i)).toBeInTheDocument();
    expect(screen.getByText(/finish your profile/i)).toBeInTheDocument();
  });
  
  test('handles input changes correctly', () => {
    render(<CreateProfile {...defaultProps} />);
    
    // Get form inputs
    const nameInput = screen.getByLabelText(/name/i);
    const ageInput = screen.getByLabelText(/age/i);
    const titleInput = screen.getByLabelText(/title/i);
    
    // Simulate typing in inputs
    fireEvent.change(nameInput, { target: { value: 'Test User' } });
    fireEvent.change(ageInput, { target: { value: '30' } });
    fireEvent.change(titleInput, { target: { value: 'Test Title' } });
    
    // Check input values were updated
    expect(nameInput.value).toBe('Test User');
    expect(ageInput.value).toBe('30');
    expect(titleInput.value).toBe('Test Title');
  });
  
  test('handles geolocation correctly', async () => {
    render(<CreateProfile {...defaultProps} />);
    
    // Click the location button
    fireEvent.click(screen.getByText(/use my location/i));
    
    // Verify confirm was called
    expect(confirm).toHaveBeenCalledWith(
      expect.stringContaining('access your location')
    );
    
    // Check that lat/long fields were updated
    await waitFor(() => {
      expect(screen.getByLabelText(/latitude/i).value).toBe('37.775');
      expect(screen.getByLabelText(/longitude/i).value).toBe('-122.419');
    });
  });
  
  test('handles details selection correctly', () => {
    render(<CreateProfile {...defaultProps} />);
    
    // Click the mock detail button
    fireEvent.click(screen.getByTestId('mock-detail-button'));
    
    // This should trigger the handleDetailsChange function
    // We'll verify this when submitting the form
  });
  
  test('handles photo selection correctly', () => {
    render(<CreateProfile {...defaultProps} />);
    
    // Click the add photo button
    fireEvent.click(screen.getByTestId('add-photo-button'));
    
    // Click the profile photo button
    fireEvent.click(screen.getByTestId('set-profile-button'));
    
    // We'll verify the photos were added when submitting the form
  });
  
  test('toggles preview correctly', () => {
    render(<CreateProfile {...defaultProps} />);
    
    // Preview should initially be hidden
    expect(screen.queryByTestId('match-list')).not.toBeInTheDocument();
    
    // Click preview button
    fireEvent.click(screen.getByText(/preview profile/i));
    
    // Preview should now be visible
    expect(screen.getByTestId('match-list')).toBeInTheDocument();
    expect(screen.getByTestId('preview-person')).toBeInTheDocument();
    
    // Click hide button
    fireEvent.click(screen.getByText(/hide preview/i));
    
    // Preview should be hidden again
    expect(screen.queryByTestId('match-list')).not.toBeInTheDocument();
  });
  
  test('handles form submission successfully', async () => {
    const mockResponse = {
      person: { id: 123, name: 'Test User' },
      upload_urls: ['new-url1', 'new-url2'],
      profile_upload_urls: ['new-profile-url']
    };
    
    // Mock successful response
    fetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(mockResponse)
    });
    
    render(<CreateProfile {...defaultProps} />);
    
    // Fill in required fields
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'Test User' } });
    fireEvent.change(screen.getByLabelText(/age/i), { target: { value: '30' } });
    fireEvent.change(screen.getByLabelText(/title/i), { target: { value: 'Test Title' } });
    fireEvent.change(screen.getByLabelText(/latitude/i), { target: { value: '37.7749' } });
    fireEvent.change(screen.getByLabelText(/longitude/i), { target: { value: '-122.4194' } });
    
    // Mock editor content
    fireEvent.change(screen.getByTestId('mock-editor'), { 
      target: { value: JSON.stringify({ ops: [{ insert: 'Test description' }] }) } 
    });
    
    // Add a photo
    fireEvent.click(screen.getByTestId('add-photo-button'));
    
    // Set profile photo
    fireEvent.click(screen.getByTestId('set-profile-button'));
    
    // Add a detail
    fireEvent.click(screen.getByTestId('mock-detail-button'));
    
    // Submit the form
    fireEvent.click(screen.getByText(/finish your profile/i));
    
    // Check that fetch was called with correct data
    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/people', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'mock-jwt-token'
        },
        body: expect.any(String)
      });
      
      // Parse the request body to check its contents
      const requestBody = JSON.parse(fetch.mock.calls[0][1].body);
      expect(requestBody).toEqual(expect.objectContaining({
        id: 123,
        name: 'Test User',
        age: 30,
        motto: 'Test Title',
        lat: 37.7749,
        long: -122.4194,
        profile: expect.objectContaining({ S3Key: expect.any(String) }),
        description: expect.any(String),
        details: expect.objectContaining({ detail1: 1 }),
        photos: expect.arrayContaining([
          expect.objectContaining({ 
            S3Key: expect.any(String),
            caption: expect.any(String)
          })
        ])
      }));
      
      // Check that state setter functions were called with correct values
      expect(mockSetLoggedInUser).toHaveBeenCalledWith(mockResponse.person);
      expect(mockSetUploadUrls).toHaveBeenCalledWith(mockResponse.upload_urls);
      expect(mockSetUploadProfileUrls).toHaveBeenCalledWith(mockResponse.profile_upload_urls);
      expect(mockSetPendingID).toHaveBeenCalledWith(null);
      
      // Check success alert
      expect(alert).toHaveBeenCalledWith('Completed your profile! Nice');
    });
  });
  
  test('handles form submission error', async () => {
    // Mock failed response
    fetch.mockResolvedValueOnce({
      ok: false
    });
    
    render(<CreateProfile {...defaultProps} />);
    
    // Fill in minimal required fields
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'Test User' } });
    
    // Submit the form
    fireEvent.click(screen.getByText(/finish your profile/i));
    
    // Check that error alert was shown
    await waitFor(() => {
      expect(alert).toHaveBeenCalledWith('Ruh roh Shaggy');
    });
    
    // Check that state setter functions were not called
    expect(mockSetLoggedInUser).not.toHaveBeenCalled();
  });
  
  test('handles network error during submission', async () => {
    // Mock network error
    fetch.mockRejectedValueOnce(new Error('Network error'));
    
    render(<CreateProfile {...defaultProps} />);
    
    // Fill in minimal required fields
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'Test User' } });
    
    // Submit the form
    fireEvent.click(screen.getByText(/finish your profile/i));
    
    // Check that error handling occurred
    await waitFor(() => {
      expect(console.error).toHaveBeenCalledWith('Error during profile creation:', expect.any(Error));
      expect(alert).toHaveBeenCalledWith('An error occurred during profile creation.');
    });
  });
  
  test('pre-fills form data when loggedInUser exists', () => {
    const userWithData = {
      id: 456,
      name: 'Existing User',
      age: 35,
      motto: 'Existing Title',
      lat: 40.7128,
      long: -74.0060,
      description: JSON.stringify({ ops: [{ insert: 'Existing description' }] }),
      details: { detail1: 2, detail2: 3 }
    };
    
    render(
      <CreateProfile 
        {...defaultProps} 
        loggedInUser={userWithData} 
      />
    );
    
    // Check that form fields are pre-filled
    expect(screen.getByLabelText(/name/i).value).toBe('Existing User');
    expect(screen.getByLabelText(/age/i).value).toBe('35');
    expect(screen.getByLabelText(/title/i).value).toBe('Existing Title');
    expect(screen.getByLabelText(/latitude/i).value).toBe('40.7128');
    expect(screen.getByLabelText(/longitude/i).value).toBe('-74.006');
    
    // Check that the submit button text is different for existing users
    expect(screen.getByText(/update your profile/i)).toBeInTheDocument();
  });
});