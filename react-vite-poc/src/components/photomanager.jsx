import { useState, useEffect } from 'react';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Croppie from 'croppie';
import 'croppie/croppie.css';

function PhotoManager({ photos, onPhotoUpdate, onProfilePhotoSelect, uploadUrls, uploadProfileUrls, setUploadUrls }) {
  // Find the profile URLs from uploadUrls
  console.log(uploadProfileUrls)
  const profileUploadUrl = uploadProfileUrls.find(urlObj => urlObj.s3key && urlObj.s3key.includes('/profile'));
  const [userPhotos, setUserPhotos] = useState(photos || []);
  const [selectedFile, setSelectedFile] = useState(null);
  const [caption, setCaption] = useState('');
  const [uploadError, setUploadError] = useState(null);
  const [usedUrls, setUsedUrls] = useState(0);
  const [croppieInstance, setCroppieInstance] = useState(null);
  const [croppingPhoto, setCroppingPhoto] = useState(null);
  const [showCroppie, setShowCroppie] = useState(false);
  const [replacePhoto, setReplacePhoto] = useState(null);

  useEffect(() => {
    // Include all photos, including profile photos, in the PhotoManager
    setUserPhotos(photos || []);
  }, [photos]);

  // Initialize Croppie when needed
  useEffect(() => {
    if (showCroppie && croppingPhoto) {
      const element = document.getElementById('croppie-container');
      if (element) {
        const croppie = new Croppie(element, {
          viewport: { width: 200, height: 200, type: 'square' },
          boundary: { width: 300, height: 300 },
          enableOrientation: true
        });
        
        croppie.bind({
          url: croppingPhoto.url
        });
        
        setCroppieInstance(croppie);
      }
    }
    
    return () => {
      if (croppieInstance) {
        croppieInstance.destroy();
      }
    };
  }, [showCroppie, croppingPhoto]);

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
    setUploadError(null);
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setUploadError('Please select a file to upload.');
      return;
    }

    try {
      let presignedUrl, s3Key;

      if (replacePhoto) {
        // Use the upload URL of the photo being replaced
        presignedUrl = replacePhoto.upload;
        s3Key = replacePhoto.s3key;
      } else {
        // Use the next available upload URL
        if (usedUrls >= uploadUrls.length) {
          setUploadError('You have reached your upload limit.');
          return;
        }
        presignedUrl = uploadUrls[usedUrls].url;
        s3Key = uploadUrls[usedUrls].s3key;
      }

      const response = await fetch(presignedUrl, {
        method: 'PUT',
        headers: {
          'Content-Type': selectedFile.type,
        },
        body: selectedFile,
      });

      if (response.ok) {
        if (replacePhoto) {
          // Replace the existing photo view
          const updatedPhotos = userPhotos.map(photo =>
            photo.s3key === replacePhoto.s3key
              ? { ...photo, url: URL.createObjectURL(selectedFile) } // Update the photo view
              : photo
          );
          setUserPhotos(updatedPhotos);
          onPhotoUpdate(updatedPhotos);
          setReplacePhoto(null); // Clear replacePhoto state
        } else {
          // Add the new photo to the photo collection
          const newPhoto = {
            s3key: s3Key,
            caption,
            url: URL.createObjectURL(selectedFile), // Temporary URL for display
            originalName: selectedFile.name, // Store original filename for reference
          };
          const updatedPhotos = [...userPhotos, newPhoto];
          setUserPhotos(updatedPhotos);
          onPhotoUpdate(updatedPhotos);
          setUsedUrls((prev) => prev + 1);
        }

        setSelectedFile(null);
        setCaption('');
        setUploadError(null);
      } else {
        setUploadError('Failed to upload the image. Please try again.');
      }
    } catch (error) {
      console.error('Error uploading image:', error);
      setUploadError('An error occurred during the upload.');
    }
  };

  const handleDeletePhoto = async (photoToDelete) => {
    try {
      // Make a DELETE request to the presigned Delete URL
      const response = await fetch(photoToDelete.delete, {
        method: 'DELETE',
      });

      if (response.ok) {
        // Remove the photo from the local state if the deletion was successful
        const updatedPhotos = userPhotos.filter(photo => photo.s3key !== photoToDelete.s3key);
        setUserPhotos(updatedPhotos);
        onPhotoUpdate(updatedPhotos);

        // Add the upload URL back to uploadProfileUrls
        setUploadUrls(prev => [
        ...prev,
        { url: photoToDelete.upload, s3key: photoToDelete.s3key }
      ]);
      } else {
        console.error('Failed to delete the photo from S3:', response.statusText);
        setUploadError('Failed to delete the photo. Please try again.');
      }
    } catch (error) {
      console.error('Error deleting photo:', error);
      setUploadError('An error occurred while deleting the photo.');
    }
  };

  const handleReplacePhoto = (photoToReplace) => {
    setReplacePhoto(photoToReplace); // Set the photo to be replaced
    setSelectedFile(null); // Clear any previously selected file
  };

  const handleMakeProfilePhoto = (photo) => {
    setCroppingPhoto({
      ...photo,
      file: selectedFile && selectedFile.name === photo.originalName ? selectedFile : null
    });
    setShowCroppie(true);
  };

  const handleCropComplete = async () => {
    if (croppieInstance) {
      try {
        // Get the cropped image as a blob
        const blob = await croppieInstance.result({
          type: 'blob',
          size: 'viewport',
          format: 'jpeg',
          quality: 1.0
        });

        // Also get base64 for preview
        const base64 = await croppieInstance.result({
          type: 'base64',
          size: 'viewport',
          format: 'jpeg',
          quality: 1.0
        });

        let profileS3Key = null;

        // Only upload raw profile after user clicks "Make Profile Photo"
        if (croppingPhoto) {
          let rawBlob;
          if (croppingPhoto.file) {
            rawBlob = croppingPhoto.file;
          } else if (croppingPhoto.url) {
            const response = await fetch(croppingPhoto.url);
            rawBlob = await response.blob();
          }
          // I left this here but handling the 'raw' profile photo was kind of much; 
          // I think most people will either a) add it to their profile or b) delete it. Keeping a copy seemed
          // like a good idea but jesus the implementation went sideways quick. Passing the photo a little more clearly will,
          // I think work better.
          // if (rawBlob) {
          //   await fetch(rawProfileUploadUrl.url, {
          //     method: 'PUT',
          //     headers: {
          //       'Content-Type': rawBlob.type || 'image/jpeg',
          //     },
          //     body: rawBlob,
          //   });
          // }
        }

        // Upload cropped profile photo
        if (profileUploadUrl) {
          const response = await fetch(profileUploadUrl.url, {
            method: 'PUT',
            headers: {
              'Content-Type': 'image/jpeg',
            },
            body: blob,
          });
          if (response.ok) {
            profileS3Key = profileUploadUrl.s3key;
          }
        }

        if (profileS3Key) {
          onProfilePhotoSelect(profileS3Key);
        } else {
          onProfilePhotoSelect(base64);
        }

        setShowCroppie(false);
        setCroppingPhoto(null);
      } catch (error) {
        console.error('Error cropping image:', error);
      }
    }
  };

  const handleCancelCrop = () => {
    setShowCroppie(false);
    setCroppingPhoto(null);
  };

  const handleCancelReplace = () => {
    setReplacePhoto(null); // Clear replacePhoto state
    setSelectedFile(null); // Clear the selected file
  };

  const handleCaptionChange = (photoIndex, newCaption) => {
    const updatedPhotos = [...userPhotos];
    updatedPhotos[photoIndex].caption = newCaption;
    setUserPhotos(updatedPhotos);
    onPhotoUpdate(updatedPhotos);
  };

  return (
    <div className="photo-manager">
      <h3>Photo Manager</h3>
      
      {showCroppie ? (
        <div className="croppie-wrapper my-3">
          <div id="croppie-container"></div>
          <div className="d-flex justify-content-center mt-3">
            <Button variant="secondary" className="me-2" onClick={handleCancelCrop}>
              Cancel
            </Button>
            <Button variant="primary" onClick={handleCropComplete}>
              Set as Profile Photo
            </Button>
          </div>
        </div>
      ) : (
        <>
          {/* Upload new photo section */}
          <div className="upload-section mb-4">
            <Form.Group controlId="photoUpload">
              <Form.Label>Upload New Photo</Form.Label>
              <Form.Control
                type="file"
                accept="image/*"
                onChange={handleFileSelect}
                disabled={usedUrls >= uploadUrls.length}
              />
            </Form.Group>
            
            {selectedFile && (
              <div className="my-2">
                <img 
                  src={URL.createObjectURL(selectedFile)} 
                  alt="Preview" 
                  style={{ maxHeight: '150px', maxWidth: '100%' }} 
                  className="mb-2"
                />
                <Form.Group controlId="photoCaption">
                  <Form.Label>Caption</Form.Label>
                  <Form.Control
                    type="text"
                    placeholder="Add a caption"
                    value={caption}
                    onChange={(e) => setCaption(e.target.value)}
                  />
                </Form.Group>
              </div>
            )}
            
            {uploadError && <p className="text-danger">{uploadError}</p>}
            
            <Button
              variant="primary"
              onClick={handleUpload}
              disabled={!selectedFile || usedUrls >= uploadUrls.length}
              className="mt-2"
            >
              {usedUrls >= uploadUrls.length ? 'Upload Limit Reached' : 
              (replacePhoto ? 'Replace Photo' : 'Upload Photo')}
            </Button>
            {replacePhoto && (
            <Button variant="secondary" onClick={handleCancelReplace} className="ms-2">
              Cancel
            </Button>
            )}
          </div>

          {/* Photo gallery */}
          <h4>Your Photos</h4>
            <Row xs={1} md={2} lg={3} className="g-4">
              {userPhotos.length === 0 ? (
                <Col>
                  <Card className="border-primary">
                    <div className="position-absolute top-0 start-0 bg-primary text-white p-1">
                      Current Profile
                    </div>
                    <Card.Img 
                      variant="top" 
                      src="./profile.svg"
                      style={{ height: '200px', objectFit: 'cover' }}
                    />
                    <Card.Body>
                      <span className="text-muted">No photos uploaded yet</span>
                    </Card.Body>
                  </Card>
                </Col>
              ) : (
                userPhotos.map((photo, index) => {
                  const isProfilePhoto = photo.s3key && photo.s3key.endsWith('/profile');
                  const isRawProfilePhoto = photo.s3key && photo.s3key.endsWith('/rawprofile');
                  return (
                    <Col key={photo.s3key || index}>
                      <Card className={isProfilePhoto ? 'border-primary' : ''}>
                        {isProfilePhoto && (
                          <div className="position-absolute top-0 start-0 bg-primary text-white p-1">
                            Current Profile
                          </div>
                        )}
                        {isRawProfilePhoto && (
                          <div className="position-absolute top-0 start-0 bg-secondary text-white p-1">
                            Raw Profile
                          </div>
                        )}
                        <Card.Img 
                          variant="top" 
                          src={photo.url || "./profile.svg"} 
                          style={{ height: '200px', objectFit: 'cover' }}
                        />
                        <Card.Body>
                          <Form.Group>
                            <Form.Control
                              type="text"
                              placeholder={isProfilePhoto ? 'Your current profile photo' : "Add a caption"}
                              value={photo.caption || ''}
                              onChange={(e) => handleCaptionChange(index, e.target.value)}
                              disabled={isProfilePhoto}
                            />
                          </Form.Group>
                          <div className="d-flex justify-content-between mt-2">
                            {!isProfilePhoto && (
                              <>
                              <Button 
                                variant="outline-primary" 
                                size="sm"
                                onClick={() => handleMakeProfilePhoto(photo)}
                              >
                                Make Profile Photo
                              </Button>

                              <Button 
                                variant="outline-warning" 
                                size="sm"
                                onClick={() => handleReplacePhoto(photo)}
                              >
                                Replace
                              </Button>
                              </>
                            )}
                            {!isProfilePhoto && (
                              <Button 
                                variant="outline-danger" 
                                size="sm"
                                onClick={() => handleDeletePhoto(photo)}
                                disabled={isRawProfilePhoto}
                              >
                                Remove
                              </Button>
                            )}
                            {isProfilePhoto && (
                              <span className="text-muted">Current profile photo</span>
                            )}
                          </div>
                        </Card.Body>
                      </Card>
                    </Col>
                  );
                })
              )}
            </Row>
        </>
      )}
    </div>
  );
}

export default PhotoManager;