import { memo } from 'react';
import Link from 'next/link';

function Footer() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-black border-t border-border/60 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <p className="text-center text-[12px] text-[#666] flex flex-wrap items-center justify-center gap-2">
          © {currentYear} PelotonLive
          <span className="text-[#444]">•</span>
          <Link href="/terms" className="hover:text-foreground transition-colors">Terms</Link>
          <span className="text-[#444]">•</span>
          <Link href="/privacy" className="hover:text-foreground transition-colors">Privacy</Link>
          <span className="text-[#444]">•</span>
          <Link href="/contact" className="hover:text-foreground transition-colors">Help</Link>
        </p>
      </div>
    </footer>
  );
}

export default memo(Footer);
