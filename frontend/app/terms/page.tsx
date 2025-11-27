import Link from 'next/link';
import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';

export const metadata = {
  title: 'Terms of Service - PelotonLive',
  description: 'Terms of Service for PelotonLive platform',
};

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 w-full">
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8 lg:p-12">
          <h1 className="text-3xl sm:text-4xl font-bold text-foreground/95 mb-4">Terms of Service</h1>
          <p className="text-muted-foreground mb-8 text-sm sm:text-base">
            Last updated: {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>

          <div className="prose prose-sm sm:prose-base max-w-none space-y-6 text-foreground/90">
            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">1. Acceptance of Terms</h2>
              <p className="text-muted-foreground">
                By accessing and using PelotonLive (&quot;the Platform&quot;), you accept and agree to be bound by the terms and provision of this agreement. If you do not agree to these Terms of Service, please do not use our service.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">2. Description of Service</h2>
              <p className="text-muted-foreground">
                PelotonLive is a streaming platform that provides live and on-demand professional cycling race content. We offer both free and paid access to race streams, depending on the specific race and subscription model.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">3. User Accounts</h2>
              <p className="text-muted-foreground">
                To access certain features of the Platform, you must create an account. You are responsible for:
              </p>
              <ul className="list-disc pl-6 mt-2 space-y-1 text-muted-foreground">
                <li>Maintaining the confidentiality of your account credentials</li>
                <li>All activities that occur under your account</li>
                <li>Providing accurate and complete information when creating your account</li>
                <li>Notifying us immediately of any unauthorized use of your account</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">4. Payment and Billing</h2>
              <p className="text-muted-foreground">
                For paid content, you agree to pay all fees associated with your purchase. All payments are processed through secure third-party payment processors. By making a purchase, you agree to:
              </p>
              <ul className="list-disc pl-6 mt-2 space-y-1 text-muted-foreground">
                <li>Provide current, complete, and accurate purchase and account information</li>
                <li>Promptly update account and payment information</li>
                <li>Pay all charges incurred by your account</li>
                <li>Comply with all applicable local, state, and federal laws regarding online payments</li>
              </ul>
              <p className="text-muted-foreground mt-3">
                All sales are final. Refunds may be provided at our sole discretion in exceptional circumstances.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">5. Content Usage</h2>
              <p className="text-muted-foreground">
                The content available on PelotonLive is protected by copyright and other intellectual property laws. You agree not to:
              </p>
              <ul className="list-disc pl-6 mt-2 space-y-1 text-muted-foreground">
                <li>Reproduce, distribute, or publicly display any content without authorization</li>
                <li>Record, download, or capture streamed content</li>
                <li>Share your account credentials with others</li>
                <li>Use the service for any commercial purpose without our written consent</li>
                <li>Reverse engineer, decompile, or disassemble any part of the Platform</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">6. User Conduct</h2>
              <p className="text-muted-foreground">
                You agree to use the Platform only for lawful purposes and in a way that does not infringe the rights of, restrict, or inhibit anyone else&apos;s use and enjoyment of the Platform. Prohibited behavior includes:
              </p>
              <ul className="list-disc pl-6 mt-2 space-y-1 text-muted-foreground">
                <li>Harassing, threatening, or abusing other users</li>
                <li>Posting or transmitting any unlawful, harmful, or offensive content</li>
                <li>Attempting to gain unauthorized access to the Platform or its systems</li>
                <li>Interfering with or disrupting the service or servers</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">7. Service Availability</h2>
              <p className="text-muted-foreground">
                We strive to provide reliable service but do not guarantee uninterrupted or error-free access. The Platform may be temporarily unavailable due to maintenance, updates, or circumstances beyond our control. We reserve the right to modify, suspend, or discontinue any part of the service at any time.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">8. Limitation of Liability</h2>
              <p className="text-muted-foreground">
                To the maximum extent permitted by law, PelotonLive shall not be liable for any indirect, incidental, special, consequential, or punitive damages, or any loss of profits or revenues, whether incurred directly or indirectly, or any loss of data, use, goodwill, or other intangible losses resulting from your use of the Platform.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">9. Changes to Terms</h2>
              <p className="text-muted-foreground">
                We reserve the right to modify these Terms of Service at any time. We will notify users of any material changes by posting the new Terms on this page and updating the &quot;Last updated&quot; date. Your continued use of the Platform after such modifications constitutes acceptance of the updated Terms.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">10. Contact Information</h2>
              <p className="text-muted-foreground">
                If you have any questions about these Terms of Service, please contact us through our{' '}
                <Link href="/contact" className="text-primary hover:underline">
                  contact page
                </Link>.
              </p>
            </section>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}

