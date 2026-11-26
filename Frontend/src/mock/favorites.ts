import type { FavoriteTeamCard } from '@/types'

// Mock data ทีมโปรด (⚠️ ยังไม่มี user auth/favorite API ใน backend)
export const mockFavorites: FavoriteTeamCard[] = [
  {
    team: {
      id: '1',
      league_id: '',
      name: 'Buriram United',
      short_name: 'BUR',
      logo_url: 'https://media.api-sports.io/football/teams/2275.png',
      created_at: '',
      updated_at: '',
    },
    venue: 'Buriram Stadium',
    next_match_date: 'เสาร์ 28 พฤศจิกายน, 19:00',
    broadcast: 'MONOMAX / AIS PLAY',
  },
  {
    team: {
      id: '2',
      league_id: '',
      name: 'Chiangrai United',
      short_name: 'CRU',
      logo_url: 'https://media.api-sports.io/football/teams/2278.png',
      created_at: '',
      updated_at: '',
    },
    venue: 'Singha Stadium',
    next_match_date: 'เสาร์ 28 พฤศจิกายน, 18:00',
    broadcast: 'MONOMAX / AIS PLAY',
  },
  {
    team: {
      id: '3',
      league_id: '',
      name: 'Bangkok United',
      short_name: 'BKK',
      logo_url: 'https://media.api-sports.io/football/teams/2276.png',
      created_at: '',
      updated_at: '',
    },
    venue: 'Thammasat Stadium',
    next_match_date: 'อาทิตย์ 29 พฤศจิกายน, 18:00',
    broadcast: 'MONOMAX / AIS PLAY',
  },
  {
    team: {
      id: '4',
      league_id: '',
      name: 'Port FC',
      short_name: 'PRT',
      logo_url: 'https://media.api-sports.io/football/teams/2282.png',
      created_at: '',
      updated_at: '',
    },
    venue: 'PAT Stadium',
    next_match_date: 'เสาร์ 28 พฤศจิกายน, 19:00',
    broadcast: 'MONOMAX / AIS PLAY',
  },
  {
    team: {
      id: '5',
      league_id: '',
      name: 'Chonburi FC',
      short_name: 'CHB',
      logo_url: 'https://media.api-sports.io/football/teams/2279.png',
      created_at: '',
      updated_at: '',
    },
    venue: 'Chonburi Stadium',
    next_match_date: 'เสาร์ 28 พฤศจิกายน, 19:30',
    broadcast: 'MONOMAX / AIS PLAY',
  },
  {
    team: {
      id: '6',
      league_id: '',
      name: 'BG Pathum United',
      short_name: 'BGP',
      logo_url: 'https://media.api-sports.io/football/teams/2277.png',
      created_at: '',
      updated_at: '',
    },
    venue: 'Leo Stadium',
    next_match_date: 'อาทิตย์ 29 พฤศจิกายน, 18:00',
    broadcast: 'MONOMAX / AIS PLAY',
  },
]
