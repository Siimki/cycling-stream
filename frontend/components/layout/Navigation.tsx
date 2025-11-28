"use client"

import { useState, useEffect, useRef, useCallback } from "react"
import Link from "next/link"
import { useRouter, usePathname } from "next/navigation"
import { Zap, Search, Bell, User, ChevronDown, Menu, LogOut, Settings } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/contexts/AuthContext"
import UserStatusBar from "@/components/user/UserStatusBar"

interface NavigationProps {
  variant?: "full" | "minimal"
}

export function Navigation({ variant = "full" }: NavigationProps) {
  const router = useRouter()
  const pathname = usePathname()
  const { user, isAuthenticated, isLoading: loadingUser, logout } = useAuth()
  const [showUserMenu, setShowUserMenu] = useState(false)
  const userMenuRef = useRef<HTMLDivElement>(null)
  
  // Helper to check if a route is active
  const isActive = (path: string) => {
    if (path === '/') {
      return pathname === '/'
    }
    return pathname?.startsWith(path)
  }

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
    <header className="sticky top-0 z-50 bg-card/95 backdrop-blur-xl border-b border-[#0D0D0D] shrink-0 shadow-sm">
      {/* Use same horizontal padding as video container for perfect alignment */}
      <div className="px-6 lg:px-8 pt-5 pb-4">
        <div className="flex items-center justify-between">
          {/* Left: Brand - Equal padding zone */}
          <div className="flex items-center flex-1">
            <Link href="/" className="flex items-center gap-2.5 hover:opacity-90 transition-opacity">
              <div className="w-9 h-9 sm:w-10 sm:h-10 rounded-lg bg-gradient-to-br from-primary to-primary/80 flex items-center justify-center shadow-lg shadow-primary/20">
                <Zap className="w-5 h-5 sm:w-6 sm:h-6 text-primary-foreground" />
              </div>
              <span className="text-xl sm:text-2xl font-bold tracking-tight">
                Peloton<span className="text-primary">Live</span>
              </span>
            </Link>
          </div>

          {/* Center: Primary Navigation - Centered */}
          {variant === "full" && (
            <nav className="hidden md:flex items-center gap-1 flex-shrink-0">
                <Link href="/">
                  <Button 
                    variant="ghost" 
                    className={`relative text-xl font-semibold h-12 px-5 py-3 transition-colors ${
                      isActive('/') 
                        ? 'text-foreground' 
                        : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
                    }`}
                  >
                    Live
                    <span className="ml-2 w-2.5 h-2.5 rounded-full bg-primary animate-live-pulse" />
                    {isActive('/') && (
                      <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full" />
                    )}
                  </Button>
                </Link>
                {isAuthenticated && (
                  <Link href="/for-you">
                    <Button
                      variant="ghost"
                      className={`relative text-xl font-semibold h-12 px-5 py-3 transition-colors ${
                        isActive('/for-you')
                          ? 'text-foreground'
                          : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
                      }`}
                    >
                      For You
                      {isActive('/for-you') && (
                        <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full" />
                      )}
                    </Button>
                  </Link>
                )}
                <Link href="/races">
                  <Button
                    variant="ghost"
                    className={`relative text-xl font-semibold h-12 px-5 py-3 transition-colors ${
                      isActive('/races')
                        ? 'text-foreground'
                        : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
                    }`}
                  >
                    Races
                    {isActive('/races') && (
                      <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full" />
                    )}
                  </Button>
                </Link>
                <Button
                  variant="ghost"
                  className="text-xl font-semibold text-muted-foreground hover:text-foreground hover:bg-muted/30 h-12 px-5 py-3 transition-colors"
                >
                  Replays
                </Button>
                <Link href="/leaderboard">
                  <Button
                    variant="ghost"
                    className={`relative text-xl font-semibold h-12 px-5 py-3 transition-colors ${
                      isActive('/leaderboard')
                        ? 'text-foreground'
                        : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
                    }`}
                  >
                    Leaderboard
                    {isActive('/leaderboard') && (
                      <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full" />
                    )}
                  </Button>
                </Link>
                {isAuthenticated && (
                  <Link href="/missions">
                    <Button
                      variant="ghost"
                      className={`relative text-xl font-semibold h-12 px-5 py-3 transition-colors ${
                        isActive('/missions')
                          ? 'text-foreground'
                          : 'text-muted-foreground hover:text-foreground hover:bg-muted/30'
                      }`}
                    >
                      Missions
                      {isActive('/missions') && (
                        <span className="absolute bottom-0 left-0 right-0 h-0.5 bg-primary rounded-full" />
                      )}
                    </Button>
                  </Link>
                )}
              </nav>
          )}

          {/* Right: Utilities + User Status - Equal padding zone */}
          <div className="flex items-center flex-1 justify-end">
            {/* Mobile menu button - Only in full variant */}
            {variant === "full" && (
              <Button variant="ghost" size="icon" className="md:hidden text-muted-foreground/90 hover:text-foreground hover:bg-muted/30 w-10 h-10 rounded-lg">
                <Menu className="w-6 h-6" />
              </Button>
            )}

            {/* Group A: Search + Bell - 16px gap between icons */}
            {variant === "full" && (
              <div className="flex items-center gap-4">
                <Button variant="ghost" size="icon" className="hidden sm:flex text-muted-foreground/90 hover:text-foreground hover:bg-muted/30 w-10 h-10 rounded-lg transition-colors">
                  <Search className="w-6 h-6" />
                </Button>
                {isAuthenticated && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-muted-foreground/90 hover:text-foreground hover:bg-muted/30 w-10 h-10 rounded-lg transition-colors"
                  >
                    <Bell className="w-6 h-6" />
                  </Button>
                )}
              </div>
            )}

            {/* Group B: User Block (Status + Avatar) - 24px gap from Group A, 16px between status and avatar */}
            {isAuthenticated && variant === "full" && (
              <div className="hidden lg:flex items-center gap-4 ml-6">
                {/* User Status Bar (Level, XP, Points, Streak) - Compact single component */}
                <UserStatusBar />
                
                {/* Profile Avatar */}
                <div className="relative" ref={userMenuRef}>
                  <button
                    onClick={() => setShowUserMenu(!showUserMenu)}
                    className="flex items-center gap-1.5 cursor-pointer group"
                    aria-label="User menu"
                    aria-expanded={showUserMenu}
                  >
                    <div className="w-10 h-10 sm:w-11 sm:h-11 rounded-full bg-muted/50 flex items-center justify-center ring-2 ring-border/50 group-hover:ring-border transition-all">
                      <User className="w-5 h-5 sm:w-6 sm:h-6 text-foreground" />
                    </div>
                    <ChevronDown className="hidden sm:block w-5 h-5 text-muted-foreground/90 group-hover:text-foreground transition-colors" />
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
              </div>
            )}

            {/* Auth buttons (when not authenticated) */}
            {!isAuthenticated && (
              <div className="flex items-center gap-2 ml-1 sm:ml-2">
                <Link href="/auth/login">
                  <Button
                    variant="ghost"
                    className="text-base font-medium h-10 px-5 rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted/50"
                  >
                    Login
                  </Button>
                </Link>
                <Link href="/auth/register">
                  <Button
                    className="text-base font-medium h-10 px-5 rounded-lg bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70 text-primary-foreground"
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
