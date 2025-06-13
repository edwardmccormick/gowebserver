import Accordion from 'react-bootstrap/Accordion';

function FAQ() {
  return (
    <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
        <Accordion.Header>Wtf even is this?</Accordion.Header>
        <Accordion.Body>
          Man I don't know. It's a dating app. We thought most dating apps sucked, and surely there was a better way.<br /><br />
          Fucked up pretty bad, didn't we?
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header>So everything is free?</Accordion.Header>
        <Accordion.Body>
          Mostly. We leaned into enshittifcation early, so we've got a lot of 'freemium' components.
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>That sucks.</Accordion.Header>
        <Accordion.Body>
          Yeah, we modeled it after your mom.
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>Rude.</Accordion.Header>
        <Accordion.Body>
          Yeah, sorry. If you think our app sucks, you should see the way we run this stupid thing as a business.
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="4">
        <Accordion.Header>Can we start over? So....what's free?</Accordion.Header>
        <Accordion.Body>
          Great question! <br />
          <br />
          It's free to sign up, free to create a profile, free to look at other profiles, free to 'match' with other people, and free to chat.<br />
          <br />
          It's free to ask somebody out on a date, with our patented Datamatic 5000 (it's basically Yelp, with an AI model).
        </Accordion.Body>
        </Accordion.Item>
      <Accordion.Item eventKey="4">
        <Accordion.Header>That....kind of sounds like everything?</Accordion.Header>
        <Accordion.Body>
          Yeah, we weren't kidding. Most things are free!<br />
            <br />
          <strong>BUT</strong> there are a few things. Free chat only lasts five days, and then you have to stop wasting our server capacity (oh and also your match's time) - shit or get off the pot brah.<br />
          <br />
          If you propose a date, and the other person accepts, there's a small cost. We're toying with $.50, although that's subject to change.
        </Accordion.Body>
      </Accordion.Item>
            <Accordion.Item eventKey="6">
        <Accordion.Header>FIFTY DOLLARS FOR A DATE?!?!? THAT'S OUTRAGEOUS!!!!</Accordion.Header>
        <Accordion.Body>
          There's a decimal there. Fifty <em>cents.</em> As in half of a dollar.<br />
            <br />
          <strong>BUT</strong> there are a few things. Free chat only lasts five days, and then you have to stop wasting our server capacity (oh and also your match's time) - shit or get off the pot brah.<br />
          <br />
          If you propose a date, and the other person accepts, there's a small cost. We're toying with $.50, although that's subject to change.
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="7">
        <Accordion.Header>Oh. That's...not much.</Accordion.Header>
        <Accordion.Body>
          There's a decimal there. Fifty <em>cents.</em> As in half of a dollar.<br />
            <br />
          <strong>BUT</strong> there are a few things. Free chat only lasts five days, and then you have to stop wasting our server capacity (oh and also your match's time) - shit or get off the pot brah.<br />
          <br />
          If you propose a date, and the other person accepts, there's a small cost. We're toying with $.50, although that's subject to change.
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

export default FAQ;