import Link from 'next/link';
import { Navigation } from '@/components/layout/Navigation';
import Footer from '@/components/layout/Footer';

export const metadata = {
  title: 'Privacy Policy - PelotonLive',
  description: 'Privacy Policy for PelotonLive platform',
};

export default function PrivacyPage() {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 w-full">
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8 lg:p-12">
          <h1 className="text-3xl sm:text-4xl font-bold text-foreground/95 mb-4">Privacy Policy</h1>
          <p className="text-muted-foreground mb-8 text-sm sm:text-base">
            Last updated: {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>

          <div className="prose prose-sm sm:prose-base max-w-none space-y-6 text-foreground/90">
            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">1. Introduction</h2>
              <p className="text-muted-foreground">
                PelotonLive (&quot;we,&quot; &quot;our,&quot; or &quot;us&quot;) is committed to protecting your privacy. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you use our streaming platform.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">2. Information We Collect</h2>
              
              <h3 className="text-lg sm:text-xl font-semibold text-foreground/95 mt-4 mb-2">2.1 Information You Provide</h3>
              <p className="text-muted-foreground">We collect information that you provide directly to us, including:</p>
              <ul className="list-disc pl-6 mt-2 space-y-1 text-muted-foreground">
                <li>Account registration information (email address, name, password)</li>
                <li>Payment information (processed securely through third-party payment processors)</li>
                <li>Profile information and preferences</li>
                <li>Communications with us (support requests, feedback)</li>
              </ul>

              <h3 className="text-lg sm:text-xl font-semibold text-foreground/95 mt-4 mb-2">2.2 Automatically Collected Information</h3>
              <p className="text-muted-foreground">When you use our Platform, we automatically collect certain information, including:</p>
              <ul className="list-disc pl-6 mt-2 space-y-1 text-muted-foreground">
                <li>Device information (IP address, browser type, operating system)</li>
                <li>Usage data (pages visited, time spent, features used)</li>
                <li>Watch session data (viewing history, watch time, race preferences)</li>
                <li>Analytics data (concurrent viewers, unique viewers per race)</li>
                <li>Cookies and similar tracking technologies</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">3. How We Use Your Information</h2>
              <p className="text-muted-foreground">We use the information we collect to:</p>
              <ul className="list-disc pl-6 mt-2 space-y-1 text-muted-foreground">
                <li>Provide, maintain, and improve our services</li>
                <li>Process transactions and send related information</li>
                <li>Authenticate users and prevent fraud</li>
                <li>Send technical notices, updates, and support messages</li>
                <li>Respond to your comments, questions, and requests</li>
                <li>Monitor and analyze usage patterns and trends</li>
                <li>Personalize your experience and content recommendations</li>
                <li>Detect, prevent, and address technical issues</li>
                <li>Comply with legal obligations</li>
              </ul>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">4. Information Sharing and Disclosure</h2>
              <p className="text-muted-foreground">We do not sell your personal information. We may share your information in the following circumstances:</p>
              
              <h3 className="text-lg sm:text-xl font-semibold text-foreground/95 mt-4 mb-2">4.1 Service Providers</h3>
              <p className="text-muted-foreground">
                We may share information with third-party service providers who perform services on our behalf, such as payment processing, data analytics, hosting, and customer support.
              </p>

              <h3 className="text-lg sm:text-xl font-semibold text-foreground/95 mt-4 mb-2">4.2 Business Transfers</h3>
              <p className="text-muted-foreground">
                If we are involved in a merger, acquisition, or sale of assets, your information may be transferred as part of that transaction.
              </p>

              <h3 className="text-lg sm:text-xl font-semibold text-foreground/95 mt-4 mb-2">4.3 Legal Requirements</h3>
              <p className="text-muted-foreground">
                We may disclose your information if required to do so by law or in response to valid requests by public authorities.
              </p>

              <h3 className="text-lg sm:text-xl font-semibold text-foreground/95 mt-4 mb-2">4.4 Aggregated Data</h3>
              <p className="text-muted-foreground">
                We may share aggregated, anonymized data that does not identify individual users for analytics, research, or business purposes.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">5. Data Security</h2>
              <p className="text-muted-foreground">
                We implement appropriate technical and organizational measures to protect your personal information against unauthorized access, alteration, disclosure, or destruction. However, no method of transmission over the Internet or electronic storage is 100% secure, and we cannot guarantee absolute security.
              </p>
              <p className="text-muted-foreground mt-3">
                Your account password is encrypted using industry-standard hashing algorithms. We recommend using a strong, unique password and not sharing your account credentials with others.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">6. Cookies and Tracking Technologies</h2>
              <p className="text-muted-foreground">
                We use cookies and similar tracking technologies to track activity on our Platform and hold certain information. Cookies are files with a small amount of data that may include an anonymous unique identifier.
              </p>
              <p className="text-muted-foreground mt-3">
                You can instruct your browser to refuse all cookies or to indicate when a cookie is being sent. However, if you do not accept cookies, you may not be able to use some portions of our Platform.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">7. Your Rights and Choices</h2>
              <p className="text-muted-foreground">Depending on your location, you may have certain rights regarding your personal information, including:</p>
              <ul className="list-disc pl-6 mt-2 space-y-1 text-muted-foreground">
                <li>The right to access your personal information</li>
                <li>The right to correct inaccurate information</li>
                <li>The right to delete your account and associated data</li>
                <li>The right to object to or restrict certain processing</li>
                <li>The right to data portability</li>
                <li>The right to withdraw consent where processing is based on consent</li>
              </ul>
              <p className="text-muted-foreground mt-3">
                To exercise these rights, please contact us through our{' '}
                <Link href="/contact" className="text-primary hover:underline">
                  contact page
                </Link>.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">8. Data Retention</h2>
              <p className="text-muted-foreground">
                We retain your personal information for as long as necessary to fulfill the purposes outlined in this Privacy Policy, unless a longer retention period is required or permitted by law. When we no longer need your information, we will securely delete or anonymize it.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">9. Children&apos;s Privacy</h2>
              <p className="text-muted-foreground">
                Our Platform is not intended for children under the age of 13. We do not knowingly collect personal information from children under 13. If you are a parent or guardian and believe your child has provided us with personal information, please contact us immediately.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">10. International Data Transfers</h2>
              <p className="text-muted-foreground">
                Your information may be transferred to and processed in countries other than your country of residence. These countries may have data protection laws that differ from those in your country. By using our Platform, you consent to the transfer of your information to these countries.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">11. Changes to This Privacy Policy</h2>
              <p className="text-muted-foreground">
                We may update this Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on this page and updating the &quot;Last updated&quot; date. You are advised to review this Privacy Policy periodically for any changes.
              </p>
            </section>

            <section>
              <h2 className="text-xl sm:text-2xl font-semibold text-foreground/95 mt-6 mb-3">12. Contact Us</h2>
              <p className="text-muted-foreground">
                If you have any questions about this Privacy Policy, please contact us through our{' '}
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

