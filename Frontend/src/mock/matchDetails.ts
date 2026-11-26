import type { MatchExtraMock } from '@/types'

// Mock ข้อมูลเหตุการณ์ในแมตช์และสถิติ (⚠️ ยังไม่มี /events และ /statistics API ใน backend)
export const mockMatchExtra: MatchExtraMock = {
  ht_home_score: 1,
  ht_away_score: 0,
  events: [
    { id: 'e1', minute: 18, type: 'goal', team: 'home', player_name: 'Guilherme Bissoli', assist_name: 'Theerathon Bunmathan' },
    { id: 'e2', minute: 34, type: 'yellow_card', team: 'away', player_name: 'Thitipan Puangjan' },
    { id: 'e3', minute: 52, type: 'goal', team: 'away', player_name: 'Muhsen Al-Ghassani' },
    { id: 'e4', minute: 64, type: 'yellow_card', team: 'home', player_name: 'Pansa Hemviboon' },
    { id: 'e5', minute: 76, type: 'red_card', team: 'away', player_name: 'Victor Cardozo' },
    { id: 'e6', minute: 83, type: 'goal', team: 'home', player_name: 'Supachai Chaided', assist_name: 'Goran Čaušić' },
  ],
  stats: [
    { label: 'การครองบอล', home: 58, away: 42, unit: '%' },
    { label: 'โอกาสยิงทั้งหมด', home: 14, away: 9 },
    { label: 'ยิงเข้ากรอบ', home: 6, away: 3 },
    { label: 'ลูกเตะมุม', home: 7, away: 4 },
    { label: 'ทำฟาวล์', home: 11, away: 15 },
    { label: 'ใบเหลือง', home: 1, away: 2 },
    { label: 'ใบแดง', home: 0, away: 1 },
  ],
}
