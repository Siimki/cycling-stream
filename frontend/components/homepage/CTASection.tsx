'use client';

import Link from 'next/link';
import { Button } from '@/components/ui/button';

export function CTASection() {
  return (
    <section className="py-24 mb-12 bg-gradient-to-b from-background to-muted/20 relative overflow-hidden">
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-primary rounded-full blur-3xl" />
        <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-primary rounded-full blur-3xl" />
      </div>

      <div className="max-w-4xl mx-auto px-6 text-center relative z-10">
        <h2 className="text-4xl sm:text-6xl font-black mb-6">Ready to Start?</h2>
        <p className="text-muted-foreground text-xl mb-10">
          Join thousands of fans supporting grassroots cycling. First 7 days free.
        </p>
        <Link href="/auth/register">
          <Button className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-bold py-5 px-12 rounded-full text-lg transition transform hover:scale-105 shadow-lg shadow-primary/20">
            Create Free Account
          </Button>
        </Link>
        <p className="text-muted-foreground text-sm mt-6">No credit card required</p>
      </div>
    </section>
  );
}

