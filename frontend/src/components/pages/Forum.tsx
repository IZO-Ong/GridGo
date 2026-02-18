"use client";
import { useState, useEffect } from "react";
import ForumCard from "@/components/forum/ForumCard";
import ForumFAB from "@/components/forum/ForumFAB";
import { useAuth } from "@/context/AuthContext";
import { getPosts } from "@/lib/api";
import { Post } from "@/types";

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

      setPosts((prev: Post[]) => {
        const existingIds = new Set(prev.map((p) => p.id));
        const uniqueNewPosts = newPosts.filter((p) => !existingIds.has(p.id));
        return [...prev, ...uniqueNewPosts];
      });

      setOffset((prev: number) => prev + 10);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-5xl mx-auto pb-20 relative min-h-screen">
      {user && <ForumFAB />}

      <div className="max-w-3xl mx-auto space-y-4 cursor-default">
        {posts.map((post) => (
          <ForumCard key={post.id} post={post} />
        ))}

        {loading && (
          <div className="p-10 text-center font-black italic animate-pulse uppercase text-[11px]">
            SYNCHRONIZING_FEED...
          </div>
        )}

        {!loading && hasMore && posts.length > 0 && (
          <button
            onClick={fetchPosts}
            className="w-full border-4 border-black p-4 font-black uppercase italic hover:bg-black hover:text-white transition-all cursor-pointer"
          >
            LOAD_MORE_POSTS
          </button>
        )}
      </div>
    </div>
  );
}
