import type { Metadata } from "next";
import { MilkeyWay } from "~/components/MilkyWay";
import { opensans } from "../styles/fonts";
import "../styles/globals.css";
import { Providers } from "./providers";

export const metadata: Metadata = {
  title: "Missing studio",
  description:
    "Missing studio is a robust open source platform to build a developer-first way to create AI application",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}): JSX.Element {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={opensans.className}>
        <MilkeyWay />
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
