import React, { useState, useMemo } from 'react';
import { Card, Form, Button, Accordion, Row, Col, Badge } from 'react-bootstrap';

// --- Inlined Data ---
// Inlining the data from details.json to resolve the import error.
const detailsData = {
  "dogs": [
    "I tell people I’m allergic, but just to be polite. Dogs are chaos gremlins.",
    "They smell weird and stare at me like I owe them money.",
    "If one jumps on me, I will file a formal complaint—with HR or God.",
    "They’re fine. From a distance. Like mountains. Or relatives.",
    "I pet them. I move on. That’s it.",
    "As long as he doesn’t shit on the floor or bark at ghosts, we’re cool.",
    "I might sneak one a piece of cheese. Maybe.",
    "I wave at dogs from my car. That’s where I’m at in life.",
    "I have dog biscuits in my pocket right now. Just in case.",
    "I show people photos of my dog before my family.",
    "I love dogs more than people. I’d run into traffic for a stranger’s golden retriever."
  ],
  "cats": [
    "They knock things off counters and my will to live.",
    "I see a cat and I see a little liar plotting something.",
    "You mean land piranhas with fur? No thanks.",
    "They’re cute on the internet. In real life? Meh.",
    "I respect their boundaries. Mostly by avoiding them.",
    "We coexist like roommates who don’t talk.",
    "I once let one nap on me for three hours. Regret nothing.",
    "I’m not a ‘cat person,’ but I have favorites.",
    "If I could nap like a cat, I’d be unstoppable.",
    "I have cat hair on everything I own and I love it.",
    "If a cat chooses me, I cancel plans. For days."
  ],
  "kids": [
    "Nope. Not even the quiet ones. Especially not the quiet ones.",
    "I hear a child scream and I start Googling vasectomy clinics.",
    "I can handle them in 5-minute increments. With noise-canceling headphones.",
    "They’re okay when they’re sleeping or not mine.",
    "I like other people’s kids... over there... doing something quiet.",
    "As long as they don’t wipe snot on me, we’re solid.",
    "Some of them are funny. I’ve met like three cool ones.",
    "They’re little weirdos, but they mean well.",
    "I will proudly wear the ‘fun aunt/uncle/friend’ title forever.",
    "I cry during kids’ talent shows. Every time.",
    "I would die for them. And I also pack snacks. Always."
  ],
  "smoking": [
    "I can smell that you lit one in 2019, and it still gives me the ick.",
    "If I wanted to taste an ashtray, I’d chew on a grill grate.",
    "It’s a no. A full-body, wrinkly-nose no.",
    "It’s your lungs, I guess. But I’m judging you silently.",
    "I’m not mad, just disappointed. And stepping three feet away.",
    "It used to be cool in movies. That’s about it.",
    "If you’re outside and it’s breezy, I might forgive you.",
    "I’ve bummed one at 2 AM in college. Haven’t we all?",
    "The occasional drag is part of my mysterious European persona.",
    "I carry lighters even when I don’t smoke. It’s a vibe.",
    "I live in a smokestack and still burn through a pack a day."
  ],
  "drinking": [
    "I’ve read the AA Big Book. Twice. For fun.",
    "The smell of tequila makes me remember things I legally can’t discuss.",
    "I’ll toast with seltzer and be in bed by nine.",
    "I drink… when forced by social obligation and peer pressure.",
    "I like a glass of wine. At weddings. Maybe.",
    "Give me one cocktail and a two-drink limit. I’m done.",
    "I drink socially, responsibly, and with snacks.",
    "A good Old Fashioned could solve 90% of my problems.",
    "Brunch without mimosas is just breakfast with judgment.",
    "I have a signature drink. And a backup. And a playlist.",
    "I have themed barware, drink recipes memorized, and stories I’ll never tell sober."
  ],
  "religion": [
    "I spontaneously combust in church parking lots.",
    "I avoid nativity scenes like they owe me money.",
    "I don’t knock it—but it’s not my thing. At all.",
    "I like the architecture. That’s about it.",
    "I’ll go for weddings, funerals, or free donuts.",
    "I’m spiritual-ish. Depends on the day and my hangover.",
    "I believe in something bigger, but I’m still figuring it out.",
    "Faith matters. So does asking questions. I like both.",
    "I pray. I volunteer. I try not to be a jerk.",
    "My life has a spiritual center. And I can quote scripture *and* memes.",
    "I love God, go to church, and yes, I’ll pray for you—and mean it."
  ],
  "politics": [
    "I vote. I read the news. I try to stay sane.",
    "I’m not a fan of politics, but I’m not blind to it.",
    "I have opinions. They’re just not for public consumption.",
    "I’m registered, but I avoid debates like the plague.",
    "I’m a moderate. I like to keep my blood pressure down.",
    "I vote based on issues, not parties. It’s complicated.",
    "I’m politically aware, but I don’t engage in Twitter wars.",
    "I believe in democracy, but I also believe in naps.",
    "I’ll discuss politics over coffee, but not at family dinners.",
    "I’m informed, but I prefer to keep my sanity intact.",
    "I care about the future, but I also care about my mental health."
  ],
  "food": [
    "I eat to live, not live to eat. Mostly.",
    "I have a favorite pizza place. And a backup.",
    "I’m not picky, but I have standards. And allergies.",
    "I’ll try anything once. Except olives. Or beets.",
    "I love food trucks. They’re like mobile art galleries.",
    "I’m a foodie, but I also enjoy a good microwave meal.",
    "I cook sometimes. Mostly when I’m hungry and there’s nothing else.",
    "I love brunch. It’s like breakfast, but with mimosas.",
    "I’ll eat leftovers for days. Waste not, want not.",
    "I have a sweet tooth. And a savory tooth. And a snack tooth.",
    "Food is life. And I’m living it one bite at a time."
  ],
  "energy_levels": [
    "I am horizontal. Emotionally and physically.",
    "I’m awake, technically. That’s it.",
    "I moved today. Once. It was overrated.",
    "If I stand up too fast, I need a snack and a nap.",
    "I do things. Slowly. With dramatic sighs.",
    "I can rally when there’s coffee, chaos, or drama.",
    "Moderate pep. Might break into dance if provoked.",
    "I’m always moving and have 3 side quests at all times.",
    "I regularly outpace golden retrievers on espresso.",
    "People ask if I’m on something. I’m not. This is natural.",
    "I radiate main-character energy at all times. No chill."
  ],
  "outdoorsy_ness": [
    "Sunlight is my sworn enemy. Nature is a scam.",
    "I like the *idea* of outdoors. From inside.",
    "I will sit on a patio. With shade. And Wi-Fi.",
    "I’ll hike to the mailbox. Maybe.",
    "I can camp—as long as there’s plumbing nearby.",
    "I go outside sometimes. Mostly for snacks and sanity.",
    "Give me trails and a solid playlist and I’m in.",
    "I own gear. Real gear. I use it too.",
    "I identify as 'sun-drenched and bug-bitten.'",
    "If it’s not muddy, is it even a real adventure?",
    "I wrestled a bear once. It was consensual. We’re friends now."
  ],
  "travel": [
    "I don’t even like leaving my bed, let alone my zip code.",
    "A staycation sounds risky. What if I lose my spot on the couch?",
    "I went to IKEA once and that was enough international travel for me.",
    "I’ll go on a trip if someone else plans it and I can bring snacks.",
    "I like weekend getaways, preferably with hot tubs and robes.",
    "Give me a three-day weekend and I’ll disappear to somewhere cute.",
    "I love road trips, train rides, weird motels and playlists.",
    "Passport ready, suitcase half-packed, I live for it.",
    "I collect passport stamps like Pokémon cards.",
    "I’ve slept in airports and danced in three continents this year.",
    "I don’t live anywhere—I just visit different places for tax reasons."
  ],
  "bougieness": [
    "I microwave water for tea and use ketchup packets as currency.",
    "Generic everything. I once reused floss.",
    "If it’s on sale *and* has a coupon, I’ll consider it.",
    "I have a Costco membership and I’m not afraid to flex it.",
    "I don’t *need* fancy, but I like a touch of nice.",
    "I’ll splurge on skincare, snacks, and soft things.",
    "Mid-level luxury: I know my wine isn’t from a box, and that’s enough.",
    "If it’s artisanal and unnecessarily expensive, I probably bought it.",
    "I have opinions about thread counts, balsamic vinegar, and skincare routines.",
    "I only cry in boutique hotels. With robes. And mood lighting.",
    "I’m one charcuterie board away from full-on Real Housewife."
  ],
  "importance_of_politics": [
    "I don’t vote. I don’t even read the labels on soup cans.",
    "Politics? That’s the thing with the flags, right?",
    "I vibe with ‘live and let live’ and ignore the news entirely.",
    "I kind of know who the president is. I think.",
    "I vote, but mostly because my grandma would haunt me if I didn’t.",
    "I have opinions, but I only argue about them with my group chat.",
    "I keep it civil, but I *will* fact-check you in real time.",
    "I read policy proposals for fun and yell at podcasts in my car.",
    "I organize, petition, and once cried at a city council meeting.",
    "I attend protests, quote legislation, and probably have a favorite senator.",
    "My political views are a core personality trait. I’m basically a campaign sticker."
  ]
};


// --- Helper Components & Style ---

// Inlined SVG icons to replace the 'react-icons' dependency
const FiSliders = ({ className }) => (<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}><line x1="4" y1="21" x2="4" y2="14"></line><line x1="4" y1="10" x2="4" y2="3"></line><line x1="12" y1="21" x2="12" y2="12"></line><line x1="12" y1="8" x2="12" y2="3"></line><line x1="20" y1="21" x2="20" y2="16"></line><line x1="20" y1="12" x2="20" y2="3"></line><line x1="1" y1="14" x2="7" y2="14"></line><line x1="9" y1="8" x2="15" y2="8"></line><line x1="17" y1="16" x2="23" y2="16"></line></svg>);
const FiRefreshCw = ({ className }) => (<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}><polyline points="23 4 23 10 17 10"></polyline><polyline points="1 20 1 14 7 14"></polyline><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path></svg>);
const FiSearch = ({ className }) => (<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>);
const FiMapPin = ({ className }) => (<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>);

// Injects custom CSS into the document head for styling.
const Styles = () => (
  <style type="text/css">
    {`
    .advanced-search-card {
      background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
      border: none;
      box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.37);
      backdrop-filter: blur(4px);
      -webkit-backdrop-filter: blur(4px);
      border-radius: 20px;
    }
    .advanced-search-card .accordion-item {
      background-color: rgba(255, 255, 255, 0.8);
      border: 1px solid rgba(255, 255, 255, 0.18);
      border-radius: 15px !important;
      margin-bottom: 1rem;
    }
    .advanced-search-card .accordion-header button {
      border-radius: 15px !important;
      background-color: transparent !important;
      color: #343a40 !important;
      font-weight: 600;
    }
    .advanced-search-card .accordion-header button:not(.collapsed) {
       box-shadow: none;
       border-bottom: 1px solid #dee2e6;
    }
    .advanced-search-card .accordion-body {
       padding: 1.5rem;
    }
    .search-btn {
      background: linear-gradient(45deg, #28a745, #218838);
      border: none;
      font-size: 1.2rem;
      font-weight: bold;
      padding: 0.75rem 1.5rem;
      border-radius: 50px;
      transition: all 0.3s ease;
      box-shadow: 0 4px 15px rgba(40, 167, 69, 0.4);
    }
    .search-btn:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(40, 167, 69, 0.6);
    }
    .form-range::-webkit-slider-thumb {
      background-color: #007bff;
    }
    .form-range::-moz-range-thumb {
      background-color: #007bff;
    }
    .form-switch .form-check-input:checked {
      background-color: #28a745;
      border-color: #28a745;
    }
    .preference-text {
        font-style: italic;
        color: #6c757d;
        background-color: #f8f9fa;
        padding: 0.75rem;
        border-radius: 10px;
        min-height: 80px;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    `}
  </style>
);


// --- Main Component ---

function GeminiAdvancedSearch({ onSearch }) {
  // Use the inlined data object
  const searchCategories = useMemo(() => ({ ...detailsData }), []);

  // Function to create the initial state for all criteria
  const getInitialCriteria = () => {
    const initialState = {};
    for (const key in searchCategories) {
      initialState[key] = {
        enabled: false,
        preference: 5, // Default to the middle of the 0-10 scale
        flexibility: 2, // Default to a flexibility of +/- 2
      };
    }
    return initialState;
  };

  const [distance, setDistance] = useState(50);
  const [criteria, setCriteria] = useState(getInitialCriteria());

  // A single handler to update any part of our criteria state
  const handleCriteriaChange = (category, field, value) => {
    setCriteria(prev => ({
      ...prev,
      [category]: {
        ...prev[category],
        [field]: value,
      },
    }));
  };

  // Resets all filters to their default state
  const handleReset = () => {
    setDistance(50);
    setCriteria(getInitialCriteria());
  };

  // Gathers all state and passes it up to the App component
  const handleSearch = () => {
    onSearch({ distance, criteria });
  };

  return (
    <>
      <Styles />
      <Card className="advanced-search-card p-3 p-md-4">
        <Card.Header className="bg-transparent border-0 d-flex justify-content-between align-items-center mb-3">
          <h2 className="mb-0 d-flex align-items-center"><FiSliders className="me-3" /> Find Your Match</h2>
          <Button variant="outline-secondary" size="sm" onClick={handleReset} className="d-flex align-items-center">
            <FiRefreshCw className="me-2" /> Reset
          </Button>
        </Card.Header>

        <Card.Body>
          {/* --- Distance Slider --- */}
          <div className="mb-5">
            <Form.Label className="fw-bold fs-5 d-flex align-items-center"><FiMapPin className="me-2"/>Distance</Form.Label>
            <Row className="align-items-center">
                <Col xs={9}>
                    <Form.Range
                        min={1}
                        max={500}
                        value={distance}
                        onChange={(e) => setDistance(Number(e.target.value))}
                    />
                </Col>
                <Col xs={3} className="text-end">
                    <Badge pill="true" bg="primary" className="fs-6 px-3 py-2">{distance} miles</Badge>
                </Col>
            </Row>
          </div>

          {/* --- Criteria Accordion --- */}
          <Accordion>
            {Object.entries(criteria).map(([key, value], index) => {
              const categoryTitle = key.replace(/_/g, ' ');
              const descriptions = searchCategories[key] || [];

              return (
                <Accordion.Item eventKey={index.toString()} key={key}>
                  <Accordion.Header>
                    <span className="text-capitalize">{categoryTitle}</span>
                    {value.enabled && (
                      <Badge bg="success" className="ms-3">
                        {value.preference} ± {value.flexibility}
                      </Badge>
                    )}
                  </Accordion.Header>
                  <Accordion.Body>
                    <Form.Group as={Row} className="mb-4 align-items-center">
                      <Col sm={9}>
                        <h5 className="mb-0">Filter by this preference?</h5>
                      </Col>
                       <Col sm={3} className="d-flex justify-content-end">
                          <Form.Check
                            type="switch"
                            id={`switch-${key}`}
                            checked={value.enabled}
                            onChange={(e) => handleCriteriaChange(key, 'enabled', e.target.checked)}
                          />
                      </Col>
                    </Form.Group>
                    
                    {value.enabled && (
                      <>
                        <p className="preference-text text-center my-3">
                            "{descriptions[value.preference] || 'No description'}"
                        </p>
                        
                        {/* Preference Slider */}
                        <Form.Group className="mb-4">
                          <Form.Label className="fw-bold">My Preference: {value.preference}</Form.Label>
                          <Form.Range
                            min={0}
                            max={descriptions.length - 1}
                            value={value.preference}
                            onChange={(e) => handleCriteriaChange(key, 'preference', Number(e.target.value))}
                          />
                        </Form.Group>

                        {/* Flexibility Slider */}
                        <Form.Group>
                          <Form.Label className="fw-bold">My Flexibility: ±{value.flexibility}</Form.Label>
                          <Form.Range
                            min={0}
                            max={5} // Max flexibility of 5 levels on each side
                            value={value.flexibility}
                            onChange={(e) => handleCriteriaChange(key, 'flexibility', Number(e.target.value))}
                          />
                        </Form.Group>
                      </>
                    )}
                  </Accordion.Body>
                </Accordion.Item>
              );
            })}
          </Accordion>
        </Card.Body>

        <Card.Footer className="bg-transparent border-0 mt-4 text-center">
          <Button className="search-btn w-75" onClick={handleSearch}>
            <FiSearch className="me-2" /> Search
          </Button>
        </Card.Footer>
      </Card>
    </>
  );
}

export default GeminiAdvancedSearch;
