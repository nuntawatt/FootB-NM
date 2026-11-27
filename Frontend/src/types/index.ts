// === API Response Types (ตรงกับ backend models) ===

export interface User {
  id: string
  username: string
  email: string
  avatar_url: string
  role: 'user' | 'admin'
  created_at: string
}

export interface AuthResponse {
  status: string
  token: string
  user: User
}

export interface League {
  id: string
  name: string
  short_name: string
  country: string
  logo_url: string
  type: 'league' | 'cup'
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Team {
  id: string
  league_id: string
  name: string
  short_name: string
  logo_url: string
  league?: League
  created_at: string
  updated_at: string
}

export interface Match {
  id: string
  league_id: string
  home_team_id: string
  away_team_id: string
  match_date: string
  home_score: number | null
  away_score: number | null
  home_penalty_score: number | null
  away_penalty_score: number | null
  minute: number | null
  round: string
  status: 'scheduled' | 'live' | 'finished' | 'postponed'
  home_team: Team
  away_team: Team
  league: League
  created_at: string
  updated_at: string
}

export interface SeasonItem {
  season: number
  label: string
  matches_count: number
}

// === Standings (GET /standings) ===

export interface ApiStandingsRow {
  rank: number
  team_id: string
  team_name: string
  team_logo: string
  played: number
  won: number
  drawn: number
  lost: number
  goals_for: number
  goals_against: number
  goal_diff: number
  points: number
  form: string[] // ["W", "D", "L"]
}

export interface StandingRow {
  position: number
  team: Team
  played: number
  won: number
  drawn: number
  lost: number
  goals_for: number
  goals_against: number
  goal_difference: number
  points: number
  form?: string[]
  zone?: 'champion' | 'acl' | 'acl2' | 'relegation'
}

// === Team Extra (GET /teams/:id/form, /coach, /squad) ===

export interface ApiFormMatch {
  match_id: string
  date: string
  opponent: string
  score: string
  result: 'W' | 'D' | 'L'
  is_home: boolean
}

export interface Coach {
  id: string
  team_id: string
  name: string
  nationality: string
  photo_url: string
}

export interface Player {
  id: string
  team_id: string
  name: string
  number: number | null
  position: string
  age: number | null
  nationality: string
  photo_url: string
}

// === Match Extra (GET /matches/:id/lineups, /venue) ===

export interface LineupPlayer {
  number: number | null
  name: string
  position: string
  rating?: string
}

export interface LineupTeam {
  team_id: string
  team_name: string
  formation: string
  coach_name: string
  starting: LineupPlayer[]
  substitutes: LineupPlayer[]
}

export interface MatchLineup {
  match_id: string
  home: LineupTeam
  away: LineupTeam
}

export interface VenueInfo {
  match_id: string
  stadium: string
  city: string
  capacity: number
  broadcast: string[]
}

// === Hero Highlights / News (GET /heroes, GET /news) ===

export interface HeroHighlight {
  id: string
  title: string
  subtitle: string
  image_url: string
  tag: string
  match_id?: string
  redirect_url?: string
}

export interface NewsArticle {
  id: string
  title: string
  summary: string
  content: string
  image_url: string
  category: string
  published_at: string
  author?: string
  redirect_url?: string
}

export interface FavoriteTeamCard {
  team: Team
  venue: string
  next_match_date: string
  broadcast: string
}

export interface SidebarLeague {
  id: string
  name: string
  logo_url: string
  type: 'league' | 'cup'
}
