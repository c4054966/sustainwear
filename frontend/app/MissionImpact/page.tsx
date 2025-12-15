import Header_main from "@/components/layout/Header_main";
import Footer_main from "@/components/layout/Footer-main";
import "./main.css";

export const metadata = {
    title: "Mission & Impact",
};

export default function MissionImpactPage() {
    return (
        <>
            <Header_main />

            <div className="headline">
                <img src="/images/headline_image.webp" alt="Sustainable Fashion" />
                <div id="headline-text">
                    <h1>Reweaving the Fabric of Society</h1>
                    <p>We are on a mission to end textile waste by connecting generous communities with the charities that need them most.</p>
                </div>
            </div>

            <div id="Content-white-bg">
                <div id="Content-white-text">
                    <h2>The Fast Fashion Crisis</h2>
                    <p>
                        Every year, the UK throws away over 300,000 tonnes of clothing into landfill. That's enough to fill Wembley Stadium. Most of these items are perfectly wearable, but a fragmented donation system means they often don't reach the right people in time.
                    </p>
                    <ul>
                        <li>10% of global carbon emissions come from the fashion industry.</li>
                        <li>    85% of all textiles go to the dump each year.</li>
                        <li>One cotton shirt takes 2,700 litres of water to make, enough for one person to drink for 2.5 years.</li>
                    </ul>
                </div>
                <div id="Content-white-img">
                    <img src="/images/headline_image.webp" alt="Sustainable Fashion" />
                </div>
            </div>

            <div id="Content-grey-bg">
                <div id="Content-grey-img">
                    <img src="/images/headline_image.webp" alt="Sustainable Fashion" />
                </div>
                <div id="Content-grey-text">
                    <h2>How Sustain Wear Closes the Loop</h2>
                    <p>
                        Sustain Wear isn't just a donation form; it's an intelligent inventory engine. By digitizing donations before they leave your home, we ensure:
                    </p>
                    <ul>
                        <li><b>Smart Sorting:</b> AI-assisted categorization means charities know exactly what they are getting before it arrives.</li>
                        <li><b>Zero Waste:</b> Items are matched to specific needs—coats for winter shelters, suits for job interview schemes.</li>
                        <li><b>Carbon Transparency:</b> We calculate the exact carbon footprint saved by extending the life of your garment.</li>
                    </ul>
                    <div className="buttons">
                        <a href=".." id="button">Start your journey today</a>
                    </div>
                </div>
            </div>

            <div id="Content-white-bg">
                <div id="Content-white-text">
                    <h2>Measuring what matters today</h2>
                    <p>
                        Transparency is at the heart of our platform. We track three key metrics for every item donated:
                    </p>
                    <ul>
                        <li><b>CO2 Avoided:</b> The emissions prevented by not manufacturing a new replacement item.</li>
                        <li><b>Landfill Diverted:</b>The physical weight of waste kept out of the ground.</li>
                        <li><b>Lives Touched:</b>The number of families or individuals who receive essential clothing.</li>
                    </ul>
                    <p>
                        We don't just guess our impact, we calculate it using industry-standard data and provide donors with a detailed breakdown of their contribution.
                    </p>
                </div>
                <div id="Content-white-img">
                    <img src="/images/headline_image.webp" alt="Sustainable Fashion" />
                </div>
            </div>

            <Footer_main />
        </>
    );
}