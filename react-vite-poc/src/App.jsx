import { useEffect, useState } from 'react';
import NavBar from './components/navbar';
import MatchList from './components/matchlist';
import SignIn from './components/login';
import { SignUpProfile } from './components/login';
import verbiage from '../../verbiage.json';
import FAQ from './components/faq';
import 'bootstrap/dist/css/bootstrap.min.css';    

function App() {
  const [people, setPeople] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loggedInUser, setLoggedInUser] = useState(null); // State to store the logged-in user
  const [jwt, setJWT] = useState(null);

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
      <NavBar
        profile={loggedInUser?.profile}
        username={loggedInUser?.name}
        setLoggedInUser={setLoggedInUser}
        setJWT={setJWT}
      />
      <h1 className='m-3 p-3 text-center'>Found {people.length} matches for you!</h1>
      <>
      <div className='container text-center'>
        <MatchList 
          peopleObject={{ people }} 
          loading={loading} />
      </div>
      </>
      {/* <WebSocketChat /> */}
      <div className="w-50 p-3 m-3">
        { loggedInUser ? 
        <SignIn 
          setLoggedInUser={setLoggedInUser} 
          setJWT={setJWT} 
          />
          : null
        }
        
      </div>
      <h2 className='m-3 p-3 text-center'>People Add</h2>
      <SignUpProfile />
      <div className='m-3 p-3'>
        <FAQ className='p-2 m-2'/>
      </div>
      { jwt ? <p>JWT: {jwt}</p> : null}
    </div>
  );
}

export default App;