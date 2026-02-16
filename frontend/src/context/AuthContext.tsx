"use client";
import { createContext, useContext, useEffect, useState } from "react";

interface AuthContextType {
  user: string | null;
  login: (token: string, username: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<string | null>(null);

  useEffect(() => {
    const storedUser = localStorage.getItem("gridgo_user");
    if (storedUser) setUser(storedUser);
  }, []);

  const login = (token: string, username: string) => {
    localStorage.setItem("gridgo_token", token);
    localStorage.setItem("gridgo_user", username);
    setUser(username);
  };

  const logout = () => {
    localStorage.removeItem("gridgo_token");
    localStorage.removeItem("gridgo_user");
    setUser(null);
    window.location.href = "/login";
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
