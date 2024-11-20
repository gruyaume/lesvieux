"use client"

import Logo from "../components/logo"

export default function Footer() {
    return (
        <footer className="l-footer--sticky p-strip--highlighted">
            <nav className="row" aria-label="Footer">
                <div className="has-cookie">
                    <Logo />
                    <p>© 2024 This blog is built with <a href="https://github.com/gruyaume/lesvieux">LesVieux</a>.</p>
                </div>
            </nav>
        </footer>
    );
}
