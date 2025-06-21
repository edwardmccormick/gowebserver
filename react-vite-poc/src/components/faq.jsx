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
        <Accordion.Header className='bg-success-subtle text-white'>So, seriously, everything is free?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Yup, you heard right, chief. Most of Urmid is genuinely free. It's almost easier to tell you what's <em>not free:</em></p>
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
            <p>That's fifty cents. As in, half a dollar. Not enough for a coffee, more like enough for an extra shot of syrup in ours. So yeah, for fifty cents, you get a chance to meet someone cool, chase true love (or whatever else you're into), and maybe even take it to Poundtown™. Pretty good deal, if you ask us.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="3">
        <Accordion.Header>Huh. What's the catch?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Nothing. We're really that stupid.</p>
            <p>If you know an angel investor, put them in touch, because at some point these bags full of money we're lighting on fire are going to start running low, and shit's going to get Lord the Flies pretty fucking quick.</p>
            
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="4">
        <Accordion.Header>So once you get to a million users or whatever, you're going to start charging us $50 a month or something, right?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Honestly, if more than a handful of people read this FAQ, we've already exceeded expectations. But no, that's not the plan.</p>
            <p>Our founder is a certified nerd who built Urmid to be lean, fast, and unbelievably cheap to run. While we're designed to scale, we know that hitting tens of thousands of users will eventually mean some real server costs.</p>
            <p>When that happens, we'll likely become ad-supported. Yeah, ads can be annoying, but you're probably already ignoring them, just like you might be ignoring this very sentence. Our hope is that tasteful advertising will cover most of our operational expenses.</p>
            <p><strong>Most.</strong></p>
            <p>But seriously: the plan is to stay free, forever. Our goal isn't to become another multi-billion-dollar dating behemoth (hi Match Group! Big fan of your investor relation reports, kthxbye!). We just want to build an app that actually helps people find happiness, and maybe – just maybe – our founder can finally quit his day job someday.</p>
          </div>
        </Accordion.Body>
        </Accordion.Item>
      <Accordion.Item eventKey="5">
        <Accordion.Header>So there's really no catch? Are you, like, selling our DNA to the Russian mob? Is this a crypto scam?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>No but the Russian mob thing sounds interesting - put their people in touch with our people, would you?</p>
            <p>Urmid is built on a promise: Free forever. You pay a tiny fee for Date Requests when you find someone awesome. You pay a small 'Success Fee' when you leave us because you found somebody even more awesome. That's the deal. No fine print.</p>
            <p>We don't sell your data, we don't sell your DNA, and honestly, if we didn't legally have to keep an email address on file, we probably wouldn't even ask for that.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
            <Accordion.Item eventKey="6">
        <Accordion.Header>What if I don't want to pay 50 cents for a Date Request? Are you going to stop me from exchanging info in chat?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>We're not.</p>
            <p>Seriously, if 50 cents is going to break your dating budget, email us at ceo@urmid.com. We'll spot you a lifetime supply of internet cool points.</p>
            <p>Now, to keep things safe and not totally chaotic, we do have some clever monitoring in our Vibe Chats. It flags obvious (and not-so-obvious) attempts to exchange contact info like emails, phone numbers, or external websites. Not because we don't want you to connect – we absolutely do! But because spammers, trolls, and scammers use those exact same tactics. It's about protecting you, not preventing your love life.</p>
            <p>So, if you want to encode your Signal username by arranging pages from War and Peace into hexadecimal values within the chat... hey man, go off. Just know you'd be impressing each other by essentially swindling your future favorite tech billionaire out of 50 cents worth of server fuel. Your call.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="7">
        <Accordion.Header>Wait. You're okay with that?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>Look, we're a sticklers for rules, especially since we're the one who wrote 'em. But I'm not your mom, or your parole officer. (Though I'm almost as fat as your mom, and definitely way sluttier than your parole officer.)</p>
            <p>It's the same deal with closing your account. You can tell us you're closing it for any reason – you're tired of our shenanigans (fair!), you've embraced a life of celibacy, you're moving to Mars. You can even tell us Urmid sucks, go get married, and yup, keep your five bucks. Spend it on your honeymoon, or on hookers, or on hookers on your honeymoon. We won't judge.</p>
            <p>We're more likely to start publicly begging for donations (shoutout to Wikipedia! We love... everything about you, actually) than to ever try to screw you over. We're about genuine connections, not nickel-and-diming you.</p>
            <p>Free. Forever. That's the promise.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="8">
        <Accordion.Header>How do you keep answering questions and I feel like I have more questions than before? HOW? WHY? WHAT, I guess?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>The 'how' is simple: our founder is an inveterate nerd who basically built Urmid on a whim. He stumbled upon some Reddit rants about how agonizing it is to find someone decent these days, and thought, "I can do better." He's a big guy (in both nerdiness and physical stature) who specializes in building incredibly snazzy, performant, and absurdly cheap cloud infrastructure. Like, borderline free.</p>
            <p>The 'why' is a bit more philosophical. The world needs more happiness, more genuine human connection, and yes, even more hooking up. Life's often a dumpster fire, but if Urmid can spark a little joy in people's lives and help them connect... that'd be pretty damn cool.</p>
            <p>And for those angel investors out there (you might want to cover your ears for this part): screw capitalism. Why should finding someone to hang out with, or even to get laid, suddenly come with a hefty price tag? It's ridiculous.</p>
            <p>So, yeah. This is genuinely an 'audacity of hope' kind of venture. Which, if you think about it, is a pretty accurate description of dating itself.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="9">
        <Accordion.Header>Okay, but nerdy is my love language. How?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>React front end, statically hosted by a CDN.</p>
            <p>Go backend, running on function-as-a-service providers who I will name drop in a fucking millisecond for some server credits (call me!), all routed through a variety of Restful and Websocket connections.</p>
            <p>Database architecture running on free/cheap managed services. Props there - if there's somebody else offering SQL on a free tier, I sure couldn't find them, but Oracle has a nice offering (and I will sing your praises more for some server credits, boys!). There are a couple of databases, though - the other is a serverless document database with a nice free tier offering.</p>
            <p>If that's interesting, or sensical to you, and you'd like to join us....well, I wouldn't hold my breath, I can't even hire <em>my own damn self</em> yet. But what the hey: ceo@urmid.com; send your resume and....whatever else it is that people do to get noticed in tech nowadays. Nudes? Some sort of AI api key? Sure. Fire away.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
      <Accordion.Item eventKey="10">
        <Accordion.Header>Okay, but nerdy is my love language. How?</Accordion.Header>
        <Accordion.Body>
          <div className="text-start">
            <p>React front end, statically hosted by a CDN.</p>
            <p>Go backend, running on function-as-a-service providers who I will name drop in a fucking millisecond for some server credits (call me!), all routed through a variety of Restful and Websocket connections.</p>
            <p>Database architecture running on free/cheap managed services. Props there - if there's somebody else offering SQL on a free tier, I sure couldn't find them, but Oracle has a nice offering (and I will sing your praises more for some server credits, boys!). There are a couple of databases, though - the other is a serverless document database with a nice free tier offering.</p>
            <p>If that's interesting, or sensical to you, and you'd like to join us....well, I wouldn't hold my breath, remember the part about not being able to hire <em>my damn self</em> yet?</p>
            <p>But whatever. ceo@urmid.com; send your resume and....whatever else it is that people do to get noticed in tech nowadays. Nudes? Some sort of AI api key? Sure. Fire away.</p>
          </div>
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
  );
}

export default FAQ;