import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import PeopleList from './components/PeopleList';
import PersonDetail from './components/PersonDetail';
import Greeting from './components/Greeting';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Greeting />} />
        <Route path="/people" element={<PeopleList />} />
        <Route path="/people/:id" element={<PersonDetail />} />
      </Routes>
    </Router>
  );
}

export default App;