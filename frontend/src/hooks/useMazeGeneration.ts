"use client";
import { useState } from "react";
import { generateMaze } from "@/lib/api";
import { MazeData } from "@/types";

export function useMazeGeneration() {
  const [maze, setMaze] = useState<MazeData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const executeGeneration = async (
    formData: FormData
  ): Promise<MazeData | null> => {
    setLoading(true);
    setError(null);

    // Retrieve the token from storage
    const token =
      typeof window !== "undefined"
        ? localStorage.getItem("gridgo_token")
        : null;

    try {
      const data = await generateMaze(formData, token);
      setMaze(data);
      return data;
    } catch (err: any) {
      const msg = err.message || "SYSTEM_FAILURE";
      setError(msg);
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { maze, loading, error, executeGeneration, setError };
}
