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
    setPost({
      ...post,
      user_vote: newValue,
      upvotes: post.upvotes + (newValue - current),
    });
    await castVote(post.id, "post", newValue);
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
    await castVote(commentId, "comment", val);
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
    <div className="max-w-5xl mx-auto space-y-8 pb-32">
      <button
        onClick={() => router.back()}
        className="flex items-center gap-2 text-[10px] font-black uppercase opacity-40 hover:opacity-100 cursor-pointer"
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

      {/* Main Post Section: items-stretch ensures Sidebar h-full works */}
      <div className="flex items-stretch border-4 border-black bg-white shadow-[12px_12px_0px_0px_rgba(0,0,0,1)]">
        <VoteSidebar
          upvotes={post.upvotes}
          userVote={post.user_vote}
          onVote={handlePostVote}
        />

        <div className="flex-1 p-8 flex items-start gap-8">
          <div className="flex-1 space-y-6 flex flex-col">
            <div className="text-xs font-black uppercase opacity-40">
              u/{post.creator?.username} •{" "}
              {new Date(post.created_at).toLocaleString()}
            </div>
            <h1 className="text-4xl font-black uppercase tracking-tighter leading-none">
              {post.title}
            </h1>
            <p className="text-lg opacity-80 whitespace-pre-wrap">
              {post.content}
            </p>
          </div>
          {post.maze?.thumbnail && (
            <Link
              href={`/solve?id=${post.maze_id}`}
              className="hidden lg:block w-64 h-64 border-4 border-black relative overflow-hidden group shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]"
            >
              <img
                src={post.maze.thumbnail}
                alt="Target"
                className="w-full h-full object-cover grayscale opacity-80 group-hover:grayscale-0 group-hover:scale-110 transition-all duration-500"
              />
            </Link>
          )}
        </div>
      </div>

      <div className="pt-8 space-y-6">
        <h3 className="text-xl font-black uppercase italic tracking-tighter">
          Response_Log ({post.comments?.length || 0})
        </h3>

        {user ? (
          <CommentForm
            onSubmit={async (c) => {
              await createComment(post.id, c);
              fetchThread();
            }}
          />
        ) : (
          <div className="border-4 border-dashed border-black p-6 text-center opacity-40 font-black uppercase">
            Authentication_Required
          </div>
        )}

        <div className="space-y-4">
          {post.comments?.map((comment) => (
            <div
              key={comment.id}
              className="border-4 border-black bg-white flex items-stretch hover:translate-x-1 hover:-translate-y-1 transition-transform overflow-hidden"
            >
              <VoteSidebar
                small
                upvotes={comment.upvotes}
                userVote={comment.user_vote}
                onVote={(v) => handleCommentVote(comment.id, v)}
              />
              <div className="flex-1 p-4 space-y-2">
                <div className="flex justify-between text-[10px] font-black uppercase opacity-40">
                  <span>u/{comment.creator?.username}</span>
                  <span>
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
