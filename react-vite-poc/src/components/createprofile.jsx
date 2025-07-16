import { useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import DetailsSelections from './detailsselections';
import { useQuillLoader, QuillEditor } from './editor';
import details from '../../../details.json';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Tooltip from 'react-bootstrap/Tooltip';
import ControlledCarousel from './carousel';

export function CreateProfile({setLoggedInUser, pendingID, setPendingID, loggedInUser, uploadUrls}) {
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
    const payload = {
      id: loggedInUser?.id || pendingID,
      age: parseInt(formData.age),
      name: formData.name,
      motto: formData.motto,
      lat: parseFloat(formData.latitude),
      long: parseFloat(formData.longitude),
      profile: formData.profile,
      description: JSON.stringify(editorDelta), // Use the Delta content from the editor
      details: formData.details,
      photos: formData.photos.map(photo => ({ 
        S3Key: photo.s3key,
        caption: photo.caption || '',
      })),
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

  const [usedUrls, setUsedUrls] = useState(0); // Track the number of used presigned URLs
  const [uploadError, setUploadError] = useState(null); // Track upload errors
  const [selectedFile, setSelectedFile] = useState(null); // Track the selected file
  const [caption, setCaption] = useState(''); // Track the caption for the image


  const handleFileSelect = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    // Validate file extension
    const validExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp'];
    const fileExtension = file.name.split('.').pop().toLowerCase();
    if (!validExtensions.includes(fileExtension)) {
      setUploadError('Invalid file type. Please upload an image file.');
      return;
    }

    setSelectedFile(file);
    setUploadError(null); // Clear any previous errors
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setUploadError('Please select a file to upload.');
      return;
    }

    // Check if there are remaining presigned URLs
    if (usedUrls >= uploadUrls.length) {
      setUploadError('You have reached your upload limit.');
      return;
    }

    try {
      const presignedUrl = uploadUrls[usedUrls].url; // Get the presigned URL for the next upload
      const s3Key = uploadUrls[usedUrls].s3key; // Use the S3 key that was generated with the presigned urls

      const response = await fetch(presignedUrl, {
        method: 'PUT',
        headers: {
          'Content-Type': selectedFile.type,
        },
        body: selectedFile,
      });

      if (response.ok) {
        setFormData((prev) => ({
          ...prev,
          photos: [
            ...prev.photos,
            { s3key: s3Key, caption }, // Persist the S3 key and caption
          ],
        }));
        setUsedUrls((prev) => prev + 1); // Increment the used URLs count
        setSelectedFile(null); // Clear the selected file
        setCaption(''); // Clear the caption
        setUploadError(null); // Clear any previous errors
        alert('Image uploaded successfully!');
      } else {
        setUploadError('Failed to upload the image. Please try again.');
      }
    } catch (error) {
      console.error('Error uploading image:', error);
      setUploadError('An error occurred during the upload.');
    }
  };

  const renderTooltip = (props) => (
    <Tooltip {...props}>
      You've reached your upload limit for images.
    </Tooltip>
  );

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

        <div className="col-3 m-1 p-1">
          <Form.Label htmlFor="profile">Profile Picture URL</Form.Label>
          <Form.Control
            id="profile"
            type="text"
            placeholder="https://wwww.example.com/profile.jpg"
            value={formData.profile}
            onChange={handleChange}
          />
        </div>

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
     {/* Upload Images Section */}
      <div className="text-center my-3">
        <Form.Group controlId="imageUpload">
          <Form.Label>Upload Images</Form.Label>
          <Form.Control
            type="file"
            accept="image/*"
            onChange={handleFileSelect}
            disabled={usedUrls >= uploadUrls.length} // Disable if all URLs are used
          />
        </Form.Group>
        {selectedFile && (
          <div className="my-3">
            <ControlledCarousel
              photos={[{ url: URL.createObjectURL(selectedFile), caption }]} // Preview the selected image
              id={loggedInUser?.id || pendingID}
            />
            <Form.Group controlId="imageCaption">
              <Form.Label>Image Caption</Form.Label>
              <Form.Control
                type="text"
                placeholder="Enter a caption for the image"
                value={caption}
                onChange={(e) => setCaption(e.target.value)}
              />
            </Form.Group>
          </div>
        )}
        {uploadError && <p className="text-danger">{uploadError}</p>}
        <OverlayTrigger
          placement="top"
          overlay={renderTooltip}
          show={usedUrls >= uploadUrls.length} // Show tooltip when limit is reached
        >
          <Button
            variant="primary"
            onClick={handleUpload}
            disabled={usedUrls >= uploadUrls.length} // Disable button if limit is reached
          >
            {usedUrls >= uploadUrls.length ? 'Upload Limit Reached' : 'Upload Image'}
          </Button>
        </OverlayTrigger>
      </div>

      {/* Submit Button */}
      <Button variant="primary" onClick={handleSubmit}>
        {loggedInUser && loggedInUser.name ? 'Update your profile' : 'Finish your profile'}
      </Button>
    </>
  );
}

export default CreateProfile;