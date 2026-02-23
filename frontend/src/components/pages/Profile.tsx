"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import ProfileHeader from "@/components/profile/ProfileHeader";
import ForumCard from "@/components/forum/ForumCard";
import MazeRepositoryCard from "@/components/profile/MazeRepositoryCard";
import {
  getProfile,
  deleteMaze as apiDeleteMaze,
  deletePost as apiDeletePost,
  deleteComment as apiDeleteComment,
} from "@/lib/api";
import ProfileCommentCard from "../profile/ProfileCommentCard";

type TabType = "mazes" | "posts" | "comments";

export default function ProfilePage() {
  const { username } = useParams();
  const { user: currentUser } = useAuth();
  const [profile, setProfile] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<TabType>("mazes");

  const isOwner = currentUser === username;

  useEffect(() => {
    getProfile(username as string)
      .then((data) => {
        setProfile(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, [username]);

  const performDelete = async (
    id: string,
    type: TabType,
    apiCall: (id: string) => Promise<boolean>
  ) => {
    if (
      !confirm(
        `CRITICAL: Purge this ${type.slice(0, -1)}? This action is irreversible.`
      )
    )
      return;

    const success = await apiCall(id);

    if (success) {
      setProfile((prev: any) => {
        if (!prev) return prev;

        let nextMazes = [...prev.mazes];
        let nextPosts = [...prev.posts];
        let nextComments = [...prev.comments];

        if (type === "mazes") {
          nextMazes = nextMazes.filter((m) => m.id !== id);

          const affectedPostIds = nextPosts
            .filter((p) => p.maze_id === id)
            .map((p) => p.id);

          nextPosts = nextPosts.filter((p) => p.maze_id !== id);

          nextComments = nextComments.filter(
            (c) => !affectedPostIds.includes(c.post_id)
          );
        } else if (type === "posts") {
          nextPosts = nextPosts.filter((p) => p.id !== id);

          nextComments = nextComments.filter((c) => c.post_id !== id);
        } else if (type === "comments") {
          nextComments = nextComments.filter((c) => c.id !== id);
        }

        return {
          ...prev,
          mazes: nextMazes,
          posts: nextPosts,
          comments: nextComments,
          stats: {
            total_mazes: nextMazes.length,
            total_posts: nextPosts.length,
            total_comments: nextComments.length,
          },
        };
      });
    }
  };

  const handlePostVoteUpdate = (
    postId: string,
    newVote: number,
    newCount: number
  ) => {
    setProfile((prev: any) => {
      if (!prev) return prev;
      return {
        ...prev,
        posts: prev.posts.map((p: any) =>
          p.id === postId ? { ...p, user_vote: newVote, upvotes: newCount } : p
        ),
      };
    });
  };

  if (loading)
    return (
      <div className="p-20 font-black italic min-h-screen animate-pulse uppercase flex items-center justify-center">
        INITIALIZING_PROFILE...
      </div>
    );

  if (!profile)
    return (
      <div className="p-20 font-black text-red-600 min-h-screen uppercase flex items-center justify-center">
        ERROR: USER_NOT_FOUND
      </div>
    );

  return (
    <div className="max-w-6xl mx-auto py-12 px-6 space-y-12 min-h-screen">
      <ProfileHeader
        username={profile.username}
        joinedAt={profile.created_at}
      />

      {/* Tab Switcher */}
      <div className="flex border-4 border-black bg-black p-1 gap-1 self-start shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
        {(["mazes", "posts", "comments"] as TabType[]).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-6 py-2 text-xs font-black uppercase transition-all ${
              activeTab === tab
                ? "bg-white text-black"
                : "text-white hover:bg-zinc-800"
            }`}
          >
            {tab} ({profile[tab]?.length || 0})
          </button>
        ))}
      </div>

      <div className="space-y-8 pb-20">
        {/* MAZES VIEW */}
        {activeTab === "mazes" && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 items-stretch">
            {profile.mazes?.map((maze: any) => (
              <MazeRepositoryCard
                key={maze.id}
                maze={maze}
                isOwner={isOwner}
                onDelete={() => performDelete(maze.id, "mazes", apiDeleteMaze)}
              />
            ))}
          </div>
        )}

        {/* POSTS VIEW */}
        {activeTab === "posts" && (
          <div className="space-y-6">
            {profile.posts?.map((post: any) => (
              <ForumCard
                key={post.id}
                post={post}
                isOwner={isOwner}
                onDelete={() => performDelete(post.id, "posts", apiDeletePost)}
                onVoteUpdate={handlePostVoteUpdate}
              />
            ))}
          </div>
        )}

        {/* COMMENTS VIEW */}
        {activeTab === "comments" && (
          <div className="space-y-6">
            {profile.comments?.map((comment: any) => (
              <ProfileCommentCard
                key={comment.id}
                comment={comment}
                isOwner={isOwner}
                onDelete={(id) =>
                  performDelete(id, "comments", apiDeleteComment)
                }
              />
            ))}
          </div>
        )}

        {/* EMPTY STATE */}
        {profile[activeTab]?.length === 0 && (
          <div className="py-32 border-4 border-dashed border-black text-center opacity-30 font-black uppercase italic">
            NO_{activeTab.toUpperCase()}_DETECTED_IN_STORAGE
          </div>
        )}
      </div>
    </div>
  );
}
