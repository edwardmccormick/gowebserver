import { useEffect, useState } from 'react';
import NavBar from './components/navbar';
import MatchList from './components/matchlist';
import FormTextExample from './components/login';
import { BasicExample } from './components/login';


function App() {
  const [people, setPeople] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/people')
      .then((res) => res.json())
      .then((data) => {
        setPeople(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  console.log(people);
  console.log(loading);

  return (

    <div className='mx-auto p-3 text-center'>
      <NavBar />
      <h1 className='m-3 p-3 text-center'>Found {people.length} matches for you!</h1>
      <>
      <div className='container text-center'>
        <MatchList 
          peopleObject={{ people }} 
          loading={loading} />
      </div>
      </>
      {/* <>
      {loading ? (
        <div className="d-flex align-items-center">
          <strong role="status">Loading...</strong>
          <div className="spinner-border ms-auto" aria-hidden="true"></div>
        </div>
      ) : (
      <Accordion className='w-100'>
        {people.map((person) => (
        <Accordion.Item className='w-100' eventKey={person.id} key={person.id}>
          <Accordion.Header className='w-100' key={person.id+100}><img src="/profile-250x250.jpg" style={{borderRadius: '50%'}} className='m-1 p-1' height='50' width='50'/><strong>{person.name}</strong> — {person.motto} - Distance: {person.distance} miles -  Likes Dogs: {person.dogs}</Accordion.Header>
          <Accordion.Body key={person.id+10000}>
            <div>
              <img src={person.profile} className='m-1 p-1 rounded float-start' height={'250'} width={'250'} />
              <p>Well uh, good, fine. No, Biff, you leave her alone. Thanks a lot, kid. Yeah, sure, okay. A block passed Maple, that's John F. Kennedy Drive.</p>
              <p>Right. Hey, not too early I sleep in on Saturday. Oh, McFly, your shoe's untied. Don't be so gullible, McFly. You got the place fixed up nice, McFly. I have you're car towed all the way to your house and all you've got for me is light beer. What are you looking at, butthead. Say hi to your mom for me. He's fine, and he's completely unaware that anything happened. As far as he's concerned the trip was instantaneous. That's why Einstein's watch is exactly one minute behind mine. He skipped over that minute to instantly arrive at this moment in time. Come here, I'll show you how it works. First, you turn the time circuits on. This readout tell you where you're going, this one tells you where you are, this one tells you where you were. You imput the destination time on this keypad. Say, you wanna see the signing of the declaration of independence, or witness the birth or Christ. Here's a red-letter date in the history of science, November 5, 1955. Yes, of course, November 5, 1955. Jennifer. Thank god I still got my hair. What on Earth is that thing I'm wearing?</p>
              <p>Right. No. You cost three-hundred buck damage to my car, you son-of-a-bitch. And I'm gonna take it out of your ass. Hold him. Hey Biff, check out this guy's life preserver, dork thinks he's gonna drown. Yeah Mom, we know, you've told us this story a million times. You felt sorry for him so you decided to go with him to The Fish Under The Sea Dance.</p>
              <div height='450' width='450' className='text-center w-fill'>
                <ControlledCarousel key={person.id+1000000}/>
              </div>
            </div>
            </Accordion.Body>
        </Accordion.Item>
        ))}
      </Accordion>
      )}
      </> */}
      {/* <div className="w-50 p-3 m-3">
        <ControlledCarousel/>
      </div> */}
      <div className="w-50 p-3 m-3">
        <FormTextExample/>
      </div>
      <h2 className='m-3 p-3 text-center'>People Add</h2>
      <BasicExample />
    </div>
  );
}

export default App;