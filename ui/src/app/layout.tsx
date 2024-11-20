import type { Metadata } from "next";
import './globals.scss'
import Footer from "./components/footer"
import { ReactQueryClientProvider } from './components/ReactQueryClientProvider'

export const metadata: Metadata = {
  title: "LesVieux",
  description: "A blogging platform",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <ReactQueryClientProvider>
      <html lang="en">
        <head>
          <link rel="icon" href="/favicon.ico" sizes="any" />
        </head>
        <body>
          {children}
          <Footer />
        </body>
      </html>
    </ReactQueryClientProvider>
  );
}
