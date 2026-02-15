import NavBar from "@/components/layout/NavBar";
import MazeMargin from "@/components/layout/MazeMargin";
import "./globals.css";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-white text-black min-h-screen relative font-mono overflow-x-hidden">
        <MazeMargin side="left" />
        <MazeMargin side="right" />

        <div className="relative z-10 max-w-5xl mx-auto flex flex-col min-h-screen">
          {/* Wrapper with p-8 governs everything */}
          <div className="p-8 flex flex-col flex-1">
            <header className="pt-4 pb-2 flex justify-between items-end border-b-4 border-black bg-white">
              <h1 className="text-4xl font-black uppercase tracking-tighter">
                GRIDGO
              </h1>
              <NavBar />
            </header>

            {/* Remove p-8 from main */}
            <main className="flex-1 pt-8">{children}</main>
          </div>
        </div>
      </body>
    </html>
  );
}
