"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import ProfileHeader from "@/components/profile/ProfileHeader";
import MazeRepositoryCard from "@/components/profile/MazeRepositoryCard";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export default function ProfilePage() {
  const { username } = useParams();
  const { user: currentUser } = useAuth();
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

  const handleDelete = async (e: React.MouseEvent, mazeId: string) => {
    e.preventDefault();
    if (!confirm(`DANGER: Purge Matrix ${mazeId} from database?`)) return;

    try {
      const res = await fetch(`${BASE_URL}/maze/delete?id=${mazeId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("gridgo_token")}`,
        },
      });

      if (res.ok) {
        setProfile({
          ...profile,
          mazes: profile.mazes.filter((m: any) => m.id !== mazeId),
          stats: {
            ...profile.stats,
            total_mazes: profile.stats.total_mazes - 1,
          },
        });
      }
    } catch (err) {
      console.error("PURGE_SEQUENCE_FAILED", err);
    }
  };

  if (loading)
    return (
      <div className="p-20 font-black italic animate-pulse uppercase">
        {" "}
        INITIALIZING_PROFILE...
      </div>
    );
  if (!profile)
    return (
      <div className="p-20 font-black text-red-600 uppercase">
        {" "}
        ERROR: USER_NOT_FOUND
      </div>
    );

  // Global Derived Stats
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
    <div className="max-w-6xl mx-auto py-12 px-6 space-y-12">
      <ProfileHeader
        username={profile.username}
        joinedAt={profile.created_at}
      />

      {/* Global Statistics Section */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[
          { label: "Mazes_Built", val: profile.stats.total_mazes },
          { label: "Total_Dead_Ends", val: totalDeadEnds },
          { label: "Avg_Complexity", val: avgComplexity.toFixed(2) },
        ].map((s, i) => (
          <div
            key={i}
            className="border-4 border-black p-6 bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]"
          >
            <span className="block text-4xl font-black italic leading-none">
              {s.val}
            </span>
            <span className="text-[10px] font-black uppercase opacity-40 tracking-widest">
              {s.label}
            </span>
          </div>
        ))}
      </div>

      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <h2 className="text-xl font-black uppercase italic tracking-tighter">
            Mazes
          </h2>
          <div className="flex-1 h-1 bg-black opacity-10"></div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {profile.mazes.map((maze: any) => (
            <MazeRepositoryCard
              key={maze.id}
              maze={maze}
              isOwner={currentUser === profile.username}
              onDelete={handleDelete}
            />
          ))}
        </div>
      </div>
    </div>
  );
}
