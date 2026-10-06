import type { Metadata } from "next";
import localFont from "next/font/local";
import { Agentation } from "agentation";
import "./globals.css";

const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
});
const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
});

export const metadata: Metadata = {
  title: "Gidocs",
  description: "Open-source tools for documents and PDFs",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable}`}>
        {children}
        {process.env.NODE_ENV === "development" && <Agentation appName="Gidocs" />}
      </body>
    </html>
  );
}
