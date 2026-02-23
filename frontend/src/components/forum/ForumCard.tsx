"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { castVote } from "@/lib/api";
import { Post } from "@/types";
import Link from "next/link";
import VoteSidebar from "./VoteSidebar";

interface ForumCardProps {
  post: Post;
  isOwner?: boolean;
  onDelete?: (id: string) => void;
  onVoteUpdate?: (postId: string, newVote: number, newCount: number) => void;
}

export default function ForumCard({
  post,
  isOwner,
  onDelete,
  onVoteUpdate,
}: ForumCardProps) {
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
    const updatedCount = voteCount + diff;

    setVoteCount(updatedCount);
    setCurrentVote(newValue);

    if (onVoteUpdate) {
      onVoteUpdate(post.id, newValue, updatedCount);
    }

    try {
      const success = await castVote(post.id, "post", newValue);
      if (!success) throw new Error();
    } catch (err) {
      console.error("Vote failed, rolling back...");
      setVoteCount(post.upvotes);
      setCurrentVote(post.user_vote ?? 0);
      if (onVoteUpdate) {
        onVoteUpdate(post.id, post.user_vote ?? 0, post.upvotes);
      }
    }
  };

  return (
    <div className="relative group transition-all duration-200 hover:-translate-y-1">
      <div
        onClick={() => router.push(`/forum/post/${post.id}`)}
        className="flex h-48 border-4 border-black bg-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] group-hover:shadow-[10px_10px_0px_0px_rgba(0,0,0,1)] transition-all cursor-pointer overflow-hidden relative"
      >
        {/* VOTE SIDEBAR */}
        <div className="flex" onClick={(e) => e.stopPropagation()}>
          <VoteSidebar
            upvotes={voteCount}
            userVote={currentVote}
            onVote={handleVote}
            small
          />
        </div>

        {/* MAIN CONTENT */}
        <div className="flex-1 flex items-center p-5 gap-6 overflow-hidden">
          <div className="flex-1 flex flex-col justify-between h-full overflow-hidden">
            <div className="space-y-2">
              <div className="flex items-center gap-2 text-[9px] font-black uppercase opacity-40">
                <Link
                  href={`/profile/${post.creator?.username}`}
                  onClick={(e) => e.stopPropagation()}
                  className="flex items-center gap-1.5 hover:opacity-100 transition-opacity group/user"
                >
                  <div className="w-4 h-4 border border-black rounded-full bg-black flex items-center justify-center shrink-0 overflow-hidden">
                    <span className="text-[7px] text-white font-black leading-none">
                      {post.creator?.username
                        ? post.creator.username[0].toUpperCase()
                        : "A"}
                    </span>
                  </div>
                  <span className="group-hover/user:underline">
                    {post.creator?.username || "anon"}
                  </span>
                </Link>
                <span>•</span>
                <span>{new Date(post.created_at).toLocaleDateString()}</span>
              </div>
              <h2 className="text-xl font-black uppercase tracking-tighter truncate">
                {post.title}
              </h2>
              <p className="text-sm opacity-70 line-clamp-3 leading-tight">
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
            </div>
          )}
        </div>

        {/* DELETE BUTTON - Nested inside the moving container */}
        {isOwner && onDelete && (
          <button
            onClick={(e) => {
              e.stopPropagation();
              onDelete(post.id);
            }}
            className="absolute top-2 right-2 bg-black text-white p-1.5 border-2 border-black opacity-0 group-hover:opacity-100 transition-all hover:bg-red-600 z-30 cursor-pointer shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:shadow-none"
          >
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="4"
            >
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        )}
      </div>
    </div>
  );
}
