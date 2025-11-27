'use client';

import { useState } from 'react';
import Link from 'next/link';
import { Navigation } from '@/components/layout/Navigation';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import Footer from '@/components/layout/Footer';

export default function ContactPage() {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    subject: '',
    message: '',
  });
  const [status, setStatus] = useState<'idle' | 'sending' | 'success' | 'error'>('idle');
  const [errorMessage, setErrorMessage] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setStatus('sending');
    setErrorMessage('');

    // Basic validation
    if (!formData.name || !formData.email || !formData.message) {
      setStatus('error');
      setErrorMessage('Please fill in all required fields.');
      return;
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(formData.email)) {
      setStatus('error');
      setErrorMessage('Please enter a valid email address.');
      return;
    }

    // Simulate form submission (in production, this would send to a backend endpoint)
    try {
      // TODO: Replace with actual API call to backend
      await new Promise((resolve) => setTimeout(resolve, 1000));
      
      setStatus('success');
      setFormData({ name: '', email: '', subject: '', message: '' });
      
      // Reset success message after 5 seconds
      setTimeout(() => {
        setStatus('idle');
      }, 5000);
    } catch {
      setStatus('error');
      setErrorMessage('Failed to send message. Please try again later.');
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navigation variant="full" />
      <main className="flex-1 max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 w-full">
        <div className="bg-card/80 backdrop-blur-sm border border-border/50 rounded-lg p-6 sm:p-8 lg:p-12">
          <h1 className="text-3xl sm:text-4xl font-bold text-foreground/95 mb-4">Contact Us</h1>
          <p className="text-muted-foreground mb-8 text-sm sm:text-base">
            Have a question or need support? We&apos;d love to hear from you. Fill out the form below and we&apos;ll get back to you as soon as possible.
          </p>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Contact Information */}
            <div className="lg:col-span-1">
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-semibold text-foreground/95 mb-2">Get in Touch</h3>
                  <p className="text-muted-foreground text-sm">
                    We typically respond within 24-48 hours during business days.
                  </p>
                </div>

                <div>
                  <h4 className="text-sm font-semibold text-foreground/95 mb-1">Support</h4>
                  <p className="text-muted-foreground text-sm">
                    For technical support, account issues, or billing questions.
                  </p>
                </div>

                <div>
                  <h4 className="text-sm font-semibold text-foreground/95 mb-1">General Inquiries</h4>
                  <p className="text-muted-foreground text-sm">
                    For partnership opportunities, media inquiries, or general questions.
                  </p>
                </div>

                <div>
                  <h4 className="text-sm font-semibold text-foreground/95 mb-1">Legal</h4>
                  <p className="text-muted-foreground text-sm">
                    For privacy concerns, data requests, or legal matters, please refer to our{' '}
                    <Link href="/privacy" className="text-primary hover:underline">
                      Privacy Policy
                    </Link>.
                  </p>
                </div>
              </div>
            </div>

            {/* Contact Form */}
            <div className="lg:col-span-2">
              <form onSubmit={handleSubmit} className="space-y-6">
                {status === 'success' && (
                  <div className="bg-primary/10 border border-primary/20 rounded-lg p-4">
                    <p className="text-primary text-sm">
                      Thank you for your message! We&apos;ll get back to you soon.
                    </p>
                  </div>
                )}

                {status === 'error' && (
                  <div className="bg-destructive/10 border border-destructive/20 rounded-lg p-4">
                    <p className="text-destructive text-sm">{errorMessage || 'An error occurred. Please try again.'}</p>
                  </div>
                )}

                <div>
                  <label htmlFor="name" className="block text-sm font-medium text-foreground mb-1">
                    Name <span className="text-destructive">*</span>
                  </label>
                  <Input
                    type="text"
                    id="name"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    className="bg-background"
                    placeholder="Your name"
                  />
                </div>

                <div>
                  <label htmlFor="email" className="block text-sm font-medium text-foreground mb-1">
                    Email <span className="text-destructive">*</span>
                  </label>
                  <Input
                    type="email"
                    id="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    required
                    className="bg-background"
                    placeholder="your.email@example.com"
                  />
                </div>

                <div>
                  <label htmlFor="subject" className="block text-sm font-medium text-foreground mb-1">
                    Subject
                  </label>
                  <select
                    id="subject"
                    name="subject"
                    value={formData.subject}
                    onChange={handleChange}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  >
                    <option value="">Select a subject</option>
                    <option value="support">Technical Support</option>
                    <option value="billing">Billing Question</option>
                    <option value="account">Account Issue</option>
                    <option value="partnership">Partnership Inquiry</option>
                    <option value="media">Media Inquiry</option>
                    <option value="other">Other</option>
                  </select>
                </div>

                <div>
                  <label htmlFor="message" className="block text-sm font-medium text-foreground mb-1">
                    Message <span className="text-destructive">*</span>
                  </label>
                  <textarea
                    id="message"
                    name="message"
                    value={formData.message}
                    onChange={handleChange}
                    required
                    rows={6}
                    className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 resize-y"
                    placeholder="Please describe your question or issue..."
                  />
                </div>

                <div>
                  <Button
                    type="submit"
                    disabled={status === 'sending'}
                    className="w-full sm:w-auto bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground font-semibold"
                  >
                    {status === 'sending' ? 'Sending...' : 'Send Message'}
                  </Button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}

