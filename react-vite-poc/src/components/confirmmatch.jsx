import { useState,useEffect } from 'react';
import Accordion from 'react-bootstrap/Accordion';
import ControlledCarousel from './carousel'
import { getPreciseDistance } from 'geolib';
import Button from 'react-bootstrap/Button';

function convertISODateToLocal(dateString) {
  const date = new Date(dateString);


  const localTime = date.toLocaleDateString('en-US', {
    month:'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: true,
  });
  return localTime; // Output: "09:03:56 PM"

}

function ConfirmMatchList({
    matches, 
    loading,
    User,
    refreshMatches,
    showConfirmMatches
  }) {

  const [submittedLikes, setSubmittedLikes] = useState({}); // Track submitted likes
  // Add distance to each person and sort by distance
  if (!loading && User) {
    matches = matches.map((match) => ({
      ...match,
      distance: Math.round(
        getPreciseDistance(
          { latitude: User.lat, longitude: User.long },
          { latitude: match.person.lat, longitude: match.person.long }
        ) / 1609.34 * 10
      ) / 10, // Convert meters to miles and round to 1 decimal place
    }));

    // Sort people by distance (smallest to largest)
  matches.sort((a, b) => a.accepted - b.accepted);;
  }

  useEffect(() => {
    console.log("Pending updated:", matches);
  }, [matches, showConfirmMatches]);
  
  const handleSubmit = async (User, person, match) => {
    const payload = {
      id: match.id,
      match_ids: [User.id, person.id],
      offered: person.id,
      offered_time: match.offered_time,
      accepted: User.id
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
        alert('Match completed successfully!');
        setSubmittedLikes((prev) => ({ ...prev, [person.id]: true })); // Mark as submitted
        refreshMatches(response); // Refresh matches in ChatSelect
      } else {
        alert('Failed to add match.');
      }
    } catch (error) {
      console.error('Error submitting form:', error);
      alert('An error occurred while submitting the form.');
    }
  };
    return (
      loading || User==undefined ? (
        <div className="d-flex align-items-center">
          <strong role="status">Loading...</strong>
          <div className="spinner-border ms-auto" aria-hidden="true"></div>
        </div>
      ) : (
      <>
      <h2>These cool folks liked your profile:</h2>
      <Accordion
       show={showConfirmMatches}
       className='w-100'
      >
        {matches.map((match) => (
        <Accordion.Item className='w-100' eventKey={match.person.id} key={match.person.id}>
          <Accordion.Header className='w-100' key={`${match.person.id}100`}>
            <img
              src={match.person.profile ? match.person.profile : '/profile.svg'}
              style={{ borderRadius: '50%' }}
              className="m-1 p-1"
              height="50"
              width="50"
              alt={`${match.person.name}'s profile`}
            />
            <strong>{match.person.name}</strong> — {match.person.motto} - Distance: { match.person.distance } miles - liked your profile on { convertISODateToLocal(match.offered)}
          </Accordion.Header>
          <Accordion.Body key={`${match.person.id}10000`}>
            <div className='text-start'>
              <img src={match.person.profile} className='m-1 p-1 rounded float-start' height={'250'} width={'250'} />
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

              <div height='450' width='450' className='text-center w-fill'>
                <ControlledCarousel
                  key={`${match.person.id}1000000`}
                  id={match.person.id}
                />
              </div>
              <div className="mx-auto w-50 d-flex flex-row justify-content-around align-items-center">
                {submittedLikes[match.person.id] ? (
                  <p className="text-success">Like submitted!</p>
                ) : (
                  <>
                    <Button
                      variant="primary"
                      className="m-2"
                      onClick={() => handleSubmit(User, match.person, match)}
                    >
                      Great match - let's chat!
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
      </>
      )
    
    )
  
}

export default ConfirmMatchList