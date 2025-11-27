"use client"

import { useState, useEffect, useRef, useCallback } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Zap, Search, Bell, User, ChevronDown, Menu, LogOut, Settings } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/contexts/AuthContext"

interface NavigationProps {
  variant?: "full" | "minimal"
}

export function Navigation({ variant = "full" }: NavigationProps) {
  const router = useRouter()
  const { user, isAuthenticated, isLoading: loadingUser, logout } = useAuth()
  const [showUserMenu, setShowUserMenu] = useState(false)
  const userMenuRef = useRef<HTMLDivElement>(null)

  // Close user menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setShowUserMenu(false)
      }
    }

    if (showUserMenu) {
      document.addEventListener('mousedown', handleClickOutside)
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [showUserMenu])

  const handleLogout = useCallback(() => {
    logout()
    setShowUserMenu(false)
    router.push('/')
  }, [logout, router])

  return (
    <header className="sticky top-0 z-50 bg-card/95 backdrop-blur-xl border-b border-border/30 shrink-0">
      <div className="px-4 sm:px-5 lg:px-6">
        <div className="flex items-center justify-between h-14 sm:h-16">
          {/* Logo */}
          <div className="flex items-center gap-6 sm:gap-8 lg:gap-10">
            <Link href="/" className="flex items-center gap-2.5 hover:opacity-90 transition-opacity">
              <div className="w-8 h-8 sm:w-9 sm:h-9 rounded-lg bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center shadow-lg shadow-primary/20">
                <Zap className="w-4 h-4 sm:w-5 sm:h-5 text-primary-foreground" />
              </div>
              <span className="text-lg sm:text-xl font-bold tracking-tight">
                Peloton<span className="text-primary">Live</span>
              </span>
            </Link>

            {/* Navigation - Only shown in full variant, hidden on mobile */}
            {variant === "full" && (
              <nav className="hidden md:flex items-center gap-1">
                <Link href="/">
                  <Button variant="ghost" className="text-foreground text-sm font-semibold h-9 px-4 rounded-lg bg-muted/50">
                    Live
                    <span className="ml-2 w-2 h-2 rounded-full bg-red-500 animate-live-pulse" />
                  </Button>
                </Link>
                <Link href="/">
                  <Button
                    variant="ghost"
                    className="text-muted-foreground hover:text-foreground hover:bg-muted/50 text-sm font-medium h-9 px-4 rounded-lg transition-colors"
                  >
                    Races
                  </Button>
                </Link>
                <Button
                  variant="ghost"
                  className="text-muted-foreground hover:text-foreground hover:bg-muted/50 text-sm font-medium h-9 px-4 rounded-lg transition-colors"
                >
                  Replays
                </Button>
              </nav>
            )}
          </div>

          {/* Right Side */}
          <div className="flex items-center gap-1.5 sm:gap-2">
            {/* Mobile menu button - Only in full variant */}
            {variant === "full" && (
              <Button variant="ghost" size="icon" className="md:hidden text-muted-foreground hover:text-foreground hover:bg-muted/50 w-9 h-9 rounded-lg">
                <Menu className="w-5 h-5" />
              </Button>
            )}

            {variant === "full" && (
              <>
                <Button variant="ghost" size="icon" className="hidden sm:flex text-muted-foreground hover:text-foreground hover:bg-muted/50 w-9 h-9 rounded-lg transition-colors">
                  <Search className="w-5 h-5" />
                </Button>
                {isAuthenticated && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-muted-foreground hover:text-foreground hover:bg-muted/50 w-9 h-9 rounded-lg transition-colors"
                  >
                    <Bell className="w-5 h-5" />
                  </Button>
                )}
              </>
            )}

            {/* Auth buttons or User menu */}
            {isAuthenticated ? (
              <div className="relative ml-1 sm:ml-2" ref={userMenuRef}>
                <button
                  onClick={() => setShowUserMenu(!showUserMenu)}
                  className="flex items-center gap-1.5 cursor-pointer group"
                  aria-label="User menu"
                  aria-expanded={showUserMenu}
                >
                  <div className="w-8 h-8 sm:w-9 sm:h-9 rounded-full bg-gradient-to-br from-primary/70 to-primary flex items-center justify-center ring-2 ring-border/50 group-hover:ring-primary/50 transition-all">
                    <User className="w-4 h-4 sm:w-5 sm:h-5 text-primary-foreground" />
                  </div>
                  <ChevronDown className="hidden sm:block w-4 h-4 text-muted-foreground group-hover:text-foreground transition-colors" />
                </button>

                {/* User dropdown menu */}
                {showUserMenu && (
                  <div className="absolute right-0 mt-2 w-56 bg-card border border-border/50 rounded-lg shadow-lg z-50">
                    {/* User info header */}
                    {loadingUser ? (
                      <div className="px-4 py-3 border-b border-border/50">
                        <div className="text-sm text-muted-foreground">Loading...</div>
                      </div>
                    ) : user ? (
                      <div className="px-4 py-3 border-b border-border/50">
                        <div className="text-sm font-semibold text-foreground truncate">
                          {user.name || user.email}
                        </div>
                        <div className="text-xs text-muted-foreground mt-1">
                          {(user.points ?? 0).toLocaleString()} points
                        </div>
                      </div>
                    ) : null}
                    
                    {/* Menu items */}
                    <div className="py-1">
                      <Link
                        href="/profile"
                        onClick={() => setShowUserMenu(false)}
                        className="block px-4 py-2 text-sm text-foreground hover:bg-muted/50 transition-colors"
                      >
                        Profile
                      </Link>
                      <Link
                        href="/settings"
                        onClick={() => setShowUserMenu(false)}
                        className="block px-4 py-2 text-sm text-foreground hover:bg-muted/50 transition-colors flex items-center gap-2"
                      >
                        <Settings className="w-4 h-4" />
                        Settings
                      </Link>
                      <button
                        onClick={handleLogout}
                        className="w-full text-left px-4 py-2 text-sm text-foreground hover:bg-muted/50 transition-colors flex items-center gap-2"
                      >
                        <LogOut className="w-4 h-4" />
                        Logout
                      </button>
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="flex items-center gap-2 ml-1 sm:ml-2">
                <Link href="/auth/login">
                  <Button
                    variant="ghost"
                    className="text-sm font-medium h-9 px-4 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted/50"
                  >
                    Login
                  </Button>
                </Link>
                <Link href="/auth/register">
                  <Button
                    className="text-sm font-medium h-9 px-4 rounded-lg bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground"
                  >
                    Sign Up
                  </Button>
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}
