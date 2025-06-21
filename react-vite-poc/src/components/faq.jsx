import Accordion from 'react-bootstrap/Accordion';

function FAQ() {
  return (
    <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
        <Accordion.Header className='bg-success text-white'>Wtf even is this?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Man I don't know. It's a dating app. We thought most dating apps sucked, and surely there was a better way.</p>
            <p>Fucked up pretty bad, didn't we?</p>
            <br />
            <p>Okay, sorry, that was rude. I'm still kind of new at this.</p>
            <br />
            <p>urmid is a dating app that we made because, well, most dating apps suck.</p>
            <p>We wanted to make something that was simple, easy to use, and didn't have a bunch of bullshit.</p>
            <p>We wanted to make something that was aligned with your interests - we think it's bullshit that most dating apps make you 'pay to win.'</p>
            <p>So our model is a little different. Everything is free. You only pay if you find a date (to exchange information with them) and when you eventually decide to quit using our site because you found what you were looking.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="1">
        <Accordion.Header className='bg-success-subtle text-white'>So everything is free?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Yup, you got it chief.</p>
            <p>There are enough free things that it's easier to talk about what's <em>not</em>free.</p>
            <br />
            <p>If you've been chatting with someone, and you want to take it to the real world (which you should! the internet sucks, mkay?) you pay to ask them out.</p>
            <br />
            <p>We're still kind of pinning down the business model, and if you're reading this we're still in pre-alpha, so even that shit is free, but, we're leaning towards, like, $.50</p>
            <p>And when you find whatever you're looking for - eternal love and happiness, someone to tie you up and call you names, or just so much pussy/dick that even sex is boring, we charge you when you break up with us.</p>
            <p>Five dollars left on the nightstand of the internet, so that we can 'buy ourselves something nice.'</p>
            <p>That's it, man.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="2">
        <Accordion.Header>FIFTY DOLLARS FOR A DATE?!?!? THAT'S OUTRAGEOUS!!!!</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Woah, man, slow down. Take a deep breath.</p>
            <p>That's fifty cents. Like, half a dollar.</p>
            <p>Yup, you get someone to hang out with, a chance at true love and/or taking it to poundtown or whatever, and we get....enough money for like, an extra shot of syrup in our coffee.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>Huh. What's the catch?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Nothing. We're really that stupid.</p>
            <p>If you know an angel investor, put them in touch, because at some point these bags full of money we're lighting on fire are going to start running low, and shit's going to get wild. Quick.</p>
            
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="4">
        <Accordion.Header>So once you get to a million users or whatever, you're going to start charging us $50 a month or something, right?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Listen, if like, 30 people read this stupid line of text...well that's all we're (I'm) really out for.</p>
            <p>Once we hit a couple of thousand users, things could get a little spicy - the guy that built this is an inveterate nerd who built it to be small, quick, and stupid cheap. But at about 10,000 users or so, it's going to start costing him money to run.</p>
            <p>We'll probably go ad supported. Yeah, it sucks too, but...you'll ignore them, just like you ignore this FAQ, and keep chasing each other like horny, rabid dogs, and mostly it looks like the advertising will pay for the server costs.</p>
            <p><strong>Mostly.</strong></p>
            <p>But seriously. The plan is to stay free, forever. Avoid the dumb bullshit that's made every other app awful, and instead of being worth a couple of billion dollars (hi Match Group! Big fan of your investor relation reports, kthxbye!) our founder might get to quit his day job in a couple of years.</p>
          </div>
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