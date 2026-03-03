"use client";
import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { getPostById, createComment, castVote } from "@/lib/api";
import { Post } from "@/types";
import VoteSidebar from "@/components/forum/VoteSidebar";
import CommentForm from "@/components/forum/CommentForm";
import Link from "next/link";

export default function PostDetailPage() {
  const { id } = useParams();
  const router = useRouter();
  const { user } = useAuth();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) fetchThread();
  }, [id]);

  const fetchThread = async () => {
    try {
      const data = await getPostById(id as string);
      setPost(data);
    } finally {
      setLoading(false);
    }
  };

  const handlePostVote = async (val: number) => {
    if (!user || !post) return alert("AUTH_REQUIRED");

    const current = post.user_vote ?? 0;
    const newValue = current === val ? 0 : val;
    const diff = newValue - current;

    setPost({
      ...post,
      user_vote: newValue,
      upvotes: post.upvotes + diff,
    });

    try {
      await castVote(post.id, "post", newValue);
    } catch (err) {
      fetchThread();
    }
  };

  const handleCommentVote = async (commentId: string, val: number) => {
    if (!user || !post) return;

    setPost({
      ...post,
      comments: post.comments.map((c) => {
        if (c.id !== commentId) return c;
        const current = c.user_vote ?? 0;
        const newValue = current === val ? 0 : val;
        return {
          ...c,
          user_vote: newValue,
          upvotes: c.upvotes + (newValue - current),
        };
      }),
    });

    try {
      await castVote(commentId, "comment", val);
    } catch (err) {
      fetchThread();
    }
  };

  if (loading)
    return (
      <div className="p-20 font-black italic animate-pulse min-h-screen uppercase text-center">
        INITIALIZING_THREAD...
      </div>
    );
  if (!post)
    return (
      <div className="p-20 font-black text-red-600 min-h-screen uppercase text-center">
        ERROR: THREAD_NOT_FOUND
      </div>
    );

  return (
    <div className="max-w-5xl mx-auto space-y-6 pb-32 px-4 md:px-0">
      <button
        onClick={() => router.back()}
        className="flex items-center gap-2 text-[10px] font-black uppercase opacity-40 hover:opacity-100 cursor-pointer transition-opacity"
      >
        <svg
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="4"
        >
          <path d="M19 12H5m7 7l-7-7 7-7" />
        </svg>
        Return_to_Feed
      </button>

      <div className="flex items-stretch border-4 border-black bg-white shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] md:shadow-[12px_12px_0px_0px_rgba(0,0,0,1)]">
        <VoteSidebar
          upvotes={post.upvotes}
          userVote={post.user_vote}
          onVote={handlePostVote}
        />

        <div className="flex-1 p-4 md:p-8 flex flex-col md:flex-row items-start gap-6 md:gap-8">
          <div className="flex-1 space-y-4 md:space-y-6 flex flex-col w-full">
            <div className="text-[10px] md:text-xs font-black uppercase opacity-40 flex items-center gap-2">
              <Link
                href={`/profile/${post.creator?.username}`}
                className="flex items-center gap-2 hover:opacity-100 transition-opacity group"
              >
                <div className="w-5 h-5 md:w-6 md:h-6 border-2 border-black rounded-full bg-black flex items-center justify-center shrink-0">
                  <span className="text-[8px] md:text-[10px] text-white font-black leading-none">
                    {post.creator?.username?.[0].toUpperCase() || "A"}
                  </span>
                </div>
                <span className="group-hover:underline truncate">
                  {post.creator?.username}
                </span>
              </Link>
              <span>•</span>
              <span className="truncate">
                {new Date(post.created_at).toLocaleString()}
              </span>
            </div>

            <h1 className="text-2xl md:text-4xl font-black uppercase tracking-tighter leading-tight md:leading-none">
              {post.title}
            </h1>
            <p className="text-sm md:text-lg opacity-80 whitespace-pre-wrap leading-tight">
              {post.content}
            </p>
          </div>

          {post.maze?.thumbnail && (
            <Link
              href={`/solve?id=${post.maze_id}`}
              className="w-full md:w-64 aspect-square border-4 border-black relative overflow-hidden group shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] md:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] hover:shadow-none transition-all hover:translate-x-1 hover:translate-y-1 shrink-0"
            >
              <img
                src={post.maze.thumbnail}
                alt="Target"
                className="w-full h-full object-cover grayscale-0 md:grayscale md:opacity-80 group-hover:grayscale-0 group-hover:scale-110 transition-all duration-500"
              />
            </Link>
          )}
        </div>
      </div>

      <div className="pt-8 space-y-6">
        <h3 className="text-lg md:text-xl font-black uppercase italic tracking-tighter">
          Comment_Section ({post.comments?.length || 0})
        </h3>

        {user ? (
          <CommentForm
            onSubmit={async (c) => {
              await createComment(post.id, c);
              fetchThread();
            }}
          />
        ) : (
          <div className="border-4 border-dashed border-black p-6 text-center opacity-40 font-black uppercase text-sm">
            Authentication_Required
          </div>
        )}

        <div className="space-y-4">
          {post.comments?.map((comment) => (
            <div
              key={comment.id}
              className="border-4 border-black bg-white flex items-stretch overflow-hidden shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]"
            >
              <VoteSidebar
                small
                upvotes={comment.upvotes}
                userVote={comment.user_vote}
                onVote={(v) => handleCommentVote(comment.id, v)}
              />
              <div className="flex-1 p-3 md:p-4 space-y-2">
                <div className="flex justify-between items-center text-[9px] md:text-[10px] font-black uppercase opacity-40">
                  <Link
                    href={`/profile/${comment.creator?.username}`}
                    className="flex items-center gap-1.5 hover:opacity-100 transition-opacity group"
                  >
                    <div className="w-4 h-4 border border-black rounded-full bg-black flex items-center justify-center shrink-0">
                      <span className="text-[7px] text-white font-black leading-none">
                        {comment.creator?.username?.[0].toUpperCase() || "A"}
                      </span>
                    </div>
                    <span className="group-hover:underline">
                      {comment.creator?.username}
                    </span>
                  </Link>
                  <span className="hidden sm:inline">
                    {new Date(comment.created_at).toLocaleDateString()}
                  </span>
                </div>

                <p className="text-sm font-medium leading-relaxed">
                  {comment.content}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
