"use client";
import React from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";

export default function NavBar() {
  const { user, logout } = useAuth();
  const pathname = usePathname();
  const router = useRouter();

  const navItems = [
    { label: "Create", href: "/" },
    { label: "Solve", href: "/solve" },
    { label: "Forum", href: "/forum" },
  ];

  const handleLogout = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    e.stopPropagation();
    logout();
    router.push("/");
  };

  return (
    // Changed: flex-wrap allows nav and identity to stack on mobile
    // Added: justify-between for mobile alignment, md:justify-end for desktop
    <div className="flex flex-wrap items-center justify-between md:justify-end gap-y-4 gap-x-6 w-full md:w-auto">
      <div className="flex items-center gap-3">
        {/* Main wrapper for the identity section */}
        <div className="flex items-center gap-3 h-10">
          {/* 1. Circle Avatar Link */}
          <Link
            href={user ? `/profile/${user}` : "/login"}
            className="flex items-center shrink-0"
          >
            <div
              className={`w-10 h-10 border-2 border-black rounded-full flex items-center justify-center transition-colors ${
                user ? "bg-black text-white" : "hover:bg-zinc-100"
              }`}
            >
              <span className="text-xs font-black">
                {user ? user[0].toUpperCase() : "G"}
              </span>
            </div>
          </Link>

          {/* 2. Vertical Container for Name (Link) and Logout (Button) */}
          <div className="flex flex-col justify-center">
            <Link
              href={user ? `/profile/${user}` : "/login"}
              className="text-[11px] font-black uppercase tracking-tighter leading-none hover:underline"
            >
              {user ? user : "Guest"}
            </Link>

            {user ? (
              <button
                onClick={handleLogout}
                className="text-[8px] font-mono opacity-40 hover:opacity-100 hover:underline text-left uppercase tracking-tight mt-1 transition-opacity"
              >
                [Logout]
              </button>
            ) : (
              <Link
                href="/login"
                className="text-[8px] font-mono opacity-40 hover:opacity-100 hover:underline text-left uppercase tracking-tight mt-1 transition-opacity"
              >
                [Sign In]
              </Link>
            )}
          </div>
        </div>
      </div>

      <nav className="flex border-2 border-black divide-x-2 divide-black bg-white shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] shrink-0">
        {navItems.map((item) => (
          <Link
            key={item.label}
            href={item.href}
            className={`px-3 md:px-4 py-1 text-[10px] md:text-xs font-black uppercase transition-colors ${
              pathname === item.href
                ? "bg-black text-white"
                : "hover:bg-zinc-100 text-black"
            }`}
          >
            {item.label}
          </Link>
        ))}
      </nav>
    </div>
  );
}
