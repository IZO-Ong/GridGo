import { Post, Comment, Maze, User } from "@/types";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

/**
 * Returns a strictly typed record for headers.
 */
const getAuthHeaders = (): Record<string, string> => {
  if (typeof window === "undefined") return {};
  const token = localStorage.getItem("gridgo_token");
  return token ? { Authorization: `Bearer ${token}` } : {};
};

// --- Auth APIs ---

export async function login(payload: any) {
  const res = await fetch(`${BASE_URL}/api/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!res.ok)
    throw new Error((await res.json()).error || "INVALID_CREDENTIALS");
  return res.json();
}

export async function register(payload: any) {
  const res = await fetch(`${BASE_URL}/api/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!res.ok)
    throw new Error((await res.json()).error || "REGISTRATION_FAILURE");
  return true;
}

export async function verifyAccount(email: string, code: string) {
  const res = await fetch(`${BASE_URL}/api/verify`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, code }),
  });
  if (!res.ok) throw new Error("VERIFICATION_FAILURE");
  return true;
}

// --- Maze APIs ---

export async function generateMaze(
  formData: FormData,
  token?: string | null
): Promise<Maze> {
  const headers: Record<string, string> = {};
  if (token) headers.Authorization = `Bearer ${token}`;

  const res = await fetch(`${BASE_URL}/api/maze/generate`, {
    method: "POST",
    body: formData,
    headers,
  });

  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function renderMazeImage(mazeData: Maze): Promise<Blob> {
  const response = await fetch(`${BASE_URL}/api/maze/render`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(mazeData),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || "SYSTEM_RENDER_FAILURE");
  }

  return response.blob();
}

export async function solveMaze(mazeData: Maze, algorithm: string) {
  const response = await fetch(`${BASE_URL}/api/maze/solve`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ maze: mazeData, algorithm }),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || "SYSTEM_SOLVE_FAILURE");
  }

  return response.json();
}

export async function getMazeById(id: string): Promise<Maze> {
  const response = await fetch(`${BASE_URL}/api/maze/get?id=${id}`);
  if (!response.ok) throw new Error("MAZE_NOT_FOUND");
  return response.json();
}

export async function deleteMaze(mazeId: string): Promise<boolean> {
  const res = await fetch(`${BASE_URL}/api/maze/delete?id=${mazeId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });
  return res.ok;
}

export async function updateThumbnail(id: string, thumbnail: string) {
  return fetch(`${BASE_URL}/api/maze/thumbnail`, {
    method: "PUT",
    headers: { "Content-Type": "application/json", ...getAuthHeaders() },
    body: JSON.stringify({ id, thumbnail }),
  });
}

// --- Forum and Profile APIs ---

export async function getProfile(username: string): Promise<User> {
  const res = await fetch(`${BASE_URL}/api/profile?username=${username}`, {
    headers: getAuthHeaders(),
  });
  if (!res.ok) throw new Error("USER_NOT_FOUND");
  return res.json();
}

export async function getMyMazes(): Promise<Maze[]> {
  const res = await fetch(`${BASE_URL}/api/maze/my-mazes`, {
    headers: getAuthHeaders(),
  });
  if (!res.ok) return [];
  return res.json();
}

export async function getPosts(offset: number = 0): Promise<Post[]> {
  const res = await fetch(`${BASE_URL}/api/forum/posts?offset=${offset}`, {
    headers: getAuthHeaders(),
  });
  if (!res.ok) return [];
  return res.json();
}

export const getPostById = async (id: string): Promise<Post> => {
  const res = await fetch(`${BASE_URL}/api/forum/post?id=${id}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(), // Fixed: Now sends token to hydrate user_vote
    },
  });

  if (!res.ok) throw new Error("THREAD_FETCH_FAILED");
  return res.json();
};

export async function createPost(postData: {
  title: string;
  content: string;
  maze_id?: string;
}): Promise<boolean> {
  const res = await fetch(`${BASE_URL}/api/forum/posts/create`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify(postData),
  });

  if (!res.ok) throw new Error(await res.text());
  return true;
}

export async function castVote(
  target_id: string,
  target_type: "post" | "comment",
  value: number
): Promise<boolean> {
  const res = await fetch(`${BASE_URL}/api/forum/vote`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify({ target_id, target_type, value }),
  });

  return res.ok;
}

export async function createComment(
  post_id: string,
  content: string
): Promise<Comment> {
  const res = await fetch(`${BASE_URL}/api/forum/comment/create`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
    },
    body: JSON.stringify({ post_id, content }),
  });

  if (!res.ok) throw new Error("COMMENT_TRANSMISSION_FAILURE");
  return res.json();
}

export async function deletePost(postId: string): Promise<boolean> {
  const res = await fetch(`${BASE_URL}/api/forum/post/delete?id=${postId}`, {
    method: "DELETE",
    headers: getAuthHeaders(),
  });

  if (!res.ok) {
    const errorText = await res.text();
    console.error("POST_DELETE_FAILURE:", errorText);
  }
  return res.ok;
}

export async function deleteComment(commentId: string): Promise<boolean> {
  const res = await fetch(
    `${BASE_URL}/api/forum/comment/delete?id=${commentId}`,
    {
      method: "DELETE",
      headers: getAuthHeaders(),
    }
  );

  if (!res.ok) {
    const errorText = await res.text();
    console.error("COMMENT_DELETE_FAILURE:", errorText);
  }
  return res.ok;
}
