// === API Response Types (ตรงกับ backend models) ===

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

// === Mock Types (ยังไม่มี API) ===

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
  zone?: 'champion' | 'acl' | 'acl2' | 'relegation'
}

export interface SquadPlayer {
  position: string
  number: number
  name: string
  rating?: number
}

export interface TeamDetail {
  coach_name: string
  coach_image: string
  form: FormMatch[]
  squad: SquadPlayer[]
}

export interface FormMatch {
  score: string
  opponent_logo: string
  result: 'win' | 'draw' | 'loss'
}

export interface HeroItem {
  id: number
  title: string
  image: string
  category: string
}

export interface FavoriteTeamCard {
  team: Team
  venue: string
  next_match_date: string
  broadcast: string
}

// === Sidebar League (Hardcode) ===

export interface SidebarLeague {
  id: string
  name: string
  logo_url: string
  type: 'league' | 'cup'
}
