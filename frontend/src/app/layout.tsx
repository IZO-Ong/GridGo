"use client";
import NavBar from "@/components/layout/NavBar";
import MazeMargin from "@/components/layout/MazeMargin";
import { AuthProvider } from "@/context/AuthContext";
import Image from "next/image";
import "./globals.css";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-white text-black relative font-mono overflow-x-hidden">
        <AuthProvider>
          <div className="hidden lg:block">
            <MazeMargin side="left" flip />
            <MazeMargin side="right" />
          </div>

          <div
            id="main-content-wrapper"
            className="relative z-10 max-w-5xl mx-auto flex flex-col"
          >
            <div className="p-4 md:p-8 flex flex-col">
              <header className="pt-4 pb-2 flex flex-col md:flex-row md:justify-between items-start md:items-end border-b-4 border-black bg-white gap-4 md:gap-0">
                <div className="flex items-center gap-4 cursor-default">
                  <h1 className="text-3xl md:text-4xl font-black uppercase tracking-tighter leading-none">
                    GRIDGO
                  </h1>

                  <div className="relative w-10 h-10 md:w-12 md:h-12">
                    <Image
                      src="/gopher-maze.png"
                      alt="GridGo Mascot"
                      fill
                      sizes="(max-width: 768px) 40px, 48px"
                      className="object-contain"
                      priority
                    />
                  </div>
                </div>

                <div className="w-full md:w-auto">
                  <NavBar />
                </div>
              </header>

              <main className="pt-8">{children}</main>
            </div>
          </div>
        </AuthProvider>
      </body>
    </html>
  );
}
