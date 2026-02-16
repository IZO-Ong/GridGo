const BASE_URL = process.env.NEXT_PUBLIC_API_URL;

export async function generateMaze(formData: FormData) {
  const response = await fetch(`${BASE_URL}/maze/generate`, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || "SYSTEM_GENERATION_FAILURE");
  }

  return response.json();
}

export async function renderMazeImage(mazeData: any) {
  const response = await fetch(`${BASE_URL}/maze/render`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(mazeData),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(errorText || "SYSTEM_RENDER_FAILURE");
  }

  const contentType = response.headers.get("content-type");
  if (!contentType || !contentType.includes("image/png")) {
    throw new Error("INVALID_PAYLOAD_FORMAT");
  }

  return response.blob();
}

export async function solveMaze(mazeData: any, algorithm: string) {
  const response = await fetch(`${BASE_URL}/maze/solve`, {
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

export async function getMazeById(id: string) {
  const response = await fetch(`${BASE_URL}/maze/get?id=${id}`);
  if (!response.ok) throw new Error("MAZE_NOT_FOUND");
  return response.json();
}
