import './../globals.scss'
import Navigation from "./nav";
import { AuthProvider } from "./auth/authContext";

export default function PortalLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
            <body>
                <AuthProvider>
                    <Navigation>
                        {children}
                    </Navigation>
                </AuthProvider>
            </body>
        </html>
    );
}