"use client";
import { createContext, useContext, useEffect, useState } from "react";
import { clearSolveSession, clearGenerateSession } from "@/lib/db";

interface AuthContextType {
  user: string | null;
  login: (token: string, username: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const isTokenExpired = (token: string | null): boolean => {
  if (!token) return true;
  try {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      window
        .atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join("")
    );
    const { exp } = JSON.parse(jsonPayload);
    return Date.now() >= exp * 1000;
  } catch (error) {
    return true;
  }
};

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<string | null>(null);

  useEffect(() => {
    const token = localStorage.getItem("gridgo_token");
    const storedUser = localStorage.getItem("gridgo_user");

    if (token && storedUser) {
      if (isTokenExpired(token)) {
        console.warn("Session expired. Logging out.");
        logout();
      } else {
        setUser(storedUser);
      }
    }
  }, []);

  const login = async (token: string, username: string) => {
    await clearSolveSession();
    await clearGenerateSession();

    localStorage.setItem("gridgo_token", token);
    localStorage.setItem("gridgo_user", username);
    setUser(username);
  };

  const logout = async () => {
    await clearSolveSession();
    await clearGenerateSession();

    localStorage.removeItem("gridgo_token");
    localStorage.removeItem("gridgo_user");
    setUser(null);
    if (typeof window !== "undefined") {
      window.location.href = "/login";
    }
  };

  return (
    <AuthContext.Provider value={{ user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within an AuthProvider");
  return context;
};
