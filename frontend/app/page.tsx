import Header_main from "@/components/layout/Header_main";
import Footer_main from "@/components/layout/Footer-main";

import "./homepage.css";

export const metadata = {
  title: "Homepage",
};


export default function Home() {
  return (
    <>
      <Header_main />
      <div className="headline">
        <div id="headline-text">
          <h1>Help Sustain Wear in achieving its mission</h1>
          <p>Together, we can close the loop on fashion waste. Join our movement to give pre-loved clothes a second life and protect our planet for future generations.</p>
          <div className="headline-buttons">
            <a href=".." id="button">Donate now</a>
            <a href="./MissionImpact" id="button">Find more about our mission</a>
          </div>
        </div>
        <img src="/images/headline_image.webp" alt="Sustainable Fashion" />
      </div>

      <div id="Content">
        <div id="Content-img">
          <img src="/images/headline_image.webp" alt="Sustainable Fashion" />
        </div>
        <div id="Content-text">
          <h2>Why Choose Sustain Wear?</h2>
          <p>
            Measurable Impact:<br></br>
            Don't just donate,know the difference you make. Our smart dashboard calculates exactly how much CO2 and landfill waste you save with every item you give.
          </p>

          <p>
            Total Transparency<br></br>
            Wondering where your old coat ended up? We bridge the gap between donors and charities, giving you visibility into the lifecycle of your donation.
          </p>

          <p>
            Empowering Charities<br></br>
            We help local charities sort and distribute inventory faster. Your structured data helps them spend less time sorting piles and more time helping people.
          </p>
        </div>
      </div>

      <div id="splitContent">
        <div className="container">
          <div className="imgcontainer">
            <a href="../MissionImpact"><img src="images/headline_image.webp" alt="Mission & Impact" /></a>
          </div>
          <h3>Mission & Impact</h3>
        </div>

        <div className="container">
          <div className="imgcontainer">
            <a href="../FindUs"><img src="./images/headline_image.webp" alt="Find us" /></a>
          </div>
          <h3>Find us</h3>
        </div>

        <div className="container">
          <div className="imgcontainer">
            <a href=".."><img src="./images/headline_image.webp" alt="My Account" /></a>
          </div>
          <h3>My Account</h3>
        </div>
      </div>

      <div className="partners">
        <h2>Trusted by UK organisations & communities</h2>
        <p>We partner with leading organisations and local shelters to ensure your donations reach those who need them most. Together, we are building a transparent, zero-waste future.</p>
        <div className="partner-logos">
          <img src="/images/headline_image.webp" alt="Partner 1" />
          <img src="/images/headline_image.webp" alt="Partner 1" />
          <img src="/images/headline_image.webp" alt="Partner 1" />
          <img src="/images/headline_image.webp" alt="Partner 1" />
          <img src="/images/headline_image.webp" alt="Partner 1" />
          <img src="/images/headline_image.webp" alt="Partner 1" />
          <img src="/images/headline_image.webp" alt="Partner 1" />
          <img src="/images/headline_image.webp" alt="Partner 1" />
          <img src="/images/headline_image.webp" alt="Partner 1" />
        </div>
      </div>
      <Footer_main />
    </>
  );
}
