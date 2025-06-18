import { useEffect, useState } from 'react';
import NavBar from './components/navbar';
import MatchList from './components/matchlist';
import SignIn from './components/login';
import { SignUpProfile } from './components/login';
import FAQ from './components/faq';
import 'bootstrap/dist/css/bootstrap.min.css';    
import ClaudeAdvancedSearch from './components/advancedsearchclaude';
import ConfirmMatchList from './components/confirmmatch';

function App() {
  const [people, setPeople] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loggedInUser, setLoggedInUser] = useState(null); // State to store the logged-in user
  const [jwt, setJWT] = useState(null);
  const [searchResults, setSearchResults] = useState([]);
  const [matches, setMatches] = useState([]);
  const [matchLoading, setMatchLoading] = useState(true);
  const [pendings, setPendings] = useState([]);
  const [offereds, setOffereds] = useState([]);
  const [showConfirmMatch, setShowConfirmMatch] = useState(false)

  useEffect(() => {
    fetch('http://localhost:8080/people')
      .then((res) => res.json())
      .then((data) => {
        setPeople(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

    // Add this useEffect hook
  useEffect(() => {
    if (loggedInUser) { // Only call refreshMatches if loggedInUser is not null
      refreshMatches();
    }
  }, [loggedInUser]); // Dependency array: this effect runs when loggedInUser changes


  const refreshMatches = () => {
      if (loggedInUser == null || loggedInUser == undefined) return; 
      else if (matches == null || matches == undefined) {console.log("This needs a fetch from refreshMatches")}
    fetch(`http://localhost:8080/matches/${loggedInUser.id}`)
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        const offered = data.filter((match) => (match.accepted_time == "0001-01-01T00:00:00Z" && match.offered == loggedInUser.id)); // This is what a null date looks like in Go
        console.log(offered);
        const accepted = data.filter((match) => match.accepted_time != "0001-01-01T00:00:00Z"); // A non-null date
        console.log(accepted);
        const pending = data.filter((match) => (match.accepted_time == "0001-01-01T00:00:00Z" && match.offered != loggedInUser.id))
        console.log(pending);
        setMatches(accepted);
        setPendings(pending);
        setOffereds(offered);
        setMatchLoading(false);
      })
      .catch(() => setMatchLoading(false));
  }; 

  // The Chat GPT component
  // const handleSearch = ({ distance, criteria }) => {
  //   const results = people.filter((person) => {
  //     const withinDistance = person.distance <= distance;
  //     const matchesCriteria = Object.entries(criteria).every(
  //       ([key, value]) => person[key] === parseInt(value)
  //     );
  //     return withinDistance && matchesCriteria;
  //   });
  //   setSearchResults(results);
  // };

  const handleSearch = (searchState) => {
    const { distance, criteria } = searchState;

    // This console.log is helpful for debugging!
    console.log("Searching with:", { distance, criteria });

    const results = people.filter((person) => {
      // 1. Check distance
      // Make sure the person object has a 'distance' property
      const personDistance = person.distance || Infinity;
      if (personDistance > distance) {
        return false;
      }

      // 2. Check all enabled criteria
      for (const key in criteria) {
        const setting = criteria[key];

        // Only filter if the criterion is enabled by the user
        if (setting.enabled) {
          const personValue = person[key]; // The person's value for this trait (e.g., 7)

          // If the person doesn't have this trait defined, they can't match.
          if (personValue === undefined) {
            return false;
          }

          const min = setting.preference - setting.flexibility;
          const max = setting.preference + setting.flexibility;

          // Check if the person's value is within the user's desired range
          if (personValue < min || personValue > max) {
            return false; // This person is outside the range, so we exclude them.
          }
        }
      }

      // 3. If the person passed the distance and all enabled criteria checks, include them!
      return true;
    });

    setSearchResults(results);
    // You might want to update the main people list to show only search results
    // setPeople(results); 
  };

  return (

    <div className='mx-auto p-3 text-center'>
      <NavBar
        User={loggedInUser}
        setLoggedInUser={setLoggedInUser}
        setJWT={setJWT}
        refreshMatches={refreshMatches}
        matches={matches}
        pendings={pendings}
        offereds={offereds}
        setShowConfirmMatch={setShowConfirmMatch}
      />

      {/* <h1 className="m-3 p-3 text-center">Advanced Search</h1>
      <AdvancedSearch onSearch={handleSearch} />
      <br />
      <br />
      <h1 className="m-3 p-3 text-center">Advanced Search - Gemini Remix</h1>
      <GeminiAdvancedSearch onSearch={handleSearch} />
      <h1 className='m-3 p-3 text-center'>Found {people.length} matches for you!</h1>
      <br />
      <br /> */}

      <ConfirmMatchList
        key={'ConfirmMatchList'}
        matches={pendings}
        loading={matchLoading}
        User={loggedInUser}
        refreshMatches={refreshMatches}
        showConfirmMatch={showConfirmMatch}
      />
      <h1 className="m-3 p-3 text-center">Advanced Search - Claude Remix</h1>
      <ClaudeAdvancedSearch onSearch={handleSearch} />
      <br />
      <>
      <div className='container text-center'>
        <MatchList 
          people={people} 
          loading={loading}
          User = {loggedInUser}
          refreshMatches= {refreshMatches
          }
        />
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