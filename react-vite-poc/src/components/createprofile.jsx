import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import DetailsSelections from './detailsselections';
import { useQuillLoader, QuillEditor } from './editor';
import details from '../../../details.json';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Tooltip from 'react-bootstrap/Tooltip';
import ControlledCarousel from './carousel';
import MatchList from './matchlist';
import PhotoManager from './photomanager';

export function CreateProfile({
  setLoggedInUser, 
  pendingID, 
  setPendingID, 
  loggedInUser, 
  uploadUrls,
  uploadProfileUrls,
}) {
  console.log(uploadProfileUrls)
  const [formData, setFormData] = useState(() => {
    // If loggedInUser exists and has a name, use its values to populate the form
    if (loggedInUser && loggedInUser.name) {
      return {
        name: loggedInUser.name || '',
        age: loggedInUser.age || '',
        description: loggedInUser.description || '',
        motto: loggedInUser.motto || '',
        latitude: loggedInUser.lat || '',
        longitude: loggedInUser.long || '',
        profile: loggedInUser.profile || '',
        details: loggedInUser.details || Object.keys(details).reduce((acc, key) => ({ ...acc, [key]: null }), {}),
        photos: loggedInUser.photos || [],
      };
    }
    // Otherwise, use empty values
    return {
      name: '',
      age: '',
      description: '',
      motto: '',
      latitude: '',
      longitude: '',
      profile: '',
      details: {},
      photos: [],
    };
  });

  const handleChange = (e) => {
    const { id, value } = e.target;
    setFormData((prev) => ({ ...prev, [id]: value }));
  };

  const handleContentChange = (content) => {
    setEditorDelta(content);
  };

  const handleDetailsChange = (key, value) => {
    setFormData((prev) => ({
      ...prev,
      details: { ...prev.details, [key]: parseInt(value, 10) },
    }));
  };

  const isQuillLoaded = useQuillLoader(); // Hook to load Quill scripts
  // State to hold the editor's content in Delta format
  const [editorDelta, setEditorDelta] = useState(() => {
    // If loggedInUser exists and has a description, try to parse it as JSON
    if (loggedInUser && loggedInUser.description) {
      try {
        return JSON.parse(loggedInUser.description);
      } catch (e) {
        console.error('Error parsing description:', e);
        return null;
      }
    }
    return null;
  });
  // State to hold the "submitted" content, which we'll then render as HTML

  const handleSubmit = async () => {
    // Prepare photos array, ensuring we don't have duplicates
    const uniquePhotos = [];
    const seenKeys = new Set();
    
    formData.photos.forEach(photo => {
      if (!seenKeys.has(photo.s3key)) {
        seenKeys.add(photo.s3key);
        uniquePhotos.push({
          S3Key: photo.s3key,
          caption: photo.caption || '',
        });
      }
    });
    
    const payload = {
      id: loggedInUser?.id || pendingID,
      age: parseInt(formData.age),
      name: formData.name,
      motto: formData.motto,
      lat: parseFloat(formData.latitude),
      long: parseFloat(formData.longitude),
      profile: { s3key: formData.profile || "" },
      description: JSON.stringify(editorDelta), // Use the Delta content from the editor
      details: formData.details,
      photos: uniquePhotos,
    };

    try {
      const response = await fetch('http://localhost:8080/people', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        const data = await response.json();
        setLoggedInUser(data); // Update the logged-in user in App.jsx
        setPendingID(null);
        console.log('Profile created successfully:', data);
        alert(loggedInUser && loggedInUser.name ? 'Updated your profile! Nice' : 'Completed your profile! Nice');
      } else {
        alert('Ruh roh Shaggy');
      }
      } catch (error) {
        console.error('Error during profile creation:', error);
        alert('An error occurred during profile creation.');
      }
  };

  const [showPreview, setShowPreview] = useState(false); // Toggle profile preview


  const handlePhotoUpdate = (updatedPhotos) => {
    setFormData(prev => ({
      ...prev,
      photos: updatedPhotos
    }));
  };

  const handleProfilePhotoSelect = (profilePhotoValue) => {
    // If the value is an S3 key (string without data:image prefix), store it as is
    // Otherwise, it's a base64 image which we'll store directly
    setFormData(prev => ({
      ...prev,
      profile: profilePhotoValue
    }));
  };

 return (
    <>
      <div className="d-flex justify-content-center align-items-center flex-wrap">


        <div className="col-2 m-1 p-1">
          <Form.Label htmlFor="name">Name</Form.Label>
          <Form.Control
            id="name"
            type="text"
            placeholder="Name"
            value={formData.name}
            onChange={handleChange}
          />
        </div>

        <div className="col-1 m-1 p-1">
          <Form.Label htmlFor="age">Age</Form.Label>
          <Form.Control
            id="age"
            type="text"
            placeholder="0"
            value={formData.age}
            onChange={handleChange}
          />
        </div>

        <div className="col-3 m-1 p-1">
          <Form.Label htmlFor="motto">Title</Form.Label>
          <Form.Control
            id="motto"
            type="text"
            placeholder="Besides your name and picture, this is the first chance you have to turn people off."
            value={formData.motto}
            onChange={handleChange}
          />
        </div>

        {/* <div className="col-3 m-1 p-1">
          <Form.Label htmlFor="profile">Profile Picture URL</Form.Label>
          <Form.Control
            id="profile"
            type="text"
            placeholder="https://wwww.example.com/profile.jpg"
            value={formData.profile}
            onChange={handleChange}
          />
        </div> */}

        <div className="col-1 m-1 p-1">
          <Form.Label htmlFor="latitude">Latitude</Form.Label>
          <Form.Control
            id="latitude"
            type="text"
            placeholder="Latitude"
            value={formData.latitude}
            onChange={handleChange}
          />
        </div>

        <div className="col-1 m-1 p-1">
          <Form.Label htmlFor="longitude">Longitude</Form.Label>
          <Form.Control
            id="longitude"
            type="text"
            placeholder="Longitude"
            value={formData.longitude}
            onChange={handleChange}
          />
        </div>
      </div>
        {/* <div className="mx-auto text-center m-1 p-1 w-75">
          <Form.Label htmlFor="description">Describe yourself</Form.Label>
          <Form.Control as="textarea" rows={5}
            id="description"
            type="textarea"
            placeholder="This is the part where you pretend to have something interesting to say."
            value={formData.description}
            onChange={handleChange}
          />
        </div> */}
      <div className="mx-auto text-center m-1 p-1 w-75" style={{ maxHeight: '300px' }}>
        {isQuillLoaded ? (
          <QuillEditor onContentChange={handleContentChange} initialContent={editorDelta} />
          ) : (
              <div className="d-flex align-items-center justify-content-center bg-secondary-subtle">
                  <p className="text-secondary-emphasis">Hold on, finding some pen and paper to write this down....</p>
              </div>
          )}
      </div>
<br />
     {/* Photo Manager Section */}
      <div className="my-4">
        <PhotoManager
          photos={formData.photos}
          onPhotoUpdate={handlePhotoUpdate}
          onProfilePhotoSelect={handleProfilePhotoSelect}
          uploadUrls={uploadUrls}
          uploadProfileUrls={uploadProfileUrls}
          userId={loggedInUser?.id || pendingID}
        />
      </div>

      <DetailsSelections
        onChange={handleDetailsChange}
        selectedValues={formData.details}
      />

      {/* Preview Toggle Button */}
      <Button 
        variant="secondary" 
        className="mx-2" 
        onClick={() => setShowPreview(!showPreview)}
      >
        {showPreview ? 'Hide Preview' : 'Preview Profile'}
      </Button>

      {/* Profile Preview */}
      {showPreview && (
        <div className="my-4 border rounded p-3">
          <h3 className="text-center mb-3">Profile Preview</h3>
          <MatchList
            people={[{
              id: loggedInUser?.id || pendingID || 'preview',
              name: formData.name || 'Your Name',
              age: formData.age || '0',
              motto: formData.motto || 'Your Title',
              lat: parseFloat(formData.latitude) || 0,
              long: parseFloat(formData.longitude) || 0,
              profile: formData.profile ? 
                formData?.profile 
                : '/profile.svg',
              description: JSON.stringify(editorDelta),
              photos: formData.photos.map(photo => ({
                url: photo.url || `http://localhost:8080/photos/${photo.s3key}`,
                caption: photo.caption || ''
              }))
            }]}
            loading={false}
            User={{
              id: loggedInUser?.id || pendingID || 'preview',
              lat: parseFloat(formData.latitude) || 0,
              long: parseFloat(formData.longitude) || 0
            }}
            refreshMatches={() => {}}
          />
        </div>
      )}

      {/* Submit Button */}
      <Button variant="primary" onClick={handleSubmit}>
        {loggedInUser && loggedInUser.name ? 'Update your profile' : 'Finish your profile'}
      </Button>
    </>
  );
}

export default CreateProfile;