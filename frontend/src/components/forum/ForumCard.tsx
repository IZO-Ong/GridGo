"use client";
import { useState, useEffect } from "react";
import Link from "next/link";
import { useAuth } from "@/context/AuthContext";
import { castVote } from "@/lib/api";

export default function ForumCard({ post }: { post: any }) {
  const { user } = useAuth();

  // Initialize state from backend data to ensure persistence
  const [currentVote, setCurrentVote] = useState<number>(post.user_vote || 0);
  const [voteCount, setVoteCount] = useState<number>(post.upvotes || 0);

  // Sync state if post prop changes (important for tab switching/navigation)
  useEffect(() => {
    setCurrentVote(post.user_vote || 0);
    setVoteCount(post.upvotes || 0);
  }, [post.user_vote, post.upvotes]);

  const handleVote = async (val: number) => {
    if (!user) return alert("AUTH_REQUIRED_TO_VOTE");
    const newValue = currentVote === val ? 0 : val;
    const diff = newValue - currentVote;

    setVoteCount((prev) => prev + diff);
    setCurrentVote(newValue);

    try {
      const success = await castVote(post.id, "post", val);
      if (!success) {
        setVoteCount((prev) => prev - diff);
        setCurrentVote(currentVote);
      }
    } catch (err) {
      setVoteCount((prev) => prev - diff);
      setCurrentVote(currentVote);
    }
  };

  return (
    <div className="flex h-48 border-4 border-black bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all group">
      {/* Sidebar: Changed to Black themed voting */}
      <div className="w-12 bg-zinc-50 border-r-4 border-black flex flex-col items-center py-4 gap-1 shrink-0">
        <button
          onClick={() => handleVote(1)}
          className={`transition-all cursor-pointer p-1 rounded ${
            currentVote === 1
              ? "text-black scale-110"
              : "text-black opacity-20 hover:opacity-100"
          }`}
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill={currentVote === 1 ? "black" : "none"} // Fill when active
            stroke="black"
            strokeWidth="4"
          >
            <path d="M18 15l-6-6-6 6" />
          </svg>
        </button>

        <span className="font-black text-sm italic text-black">
          {voteCount}
        </span>

        <button
          onClick={() => handleVote(-1)}
          className={`transition-all cursor-pointer p-1 rounded ${
            currentVote === -1
              ? "text-black scale-110"
              : "text-black opacity-20 hover:opacity-100"
          }`}
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill={currentVote === -1 ? "black" : "none"} // Fill when active
            stroke="black"
            strokeWidth="4"
          >
            <path d="M6 9l6 6 6-6" />
          </svg>
        </button>
      </div>

      {/* Content Area */}
      <div className="flex-1 p-4 flex flex-col justify-between overflow-hidden">
        <div className="space-y-2">
          <div className="flex items-center gap-2 text-[9px] font-black uppercase opacity-40">
            <span>{post.creator?.username || "anonymous"}</span>
            <span>•</span>
            <span>{new Date(post.created_at).toLocaleDateString()}</span>
          </div>
          <Link href={`/forum/post/${post.id}`}>
            <h2 className="text-xl font-black uppercase tracking-tighter group-hover:underline decoration-4 truncate">
              {post.title}
            </h2>
          </Link>
          <p className="text-sm opacity-70 line-clamp-4">{post.content}</p>
        </div>
      </div>

      {post.maze && post.maze.thumbnail && (
        <Link
          href={`/solve?id=${post.maze_id}`}
          className="hidden md:block w-48 border-l-4 border-black relative shrink-0 overflow-hidden bg-zinc-100"
        >
          <img
            src={post.maze.thumbnail}
            alt="Maze Preview"
            className="w-full h-full object-cover opacity-90 group-hover:opacity-100 group-hover:scale-105 transition-all duration-300"
          />
        </Link>
      )}
    </div>
  );
}
