import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/contexts/AuthContext";
import { ExperienceProvider } from "@/contexts/ExperienceContext";
import { SoundProvider } from "@/components/providers/SoundProvider";
import { AchievementProvider } from "@/components/providers/AchievementProvider";
import { OnboardingGuard } from "@/components/onboarding/OnboardingGuard";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "CyclingStream - Pro Cycling Streaming Platform",
  description: "Watch live professional cycling races and events",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased flex flex-col min-h-screen bg-background text-foreground`}
      >
        <AuthProvider>
          <ExperienceProvider>
            <SoundProvider>
              <AchievementProvider>
                <OnboardingGuard>
                  {children}
                </OnboardingGuard>
              </AchievementProvider>
            </SoundProvider>
          </ExperienceProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
