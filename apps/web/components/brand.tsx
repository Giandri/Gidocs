import Image from "next/image";

export function Brand() {
  return (
    <>
      <Image alt="" aria-hidden="true" className="h-[18px] w-auto" height={232} src="/gidocs.png" width={190} />
      <p className="mt-0.5 text-[17px] font-bold tracking-[-0.08em] text-[#f3f1f1] sm:text-[20px]">
        Gi<span className="text-accent">Docs</span>
      </p>
    </>
  );
}
