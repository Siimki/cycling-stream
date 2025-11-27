"use client"

import { useState, useRef, useEffect } from "react"
import { Send, Settings, Users } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

interface ChatMessage {
  id: string
  user: string
  message: string
  color: string
  badge?: "mod" | "vip" | "sub" | "founder"
}

const userColors = ["#22d3ee", "#a78bfa", "#fb7185", "#4ade80", "#fbbf24", "#f472b6", "#60a5fa", "#34d399"]

const initialMessages: ChatMessage[] = [
  {
    id: "1",
    user: "VeloMaster",
    message: "Pogačar is absolutely flying today",
    color: userColors[0],
    badge: "founder",
  },
  { id: "2", user: "TourFan2024", message: "3 min gap already!", color: userColors[1] },
  { id: "3", user: "AlpineRider", message: "This gradient is insane", color: userColors[2], badge: "sub" },
  { id: "4", user: "ProCyclist", message: "Vingegaard needs to respond soon", color: userColors[3], badge: "mod" },
  { id: "5", user: "MtnKing", message: "Galibier never disappoints", color: userColors[4] },
  { id: "6", user: "SprintFan", message: "who's taking the stage?", color: userColors[5], badge: "vip" },
  { id: "7", user: "PelotonPro", message: "UAE dominating the mountains", color: userColors[6], badge: "sub" },
  { id: "8", user: "CyclingNerd", message: "what an attack!!", color: userColors[7] },
]

const randomMessages = [
  "insane pace",
  "lets gooo",
  "what a climb",
  "GC battle heating up",
  "amazing scenery",
  "that attack tho",
  "who wins today?",
  "brutal gradient",
  "peloton is suffering",
  "this is incredible",
  "champion move",
  "legend",
  "pogacar insane",
  "great stage",
]

const randomUsers = [
  "RoadWarrior",
  "VeloFan",
  "ClimbKing",
  "SprintPro",
  "TourLover",
  "MtnClimber",
  "StageHunter",
  "GCRider",
  "PelotonFan",
  "CycleNerd",
]

export function LiveChat() {
  const [messages, setMessages] = useState<ChatMessage[]>(initialMessages)
  const [newMessage, setNewMessage] = useState("")
  const scrollRef = useRef<HTMLDivElement>(null)
  const [viewerCount] = useState(847)

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [messages])

  useEffect(() => {
    const interval = setInterval(() => {
      const msg = randomMessages[Math.floor(Math.random() * randomMessages.length)]
      const user = randomUsers[Math.floor(Math.random() * randomUsers.length)]
      const color = userColors[Math.floor(Math.random() * userColors.length)]
      const badges: (ChatMessage["badge"] | undefined)[] = [undefined, undefined, undefined, "sub", "vip"]
      const badge = badges[Math.floor(Math.random() * badges.length)]

      setMessages((prev) => [...prev.slice(-100), { id: Date.now().toString(), user, message: msg, color, badge }])
    }, 2500)

    return () => clearInterval(interval)
  }, [])

  const handleSend = () => {
    if (!newMessage.trim()) return
    setMessages((prev) => [
      ...prev,
      { id: Date.now().toString(), user: "You", message: newMessage, color: "#22d3ee", badge: "sub" },
    ])
    setNewMessage("")
  }

  const getBadge = (badge?: ChatMessage["badge"]) => {
    const baseClass = "px-1 py-px text-[9px] font-semibold uppercase rounded mr-1"
    switch (badge) {
      case "mod":
        return <span className={`${baseClass} bg-green-500/20 text-green-400`}>mod</span>
      case "vip":
        return <span className={`${baseClass} bg-pink-500/20 text-pink-400`}>vip</span>
      case "sub":
        return <span className={`${baseClass} bg-primary/20 text-primary`}>sub</span>
      case "founder":
        return <span className={`${baseClass} bg-amber-500/20 text-amber-400`}>og</span>
      default:
        return null
    }
  }

  return (
    <div className="flex flex-col h-full bg-card/50">
      {/* Header */}
      <div className="px-3 py-2 border-b border-border/50 flex items-center justify-between shrink-0">
        <div className="flex items-center gap-2">
          <span className="text-xs font-medium text-foreground">Stream Chat</span>
        </div>
        <div className="flex items-center gap-1">
          <div className="flex items-center gap-1 text-muted-foreground mr-1">
            <Users className="w-3 h-3" />
            <span className="text-[10px] tabular-nums">{viewerCount}</span>
          </div>
          <Button variant="ghost" size="icon" className="h-6 w-6 text-muted-foreground hover:text-foreground">
            <Settings className="w-3 h-3" />
          </Button>
        </div>
      </div>

      {/* Messages */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto chat-scroll px-3 py-2">
        <div className="space-y-1">
          {messages.map((msg) => (
            <div key={msg.id} className="leading-tight py-0.5 hover:bg-muted/30 -mx-1 px-1 rounded">
              {getBadge(msg.badge)}
              <span className="text-[11px] font-medium cursor-pointer hover:underline" style={{ color: msg.color }}>
                {msg.user}
              </span>
              <span className="text-muted-foreground text-[11px]">: </span>
              <span className="text-[11px] text-foreground/90">{msg.message}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Input */}
      <div className="p-2 border-t border-border/50 shrink-0">
        <div className="flex items-center gap-1.5">
          <Input
            placeholder="Send a message"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSend()}
            className="flex-1 h-7 bg-muted/50 border-border/50 text-[11px] placeholder:text-muted-foreground/50 focus-visible:ring-1 focus-visible:ring-primary/50"
          />
          <Button size="icon" className="h-7 w-7 bg-primary/90 hover:bg-primary shrink-0" onClick={handleSend}>
            <Send className="w-3 h-3" />
          </Button>
        </div>
        <p className="text-[9px] text-muted-foreground/60 mt-1.5 px-0.5">Send a message to earn +5 bonus points</p>
      </div>
    </div>
  )
}
