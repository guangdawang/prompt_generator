import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: '提示词模板系统',
  description: '强大的提示词模板生成和管理系统',
  icons: {
    icon: '/favicon.ico',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="zh-CN">
      <body>{children}</body>
    </html>
  )
}
