import { useEffect, useState } from 'react';
import NavBar from './components/navbar';
import MatchList from './components/matchlist';
import SignIn from './components/login';
import { SignUpProfile } from './components/login';
import FAQ from './components/faq';
import 'bootstrap/dist/css/bootstrap.min.css';    
import ClaudeAdvancedSearch from './components/advancedsearchclaude';
import ConfirmMatchList from './components/confirmmatch';
import SignUp from './components/signup'; 
import { LogIn } from 'lucide-react';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css'; // Import your CSS file

function App() {
  const [people, setPeople] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loggedInUser, setLoggedInUser] = useState(null); // State to store the logged-in user
  const [jwt, setJWT] = useState(null);
  const [searchResults, setSearchResults] = useState([]);
  const [matches, setMatches] = useState([]);
  const [matchLoading, setMatchLoading] = useState(false);
  const [pendings, setPendings] = useState([]);
  const [offereds, setOffereds] = useState([]);
  const [showConfirmMatch, setShowConfirmMatch] = useState(true)
  const [showAdvancedSearch, setShowAdvancedSearch] = useState(false); // Toggle ClaudeAdvancedSearch visibility
  const [showFAQ, setShowFAQ] = useState(false);
  const [animationStarted, setAnimationStarted] = useState(false); // Controls the logo animation
  const [showText, setShowText] = useState(false); // Controls the visibility of the <h1>
  const [signUpFlow, setSignUpFlow] = useState(false); // Controls the visibility of the sign-up flow

  useEffect(() => {
    // Start the animation after 3 seconds
    const timer = setTimeout(() => {
      setAnimationStarted(true);
      // Show the text after the animation starts
      setTimeout(() => setShowText(true), 1000); // Delay for the text to appear
    }, 3000);

    return () => clearTimeout(timer); // Cleanup the timer
  }, []);

  useEffect(() => {
    if (loggedInUser == null || loggedInUser == undefined) {
      setLoading(false);
      setMatchLoading(false);
      setPeople([]);
      setMatches([]);
      setPendings([]);
      setOffereds([]);
        return;
      }; 
    fetch('http://localhost:8080/people')
      .then((res) => res.json())
      .then((data) => {
        setPeople(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [loggedInUser]);

    // Add this useEffect hook
  useEffect(() => {
    
    if (loggedInUser) { 
      // Scroll to the top of the page when the user logs in
      window.scrollTo({ top: 0, behavior: 'smooth' });
      // Only call refreshMatches if loggedInUser is not null
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

  return ( !loggedInUser ? (
    <>
      <div className="mx-auto p-3 text-center bg-black h-100" style={{height: '100vh', width: '100vw'}}>
      
      
        <img src={"./urmid.svg"} className={`w-50 mx-auto text-center frontlogo ${animationStarted ? 'swoop' : ''}`} /> 
         {showText && (
        <h1 className="text-white animated-text">
          urmid - find love so you Go. Away.
        </h1>
      )}
          <div className='m-2 p-2 w-50 bg-white mx-auto text-center rounded form-container'>
            <div className="form-content">
            <h1 className='bg-white'>And we all talk about it behind your back.</h1>
              <h6 className='bg-white'>Look, no one is excited about this, so just make it quick. <br />
              Two roads diverged in the woods....yadda yadda. Pick what fits you best:</h6>
              <div className='d-flex justify-content-around align-items-center'>
                <div className='col-5'>
                  <SignIn 
                    setLoggedInUser={setLoggedInUser} 
                    setJWT={setJWT}
                  />   
                </div>   
                <div><h4>Or</h4></div>
                <div className='col-5'>
                  <SignUp />
                </div>
              </div></div>
          </div>
                 <br />
       <br />
       <br />
       <br />
       <br />
       <br />
       <p className="text-black">God this shit is so stupid, does it even work? A free dating site. Yes, totally fucking free. We use ads because honestly, you kind of piss us off and we're hoping you find the love of your life and get married. Because then you'll leave us alone. Asshole.</p>
       </div>

  </>
  ) : (

  <>
    <div className='mx-auto text-center'>
      <NavBar
        User={loggedInUser}
        setLoggedInUser={setLoggedInUser}
        setJWT={setJWT}
        jwt={jwt}
        refreshMatches={refreshMatches}
        matches={matches}
        pendings={pendings}
        offereds={offereds}
        setShowConfirmMatch={setShowConfirmMatch}
        onSearchClick={() => {
          setShowAdvancedSearch(true)
          setShowConfirmMatch(false);
          setShowFAQ(false);
        }}
        onMeetClick={() => {
          setShowAdvancedSearch(false)
          setShowConfirmMatch(true);
          setShowFAQ(false);
        }}
        onFAQClick={() => {
          setShowFAQ(true)
          setShowConfirmMatch(false);
          setShowAdvancedSearch(false);
        }}
      />

       <div className={`text-center mx-auto container m-4 p-2 fade-container ${!showConfirmMatch ? 'hidden' : 'visible'}`}>
        <ConfirmMatchList
          key='ConfirmMatchList'
          matches={pendings}
          loading={matchLoading}
          User={loggedInUser}
          refreshMatches={refreshMatches}
          className='m-5 p-4'
        />

        <br />
        
        <br />
        
        <br />

        <MatchList 
          people={people} 
          loading={loading}
          User = {loggedInUser}
          refreshMatches= {refreshMatches}
        />
      </div>



        <>
      <div className={`text-center mx-auto container m-5 p-4 fade-container ${!showAdvancedSearch ? 'hidden' : 'visible'}`}>
        <h1 className='m-3 p-3 text-center'>Advanced Search - Claude Remix</h1>
        <ClaudeAdvancedSearch onSearch={handleSearch} />
      </div>
      </>

    <>
      <h2 className='m-3 p-3 text-center'>People Add</h2>
      <SignUpProfile />
      <div className={`m-3 p-3fade-container ${!showFAQ ? 'hidden' : 'visible'}`}>
        
        <FAQ className='p-2 m-2'/>
      </div>
      </>
      { jwt ? <p>JWT: {jwt}</p> : null }
    </div>
  </>
  
      
    
  ));
}

export default App;