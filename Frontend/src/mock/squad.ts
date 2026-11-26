import type { SquadPlayer } from '@/types'

// Mock data squad/lineup (⚠️ ยังไม่มี squad API ใน backend)
export const mockSquad: SquadPlayer[] = [
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
]

export const mockLineup = {
  home: mockSquad,
  away: [
    { position: 'GK', number: 1, name: 'Kampon Phanuphong', rating: 6 },
    { position: 'DC', number: 4, name: 'Adison Promruk', rating: 6 },
    { position: 'DC', number: 5, name: 'Victor Cardozo', rating: 6 },
    { position: 'DC', number: 22, name: 'Wattana Sakha', rating: 5 },
    { position: 'RM', number: 7, name: 'Thitipan Puangjan', rating: 7 },
    { position: 'CM', number: 8, name: 'Sarach Yooyen', rating: 6 },
    { position: 'CM', number: 10, name: 'Chanathip Songkrasin', rating: 7 },
    { position: 'LM', number: 11, name: 'Suphanat Mueanta', rating: 7 },
    { position: 'RW', number: 17, name: 'Ekanit Panya', rating: 6 },
    { position: 'LW', number: 20, name: 'Sumanya Purisai', rating: 6 },
    { position: 'ST', number: 9, name: 'Teerasil Dangda', rating: 7 },
  ],
}
