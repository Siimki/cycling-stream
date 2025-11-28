'use client';

export function FairModelSection() {
  return (
    <section id="how-it-works" className="py-12 mb-12">
      <div className="max-w-5xl mx-auto text-center">
        <div className="inline-block bg-primary/20 text-primary text-sm font-bold px-3 py-1 rounded-full mb-4">
          THE FAIRMUS MODEL
        </div>
        <h2 className="text-4xl sm:text-5xl font-black mb-6">Your View is Your Vote</h2>
        <p className="text-muted-foreground text-xl mb-16 max-w-3xl mx-auto">
          Support the sport you love. Revenue is distributed based on what you actually watch—no corporate gatekeepers deciding which races matter.
        </p>

        {/* Flow Diagram */}
        <div className="flex flex-col md:flex-row items-center justify-center gap-8 md:gap-6">
          <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-2xl p-8 flex-1 max-w-xs">
            <div className="w-16 h-16 bg-primary/20 rounded-full flex items-center justify-center mx-auto mb-4">
              <span className="text-3xl">💳</span>
            </div>
            <h3 className="font-bold text-lg mb-2">You Pay $5/month</h3>
            <p className="text-muted-foreground text-sm">Unlimited access to Tier 2/3 races worldwide</p>
          </div>

          <div className="text-primary text-3xl hidden md:block">→</div>

          <div className="bg-card/80 backdrop-blur-sm border border-primary rounded-2xl p-8 flex-1 max-w-xs shadow-lg shadow-primary/20">
            <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center mx-auto mb-4">
              <span className="text-3xl">📊</span>
            </div>
            <h3 className="font-bold text-lg mb-2">Platform (50%)</h3>
            <p className="text-muted-foreground text-sm">Streaming, hosting, and development costs</p>
          </div>

          <div className="text-primary text-3xl">+</div>

          <div className="bg-card/80 backdrop-blur-sm border border-primary rounded-2xl p-8 flex-1 max-w-xs shadow-lg shadow-primary/20">
            <div className="w-16 h-16 bg-primary rounded-full flex items-center justify-center mx-auto mb-4">
              <span className="text-3xl">🏆</span>
            </div>
            <h3 className="font-bold text-lg mb-2">Race Organizers (50%)</h3>
            <p className="text-muted-foreground text-sm">Split proportionally by your watch time</p>
          </div>
        </div>

        <div className="mt-12 bg-card/80 backdrop-blur-sm border border-border/50 rounded-xl p-6 max-w-2xl mx-auto">
          <p className="text-foreground text-sm">
            <span className="font-bold text-primary">Example:</span> If you watch 80% Tour of Estonia and 20% Flanders, then 80% of your organizer share goes to Estonia, 20% to Flanders.
          </p>
        </div>
      </div>
    </section>
  );
}

