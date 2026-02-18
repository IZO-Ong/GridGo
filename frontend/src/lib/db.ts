import { MazeData } from "@/types";

export const DB_NAME = "GridGoDB";
export const STORE_NAME = "mazes";

export const initDB = (): Promise<IDBDatabase> => {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, 1);
    request.onupgradeneeded = () => {
      if (!request.result.objectStoreNames.contains(STORE_NAME)) {
        request.result.createObjectStore(STORE_NAME);
      }
    };
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
};

export const saveGenerateSession = async (maze: MazeData) => {
  const db = await initDB();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, "readwrite");
    const store = transaction.objectStore(STORE_NAME);
    const request = store.put(maze, "last_generated");
    request.onsuccess = () => resolve(true);
    request.onerror = () => reject(request.error);
  });
};

export const loadGenerateSession = async (): Promise<MazeData> => {
  const db = await initDB();
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, "readonly");
    const store = transaction.objectStore(STORE_NAME);
    const request = store.get("last_generated");
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
};

export const saveSolveSession = async (maze: MazeData) => {
  const db = await initDB();
  const transaction = db.transaction(STORE_NAME, "readwrite");
  transaction.objectStore(STORE_NAME).put(maze, "solve_session");
};

export const loadSolveSession = async (): Promise<MazeData> => {
  const db = await initDB();
  const transaction = db.transaction(STORE_NAME, "readonly");
  const request = transaction.objectStore(STORE_NAME).get("solve_session");
  return new Promise((resolve) => {
    request.onsuccess = () => resolve(request.result);
  });
};

export const savePreferences = async (
  key: "gen_prefs" | "solve_prefs",
  prefs: any
) => {
  const db = await initDB();
  const transaction = db.transaction(STORE_NAME, "readwrite");
  transaction.objectStore(STORE_NAME).put(prefs, key);
};

export const loadPreferences = async (
  key: "gen_prefs" | "solve_prefs"
): Promise<any> => {
  const db = await initDB();
  const transaction = db.transaction(STORE_NAME, "readonly");
  const request = transaction.objectStore(STORE_NAME).get(key);
  return new Promise((resolve) => {
    request.onsuccess = () => resolve(request.result);
  });
};
