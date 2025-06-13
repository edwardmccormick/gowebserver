import { useEffect, useState } from 'react';
import NavBar from './components/navbar';
import MatchList from './components/matchlist';
import FormTextExample from './components/login';
import { SignUpProfile } from './components/login';
import verbiage from '../../verbiage.json';
import FAQ from './components/faq';
import 'bootstrap/dist/css/bootstrap.min.css';    

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
      <NavBar
        profile='https://avatars.githubusercontent.com/u/102410?v=4'
        username='urmid'

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
        <FormTextExample/>
      </div>
      <h2 className='m-3 p-3 text-center'>People Add</h2>
      <SignUpProfile />
      <div className='m-3 p-3'>
        <FAQ className='p-2 m-2'/>
      </div>
    </div>
  );
}

export default App;