"use client";
import { useState } from "react";
import "./Header_homepage.css";

export default function Header() {
  const [menuOpen, setMenuOpen] = useState(false);

  const toggleMenu = () => {
    setMenuOpen(!menuOpen);
  };

  return (
    <header>
      <div id="topbar">
        <a href=".."><img src="/Logo.webp" /></a>
        <div id="account">
          <img src="/icons/Account_icon.webp" />
          <a href="/login">My Account</a>        </div>
        <img src="/icons/menu.webp" id='menubtn' onClick={toggleMenu}></img>
      </div>
      <nav className={menuOpen ? "open" : " "}>
        <ul>
          <li><a href="..">HOME</a></li>
          <li><a href="../MissionImpact">MISSION & IMPACT</a></li>
          <li><a href="../FindUs">FIND US</a></li>
          <li id="moblink"><a href="..">MY ACCOUNT</a></li>
        </ul>
      </nav>
    </header>
  );
}