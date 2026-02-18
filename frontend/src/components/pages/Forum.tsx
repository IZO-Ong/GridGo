"use client";
import { useState, useEffect } from "react";
import ForumCard from "@/components/forum/ForumCard";
import { useAuth } from "@/context/AuthContext";
import { getPosts } from "@/lib/api";
import Link from "next/link";

interface Post {
  id: string;
  title: string;
  content: string;
  maze_id?: string;
  upvotes: number;
  created_at: string;
  creator: {
    username: string;
  };
}

export default function ForumPage() {
  const { user } = useAuth();
  const [posts, setPosts] = useState<Post[]>([]);
  const [offset, setOffset] = useState(0);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    fetchPosts();
  }, []);

  const fetchPosts = async () => {
    if (loading) return;
    setLoading(true);
    try {
      const newPosts: Post[] = await getPosts(offset);
      if (newPosts.length < 10) setHasMore(false);

      setPosts((prev) => {
        const existingIds = new Set(prev.map((p) => p.id));
        const uniqueNewPosts = newPosts.filter((p) => !existingIds.has(p.id));
        return [...prev, ...uniqueNewPosts];
      });

      setOffset((prev) => prev + 10);
    } catch (err) {
      console.error("FEED_ERROR", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto space-y-6 pb-20">
      {user && (
        <Link
          href="/forum/new"
          className="block border-4 border-black p-4 bg-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:translate-x-1 hover:translate-y-1 hover:shadow-none transition-all"
        >
          <p className="font-bold uppercase text-xs opacity-50 italic">
            {" "}
            START_NEW_THREAD...
          </p>
        </Link>
      )}

      <div className="space-y-4">
        {posts.map((post) => (
          <ForumCard key={post.id} post={post} />
        ))}
      </div>

      {!loading && posts.length === 0 && (
        <div className="p-20 border-4 border-dashed border-zinc-300 text-center space-y-4">
          <p className="opacity-40 tracking-[0.3em] font-black uppercase text-xl">
            Nothing to see here
          </p>
          <p className="text-[10px] font-bold uppercase opacity-30">
            How about creating a post?
          </p>
        </div>
      )}

      {loading && (
        <div className="p-10 text-center font-black italic animate-pulse uppercase text-[11px]">
          SYNCHRONIZING_FEED...
        </div>
      )}

      {!loading && hasMore && posts.length > 0 && (
        <button
          onClick={fetchPosts}
          className="w-full border-4 border-black p-4 font-black uppercase italic hover:bg-black hover:text-white transition-all"
        >
          LOAD_MORE_POSTS
        </button>
      )}
    </div>
  );
}
