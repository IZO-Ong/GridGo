export default function ProfileHeader({
  username,
  joinedAt,
}: {
  username: string;
  joinedAt: string;
}) {
  return (
    <section className="relative border-4 border-black p-8 bg-white shadow-[12px_12px_0px_0px_rgba(0,0,0,1)] flex flex-col md:flex-row items-center gap-8">
      <div className="w-32 h-32 bg-black text-white border-4 border-black flex items-center justify-center text-5xl font-black shrink-0 shadow-[4px_4px_0px_0px_rgba(0,0,0,0.2)]">
        {username[0].toUpperCase()}
      </div>
      <div className="flex-1 text-center md:text-left">
        <h1 className="text-5xl font-black uppercase tracking-tighter leading-none mb-2">
          {username}
        </h1>
        <p className="text-sm font-bold opacity-40 uppercase">
          JOINED // {new Date(joinedAt).toLocaleDateString()}
        </p>
      </div>
    </section>
  );
}
