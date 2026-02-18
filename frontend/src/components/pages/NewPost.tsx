"use client";
import { useState, useEffect, useMemo } from "react";
import { useRouter } from "next/navigation";
import { loadGenerateSession } from "@/lib/db";
import { useAuth } from "@/context/AuthContext";
import { getMyMazes, createPost } from "@/lib/api"; // Integrated lib/api

export default function NewPost() {
  const { user } = useAuth();
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [userMazes, setUserMazes] = useState<any[]>([]);
  const [selectedMazeId, setSelectedMazeId] = useState<string>("");

  useEffect(() => {
    if (!user) router.push("/forum");

    const init = async () => {
      const currentSession = await loadGenerateSession();
      if (currentSession) setSelectedMazeId(currentSession.id);

      try {
        const data = await getMyMazes();
        setUserMazes(data || []);
      } catch (err) {
        console.error("MAZE_FETCH_FAILED", err);
      }
    };
    init();
  }, [user, router]);

  const selectedMaze = useMemo(
    () => userMazes.find((m) => m.id === selectedMazeId),
    [selectedMazeId, userMazes]
  );

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (loading) return;

    setLoading(true);
    const formData = new FormData(e.currentTarget);

    try {
      const success = await createPost({
        title: formData.get("title") as string,
        content: formData.get("content") as string,
        maze_id: selectedMazeId,
      });

      if (success) router.push("/forum");
    } catch (err) {
      console.error("POST_FAILED", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto pb-20">
      <h1 className="text-3xl font-black uppercase tracking-tighter mb-8 italic">
        INITIALIZE_THREAD
      </h1>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-2">
          <label className="block text-[10px] font-black uppercase opacity-40">
            Thread_Title
          </label>
          <input
            name="title"
            required
            placeholder="A fascinating 200x300 matrix analysis..."
            className="w-full border-4 border-black p-4 font-bold text-sm focus:bg-zinc-50 outline-none transition-colors"
          />
        </div>

        <div className="space-y-2">
          <label className="block text-[10px] font-black uppercase opacity-40">
            Content_Payload
          </label>
          <textarea
            name="content"
            required
            rows={6}
            placeholder="Explain the algorithmic complexity or visual patterns..."
            className="w-full border-4 border-black p-4 font-bold text-sm focus:bg-zinc-50 outline-none transition-colors"
          />
        </div>

        <div className="space-y-2">
          <label className="block text-[10px] font-black uppercase opacity-40">
            Attach_Matrix_Reference
          </label>

          {/* Side-by-side Selection and Thumbnail */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 h-full md:h-40">
            <div className="md:col-span-2 border-4 border-black bg-zinc-50 flex flex-col justify-between">
              <select
                value={selectedMazeId}
                onChange={(e) => setSelectedMazeId(e.target.value)}
                className="w-full p-4 font-black uppercase text-xs bg-transparent outline-none appearance-none cursor-pointer h-full"
              >
                <option value="">-- NO_ATTACHMENT --</option>
                {userMazes.map((m) => (
                  <option key={m.id} value={m.id}>
                    {m.id} ({m.rows}x{m.cols}) -{" "}
                    {new Date(m.created_at).toLocaleDateString()}
                  </option>
                ))}
              </select>

              {selectedMazeId && (
                <div className="border-t-2 border-black p-2 bg-white flex items-center justify-between shrink-0">
                  <span className="text-[9px] font-black uppercase text-green-600">
                    READY_FOR_LINKING
                  </span>
                  <button
                    type="button"
                    onClick={() => setSelectedMazeId("")}
                    className="text-[9px] font-black uppercase underline hover:text-red-600 transition-colors"
                  >
                    Detach
                  </button>
                </div>
              )}
            </div>

            {/* Thumbnail Preview Area */}
            <div className="border-4 border-black bg-white flex items-center justify-center relative overflow-hidden group min-h-[150px] md:min-h-0">
              {selectedMaze?.thumbnail ? (
                <img
                  src={selectedMaze.thumbnail}
                  alt="Preview"
                  className="w-full h-full object-cover"
                />
              ) : (
                <div className="text-[8px] font-black uppercase opacity-20 text-center p-4">
                  NO_PREVIEW_AVAILABLE
                </div>
              )}
            </div>
          </div>
        </div>

        <button
          disabled={loading}
          type="submit"
          className="w-full border-4 border-black p-4 bg-black text-white font-black uppercase italic hover:bg-zinc-800 disabled:opacity-50 transition-all shadow-[6px_6px_0px_0px_rgba(0,0,0,0.3)] active:shadow-none active:translate-x-1 active:translate-y-1"
        >
          {loading ? "TRANSMITTING..." : "PUBLISH_POST"}
        </button>
      </form>
    </div>
  );
}
