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
  title: "Gidoc — Open-source tools for documents and PDFs",
  description:
    "Explore open-source tools for working with PDF documents and files.",
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
        {process.env.NODE_ENV === "development" && (
          <Agentation appName="Gidoc" />
        )}
      </body>
    </html>
  );
}
