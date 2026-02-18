"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { castVote } from "@/lib/api";
import { Post } from "@/types";
import Link from "next/link";
import VoteSidebar from "./VoteSidebar";

export default function ForumCard({ post }: { post: Post }) {
  const { user } = useAuth();
  const router = useRouter();
  const [currentVote, setCurrentVote] = useState(post.user_vote ?? 0);
  const [voteCount, setVoteCount] = useState(post.upvotes);

  useEffect(() => {
    setCurrentVote(post.user_vote ?? 0);
    setVoteCount(post.upvotes);
  }, [post.user_vote, post.upvotes]);

  const handleVote = async (val: number) => {
    if (!user) return alert("AUTH_REQUIRED");
    const current = currentVote;
    const newValue = current === val ? 0 : val;
    const diff = newValue - current;

    setVoteCount((prev) => prev + diff);
    setCurrentVote(newValue);

    const success = await castVote(post.id, "post", newValue);
    if (!success) {
      setVoteCount(post.upvotes);
      setCurrentVote(post.user_vote ?? 0);
    }
  };

  return (
    <div
      onClick={() => router.push(`/forum/post/${post.id}`)}
      className="flex h-48 border-4 border-black bg-white shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all group cursor-pointer overflow-hidden"
    >
      <div className="flex" onClick={(e) => e.stopPropagation()}>
        <VoteSidebar
          upvotes={voteCount}
          userVote={currentVote}
          onVote={handleVote}
          small
        />
      </div>

      <div className="flex-1 flex items-center p-5 gap-6 overflow-hidden">
        <div className="flex-1 flex flex-col justify-between h-full overflow-hidden">
          <div className="space-y-2">
            <div className="flex items-center gap-2 text-[9px] font-black uppercase opacity-40">
              <span>u/{post.creator?.username || "anon"}</span>
              <span>•</span>
              <span>{new Date(post.created_at).toLocaleDateString()}</span>
            </div>
            <h2 className="text-xl font-black uppercase tracking-tighter group-hover:underline decoration-4 truncate">
              {post.title}
            </h2>
            <p className="text-sm opacity-70 line-clamp-3 leading-tight select-none">
              {post.content}
            </p>
          </div>

          <div className="flex items-center gap-4 shrink-0 text-[10px] font-black uppercase">
            {post.maze_id && (
              <Link
                href={`/solve?id=${post.maze_id}`}
                onClick={(e) => e.stopPropagation()}
                className="hover:underline flex items-center gap-1"
              >
                <span className="opacity-30">REF:</span>{" "}
                {post.maze_id.slice(0, 8)}
              </Link>
            )}
            <span className="opacity-30">|</span>
            <span className="opacity-40">
              {post.comments?.length || 0} Comments
            </span>
          </div>
        </div>

        {post.maze?.thumbnail && (
          <div className="hidden md:block h-32 w-32 border-4 border-black relative shrink-0 overflow-hidden bg-zinc-100 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]">
            <img
              src={post.maze.thumbnail}
              alt="Preview"
              className="w-full h-full object-cover opacity-90 group-hover:opacity-100 group-hover:scale-110 transition-all duration-500"
            />
            <div className="absolute inset-0 bg-black opacity-0 group-hover:opacity-10 transition-opacity"></div>
          </div>
        )}
      </div>
    </div>
  );
}
