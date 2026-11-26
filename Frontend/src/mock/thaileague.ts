export interface TopScorer {
  rank: number
  name: string
  team: string
  team_logo: string
  goals: number
}

export interface NewsItem {
  id: number
  title: string
  summary: string
  image: string
  category: string
  date: string
}

// Mock ดาวซัลโว (⚠️ ยังไม่มี top scorers API ใน backend)
export const mockTopScorers: TopScorer[] = [
  {
    rank: 1,
    name: 'Guilherme Bissoli',
    team: 'Buriram United',
    team_logo: 'https://media.api-sports.io/football/teams/2275.png',
    goals: 12,
  },
  {
    rank: 2,
    name: 'Muhsen Al-Ghassani',
    team: 'Bangkok United',
    team_logo: 'https://media.api-sports.io/football/teams/2276.png',
    goals: 9,
  },
  {
    rank: 3,
    name: 'Supachai Chaided',
    team: 'Buriram United',
    team_logo: 'https://media.api-sports.io/football/teams/2275.png',
    goals: 8,
  },
  {
    rank: 4,
    name: 'Lonsana Doumbouya',
    team: 'Port FC',
    team_logo: 'https://media.api-sports.io/football/teams/2282.png',
    goals: 7,
  },
  {
    rank: 5,
    name: 'Teerasil Dangda',
    team: 'BG Pathum United',
    team_logo: 'https://media.api-sports.io/football/teams/2277.png',
    goals: 6,
  },
]

// Mock ข่าวสารไทยลีก (⚠️ ยังไม่มี news API ใน backend)
export const mockNews: NewsItem[] = [
  {
    id: 1,
    title: 'ไทยลีก 1 เตรียมใช้ VAR เต็มรูปแบบทุกสนามในเลกสอง',
    summary: 'สมาคมกีฬาฟุตบอลฯ ยืนยันความพร้อมระบบ VAR เพื่อยกระดับมาตรฐานการตัดสินฟุตบอลไทยลีก',
    image: 'https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=600&h=380&fit=crop',
    category: 'Thai League 1',
    date: '26 พ.ย. 2026',
  },
  {
    id: 2,
    title: 'บุรีรัมย์ฯ ปะทะ แบงค็อกฯ บิ๊กแมตช์ชี้ชะตาจ่าฝูงสุดสัปดาห์นี้',
    summary: 'วิเคราะห์ความพร้อมก่อนเกมคู่เอก ไทยลีก 1 พร้อมสถิติการพบกัน 5 นัดหลังสุด',
    image: 'https://images.unsplash.com/photo-1431324155629-1a6deb1dec8d?w=600&h=380&fit=crop',
    category: 'พรีวิว',
    date: '25 พ.ย. 2026',
  },
  {
    id: 3,
    title: 'ช้าง เอฟเอ คัพ รอบ 16 ทีมสุดท้าย เปิดโผประกบคู่สุดเดือด',
    summary: 'ทีมยักษ์ใหญ่จากไทยลีก 1 โคจรมาพบกันหลายคู่ เตรียมฟาดแข้งเดือนธันวาคมนี้',
    image: 'https://images.unsplash.com/photo-1579952363873-27f3bade9f55?w=600&h=380&fit=crop',
    category: 'FA CUP',
    date: '24 พ.ย. 2026',
  },
]
