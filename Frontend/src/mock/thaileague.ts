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
  content: string
  image: string
  category: string
  date: string
  author: string
  redirect_url?: string
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

// Mock ข่าวสารไทยลีก (⚠️ ยังไม่มี /news API ใน backend)
export const mockNews: NewsItem[] = [
  {
    id: 1,
    title: 'ไทยลีก 1 เตรียมใช้ VAR เต็มรูปแบบทุกสนามในเลกสอง',
    summary: 'สมาคมกีฬาฟุตบอลฯ ยืนยันความพร้อมระบบ VAR เพื่อยกระดับมาตรฐานการตัดสินฟุตบอลไทยลีก',
    content: 'สมาคมกีฬาฟุตบอลแห่งประเทศไทยฯ ร่วมกับบริษัท ไทยลีก จำกัด เตรียมนำเทคโนโลยี VAR (Video Assistant Referee) มาใช้งานเต็มรูปแบบครบทุกสนามแข่งขันในศึกไทยลีก 1 เลกที่สอง หลังผ่านการอบรมผู้ตัดสินและทดสอบระบบในสนามเหย้าของทุกสโมสรเรียบร้อยแล้ว เพื่อความโปร่งใสและยกระดับมาตรฐานการแข่งขันให้ทัดเทียมระดับเอเชีย',
    image: 'https://images.unsplash.com/photo-1517466787929-bc90951d0974?w=800&h=480&fit=crop',
    category: 'Thai League 1',
    date: '26 พ.ย. 2026',
    author: 'กองบรรณาธิการ FOOTB',
  },
  {
    id: 2,
    title: 'บุรีรัมย์ฯ ปะทะ แบงค็อกฯ บิ๊กแมตช์ชี้ชะตาจ่าฝูงสุดสัปดาห์นี้',
    summary: 'วิเคราะห์ความพร้อมก่อนเกมคู่เอก ไทยลีก 1 พร้อมสถิติการพบกัน 5 นัดหลังสุด',
    content: 'ศึกบิ๊กแมตช์ประจำสัปดาห์ที่แฟนบอลไทยรอคอย บุรีรัมย์ ยูไนเต็ด เตรียมเปิดบ้านต้อนรับการมาเยือนของ ทรู แบงค็อก ยูไนเต็ด โดยทั้งสองทีมต่างกำลังอยู่ในฟอร์มยอดเยี่ยมและมีคะแนนเบียดแย่งจ่าฝูงกันอย่างสูสี คาดว่าทั้งคู่จะจัดผู้เล่นชุดใหญ่ลงสนามเต็มอัตราศึก',
    image: 'https://images.unsplash.com/photo-1431324155629-1a6deb1dec8d?w=800&h=480&fit=crop',
    category: 'พรีวิว',
    date: '25 พ.ย. 2026',
    author: 'ทีมข่าวไทยลีก',
  },
  {
    id: 3,
    title: 'ช้าง เอฟเอ คัพ รอบ 16 ทีมสุดท้าย เปิดโผประกบคู่สุดเดือด',
    summary: 'ทีมยักษ์ใหญ่จากไทยลีก 1 โคจรมาพบกันหลายคู่ เตรียมฟาดแข้งเดือนธันวาคมนี้',
    content: 'ผลการจับสลากประกบคู่ฟุตบอลถ้วย ช้าง เอฟเอ คัพ รอบ 16 ทีมสุดท้าย ออกมาเป็นที่เรียบร้อย โดยมีหลายคู่บิ๊กแมตช์ที่ทีมจากไทยลีก 1 ต้องโคจรมาพบกันเอง ขณะที่ทีมม้ามืดจากไทยลีก 2 ก็พร้อมสร้างเซอร์ไพรส์ในบ้านตัวเอง',
    image: 'https://images.unsplash.com/photo-1579952363873-27f3bade9f55?w=800&h=480&fit=crop',
    category: 'FA CUP',
    date: '24 พ.ย. 2026',
    author: 'กองบรรณาธิการ FOOTB',
  },
  {
    id: 4,
    title: 'ไทยลีก 2 ลุ้นตั๋วเลื่อนชั้นสุดเข้มข้น 4 ทีมหัวตารางแต้มห่างแค่ 2 คะแนน',
    summary: 'สรุปสถานการณ์ลีกรองหลังผ่านครึ่งทางแรกของฤดูกาล 2026/2027',
    content: 'การแข่งขันฟุตบอลไทยลีก 2 ฤดูกาลนี้ดุเดือดกว่าทุกปี เมื่อ 4 ทีมหัวตารางมีคะแนนไล่เลี่ยกันห่างเพียง 2 แต้ม ทำให้ทุกนัดต่อจากนี้มีความหมายต่อการลุ้นโควตาเลื่อนชั้นสู่ลีกสูงสุดโดยตรง',
    image: 'https://images.unsplash.com/photo-1508098682722-e99c43a406b2?w=800&h=480&fit=crop',
    category: 'Thai League 2',
    date: '23 พ.ย. 2026',
    author: 'ทีมข่าวไทยลีก',
  },
]
