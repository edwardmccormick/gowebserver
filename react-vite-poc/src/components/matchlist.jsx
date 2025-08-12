import { useState,useEffect } from 'react';
import Accordion from 'react-bootstrap/Accordion';
import ControlledCarousel from './carousel'
import { getPreciseDistance } from 'geolib';
import Button from 'react-bootstrap/Button';
import { useQuillLoader, QuillEditor } from './editor';
import MatchMap from './MatchMap';
import DetailFlags from './DetailFlags';

function MatchList({
    people, 
    loading,
    User,
    refreshMatches
  }) {

  const [submittedLikes, setSubmittedLikes] = useState({}); // Track submitted likes
  const [showMap, setShowMap] = useState(false); // Toggle between list and map view
  const [isFlipping, setIsFlipping] = useState(false); // Animation state
  const isQuillLoaded = useQuillLoader(); // Load Quill at the top level
  // let people = peopleObject.people;
  // Add distance to each person and sort by distance
  if (!loading && User) {
    people = people.map((person) => ({
      ...person,
      distance: Math.round(
        getPreciseDistance(
          { latitude: User.lat, longitude: User.long },
          { latitude: person.lat, longitude: person.long }
        ) / 1609.34 * 10
      ) / 10, // Convert meters to miles and round to 1 decimal place
    }));

    // Sort people by distance (smallest to largest)
  people.sort((a, b) => a.distance - b.distance);
  }

  useEffect(() => {
    console.log("People updated:", people);
  }, [people]);
  
  const handleSubmit = async (User, person) => {
    const payload = {
      offered: User.id,
      OfferedProfile: User,
      accepted: person.id,
      AcceptedProfile: person
    };
    try {
      const response = await fetch('http://localhost:8080/matches', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        alert('Match added successfully!');
        setSubmittedLikes((prev) => ({ ...prev, [person.id]: true })); // Mark as submitted
        refreshMatches(); // Refresh matches in ChatSelect
      } else {
        alert('Failed to add match.');
      }
    } catch (error) {
      console.error('Error submitting form:', error);
      alert('An error occurred while submitting the form.');
    }
  };

  // Refactored renderDeltaAsHtml function
  const renderDeltaAsHtml = (deltaString) => {
    if (!deltaString) {
      return { __html: '<p>Nothing to show yet.</p>' };
    }

    if (!isQuillLoaded) {
      return { __html: '<p>Loading content...</p>' }; // Handle case where Quill is not yet loaded
    }

    try {
      // Parse the string into a Delta JSON object
      const delta = JSON.parse(deltaString);

      // Create a temporary, non-rendered Quill instance to convert Delta to HTML
      const tempContainer = document.createElement('div');
      const tempQuill = new window.Quill(tempContainer);
      tempQuill.setContents(delta); // Load the Delta
      return { __html: tempContainer.querySelector('.ql-editor').innerHTML };
    } catch (error) {
      console.error('Error parsing Delta string:', error);
      return { __html: '<p>Invalid content format.</p>' };
    }
  };

  const handleViewToggle = () => {
    setIsFlipping(true);
    setTimeout(() => {
      setShowMap(!showMap);
      setIsFlipping(false);
    }, 300);
  };
    return (
      loading || User==undefined ? (
        <div className="d-flex align-items-center">
          <strong role="status">Loading...</strong>
          <div className="spinner-border ms-auto" aria-hidden="true"></div>
        </div>
      ) : (
      <>
      <div className="d-flex justify-content-between align-items-center mb-3">
        <h2>And these {people ? `${people.length} ` : null }cool folks want to be liked by you:</h2>
        <Button 
          variant="outline-primary" 
          onClick={handleViewToggle}
          disabled={isFlipping}
        >
          {showMap ? '📋 List View' : '🗺️ Map View'}
        </Button>
      </div>
      <div 
        className={`view-container ${isFlipping ? 'flipping' : ''}`}
        style={{
          transform: isFlipping ? 'rotateY(90deg)' : 'rotateY(0deg)',
          transition: 'transform 0.3s ease-in-out',
          transformStyle: 'preserve-3d'
        }}
      >
        {showMap ? (
          <MatchMap people={people} User={User} />
        ) : (
      <Accordion className='w-100'>
        {people.map((person) => (
        <Accordion.Item className='w-100' eventKey={person.id} key={person.id}>
          <Accordion.Header className='w-100' key={`${person.id}100`}>
            <img
              src={person.profile.url !== null && person.profile.url !== "" ? person.profile.url : '/profile.svg'}
              style={{ borderRadius: '50%' }}
              className="m-1 p-1"
              height="50"
              width="50"
              alt={`${person.name}'s profile`}
            />
            <strong>{person.name}</strong> <span>{person.motto} - Distance: { person.distance } miles</span>
            {person.details && <DetailFlags details={person.details} limit={3} topValues={true} />}
          </Accordion.Header>
          <Accordion.Body key={`${person.id}10000`}>
            <div className='text-start'>
              <img src={person.profile.url !== null && person.profile.url !== "" ? person.profile.url : '/profile.svg'} className='m-1 p-1 rounded float-start' height={'250'} width={'auto'} />
              {person.details && (
                <div className="mb-3">
                  <h5>Personality Indicators:</h5>
                  <DetailFlags details={person.details} />
                </div>
              )}
                { person?.description ? 
                (
                  <div
                    dangerouslySetInnerHTML={renderDeltaAsHtml(
                      person.description
                    )}
                  style={{ minHeight: '275px' }}
                  />
                )
                : (<div> 
                  <p>(Bender) There we were in the park when suddenly some old lady says I stole her purse. I chucked the professor at her but she kept coming. So I had to hit her with this purse I found.</p>
                    <p>(Bender) Boy, who knew a cooler could also make a handy wang coffin?!</p>
                    <p>Leela Futurama Quotes: (Amy After Bender destroys Fry's tent) Bender, wasn't that Fry's Tent? (Bender Responds Scoffing) Bender, Mominey mum meh. (Leela) Bender Raises a good point. Where is Fry?</p>
                    <p>Farnsworth Futurama Quotes: (Farnsworth to group) As new employees, I'd like your opinion on our commercial. I've paid to have it aired during the Super Bowl (Fry) Wow. (Farnsworth) Not on the same channel of course.</p>
                    <p>(Zap) The spirit is willing, but the flesh is spongy and bruised.</p>
                    <p>(President Truman) Bush-wah! Now, what's your mission? Are you here to make some sort of alien-human hybrid? (Zoidberg) Are you coming on to me? President Truman: Hot Crackers! I take exception to that! (Zoidberg) (leering) I'm not hearing a no...</p>
                    <p>(Zoidberg) Did you see me escaping? I was all like, "WOO WOO WOO WOO!"</p>
                    <p>(Lucy Liu-bot) You're cute!  (Fry) No, you are!  (Lucy Liu-bot) No, you!  (Fry) No, you! (Lucy Liu-bot) No, you! (Fry) No, you!  (Professor Hubert Farnsworth) Oh dear, she's stuck in an infinite loop and he's an idiot!</p>
                    <p>(Zap) Prepare to continue the epic struggle between good and neutral.</p>
                    <p>(Fry)I did do the nasty in the pasty!</p>
                  </div>
                )}

              <div height='450' width='450' className='text-center w-fill'>
                <ControlledCarousel
                  key={`${person.id}1000000`}
                  photos={person.photos}
                  id={person.id}
                />
              </div>
              <div className="mx-auto w-50 d-flex flex-row justify-content-around align-items-center">
                {submittedLikes[person.id] ? (
                  <p className="text-success">Like submitted!</p>
                ) : (
                  <>
                    <Button
                      variant="primary"
                      className="m-2"
                      onClick={() => handleSubmit(User, person)}
                    >
                      Oh yeah, that's what I like! Match!
                    </Button>
                    <button className="btn btn-danger m-2">
                      Show me less like this person
                    </button>
                  </>
                )}
              </div>
            </div>
            </Accordion.Body>
        </Accordion.Item>
        ))}
      </Accordion>
        )}
      </div>
    </>
      )
    
    )
  
}

export default MatchList