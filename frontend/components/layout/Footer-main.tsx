import "./Footer_main.css";

export default function Header() {
    return (

        <footer>
                    <div id="footerTop">
                        <img src="/logo.webp"></img>
                        <div className="footerContent">
                            <h3>Sustain Wear</h3>
                            <ul>
                                <li><a href=".."><p>Home</p></a></li>
                                <li><a href="../MissionImpact"><p>Mission & Impact</p></a></li>
                                <li><a href="../FindUs"><p>Find Us</p></a></li>
                                <li><a href="../MyAccount"><p>My Account</p></a></li>
                            </ul>
                        </div>
                        <div className="footerContent">
                            <h3>Contact</h3>
                            <ul>
                                <li><p>cantorhelp@cc.ac.uk</p></li>
                                <li><p>internationalEXP@cc.ac.uk</p></li>
                                <li><p>+44-(0)144 234 0235</p></li>
                                <li><p>+44-(0)744 234 0236</p></li>
                            </ul>
                        </div>
                        <div className="footerContent">
                            <h3>Address</h3>
                            <ul>
                                <li><p>Cantor College</p></li>
                                <li><p>Main Street</p></li>
                                <li><p>Sheffield</p></li>
                                <li><p>SC4 2BB</p></li>
                            </ul>
                        </div>
                    </div>
                    <div id="footerBottom">
                        <p>Copright &copy; {new Date().getFullYear()} , ALL RIGHTS RESERVED to <b>Sustain Wear</b></p>
                    </div>
                </footer>

  );
}