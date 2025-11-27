import { memo } from 'react';
import Link from 'next/link';

function Footer() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-card/80 backdrop-blur-sm border-t border-border/50 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Brand Section */}
          <div>
            <h3 className="text-foreground text-lg font-bold mb-4">PelotonLive</h3>
            <p className="text-sm text-muted-foreground">
              Professional cycling streaming platform. Watch live races and events from around the world.
            </p>
          </div>

          {/* Legal Links */}
          <div>
            <h4 className="text-foreground text-sm font-semibold mb-4 uppercase tracking-wider">Legal</h4>
            <ul className="space-y-2">
              <li>
                <Link
                  href="/terms"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Terms of Service
                </Link>
              </li>
              <li>
                <Link
                  href="/privacy"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Privacy Policy
                </Link>
              </li>
              <li>
                <Link
                  href="/contact"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Contact Us
                </Link>
              </li>
            </ul>
          </div>

          {/* Quick Links */}
          <div>
            <h4 className="text-foreground text-sm font-semibold mb-4 uppercase tracking-wider">Quick Links</h4>
            <ul className="space-y-2">
              <li>
                <Link
                  href="/"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Home
                </Link>
              </li>
              <li>
                <Link
                  href="/"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Browse Races
                </Link>
              </li>
              <li>
                <Link
                  href="/auth/login"
                  className="text-sm text-muted-foreground hover:text-primary transition-colors"
                >
                  Login
                </Link>
              </li>
            </ul>
          </div>
        </div>

        {/* Copyright */}
        <div className="mt-8 pt-8 border-t border-border/50">
          <p className="text-center text-sm text-muted-foreground">
            © {currentYear} PelotonLive. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}

export default memo(Footer);

