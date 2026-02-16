"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export default function ProfilePage() {
  const { username } = useParams();
  const [profile, setProfile] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${BASE_URL}/profile?username=${username}`)
      .then((res) => res.json())
      .then((data) => {
        setProfile(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [username]);

  if (loading)
    return (
      <div className="p-20 font-black italic animate-pulse uppercase">
        {" "}
        INITIALIZING_USER_PROFILE...
      </div>
    );
  if (!profile)
    return (
      <div className="p-20 font-black text-red-600 uppercase">
        {" "}
        ERROR: USER_NOT_FOUND{" "}
      </div>
    );

  const totalDeadEnds = profile.mazes.reduce(
    (acc: number, m: any) => acc + (m.dead_ends || 0),
    0
  );
  const avgComplexity =
    profile.stats.total_mazes > 0
      ? profile.mazes.reduce(
          (acc: number, m: any) => acc + (m.complexity || 0),
          0
        ) / profile.stats.total_mazes
      : 0;

  return (
    <div className="max-w-6xl mx-auto py-12 px-6 space-y-8">
      {/* Header Section */}
      <section className="relative border-4 border-black p-8 bg-white shadow-[12px_12px_0px_0px_rgba(0,0,0,1)] flex flex-col md:flex-row items-center gap-8">
        <div className="w-32 h-32 bg-black text-white border-4 border-black flex items-center justify-center text-5xl font-black shrink-0 shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)]">
          {profile.username[0].toUpperCase()}
        </div>

        <div className="flex-1 text-center md:text-left">
          <h1 className="text-5xl font-black uppercase tracking-tighter leading-none mb-2">
            {profile.username}
          </h1>
          <div className="flex flex-wrap justify-center md:justify-start gap-4">
            <span className="bg-zinc-100 px-3 py-1 text-[10px] font-black uppercase border-2 border-black">
              ID: {profile.username.toLowerCase()}
            </span>
            <span className="text-sm font-bold opacity-40 uppercase py-1">
              JOINED // {new Date(profile.created_at).toLocaleDateString()}
            </span>
          </div>
        </div>

        {/* Global Stats Grid */}
        <div className="grid grid-cols-2 gap-4 w-full md:w-auto border-t-4 md:border-t-0 md:border-l-4 border-black pt-8 md:pt-0 md:pl-8">
          <div className="text-center md:text-left">
            <span className="block text-3xl font-black leading-none">
              {profile.stats.total_mazes}
            </span>
            <span className="text-[9px] font-black uppercase opacity-50 tracking-tighter">
              Mazes_Built
            </span>
          </div>
          <div className="text-center md:text-left">
            <span className="block text-3xl font-black leading-none">
              {totalDeadEnds}
            </span>
            <span className="text-[9px] font-black uppercase opacity-50 tracking-tighter">
              Total_Dead_Ends
            </span>
          </div>
          <div className="text-center md:text-left">
            <span className="block text-3xl font-black leading-none">
              {avgComplexity.toFixed(2)}
            </span>
            <span className="text-[9px] font-black uppercase opacity-50 tracking-tighter">
              Avg_Complexity
            </span>
          </div>
        </div>
      </section>

      {/* Mazes Gallery */}
      <div className="space-y-4">
        <div className="flex items-center gap-4">
          <h2 className="text-xl font-black uppercase italic tracking-tight">
            Public_Repositories
          </h2>
          <div className="flex-1 h-1 bg-black opacity-10"></div>
        </div>

        {profile.mazes.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {profile.mazes.map((maze: any) => (
              <Link key={maze.id} href={`/solve?id=${maze.id}`}>
                <div className="border-4 border-black p-4 bg-white hover:-translate-y-1 hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] transition-all group flex flex-col h-full">
                  <div
                    className="aspect-video bg-zinc-50 border-2 border-black mb-4 flex items-center justify-center overflow-hidden relative"
                    style={{
                      backgroundImage: `url(${maze.thumbnail})`,
                      backgroundSize: "cover",
                      backgroundPosition: "center",
                    }}
                  >
                    <div className="absolute inset-0 bg-white/80 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center backdrop-blur-[2px]">
                      <span className="text-[10px] font-black uppercase tracking-widest z-10">
                        VIEW_MATRIX_{maze.id}
                      </span>
                    </div>
                  </div>

                  <div className="mt-auto space-y-3">
                    <div className="flex justify-between items-center">
                      <span className="font-black text-lg italic tracking-tighter">
                        {maze.rows}x{maze.cols}
                      </span>
                      <span className="bg-black text-white text-[9px] font-black px-2 py-0.5">
                        COMPLEXITY: {maze.complexity?.toFixed(2) ?? "0.00"}
                      </span>
                    </div>
                    <div className="text-[9px] font-bold opacity-30 flex justify-between border-t border-black border-dotted pt-2">
                      <span>NODES: {maze.rows * maze.cols}</span>
                      <span>
                        {new Date(maze.created_at).toLocaleDateString()}
                      </span>
                    </div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        ) : (
          <div className="border-4 border-black border-dashed p-12 text-center">
            <p className="font-black opacity-20 uppercase tracking-[0.3em]">
              No_Mazes_Generated_Yet
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
