import Accordion from 'react-bootstrap/Accordion';
import ControlledCarousel from './carousel'
import { getPreciseDistance } from 'geolib';
import { User } from 'lucide-react';

function MatchList({
    peopleObject, 
    loading,
    User
}) {let people = peopleObject.people;
    console.log(people);
    return (
      loading || User==undefined ? (
        <div className="d-flex align-items-center">
          <strong role="status">Loading...</strong>
          <div className="spinner-border ms-auto" aria-hidden="true"></div>
        </div>
      ) : (
      <Accordion className='w-100'>
        {people.map((person) => (
        <Accordion.Item className='w-100' eventKey={person.id} key={person.id}>
          <Accordion.Header className='w-100' key={person.id+100}>
            <img
              src={person.profile ? person.profile : '/profile.svg'}
              style={{ borderRadius: '50%' }}
              className="m-1 p-1"
              height="50"
              width="50"
              alt={`${person.name}'s profile`}
            />
            <strong>{person.name}</strong> — {person.motto} - Distance: { Math.round(getPreciseDistance( {latitude: User.lat, longitude: User.long}, {latitude: person.lat, longitude: person.long})/1609.34*10)/10} miles</Accordion.Header>
          <Accordion.Body key={person.id+10000}>
            <div className='text-start'>
              <img src={person.profile} className='m-1 p-1 rounded float-start' height={'250'} width={'250'} />
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
                  key={`${person.id}1000000`}
                  id={person.id}
                />
              </div>
            </div>
            </Accordion.Body>
        </Accordion.Item>
        ))}
      </Accordion>
      )
    
    )
}

export default MatchList