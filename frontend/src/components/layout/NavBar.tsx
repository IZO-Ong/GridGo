"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/context/AuthContext"; // 1. Import the hook

export default function NavBar() {
  const { user, logout } = useAuth(); // 2. Use global state instead of local state
  const pathname = usePathname();

  const navItems = [
    { label: "Create", href: "/" },
    { label: "Solve", href: "/solve" },
    { label: "Forum", href: "/forum" },
  ];

  // REMOVED: useEffect and useState for 'user'

  const handleLogout = () => {
    logout(); // 3. Use the context logout function
  };

  return (
    <div className="flex items-center gap-6">
      <div className="flex items-center gap-3">
        <Link
          href={user ? `/profile/${user}` : "/login"}
          className="flex items-center gap-3 group h-10"
        >
          <div
            className={`w-10 h-10 border-2 border-black rounded-full flex items-center justify-center shrink-0 transition-colors ${
              user ? "bg-black text-white" : "group-hover:bg-zinc-100"
            }`}
          >
            <span className="text-xs font-black">
              {user ? user[0].toUpperCase() : "G"}
            </span>
          </div>
          <div className="flex flex-col justify-center">
            <span className="text-[11px] font-black uppercase tracking-tighter leading-none">
              {user ? user : "Guest_Session"}
            </span>
            <span className="text-[8px] font-mono opacity-40 uppercase tracking-tight mt-1">
              {user ? "Authenticated" : "Click to Sign_In"}
            </span>
          </div>
        </Link>

        {user && (
          <button
            onClick={handleLogout}
            title="Terminate Session"
            className="w-6 h-6 border-2 border-black flex items-center justify-center hover:bg-black hover:text-white transition-colors cursor-pointer group"
          >
            <span className="text-[10px] font-black">X</span>
          </button>
        )}
      </div>

      <nav className="flex border-2 border-black divide-x-2 divide-black bg-white shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]">
        {navItems.map((item) => (
          <Link
            key={item.label}
            href={item.href}
            className={`px-4 py-1 text-xs font-black uppercase transition-colors ${
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
