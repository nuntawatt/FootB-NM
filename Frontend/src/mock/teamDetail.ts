import type { TeamDetail } from '@/types'

// Mock data ข้อมูลทีม (⚠️ ยังไม่มี coach/form/squad API ใน backend)
export const mockTeamDetail: TeamDetail = {
  coach_name: 'Raúl Caneda',
  coach_image: 'https://media.api-sports.io/football/coachs/19234.png',
  form: [
    { score: '1-1', opponent_logo: 'https://media.api-sports.io/football/teams/2279.png', result: 'draw' },
    { score: '0-0', opponent_logo: 'https://media.api-sports.io/football/teams/2282.png', result: 'draw' },
    { score: '1-2', opponent_logo: 'https://media.api-sports.io/football/teams/2284.png', result: 'loss' },
    { score: '1-0', opponent_logo: 'https://media.api-sports.io/football/teams/2288.png', result: 'win' },
    { score: '3-1', opponent_logo: 'https://media.api-sports.io/football/teams/2287.png', result: 'win' },
  ],
  squad: [
    { position: 'GK', number: 13, name: 'Neil Etheridge', rating: 6 },
    { position: 'RDC', number: 39, name: 'Jakkaphan Praisuwan', rating: 7 },
    { position: 'DC', number: 16, name: 'Kenneth Dougall', rating: 6 },
    { position: 'LDC', number: 3, name: 'Pansa Hemviboon', rating: 6 },
    { position: 'RWB', number: 19, name: 'Kingsley Schindler', rating: 7 },
    { position: 'CM', number: 23, name: 'Goran Čaušić', rating: 7 },
    { position: 'CM', number: 37, name: 'Darlan', rating: 6 },
    { position: 'CM', number: 55, name: 'Theerathon Bunmathan', rating: 7 },
    { position: 'LWB', number: 2, name: 'Sasalak Haiprakhon', rating: 7 },
    { position: 'ST', number: 30, name: 'Stefan Schimmer', rating: 7 },
    { position: 'ST', number: 7, name: 'Bissoli', rating: 8 },
  ],
}
